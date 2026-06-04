<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { getMessages, sendMessage } from '../api/messages'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const messages = ref([])
const loading = ref(false)
const newMessage = ref('')
const messagesContainer = ref(null)

const exchangeId = route.params.exchangeId

onMounted(loadMessages)

async function loadMessages() {
  loading.value = true
  try {
    const res = await getMessages(exchangeId)
    messages.value = res.messages || res || []
    await nextTick()
    scrollToBottom()
  } catch (e) {
    appStore.showToast('加载消息失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleSend() {
  if (!newMessage.value.trim()) return
  try {
    await sendMessage({ exchangeId, content: newMessage.value.trim() })
    newMessage.value = ''
    await loadMessages()
  } catch (e) {
    appStore.showToast('发送失败', 'error')
  }
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function isMyMessage(msg) {
  const senderId = msg.sender?._id || msg.sender?.id || msg.sender
  const userId = userStore.user?._id || userStore.user?.id
  return senderId === userId
}
</script>

<template>
  <div class="container messages-page">
    <h1 class="page-title">消息沟通</h1>

    <div class="chat-container card">
      <div v-if="loading" class="empty-state">加载中...</div>

      <div v-else ref="messagesContainer" class="chat-messages">
        <div v-if="!messages.length" class="empty-state">暂无消息，发送第一条消息吧</div>
        <div
          v-for="msg in messages"
          :key="msg._id || msg.id"
          class="chat-bubble"
          :class="{ mine: isMyMessage(msg), other: !isMyMessage(msg) }"
        >
          <div class="bubble-name">{{ isMyMessage(msg) ? '我' : (msg.sender?.nickname || msg.sender?.username || '对方') }}</div>
          <div class="bubble-content">{{ msg.content }}</div>
          <div class="bubble-time">{{ new Date(msg.createdAt).toLocaleString() }}</div>
        </div>
      </div>

      <div class="chat-input">
        <input
          v-model="newMessage"
          class="form-input"
          placeholder="输入消息..."
          @keyup.enter="handleSend"
        />
        <button class="btn btn-primary" @click="handleSend">发送</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.messages-page {
  max-width: 700px;
  margin: 0 auto;
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 70vh;
  padding: 0;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-bubble {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: var(--radius);
}

.chat-bubble.mine {
  align-self: flex-end;
  background: var(--primary);
  color: #fff;
}

.chat-bubble.other {
  align-self: flex-start;
  background: var(--gray-100);
  color: var(--gray-800);
}

.bubble-name {
  font-size: 11px;
  opacity: 0.7;
  margin-bottom: 4px;
}

.bubble-content {
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.bubble-time {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 4px;
  text-align: right;
}

.chat-input {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid var(--gray-100);
}

.chat-input .form-input {
  flex: 1;
}
</style>

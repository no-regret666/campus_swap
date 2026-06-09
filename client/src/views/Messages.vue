<script setup>
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMessages, sendMessage, markMessagesRead, deleteMessage, uploadImage } from '../api/messages'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import BackButton from '../components/BackButton.vue'

const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const messages = ref([])
const loading = ref(false)
const newMessage = ref('')
const messagesContainer = ref(null)
const fileInput = ref(null)
const uploading = ref(false)
const selectedImages = ref([]) // 存储选中的图片预览
const previewImage = ref(null) // 当前预览的图片

const exchangeId = route.params.exchangeId

let pollTimer = null
let notificationPermission = 'default'

onMounted(() => {
  requestNotificationPermission()
  loadMessages()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})

function startPolling() {
  pollTimer = setInterval(() => {
    loadMessages(true) // silent poll for new messages
  }, 5000)
}

async function requestNotificationPermission() {
  if ('Notification' in window) {
    if (Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      notificationPermission = permission
    } else {
      notificationPermission = Notification.permission
    }
  }
}

function showBrowserNotification(title, body, icon) {
  if (notificationPermission === 'granted') {
    try {
      new Notification(title, {
        body: body,
        icon: icon || '/favicon.ico',
        badge: '/favicon.ico',
        tag: 'message-notification',
        renotify: true
      })
    } catch (e) {
      console.log('Browser notification not supported')
    }
  }
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

async function loadMessages(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await getMessages(exchangeId)
    const oldCount = messages.value.length
    messages.value = res.messages || res || []
    await nextTick()
    scrollToBottom()

    // 自动标记已读
    if (!silent && messages.value.length > 0) {
      await markMessagesRead(exchangeId).catch(() => {})
    }

    // 新消息提醒
    if (!silent && messages.value.length > oldCount) {
      const hasNewFromOther = messages.value.slice(oldCount).some(msg => !isMyMessage(msg))
      if (hasNewFromOther) {
        playNotificationSound()
        // 发送浏览器通知
        const senderName = messages.value[messages.value.length - 1]?.sender?.nickname || '对方'
        const preview = messages.value[messages.value.length - 1]?.content || '[图片]'
        showBrowserNotification('新消息', `${senderName}: ${preview.substring(0, 50)}`)
      }
    }
  } catch (e) {
    if (!silent) appStore.showToast('加载消息失败', 'error')
  } finally {
    if (!silent) loading.value = false
  }
}

function playNotificationSound() {
  try {
    const audioContext = new (window.AudioContext || window.webkitAudioContext)()
    const oscillator = audioContext.createOscillator()
    const gainNode = audioContext.createGain()
    oscillator.connect(gainNode)
    gainNode.connect(audioContext.destination)
    oscillator.frequency.value = 800
    oscillator.type = 'sine'
    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3)
    oscillator.start(audioContext.currentTime)
    oscillator.stop(audioContext.currentTime + 0.3)
  } catch (e) {
    console.log('Audio notification not supported')
  }
}

async function handleSend() {
  const hasText = newMessage.value.trim()
  const hasImages = selectedImages.value.length > 0

  if (!hasText && !hasImages) return

  // 如果有图片，先上传所有图片
  if (hasImages) {
    uploading.value = true
    try {
      const imageUrls = []
      for (const image of selectedImages.value) {
        if (image.isLocal) {
          const res = await uploadImage(image.file)
          imageUrls.push(res.url)
        }
      }
      // 合并为一条消息发送
      await sendMessage({
        exchangeId,
        content: hasText ? newMessage.value.trim() : '',
        images: imageUrls
      })
      selectedImages.value = [] // 清空已选图片
      newMessage.value = ''
      await loadMessages(true)
      appStore.showToast('发送成功', 'success')
    } catch (e) {
      appStore.showToast('发送失败', 'error')
    } finally {
      uploading.value = false
    }
    return
  }

  // 只有文字
  try {
    await sendMessage({ exchangeId, content: newMessage.value.trim() })
    newMessage.value = ''
    await loadMessages(true)
  } catch (e) {
    appStore.showToast('发送失败', 'error')
  }
}

async function handleRecall(msg) {
  try {
    await deleteMessage(msg.id || msg._id)
    await loadMessages(true)
    appStore.showToast('消息已撤回', 'success')
  } catch (e) {
    appStore.showToast('撤回失败', 'error')
  }
}

// 处理图片选择（支持多选）
function triggerFileInput() {
  fileInput.value?.click()
}

function handleFileChange(event) {
  const files = Array.from(event.target.files)
  if (!files.length) return

  for (const file of files) {
    // 验证文件类型和大小
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
    if (!allowedTypes.includes(file.type)) {
      appStore.showToast(`${file.name} 格式不支持`, 'error')
      continue
    }
    if (file.size > 5 * 1024 * 1024) {
      appStore.showToast(`${file.name} 超过5MB限制`, 'error')
      continue
    }

    // 生成本地预览 URL
    const previewUrl = URL.createObjectURL(file)
    selectedImages.value.push({
      id: Date.now() + Math.random(),
      file,
      previewUrl,
      isLocal: true
    })
  }

  event.target.value = '' // 清空 input 以便下次选择相同文件
}

// 删除选中的图片
function removeImage(image) {
  URL.revokeObjectURL(image.previewUrl) // 释放预览 URL
  selectedImages.value = selectedImages.value.filter(img => img.id !== image.id)
}

// 打开图片预览
function openImagePreview(imgSrc) {
  previewImage.value = imgSrc
}

// 关闭图片预览
function closeImagePreview() {
  previewImage.value = null
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function isMyMessage(msg) {
  const senderId = msg.senderId || msg.sender?._id || msg.sender?.id || msg.sender
  const userId = userStore.user?._id || userStore.user?.id
  return senderId === userId
}

// 格式化时间显示
function formatTime(timestamp) {
  const date = new Date(timestamp)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const messageDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())

  if (messageDate.getTime() === today.getTime()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else if (messageDate.getTime() === today.getTime() - 86400000) {
    return '昨天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else {
    return date.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
  }
}

// 按日期分组消息
const groupedMessages = computed(() => {
  const groups = []
  let currentDate = null

  messages.value.forEach(msg => {
    const msgDate = new Date(msg.createdAt)
    const dateKey = msgDate.toDateString()

    if (dateKey !== currentDate) {
      currentDate = dateKey
      const today = new Date().toDateString()
      const yesterday = new Date(Date.now() - 86400000).toDateString()

      let dateLabel
      if (dateKey === today) dateLabel = '今天'
      else if (dateKey === yesterday) dateLabel = '昨天'
      else dateLabel = msgDate.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })

      groups.push({ type: 'date', label: dateLabel, dateKey })
    }
    groups.push({ type: 'message', ...msg })
  })

  return groups
})

// 检查消息是否在5分钟内可以撤回
function canRecall(msg) {
  if (!isMyMessage(msg)) return false
  if (msg.isDeleted) return false
  const createdTime = new Date(msg.createdAt)
  return (Date.now() - createdTime.getTime()) < 5 * 60 * 1000
}
</script>

<template>
  <div class="container messages-page">
    <BackButton />
    <h1 class="page-title">消息沟通</h1>

    <div class="chat-container card">
      <div v-if="loading" class="empty-state">加载中...</div>

      <div v-else ref="messagesContainer" class="chat-messages">
        <div v-if="!messages.length" class="empty-state">暂无消息，发送第一条消息吧</div>
        <template v-for="(item, index) in groupedMessages" :key="index">
          <!-- 日期分隔线 -->
          <div v-if="item.type === 'date'" class="date-divider">
            <span>{{ item.label }}</span>
          </div>
          <!-- 消息气泡 -->
          <div
            v-else
            class="chat-bubble"
            :class="{ mine: isMyMessage(item), other: !isMyMessage(item), deleted: item.isDeleted }"
            @dblclick="canRecall(item) && handleRecall(item)"
          >
            <div class="bubble-name">{{ isMyMessage(item) ? '我' : (item.sender?.nickname || item.sender?.username || '对方') }}</div>
            <div class="bubble-content">
              <!-- 支持多图片 -->
              <div v-if="item.images && item.images.length > 0" class="message-images" :class="`count-${item.images.length}`">
                <img v-for="(img, idx) in item.images" :key="idx" :src="img" class="message-image" @click="openImagePreview(img)" />
              </div>
              <span v-if="item.content">{{ item.content }}</span>
            </div>
            <div class="bubble-footer">
              <span class="bubble-time">{{ formatTime(item.createdAt) }}</span>
              <span v-if="item.isDeleted" class="recall-hint">（已撤回，双击撤回）</span>
              <span v-else-if="isMyMessage(item) && canRecall(item)" class="recall-hint">（双击撤回）</span>
              <!-- 已读/未读状态 -->
              <span v-if="isMyMessage(item) && !item.isDeleted" class="read-status" :class="{ unread: !item.isRead }">
                {{ item.isRead ? '已读' : '未读' }}
              </span>
            </div>
          </div>
        </template>
      </div>

      <!-- 图片预览区域 -->
      <div v-if="selectedImages.length > 0" class="image-preview-bar">
        <div v-for="img in selectedImages" :key="img.id" class="preview-item">
          <img :src="img.previewUrl" class="preview-thumb" />
          <button class="preview-remove" @click="removeImage(img)" title="删除">×</button>
        </div>
      </div>

      <div class="chat-input">
        <input type="file" ref="fileInput" accept="image/*" multiple style="display:none" @change="handleFileChange" />
        <button class="btn btn-icon" @click="triggerFileInput" :disabled="uploading" title="发送图片">
          {{ uploading ? '...' : '📷' }}
        </button>
        <input
          v-model="newMessage"
          class="form-input"
          placeholder="输入消息...（双击自己的消息可撤回）"
          @keyup.enter="handleSend"
        />
 <button class="btn btn-primary" @click="handleSend">发送</button>
      </div>
    </div>

    <!-- 图片预览弹窗 -->
    <div v-if="previewImage" class="image-preview-modal" @click="closeImagePreview">
      <button class="preview-close" @click="closeImagePreview">×</button>
      <img :src="previewImage" class="preview-full" @click.stop />
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

/* 日期分隔线 */
.date-divider {
  text-align: center;
  padding: 8px 0;
  color: var(--gray-400);
  font-size: 12px;
}

.date-divider span {
  background: var(--gray-100);
  padding: 4px 12px;
  border-radius: 10px;
}

/* 消息气泡 */
.chat-bubble {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: var(--radius);
  cursor: default;
}

.chat-bubble.deleted {
  opacity: 0.5;
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

/* 多图片网格布局 */
.message-images {
  display: grid;
  gap: 4px;
  margin-bottom: 8px;
}

.message-images.count-1 {
  grid-template-columns: 1fr;
}

.message-images.count-2 {
  grid-template-columns: repeat(2, 1fr);
}

.message-images.count-3,
.message-images.count-4 {
  grid-template-columns: repeat(2, 1fr);
}

.message-images.count-5 {
  grid-template-columns: repeat(3, 1fr);
}

.message-images.count-6,
.message-images.count-7,
.message-images.count-8,
.message-images.count-9 {
  grid-template-columns: repeat(3, 1fr);
}

.message-image {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
}

.message-image:hover {
  transform: scale(1.02);
}

.message-images.count-1 .message-image {
  width: 200px;
  height: 200px;
}

.bubble-footer {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 4px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}

.bubble-time {
  opacity: 0.6;
}

.recall-hint {
  opacity: 0.5;
  font-size: 10px;
}

.read-status {
  font-size: 10px;
  opacity: 0.5;
  padding: 1px 4px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.2);
}

.read-status.unread {
  opacity: 0.8;
  background: rgba(255, 255, 255, 0.3);
}

.chat-input {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid var(--gray-100);
  align-items: center;
}

.chat-input .form-input {
  flex: 1;
}

.btn-icon {
  background: var(--gray-100);
  border: none;
  padding: 8px 12px;
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 18px;
}

.btn-icon:hover {
  background: var(--gray-200);
}

.btn-icon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 图片预览区域 */
.image-preview-bar {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--gray-100);
  background: var(--gray-50);
  overflow-x: auto;
  flex-wrap: wrap;
}

.preview-item {
  position: relative;
  width: 60px;
  height: 60px;
  flex-shrink: 0;
}

.preview-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  border: 2px solid var(--gray-200);
}

.preview-remove {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--danger, #dc3545);
  color: #fff;
  border: 2px solid #fff;
  font-size: 14px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.preview-remove:hover {
  background: #c82333;
  transform: scale(1.1);
}

/* 图片预览弹窗 */
.image-preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  cursor: pointer;
}

.preview-full {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  cursor: default;
}

.preview-close {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border: none;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-close:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>

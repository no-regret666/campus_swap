<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getNotifications, markNotificationRead, markAllNotificationsRead } from '../api/notifications'
import { useAppStore } from '../stores/app'
import BackButton from '../components/BackButton.vue'

const router = useRouter()
const appStore = useAppStore()

const notifications = ref([])
const loading = ref(false)

onMounted(loadNotifications)

async function loadNotifications() {
  loading.value = true
  try {
    const res = await getNotifications()
    notifications.value = res.notifications || []
  } catch (e) {
    appStore.showToast('加载通知失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleMarkAllRead() {
  try {
    await markAllNotificationsRead()
    appStore.showToast('已全部标记已读', 'success')
    loadNotifications()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleClick(notif) {
  // 标记已读
  if (!notif.isRead) {
    await markNotificationRead(notif.id || notif._id).catch(() => {})
  }

  // 跳转到相关页面
  if (notif.type === 'exchange_request' || notif.type === 'exchange_response') {
    router.push('/exchanges')
  } else if (notif.type === 'new_message' || notif.type === 'message') {
    router.push(`/messages/${notif.relatedId}`)
  }
}

function formatTime(timestamp) {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff/60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff/3600000)}小时前`
  return date.toLocaleDateString('zh-CN')
}

function getTypeIcon(type) {
  const icons = {
    exchange_request: '📩',
    exchange_response: '📬',
    exchange_reminder: '⏰',
    message: '💬'
  }
  return icons[type] || '🔔'
}
</script>

<template>
  <div class="container">
    <BackButton />
    <div class="page-header">
      <h1 class="page-title">通知</h1>
      <button v-if="notifications.length" class="btn btn-secondary btn-sm" @click="handleMarkAllRead">
        全部已读
      </button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="!notifications.length" class="empty-state">
      <div class="empty-icon">🔔</div>
      <p>暂无通知</p>
    </div>

    <div v-else class="notification-list">
      <div
        v-for="notif in notifications"
        :key="notif.id || notif._id"
        class="notification-item card"
        :class="{ unread: !notif.isRead }"
        @click="handleClick(notif)"
      >
        <div class="notification-icon">{{ getTypeIcon(notif.type) }}</div>
        <div class="notification-content">
          <div class="notification-title">{{ notif.title }}</div>
          <div class="notification-text">{{ notif.content }}</div>
          <div class="notification-time">{{ formatTime(notif.createdAt) }}</div>
        </div>
        <div v-if="!notif.isRead" class="unread-dot"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.notification-item:hover {
  background: var(--gray-50);
}

.notification-item.unread {
  background: #f0f7ff;
}

.notification-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-weight: 600;
  color: var(--gray-800);
  margin-bottom: 4px;
}

.notification-text {
  font-size: 14px;
  color: var(--gray-600);
  margin-bottom: 6px;
  line-height: 1.4;
}

.notification-time {
  font-size: 12px;
  color: var(--gray-400);
}

.unread-dot {
  width: 8px;
  height: 8px;
  background: var(--primary);
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 6px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--gray-400);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}
</style>
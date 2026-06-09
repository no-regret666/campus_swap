<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import { markNotificationRead } from '../api/misc'

const userStore = useUserStore()
const appStore = useAppStore()
const router = useRouter()
const mobileMenuOpen = ref(false)
const showNotifPanel = ref(false)

// 登录后自动加载通知
onMounted(() => {
  if (userStore.isLoggedIn) {
    appStore.fetchNotifications()
  }
})

watch(() => userStore.isLoggedIn, (val) => {
  if (val) appStore.fetchNotifications()
})

const unreadCount = computed(() => {
  return appStore.notifications.filter(n => !n.isRead).length
})

function handleLogout() {
  userStore.logout()
}

async function markRead(notif) {
  if (notif.isRead) return
  try {
    await markNotificationRead(notif.id)
    notif.isRead = true
  } catch (e) {
    // ignore
  }
}

async function clickNotif(notif) {
  await markRead(notif)
  showNotifPanel.value = false
  // 根据类型跳转
  if (notif.type === 'exchange_created' || notif.type === 'exchange_status') {
    router.push('/exchanges')
  } else if (notif.type === 'new_message' && notif.relatedId) {
    router.push(`/messages/${notif.relatedId}`)
  }
}

async function markAllRead() {
  for (const n of appStore.notifications) {
    if (!n.isRead) {
      await markRead(n)
    }
  }
}

function formatNotifTime(timeStr) {
  if (!timeStr) return ''
  try {
    const d = new Date(timeStr)
    const now = new Date()
    const diff = now - d
    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
    if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
    return d.toLocaleDateString()
  } catch {
    return timeStr
  }
}
</script>

<template>
  <header class="navbar">
    <div class="navbar-inner container">
      <router-link to="/" class="navbar-brand">
        <span class="brand-icon">🔄</span>
        <span class="brand-text">闲置交换</span>
      </router-link>

      <button class="mobile-toggle" @click="mobileMenuOpen = !mobileMenuOpen">
        <span></span><span></span><span></span>
      </button>

      <nav class="navbar-nav" :class="{ open: mobileMenuOpen }" @click="mobileMenuOpen = false">
        <router-link to="/" class="nav-link">广场</router-link>
        <router-link to="/recommendations" class="nav-link">推荐</router-link>
        <template v-if="userStore.isLoggedIn">
          <router-link to="/publish" class="nav-link">发布</router-link>
          <router-link to="/exchanges" class="nav-link">交换</router-link>
          <router-link to="/favorites" class="nav-link">收藏</router-link>
          <router-link to="/admin" class="nav-link" v-if="userStore.isAdmin">后台</router-link>

          <!-- 通知铃铛 -->
          <div class="notif-wrapper" v-if="userStore.isLoggedIn">
            <button class="nav-link notif-btn" @click.stop="showNotifPanel = !showNotifPanel">
              🔔
              <span v-if="unreadCount" class="notif-badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
            </button>
            <div v-if="showNotifPanel" class="notif-panel" @click.stop>
              <div class="notif-header">
                <span>通知</span>
                <button v-if="unreadCount" class="notif-mark-all" @click="markAllRead">全部已读</button>
              </div>
              <div class="notif-list">
                <div v-if="!appStore.notifications.length" class="notif-empty">暂无通知</div>
                <div
                  v-for="n in [...appStore.notifications].reverse().slice(0, 20)"
                  :key="n.id"
                  class="notif-item"
                  :class="{ unread: !n.isRead }"
                  @click="clickNotif(n)"
                >
                  <div class="notif-dot" v-if="!n.isRead"></div>
                  <div class="notif-content">
                    <div class="notif-title">{{ n.title }}</div>
                    <div class="notif-desc">{{ n.content }}</div>
                    <div class="notif-time">{{ formatNotifTime(n.createdAt) }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <router-link to="/profile" class="nav-link nav-user">
            <span class="user-avatar">{{ userStore.user?.name?.[0] || '我' }}</span>
            <span>{{ userStore.user?.name || '个人中心' }}</span>
          </router-link>
          <button class="nav-link logout-btn" @click="handleLogout">退出</button>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link btn-login">登录</router-link>
        </template>
      </nav>
    </div>
  </header>

  <!-- 点击其他区域关闭通知面板 -->
  <div v-if="showNotifPanel" class="notif-overlay" @click="showNotifPanel = false"></div>
</template>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: white;
  border-bottom: 1px solid var(--gray-200);
  z-index: 1000;
  backdrop-filter: blur(8px);
  background: rgba(255,255,255,0.95);
}

.navbar-inner {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  font-size: 18px;
  color: var(--primary);
}

.brand-icon {
  font-size: 24px;
}

.navbar-nav {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-link {
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  color: var(--gray-600);
  transition: all 0.2s;
  background: none;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--primary);
  background: var(--primary-light);
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 6px;
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.btn-login {
  background: var(--primary) !important;
  color: white !important;
  border-radius: var(--radius-sm);
}

.btn-login:hover {
  background: var(--primary-hover) !important;
}

.logout-btn {
  color: var(--gray-400);
  font-size: 13px;
}

.logout-btn:hover {
  color: var(--danger);
}

/* 通知铃铛 */
.notif-wrapper {
  position: relative;
}

.notif-btn {
  position: relative;
  font-size: 16px;
  padding: 8px 10px;
}

.notif-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 16px;
  height: 16px;
  line-height: 16px;
  font-size: 10px;
  font-weight: 700;
  background: #e53935;
  color: white;
  border-radius: 8px;
  text-align: center;
  padding: 0 4px;
}

.notif-panel {
  position: absolute;
  top: 100%;
  right: 0;
  width: 340px;
  max-height: 420px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.12);
  border: 1px solid var(--gray-200);
  overflow: hidden;
  z-index: 2000;
}

.notif-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px 10px;
  font-weight: 600;
  font-size: 15px;
  color: var(--gray-800);
  border-bottom: 1px solid var(--gray-100);
}

.notif-mark-all {
  background: none;
  border: none;
  color: var(--primary);
  font-size: 13px;
  cursor: pointer;
  padding: 2px 6px;
}

.notif-mark-all:hover {
  text-decoration: underline;
}

.notif-list {
  max-height: 350px;
  overflow-y: auto;
}

.notif-empty {
  padding: 32px 16px;
  text-align: center;
  color: var(--gray-400);
  font-size: 14px;
}

.notif-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid var(--gray-50);
}

.notif-item:hover {
  background: var(--gray-50);
}

.notif-item.unread {
  background: #f0f7ff;
}

.notif-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--primary);
  flex-shrink: 0;
  margin-top: 6px;
}

.notif-content {
  flex: 1;
  min-width: 0;
}

.notif-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--gray-800);
  margin-bottom: 2px;
}

.notif-desc {
  font-size: 12px;
  color: var(--gray-500);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notif-time {
  font-size: 11px;
  color: var(--gray-400);
}

.notif-overlay {
  position: fixed;
  inset: 0;
  z-index: 999;
}

/* 移动端 */
.mobile-toggle {
  display: none;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: none;
  border: none;
}

.mobile-toggle span {
  display: block;
  width: 20px;
  height: 2px;
  background: var(--gray-600);
  border-radius: 2px;
  transition: 0.3s;
}

@media (max-width: 768px) {
  .mobile-toggle {
    display: flex;
  }
  .navbar-nav {
    position: fixed;
    top: 64px;
    left: 0;
    right: 0;
    background: white;
    flex-direction: column;
    padding: 16px;
    border-bottom: 1px solid var(--gray-200);
    box-shadow: var(--shadow-lg);
    transform: translateY(-100%);
    opacity: 0;
    pointer-events: none;
    transition: all 0.3s ease;
  }
  .navbar-nav.open {
    transform: translateY(0);
    opacity: 1;
    pointer-events: auto;
  }
  .nav-link {
    width: 100%;
    padding: 12px 14px;
  }
  .notif-panel {
    width: 300px;
    right: -40px;
  }
}
</style>

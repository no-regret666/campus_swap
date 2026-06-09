<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import { getUnreadCount, getAllUnreadByExchange } from '../api/messages'
import { getUnreadNotificationCount } from '../api/notifications'

const userStore = useUserStore()
const appStore = useAppStore()
const router = useRouter()
const mobileMenuOpen = ref(false)
const unreadCount = ref(0)
const notifCount = ref(0)

// 获取未读消息数和通知数
onMounted(() => {
  if (userStore.isLoggedIn) {
    loadUnreadCounts()
  }
})

async function loadUnreadCounts() {
  try {
    const [msgRes, notifRes] = await Promise.all([
      getUnreadCount(),
      getUnreadNotificationCount().catch(() => ({ unreadCount: 0 }))
    ])
    unreadCount.value = msgRes.unreadCount || 0
    notifCount.value = notifRes.unreadCount || 0
  } catch (e) {
    // ignore
  }
}

function handleLogout() {
  userStore.logout()
}

// 刷新未读数（当从其他页面返回时）
router.afterEach(() => {
  if (userStore.isLoggedIn) {
    loadUnreadCounts()
  }
})
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
        <template v-if="userStore.isLoggedIn">
          <router-link to="/publish" class="nav-link">发布</router-link>
          <router-link to="/exchanges" class="nav-link">
            交换<span v-if="unreadCount" class="badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
          </router-link>
          <router-link to="/notifications" class="nav-link nav-icon-link">
            🔔<span v-if="notifCount" class="badge">{{ notifCount > 99 ? '99+' : notifCount }}</span>
          </router-link>
          <router-link to="/favorites" class="nav-link">收藏</router-link>
          <router-link to="/admin" class="nav-link" v-if="userStore.isAdmin">后台</router-link>
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
  position: relative;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--primary);
  background: var(--primary-light);
}

.nav-icon-link {
  font-size: 16px;
}

.badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: var(--danger);
  color: white;
  border-radius: 8px;
  font-size: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
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
}
</style>

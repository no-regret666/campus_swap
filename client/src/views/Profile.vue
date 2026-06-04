<script setup>
import { ref, onMounted } from 'vue'
import { getItems } from '../api/items'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'

const userStore = useUserStore()
const appStore = useAppStore()

const myItems = ref([])
const loading = ref(false)

onMounted(async () => {
  if (!userStore.user) await userStore.fetchUser()
  loadMyItems()
})

async function loadMyItems() {
  loading.value = true
  try {
    const userId = userStore.user?._id || userStore.user?.id
    const res = await getItems({ owner: userId })
    myItems.value = res.items || res || []
  } catch (e) {
    appStore.showToast('加载物品失败', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">个人中心</h1>

    <!-- 用户信息 -->
    <div class="profile-card card" v-if="userStore.user">
      <div class="profile-header">
        <div class="avatar">{{ (userStore.user.nickname || userStore.user.username || '?')[0] }}</div>
        <div class="profile-info">
          <h2>{{ userStore.user.nickname || userStore.user.username }}</h2>
          <p class="profile-meta">
            <span>👤 {{ userStore.user.username }}</span>
            <span>📍 {{ userStore.user.campus || '未设置' }}</span>
            <span class="badge" :class="userStore.isAdmin ? 'badge-danger' : 'badge-primary'">
              {{ userStore.isAdmin ? '管理员' : '普通用户' }}
            </span>
          </p>
        </div>
      </div>
      <div class="profile-stats" v-if="userStore.user.creditScore !== undefined">
        <div class="stat-item">
          <span class="stat-value">{{ userStore.user.creditScore || '-' }}</span>
          <span class="stat-label">信用评分</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ myItems.length }}</span>
          <span class="stat-label">发布物品</span>
        </div>
      </div>
    </div>

    <!-- 我发布的物品 -->
    <h2 class="section-title">我发布的物品</h2>
    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="myItems.length" class="grid-items">
      <ItemCard v-for="item in myItems" :key="item._id || item.id" :item="item" />
    </div>
    <div v-else class="empty-state">
      <p>暂无发布的物品</p>
    </div>
  </div>
</template>

<style scoped>
.profile-card {
  padding: 24px;
  margin-bottom: 32px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
}

.profile-info h2 {
  font-size: 20px;
  margin-bottom: 4px;
}

.profile-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  font-size: 14px;
  color: var(--gray-500);
}

.profile-stats {
  display: flex;
  gap: 32px;
  padding-top: 16px;
  border-top: 1px solid var(--gray-100);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--primary);
}

.stat-label {
  font-size: 13px;
  color: var(--gray-500);
}

.section-title {
  font-size: 18px;
  margin-bottom: 16px;
  color: var(--gray-800);
}
</style>

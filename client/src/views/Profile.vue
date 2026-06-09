<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getItems, deleteItem } from '../api/items'
import { getRatings } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'
import BackButton from '../components/BackButton.vue'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const myItems = ref([])
const loading = ref(false)
const ratings = ref([])
const avgScore = ref(0)
const ratingCount = ref(0)
const ratingsLoading = ref(false)

onMounted(async () => {
  if (!userStore.user) await userStore.fetchUser()
  loadMyItems()
  loadMyRatings()
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

async function loadMyRatings() {
  const userId = userStore.user?._id || userStore.user?.id
  if (!userId) return
  ratingsLoading.value = true
  try {
    const res = await getRatings({ userId })
    ratings.value = res.ratings || []
    avgScore.value = res.avgScore ?? 0
    ratingCount.value = res.ratingCount ?? 0
  } catch (e) {
    // 静默处理
  } finally {
    ratingsLoading.value = false
  }
}

function goEdit(item) {
  router.push(`/item/${item._id || item.id}`)
}

async function handleDelete(item) {
  if (!confirm(`确定删除「${item.title}」吗？删除后不可恢复。`)) return
  try {
    await deleteItem(item._id || item.id)
    appStore.showToast('删除成功', 'success')
    myItems.value = myItems.value.filter(i => (i._id || i.id) !== (item._id || item.id))
  } catch (e) {
    appStore.showToast('删除失败', 'error')
  }
}

function stars(score) {
  return '★'.repeat(score) + '☆'.repeat(5 - score)
}

function formatRatingTime(timeStr) {
  if (!timeStr) return ''
  try {
    return new Date(timeStr).toLocaleDateString()
  } catch {
    return timeStr
  }
}
</script>

<template>
  <div class="container">
    <BackButton />
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
      <div class="profile-stats">
        <div class="stat-item">
          <span class="stat-value">{{ userStore.user.creditScore || '-' }}</span>
          <span class="stat-label">信用评分</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ myItems.length }}</span>
          <span class="stat-label">发布物品</span>
        </div>
        <div class="stat-item">
          <span class="stat-value rating-stars-main">{{ avgScore > 0 ? avgScore.toFixed(1) : '-' }}</span>
          <span class="stat-label">平均评分 ({{ ratingCount }}条)</span>
        </div>
      </div>
    </div>

    <!-- 我收到的评价 -->
    <h2 class="section-title">我收到的评价</h2>
    <div v-if="ratingsLoading" class="empty-state">加载中...</div>
    <div v-else-if="ratings.length" class="ratings-list">
      <div v-for="r in ratings" :key="r.id" class="rating-card card">
        <div class="rating-header">
          <span class="rating-score">{{ stars(r.score) }}</span>
          <span class="rating-date">{{ formatRatingTime(r.createdAt) }}</span>
        </div>
        <p v-if="r.comment" class="rating-comment">{{ r.comment }}</p>
        <p v-else class="rating-comment empty-comment">未留下评价</p>
      </div>
    </div>
    <div v-else class="empty-state">
      <p>暂无评价</p>
    </div>

    <!-- 我发布的物品 -->
    <h2 class="section-title">我发布的物品</h2>
    <div v-if="loading" class="empty-state">加载中...</div>
    <div v-else-if="myItems.length" class="my-items-list">
      <div v-for="item in myItems" :key="item._id || item.id" class="my-item-row card">
        <div class="my-item-left" @click="goEdit(item)">
          <img v-if="item.images && item.images.length" :src="item.images[0]" class="my-item-img" />
          <div v-else class="my-item-img placeholder">📦</div>
          <div class="my-item-info">
            <h4>{{ item.title }}</h4>
            <p>{{ item.category }} · {{ item.condition || '未标注' }} · {{ item.campus }}</p>
            <span class="badge" :class="item.status === 'available' ? 'badge-success' : 'badge-warning'">
              {{ item.status === 'available' ? '可交换' : '已交换' }}
            </span>
          </div>
        </div>
        <div class="my-item-actions">
          <button class="btn btn-secondary btn-sm" @click="goEdit(item)">✏️ 编辑</button>
          <button class="btn btn-danger btn-sm" @click="handleDelete(item)">🗑️ 删除</button>
        </div>
      </div>
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

.rating-stars-main {
  color: #f5a623;
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

/* 评价列表 */
.ratings-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 32px;
}

.rating-card {
  padding: 16px 20px;
}

.rating-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.rating-score {
  color: #f5a623;
  font-size: 16px;
  letter-spacing: 2px;
}

.rating-date {
  font-size: 12px;
  color: var(--gray-400);
}

.rating-comment {
  font-size: 14px;
  color: var(--gray-600);
  line-height: 1.5;
}

.empty-comment {
  color: var(--gray-400);
  font-style: italic;
}

/* 物品列表 */
.my-items-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.my-item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
}

.my-item-left {
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  flex: 1;
  min-width: 0;
}

.my-item-img {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  object-fit: cover;
  flex-shrink: 0;
}

.my-item-img.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gray-100);
  font-size: 24px;
}

.my-item-info {
  min-width: 0;
}

.my-item-info h4 {
  font-size: 15px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.my-item-info p {
  font-size: 12px;
  color: var(--gray-500);
  margin-bottom: 4px;
}

.my-item-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  margin-left: 12px;
}
</style>

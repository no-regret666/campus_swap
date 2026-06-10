<script setup>
import { ref, onMounted, computed } from 'vue'
import { getItems } from '../api/items'
import { getRecommendations } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'

const appStore = useAppStore()
const userStore = useUserStore()

const items = ref([])
const recommendations = ref([])
const loading = ref(false)
const recLoading = ref(false)
const keyword = ref('')
const category = ref('')
const campus = ref('')
const sortBy = ref('newest')

const campusOptions = ['', '南湖校区', '东湖校区', '线上']
const sortOptions = [
  { value: 'newest', label: '最新发布' },
  { value: 'views', label: '最多浏览' }
]

// 未登录时的热门推荐：取浏览量前6的可用物品
const hotItems = computed(() => {
  return [...items.value]
    .filter(i => i.status === 'available')
    .sort((a, b) => (b.views || 0) - (a.views || 0))
    .slice(0, 6)
})

onMounted(() => {
  if (!appStore.categories.length) appStore.fetchCategories()
  loadItems()
  loadRecommendations()
})

async function loadItems() {
  loading.value = true
  try {
    const params = {}
    if (keyword.value) params.keyword = keyword.value
    if (category.value) params.category = category.value
    if (campus.value) params.campus = campus.value
    if (sortBy.value) params.sortBy = sortBy.value
    const res = await getItems(params)
    items.value = res.items || res || []
  } catch (e) {
    appStore.showToast('加载物品失败', 'error')
  } finally {
    loading.value = false
  }
}

async function loadRecommendations() {
  if (!userStore.isLoggedIn) return
  recLoading.value = true
  try {
    const res = await getRecommendations()
    recommendations.value = res.items || res || []
  } catch (e) {
    // 静默处理
  } finally {
    recLoading.value = false
  }
}

function handleSearch() {
  loadItems()
}

function resetFilters() {
  keyword.value = ''
  category.value = ''
  campus.value = ''
  sortBy.value = 'newest'
  loadItems()
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">闲置广场</h1>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input v-model="keyword" class="form-input" placeholder="搜索物品..." @keyup.enter="handleSearch" />
      <select v-model="category" class="form-input" @change="handleSearch">
        <option value="">全部分类</option>
        <option v-for="cat in appStore.categories" :key="cat._id || cat.id || cat.name" :value="cat._id || cat.id || cat.name">
          {{ cat.name || cat }}
        </option>
      </select>
      <select v-model="campus" class="form-input" @change="handleSearch">
        <option value="">全部校区</option>
        <option v-for="c in campusOptions.slice(1)" :key="c" :value="c">{{ c }}</option>
      </select>
      <select v-model="sortBy" class="form-input" @change="handleSearch">
        <option v-for="s in sortOptions" :key="s.value" :value="s.value">{{ s.label }}</option>
      </select>
      <button class="btn btn-primary" @click="handleSearch">搜索</button>
      <button class="btn btn-secondary" @click="resetFilters">重置</button>
    </div>

    <!-- 为你推荐（已登录） -->
    <div v-if="userStore.isLoggedIn && recommendations.length" class="recommend-section">
      <div class="section-header">
        <h2 class="section-title">✨ 为你推荐</h2>
        <span class="section-hint">基于你的浏览和收藏偏好</span>
      </div>
      <div class="rec-scroll">
        <div v-for="item in recommendations" :key="item._id || item.id" class="rec-item">
          <ItemCard :item="item" />
        </div>
      </div>
    </div>

    <!-- 热门推荐（未登录） -->
    <div v-if="!userStore.isLoggedIn && hotItems.length" class="recommend-section hot-section">
      <div class="section-header">
        <h2 class="section-title">🔥 热门物品</h2>
        <span class="section-hint">登录后查看个性化推荐</span>
      </div>
      <div class="rec-scroll">
        <div v-for="item in hotItems" :key="item._id || item.id" class="rec-item">
          <ItemCard :item="item" />
        </div>
      </div>
    </div>

    <!-- 全部物品 -->
    <div class="all-items-header">
      <h2 class="section-title">📦 全部物品</h2>
      <span class="item-count">共 {{ items.length }} 件</span>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="empty-state">加载中...</div>

    <!-- 物品列表 -->
    <div v-else-if="items.length" class="grid-items">
      <ItemCard v-for="item in items" :key="item._id || item.id" :item="item" />
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <p>暂无物品，快来发布第一件闲置吧！</p>
      <router-link v-if="userStore.isLoggedIn" to="/publish" class="btn btn-primary" style="margin-top:12px">去发布</router-link>
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-bar .form-input {
  flex: 1;
  min-width: 150px;
}

.search-bar select {
  flex: 0 0 140px;
}

/* 推荐区块 */
.recommend-section {
  margin-bottom: 32px;
  padding: 20px;
  background: linear-gradient(135deg, #f0f7ff 0%, #f5f0ff 100%);
  border-radius: var(--radius);
  border: 1px solid #e0e8f5;
}

.hot-section {
  background: linear-gradient(135deg, #fff8f0 0%, #fff0f0 100%);
  border-color: #f5e0d8;
}

.section-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--gray-800);
}

.section-hint {
  font-size: 12px;
  color: var(--gray-400);
}

.rec-scroll {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding-bottom: 8px;
  scroll-snap-type: x mandatory;
}

.rec-scroll::-webkit-scrollbar {
  height: 4px;
}

.rec-scroll::-webkit-scrollbar-thumb {
  background: var(--gray-300);
  border-radius: 2px;
}

.rec-item {
  flex: 0 0 220px;
  scroll-snap-align: start;
}

.all-items-header {
  margin-bottom: 16px;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.item-count {
  font-size: 13px;
  color: var(--gray-400);
}
</style>

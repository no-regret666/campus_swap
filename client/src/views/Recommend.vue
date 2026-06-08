<script setup>
import { ref, computed, onMounted } from 'vue'
import { getRecommendations, getCategories } from '../api/misc'
import { useAppStore } from '../stores/app'
import { useUserStore } from '../stores/user'
import ItemCard from '../components/ItemCard.vue'

const appStore = useAppStore()
const userStore = useUserStore()

const items = ref([])
const categories = ref([])
const selectedCategory = ref('')
const loading = ref(false)

onMounted(async () => {
  await loadCategories()
  await loadRecommendations()
})

async function loadCategories() {
  try {
    const res = await getCategories()
    categories.value = res.categories || []
  } catch (e) {
    // 静默处理
  }
}

async function loadRecommendations() {
  loading.value = true
  try {
    const res = await getRecommendations()
    items.value = res.recommendations || res.items || res || []
  } catch (e) {
    appStore.showToast('加载推荐失败', 'error')
  } finally {
    loading.value = false
  }
}

const filteredItems = computed(() => {
  if (!selectedCategory.value) return items.value
  return items.value.filter(item => {
    const cat = item.category || item.categoryId || ''
    return cat === selectedCategory.value
  })
})

function getCategoryName(catId) {
  const cat = categories.value.find(c => c.id === catId)
  return cat ? cat.name : catId
}

function refresh() {
  selectedCategory.value = ''
  loadRecommendations()
}
</script>

<template>
  <div class="container">
    <div class="page-header">
      <h1 class="page-title">智能推荐</h1>
      <button class="btn btn-secondary" @click="refresh" :disabled="loading">
        🔄 {{ loading ? '加载中...' : '刷新推荐' }}
      </button>
    </div>

    <!-- 分类筛选 -->
    <div class="filter-bar" v-if="categories.length">
      <button
        class="filter-chip"
        :class="{ active: !selectedCategory }"
        @click="selectedCategory = ''"
      >全部</button>
      <button
        v-for="cat in categories"
        :key="cat.id"
        class="filter-chip"
        :class="{ active: selectedCategory === cat.id }"
        @click="selectedCategory = cat.id"
      >{{ cat.icon }} {{ cat.name }}</button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="filteredItems.length" class="grid-items">
      <div v-for="item in filteredItems" :key="item._id || item.id" class="recommend-item">
        <ItemCard :item="item" />
        <div v-if="item.reason" class="recommend-reason">
          <span class="reason-icon">💡</span>
          <span class="reason-text">{{ item.reason }}</span>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon">🔍</div>
      <p v-if="selectedCategory">该分类暂无推荐物品</p>
      <p v-else>暂无推荐内容，试试收藏一些物品或浏览更多吧</p>
      <router-link v-if="!userStore.isLoggedIn" to="/login" class="btn btn-primary btn-sm" style="margin-top: 12px;">登录获取个性化推荐</router-link>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header .page-title {
  margin-bottom: 0;
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 24px;
}

.filter-chip {
  padding: 6px 16px;
  border-radius: 100px;
  font-size: 13px;
  font-weight: 500;
  background: white;
  color: var(--gray-600);
  border: 1.5px solid var(--gray-200);
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-chip:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.filter-chip.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.recommend-item {
  position: relative;
}

.recommend-reason {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--primary-light);
  border-radius: 0 0 var(--radius) var(--radius);
  font-size: 12px;
  color: var(--primary);
  margin-top: -4px;
}

.reason-icon {
  font-size: 14px;
}

.reason-text {
  font-weight: 500;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}
</style>

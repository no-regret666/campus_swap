<script setup>
import { ref, onMounted } from 'vue'
import { getRecommendations } from '../api/misc'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'

const appStore = useAppStore()

const items = ref([])
const loading = ref(false)

onMounted(loadRecommendations)

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

function refresh() {
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

    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="items.length" class="grid-items">
      <ItemCard v-for="item in items" :key="item._id || item.id" :item="item" />
    </div>

    <div v-else class="empty-state">
      <p>暂无推荐内容，试试发布或浏览更多物品吧</p>
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
</style>

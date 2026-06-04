<script setup>
import { ref, onMounted } from 'vue'
import { getItems } from '../api/items'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'

const appStore = useAppStore()

const items = ref([])
const loading = ref(false)
const keyword = ref('')
const category = ref('')
const campus = ref('')

const campusOptions = ['', '南湖校区', '东湖校区', '线上']

onMounted(() => {
  if (!appStore.categories.length) appStore.fetchCategories()
  loadItems()
})

async function loadItems() {
  loading.value = true
  try {
    const params = {}
    if (keyword.value) params.keyword = keyword.value
    if (category.value) params.category = category.value
    if (campus.value) params.campus = campus.value
    const res = await getItems(params)
    items.value = res.items || res || []
  } catch (e) {
    appStore.showToast('加载物品失败', 'error')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  loadItems()
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">闲置广场</h1>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input v-model="keyword" class="form-input" placeholder="搜索物品..." @keyup.enter="handleSearch" />
      <select v-model="category" class="form-input">
        <option value="">全部分类</option>
        <option v-for="cat in appStore.categories" :key="cat._id || cat.id || cat.name" :value="cat._id || cat.id || cat.name">
          {{ cat.name || cat }}
        </option>
      </select>
      <select v-model="campus" class="form-input">
        <option value="">全部校区</option>
        <option v-for="c in campusOptions.slice(1)" :key="c" :value="c">{{ c }}</option>
      </select>
      <button class="btn btn-primary" @click="handleSearch">搜索</button>
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
</style>

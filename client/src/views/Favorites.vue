<script setup>
import { ref, onMounted } from 'vue'
import { getFavorites, removeFavorite } from '../api/favorites'
import { useAppStore } from '../stores/app'
import ItemCard from '../components/ItemCard.vue'
import BackButton from '../components/BackButton.vue'

const appStore = useAppStore()

const favorites = ref([])
const loading = ref(false)

onMounted(loadFavorites)

async function loadFavorites() {
  loading.value = true
  try {
    const res = await getFavorites()
    favorites.value = res.favorites || res || []
  } catch (e) {
    appStore.showToast('加载收藏列表失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleRemove(item) {
  try {
    await removeFavorite(item._id || item.id)
    favorites.value = favorites.value.filter(f => (f._id || f.id) !== (item._id || item.id))
    appStore.showToast('已取消收藏', 'success')
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}
</script>

<template>
  <div class="container">
    <BackButton />
    <h1 class="page-title">我的收藏</h1>

    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="favorites.length" class="grid-items">
      <div v-for="item in favorites" :key="item._id || item.id" class="fav-item">
        <ItemCard :item="item.item || item" />
        <button class="btn btn-danger btn-sm remove-btn" @click="handleRemove(item.item || item)">
          取消收藏
        </button>
      </div>
    </div>

    <div v-else class="empty-state">
      <p>暂无收藏的物品</p>
    </div>
  </div>
</template>

<style scoped>
.fav-item {
  position: relative;
}

.remove-btn {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 2;
}
</style>

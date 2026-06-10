<script setup>
import { ref, onMounted, computed } from 'vue'
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

// 过滤掉物品已被删除的收藏
const validFavorites = computed(() => {
  return favorites.value.filter(f => {
    const it = f.item || f
    return it && it.id
  })
})

const removedCount = computed(() => favorites.value.length - validFavorites.value.length)

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

    <div v-else-if="validFavorites.length" class="grid-items">
      <div v-for="item in validFavorites" :key="item._id || item.id" class="fav-item">
        <ItemCard :item="item.item || item" />
        <button class="btn btn-danger btn-sm remove-btn" @click="handleRemove(item.item || item)">
          取消收藏
        </button>
      </div>
    </div>

    <div v-if="removedCount > 0" class="removed-hint">
      有 {{ removedCount }} 件收藏物品已下架
    </div>

    <div v-if="!loading && !validFavorites.length" class="empty-state">
      <p>暂无收藏的物品</p>
      <p v-if="removedCount > 0" class="sub-hint">部分收藏物品已被发布者删除</p>
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

.removed-hint {
  margin-top: 16px;
  padding: 10px 14px;
  background: var(--gray-100);
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--gray-500);
  text-align: center;
}

.sub-hint {
  font-size: 13px;
  color: var(--gray-400);
  margin-top: 8px;
}
</style>

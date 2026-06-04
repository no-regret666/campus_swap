<script setup>
import { computed } from 'vue'
import { useUserStore } from '../stores/user'

const props = defineProps({
  item: { type: Object, required: true }
})

const userStore = useUserStore()

const isOwner = computed(() => {
  if (!userStore.user || !props.item) return false
  const itemOwnerId = props.item.ownerId || props.item.owner?.id || props.item.owner
  const userId = userStore.user.id || userStore.user._id
  return itemOwnerId === userId
})
</script>

<template>
  <router-link :to="`/item/${item._id || item.id}`" class="item-card card">
    <div class="card-img">
      <img v-if="item.images && item.images.length" :src="item.images[0]" :alt="item.title" />
      <div v-else class="card-img-placeholder">📦</div>
      <span v-if="isOwner" class="card-status badge badge-primary">我的物品</span>
      <span v-else class="card-status badge" :class="item.status === 'available' ? 'badge-success' : 'badge-warning'">
        {{ item.status === 'available' ? '可交换' : '已交换' }}
      </span>
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ item.title }}</h3>
      <p class="card-desc">{{ item.description }}</p>
      <div class="card-meta">
        <span class="card-campus">📍 {{ item.campus || '未知校区' }}</span>
        <span v-if="isOwner" class="card-mine-tag">🚫 不可交换</span>
        <span v-else class="card-category">{{ item.categoryName || item.category || '' }}</span>
      </div>
    </div>
  </router-link>
</template>

<style scoped>
.item-card {
  display: block;
  cursor: pointer;
}

.card-img {
  position: relative;
  height: 200px;
  background: var(--gray-100);
  overflow: hidden;
}

.card-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.item-card:hover .card-img img {
  transform: scale(1.05);
}

.card-img-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  background: var(--gray-100);
}

.card-status {
  position: absolute;
  top: 12px;
  right: 12px;
}

.card-owner {
  position: absolute;
  top: 12px;
  left: 12px;
}

.card-mine-tag {
  color: var(--danger);
  font-weight: 500;
  font-size: 12px;
}

.card-body {
  padding: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--gray-800);
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-desc {
  font-size: 13px;
  color: var(--gray-500);
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--gray-400);
}

.card-category {
  background: var(--gray-100);
  padding: 2px 8px;
  border-radius: 4px;
}
</style>

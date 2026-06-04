<script setup>
defineProps({
  exchange: { type: Object, required: true }
})
defineEmits(['accept', 'reject'])

function statusLabel(status) {
  const map = { pending: '待确认', accepted: '已接受', rejected: '已拒绝', completed: '已完成' }
  return map[status] || status
}
function statusClass(status) {
  const map = { pending: 'badge-warning', accepted: 'badge-success', rejected: 'badge-danger', completed: 'badge-primary' }
  return map[status] || 'badge-primary'
}
</script>

<template>
  <div class="exchange-card card">
    <div class="exchange-header">
      <span class="badge" :class="statusClass(exchange.status)">{{ statusLabel(exchange.status) }}</span>
      <span class="exchange-time">{{ new Date(exchange.createdAt).toLocaleDateString() }}</span>
    </div>
    <div class="exchange-body">
      <div class="exchange-info">
        <h4>{{ exchange.item?.title || '物品' }}</h4>
        <p>交换物品：{{ exchange.offeredItem }}</p>
        <p v-if="exchange.meetPlace">见面地点：{{ exchange.meetPlace }}</p>
        <p v-if="exchange.message" class="exchange-msg">留言：{{ exchange.message }}</p>
      </div>
    </div>
    <div class="exchange-actions" v-if="exchange.status === 'pending' && exchange.isOwner">
      <button class="btn btn-success btn-sm" @click="$emit('accept', exchange)">接受</button>
      <button class="btn btn-danger btn-sm" @click="$emit('reject', exchange)">拒绝</button>
    </div>
    <div class="exchange-actions">
      <router-link :to="`/messages/${exchange._id || exchange.id}`" class="btn btn-secondary btn-sm">
        💬 沟通
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.exchange-card {
  padding: 20px;
  margin-bottom: 16px;
}

.exchange-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.exchange-time {
  font-size: 12px;
  color: var(--gray-400);
}

.exchange-body {
  margin-bottom: 12px;
}

.exchange-info h4 {
  font-size: 16px;
  margin-bottom: 6px;
  color: var(--gray-800);
}

.exchange-info p {
  font-size: 13px;
  color: var(--gray-500);
  margin-bottom: 4px;
}

.exchange-msg {
  font-style: italic;
}

.exchange-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}
</style>

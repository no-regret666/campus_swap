<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}
</script>

<template>
  <button class="back-btn" @click="goBack" title="返回">
    <span class="back-arrow">←</span>
    <span class="back-text">返回</span>
  </button>
</template>

<style scoped>
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--gray-100);
  border: none;
  border-radius: var(--radius-sm);
  font-size: 14px;
  color: var(--gray-700);
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 20px;
}

.back-btn:hover {
  background: var(--gray-200);
  color: var(--gray-900);
}

.back-arrow {
  font-size: 16px;
  font-weight: 600;
}

@media (max-width: 768px) {
  .back-btn {
    padding: 6px 12px;
    font-size: 13px;
  }
  
  .back-text {
    display: none;
  }
  
  .back-btn {
    padding: 8px;
  }
}
</style>
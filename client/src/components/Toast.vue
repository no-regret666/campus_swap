<script setup>
import { useAppStore } from '../stores/app'

const appStore = useAppStore()
</script>

<template>
  <Transition name="toast">
    <div v-if="appStore.toast.show" class="toast" :class="'toast-' + appStore.toast.type">
      <span class="toast-icon">
        {{ appStore.toast.type === 'success' ? '✓' : appStore.toast.type === 'error' ? '✕' : 'ℹ' }}
      </span>
      <span>{{ appStore.toast.message }}</span>
    </div>
  </Transition>
</template>

<style scoped>
.toast {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  z-index: 9999;
  box-shadow: var(--shadow-lg);
}

.toast-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.toast-success {
  background: #ecfdf5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}
.toast-success .toast-icon { background: var(--success); color: white; }

.toast-error {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecaca;
}
.toast-error .toast-icon { background: var(--danger); color: white; }

.toast-info {
  background: #eff6ff;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}
.toast-info .toast-icon { background: var(--primary); color: white; }

.toast-enter-active, .toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}
</style>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getItem } from '../api/items'
import { createExchange } from '../api/exchanges'
import { addFavorite, removeFavorite } from '../api/favorites'
import { report } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const item = ref(null)
const loading = ref(false)
const isFavorited = ref(false)
const ownerName = ref('')

// 交换弹窗
const showExchangeModal = ref(false)
const exchangeForm = ref({ offeredItem: '', message: '' })

// 举报弹窗
const showReportModal = ref(false)
const reportReason = ref('')

const isOwner = computed(() => {
  if (!item.value) return true
  if (!userStore.user) return false // 用户信息未加载时不隐藏按钮
  const itemOwnerId = item.value.ownerId || item.value.owner?.id || item.value.owner
  const userId = userStore.user.id || userStore.user._id
  return itemOwnerId === userId
})

onMounted(async () => {
  loading.value = true
  try {
    const res = await getItem(route.params.id)
    item.value = res.item || res
    ownerName.value = res.ownerName || ''
    isFavorited.value = item.value.isFavorited || false
  } catch (e) {
    appStore.showToast('加载物品详情失败', 'error')
  } finally {
    loading.value = false
  }
})

async function handleExchange() {
  if (!exchangeForm.value.offeredItem) {
    appStore.showToast('请填写您要交换的物品描述', 'warning')
    return
  }
  try {
    await createExchange({
      itemId: item.value._id || item.value.id,
      offeredItem: exchangeForm.value.offeredItem,
      message: exchangeForm.value.message
    })
    appStore.showToast('交换申请已发送', 'success')
    showExchangeModal.value = false
    exchangeForm.value = { offeredItem: '', message: '' }
  } catch (e) {
    appStore.showToast(e.response?.data?.message || '申请失败', 'error')
  }
}

async function toggleFavorite() {
  if (!userStore.isLoggedIn) {
    appStore.showToast('请先登录', 'warning')
    return
  }
  try {
    const id = item.value._id || item.value.id
    if (isFavorited.value) {
      await removeFavorite(id)
      isFavorited.value = false
      appStore.showToast('已取消收藏', 'success')
    } else {
      await addFavorite(id)
      isFavorited.value = true
      appStore.showToast('已收藏', 'success')
    }
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleReport() {
  if (!reportReason.value) {
    appStore.showToast('请填写举报理由', 'warning')
    return
  }
  try {
    await report({
      targetType: 'item',
      targetId: item.value._id || item.value.id,
      reason: reportReason.value
    })
    appStore.showToast('举报已提交', 'success')
    showReportModal.value = false
    reportReason.value = ''
  } catch (e) {
    appStore.showToast('举报失败', 'error')
  }
}
</script>

<template>
  <div class="container">
    <div v-if="loading" class="empty-state">加载中...</div>

    <div v-else-if="item" class="detail-page">
      <!-- 图片 -->
      <div class="detail-image">
        <img v-if="item.images && item.images.length" :src="item.images[0]" :alt="item.title" />
        <div v-else class="img-placeholder">📦</div>
      </div>

      <!-- 信息 -->
      <div class="detail-info card">
        <h1 class="page-title">{{ item.title }}</h1>
        <div class="detail-meta">
          <span class="badge badge-primary">{{ item.categoryName || item.category }}</span>
          <span class="badge badge-success">{{ item.condition || '未标注' }}</span>
          <span>📍 {{ item.campus || '未知校区' }}</span>
          <span>👁 {{ item.views || 0 }} 次浏览</span>
          <span>📅 {{ new Date(item.createdAt).toLocaleDateString() }}</span>
        </div>
        <div class="detail-desc">
          <h3>物品描述</h3>
          <p>{{ item.description }}</p>
        </div>
        <div v-if="item.wantExchange" class="detail-desc">
          <h3>期望交换</h3>
          <p>{{ item.wantExchange }}</p>
        </div>

        <!-- 物主信息 -->
        <div class="owner-info" v-if="ownerName">
          <h3>发布者</h3>
          <p>{{ ownerName }}</p>
        </div>

        <!-- 操作按钮 -->
        <div class="detail-actions" v-if="userStore.isLoggedIn">
          <button class="btn" :class="isFavorited ? 'btn-warning' : 'btn-secondary'" @click="toggleFavorite">
            {{ isFavorited ? '❤️ 已收藏' : '🤍 收藏' }}
          </button>
          <button v-if="!isOwner" class="btn btn-primary" @click="showExchangeModal = true">
            🔄 申请交换
          </button>
          <button v-if="isOwner" class="btn btn-secondary" disabled>
            🚫 自己的物品，不可交换
          </button>
          <button class="btn btn-danger" @click="showReportModal = true">
            ⚠️ 举报
          </button>
        </div>
      </div>
    </div>

    <!-- 交换弹窗 -->
    <div v-if="showExchangeModal" class="modal-overlay" @click.self="showExchangeModal = false">
      <div class="modal card">
        <h3>申请交换</h3>
        <div class="form-group">
          <label>您要交换的物品描述</label>
          <input v-model="exchangeForm.offeredItem" class="form-input" placeholder="描述您要交换的物品" />
        </div>
        <div class="form-group">
          <label>留言（可选）</label>
          <textarea v-model="exchangeForm.message" class="form-input" rows="3" placeholder="给对方留言..."></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="showExchangeModal = false">取消</button>
          <button class="btn btn-primary" @click="handleExchange">提交申请</button>
        </div>
      </div>
    </div>

    <!-- 举报弹窗 -->
    <div v-if="showReportModal" class="modal-overlay" @click.self="showReportModal = false">
      <div class="modal card">
        <h3>举报物品</h3>
        <div class="form-group">
          <label>举报理由</label>
          <textarea v-model="reportReason" class="form-input" rows="3" placeholder="请描述举报原因..."></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="showReportModal = false">取消</button>
          <button class="btn btn-danger" @click="handleReport">提交举报</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-page {
  max-width: 800px;
  margin: 0 auto;
}

.detail-image {
  width: 100%;
  height: 350px;
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--gray-100);
  margin-bottom: 24px;
}

.detail-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.img-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 80px;
}

.detail-info {
  padding: 24px;
}

.detail-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
  font-size: 14px;
  color: var(--gray-500);
}

.detail-desc {
  margin-bottom: 20px;
}

.detail-desc h3 {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--gray-700);
}

.detail-desc p {
  color: var(--gray-600);
  line-height: 1.6;
}

.owner-info {
  margin-bottom: 20px;
  padding: 12px;
  background: var(--gray-50);
  border-radius: var(--radius-sm);
}

.owner-info h3 {
  font-size: 14px;
  color: var(--gray-500);
  margin-bottom: 4px;
}

.detail-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  padding-top: 16px;
  border-top: 1px solid var(--gray-100);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  width: 90%;
  max-width: 450px;
  padding: 24px;
}

.modal h3 {
  margin-bottom: 16px;
  font-size: 18px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

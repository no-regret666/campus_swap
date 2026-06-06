<script setup>
import { ref, computed, onMounted } from 'vue'
import { getExchanges, updateExchange } from '../api/exchanges'
import { rate } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

const userStore = useUserStore()
const appStore = useAppStore()

const exchanges = ref([])
const loading = ref(false)
const activeTab = ref('sent')

// 评价弹窗
const showRateModal = ref(false)
const rateForm = ref({ score: 5, comment: '' })
const ratingExchange = ref(null)

const sentExchanges = computed(() =>
  exchanges.value.filter(e => {
    const userId = userStore.user?.id
    return e.requesterId === userId
  })
)

const receivedExchanges = computed(() =>
  exchanges.value.filter(e => {
    const userId = userStore.user?.id
    return e.ownerId === userId
  })
)

onMounted(loadExchanges)

async function loadExchanges() {
  loading.value = true
  try {
    const res = await getExchanges()
    exchanges.value = res.exchanges || res || []
  } catch (e) {
    appStore.showToast('加载交换记录失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleAccept(exchange) {
  try {
    await updateExchange(exchange.id, { status: 'accepted' })
    appStore.showToast('已接受', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleReject(exchange) {
  try {
    await updateExchange(exchange.id, { status: 'rejected' })
    appStore.showToast('已拒绝', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleComplete(exchange) {
  try {
    await updateExchange(exchange.id, { status: 'completed' })
    appStore.showToast('交换已完成', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

function openRate(exchange) {
  ratingExchange.value = exchange
  rateForm.value = { score: 5, comment: '' }
  showRateModal.value = true
}

async function handleRate() {
  try {
    await rate({
      exchangeId: ratingExchange.value.id,
      score: rateForm.value.score,
      comment: rateForm.value.comment
    })
    appStore.showToast('评价成功', 'success')
    showRateModal.value = false
  } catch (e) {
    appStore.showToast('评价失败', 'error')
  }
}

function statusLabel(status) {
  const map = { pending: '⏳ 待确认', accepted: '✅ 已接受', rejected: '❌ 已拒绝', completed: '🎉 已完成', cancelled: '🚫 已取消' }
  return map[status] || status
}

function statusClass(status) {
  const map = { pending: 'status-pending', accepted: 'status-accepted', rejected: 'status-rejected', completed: 'status-completed' }
  return map[status] || ''
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">交换管理</h1>

    <div class="tabs">
      <button class="tab-btn" :class="{ active: activeTab === 'sent' }" @click="activeTab = 'sent'">
        我发起的 ({{ sentExchanges.length }})
      </button>
      <button class="tab-btn" :class="{ active: activeTab === 'received' }" @click="activeTab = 'received'">
        我收到的 ({{ receivedExchanges.length }})
      </button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

    <template v-else>
      <!-- 我发起的 -->
      <div v-if="activeTab === 'sent'">
        <div v-if="sentExchanges.length">
          <div v-for="ex in sentExchanges" :key="ex.id" class="exchange-card card">
            <div class="exchange-header">
              <span class="status-badge" :class="statusClass(ex.status)">{{ statusLabel(ex.status) }}</span>
              <span class="exchange-time">{{ new Date(ex.createdAt).toLocaleString() }}</span>
            </div>
            <div class="exchange-body">
              <div class="exchange-item-info">
                <img v-if="ex.itemImage" :src="ex.itemImage" class="exchange-item-img" />
                <div class="exchange-item-img placeholder" v-else>📦</div>
                <div class="exchange-detail">
                  <h4>{{ ex.itemTitle || '未知物品' }}</h4>
                  <p class="detail-row"><span class="label">物品所有者：</span>{{ ex.ownerName }}</p>
                  <p class="detail-row"><span class="label">我提供交换：</span>{{ ex.offerDesc || '未填写' }}</p>
                  <p class="detail-row" v-if="ex.message"><span class="label">留言：</span>{{ ex.message }}</p>
                  <p class="detail-row" v-if="ex.itemCategory"><span class="label">分类：</span>{{ ex.itemCategory }}</p>
                </div>
              </div>
            </div>
            <div class="exchange-actions">
              <router-link :to="`/messages/${ex.id}`" class="btn btn-secondary btn-sm">💬 沟通</router-link>
              <button v-if="ex.status === 'accepted'" class="btn btn-success btn-sm" @click="handleComplete(ex)">确认完成</button>
              <button v-if="ex.status === 'completed'" class="btn btn-primary btn-sm" @click="openRate(ex)">⭐ 评价</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">暂无发起的交换申请</div>
      </div>

      <!-- 我收到的 -->
      <div v-if="activeTab === 'received'">
        <div v-if="receivedExchanges.length">
          <div v-for="ex in receivedExchanges" :key="ex.id" class="exchange-card card">
            <div class="exchange-header">
              <span class="status-badge" :class="statusClass(ex.status)">{{ statusLabel(ex.status) }}</span>
              <span class="exchange-time">{{ new Date(ex.createdAt).toLocaleString() }}</span>
            </div>
            <div class="exchange-body">
              <div class="exchange-item-info">
                <img v-if="ex.itemImage" :src="ex.itemImage" class="exchange-item-img" />
                <div class="exchange-item-img placeholder" v-else>📦</div>
                <div class="exchange-detail">
                  <h4>{{ ex.itemTitle || '未知物品' }}</h4>
                  <p class="detail-row"><span class="label">申请人：</span>{{ ex.requesterName }}</p>
                  <p class="detail-row"><span class="label">对方提供交换：</span>{{ ex.offerDesc || '未填写' }}</p>
                  <p class="detail-row" v-if="ex.message"><span class="label">留言：</span>{{ ex.message }}</p>
                  <p class="detail-row" v-if="ex.itemCategory"><span class="label">分类：</span>{{ ex.itemCategory }}</p>
                </div>
              </div>
            </div>
            <div class="exchange-actions">
              <router-link :to="`/messages/${ex.id}`" class="btn btn-secondary btn-sm">💬 沟通</router-link>
              <template v-if="ex.status === 'pending'">
                <button class="btn btn-success btn-sm" @click="handleAccept(ex)">✓ 接受</button>
                <button class="btn btn-danger btn-sm" @click="handleReject(ex)">✗ 拒绝</button>
              </template>
              <button v-if="ex.status === 'accepted'" class="btn btn-success btn-sm" @click="handleComplete(ex)">确认完成</button>
              <button v-if="ex.status === 'completed'" class="btn btn-primary btn-sm" @click="openRate(ex)">⭐ 评价</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">暂无收到的交换申请</div>
      </div>
    </template>

    <!-- 评价弹窗 -->
    <div v-if="showRateModal" class="modal-overlay" @click.self="showRateModal = false">
      <div class="modal card">
        <h3>评价交换</h3>
        <div class="form-group">
          <label>评分（1-5）</label>
          <select v-model.number="rateForm.score" class="form-input">
            <option v-for="s in 5" :key="s" :value="s">{{ s }} 星</option>
          </select>
        </div>
        <div class="form-group">
          <label>评论</label>
          <textarea v-model="rateForm.comment" class="form-input" rows="3" placeholder="写一段评价..."></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="showRateModal = false">取消</button>
          <button class="btn btn-primary" @click="handleRate">提交评价</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tabs {
  display: flex;
  margin-bottom: 24px;
  border-bottom: 2px solid var(--gray-100);
}

.tab-btn {
  flex: 1;
  padding: 12px;
  border: none;
  background: none;
  font-size: 15px;
  font-weight: 500;
  color: var(--gray-500);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
}

.tab-btn.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

/* 交换卡片 */
.exchange-card {
  padding: 20px;
  margin-bottom: 16px;
}

.exchange-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.status-badge {
  font-size: 13px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 12px;
}

.status-pending { background: #fff3cd; color: #856404; }
.status-accepted { background: #d4edda; color: #155724; }
.status-rejected { background: #f8d7da; color: #721c24; }
.status-completed { background: #cce5ff; color: #004085; }

.exchange-time {
  font-size: 12px;
  color: var(--gray-400);
}

.exchange-item-info {
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.exchange-item-img {
  width: 72px;
  height: 72px;
  border-radius: 8px;
  object-fit: cover;
  flex-shrink: 0;
}

.exchange-item-img.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gray-100);
  font-size: 28px;
}

.exchange-detail {
  flex: 1;
}

.exchange-detail h4 {
  font-size: 16px;
  margin-bottom: 8px;
  color: var(--gray-800);
}

.detail-row {
  font-size: 13px;
  color: var(--gray-600);
  margin-bottom: 4px;
}

.detail-row .label {
  color: var(--gray-400);
  margin-right: 4px;
}

.exchange-actions {
  display: flex;
  gap: 8px;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--gray-100);
}

/* 弹窗 */
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
  max-width: 420px;
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

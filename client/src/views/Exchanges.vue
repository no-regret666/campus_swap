<script setup>
import { ref, computed, onMounted } from 'vue'
import { getExchanges, updateExchange, getOverdueExchanges } from '../api/exchanges'
import { rate } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import BackButton from '../components/BackButton.vue'

const userStore = useUserStore()
const appStore = useAppStore()

const exchanges = ref([])
const loading = ref(false)
const activeTab = ref('sent')

// 搜索和筛选
const searchKeyword = ref('')
const statusFilter = ref('')
const startDate = ref('')
const endDate = ref('')

// 逾期提醒
const overdueExchanges = ref([])

// meet时间/地点弹窗
const showMeetModal = ref(false)
const meetExchange = ref(null)
const meetForm = ref({ meetTime: '', meetPlace: '' })

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

onMounted(() => {
  loadExchanges()
  loadOverdue()
})

async function loadExchanges() {
  loading.value = true
  try {
    const params = {}
    if (searchKeyword.value) params.search = searchKeyword.value
    if (statusFilter.value) params.status = statusFilter.value
    if (startDate.value) params.startDate = startDate.value
    if (endDate.value) params.endDate = endDate.value
    const res = await getExchanges(params)
    exchanges.value = res.exchanges || res || []
  } catch (e) {
    appStore.showToast('加载交换记录失败', 'error')
  } finally {
    loading.value = false
  }
}

async function loadOverdue() {
  try {
    const res = await getOverdueExchanges()
    overdueExchanges.value = res.overdue || []
  } catch (e) {
    // 忽略错误
  }
}

function openMeetModal(exchange) {
  meetExchange.value = exchange
  meetForm.value = {
    meetTime: exchange.meetTime || '',
    meetPlace: exchange.meetPlace || ''
  }
  showMeetModal.value = true
}

async function handleSaveMeet() {
  try {
    await updateExchange(meetExchange.value.id, {
      meetTime: meetForm.value.meetTime,
      meetPlace: meetForm.value.meetPlace
    })
    appStore.showToast('见面信息已保存', 'success')
    showMeetModal.value = false
    loadExchanges()
  } catch (e) {
    appStore.showToast('保存失败', 'error')
  }
}

async function handleCancel(exchange) {
  if (!confirm('确定要取消这个交换申请吗？')) return
  try {
    await updateExchange(exchange.id, { status: 'cancelled' })
    appStore.showToast('交换已取消', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
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
  const map = {
    pending: '⏳ 待确认',
    accepted: '✅ 已接受',
    rejected: '❌ 已拒绝',
    completed: '🎉 已完成',
    cancelled: '🚫 已取消'
  }
  return map[status] || status
}

function statusClass(status) {
  const map = {
    pending: 'status-pending',
    accepted: 'status-accepted',
    rejected: 'status-rejected',
    completed: 'status-completed',
    cancelled: 'status-cancelled'
  }
  return map[status] || ''
}

function formatOverdue(hours) {
  if (hours >= 72) return `${Math.floor(hours/24)}天`
  return `${hours}小时`
}

function clearDateFilter() {
  startDate.value = ''
  endDate.value = ''
  loadExchanges()
}
</script>

<template>
  <div class="container">
    <BackButton />
    <h1 class="page-title">交换管理</h1>

    <!-- 逾期提醒 -->
    <div v-if="overdueExchanges.length > 0" class="overdue-warning card">
      <div class="overdue-header">
        <span class="overdue-icon">⚠️</span>
        <span>您有 {{ overdueExchanges.length }} 个交换待处理（超过24小时）</span>
      </div>
      <div v-for="ex in overdueExchanges" :key="ex.id" class="overdue-item">
        <span>{{ ex.itemTitle }}</span>
        <span class="overdue-time">已逾期 {{ formatOverdue(ex.overdueHours) }}</span>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-bar">
      <input
        v-model="searchKeyword"
        class="form-input search-input"
        placeholder="搜索物品名称..."
        @keyup.enter="loadExchanges"
      />
      <select v-model="statusFilter" class="form-input status-select" @change="loadExchanges">
        <option value="">全部状态</option>
        <option value="pending">待确认</option>
        <option value="accepted">已接受</option>
        <option value="completed">已完成</option>
        <option value="rejected">已拒绝</option>
        <option value="cancelled">已取消</option>
      </select>
      <button class="btn btn-secondary" @click="loadExchanges">搜索</button>
    </div>

    <!-- 日期筛选 -->
    <div class="date-filter">
      <span class="filter-label">筛选日期：</span>
      <input type="date" v-model="startDate" class="form-input date-input" @change="loadExchanges" />
      <span class="filter-separator">至</span>
      <input type="date" v-model="endDate" class="form-input date-input" @change="loadExchanges" />
      <button v-if="startDate || endDate" class="btn btn-outline btn-sm" @click="clearDateFilter">清除</button>
    </div>

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
                  <p class="detail-row" v-if="ex.meetTime || ex.meetPlace">
                    <span class="label">约定：</span>
                    {{ ex.meetTime || '' }}{{ ex.meetTime && ex.meetPlace ? ' @ ' : '' }}{{ ex.meetPlace || '' }}
                  </p>
                </div>
              </div>
            </div>
            <div class="exchange-actions">
              <router-link :to="`/messages/${ex.id}`" class="btn btn-secondary btn-sm">
                💬 沟通<span v-if="ex.unreadCount" class="unread-badge">{{ ex.unreadCount }}</span>
              </router-link>
              <button v-if="ex.status === 'pending'" class="btn btn-outline btn-sm" @click="handleCancel(ex)">取消</button>
              <button v-if="ex.status === 'accepted' && (ex.meetTime || ex.meetPlace)" class="btn btn-outline btn-sm" @click="openMeetModal(ex)">📍 见面信息</button>
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
                  <p class="detail-row" v-if="ex.meetTime || ex.meetPlace">
                    <span class="label">约定：</span>
                    {{ ex.meetTime || '' }}{{ ex.meetTime && ex.meetPlace ? ' @ ' : '' }}{{ ex.meetPlace || '' }}
                  </p>
                </div>
              </div>
            </div>
            <div class="exchange-actions">
              <router-link :to="`/messages/${ex.id}`" class="btn btn-secondary btn-sm">
                💬 沟通<span v-if="ex.unreadCount" class="unread-badge">{{ ex.unreadCount }}</span>
              </router-link>
              <template v-if="ex.status === 'pending'">
                <button class="btn btn-success btn-sm" @click="handleAccept(ex)">✓ 接受</button>
                <button class="btn btn-danger btn-sm" @click="handleReject(ex)">✗ 拒绝</button>
              </template>
              <button v-if="ex.status === 'accepted'" class="btn btn-outline btn-sm" @click="openMeetModal(ex)">📍 设置见面</button>
              <button v-if="ex.status === 'accepted'" class="btn btn-success btn-sm" @click="handleComplete(ex)">确认完成</button>
              <button v-if="ex.status === 'completed'" class="btn btn-primary btn-sm" @click="openRate(ex)">⭐ 评价</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">暂无收到的交换申请</div>
      </div>
    </template>

    <!-- 见面信息弹窗 -->
    <div v-if="showMeetModal" class="modal-overlay" @click.self="showMeetModal = false">
      <div class="modal card">
        <h3>设置见面信息</h3>
        <div class="form-group">
          <label>见面时间</label>
          <input v-model="meetForm.meetTime" class="form-input" placeholder="例如：明天 14:00" />
        </div>
        <div class="form-group">
          <label>见面地点</label>
          <input v-model="meetForm.meetPlace" class="form-input" placeholder="例如：图书馆门口" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="showMeetModal = false">取消</button>
          <button class="btn btn-primary" @click="handleSaveMeet">保存</button>
        </div>
      </div>
    </div>

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
/* 搜索栏 */
.search-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-input {
  flex: 1;
}

.status-select {
  width: 120px;
}

/* 日期筛选 */
.date-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 13px;
  color: var(--gray-500);
}

.filter-separator {
  font-size: 13px;
  color: var(--gray-400);
}

.date-input {
  width: 140px;
  padding: 8px 12px;
  font-size: 13px;
}

/* 逾期提醒 */
.overdue-warning {
  background: #fff3cd;
  border: 1px solid #ffc107;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.overdue-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #856404;
  margin-bottom: 8px;
}

.overdue-icon {
  font-size: 18px;
}

.overdue-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 13px;
  color: #856404;
}

.overdue-time {
  color: #d63384;
  font-weight: 600;
}

/* 未读角标 */
.unread-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: #dc3545;
  color: #fff;
  border-radius: 9px;
  font-size: 11px;
  margin-left: 4px;
}

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
.status-cancelled { background: var(--gray-200); color: var(--gray-600); }

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
  flex-wrap: wrap;
}

/* 按钮样式 */
.btn-outline {
  background: transparent;
  border: 1px solid var(--gray-300);
  color: var(--gray-700);
}

.btn-outline:hover {
  background: var(--gray-100);
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

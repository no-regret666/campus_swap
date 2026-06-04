<script setup>
import { ref, computed, onMounted } from 'vue'
import { getExchanges, updateExchange } from '../api/exchanges'
import { rate } from '../api/misc'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import ExchangeCard from '../components/ExchangeCard.vue'

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
    const requesterId = e.requester?._id || e.requester?.id || e.requester
    const userId = userStore.user?._id || userStore.user?.id
    return requesterId === userId
  })
)

const receivedExchanges = computed(() =>
  exchanges.value.filter(e => {
    const requesterId = e.requester?._id || e.requester?.id || e.requester
    const userId = userStore.user?._id || userStore.user?.id
    return requesterId !== userId
  }).map(e => ({ ...e, isOwner: true }))
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
    await updateExchange(exchange._id || exchange.id, { status: 'accepted' })
    appStore.showToast('已接受', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleReject(exchange) {
  try {
    await updateExchange(exchange._id || exchange.id, { status: 'rejected' })
    appStore.showToast('已拒绝', 'success')
    loadExchanges()
  } catch (e) {
    appStore.showToast('操作失败', 'error')
  }
}

async function handleComplete(exchange) {
  try {
    await updateExchange(exchange._id || exchange.id, { status: 'completed' })
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
      exchangeId: ratingExchange.value._id || ratingExchange.value.id,
      score: rateForm.value.score,
      comment: rateForm.value.comment
    })
    appStore.showToast('评价成功', 'success')
    showRateModal.value = false
  } catch (e) {
    appStore.showToast('评价失败', 'error')
  }
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">交换管理</h1>

    <div class="tabs">
      <button class="tab-btn" :class="{ active: activeTab === 'sent' }" @click="activeTab = 'sent'">我发起的</button>
      <button class="tab-btn" :class="{ active: activeTab === 'received' }" @click="activeTab = 'received'">我收到的</button>
    </div>

    <div v-if="loading" class="empty-state">加载中...</div>

    <template v-else>
      <!-- 我发起的 -->
      <div v-if="activeTab === 'sent'">
        <div v-if="sentExchanges.length">
          <div v-for="ex in sentExchanges" :key="ex._id || ex.id">
            <ExchangeCard :exchange="ex" />
            <div class="extra-actions" v-if="ex.status === 'accepted'">
              <button class="btn btn-success btn-sm" @click="handleComplete(ex)">确认完成</button>
            </div>
            <div class="extra-actions" v-if="ex.status === 'completed'">
              <button class="btn btn-primary btn-sm" @click="openRate(ex)">⭐ 评价</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">暂无发起的交换</div>
      </div>

      <!-- 我收到的 -->
      <div v-if="activeTab === 'received'">
        <div v-if="receivedExchanges.length">
          <div v-for="ex in receivedExchanges" :key="ex._id || ex.id">
            <ExchangeCard :exchange="ex" @accept="handleAccept" @reject="handleReject" />
            <div class="extra-actions" v-if="ex.status === 'accepted'">
              <button class="btn btn-success btn-sm" @click="handleComplete(ex)">确认完成</button>
            </div>
            <div class="extra-actions" v-if="ex.status === 'completed'">
              <button class="btn btn-primary btn-sm" @click="openRate(ex)">⭐ 评价</button>
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

.extra-actions {
  display: flex;
  gap: 8px;
  margin-top: -8px;
  margin-bottom: 16px;
  padding-left: 20px;
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

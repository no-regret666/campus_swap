<script setup>
import { ref, onMounted } from 'vue'
import { getDashboard, getReports, updateReport } from '../api/misc'
import { useAppStore } from '../stores/app'
import { useUserStore } from '../stores/user'

const appStore = useAppStore()
const userStore = useUserStore()

const dashboard = ref(null)
const reports = ref([])
const loading = ref(false)
const reportsLoading = ref(false)
const activeTab = ref('overview')

onMounted(loadData)

async function loadData() {
  loading.value = true
  reportsLoading.value = true
  try {
    const res = await getDashboard()
    dashboard.value = res.dashboard || res.stats || res
  } catch (e) {
    appStore.showToast('加载看板数据失败', 'error')
  } finally {
    loading.value = false
  }

  try {
    const res = await getReports()
    reports.value = res.reports || []
  } catch (e) {
    // 非管理员可能403，静默处理
  } finally {
    reportsLoading.value = false
  }
}

async function handleReport(reportItem, status) {
  const action = status === 'resolved' ? '处理' : '驳回'
  try {
    await updateReport(reportItem.id, { status })
    appStore.showToast(`举报已${action}`, 'success')
    // 刷新数据
    const res = await getReports()
    reports.value = res.reports || []
    // 刷新看板
    const dashRes = await getDashboard()
    dashboard.value = dashRes.dashboard || dashRes.stats || dashRes
  } catch (e) {
    appStore.showToast(`操作失败`, 'error')
  }
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  try {
    return new Date(timeStr).toLocaleString()
  } catch {
    return timeStr
  }
}

function reportTypeLabel(type) {
  return type === 'item' ? '物品举报' : '用户举报'
}

function reportStatusLabel(status) {
  const map = { pending: '⏳ 待处理', resolved: '✅ 已处理', dismissed: '❌ 已驳回' }
  return map[status] || status
}

function reportStatusClass(status) {
  const map = { pending: 'status-pending', resolved: 'status-resolved', dismissed: 'status-dismissed' }
  return map[status] || ''
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">管理后台</h1>

    <!-- 标签切换 -->
    <div class="tabs">
      <button class="tab-btn" :class="{ active: activeTab === 'overview' }" @click="activeTab = 'overview'">
        📊 数据概览
      </button>
      <button class="tab-btn" :class="{ active: activeTab === 'reports' }" @click="activeTab = 'reports'">
        ⚠️ 举报管理
        <span v-if="dashboard?.pendingReports" class="badge badge-danger" style="margin-left: 6px;">{{ dashboard.pendingReports }}</span>
      </button>
    </div>

    <!-- 数据概览 -->
    <div v-if="activeTab === 'overview'">
      <div v-if="loading" class="empty-state">加载中...</div>

      <template v-else-if="dashboard">
        <!-- 统计卡片 -->
        <div class="stats-grid">
          <div class="stat-card card">
            <div class="stat-icon">👥</div>
            <div class="stat-number">{{ dashboard.userCount ?? 0 }}</div>
            <div class="stat-title">用户总数</div>
          </div>
          <div class="stat-card card">
            <div class="stat-icon">📦</div>
            <div class="stat-number">{{ dashboard.itemCount ?? 0 }}</div>
            <div class="stat-title">物品总数</div>
          </div>
          <div class="stat-card card">
            <div class="stat-icon">🔄</div>
            <div class="stat-number">{{ dashboard.exchangeCount ?? 0 }}</div>
            <div class="stat-title">交换总数</div>
          </div>
          <div class="stat-card card">
            <div class="stat-icon">✅</div>
            <div class="stat-number">{{ dashboard.exchangeRate?.toFixed(1) ?? 0 }}%</div>
            <div class="stat-title">交换完成率</div>
          </div>
          <div class="stat-card card">
            <div class="stat-icon">⚠️</div>
            <div class="stat-number">{{ dashboard.pendingReports ?? 0 }}</div>
            <div class="stat-title">待处理举报</div>
          </div>
          <div class="stat-card card">
            <div class="stat-icon">💬</div>
            <div class="stat-number">{{ dashboard.messageCount ?? 0 }}</div>
            <div class="stat-title">消息总数</div>
          </div>
        </div>

        <!-- 热门分类排行 -->
        <div class="section card" v-if="dashboard.hotCategories && dashboard.hotCategories.length">
          <h3>热门分类排行</h3>
          <div class="rank-list">
            <div class="rank-item" v-for="(cat, index) in dashboard.hotCategories" :key="index">
              <span class="rank-num" :class="{ 'rank-top': index < 3 }">{{ index + 1 }}</span>
              <span class="rank-name">{{ cat.name || cat.category }}</span>
              <span class="rank-bar-wrapper">
                <span class="rank-bar" :style="{ width: Math.min(cat.count / (dashboard.hotCategories[0]?.count || 1) * 100, 100) + '%' }"></span>
              </span>
              <span class="rank-count badge badge-primary">{{ cat.count }} 件</span>
            </div>
          </div>
        </div>

        <!-- 最近审计日志 -->
        <div class="section card" v-if="dashboard.auditLogs && dashboard.auditLogs.length">
          <h3>最近审计日志</h3>
          <div class="log-list">
            <div class="log-item" v-for="log in dashboard.auditLogs" :key="log._id || log.id">
              <span class="log-action badge badge-warning">{{ log.action }}</span>
              <span class="log-desc">{{ log.description || log.detail }}</span>
              <span class="log-time">{{ formatTime(log.createdAt) }}</span>
            </div>
          </div>
        </div>

        <div v-if="(!dashboard.hotCategories || !dashboard.hotCategories.length) && (!dashboard.auditLogs || !dashboard.auditLogs.length)" class="empty-state">
          暂无更多运营数据
        </div>
      </template>

      <div v-else class="empty-state">暂无数据</div>
    </div>

    <!-- 举报管理 -->
    <div v-if="activeTab === 'reports'">
      <div v-if="reportsLoading" class="empty-state">加载中...</div>

      <template v-else-if="reports.length">
        <div v-for="rpt in reports" :key="rpt.id" class="report-card card">
          <div class="report-header">
            <div class="report-meta">
              <span class="report-type badge" :class="rpt.type === 'item' ? 'badge-warning' : 'badge-danger'">
                {{ reportTypeLabel(rpt.type) }}
              </span>
              <span class="report-status" :class="reportStatusClass(rpt.status)">
                {{ reportStatusLabel(rpt.status) }}
              </span>
            </div>
            <span class="report-time">{{ formatTime(rpt.createdAt) }}</span>
          </div>
          <div class="report-body">
            <p class="report-reason"><strong>举报原因：</strong>{{ rpt.reason }}</p>
            <p class="report-target"><strong>目标ID：</strong>{{ rpt.targetId }}</p>
          </div>
          <div class="report-actions" v-if="rpt.status === 'pending'">
            <button class="btn btn-success btn-sm" @click="handleReport(rpt, 'resolved')">✓ 处理（下架目标）</button>
            <button class="btn btn-secondary btn-sm" @click="handleReport(rpt, 'dismissed')">✗ 驳回</button>
          </div>
        </div>
      </template>

      <div v-else class="empty-state">暂无举报记录</div>
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
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.tab-btn.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  padding: 24px 16px;
  text-align: center;
}

.stat-icon {
  font-size: 28px;
  margin-bottom: 8px;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: var(--primary);
}

.stat-title {
  font-size: 13px;
  color: var(--gray-500);
  margin-top: 4px;
}

.section {
  padding: 24px;
  margin-bottom: 24px;
}

.section h3 {
  font-size: 16px;
  margin-bottom: 16px;
  color: var(--gray-800);
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rank-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--gray-50);
}

.rank-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--gray-100);
  color: var(--gray-500);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.rank-num.rank-top {
  background: var(--primary);
  color: white;
}

.rank-name {
  width: 80px;
  font-size: 14px;
  flex-shrink: 0;
}

.rank-bar-wrapper {
  flex: 1;
  height: 8px;
  background: var(--gray-100);
  border-radius: 4px;
  overflow: hidden;
}

.rank-bar {
  height: 100%;
  background: var(--primary);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.rank-count {
  flex-shrink: 0;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--gray-50);
  font-size: 13px;
}

.log-desc {
  flex: 1;
  color: var(--gray-600);
}

.log-time {
  font-size: 12px;
  color: var(--gray-400);
  flex-shrink: 0;
}

/* 举报管理 */
.report-card {
  padding: 20px;
  margin-bottom: 16px;
}

.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.report-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.report-status {
  font-size: 13px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 12px;
}

.status-pending { background: #fff3cd; color: #856404; }
.status-resolved { background: #d4edda; color: #155724; }
.status-dismissed { background: #f8d7da; color: #721c24; }

.report-time {
  font-size: 12px;
  color: var(--gray-400);
}

.report-body {
  margin-bottom: 12px;
}

.report-reason, .report-target {
  font-size: 14px;
  color: var(--gray-700);
  margin-bottom: 4px;
}

.report-target {
  font-size: 12px;
  color: var(--gray-400);
}

.report-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid var(--gray-100);
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .rank-name {
    width: 60px;
  }
}
</style>

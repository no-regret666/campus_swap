<script setup>
import { ref, onMounted } from 'vue'
import { getDashboard } from '../api/misc'
import { useAppStore } from '../stores/app'
import BackButton from '../components/BackButton.vue'

const appStore = useAppStore()

const dashboard = ref(null)
const loading = ref(false)

onMounted(loadDashboard)

async function loadDashboard() {
  loading.value = true
  try {
    const res = await getDashboard()
    // 后端返回 { stats: {...}, dashboard: {...} }，取 stats 数据
    dashboard.value = res.stats || res.dashboard || res
  } catch (e) {
    appStore.showToast('加载看板数据失败', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container">
    <BackButton />
    <h1 class="page-title">管理后台</h1>

    <div v-if="loading" class="empty-state">加载中...</div>

    <template v-else-if="dashboard">
      <!-- 统计卡片 -->
      <div class="stats-grid">
        <div class="stat-card card">
          <div class="stat-icon">👥</div>
          <div class="stat-number">{{ dashboard.userCount ?? dashboard.users ?? 0 }}</div>
          <div class="stat-title">用户总数</div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">📦</div>
          <div class="stat-number">{{ dashboard.itemCount ?? dashboard.items ?? 0 }}</div>
          <div class="stat-title">物品总数</div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">🔄</div>
          <div class="stat-number">{{ dashboard.exchangeCount ?? dashboard.exchanges ?? 0 }}</div>
          <div class="stat-title">交换总数</div>
        </div>
        <div class="stat-card card">
          <div class="stat-icon">⚠️</div>
          <div class="stat-number">{{ dashboard.reportCount ?? dashboard.pendingReports ?? 0 }}</div>
          <div class="stat-title">待处理举报</div>
        </div>
      </div>

      <!-- 热门分类 -->
      <div class="section card" v-if="dashboard.hotCategories && dashboard.hotCategories.length">
        <h3>热门分类排行</h3>
        <div class="rank-list">
          <div class="rank-item" v-for="(cat, index) in dashboard.hotCategories" :key="index">
            <span class="rank-num">{{ index + 1 }}</span>
            <span class="rank-name">{{ cat.name || cat.category }}</span>
            <span class="rank-count badge badge-primary">{{ cat.count }} 件</span>
          </div>
        </div>
      </div>

      <!-- 审计日志 -->
      <div class="section card" v-if="dashboard.auditLogs && dashboard.auditLogs.length">
        <h3>最近审计日志</h3>
        <div class="log-list">
          <div class="log-item" v-for="log in dashboard.auditLogs" :key="log._id || log.id">
            <span class="log-action badge badge-warning">{{ log.action }}</span>
            <span class="log-desc">{{ log.description || log.detail }}</span>
            <span class="log-time">{{ new Date(log.createdAt).toLocaleString() }}</span>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="empty-state">暂无数据</div>
  </div>
</template>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  padding: 24px;
  text-align: center;
}

.stat-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: var(--primary);
}

.stat-title {
  font-size: 14px;
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
  gap: 10px;
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
  background: var(--primary-light);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.rank-name {
  flex: 1;
  font-size: 14px;
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
}
</style>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login, register } from '../api/auth'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const activeTab = ref('login')
const loading = ref(false)

const loginForm = ref({ username: '', password: '' })
const registerForm = ref({ username: '', password: '', nickname: '', campus: '南湖校区' })

const campusOptions = ['南湖校区', '东湖校区', '线上']

async function handleLogin() {
  if (!loginForm.value.username || !loginForm.value.password) {
    appStore.showToast('请填写用户名和密码', 'warning')
    return
  }
  loading.value = true
  try {
    const res = await login(loginForm.value)
    userStore.setLogin(res)
    appStore.showToast('登录成功', 'success')
    router.push('/')
  } catch (e) {
    appStore.showToast(e.response?.data?.message || '登录失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!registerForm.value.username || !registerForm.value.password || !registerForm.value.nickname) {
    appStore.showToast('请填写完整信息', 'warning')
    return
  }
  loading.value = true
  try {
    const res = await register(registerForm.value)
    userStore.setLogin(res)
    appStore.showToast('注册成功', 'success')
    router.push('/')
  } catch (e) {
    appStore.showToast(e.response?.data?.message || '注册失败', 'error')
  } finally {
    loading.value = false
  }
}

function quickLogin(username, password) {
  loginForm.value = { username, password }
  handleLogin()
}
</script>

<template>
  <div class="container login-page">
    <div class="login-card card">
      <h1 class="page-title" style="text-align:center;">校园闲置交换平台</h1>

      <div class="tabs">
        <button class="tab-btn" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">登录</button>
        <button class="tab-btn" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">注册</button>
      </div>

      <!-- 登录表单 -->
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="loginForm.username" class="form-input" placeholder="请输入用户名" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="loginForm.password" type="password" class="form-input" placeholder="请输入密码" />
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <!-- 注册表单 -->
      <form v-else @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label>用户名（学号）</label>
          <input v-model="registerForm.username" class="form-input" placeholder="请输入学号" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="registerForm.password" type="password" class="form-input" placeholder="请输入密码" />
        </div>
        <div class="form-group">
          <label>昵称</label>
          <input v-model="registerForm.nickname" class="form-input" placeholder="请输入昵称" />
        </div>
        <div class="form-group">
          <label>校区</label>
          <select v-model="registerForm.campus" class="form-input">
            <option v-for="c in campusOptions" :key="c" :value="c">{{ c }}</option>
          </select>
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <!-- 快速登录 -->
      <div class="quick-login">
        <p>快速体验（演示账号）：</p>
        <div class="quick-btns">
          <button class="btn btn-secondary btn-sm" @click="quickLogin('202401001', '123456')">学生1</button>
          <button class="btn btn-secondary btn-sm" @click="quickLogin('202401002', '123456')">学生2</button>
          <button class="btn btn-secondary btn-sm" @click="quickLogin('admin', 'admin123')">管理员</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
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

.login-form {
  margin-bottom: 24px;
}

.quick-login {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid var(--gray-100);
}

.quick-login p {
  font-size: 13px;
  color: var(--gray-500);
  margin-bottom: 12px;
}

.quick-btns {
  display: flex;
  gap: 8px;
  justify-content: center;
}
</style>

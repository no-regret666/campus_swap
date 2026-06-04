import axios from 'axios'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

const http = axios.create({
  baseURL: '/api',
  timeout: 10000
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  res => res.data,
  err => {
    const msg = err.response?.data?.message || '请求失败'
    const appStore = useAppStore()
    appStore.showToast(msg, 'error')
    if (err.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
    }
    return Promise.reject(err)
  }
)

export default http

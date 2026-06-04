import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// 如果有 token，自动加载用户信息
import { useUserStore } from './stores/user'
const userStore = useUserStore()
if (userStore.token) {
  userStore.fetchUser()
}

app.mount('#app')

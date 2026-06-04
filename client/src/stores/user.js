import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getMe } from '../api/auth'
import router from '../router'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function fetchUser() {
    if (!token.value) return
    try {
      const res = await getMe()
      user.value = res.user || res
    } catch {
      logout()
    }
  }

  function setLogin(data) {
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    router.push('/login')
  }

  return { user, token, isLoggedIn, isAdmin, fetchUser, setLogin, logout }
})

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories, getNotifications } from '../api/misc'

export const useAppStore = defineStore('app', () => {
  const categories = ref([])
  const notifications = ref([])
  const toast = ref({ show: false, message: '', type: 'info' })

  let toastTimer = null

  function showToast(message, type = 'info') {
    toast.value = { show: true, message, type }
    if (toastTimer) clearTimeout(toastTimer)
    toastTimer = setTimeout(() => {
      toast.value.show = false
    }, 3000)
  }

  async function fetchCategories() {
    try {
      const res = await getCategories()
      categories.value = res.categories || res || []
    } catch {
      // ignore
    }
  }

  async function fetchNotifications() {
    try {
      const res = await getNotifications()
      notifications.value = res.notifications || res || []
    } catch {
      // ignore
    }
  }

  return { categories, notifications, toast, showToast, fetchCategories, fetchNotifications }
})

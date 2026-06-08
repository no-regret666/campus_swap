import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { guest: true } },
  { path: '/', name: 'Home', component: () => import('../views/Home.vue') },
  { path: '/item/:id', name: 'ItemDetail', component: () => import('../views/ItemDetail.vue') },
  { path: '/publish', name: 'Publish', component: () => import('../views/Publish.vue'), meta: { auth: true } },
  { path: '/exchanges', name: 'Exchanges', component: () => import('../views/Exchanges.vue'), meta: { auth: true } },
  { path: '/messages/:exchangeId', name: 'Messages', component: () => import('../views/Messages.vue'), meta: { auth: true } },
  { path: '/favorites', name: 'Favorites', component: () => import('../views/Favorites.vue'), meta: { auth: true } },
  { path: '/profile', name: 'Profile', component: () => import('../views/Profile.vue'), meta: { auth: true } },
  { path: '/admin', name: 'Admin', component: () => import('../views/Admin.vue'), meta: { auth: true, admin: true } },
  { path: '/recommendations', name: 'Recommend', component: () => import('../views/Recommend.vue'), meta: { auth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.auth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.admin && !userStore.isAdmin) {
    next('/')
  } else if (to.meta.guest && userStore.isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router

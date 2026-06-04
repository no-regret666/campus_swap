import http from './index'

export const report = (data) => http.post('/reports', data)
export const rate = (data) => http.post('/ratings', data)
export const getDashboard = () => http.get('/dashboard')
export const getRecommendations = () => http.get('/recommendations')
export const getCategories = () => http.get('/categories')
export const getNotifications = () => http.get('/notifications')
export const markNotificationRead = (id) => http.patch(`/notifications/${id}/read`)

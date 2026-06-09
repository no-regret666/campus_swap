import http from './index'

export const getNotifications = () => http.get('/notifications')
export const markNotificationRead = (id) => http.patch(`/notifications/${id}/read`)
export const markAllNotificationsRead = () => http.post('/notifications/read-all')
export const getUnreadNotificationCount = () => http.get('/notifications/unread-count')
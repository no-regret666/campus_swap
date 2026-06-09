import http from './index'

export const getMessages = (exchangeId) => http.get(`/messages/${exchangeId}`)
export const sendMessage = (data) => http.post('/messages', data)
export const markMessagesRead = (exchangeId) => http.post(`/messages/${exchangeId}/read`)
export const deleteMessage = (messageId) => http.delete(`/messages/${messageId}`)
export const getUnreadCount = () => http.get('/messages/unread/count')
export const getAllUnreadByExchange = () => http.get('/messages/unread/by-exchange')

// 上传图片
export const uploadImage = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return http.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

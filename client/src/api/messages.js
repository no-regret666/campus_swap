import http from './index'

export const getMessages = (exchangeId) => http.get(`/messages/${exchangeId}`)
export const sendMessage = (data) => http.post('/messages', data)

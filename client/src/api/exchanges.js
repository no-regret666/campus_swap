import http from './index'

export const getExchanges = (params) => http.get('/exchanges', { params })
export const createExchange = (data) => http.post('/exchanges', data)
export const updateExchange = (id, data) => http.patch(`/exchanges/${id}`, data)
export const getOverdueExchanges = () => http.get('/exchanges/overdue')

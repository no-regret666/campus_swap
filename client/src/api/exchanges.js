import http from './index'

export const getExchanges = () => http.get('/exchanges')
export const createExchange = (data) => http.post('/exchanges', data)
export const updateExchange = (id, data) => http.patch(`/exchanges/${id}`, data)

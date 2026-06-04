import http from './index'

export const getItems = (params) => http.get('/items', { params })
export const getItem = (id) => http.get(`/items/${id}`)
export const createItem = (data) => http.post('/items', data)

import http from './index'

export const getItems = (params) => http.get('/items', { params })
export const getItem = (id) => http.get(`/items/${id}`)
export const createItem = (data) => http.post('/items', data)
export const updateItem = (id, data) => http.put(`/items/${id}`, data)
export const deleteItem = (id) => http.delete(`/items/${id}`)

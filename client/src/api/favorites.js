import http from './index'

export const getFavorites = () => http.get('/favorites')
export const addFavorite = (itemId) => http.post('/favorites', { itemId })
export const removeFavorite = (itemId) => http.delete(`/favorites/${itemId}`)

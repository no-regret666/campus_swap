import http from './index'

// 上传图片到对象存储
export const uploadImage = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  // http拦截器已经解包res.data，这里直接返回
  const data = await http.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return data
}

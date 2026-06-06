<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createItem } from '../api/items'
import { uploadImage } from '../api/upload'
import { useAppStore } from '../stores/app'

const router = useRouter()
const appStore = useAppStore()

const loading = ref(false)
const uploading = ref(false)
const imageList = ref([]) // 已上传图片列表 [{url, filename}]
const form = ref({
  title: '',
  description: '',
  category: '',
  campus: '南湖校区',
  condition: '九成新',
  wantedExchange: ''
})

const campusOptions = ['南湖校区', '东湖校区', '线上']
const conditionOptions = ['全新', '九成新', '八成新', '七成新', '六成新及以下']

onMounted(() => {
  if (!appStore.categories.length) appStore.fetchCategories()
})

// 选择本地图片并上传
async function handleImageSelect(e) {
  const files = e.target.files
  if (!files || !files.length) return

  for (const file of files) {
    // 前端校验
    if (file.size > 5 * 1024 * 1024) {
      appStore.showToast(`${file.name} 超过5MB限制`, 'warning')
      continue
    }
    const ext = file.name.split('.').pop().toLowerCase()
    if (!['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)) {
      appStore.showToast(`${file.name} 格式不支持`, 'warning')
      continue
    }

    // 上传到服务器（对象存储）
    uploading.value = true
    try {
      const res = await uploadImage(file)
      imageList.value.push({
        url: res.url,
        filename: res.filename,
        // 本地预览
        preview: URL.createObjectURL(file)
      })
    } catch (err) {
      appStore.showToast(`${file.name} 上传失败`, 'error')
    }
  }
  uploading.value = false
  // 清空input，允许重复选择相同文件
  e.target.value = ''
}

// 移除已上传图片
function removeImage(index) {
  const img = imageList.value[index]
  if (img.preview) URL.revokeObjectURL(img.preview)
  imageList.value.splice(index, 1)
}

async function handleSubmit() {
  if (!form.value.title || !form.value.description || !form.value.category) {
    appStore.showToast('请填写标题、描述和分类', 'warning')
    return
  }
  loading.value = true
  try {
    const data = {
      ...form.value,
      images: imageList.value.map(img => img.url)
    }
    if (!data.images.length) delete data.images
    await createItem(data)
    appStore.showToast('发布成功', 'success')
    router.push('/')
  } catch (e) {
    appStore.showToast(e.response?.data?.message || '发布失败', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container">
    <h1 class="page-title">发布闲置物品</h1>

    <div class="publish-card card">
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>物品标题 *</label>
          <input v-model="form.title" class="form-input" placeholder="请输入物品标题" />
        </div>

        <div class="form-group">
          <label>物品描述 *</label>
          <textarea v-model="form.description" class="form-input" rows="4" placeholder="详细描述物品信息..."></textarea>
        </div>

        <div class="form-group">
          <label>分类 *</label>
          <select v-model="form.category" class="form-input">
            <option value="">请选择分类</option>
            <option v-for="cat in appStore.categories" :key="cat._id || cat.id || cat.name" :value="cat._id || cat.id || cat.name">
              {{ cat.name || cat }}
            </option>
          </select>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>校区</label>
            <select v-model="form.campus" class="form-input">
              <option v-for="c in campusOptions" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>

          <div class="form-group">
            <label>成色</label>
            <select v-model="form.condition" class="form-input">
              <option v-for="c in conditionOptions" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>
        </div>

        <div class="form-group">
          <label>期望交换物品</label>
          <input v-model="form.wantedExchange" class="form-input" placeholder="描述您期望交换的物品（可选）" />
        </div>

        <!-- 图片上传区域 -->
        <div class="form-group">
          <label>物品图片（最多3张，支持jpg/png/gif/webp，单张不超过5MB）</label>
          
          <div class="image-upload-area">
            <!-- 已上传的图片预览 -->
            <div v-for="(img, index) in imageList" :key="index" class="image-preview-item">
              <img :src="img.preview || img.url" alt="物品图片" />
              <button type="button" class="image-remove-btn" @click="removeImage(index)">×</button>
            </div>

            <!-- 添加图片按钮 -->
            <label v-if="imageList.length < 3" class="image-add-btn" :class="{ disabled: uploading }">
              <input
                type="file"
                accept="image/jpeg,image/png,image/gif,image/webp"
                multiple
                :disabled="uploading"
                @change="handleImageSelect"
                style="display:none"
              />
              <span v-if="uploading" class="upload-loading">上传中...</span>
              <span v-else class="upload-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
                <br/>选择图片
              </span>
            </label>
          </div>
        </div>

        <button type="submit" class="btn btn-primary btn-block" :disabled="loading || uploading">
          {{ loading ? '发布中...' : '发布物品' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.publish-card {
  max-width: 600px;
  margin: 0 auto;
  padding: 32px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .form-group {
  flex: 1;
}

/* 图片上传区域 */
.image-upload-area {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 8px;
}

.image-preview-item {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e0e0e0;
}

.image-preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-remove-btn {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  border: none;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-add-btn {
  width: 100px;
  height: 100px;
  border: 2px dashed #ccc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #999;
  font-size: 12px;
  text-align: center;
  transition: border-color 0.2s, color 0.2s;
}

.image-add-btn:hover {
  border-color: #4caf50;
  color: #4caf50;
}

.image-add-btn.disabled {
  pointer-events: none;
  opacity: 0.5;
}

.upload-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.upload-loading {
  font-size: 12px;
  color: #666;
}
</style>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createItem } from '../api/items'
import { useAppStore } from '../stores/app'

const router = useRouter()
const appStore = useAppStore()

const loading = ref(false)
const form = ref({
  title: '',
  description: '',
  category: '',
  campus: '南湖校区',
  condition: '九成新',
  wantedExchange: '',
  images: ''
})

const campusOptions = ['南湖校区', '东湖校区', '线上']
const conditionOptions = ['全新', '九成新', '八成新', '七成新', '六成新及以下']

onMounted(() => {
  if (!appStore.categories.length) appStore.fetchCategories()
})

async function handleSubmit() {
  if (!form.value.title || !form.value.description || !form.value.category) {
    appStore.showToast('请填写标题、描述和分类', 'warning')
    return
  }
  loading.value = true
  try {
    const data = { ...form.value }
    if (data.images) {
      data.images = [data.images]
    } else {
      delete data.images
    }
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

        <div class="form-group">
          <label>图片URL（可选）</label>
          <input v-model="form.images" class="form-input" placeholder="输入图片链接" />
        </div>

        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
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
</style>

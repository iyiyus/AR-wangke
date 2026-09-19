<template>
  <div class="taowa-notice-page">
    <div class="art-card p-5">
      <div class="flex-cb mb-4">
        <h4 class="font-bold">源台公告</h4>
        <ElButton size="small" :loading="loading" @click="load">刷新</ElButton>
      </div>
      <div v-if="loading" class="py-10 text-center text-g-500">加载中...</div>
      <div v-else-if="html" class="notice-body" v-html="html"></div>
      <div v-else class="py-10 text-center text-g-500">
        <p class="mb-2">{{ text }}</p>
        <p class="text-sm">（源台公告加载失败或未配置，展示站内公告）</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { getTaowaNotice } from '@/api/wk'

  defineOptions({ name: 'WkTaowaNotice' })

  const loading = ref(false)
  const html = ref('')
  const text = ref('')

  const load = async () => {
    loading.value = true
    try {
      const res: any = await getTaowaNotice()
      if (typeof res === 'string') {
        html.value = res
      } else {
        const n = res?.notice || ''
        if (n) html.value = n.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/\n/g, '<br>')
        else text.value = '暂无公告'
      }
    } catch (e) {
      text.value = '暂无公告'
    } finally {
      loading.value = false
    }
  }

  onMounted(() => load())
</script>

<style scoped>
  .notice-body :deep(p) { margin-bottom: 0.5rem; }
  .notice-body :deep(img) { max-width: 100%; }
</style>

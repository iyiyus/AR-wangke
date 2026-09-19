<template>
  <div class="art-card h-128 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>消费动态</h4>
        <p>近 8 条记录</p>
      </div>
    </div>
    <div class="h-9/10 mt-2 overflow-hidden">
      <ElScrollbar>
        <div
          class="h-17.5 leading-17.5 border-b border-g-300 text-sm overflow-hidden last:border-b-0"
          v-for="(item, index) in list"
          :key="index"
        >
          <span class="text-g-800 font-medium">{{ item.type }}</span>
          <span class="mx-2 text-g-600 text-ellipsis overflow-hidden whitespace-nowrap">{{ item.text }}</span>
          <span class="text-theme">{{ item.money }}</span>
        </div>
      </ElScrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { getAdminLog } from '@/api/wk'
  import { useUserStore } from '@/store/modules/user'

  const list = ref<any[]>([])

  onMounted(async () => {
    const uid = (useUserStore().info as any)?.uid
    try {
      const res: any = await getAdminLog({ current: 1, size: 8, uid: String(uid) })
      list.value = res.list || []
    } catch {}
  })
</script>

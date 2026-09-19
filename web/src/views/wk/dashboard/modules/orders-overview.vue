<template>
  <div class="art-card h-105 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>近 7 天订单数</h4>
        <p>合计 <span class="text-success ml-1">{{ total }}</span> 单</p>
      </div>
    </div>
    <ArtLineChart
      height="calc(100% - 56px)"
      :data="data"
      :xAxisData="xAxisData"
      :showAreaColor="true"
      :showAxisLine="false"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import request from '@/utils/http'

  const data = ref<number[]>([])
  const xAxisData = ref<string[]>([])
  const total = ref(0)

  onMounted(async () => {
    const res: any = await request.get({ url: '/api/order/monthly' })
    xAxisData.value = res.labels
    data.value = res.values
    total.value = res.values.reduce((s: number, v: number) => s + v, 0)
  })
</script>

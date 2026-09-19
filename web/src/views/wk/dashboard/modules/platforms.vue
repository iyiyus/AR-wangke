<template>
  <div class="art-card h-105 p-4 box-border mb-5 max-sm:mb-4">
    <ArtBarChart
      v-if="chartData.length"
      class="box-border p-2"
      barWidth="50%"
      height="13.7rem"
      :showAxisLine="false"
      :data="chartData"
      :xAxisData="xAxisLabels"
    />
    <div v-else class="h-[13.7rem] flex-cc text-g-500 text-sm">暂无订单数据</div>
    <div class="ml-1">
      <h3 class="mt-5 text-lg font-medium">平台概览</h3>
      <p class="mt-1 text-sm">在售平台 <span class="text-success font-medium">{{ activeCount }}</span></p>
      <p class="mt-1 text-sm">展示订单量前 9 名网课平台</p>
    </div>
    <div class="flex-b mt-2">
      <div class="flex-1" v-for="(item, index) in list" :key="index">
        <p class="text-2xl text-g-900">{{ item.num }}</p>
        <p class="text-xs text-g-500">{{ item.name }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { getClassAll } from '@/api/wk'
  import request from '@/utils/http'

  const xAxisLabels = ref<string[]>([])
  const chartData = ref<number[]>([])
  const activeCount = ref(0)

  const list = ref([
    { name: '总平台数', num: 0 as number | string },
    { name: '已上架', num: 0 as number | string },
    { name: '已下架', num: 0 as number | string },
    { name: '平均价格', num: '¥0' as number | string }
  ])

  onMounted(async () => {
    const [classes, byPlat]: [any[], any] = await Promise.all([
      getClassAll(),
      request.get({ url: '/api/order/by-platform' })
    ])
    xAxisLabels.value = byPlat.labels || []
    chartData.value = byPlat.values || []
    const total = classes.length
    const active = classes.filter((c) => c.status === 1).length
    activeCount.value = active
    const avg = total ? classes.reduce((s, c) => s + Number(c.price), 0) / total : 0
    list.value = [
      { name: '总平台数', num: total },
      { name: '已上架', num: active },
      { name: '已下架', num: total - active },
      { name: '平均价格', num: `¥${avg.toFixed(1)}` }
    ]
  })
</script>

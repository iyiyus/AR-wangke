<template>
  <div>
    <!-- 用户数据 -->
    <ElRow :gutter="20" class="flex">
      <ElCol :sm="12" :md="6">
        <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
          <span class="text-g-700 text-sm">总用户数</span>
          <ArtCountTo class="text-[26px] font-medium mt-2" :target="data.users?.total || 0" :duration="1300" />
          <div class="flex-c mt-1">
            <span class="text-xs text-g-600">活跃</span>
            <span class="ml-1 text-xs font-semibold text-success">{{ data.users?.active || 0 }}</span>
          </div>
          <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
            <ArtSvgIcon icon="ri:group-line" class="text-xl text-theme" />
          </div>
        </div>
      </ElCol>
      <ElCol :sm="12" :md="6">
        <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
          <span class="text-g-700 text-sm">总订单数</span>
          <ArtCountTo class="text-[26px] font-medium mt-2" :target="data.orders?.total || 0" :duration="1300" />
          <div class="flex-c mt-1">
            <span class="text-xs text-g-600">今日</span>
            <span class="ml-1 text-xs font-semibold text-success">+{{ data.orders?.today || 0 }}</span>
          </div>
          <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
            <ArtSvgIcon icon="ri:list-check-2" class="text-xl text-theme" />
          </div>
        </div>
      </ElCol>
      <ElCol :sm="12" :md="6">
        <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
          <span class="text-g-700 text-sm">总充值</span>
          <ArtCountTo class="text-[26px] font-medium mt-2" prefix="¥ " :decimals="2" :target="Number(data.revenue?.total) || 0" :duration="1300" />
          <div class="flex-c mt-1">
            <span class="text-xs text-g-600">今日</span>
            <span class="ml-1 text-xs font-semibold text-success">¥{{ Number(data.revenue?.today || 0).toFixed(2) }}</span>
          </div>
          <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
            <ArtSvgIcon icon="ri:wallet-3-line" class="text-xl text-theme" />
          </div>
        </div>
      </ElCol>
      <ElCol :sm="12" :md="6">
        <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
          <span class="text-g-700 text-sm">总消费</span>
          <ArtCountTo class="text-[26px] font-medium mt-2" prefix="¥ " :decimals="2" :target="Number(data.revenue?.consume) || 0" :duration="1300" />
          <div class="flex-c mt-1">
            <span class="text-xs text-g-600">本月</span>
            <span class="ml-1 text-xs font-semibold text-success">{{ data.orders?.month || 0 }} 单</span>
          </div>
          <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
            <ArtSvgIcon icon="ri:money-cny-circle-line" class="text-xl text-theme" />
          </div>
        </div>
      </ElCol>
    </ElRow>

    <!-- 订单状态分布 -->
    <ElRow :gutter="20">
      <ElCol :md="12">
        <div class="art-card p-5 mb-5 h-105">
          <div class="art-card-header">
            <div class="title">
              <h4>订单状态分布</h4>
              <p>各状态订单数</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4 mt-6">
            <div class="art-card p-4">
              <div class="text-g-600 text-sm">待处理</div>
              <div class="text-2xl font-medium mt-2 text-warning">{{ data.orders?.pending || 0 }}</div>
            </div>
            <div class="art-card p-4">
              <div class="text-g-600 text-sm">已完成</div>
              <div class="text-2xl font-medium mt-2 text-success">{{ data.orders?.done || 0 }}</div>
            </div>
            <div class="art-card p-4">
              <div class="text-g-600 text-sm">异常</div>
              <div class="text-2xl font-medium mt-2 text-danger">{{ data.orders?.error || 0 }}</div>
            </div>
            <div class="art-card p-4">
              <div class="text-g-600 text-sm">已退款</div>
              <div class="text-2xl font-medium mt-2 text-g-700">{{ data.orders?.refunded || 0 }}</div>
            </div>
          </div>
        </div>
      </ElCol>
      <ElCol :md="12">
        <div class="art-card p-5 mb-5 h-105">
          <div class="art-card-header">
            <div class="title">
              <h4>近 7 天订单趋势</h4>
              <p>每日订单数</p>
            </div>
          </div>
          <ArtLineChart
            v-if="chartData.length"
            height="calc(100% - 56px)"
            :data="chartData"
            :xAxisData="chartLabels"
            :showAreaColor="true"
            :showAxisLine="false"
          />
        </div>
      </ElCol>
    </ElRow>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { getAdminDashboard } from '@/api/wk'
  import request from '@/utils/http'

  defineOptions({ name: 'WkAdminData' })

  const data = ref<any>({})
  const chartLabels = ref<string[]>([])
  const chartData = ref<number[]>([])

  onMounted(async () => {
    data.value = await getAdminDashboard()
    const res: any = await request.get({ url: '/api/order/monthly' })
    chartLabels.value = res.labels || []
    chartData.value = res.values || []
  })
</script>

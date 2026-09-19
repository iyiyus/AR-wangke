<template>
  <div class="art-card p-5 flex-b mb-5 max-sm:mb-4">
    <div>
      <h2 class="text-2xl font-medium">关于平台</h2>
      <p class="text-g-700 mt-1">{{ systemName }} 是一款多平台网课代刷管理系统</p>
      <p class="text-g-700 mt-1">支持多平台聚合下单、代理分销、密价定价、API 对接等</p>

      <div class="flex flex-wrap gap-3.5 max-w-150 mt-9">
        <div
          class="w-60 flex-cb h-12.5 px-3.5 border border-g-300 c-p rounded-lg text-sm bg-g-100 duration-300 hover:-translate-y-1 max-sm:w-full"
          v-for="link in linkList"
          :key="link.label"
          @click="$router.push(link.url)"
        >
          <span class="text-g-700">{{ link.label }}</span>
          <ArtSvgIcon icon="ri:arrow-right-s-line" class="text-lg text-g-600" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import AppConfig from '@/config'

  const systemName = AppConfig.systemInfo.name

  const linkList = [
    { label: '提交任务', url: '/wk/submit' },
    { label: '我的订单', url: '/wk/order' },
    { label: '账户充值', url: '/wk/pay' },
    { label: '联系上级', url: '/wk/charge' }
  ]
</script>

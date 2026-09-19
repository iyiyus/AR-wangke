<template>
  <div class="art-card h-128 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>快捷入口</h4>
        <p>常用功能一键直达</p>
      </div>
    </div>
    <div class="grid grid-cols-2 gap-3 mt-4">
      <div
        v-for="item in actions"
        :key="item.path"
        class="art-card flex flex-col items-center justify-center h-24 cursor-pointer duration-300 hover:-translate-y-1"
        @click="$router.push(item.path)"
      >
        <div class="size-10 rounded-xl flex-cc bg-theme/10">
          <ArtSvgIcon :icon="item.icon" class="text-xl text-theme" />
        </div>
        <span class="mt-2 text-sm">{{ item.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  const actions = [
    { label: '提交任务', icon: 'ri:add-circle-line', path: '/wk/submit' },
    { label: '账户充值', icon: 'ri:wallet-3-line', path: '/wk/pay' },
    { label: '任务列表', icon: 'ri:list-check-2', path: '/wk/order' },
    { label: '工单系统', icon: 'ri:customer-service-2-line', path: '/wk/gongdan' },
    { label: '我的价格', icon: 'ri:price-tag-3-line', path: '/wk/myprice' },
    { label: '联系上级', icon: 'ri:contacts-line', path: '/wk/charge' }
  ]
</script>

<template>
  <div>
    <CardList />

    <ElRow :gutter="20">
      <ElCol :sm="24" :md="12" :lg="10">
        <Platforms />
      </ElCol>
      <ElCol :sm="24" :md="12" :lg="14">
        <OrdersOverview />
      </ElCol>
    </ElRow>

    <ElRow :gutter="20">
      <ElCol :sm="24" :md="24" :lg="12">
        <RecentOrders />
      </ElCol>
      <ElCol :sm="24" :md="12" :lg="6">
        <DynamicStats />
      </ElCol>
      <ElCol :sm="24" :md="12" :lg="6">
        <QuickActions />
      </ElCol>
    </ElRow>

    <AboutProject />
  </div>
</template>

<script setup lang="ts">
  import CardList from './modules/card-list.vue'
  import Platforms from './modules/platforms.vue'
  import OrdersOverview from './modules/orders-overview.vue'
  import RecentOrders from './modules/recent-orders.vue'
  import DynamicStats from './modules/dynamic-stats.vue'
  import QuickActions from './modules/quick-actions.vue'
  import AboutProject from './modules/about-project.vue'

  defineOptions({ name: 'WkDashboard' })
</script>

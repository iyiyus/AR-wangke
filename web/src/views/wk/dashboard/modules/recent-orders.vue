<template>
  <div class="art-card p-5 h-128 overflow-hidden mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>最近订单</h4>
        <p>共 {{ list.length }} 条</p>
      </div>
      <ElButton size="small" @click="$router.push('/wk/order')">查看全部</ElButton>
    </div>
    <ArtTable
      class="w-full"
      :data="list"
      style="width: 100%"
      size="large"
      :border="false"
      :stripe="false"
      :header-cell-style="{ background: 'transparent' }"
    >
      <template #default>
        <ElTableColumn label="订单ID" prop="oid" width="80" />
        <ElTableColumn label="平台" prop="ptname" />
        <ElTableColumn label="账号" prop="user" />
        <ElTableColumn label="课程" prop="kcname" show-overflow-tooltip />
        <ElTableColumn label="金额" prop="fees" width="80" />
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="statusType(row.status)" size="small">{{ row.status }}</ElTag>
          </template>
        </ElTableColumn>
      </template>
    </ArtTable>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElTag } from 'element-plus'
  import { getOrderList } from '@/api/wk'

  const list = ref<any[]>([])

  const statusType = (s: string) => (
    { 待处理: 'info', 进行中: 'warning', 已完成: 'success', 异常: 'danger' }[s] || 'info'
  ) as any

  onMounted(async () => {
    const res: any = await getOrderList({ current: 1, size: 6 })
    list.value = res.list || []
  })
</script>

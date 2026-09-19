<template>
  <div class="paylist-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-g-600 text-sm">充值记录</span>
        </template>
      </ArtTableHeader>
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ElTag } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import { getPayList } from '@/api/wk'

  defineOptions({ name: 'WkPayList' })

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getPayList,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'out_trade_no', label: '订单号', minWidth: 200, showOverflowTooltip: true },
        { prop: 'type', label: '支付方式', width: 100 },
        {
          prop: 'money', label: '金额', width: 100,
          formatter: (row) => `¥ ${row.money}`
        },
        {
          prop: 'status', label: '状态', width: 90,
          formatter: (row) => h(ElTag, { type: row.status === 1 ? 'success' : 'info', size: 'small' },
            () => row.status === 1 ? '已完成' : '待支付')
        },
        { prop: 'addtime', label: '时间', width: 180 }
      ]
    }
  })
</script>

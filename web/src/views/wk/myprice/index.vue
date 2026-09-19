<template>
  <div class="myprice-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-g-600 text-sm">我的价格（密价）</span>
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
  import { useTable } from '@/hooks/core/useTable'
  import { getMyPriceList } from '@/api/wk'

  defineOptions({ name: 'WkMyPrice' })

  const modeLabel = (m: number) => ['价格扣除', '倍数扣除', '直接定价'][m] || '未知'

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getMyPriceList,
      apiParams: { current: 1, size: 50 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'cid', label: '平台ID', width: 80 },
        { prop: 'mode', label: '模式', width: 100, formatter: (row) => modeLabel(row.mode) },
        { prop: 'price', label: '价格', width: 100 },
        { prop: 'addtime', label: '设置时间', width: 160 }
      ]
    }
  })
</script>

<template>
  <div class="log-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <span class="text-g-600 text-sm">消费记录</span>
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
  import { getAdminLog } from '@/api/wk'
  import { useUserStore } from '@/store/modules/user'

  defineOptions({ name: 'WkLog' })

  const userStore = useUserStore()
  const uid = (userStore.info as any)?.uid

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getAdminLog,
      apiParams: { current: 1, size: 15, uid: String(uid) },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'type', label: '类型', width: 100 },
        { prop: 'text', label: '说明', minWidth: 200, showOverflowTooltip: true },
        {
          prop: 'money', label: '金额', width: 100,
          formatter: (row) => h(ElTag, {
            type: parseFloat(row.money) >= 0 ? 'success' : 'danger', size: 'small'
          }, () => row.money)
        },
        { prop: 'smoney', label: '余额', width: 100 },
        { prop: 'addtime', label: '时间', width: 180 }
      ]
    }
  })
</script>

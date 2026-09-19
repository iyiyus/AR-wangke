<template>
  <div class="admin-log-page art-full-height">
    <ElCard shadow="never" class="mb-3">
      <ElForm :model="searchForm" inline>
        <ElFormItem label="UID">
          <ElInput v-model="searchForm.uid" placeholder="UID" clearable style="width:120px" />
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch(searchForm)">查询</ElButton>
          <ElButton @click="resetSearchParams">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import { getAdminLog, batchDeleteLog } from '@/api/wk'

  defineOptions({ name: 'WkAdminLog' })

  const searchForm = ref({ uid: '' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.id)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条日志？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteLog(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination, searchParams,
    resetSearchParams, handleSizeChange, handleCurrentChange, refreshData, getData } = useTable({
    core: {
      apiFn: getAdminLog,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'uid', label: 'UID', width: 100 },
        { prop: 'type', label: '类型', width: 140 },
        { prop: 'text', label: '说明', minWidth: 240, showOverflowTooltip: true },
        {
          prop: 'money', label: '金额', width: 110,
          formatter: (row: any) => h(ElTag, {
            type: parseFloat(row.money) >= 0 ? 'success' : 'danger', size: 'small'
          }, () => row.money)
        },
        { prop: 'smoney', label: '余额', width: 110 },
        { prop: 'ip', label: 'IP', width: 150 },
        { prop: 'addtime', label: '时间', width: 200 }
      ]
    }
  })

  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }
</script>

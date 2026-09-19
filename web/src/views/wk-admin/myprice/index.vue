<template>
  <div class="admin-myprice-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="dialogVisible = true">设置密价</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <ElDialog v-model="dialogVisible" title="设置密价" width="460px" align-center>
      <ElForm :model="form" label-width="90px">
        <ElFormItem label="用户UID"><ElInput v-model="form.uid" type="number" /></ElFormItem>
        <ElFormItem label="平台CID"><ElInput v-model="form.cid" type="number" /></ElFormItem>
        <ElFormItem label="模式">
          <ElSelect v-model="form.mode" style="width:100%">
            <ElOption label="价格基础上扣除" :value="0" />
            <ElOption label="倍数基础上扣除" :value="1" />
            <ElOption label="直接定价" :value="2" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="价格"><ElInput v-model="form.price" type="number" step="0.01" /></ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleSave">保存</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { getMyPriceList, setMyPrice, deleteMyPrice, batchDeleteMyPrice } from '@/api/wk'

  defineOptions({ name: 'WkAdminMyPrice' })

  const dialogVisible = ref(false)
  const submitting = ref(false)
  const form = ref({ uid: 0, cid: 0, mode: 0, price: '' })
  const selectedIds = ref<number[]>([])
  const modeLabel = (m: number) => ['价格扣除', '倍数扣除', '直接定价'][m] || '未知'

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.mid)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条密价？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteMyPrice(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getMyPriceList,
      apiParams: { current: 1, size: 50 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'uid', label: 'UID', width: 100 },
        { prop: 'cid', label: '平台ID', width: 100 },
        { prop: 'mode', label: '模式', minWidth: 140, formatter: (row) => modeLabel(row.mode) },
        { prop: 'price', label: '价格', width: 120 },
        { prop: 'addtime', label: '设置时间', minWidth: 200 },
        {
          prop: 'operation', label: '操作', width: 100, fixed: 'right',
          formatter: (row: any) => h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
        }
      ]
    }
  })

  const handleSave = async () => {
    submitting.value = true
    try {
      await setMyPrice(form.value)
      ElMessage.success('设置成功')
      dialogVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleDelete = async (row: any) => {
    await ElMessageBox.confirm('确认删除该密价？', '提示', { type: 'warning' })
    await deleteMyPrice(row.mid)
    ElMessage.success('删除成功')
    refreshData()
  }
</script>

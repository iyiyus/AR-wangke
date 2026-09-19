<template>
  <div class="admin-dengji-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="openDialog()">添加等级</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <ElDialog v-model="dialogVisible" :title="editId ? '编辑等级' : '添加等级'" width="420px" align-center>
      <ElForm :model="form" label-width="80px">
        <ElFormItem label="等级名称"><ElInput v-model="form.name" /></ElFormItem>
        <ElFormItem label="费率"><ElInput v-model="form.rate" type="number" step="0.01" /></ElFormItem>
        <ElFormItem label="排序"><ElInput v-model="form.sort" type="number" /></ElFormItem>
        <ElFormItem label="状态">
          <ElSwitch v-model="form.status" active-value="1" inactive-value="0" />
        </ElFormItem>
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
  import { getDengjiList, addDengji, updateDengji, deleteDengji, batchDeleteDengji } from '@/api/wk'

  defineOptions({ name: 'WkAdminDengji' })

  const dialogVisible = ref(false)
  const submitting = ref(false)
  const editId = ref<number | null>(null)
  const form = ref<Partial<Api.Dengji.Item>>({ status: '1', sort: '10' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.id)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条等级？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteDengji(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getDengjiList,
      apiParams: { current: 1, size: 50 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'name', label: '等级名称', minWidth: 200 },
        { prop: 'rate', label: '费率', width: 120 },
        { prop: 'sort', label: '排序', width: 100 },
        {
          prop: 'operation', label: '操作', width: 120, fixed: 'right',
          formatter: (row: any) => h('div', { class: 'flex gap-1' }, [
            h(ArtButtonTable, { type: 'edit', onClick: () => openDialog(row) }),
            h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
          ])
        }
      ]
    }
  })

  const openDialog = (row?: Api.Dengji.Item) => {
    editId.value = row?.id ?? null
    form.value = row ? { ...row } : { status: '1', sort: '10' }
    dialogVisible.value = true
  }

  const handleSave = async () => {
    submitting.value = true
    try {
      if (editId.value) {
        await updateDengji(editId.value, form.value)
      } else {
        await addDengji(form.value)
      }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleDelete = async (row: Api.Dengji.Item) => {
    await ElMessageBox.confirm(`确认删除等级「${row.name}」？`, '提示', { type: 'warning' })
    await deleteDengji(row.id)
    ElMessage.success('删除成功')
    refreshData()
  }
</script>

<template>
  <div class="admin-fenlei-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="openDialog()">添加分类</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <ElDialog v-model="dialogVisible" :title="editId ? '编辑分类' : '添加分类'" width="420px" align-center>
      <ElForm :model="form" label-width="80px">
        <ElFormItem label="分类ID">
          <ElInput v-model="form.id" type="number"
            :disabled="!!editId"
            placeholder="可填写与对方平台一致的 ID，留空则自增" />
        </ElFormItem>
        <ElFormItem label="分类名称"><ElInput v-model="form.name" /></ElFormItem>
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
  import {
    getFenLeiList, addFenLeiAdmin, updateFenLeiAdmin, deleteFenLeiAdmin,
    batchDeleteFenLei
  } from '@/api/wk'

  defineOptions({ name: 'WkAdminFenLei' })

  const dialogVisible = ref(false)
  const submitting = ref(false)
  const editId = ref<number | null>(null)
  const form = ref<any>({ status: '1', sort: '0' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.id)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条分类？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteFenLei(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getFenLeiList,
      apiParams: { current: 1, size: 50 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'id', label: 'ID', width: 80 },
        { prop: 'name', label: '分类名称', minWidth: 200 },
        { prop: 'sort', label: '排序', width: 100 },
        { prop: 'time', label: '创建时间', minWidth: 200 },
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

  const openDialog = (row?: any) => {
    editId.value = row?.id ?? null
    form.value = row ? { ...row } : { status: '1', sort: '0' }
    dialogVisible.value = true
  }

  const handleSave = async () => {
    submitting.value = true
    try {
      if (editId.value) {
        await updateFenLeiAdmin(editId.value, form.value)
      } else {
        await addFenLeiAdmin(form.value)
      }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleDelete = async (row: any) => {
    await ElMessageBox.confirm(`确认删除分类「${row.name}」？`, '提示', { type: 'warning' })
    await deleteFenLeiAdmin(row.id)
    ElMessage.success('删除成功')
    refreshData()
  }
</script>

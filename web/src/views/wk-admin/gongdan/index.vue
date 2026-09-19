<template>
  <div class="admin-gongdan-page art-full-height">
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

    <ElDialog v-model="replyVisible" title="回复工单" width="500px" align-center>
      <ElForm :model="replyForm" label-width="80px">
        <ElFormItem label="工单内容">
          <div class="text-g-600 text-sm">{{ replyForm.content }}</div>
        </ElFormItem>
        <ElFormItem label="回复内容">
          <ElInput v-model="replyForm.answer" type="textarea" :rows="4" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSelect v-model="replyForm.state">
            <ElOption label="待处理" value="待处理" />
            <ElOption label="处理中" value="处理中" />
            <ElOption label="已解决" value="已解决" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="replyVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleReply">提交回复</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { getGongDanList, replyGongDan, batchDeleteGongDan } from '@/api/wk'

  defineOptions({ name: 'WkAdminGongDan' })

  const replyVisible = ref(false)
  const submitting = ref(false)
  const replyForm = ref({ gid: 0, content: '', answer: '', state: '处理中' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.gid)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条工单？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteGongDan(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const stateType = (s: string) => ({ 待处理: 'info', 处理中: 'warning', 已解决: 'success' }[s] || 'info') as any

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getGongDanList,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'uid', label: 'UID', width: 80 },
        { prop: 'region', label: '类型', width: 110 },
        { prop: 'title', label: '标题', width: 180, showOverflowTooltip: true },
        { prop: 'content', label: '内容', minWidth: 200, showOverflowTooltip: true },
        { prop: 'answer', label: '回复', minWidth: 180, showOverflowTooltip: true },
        {
          prop: 'state', label: '状态', width: 100,
          formatter: (row: any) => h(ElTag, { type: stateType(row.state), size: 'small' }, () => row.state)
        },
        { prop: 'addtime', label: '时间', width: 180 },
        {
          prop: 'operation', label: '操作', width: 100, fixed: 'right',
          formatter: (row: any) => h(ArtButtonTable, { type: 'edit', title: '回复', onClick: () => openReply(row) })
        }
      ]
    }
  })

  const openReply = (row: Api.GongDan.Item) => {
    replyForm.value = { gid: row.gid, content: row.content, answer: row.answer || '', state: row.state || '处理中' }
    replyVisible.value = true
  }

  const handleReply = async () => {
    submitting.value = true
    try {
      await replyGongDan({ gid: replyForm.value.gid, answer: replyForm.value.answer, state: replyForm.value.state })
      ElMessage.success('回复成功')
      replyVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }
</script>

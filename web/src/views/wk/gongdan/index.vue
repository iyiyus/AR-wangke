<template>
  <div class="gongdan-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="addVisible = true">提交工单</ElButton>
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

    <!-- 提交工单弹窗 -->
    <ElDialog v-model="addVisible" title="提交工单" width="500px" align-center>
      <ElForm :model="addForm" label-width="80px">
        <ElFormItem label="工单类型">
          <ElSelect v-model="addForm.region" style="width:100%">
            <ElOption label="订单问题" value="订单问题" />
            <ElOption label="充值问题" value="充值问题" />
            <ElOption label="账号问题" value="账号问题" />
            <ElOption label="其他" value="其他" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="标题">
          <ElInput v-model="addForm.title" placeholder="工单标题" />
        </ElFormItem>
        <ElFormItem label="内容">
          <ElInput v-model="addForm.content" type="textarea" :rows="4" placeholder="请详细描述问题" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="addVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleAdd">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElTag, ElMessage } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import { getGongDanList, addGongDan } from '@/api/wk'

  defineOptions({ name: 'WkGongDan' })

  const addVisible = ref(false)
  const submitting = ref(false)
  const addForm = ref({ region: '订单问题', title: '', content: '' })

  const stateType = (s: string) => ({ '待处理': 'info', '处理中': 'warning', '已解决': 'success' }[s] || 'info') as any

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getGongDanList,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'region', label: '类型', width: 100 },
        { prop: 'title', label: '标题', minWidth: 150, showOverflowTooltip: true },
        { prop: 'content', label: '内容', minWidth: 200, showOverflowTooltip: true },
        { prop: 'answer', label: '回复', minWidth: 150, showOverflowTooltip: true },
        {
          prop: 'state', label: '状态', width: 90,
          formatter: (row: any) => h(ElTag, { type: stateType(row.state), size: 'small' }, () => row.state)
        },
        { prop: 'addtime', label: '时间', width: 160 }
      ]
    }
  })

  const handleAdd = async () => {
    if (!addForm.value.content) { ElMessage.warning('请填写内容'); return }
    submitting.value = true
    try {
      await addGongDan(addForm.value)
      ElMessage.success('提交成功')
      addVisible.value = false
      addForm.value = { region: '订单问题', title: '', content: '' }
      refreshData()
    } finally {
      submitting.value = false
    }
  }
</script>

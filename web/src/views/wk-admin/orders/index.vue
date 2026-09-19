<template>
  <div class="admin-orders-page art-full-height">
    <ElCard shadow="never" class="mb-3">
      <ElForm :model="searchForm" inline>
        <ElFormItem label="订单ID">
          <ElInput v-model="searchForm.oid" clearable style="width:100px" />
        </ElFormItem>
        <ElFormItem label="UID">
          <ElInput v-model="searchForm.uid" clearable style="width:90px" />
        </ElFormItem>
        <ElFormItem label="账号">
          <ElInput v-model="searchForm.qq" clearable style="width:140px" />
        </ElFormItem>
        <ElFormItem label="任务状态">
          <ElSelect v-model="searchForm.status_text" clearable style="width:110px">
            <ElOption v-for="s in statusOptions" :key="s" :label="s" :value="s" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="处理状态">
          <ElSelect v-model="searchForm.dock" clearable style="width:110px">
            <ElOption label="待处理" value="0" />
            <ElOption label="已完成" value="1" />
            <ElOption label="处理失败" value="2" />
            <ElOption label="重复下单" value="3" />
            <ElOption label="已取消" value="4" />
            <ElOption label="我的" value="99" />
          </ElSelect>
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
          <ElSpace>
            <ElButton @click="batchSetStatus('已完成')" :disabled="!selectedOids.length">标记完成</ElButton>
            <ElButton @click="batchSetStatus('待处理')" :disabled="!selectedOids.length">标记待处理</ElButton>
            <ElButton type="danger" @click="handleRefund" :disabled="!selectedOids.length">退款</ElButton>
            <ElButton type="danger" @click="handleBatchDelete" :disabled="!selectedOids.length">
              硬删除 ({{ selectedOids.length }})
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="data" :columns="columns" :pagination="pagination"
        @selection-change="selectedOids = $event.map((r: any) => r.oid)"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import { getOrderList, batchOrderStatus, refundOrder, batchDeleteOrder } from '@/api/wk'

  defineOptions({ name: 'WkAdminOrders' })

  const selectedOids = ref<number[]>([])

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedOids.value.length} 个订单？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteOrder(selectedOids.value)
    ElMessage.success(`已删除 ${res.deleted} 个`)
    selectedOids.value = []
    refreshData()
  }
  const searchForm = ref({ oid: '', uid: '', qq: '', status_text: '', dock: '' })
  const statusOptions = ['待处理', '进行中', '已完成', '补刷中', '异常', '已取消', '已退款']

  const statusTagType = (s: string) => ({ 待处理: 'info', 进行中: 'warning', 已完成: 'success', 异常: 'danger' }[s] || 'info') as any
  const dockLabel = (d: string) => ({ '0': '待处理', '1': '已完成', '2': '处理失败', '3': '重复', '4': '已取消', '99': '我的' }[d] || d)
  const dockTagType = (d: string) => ({ '0': 'info', '1': 'success', '2': 'danger', '3': 'warning', '99': 'warning' }[d] || 'info') as any

  const { columns, columnChecks, data, loading, pagination, searchParams,
    resetSearchParams, handleSizeChange, handleCurrentChange, refreshData, getData } = useTable({
    core: {
      apiFn: getOrderList,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'selection' },
        { prop: 'oid', label: 'ID', width: 80 },
        { prop: 'uid', label: 'UID', width: 70 },
        { prop: 'ptname', label: '平台', width: 120 },
        { prop: 'user', label: '账号', width: 140 },
        { prop: 'kcname', label: '课程', minWidth: 180, showOverflowTooltip: true },
        { prop: 'fees', label: '金额', width: 80 },
        {
          prop: 'status', label: '任务状态', width: 100,
          formatter: (row: any) => h(ElTag, { type: statusTagType(row.status), size: 'small' }, () => row.status)
        },
        {
          prop: 'dockstatus', label: '处理状态', width: 100,
          formatter: (row: any) => h(ElTag, { type: dockTagType(row.dockstatus), size: 'small' }, () => dockLabel(row.dockstatus))
        },
        { prop: 'process', label: '进度', width: 120, showOverflowTooltip: true },
        { prop: 'remarks', label: '备注', width: 120, showOverflowTooltip: true },
        { prop: 'addtime', label: '时间', width: 180 }
      ]
    }
  })

  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }

  const batchSetStatus = async (status: string) => {
    await batchOrderStatus(selectedOids.value, status)
    ElMessage.success('操作成功')
    refreshData()
  }

  const handleRefund = async () => {
    await ElMessageBox.confirm(`确认对选中的 ${selectedOids.value.length} 条订单退款？`, '提示', { type: 'warning' })
    await refundOrder(selectedOids.value)
    ElMessage.success('退款成功')
    refreshData()
  }
</script>

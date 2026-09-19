<template>
  <div class="order-page art-full-height">
    <!-- 搜索栏 -->
    <ElCard shadow="never" class="mb-3">
      <ElForm :model="searchForm" inline>
        <ElFormItem label="订单ID">
          <ElInput v-model="searchForm.oid" placeholder="订单ID" clearable style="width:110px" />
        </ElFormItem>
        <ElFormItem label="账号">
          <ElInput v-model="searchForm.qq" placeholder="学号/手机号" clearable style="width:150px" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSelect v-model="searchForm.status_text" placeholder="全部" clearable style="width:110px">
            <ElOption v-for="s in statusOptions" :key="s" :label="s" :value="s" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch(searchForm)">查询</ElButton>
          <ElButton @click="resetSearchParams">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData" />

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />

      <!-- 详情弹窗 -->
      <ElDialog v-model="detailVisible" title="订单详情" width="600px" align-center>
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="订单ID">{{ detail.oid }}</ElDescriptionsItem>
          <ElDescriptionsItem label="平台">{{ detail.ptname }}</ElDescriptionsItem>
          <ElDescriptionsItem label="学校">{{ detail.school }}</ElDescriptionsItem>
          <ElDescriptionsItem label="账号">{{ detail.user }}</ElDescriptionsItem>
          <ElDescriptionsItem label="密码">{{ detail.pass }}</ElDescriptionsItem>
          <ElDescriptionsItem label="课程">{{ detail.kcname }}</ElDescriptionsItem>
          <ElDescriptionsItem label="状态">{{ detail.status }}</ElDescriptionsItem>
          <ElDescriptionsItem label="进度">{{ detail.process }}</ElDescriptionsItem>
          <ElDescriptionsItem label="备注" :span="2">{{ detail.remarks }}</ElDescriptionsItem>
          <ElDescriptionsItem label="提交时间" :span="2">{{ detail.addtime }}</ElDescriptionsItem>
        </ElDescriptions>
      </ElDialog>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElTag, ElMessageBox, ElMessage } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import { getOrderList, cancelOrder, restartOrder } from '@/api/wk'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'

  defineOptions({ name: 'WkOrder' })

  type OrderItem = Api.Order.Item

  const detailVisible = ref(false)
  const detail = ref<Partial<OrderItem>>({})

  const searchForm = ref({ oid: '', qq: '', status_text: '' })
  const statusOptions = ['待处理', '进行中', '已完成', '补刷中', '异常', '已取消', '已退款']

  const statusTagType = (s: string) => {
    const map: Record<string, string> = {
      待处理: 'info', 进行中: 'warning', 已完成: 'success',
      补刷中: 'warning', 异常: 'danger', 已取消: 'info', 已退款: 'info'
    }
    return (map[s] || 'info') as any
  }

  const { columns, columnChecks, data, loading, pagination, searchParams,
    resetSearchParams, handleSizeChange, handleCurrentChange, refreshData, getData } = useTable({
    core: {
      apiFn: getOrderList,
      apiParams: { current: 1, size: 15, ...searchForm.value },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'oid', label: 'ID', width: 80 },
        { prop: 'ptname', label: '平台', width: 120 },
        { prop: 'user', label: '账号', width: 140 },
        { prop: 'kcname', label: '课程', minWidth: 180, showOverflowTooltip: true },
        { prop: 'fees', label: '金额', width: 80 },
        {
          prop: 'status', label: '状态', width: 100,
          formatter: (row: any) => h(ElTag, { type: statusTagType(row.status), size: 'small' }, () => row.status)
        },
        { prop: 'process', label: '进度', width: 120, showOverflowTooltip: true },
        { prop: 'remarks', label: '备注', width: 120, showOverflowTooltip: true },
        { prop: 'addtime', label: '时间', width: 180 },
        {
          prop: 'operation', label: '操作', width: 130, fixed: 'right',
          formatter: (row: any) => h('div', { class: 'flex gap-1' }, [
            h(ArtButtonTable, { type: 'edit', title: '详情', onClick: () => { detail.value = row; detailVisible.value = true } }),
            h(ArtButtonTable, { type: 'delete', title: '取消', onClick: () => handleCancel(row) })
          ])
        }
      ]
    }
  })

  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }

  const handleCancel = async (row: OrderItem) => {
    await ElMessageBox.confirm('确认取消该订单？', '提示', { type: 'warning' })
    await cancelOrder(row.oid)
    ElMessage.success('已取消')
    refreshData()
  }
</script>

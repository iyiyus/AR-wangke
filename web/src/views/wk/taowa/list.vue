<template>
  <div class="taowa-list-page">
    <div class="art-card p-5">
      <div class="flex-cb mb-4">
        <div class="flex items-center gap-3">
          <h4 class="font-bold">盖章任务列表</h4>
          <ElTag type="primary" size="small">实习助手</ElTag>
        </div>
        <div class="flex gap-2">
          <ElSelect v-model="filter.status" placeholder="状态筛选" clearable size="small" style="width: 130px" @change="load(1)">
            <ElOption v-for="s in statusList" :key="s" :label="s" :value="s" />
          </ElSelect>
          <ElButton size="small" :loading="loading" @click="load(pagination.current)">刷新</ElButton>
        </div>
      </div>

      <ElTable :data="list" v-loading="loading" stripe>
        <ElTableColumn prop="oid" label="订单号" width="80" />
        <ElTableColumn prop="ptname" label="公司" min-width="120" show-overflow-tooltip />
        <ElTableColumn prop="school" label="姓名" width="90" />
        <ElTableColumn prop="user" label="电话" width="120" />
        <ElTableColumn prop="pass" label="收货地址" min-width="150" show-overflow-tooltip />
        <ElTableColumn prop="kcname" label="规格" min-width="140" show-overflow-tooltip />
        <ElTableColumn label="文件" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.courseEndTime" class="text-g-600">{{ row.courseEndTime }}</span>
            <span v-else class="text-g-400">无</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="fees" label="金额" width="90">
          <template #default="{ row }">¥{{ row.fees }}</template>
        </ElTableColumn>
        <ElTableColumn prop="status" label="状态" width="90">
          <template #default="{ row }">
            <ElTag :type="statusType(row.status)" size="small">{{ row.status }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="addtime" label="下单时间" width="160" />
        <ElTableColumn label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <ElButton size="small" link type="primary" @click="handleCancel(row)">取消</ElButton>
            <ElButton size="small" link type="primary" @click="handleGxx(row)">补充文件</ElButton>
            <ElButton size="small" link type="primary" @click="handleRemark(row)">备注</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <div class="flex justify-end mt-4">
        <ElPagination background layout="total, prev, pager, next" :total="total" :page-size="pagination.pageSize"
          :current-page="pagination.current" @current-change="load" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getTaowaOrderList, cancelTaowaOrder, taowaOrderRemark, taowaOrderGxx } from '@/api/wk'

  defineOptions({ name: 'WkTaowaList' })

  const loading = ref(false)
  const list = ref<any[]>([])
  const total = ref(0)
  const filter = reactive({ status: '' })
  const pagination = reactive({ current: 1, pageSize: 10 })
  const statusList = ['待处理', '处理中', '已完成', '已取消', '已退款']

  const statusType = (s: string) => {
    if (s === '已完成') return 'success'
    if (s === '已取消' || s === '已退款') return 'danger'
    if (s === '处理中') return 'primary'
    return 'warning'
  }

  const load = async (page = 1) => {
    pagination.current = page
    loading.value = true
    try {
      const res: any = await getTaowaOrderList({
        type: 'gz', page: pagination.current, pageSize: pagination.pageSize, status: filter.status
      })
      list.value = res?.list || []
      total.value = res?.total || 0
    } catch (e: any) {
      ElMessage.error(e?.message || '加载失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = async (row: any) => {
    await ElMessageBox.confirm('确定取消该订单？取消后金额将退回余额。', '取消订单', { type: 'warning' })
    const res: any = await cancelTaowaOrder(row.oid)
    ElMessage.success(res?.msg || '已取消')
    load()
  }

  const handleRemark = async (row: any) => {
    const { value } = await ElMessageBox.prompt('输入新备注', '修改备注', { inputValue: row.remarks || '' })
    await taowaOrderRemark(row.oid, value)
    ElMessage.success('已修改')
    load()
  }

  const handleGxx = async (row: any) => {
    const { value } = await ElMessageBox.prompt('输入补充文件名（多个用逗号分隔）', '补充文件')
    if (!value?.trim()) return
    await taowaOrderGxx(row.oid, value.trim())
    ElMessage.success('已补充')
    load()
  }

  onMounted(() => load())
</script>

<template>
  <div class="admin-huoyuan-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="openDialog()">添加货源</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <ElDialog v-model="dialogVisible" :title="editId ? '编辑货源' : '添加货源'" width="500px" align-center>
      <ElForm :model="form" label-width="100px">
        <ElFormItem label="平台标识">
          <ElSelect v-model="form.pt" placeholder="选择平台类型">
            <ElOption label="27 (HTTP)" value="27" />
            <ElOption label="27s (HTTPS)" value="27s" />
            <ElOption label="29 (HTTP) 标准网课" value="29" />
            <ElOption label="29s (HTTPS) 标准网课" value="29s" />
            <ElOption label="小月 xy" value="xy" />
            <ElOption label="oligei / benz" value="oligei" />
            <ElOption label="哥斯拉 2.0" value="gsl2.0" />
            <ElOption label="网课联盟 (考试)" value="wklm" />
            <ElOption label="网课联盟 (不考试)" value="wklmbks" />
            <ElOption label="欧巴 (普通)" value="ouba" />
            <ElOption label="欧巴 (考试)" value="obks" />
            <ElOption label="欧巴 (视频)" value="obsp" />
            <ElOption label="00 隐藏接口" value="00" />
            <ElOption label="止水接口" value="zs" />
            <ElOption label="捐赠接口" value="jz" />
            <ElOption label="鸡腿接口" value="jt" />
            <ElOption label="牛牛接口" value="niuniu" />
            <ElOption label="指尖接口" value="zijian" />
            <ElOption label="学习通官方" value="xxtgf" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="名称"><ElInput v-model="form.name" /></ElFormItem>
        <ElFormItem label="URL"><ElInput v-model="form.url" placeholder="不带 http 顶级域名" /></ElFormItem>
        <ElFormItem label="账号"><ElInput v-model="form.user" /></ElFormItem>
        <ElFormItem label="密码"><ElInput v-model="form.pass" /></ElFormItem>
        <ElFormItem label="Token"><ElInput v-model="form.token" type="textarea" :rows="2" /></ElFormItem>
        <ElFormItem label="Cookie"><ElInput v-model="form.cookie" type="textarea" :rows="2" /></ElFormItem>
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
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { getHuoYuanList, addHuoYuan, updateHuoYuan, deleteHuoYuan, batchDeleteHuoYuan } from '@/api/wk'

  defineOptions({ name: 'WkAdminHuoYuan' })

  const dialogVisible = ref(false)
  const submitting = ref(false)
  const editId = ref<number | null>(null)
  const form = ref<any>({ status: '1' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.hid)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条货源？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteHuoYuan(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getHuoYuanList,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'hid', label: 'ID', width: 80 },
        { prop: 'pt', label: '平台类型', width: 110 },
        { prop: 'name', label: '名称', width: 150 },
        { prop: 'url', label: 'URL', minWidth: 200, showOverflowTooltip: true },
        { prop: 'user', label: '账号', width: 140 },
        { prop: 'money', label: '余额', width: 100 },
        {
          prop: 'status', label: '状态', width: 90,
          formatter: (row: any) => h(ElTag, { type: row.status === '1' ? 'success' : 'info', size: 'small' },
            () => row.status === '1' ? '正常' : '禁用')
        },
        { prop: 'addtime', label: '添加时间', width: 180 },
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
    editId.value = row?.hid ?? null
    form.value = row ? { ...row } : { status: '1' }
    dialogVisible.value = true
  }

  const handleSave = async () => {
    submitting.value = true
    try {
      if (editId.value) {
        await updateHuoYuan(editId.value, form.value)
      } else {
        await addHuoYuan(form.value)
      }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleDelete = async (row: any) => {
    await ElMessageBox.confirm(`确认删除货源「${row.name}」？`, '提示', { type: 'warning' })
    await deleteHuoYuan(row.hid)
    ElMessage.success('删除成功')
    refreshData()
  }
</script>

<template>
  <div class="admin-classes-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="openDialog()">添加平台</ElButton>
          <ElButton type="success" @click="importVisible = true">从货源批量导入</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <!-- 批量导入弹窗 -->
    <ElDialog v-model="importVisible" title="从货源批量导入平台" width="1000px" align-center>
      <ElForm :model="importForm" label-width="120px" inline>
        <ElFormItem label="选择货源">
          <ElSelect v-model="importForm.hid" placeholder="选择要拉取的货源" style="width:240px" filterable @change="handleFetch">
            <ElOption v-for="h in huoYuanList" :key="h.hid" :label="`[${h.pt}] ${h.name}`" :value="h.hid" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="使用远端数据">
          <ElSwitch v-model="importForm.use_remote" />
        </ElFormItem>
        <ElFormItem v-if="importForm.use_remote" label="加价方式">
          <ElRadioGroup v-model="importForm.markup_mode">
            <ElRadio value="multiply">倍数</ElRadio>
            <ElRadio value="add">加价</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="importForm.use_remote" :label="importForm.markup_mode === 'add' ? '加价金额' : '加价倍数'">
          <ElInput v-model="importForm.price_markup" style="width:150px"
            :placeholder="importForm.markup_mode === 'add' ? '如 0.5' : '如 1.2'" />
        </ElFormItem>
        <ElFormItem v-if="!importForm.use_remote" label="导入分类">
          <ElSelect v-model="importForm.fenlei" placeholder="选择分类" style="width:200px" clearable filterable>
            <ElOption v-for="f in fenLeiList" :key="f.id" :label="f.name" :value="String(f.id)" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="!importForm.use_remote" label="统一定价">
          <ElInput v-model="importForm.price" placeholder="¥1.50" style="width:120px" />
        </ElFormItem>
      </ElForm>

      <!-- 筛选 + 全选 -->
      <div class="flex gap-2 mb-2">
        <ElInput v-model="filterKeyword" placeholder="按名称搜索" clearable size="small" style="width:200px" />
        <ElSelect v-model="filterFenLei" placeholder="按远端分类筛选" clearable size="small" style="width:160px">
          <ElOption v-for="fl in remoteFenLeiOptions" :key="fl" :label="fl" :value="fl" />
        </ElSelect>
        <ElButton size="small" @click="selectAllFiltered">全选筛选结果（{{ filteredRemote.length }}）</ElButton>
        <ElButton size="small" @click="clearSelection">清空选择</ElButton>
        <span class="ml-auto text-sm text-g-500 self-center">
          已选 <b class="text-theme">{{ selectedItems.length }}</b> / {{ remoteList.length }}
        </span>
      </div>

      <ElTable v-loading="fetching" :data="filteredRemote" max-height="450"
        ref="remoteTableRef" row-key="cid"
        @selection-change="(rows: any[]) => selectedItems = rows"
        :show-overflow-tooltip="true" border>
        <ElTableColumn type="selection" width="40" reserve-selection />
        <ElTableColumn prop="cid" label="远端CID" width="80" />
        <ElTableColumn prop="name" label="平台名称" min-width="200" />
        <ElTableColumn prop="price" label="远端价格" width="90" />
        <ElTableColumn prop="fenlei" label="远端分类" width="80" />
        <ElTableColumn prop="noun" label="对接编号" width="90" />
      </ElTable>
      <template #footer>
        <ElButton @click="importVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="importing" @click="handleImport" :disabled="!selectedItems.length">
          导入选中 ({{ selectedItems.length }})
        </ElButton>
      </template>
    </ElDialog>

    <ElDialog v-model="dialogVisible" :title="editId ? '编辑平台' : '添加平台'" width="560px" align-center>
      <ElForm :model="form" label-width="100px">
        <ElFormItem label="平台名称"><ElInput v-model="form.name" /></ElFormItem>
        <ElFormItem label="定价"><ElInput v-model="form.price" type="number" /></ElFormItem>
        <ElFormItem label="查询参数"><ElInput v-model="form.getnoun" placeholder="对接货源的平台编号" /></ElFormItem>
        <ElFormItem label="对接参数"><ElInput v-model="form.noun" placeholder="对接货源的平台编号" /></ElFormItem>
        <ElFormItem label="查询货源">
          <ElSelect v-model="form.queryplat" placeholder="选择查询用的货源" style="width:100%" filterable>
            <ElOption v-for="h in huoYuanList" :key="h.hid" :label="`[${h.pt}] ${h.name}`" :value="String(h.hid)" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="对接货源">
          <ElSelect v-model="form.docking" placeholder="选择下单用的货源" style="width:100%" filterable>
            <ElOption label="自营 (不对接)" value="0" />
            <ElOption v-for="h in huoYuanList" :key="h.hid" :label="`[${h.pt}] ${h.name}`" :value="String(h.hid)" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="费率运算">
          <ElSelect v-model="form.yunsuan">
            <ElOption label="乘法 (*)" value="*" />
            <ElOption label="加法 (+)" value="+" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="分类">
          <ElSelect v-model="form.fenlei" placeholder="选择分类" style="width:100%" clearable>
            <ElOption v-for="f in fenLeiList" :key="f.id" :label="f.name" :value="String(f.id)" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="说明"><ElInput v-model="form.content" type="textarea" :rows="2" /></ElFormItem>
        <ElFormItem label="排序"><ElInput v-model="form.sort" type="number" /></ElFormItem>
        <ElFormItem label="状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
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
  import { ref, onMounted, computed } from 'vue'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import {
    getClassList, addClass, updateClass, deleteClass,
    getHuoYuanList, getFenLeiList,
    fetchRemoteClass, batchImportClass,
    batchDeleteClass
  } from '@/api/wk'

  defineOptions({ name: 'WkAdminClasses' })

  const dialogVisible = ref(false)
  const submitting = ref(false)
  const editId = ref<number | null>(null)
  const form = ref<Partial<Api.Class.Item>>({ yunsuan: '*', status: 1, sort: 10 })
  const huoYuanList = ref<any[]>([])
  const fenLeiList = ref<any[]>([])
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.cid)
  }

  const handleBatchDelete = async () => {
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 条平台？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteClass(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 条`)
    selectedIds.value = []
    refreshData()
  }

  // 批量导入相关
  const importVisible = ref(false)
  const fetching = ref(false)
  const importing = ref(false)
  const importForm = ref({
    hid: 0, fenlei: '', price: '1.00',
    use_remote: true, markup_mode: 'multiply', price_markup: '1.2'
  })
  const remoteList = ref<any[]>([])
  const selectedItems = ref<any[]>([])
  const remoteTableRef = ref<any>()
  const filterKeyword = ref('')
  const filterFenLei = ref('')

  const filteredRemote = computed(() => {
    let list = remoteList.value
    if (filterKeyword.value.trim()) {
      const kw = filterKeyword.value.trim().toLowerCase()
      list = list.filter(r => r.name?.toLowerCase().includes(kw))
    }
    if (filterFenLei.value) {
      list = list.filter(r => r.fenlei === filterFenLei.value)
    }
    return list
  })

  const remoteFenLeiOptions = computed(() => {
    const set = new Set<string>()
    remoteList.value.forEach(r => r.fenlei && set.add(r.fenlei))
    return Array.from(set).sort()
  })

  const selectAllFiltered = () => {
    filteredRemote.value.forEach(r => remoteTableRef.value?.toggleRowSelection(r, true))
  }
  const clearSelection = () => {
    remoteTableRef.value?.clearSelection()
  }

  onMounted(async () => {
    const [hy, fl]: [any, any[]] = await Promise.all([
      getHuoYuanList({ current: 1, size: 200 }),
      getFenLeiList()
    ])
    huoYuanList.value = hy.list || []
    fenLeiList.value = fl || []
  })

  const handleFetch = async () => {
    if (!importForm.value.hid) return
    fetching.value = true
    try {
      const res = await fetchRemoteClass(importForm.value.hid)
      remoteList.value = res || []
      ElMessage.success(`拉取到 ${remoteList.value.length} 个平台`)
    } finally { fetching.value = false }
  }

  const handleImport = async () => {
    importing.value = true
    try {
      const res = await batchImportClass({
        hid: importForm.value.hid,
        fenlei: importForm.value.fenlei,
        price: importForm.value.price,
        use_remote: importForm.value.use_remote,
        markup_mode: importForm.value.markup_mode,
        price_markup: importForm.value.price_markup,
        items: selectedItems.value
      })
      ElMessage.success(`导入成功 ${res.imported} 个，跳过 ${res.total - res.imported} 个`)
      importVisible.value = false
      remoteList.value = []
      selectedItems.value = []
      refreshData()
    } finally { importing.value = false }
  }

  const { columns, columnChecks, data, loading, pagination,
    handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: getClassList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'cid', label: 'ID', width: 80 },
        { prop: 'name', label: '平台名称', minWidth: 200 },
        { prop: 'price', label: '定价', width: 100 },
        { prop: 'yunsuan', label: '费率', width: 90 },
        { prop: 'content', label: '说明', minWidth: 200, showOverflowTooltip: true },
        {
          prop: 'status', label: '状态', width: 90,
          formatter: (row: any) => h(ElTag, { type: row.status === 1 ? 'success' : 'info', size: 'small' },
            () => row.status === 1 ? '上架' : '下架')
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

  const openDialog = (row?: Api.Class.Item) => {
    editId.value = row?.cid ?? null
    form.value = row ? { ...row } : { yunsuan: '*', status: 1, sort: 10 }
    dialogVisible.value = true
  }

  const handleSave = async () => {
    submitting.value = true
    try {
      if (editId.value) {
        await updateClass(editId.value, form.value)
      } else {
        await addClass(form.value)
      }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleDelete = async (row: Api.Class.Item) => {
    await ElMessageBox.confirm(`确认删除平台「${row.name}」？`, '提示', { type: 'warning' })
    await deleteClass(row.cid)
    ElMessage.success('删除成功')
    refreshData()
  }
</script>

<template>
  <div class="admin-users-page art-full-height">
    <ElCard shadow="never" class="mb-3">
      <ElForm :model="searchForm" inline>
        <ElFormItem>
          <ElSelect v-model="searchForm.type" style="width:110px">
            <ElOption label="UID" value="1" />
            <ElOption label="用户名" value="2" />
            <ElOption label="邀请码" value="3" />
            <ElOption label="昵称" value="4" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElInput v-model="searchForm.qq" placeholder="搜索..." clearable style="width:200px" />
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
          <ElButton type="primary" @click="addVisible = true">添加用户</ElButton>
          <ElButton type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
            批量删除 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @selection-change="onSelectionChange"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <!-- 添加用户 -->
    <ElDialog v-model="addVisible" title="添加用户" width="460px" align-center>
      <ElForm :model="addForm" label-width="80px">
        <ElFormItem label="昵称"><ElInput v-model="addForm.name" /></ElFormItem>
        <ElFormItem label="账号"><ElInput v-model="addForm.user" /></ElFormItem>
        <ElFormItem label="密码"><ElInput v-model="addForm.pass" show-password /></ElFormItem>
        <ElFormItem label="等级">
          <ElSelect v-model="addForm.addprice" style="width:100%">
            <ElOption v-for="d in dengjiList" :key="d.id" :label="`${d.name} [${d.rate}]`" :value="d.rate" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="addVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleAdd">确认</ElButton>
      </template>
    </ElDialog>

    <!-- 充值弹窗 -->
    <ElDialog v-model="czVisible" title="账户充值" width="400px" align-center>
      <ElForm label-width="80px">
        <ElFormItem label="充值金额">
          <ElInput v-model="czForm.money" type="number" placeholder="请输入金额" />
        </ElFormItem>
        <ElFormItem label="用户">{{ czForm.name }}</ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="czVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleRecharge">确认</ElButton>
      </template>
    </ElDialog>

    <!-- 修改等级 -->
    <ElDialog v-model="levelVisible" title="修改等级" width="400px" align-center>
      <ElForm label-width="80px">
        <ElFormItem label="等级">
          <ElSelect v-model="levelForm.addprice" style="width:100%">
            <ElOption v-for="d in dengjiList" :key="d.id" :label="`${d.name} [${d.rate}]`" :value="d.rate" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="levelVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleSetLevel">确认</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import {
    getAdminUsers, adminAddUser, adminRecharge, adminBanUser,
    adminResetPass, adminSetLevel, getDengjiList,
    batchDeleteUsers
  } from '@/api/wk'

  defineOptions({ name: 'WkAdminUsers' })

  type UserItem = Api.User.Item

  const addVisible = ref(false)
  const czVisible = ref(false)
  const levelVisible = ref(false)
  const submitting = ref(false)
  const dengjiList = ref<Api.Dengji.Item[]>([])
  const addForm = ref({ name: '', user: '', pass: '', addprice: 1 })
  const czForm = ref({ uid: 0, name: '', money: '' })
  const levelForm = ref({ uid: 0, addprice: 1 })
  const searchForm = ref({ type: '1', qq: '' })
  const selectedIds = ref<number[]>([])

  const onSelectionChange = (rows: any[]) => {
    selectedIds.value = rows.map(r => r.uid)
  }

  const handleBatchDelete = async () => {
    if (selectedIds.value.includes(1)) {
      ElMessage.warning('不能删除站长账号 (UID=1)')
      return
    }
    await ElMessageBox.confirm(`确认硬删除选中的 ${selectedIds.value.length} 个用户？此操作不可恢复！`, '提示', { type: 'warning' })
    const res = await batchDeleteUsers(selectedIds.value)
    ElMessage.success(`已删除 ${res.deleted} 个`)
    selectedIds.value = []
    refreshData()
  }

  const { columns, columnChecks, data, loading, pagination, searchParams,
    resetSearchParams, handleSizeChange, handleCurrentChange, refreshData, getData } = useTable({
    core: {
      apiFn: getAdminUsers,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'selection', width: 40 },
        { type: 'index', width: 60, label: '序号' },
        { prop: 'uid', label: 'UID', width: 80 },
        { prop: 'user', label: '账号', width: 140 },
        { prop: 'name', label: '昵称', minWidth: 140 },
        { prop: 'money', label: '余额', width: 110, formatter: (row: any) => `¥${row.money}` },
        { prop: 'addprice', label: '费率', width: 90 },
        {
          prop: 'active', label: '状态', width: 90,
          formatter: (row: any) => h(ElTag, {
            type: row.active === '1' ? 'success' : 'danger', size: 'small',
            style: 'cursor:pointer',
            onClick: () => handleBan(row)
          }, () => row.active === '1' ? '正常' : '封禁')
        },
        { prop: 'endtime', label: '最近在线', width: 180 },
        {
          prop: 'operation', label: '操作', width: 160, fixed: 'right',
          formatter: (row: any) => h('div', { class: 'flex gap-1' }, [
            h(ArtButtonTable, { type: 'edit', title: '充值', onClick: () => openCZ(row) }),
            h(ArtButtonTable, { type: 'edit', title: '等级', onClick: () => openLevel(row) }),
            h(ArtButtonTable, { type: 'delete', title: '重置密码', onClick: () => handleResetPass(row) })
          ])
        }
      ]
    }
  })

  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openCZ = (row: UserItem) => {
    czForm.value = { uid: row.uid, name: row.name, money: '' }
    czVisible.value = true
  }

  const openLevel = (row: UserItem) => {
    levelForm.value = { uid: row.uid, addprice: row.addprice }
    levelVisible.value = true
  }

  const handleAdd = async () => {
    submitting.value = true
    try {
      await adminAddUser(addForm.value)
      ElMessage.success('添加成功')
      addVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleRecharge = async () => {
    submitting.value = true
    try {
      await adminRecharge(czForm.value.uid, parseFloat(czForm.value.money))
      ElMessage.success('充值成功')
      czVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleSetLevel = async () => {
    submitting.value = true
    try {
      await adminSetLevel(levelForm.value.uid, levelForm.value.addprice)
      ElMessage.success('修改成功')
      levelVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleBan = async (row: UserItem) => {
    const newActive = row.active === '1' ? '0' : '1'
    const label = newActive === '0' ? '封禁' : '解封'
    await ElMessageBox.confirm(`确认${label}该用户？`, '提示', { type: 'warning' })
    await adminBanUser(row.uid, newActive)
    ElMessage.success(`${label}成功`)
    refreshData()
  }

  const handleResetPass = async (row: UserItem) => {
    const { value } = await ElMessageBox.prompt('请输入新密码', '重置密码', {
      confirmButtonText: '确认', cancelButtonText: '取消'
    })
    await adminResetPass(row.uid, value)
    ElMessage.success('重置成功')
  }

  onMounted(async () => {
    dengjiList.value = await getDengjiList()
  })
</script>

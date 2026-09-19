<template>
  <div class="agents-page art-full-height">
    <ElCard shadow="never" class="mb-3">
      <ElForm :model="searchForm" inline>
        <ElFormItem>
          <ElSelect v-model="searchForm.type" style="width:110px">
            <ElOption label="UID" value="1" />
            <ElOption label="账号" value="2" />
            <ElOption label="邀请码" value="3" />
            <ElOption label="昵称" value="4" />
            <ElOption label="费率" value="5" />
            <ElOption label="余额" value="6" />
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
          <ElButton type="primary" @click="openAdd">开户</ElButton>
        </template>
      </ArtTableHeader>
      <ArtTable :loading="loading" :data="(data as any)" :columns="columns" :pagination="pagination"
        @pagination:size-change="handleSizeChange" @pagination:current-change="handleCurrentChange" />
    </ElCard>

    <!-- 开户弹窗 -->
    <ElDialog v-model="addVisible" title="开户" width="460px" align-center>
      <ElForm :model="addForm" label-width="80px">
        <ElFormItem label="昵称"><ElInput v-model="addForm.name" placeholder="请输入昵称" /></ElFormItem>
        <ElFormItem label="QQ账号"><ElInput v-model="addForm.user" placeholder="请输入QQ号" /></ElFormItem>
        <ElFormItem label="密码"><ElInput v-model="addForm.pass" show-password placeholder="至少6位" /></ElFormItem>
        <ElFormItem label="费率">
          <ElInput v-model="addForm.addprice" type="number" placeholder="0.05的倍数，且不低于自己的费率" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="addVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleAdd">确认</ElButton>
      </template>
    </ElDialog>

    <!-- 充值弹窗 -->
    <ElDialog v-model="czVisible" title="下级充值" width="400px" align-center>
      <ElForm label-width="80px">
        <ElFormItem label="下级">{{ czForm.name }}</ElFormItem>
        <ElFormItem label="充值金额">
          <ElInput v-model="czForm.money" type="number" placeholder="请输入金额（按其费率换算扣费）" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="czVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleRecharge">确认</ElButton>
      </template>
    </ElDialog>

    <!-- 改价弹窗 -->
    <ElDialog v-model="priceVisible" title="修改费率" width="400px" align-center>
      <ElForm label-width="80px">
        <ElFormItem label="下级">{{ priceForm.name }}</ElFormItem>
        <ElFormItem label="新费率">
          <ElInput v-model="priceForm.addprice" type="number" placeholder="只能上调，不能低于当前费率" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="priceVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleSetPrice">确认</ElButton>
      </template>
    </ElDialog>

    <!-- 设邀请码弹窗 -->
    <ElDialog v-model="yqmVisible" title="设置邀请码" width="400px" align-center>
      <ElForm label-width="80px">
        <ElFormItem label="下级">{{ yqmForm.name }}</ElFormItem>
        <ElFormItem label="邀请码">
          <ElInput v-model="yqmForm.yqm" placeholder="请输入邀请码" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="yqmVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleSetYQM">确认</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, h } from 'vue'
  import { ElMessage, ElMessageBox, ElTag } from 'element-plus'
  import { useTable } from '@/hooks/core/useTable'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import {
    getProxyUsers, proxyAddUser, proxyRecharge, proxySetPrice, proxyResetPass,
    proxyBanUser, proxySetYQM, proxyOpenKey
  } from '@/api/wk'

  defineOptions({ name: 'WkAgents' })

  const addVisible = ref(false)
  const czVisible = ref(false)
  const priceVisible = ref(false)
  const yqmVisible = ref(false)
  const submitting = ref(false)
  const addForm = ref({ name: '', user: '', pass: '', addprice: '' })
  const czForm = ref({ uid: 0, name: '', money: '' })
  const priceForm = ref({ uid: 0, name: '', addprice: '' })
  const yqmForm = ref({ uid: 0, name: '', yqm: '' })
  const searchForm = ref({ type: '1', qq: '' })

  const copyText = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
    } catch {
      const input = document.createElement('textarea')
      input.value = text
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      document.body.removeChild(input)
    }
    ElMessage.success('已复制')
  }

  const { columns, columnChecks, data, loading, pagination, searchParams,
    resetSearchParams, handleSizeChange, handleCurrentChange, refreshData, getData } = useTable({
    core: {
      apiFn: getProxyUsers,
      apiParams: { current: 1, size: 15 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'uid', label: 'UID', width: 80 },
        { prop: 'user', label: '账号', width: 140 },
        { prop: 'name', label: '昵称', minWidth: 120 },
        { prop: 'money', label: '余额', width: 110, formatter: (row: any) => `¥${row.money}` },
        { prop: 'addprice', label: '费率', width: 90 },
        { prop: 'yqm', label: '邀请码', width: 110, formatter: (row: any) => row.yqm || '-' },
        {
          prop: 'key', label: 'API', width: 140,
          formatter: (row: any) => {
            if (!row.key || row.key === '0') return '未开通'
            return h('span', {
              style: 'cursor:pointer;color:var(--el-color-primary)',
              title: '点击复制',
              onClick: () => copyText(row.key)
            }, () => `${row.key.slice(0, 8)}...`)
          }
        },
        {
          prop: 'active', label: '状态', width: 90,
          formatter: (row: any) => h(ElTag, {
            type: row.active === '1' ? 'success' : 'danger', size: 'small'
          }, () => row.active === '1' ? '正常' : '封禁')
        },
        {
          prop: 'operation', label: '操作', width: 260, fixed: 'right',
          formatter: (row: any) => h('div', { class: 'flex gap-1' }, [
            h(ArtButtonTable, { type: 'edit', title: '充值', onClick: () => openCZ(row) }),
            h(ArtButtonTable, { type: 'edit', title: '改价', onClick: () => openPrice(row) }),
            h(ArtButtonTable, { type: 'delete', title: '重置密码', onClick: () => handleResetPass(row) }),
            h(ArtButtonTable, { type: 'edit', title: row.active === '1' ? '封禁' : '解封', onClick: () => handleBan(row) }),
            h(ArtButtonTable, { type: 'edit', title: '邀请码', onClick: () => openYQM(row) }),
            h(ArtButtonTable, { type: 'edit', title: '开API', onClick: () => handleOpenKey(row) })
          ])
        }
      ]
    }
  })

  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }

  const openAdd = () => {
    addForm.value = { name: '', user: '', pass: '', addprice: '' }
    addVisible.value = true
  }

  const openCZ = (row: any) => {
    czForm.value = { uid: row.uid, name: `${row.name} (${row.user})`, money: '' }
    czVisible.value = true
  }

  const openPrice = (row: any) => {
    priceForm.value = { uid: row.uid, name: `${row.name} (当前费率 ${row.addprice})`, addprice: '' }
    priceVisible.value = true
  }

  const openYQM = (row: any) => {
    yqmForm.value = { uid: row.uid, name: `${row.name} (${row.user})`, yqm: row.yqm || '' }
    yqmVisible.value = true
  }

  const handleAdd = async () => {
    if (!/^\d{5,11}$/.test(addForm.value.user)) {
      ElMessage.warning('账号必须为QQ号格式')
      return
    }
    if (addForm.value.pass.length < 6) {
      ElMessage.warning('密码至少6位')
      return
    }
    const price = parseFloat(addForm.value.addprice)
    if (isNaN(price) || Math.round(price * 100) % 5 !== 0) {
      ElMessage.warning('费率必须为0.05的倍数')
      return
    }
    submitting.value = true
    try {
      await proxyAddUser({ ...addForm.value, addprice: price })
      ElMessage.success('开户成功')
      addVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleRecharge = async () => {
    const money = parseFloat(czForm.value.money)
    if (isNaN(money) || money <= 0) {
      ElMessage.warning('请输入正确的金额')
      return
    }
    submitting.value = true
    try {
      await proxyRecharge(czForm.value.uid, money)
      ElMessage.success('充值成功')
      czVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleSetPrice = async () => {
    const addprice = parseFloat(priceForm.value.addprice)
    if (isNaN(addprice) || Math.round(addprice * 100) % 5 !== 0) {
      ElMessage.warning('费率必须为0.05的倍数')
      return
    }
    submitting.value = true
    try {
      await proxySetPrice(priceForm.value.uid, addprice)
      ElMessage.success('修改成功')
      priceVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleResetPass = async (row: any) => {
    await ElMessageBox.confirm(`确认将 ${row.name} 的密码重置为 123456？`, '重置密码', { type: 'warning' })
    await proxyResetPass(row.uid)
    ElMessage.success('重置成功')
  }

  const handleBan = async (row: any) => {
    const newActive = row.active === '1' ? '0' : '1'
    const label = newActive === '0' ? '封禁' : '解封'
    await ElMessageBox.confirm(`确认${label}该下级？`, '提示', { type: 'warning' })
    await proxyBanUser(row.uid, newActive)
    ElMessage.success(`${label}成功`)
    refreshData()
  }

  const handleSetYQM = async () => {
    if (!yqmForm.value.yqm) {
      ElMessage.warning('请输入邀请码')
      return
    }
    submitting.value = true
    try {
      await proxySetYQM(yqmForm.value.uid, yqmForm.value.yqm)
      ElMessage.success('设置成功')
      yqmVisible.value = false
      refreshData()
    } finally { submitting.value = false }
  }

  const handleOpenKey = async (row: any) => {
    const res = await proxyOpenKey(row.uid)
    await ElMessageBox.alert(`API Key：${res.key}`, '开通成功', {
      confirmButtonText: '复制并关闭',
      type: 'success'
    })
    copyText(res.key)
  }
</script>

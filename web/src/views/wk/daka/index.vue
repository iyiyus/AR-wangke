<template>
  <div class="daka-page">
    <!-- 平台单价 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">实习打卡</h4>
        <span class="text-g-500 text-sm">余额：<span class="text-theme font-bold">¥ {{ userInfo.money?.toFixed(2) }}</span></span>
      </div>
      <div class="grid grid-cols-3 md:grid-cols-6 lg:grid-cols-11 gap-2">
        <div v-for="(p, key) in prices" :key="key"
          class="px-2 py-2 text-center rounded-md border cursor-pointer duration-200 hover:border-theme hover:bg-theme/5"
          :class="{ '!border-theme bg-theme/10': form.platform === key }"
          @click="selectPlatform(key)">
          <div class="text-xs truncate">{{ p.name }}</div>
          <div class="text-theme font-bold text-sm mt-1">¥{{ p.price }}</div>
        </div>
      </div>
      <div v-if="form.platform" class="mt-3 px-4 py-3 bg-theme/10 rounded-md text-sm text-g-700">
        已选：<b>{{ prices[form.platform]?.name }}</b>（单价 ¥{{ prices[form.platform]?.price }}/天）
      </div>
    </div>

    <!-- 添加/修改订单 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">{{ editMode ? '修改订单' : '添加订单' }}</h4>
        <ElButton v-if="editMode" size="small" link @click="resetForm">返回添加</ElButton>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <ElFormItem label="手机号" class="!mb-0">
          <ElInput v-model="form.phone" placeholder="登录手机号" :disabled="editMode" />
        </ElFormItem>
        <ElFormItem label="登录密码" class="!mb-0">
          <ElInput v-model="form.password" placeholder="登录密码" show-password />
        </ElFormItem>
        <ElFormItem label="姓名" class="!mb-0">
          <ElInput v-model="form.name" placeholder="真实姓名（选填）" />
        </ElFormItem>
        <ElFormItem label="打卡地址" class="!mb-0 md:col-span-2">
          <ElInput v-model="form.address" placeholder="打卡定位地址" />
        </ElFormItem>
        <ElFormItem label="运行方式" v-if="form.platform === 'xyb'" class="!mb-0">
          <ElRadioGroup v-model="form.runType">
            <ElRadio :value="1">抓包</ElRadio>
            <ElRadio :value="2">账密1号</ElRadio>
            <ElRadio :value="3">绑微1号</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="上班打卡时间" class="!mb-0">
          <ElTimePicker v-model="form.up_check_time" value-format="HH:mm:ss" placeholder="选择时间" style="width: 100%" />
        </ElFormItem>
        <ElFormItem label="下班打卡" class="!mb-0">
          <div class="flex items-center gap-2">
            <ElSwitch v-model="formOther.down_check" />
            <ElTimePicker v-if="formOther.down_check" v-model="form.down_check_time" value-format="HH:mm:ss"
              placeholder="下班时间" style="width: 140px" />
          </div>
        </ElFormItem>
        <ElFormItem label="打卡周期" class="!mb-0">
          <ElCheckboxGroup v-model="form.check_week">
            <ElCheckbox v-for="(w, i) in weekNames" :key="i" :value="i">{{ w }}</ElCheckbox>
          </ElCheckboxGroup>
        </ElFormItem>
        <ElFormItem label="结束时间" class="!mb-0">
          <ElDatePicker v-model="form.end_time" type="date" value-format="YYYY-MM-DD" placeholder="选择结束日期" style="width: 100%" />
        </ElFormItem>
        <ElFormItem label="日报" class="!mb-0">
          <ElSwitch v-model="form.day_paper" active-value="1" inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="周报" class="!mb-0">
          <ElSwitch v-model="form.week_paper" active-value="1" inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="月报" class="!mb-0">
          <ElSwitch v-model="form.month_paper" active-value="1" inactive-value="0" />
        </ElFormItem>
      </div>
      <div class="mt-4 flex justify-end">
        <ElButton type="primary" :loading="submitting" :disabled="!form.platform" @click="handleSubmit">
          {{ editMode ? '确认修改' : '确认下单' }}
        </ElButton>
      </div>
    </div>

    <!-- 订单列表 -->
    <div class="art-card p-5">
      <div class="flex-cb mb-4">
        <h4 class="font-bold">我的打卡订单</h4>
        <ElButton size="small" :loading="loading" @click="loadOrders">刷新</ElButton>
      </div>
      <ElTable :data="orders" v-loading="loading" stripe>
        <ElTableColumn prop="id" label="ID" width="60" />
        <ElTableColumn label="平台" width="100">
          <template #default="{ row }">{{ platformName(row.platform) }}</template>
        </ElTableColumn>
        <ElTableColumn prop="phone" label="手机号" width="120" />
        <ElTableColumn prop="address" label="打卡地址" min-width="140" show-overflow-tooltip />
        <ElTableColumn label="打卡时间" width="130">
          <template #default="{ row }">{{ row.up_check_time || '-' }}<template v-if="row.down_check_time"> / {{ row.down_check_time }}</template></template>
        </ElTableColumn>
        <ElTableColumn label="周期" width="90">
          <template #default="{ row }">{{ weekText(row.check_week) }}</template>
        </ElTableColumn>
        <ElTableColumn prop="end_time" label="结束时间" width="110" />
        <ElTableColumn label="剩余" width="70">
          <template #default="{ row }">
            <span :class="row.expired ? 'text-g-400' : 'text-theme font-bold'">{{ row.expired ? '已到期' : row.day + '天' }}</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="报告" width="100">
          <template #default="{ row }">
            {{ row.day_paper ? '日' : '' }}{{ row.week_paper ? '周' : '' }}{{ row.month_paper ? '月' : '' }}
            <span v-if="!row.day_paper && !row.week_paper && !row.month_paper" class="text-g-400">无</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <ElButton size="small" link type="primary" @click="nowCheck(row)">打卡</ElButton>
            <ElButton size="small" link type="primary" @click="openBuPapers(row)">补报告</ElButton>
            <ElButton size="small" link type="primary" @click="openEdit(row)">改单</ElButton>
            <ElButton size="small" link type="primary" @click="openLog(row)">日志</ElButton>
            <ElButton size="small" link type="danger" @click="delOrder(row)">删单</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
    </div>

    <!-- 补报告弹窗 -->
    <ElDialog v-model="buPapersVisible" title="补报告" width="420">
      <ElForm label-width="90px">
        <ElFormItem label="类型">
          <ElSelect v-model="buPapersForm.levelName" style="width: 100%">
            <ElOption label="日报" value="日报" />
            <ElOption label="周报" value="周报" />
            <ElOption label="月报" value="月报" />
            <ElOption label="总结" value="总结" />
            <ElOption label="上班打卡" value="上班打卡" />
            <ElOption label="下班打卡" value="下班打卡" />
            <ElOption label="上下班打卡" value="上下班打卡" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="开始时间"><ElDatePicker v-model="buPapersForm.startTime" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></ElFormItem>
        <ElFormItem label="结束时间"><ElDatePicker v-model="buPapersForm.endTime" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="buPapersVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="buPapersLoading" @click="submitBuPapers">提交</ElButton>
      </template>
    </ElDialog>

    <!-- 打卡日志弹窗 -->
    <ElDialog v-model="logVisible" title="打卡日志" width="560">
      <div v-if="logLoading" class="py-6 text-center text-g-500">加载中...</div>
      <pre v-else class="log-body">{{ logText }}</pre>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { dakaApi, getUserInfo } from '@/api/wk'

  defineOptions({ name: 'WkDaka' })

  const userInfo = ref<Partial<Api.Auth.UserInfo>>({})
  const prices = ref<Record<string, any>>({})
  const weekNames = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const platforms: Record<string, string> = {
    zxjy: '职校家园', qzt: '黔职通', xyb: '校友帮', gxy: '工学云', xxy: '习讯云',
    xxt: '学习通', hzj: '慧职教', gxzy: '广西职业', jxzhjy: '江西智慧教育', cxy: '成学云', bx: '博行'
  }

  const defaultForm = () => ({
    platform: '', phone: '', password: '', name: '', address: '',
    up_check_time: '', down_check_time: '', check_week: [] as number[],
    end_time: '', day_paper: '0', week_paper: '0', month_paper: '0', runType: 1
  })
  const form = reactive(defaultForm())
  const formOther = reactive({ down_check: false })
  const editMode = ref(false)
  const editId = ref(0)

  const orders = ref<any[]>([])
  const loading = ref(false)
  const submitting = ref(false)

  const buPapersVisible = ref(false)
  const buPapersLoading = ref(false)
  const buPapersForm = reactive({ id: 0, levelName: '日报', startTime: '', endTime: '' })
  const logVisible = ref(false)
  const logLoading = ref(false)
  const logText = ref('')

  const platformName = (k: string) => platforms[k] || k

  const weekText = (s: string) => {
    if (!s) return '-'
    return s.split(',').map((n: string) => weekNames[Number(n)]).filter(Boolean).join('、') || '-'
  }

  const selectPlatform = (key: string) => {
    form.platform = key
    formOther.down_check = false
  }

  const loadPrices = async () => {
    try {
      const res: any = await dakaApi('price')
      if (res?.code === 0) prices.value = res.data
    } catch (e) { /* 静默 */ }
  }

  const loadOrders = async () => {
    loading.value = true
    try {
      const res: any = await dakaApi('order')
      if (res?.code === 0) orders.value = res.data || []
      else ElMessage.error(res?.msg || '加载失败')
    } catch (e: any) {
      ElMessage.error(e?.message || '加载失败')
    } finally {
      loading.value = false
    }
  }

  const resetForm = () => {
    Object.assign(form, defaultForm())
    editMode.value = false
    editId.value = 0
    formOther.down_check = false
  }

  const handleSubmit = async () => {
    if (!form.platform) return ElMessage.warning('请选择平台')
    if (!form.phone || !form.password) return ElMessage.warning('请填写手机号和密码')
    if (!form.address) return ElMessage.warning('请填写打卡地址')
    if (!form.up_check_time && form.platform !== 'zxjy') return ElMessage.warning('请选择上班打卡时间')
    if (!form.check_week.length) return ElMessage.warning('请选择打卡周期')
    if (!form.end_time) return ElMessage.warning('请选择结束时间')
    submitting.value = true
    try {
      const params: Record<string, any> = {
        platform: form.platform,
        phone: form.phone,
        password: form.password,
        name: form.name,
        address: form.address,
        up_check_time: form.platform === 'zxjy' ? '' : form.up_check_time,
        down_check_time: formOther.down_check ? form.down_check_time : '',
        check_week: form.check_week.join(','),
        end_time: form.end_time,
        day_paper: Number(form.day_paper),
        week_paper: Number(form.week_paper),
        month_paper: Number(form.month_paper),
        runType: form.runType
      }
      let res: any
      if (editMode.value) {
        params.id = editId.value
        res = await dakaApi('edit', { form: params })
      } else {
        res = await dakaApi('add', { form: params })
      }
      if (res?.code === 0) {
        ElMessage.success(res.msg || (editMode.value ? '修改成功' : '添加成功'))
        resetForm()
        loadOrders()
        userInfo.value = await getUserInfo()
      } else {
        ElMessage.error(res?.msg || '操作失败')
      }
    } catch (e: any) {
      ElMessage.error(e?.message || '操作失败')
    } finally {
      submitting.value = false
    }
  }

  const nowCheck = async (row: any) => {
    await ElMessageBox.confirm(`立即为「${platformName(row.platform)} ${row.phone}」打卡？将按单价扣费。`, '立即打卡', { type: 'warning' })
    const res: any = await dakaApi('nowCheck', { id: row.id, platform: row.platform })
    if (res?.code === 0) {
      ElMessage.success(res.msg || '打卡成功')
      loadOrders()
      userInfo.value = await getUserInfo()
    } else {
      ElMessage.error(res?.msg || '打卡失败')
    }
  }

  const openBuPapers = (row: any) => {
    buPapersForm.id = row.id
    buPapersForm.levelName = '日报'
    buPapersForm.startTime = ''
    buPapersForm.endTime = ''
    buPapersVisible.value = true
  }

  const submitBuPapers = async () => {
    if (!buPapersForm.startTime || !buPapersForm.endTime) return ElMessage.warning('请选择时间范围')
    buPapersLoading.value = true
    try {
      const res: any = await dakaApi('buPapers', {
        id: buPapersForm.id, levelName: buPapersForm.levelName,
        startTime: buPapersForm.startTime, endTime: buPapersForm.endTime
      })
      if (res?.code === 0) {
        ElMessage.success('补报告成功')
        buPapersVisible.value = false
      } else {
        ElMessage.error(res?.msg || '补报告失败')
      }
    } catch (e: any) {
      ElMessage.error(e?.message || '补报告失败')
    } finally {
      buPapersLoading.value = false
    }
  }

  const openEdit = (row: any) => {
    editMode.value = true
    editId.value = row.id
    form.platform = row.platform
    form.phone = row.phone
    form.password = row.password
    form.name = row.name || ''
    form.address = row.address || ''
    form.up_check_time = row.up_check_time || ''
    form.down_check_time = row.down_check_time || ''
    formOther.down_check = !!row.down_check_time
    form.check_week = row.check_week ? row.check_week.split(',').map(Number) : []
    form.end_time = row.end_time || ''
    form.day_paper = String(row.day_paper || 0)
    form.week_paper = String(row.week_paper || 0)
    form.month_paper = String(row.month_paper || 0)
  }

  const openLog = async (row: any) => {
    logVisible.value = true
    logLoading.value = true
    logText.value = ''
    try {
      const res: any = await dakaApi('getLog', { id: row.id })
      logText.value = JSON.stringify(res, null, 2)
    } catch (e: any) {
      logText.value = e?.message || '日志获取失败'
    } finally {
      logLoading.value = false
    }
  }

  const delOrder = async (row: any) => {
    await ElMessageBox.confirm('删除后源台订单将取消，剩余天数将退款。确认删除？', '删除订单', { type: 'warning' })
    const res: any = await dakaApi('del', { id: row.id })
    if (res?.code === 0) {
      ElMessage.success(res.msg || '删除成功')
      loadOrders()
      userInfo.value = await getUserInfo()
    } else {
      ElMessage.error(res?.msg || '删除失败')
    }
  }

  onMounted(async () => {
    try { userInfo.value = await getUserInfo() } catch (e) { /* 静默 */ }
    loadPrices()
    loadOrders()
  })
</script>

<style scoped>
  .log-body {
    max-height: 420px;
    overflow: auto;
    white-space: pre-wrap;
    word-break: break-all;
    background: #f7f8fa;
    padding: 12px;
    border-radius: 8px;
    font-size: 12px;
  }
</style>

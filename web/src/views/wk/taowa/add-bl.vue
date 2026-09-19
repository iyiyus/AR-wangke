<template>
  <div class="taowa-addbl-page">
    <!-- ① 选择模板 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">① 选择病历模板</h4>
        <span class="text-g-500 text-sm">余额：<span class="text-theme font-bold">¥ {{ userInfo.money?.toFixed(2) }}</span></span>
      </div>
      <div v-if="!configured" class="px-4 py-3 rounded-md bg-amber-50 text-amber-600 text-sm mb-3">
        对接尚未配置，请联系管理员在「系统设置 - 实习助手对接」中填写源台 UID / KEY。
      </div>
      <div class="flex-cb mb-3">
        <ElButton type="primary" plain @click="pdfVisible = true">预览模板样式</ElButton>
      </div>
      <ElInput v-model="search" placeholder="搜索模板" clearable size="small" style="width: 280px" class="mb-3">
        <template #prefix><ElIcon><Search /></ElIcon></template>
      </ElInput>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 max-h-72 overflow-auto">
        <div v-if="!filteredTemplates.length" class="col-span-full p-6 text-center text-g-500">无匹配模板</div>
        <div v-for="t in filteredTemplates" :key="t.name"
          class="px-3 py-2 cursor-pointer rounded-md border border-g-300 duration-200 hover:border-theme hover:bg-theme/5"
          :class="{ '!border-theme bg-theme/10': form.templateName === t.name }"
          @click="selectTemplate(t)">
          <div class="text-sm truncate">{{ t.name }}</div>
          <div class="flex-cb mt-1">
            <span class="text-xs text-g-500">¥{{ t.price }} 起</span>
          </div>
        </div>
      </div>
      <div v-if="currentTemplate?.content" class="mt-3 px-4 py-3 bg-theme/10 rounded-md text-sm text-g-700">
        {{ currentTemplate.content }}
      </div>
    </div>

    <!-- ② 患者信息 -->
    <div class="art-card p-5 mb-5">
      <h4 class="font-bold mb-3">② 患者信息</h4>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <ElFormItem label="患者姓名" class="!mb-0">
          <ElInput v-model="form.patientName" placeholder="患者真实姓名" />
        </ElFormItem>
        <ElFormItem label="性别" class="!mb-0">
          <ElSelect v-model="form.gender" style="width: 100%">
            <ElOption label="男" value="男" /><ElOption label="女" value="女" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="年龄" class="!mb-0">
          <ElInput v-model="form.age" placeholder="如：28" />
        </ElFormItem>
        <ElFormItem label="科室" class="!mb-0">
          <ElInput v-model="form.department" placeholder="如：内科" />
        </ElFormItem>
        <ElFormItem label="病情描述" class="!mb-0">
          <ElInput v-model="form.condition" placeholder="病情/诊断情况" />
        </ElFormItem>
        <ElFormItem label="治疗方式" class="!mb-0">
          <ElInput v-model="form.treatment" placeholder="治疗/用药情况" />
        </ElFormItem>
        <ElFormItem label="诊断日期" class="!mb-0">
          <ElDatePicker v-model="form.diagnosisDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" style="width: 100%" />
        </ElFormItem>
        <ElFormItem label="就诊医院" class="!mb-0">
          <ElInput v-model="form.hospital" placeholder="医院名称" />
        </ElFormItem>
        <ElFormItem label="QQ邮箱" class="!mb-0">
          <ElInput v-model="form.email" placeholder="用于接收病历（QQ邮箱）" />
        </ElFormItem>
      </div>
    </div>

    <!-- 模板样式预览弹窗 -->
    <ElDialog v-model="pdfVisible" title="模板样式预览" width="900px" align-center>
      <div class="border border-g-300 rounded-md overflow-hidden" style="height: 700px;">
        <iframe src="https://www.documenticloud.cfd/%E6%A8%A1%E6%9D%BF%E6%A0%B7%E5%BC%8F.pdf"
          style="width:100%;height:100%;border:none;" />
      </div>
    </ElDialog>

    <!-- 底部提交 -->
    <div class="art-card p-5 flex-cb sticky bottom-4">
      <div class="text-sm text-g-700">
        {{ currentTemplate ? `模板：${currentTemplate.name}` : '请先选择模板' }}
        <span v-if="currentTemplate" class="ml-3">预计费用：
          <span class="text-theme font-bold text-lg">¥ {{ estimatedFee }}</span>
        </span>
      </div>
      <ElButton type="primary" size="large" :loading="submitting" :disabled="!form.templateName" @click="handleSubmit">
        确认提交
      </ElButton>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { ElMessage, ElIcon } from 'element-plus'
  import { Search } from '@element-plus/icons-vue'
  import { getUserInfo, getTaowaTemplates, addTaowaOrderBL } from '@/api/wk'

  defineOptions({ name: 'WkTaowaAddBl' })

  const userInfo = ref<Partial<Api.Auth.UserInfo>>({})
  const configured = ref(false)
  const templates = ref<any[]>([])
  const search = ref('')
  const currentTemplate = ref<any>(null)
  const form = ref({
    templateName: '', patientName: '', gender: '男', age: '', department: '',
    condition: '', treatment: '', diagnosisDate: '', hospital: '', email: ''
  })
  const submitting = ref(false)
  const pdfVisible = ref(false)

  const estimatedFee = computed(() => {
    if (!currentTemplate.value) return '0.00'
    return (parseFloat(currentTemplate.value.price || '0') * 2).toFixed(2)
  })

  const filteredTemplates = computed(() => {
    if (!search.value.trim()) return templates.value
    const kw = search.value.trim().toLowerCase()
    return templates.value.filter(t => t.name.toLowerCase().includes(kw))
  })

  const selectTemplate = (t: any) => {
    form.value.templateName = t.name
    currentTemplate.value = t
  }

  const handleSubmit = async () => {
    if (!form.value.templateName) return ElMessage.warning('请先选择模板')
    const f = form.value
    if (!f.patientName || !f.age || !f.department || !f.condition || !f.treatment || !f.diagnosisDate || !f.hospital || !f.email) {
      return ElMessage.warning('请填写所有必填字段')
    }
    if (!/^[^\s@]+@qq\.com$/i.test(f.email)) {
      return ElMessage.warning('请输入有效的QQ邮箱地址')
    }
    submitting.value = true
    try {
      const res: any = await addTaowaOrderBL({
        templateName: f.templateName, patientName: f.patientName, gender: f.gender,
        age: f.age, department: f.department, condition: f.condition, treatment: f.treatment,
        diagnosisDate: f.diagnosisDate, hospital: f.hospital, email: f.email,
        price: parseFloat(currentTemplate.value.price || '0')
      })
      ElMessage.success(res?.msg || '提交成功')
      form.value = { templateName: '', patientName: '', gender: '男', age: '', department: '', condition: '', treatment: '', diagnosisDate: '', hospital: '', email: '' }
      currentTemplate.value = null
      userInfo.value = await getUserInfo()
    } catch (e: any) {
      ElMessage.error(e?.message || '提交失败')
    } finally {
      submitting.value = false
    }
  }

  onMounted(async () => {
    try {
      const [info, data]: any[] = await Promise.all([getUserInfo(), getTaowaTemplates()])
      userInfo.value = info
      templates.value = data?.templates || []
      configured.value = data?.configured !== false
    } catch (e) { /* 未配置时静默 */ }
  })
</script>

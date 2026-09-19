<template>
  <div class="taowa-add-page">
    <!-- 下单流程 & 注意事项 -->
    <div class="art-card p-5 mb-5">
      <ElCollapse>
        <ElCollapseItem title="查看下单流程" name="flow">
          <div class="text-sm text-g-700 space-y-2 leading-relaxed">
            <p>填写学生姓名 → 填写学生收件手机号 → 填写详细地址 → 如果学生有学校发的材料，填写寄到工作室的快递单号</p>
            <p>例如：</p>
            <div class="px-4 py-3 bg-g-50 rounded-md font-mono text-sm space-y-1">
              <div>张三</div>
              <div>18888888888</div>
              <div>山西省晋中市榆次区翻斗花园3号楼三单元101号</div>
              <div>+快递单号SF121278183311</div>
            </div>
          </div>
        </ElCollapseItem>
        <ElCollapseItem title="查看注意事项" name="tips">
          <div class="text-sm text-g-700 space-y-2 leading-relaxed">
            <p>1. 可以在材料文件夹里面添加一个文本文档来说明具体盖章要求。</p>
            <p>2. 如果学生不用把材料寄到工作室，就不用填写单号。</p>
            <p>3. 下单时请正确上传文件。</p>
          </div>
        </ElCollapseItem>
      </ElCollapse>
    </div>

    <!-- ① 选择公司 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">① 选择公司</h4>
        <span class="text-g-500 text-sm">余额：<span class="text-theme font-bold">¥ {{ userInfo.money?.toFixed(2) }}</span></span>
      </div>
      <div v-if="!configured" class="px-4 py-3 rounded-md bg-amber-50 text-amber-600 text-sm mb-3">
        对接尚未配置，请联系管理员在「系统设置 - 实习助手对接」中填写源台 UID / KEY。
      </div>
      <ElInput v-model="search" placeholder="搜索公司" clearable size="small" style="width: 280px" class="mb-3">
        <template #prefix><ElIcon><Search /></ElIcon></template>
      </ElInput>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 max-h-72 overflow-auto">
        <div v-if="!filteredCompanies.length" class="col-span-full p-6 text-center text-g-500">无匹配公司</div>
        <div v-for="c in filteredCompanies" :key="c.name"
          class="px-3 py-2 cursor-pointer rounded-md border border-g-300 duration-200 hover:border-theme hover:bg-theme/5"
          :class="{ '!border-theme bg-theme/10': form.companyName === c.name }"
          @click="selectCompany(c)">
          <div class="text-sm truncate">{{ c.name }}</div>
          <div class="flex-cb mt-1">
            <span class="text-xs text-g-500">¥{{ c.price }} 起</span>
          </div>
        </div>
      </div>
      <div v-if="currentCompany?.content" class="mt-3 px-4 py-3 bg-theme/10 rounded-md text-sm text-g-700">
        {{ currentCompany.content }}
      </div>
    </div>

    <!-- ② 收货信息 -->
    <div class="art-card p-5 mb-5">
      <h4 class="font-bold mb-3">② 收货信息</h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <ElFormItem label="客户姓名" class="!mb-0">
          <ElInput v-model="form.name" placeholder="客户真实姓名" />
        </ElFormItem>
        <ElFormItem label="联系电话" class="!mb-0">
          <ElInput v-model="form.phone" placeholder="11位手机号" maxlength="15" />
        </ElFormItem>
        <ElFormItem label="收货地址" class="!mb-0 md:col-span-2">
          <ElInput v-model="form.address" placeholder="省市区县详细地址" />
        </ElFormItem>
        <ElFormItem label="其他备注" class="!mb-0 md:col-span-2">
          <ElInput v-model="form.remark" placeholder="选填，如：加急、颜色要求等" />
        </ElFormItem>
      </div>
    </div>

    <!-- ③ 文件上传 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">③ 文件上传</h4>
        <ElSwitch v-model="fileEnabled" active-text="上传文件" inactive-text="仅交单" />
      </div>
      <div v-if="fileEnabled">
        <div class="flex items-center gap-3">
          <ElButton type="primary" @click="triggerUpload">选择文件</ElButton>
          <span class="text-g-500 text-sm">支持多选，每个文件可设置份数 / 黑白彩印 / 尺寸 / 单双面</span>
        </div>
        <input ref="fileInputRef" type="file" multiple style="display: none" @change="onFileChange" />
        <div v-if="selectedFiles.length" class="mt-4 space-y-3">
          <div v-for="(file, i) in selectedFiles" :key="file.name"
            class="flex flex-wrap items-center gap-3 px-4 py-3 rounded-md border border-g-300">
            <div class="flex-1 min-w-40">
              <div class="text-sm truncate">{{ file.name }}</div>
            </div>
            <ElInputNumber v-model="file.copy" :min="1" :max="999" size="small" />
            <ElSelect v-model="file.type" size="small" style="width: 90px">
              <ElOption label="黑白" value="黑白" /><ElOption label="彩印" value="彩印" />
            </ElSelect>
            <ElSelect v-model="file.size" size="small" style="width: 80px">
              <ElOption label="A4" value="A4" /><ElOption label="A3" value="A3" /><ElOption label="A5" value="A5" />
            </ElSelect>
            <ElSelect v-model="file.sided" size="small" style="width: 80px">
              <ElOption label="单面" value="单面" /><ElOption label="双面" value="双面" />
            </ElSelect>
            <ElButton link type="danger" @click="removeFile(i)">移除</ElButton>
          </div>
        </div>
        <div v-if="fileName" class="mt-3 px-4 py-2 rounded-md bg-theme/10 text-sm text-g-700 break-all">{{ fileName }}</div>
      </div>
    </div>

    <!-- ④ 选择规格 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">④ 选择规格</h4>
        <div class="flex gap-2">
          <ElButton :loading="specLoading" :disabled="!form.companyName || !form.name || !form.phone" @click="loadSpec">
            查询规格
          </ElButton>
        </div>
      </div>
      <div v-if="!specList.length" class="text-g-500 text-sm py-2">
        {{ form.companyName ? '点击「查询规格」获取可选项' : '请先选择公司' }}
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
        <div v-for="(s, i) in specList" :key="i"
          class="px-3 py-2 cursor-pointer rounded-md border border-g-300 duration-200 hover:border-theme hover:bg-theme/5"
          :class="{ '!border-theme bg-theme/10': checkedSpecs.includes(i) }"
          @click="toggleSpec(i)">
          <div class="flex-cb">
            <span class="text-sm">{{ s.name }}</span>
            <span class="text-theme font-bold text-sm">¥{{ (parseFloat(s.price || '0') * multiple).toFixed(2) }}</span>
          </div>
          <div v-if="s.id" class="text-xs text-g-500 mt-1">ID: {{ s.id }}</div>
        </div>
      </div>
    </div>

    <!-- 底部提交 -->
    <div class="art-card p-5 flex-cb sticky bottom-4">
      <div class="text-sm text-g-700">
        已选 <b class="text-theme">{{ checkedSpecs.length }}</b> 个规格，
        预计费用：<span class="text-theme font-bold text-lg">¥ {{ totalFee }}</span>
      </div>
      <ElButton type="primary" size="large" :loading="submitting" :disabled="!checkedSpecs.length" @click="handleSubmit">
        确认下单
      </ElButton>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import axios from 'axios'
  import { ElMessage, ElIcon } from 'element-plus'
  import { Search } from '@element-plus/icons-vue'
  import { getUserInfo, getTaowaCompanies, getTaowaSpec, addTaowaOrder } from '@/api/wk'

  defineOptions({ name: 'WkTaowaAdd' })

  const userInfo = ref<Partial<Api.Auth.UserInfo>>({})
  const configured = ref(false)
  const companies = ref<any[]>([])
  const search = ref('')
  const form = ref({ companyName: '', name: '', phone: '', address: '', remark: '' })
  const currentCompany = ref<any>(null)
  const specList = ref<any[]>([])
  const checkedSpecs = ref<number[]>([])
  const specLoading = ref(false)
  const submitting = ref(false)
  const fileEnabled = ref(false)
  const selectedFiles = ref<any[]>([])
  const fileInputRef = ref<HTMLInputElement>()
  const fileName = ref('')

  const multiple = computed(() => 2) // 盖章价格倍数（后端实际扣费为准）

  const filteredCompanies = computed(() => {
    if (!search.value.trim()) return companies.value
    const kw = search.value.trim().toLowerCase()
    return companies.value.filter(c => c.name.toLowerCase().includes(kw))
  })

  const totalFee = computed(() => {
    let total = 0
    for (const i of checkedSpecs.value) {
      total += parseFloat(specList.value[i]?.price || '0') * multiple.value
    }
    return total.toFixed(2)
  })

  const selectCompany = (c: any) => {
    form.value.companyName = c.name
    currentCompany.value = c
    specList.value = []
    checkedSpecs.value = []
  }

  const loadSpec = async () => {
    if (!form.value.companyName) return ElMessage.warning('请先选择公司')
    if (!form.value.name || !form.value.phone) return ElMessage.warning('请填写姓名和电话')
    specLoading.value = true
    try {
      const data: any = await getTaowaSpec({
        companyName: form.value.companyName,
        userinfo: `${form.value.name} ${form.value.phone} ${form.value.address}`
      })
      if (Array.isArray(data)) {
        specList.value = data
      } else if (data?.data && Array.isArray(data.data)) {
        specList.value = data.data
      } else {
        specList.value = []
        ElMessage.warning('未获取到规格')
      }
    } catch (e: any) {
      ElMessage.error(e?.message || '查询失败')
    } finally {
      specLoading.value = false
    }
  }

  const toggleSpec = (i: number) => {
    const idx = checkedSpecs.value.indexOf(i)
    if (idx >= 0) checkedSpecs.value.splice(idx, 1)
    else checkedSpecs.value.push(i)
  }

  const triggerUpload = () => fileInputRef.value?.click()

  const onFileChange = (e: Event) => {
    const files = (e.target as HTMLInputElement).files
    if (!files?.length) return
    for (const f of Array.from(files)) {
      if (!selectedFiles.value.some(x => x.name === f.name)) {
        selectedFiles.value.push({ file: f, name: f.name, copy: 1, type: '黑白', size: 'A4', sided: '单面' })
      }
    }
    ;(e.target as HTMLInputElement).value = ''
  }

  const removeFile = (i: number) => {
    selectedFiles.value.splice(i, 1)
    fileName.value = ''
  }

  // 上传文件到文档云，生成文件名（与源台格式一致）
  const uploadFiles = async () => {
    const formData = new FormData()
    const names: string[] = []
    for (const f of selectedFiles.value) {
      const prefix = Math.random().toString(36).substring(2, 5)
      const newName = `${prefix}-${f.copy}份-${f.type}-${f.size}-${f.sided}-${f.name}`
      names.push(newName)
      formData.append('uploaded_file[]', f.file, newName)
    }
    const res: any = await axios.post('https://www.documenticloud.cfd/upload.php', formData, { timeout: 120000 })
    if (res.data?.code !== 1) throw new Error('文件上传失败')
    fileName.value = names.join(',')
  }

  const handleSubmit = async () => {
    if (!checkedSpecs.value.length) return
    submitting.value = true
    try {
      if (fileEnabled.value) {
        if (!selectedFiles.value.length) {
          ElMessage.warning('请选择要上传的文件')
          return
        }
        await uploadFiles()
      }
      const data = checkedSpecs.value.map(i => {
        const s = specList.value[i]
        return { userinfo: form.value.name, userName: form.value.name, data: s }
      })
      const res: any = await addTaowaOrder({
        companyName: form.value.companyName,
        name: form.value.name,
        phone: form.value.phone,
        address: form.value.address,
        remark: form.value.remark,
        data,
        fileName: fileName.value
      })
      ElMessage.success(res?.msg || '下单成功')
      checkedSpecs.value = []
      specList.value = []
      selectedFiles.value = []
      fileName.value = ''
      userInfo.value = await getUserInfo()
    } catch (e: any) {
      ElMessage.error(e?.message || '下单失败')
    } finally {
      submitting.value = false
    }
  }

  onMounted(async () => {
    try {
      const [info, companiesData]: any[] = await Promise.all([getUserInfo(), getTaowaCompanies()])
      userInfo.value = info
      companies.value = companiesData?.companies || []
      configured.value = companiesData?.configured !== false
    } catch (e) { /* 未配置时静默 */ }
  })
</script>

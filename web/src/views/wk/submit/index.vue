<template>
  <div class="submit-page">
    <!-- 1. 课程分类 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">① 课程分类</h4>
        <span class="text-g-500 text-sm">余额：<span class="text-theme font-bold">¥ {{ userInfo.money?.toFixed(2) }}</span></span>
      </div>
      <div class="flex flex-wrap gap-2">
        <ElTag :type="filter.fenlei === '' ? 'primary' : 'info'" :effect="filter.fenlei === '' ? 'dark' : 'plain'"
          class="!cursor-pointer" @click="onFenLeiChange('')">全部</ElTag>
        <ElTag v-for="f in fenLeiList" :key="f.id"
          :type="filter.fenlei === String(f.id) ? 'primary' : 'info'"
          :effect="filter.fenlei === String(f.id) ? 'dark' : 'plain'"
          class="!cursor-pointer" @click="onFenLeiChange(String(f.id))">{{ f.name }}</ElTag>
      </div>
    </div>

    <!-- 2. 课程平台 -->
    <div class="art-card p-5 mb-5">
      <div class="flex-cb mb-3">
        <h4 class="font-bold">② 课程平台</h4>
        <ElInput v-model="search" placeholder="搜索平台" clearable size="small" style="width:240px">
          <template #prefix><ElIcon><Search /></ElIcon></template>
        </ElInput>
      </div>
      <div v-if="currentClass" class="mb-3 px-4 py-3 bg-theme/10 rounded-md flex-cb">
        <span>已选：<b>{{ currentClass.name }}</b>
          <span class="text-theme font-bold ml-2">¥{{ calcPrice(currentClass) }}</span>
          <span v-if="currentClass.content" class="text-g-500 text-sm ml-3">{{ currentClass.content }}</span>
        </span>
        <ElButton size="small" link @click="filter.cid = undefined">清除</ElButton>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 max-h-72 overflow-auto">
        <div v-if="!filteredClasses.length" class="col-span-full p-6 text-center text-g-500">无匹配平台</div>
        <div v-for="c in filteredClasses" :key="c.cid"
          class="px-3 py-2 cursor-pointer rounded-md border border-g-300 duration-200 hover:border-theme hover:bg-theme/5"
          :class="{ '!border-theme bg-theme/10': filter.cid === c.cid }"
          @click="filter.cid = c.cid">
          <div class="text-sm truncate">{{ c.name }}</div>
          <div class="flex-cb mt-1">
            <span class="text-xs text-g-500">¥{{ calcPrice(c) }}</span>
            <span v-if="c.fenlei" class="text-xs text-g-400">#{{ c.fenlei }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. 下单信息 -->
    <div class="art-card p-5 mb-5">
      <h4 class="font-bold mb-3">③ 下单信息</h4>
      <ElInput v-model="userinfo" type="textarea" :rows="6"
        placeholder="格式：学校 账号 密码（学校可为空）&#10;多账号请换行，比如：&#10;清华大学 12345678 password&#10;13800138000 mypassword" />
      <div class="mt-4 flex-cb">
        <div class="text-sm text-g-700">
          <span v-if="checkedItems.length">
            已勾选 <b class="text-theme">{{ checkedItems.length }}</b> 门课程，预计费用：
            <span class="text-theme font-bold text-lg">¥ {{ totalFee }}</span>
          </span>
        </div>
        <div class="flex gap-2">
          <ElButton type="success" :loading="querying" @click="handleQuery">查询课程</ElButton>
          <ElButton type="primary" :loading="submitting" :disabled="!checkedItems.length" @click="handleSubmit">
            确认下单 ({{ checkedItems.length }})
          </ElButton>
        </div>
      </div>
    </div>

    <!-- 4. 查询结果 -->
    <div v-if="!queryResults.length" class="art-card p-10 text-center text-g-500">
      <ElIcon size="40"><Search /></ElIcon>
      <p class="mt-2">填写下单信息后点「查询课程」，结果将显示在这里</p>
    </div>
    <div v-else>
      <h4 class="font-bold mb-3">④ 查询结果</h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="(rs, idx) in queryResults" :key="idx" class="art-card p-4">
          <div class="flex-cb mb-2 pb-2 border-b border-g-300">
            <div>
              <div class="font-bold text-sm">{{ rs.userinfo }}</div>
              <div v-if="rs.userName" class="text-xs text-g-500 mt-1">姓名：{{ rs.userName }}</div>
            </div>
            <ElTag :type="rs.code === 1 ? 'success' : 'danger'" size="small">{{ rs.msg }}</ElTag>
          </div>
          <ElCheckboxGroup v-if="rs.code === 1 && rs.data?.length"
            v-model="rs.checked" @change="updateChecked">
            <div v-for="(c, ci) in rs.data" :key="ci" class="mb-1">
              <ElCheckbox :value="ci">
                <span class="text-sm">{{ c.name }}</span>
                <span v-if="c.id" class="text-g-500 text-xs ml-1">[{{ c.id }}]</span>
              </ElCheckbox>
            </div>
          </ElCheckboxGroup>
          <div v-else-if="rs.code !== 1" class="text-g-500 text-sm py-2">{{ rs.msg }}</div>
          <div v-else class="text-g-500 text-sm py-2">暂无可学课程</div>
          <div v-if="rs.code === 1 && rs.data?.length" class="mt-2 flex gap-2 justify-end">
            <ElButton size="small" @click="checkAll(rs)">全选</ElButton>
            <ElButton size="small" @click="rs.checked = []; updateChecked()">清空</ElButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { ElMessage, ElIcon } from 'element-plus'
  import { Search } from '@element-plus/icons-vue'
  import { getClassAll, getFenLeiList, getUserInfo, queryCourse, addOrder } from '@/api/wk'

  defineOptions({ name: 'WkSubmit' })

  const userInfo = ref<Partial<Api.Auth.UserInfo>>({})
  const classes = ref<Api.Class.Item[]>([])
  const fenLeiList = ref<any[]>([])
  const filter = ref({ fenlei: '', cid: undefined as number | undefined })
  const search = ref('')
  const userinfo = ref('')
  const querying = ref(false)
  const submitting = ref(false)

  interface QueryRow {
    userinfo: string
    school: string
    user: string
    pass: string
    code: number
    msg: string
    userName?: string
    data: { id: string; name: string }[]
    checked: number[]
  }

  const queryResults = ref<QueryRow[]>([])
  const checkedItems = ref<{ school: string; user: string; pass: string; kcid: string; kcname: string }[]>([])

  const filteredClasses = computed(() => {
    let list = classes.value
    if (filter.value.fenlei) list = list.filter(c => c.fenlei === filter.value.fenlei)
    if (search.value.trim()) {
      const kw = search.value.trim().toLowerCase()
      list = list.filter(c => c.name.toLowerCase().includes(kw))
    }
    return list
  })

  const calcPrice = (c: Api.Class.Item) =>
    (parseFloat(c.price) * (userInfo.value.addprice || 1)).toFixed(2)

  const currentClass = computed(() => classes.value.find(c => c.cid === filter.value.cid))

  const totalFee = computed(() => {
    if (!currentClass.value) return '0.00'
    return (parseFloat(calcPrice(currentClass.value)) * checkedItems.value.length).toFixed(2)
  })

  const onFenLeiChange = (val: string) => {
    filter.value.fenlei = val
    filter.value.cid = undefined
  }

  const parseLine = (line: string) => {
    const parts = line.trim().split(/\s+/)
    if (parts.length === 2) return { school: '', user: parts[0], pass: parts[1] }
    if (parts.length >= 3) return { school: parts[0], user: parts[1], pass: parts[2] }
    return null
  }

  const handleQuery = async () => {
    if (!filter.value.cid) {
      ElMessage.warning('请选择课程平台')
      return
    }
    if (!userinfo.value.trim()) {
      ElMessage.warning('请输入下单信息')
      return
    }
    querying.value = true
    queryResults.value = []
    checkedItems.value = []
    try {
      const lines = userinfo.value.split(/\r?\n/).filter(l => l.trim())
      for (const line of lines) {
        const info = parseLine(line)
        if (!info) {
          queryResults.value.push({
            userinfo: line, school: '', user: '', pass: '', code: -1,
            msg: '格式错误', data: [], checked: []
          })
          continue
        }
        try {
          const res: any = await queryCourse({
            cid: filter.value.cid!,
            school: info.school || '1',
            user: info.user,
            pass: info.pass
          })
          queryResults.value.push({
            userinfo: line,
            school: info.school || '1',
            user: info.user,
            pass: info.pass,
            code: res.code === 0 || res.code === 1 ? 1 : res.code,
            msg: res.msg || '查询成功',
            userName: res.userName,
            data: res.data || [],
            checked: []
          })
        } catch (e: any) {
          queryResults.value.push({
            userinfo: line, ...info, code: -1,
            msg: e?.message || '查询失败', data: [], checked: []
          })
        }
      }
    } finally {
      querying.value = false
    }
  }

  const checkAll = (rs: QueryRow) => {
    rs.checked = rs.data.map((_, i) => i)
    updateChecked()
  }

  const updateChecked = () => {
    const items: typeof checkedItems.value = []
    for (const rs of queryResults.value) {
      for (const idx of rs.checked) {
        const c = rs.data[idx]
        if (c) items.push({
          school: rs.school, user: rs.user, pass: rs.pass,
          kcid: c.id || '', kcname: c.name
        })
      }
    }
    checkedItems.value = items
  }

  const handleSubmit = async () => {
    if (!checkedItems.value.length) return
    submitting.value = true
    try {
      const grouped: Record<string, typeof checkedItems.value> = {}
      for (const it of checkedItems.value) {
        const key = `${it.school}|${it.user}|${it.pass}`
        ;(grouped[key] ||= []).push(it)
      }
      let success = 0
      for (const list of Object.values(grouped)) {
        await addOrder({
          cid: filter.value.cid!,
          school: list[0].school,
          user: list[0].user,
          pass: list[0].pass,
          kcid: list.map(x => x.kcid).join(','),
          kcname: list.map(x => x.kcname).join(',')
        })
        success += list.length
      }
      ElMessage.success(`提交成功 ${success} 门课程`)
      queryResults.value = []
      checkedItems.value = []
      userinfo.value = ''
      userInfo.value = await getUserInfo()
    } finally {
      submitting.value = false
    }
  }

  onMounted(async () => {
    [classes.value, fenLeiList.value, userInfo.value] = await Promise.all([
      getClassAll(), getFenLeiList(), getUserInfo()
    ])
  })
</script>

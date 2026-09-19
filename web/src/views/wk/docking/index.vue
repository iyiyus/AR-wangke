<template>
  <div class="docking-page art-full-height">
    <ElCard shadow="never" class="mb-3" v-loading="loading">
      <template #header>我的对接信息</template>
      <ElDescriptions :column="2" border>
        <ElDescriptionsItem label="UID">{{ info.uid ?? '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="API Key">
          <template v-if="hasKey">
            {{ info.key }}
            <ElLink type="primary" class="ml-2" @click="copyText(info.key)">复制</ElLink>
          </template>
          <template v-else>
            <ElTag type="danger">未开通</ElTag>
            <span class="ml-2 text-g-600">请前往用户中心开通API</span>
          </template>
        </ElDescriptionsItem>
      </ElDescriptions>
      <ElAlert v-if="info.key === '0'" title="您尚未开通API功能，请先前往用户中心开通API后再对接" type="warning" show-icon class="mt-3" :closable="false" />
    </ElCard>

    <ElCard shadow="never" class="art-table-card">
      <template #header>
        接口说明（均为 POST 请求，返回 {code:200,msg,data} 格式）
      </template>
      <ElTable :data="apiList" border>
        <ElTableColumn prop="name" label="接口名称" width="120" />
        <ElTableColumn label="接口地址" min-width="260">
          <template #default="{ row }">
            <span class="code-text">{{ row.url }}</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="params" label="参数" min-width="240" />
        <ElTableColumn prop="desc" label="说明" min-width="220" />
      </ElTable>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getUserInfo } from '@/api/wk'

  defineOptions({ name: 'WkDocking' })

  const loading = ref(false)
  const info = ref<any>({})
  const hasKey = computed(() => !!info.value.key && info.value.key !== '0')

  const origin = window.location.origin
  const apiList = [
    {
      name: '查询余额',
      url: `${origin}/api.php?act=getmoney`,
      params: 'uid, key',
      desc: '查询当前账号余额'
    },
    {
      name: 'API查课',
      url: `${origin}/api.php?act=get`,
      params: 'uid, key, platform, school, user, pass',
      desc: '查询学习平台课程信息'
    },
    {
      name: 'API下单',
      url: `${origin}/api.php?act=add`,
      params: 'uid, key, platform, school, user, pass, kcid, kcname',
      desc: '提交课程下单任务'
    },
    {
      name: '查询下单',
      url: `${origin}/api.php?act=getadd`,
      params: 'uid, key, platform, school, user, pass, kcname',
      desc: '查询已下单课程信息'
    },
    {
      name: '查单',
      url: `${origin}/api.php?act=chadan`,
      params: 'username',
      desc: '按账号查询订单进度'
    },
    {
      name: '补刷',
      url: `${origin}/api.php?act=budan`,
      params: 'id',
      desc: '对指定任务进行补刷'
    },
    {
      name: '课程列表',
      url: `${origin}/api.php?act=getclass`,
      params: '无',
      desc: '获取可对接的课程列表'
    }
  ]

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

  onMounted(async () => {
    loading.value = true
    try {
      info.value = await getUserInfo()
    } finally {
      loading.value = false
    }
  })
</script>

<style scoped>
  .code-text {
    font-family: monospace;
    word-break: break-all;
  }
</style>

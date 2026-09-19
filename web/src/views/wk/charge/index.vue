<template>
  <div class="charge-page">
    <ElCard shadow="never" v-loading="loading">
      <template #header>联系上级</template>
      <ElDescriptions :column="2" border class="mb-5">
        <ElDescriptionsItem label="上级账号">{{ boss.user }}</ElDescriptionsItem>
        <ElDescriptionsItem label="上级昵称">{{ boss.name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="邀请码">{{ boss.yqm || '无' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="联系方式">
          <ElLink :href="`https://api.btstu.cn/qqtalk/api.php?qq=${boss.user}`" target="_blank" type="primary">
            发起会话
          </ElLink>
        </ElDescriptionsItem>
      </ElDescriptions>

      <ElTabs>
        <ElTabPane label="上级公告">
          <div class="p-2 text-g-600" v-html="boss.notice || '暂无公告'" />
        </ElTabPane>
        <ElTabPane label="我的公告">
          <ElInput v-model="myNotice" type="textarea" :rows="5" placeholder="请输入公告内容" class="mb-3" />
          <ElButton type="primary" :loading="saving" @click="handleSaveNotice">保存公告</ElButton>
        </ElTabPane>
      </ElTabs>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getUserInfo, updateNotice } from '@/api/wk'
  import request from '@/utils/http'

  defineOptions({ name: 'WkCharge' })

  const loading = ref(false)
  const saving = ref(false)
  const boss = ref<any>({})
  const myNotice = ref('')

  const load = async () => {
    loading.value = true
    try {
      const info = await getUserInfo()
      myNotice.value = info.notice || ''
      // 获取上级信息
      if (info.uuid && info.uuid !== info.uid) {
        const bossInfo = await request.get<any>({ url: `/api/user/boss/${info.uuid}` })
        boss.value = bossInfo
      } else {
        boss.value = { user: info.user, name: info.name, yqm: info.yqm, notice: '' }
      }
    } finally {
      loading.value = false
    }
  }

  const handleSaveNotice = async () => {
    saving.value = true
    try {
      await updateNotice(myNotice.value)
      ElMessage.success('保存成功')
    } finally {
      saving.value = false
    }
  }

  onMounted(load)
</script>

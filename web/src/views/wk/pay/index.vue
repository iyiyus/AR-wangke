<template>
  <div class="pay-page">
    <ElRow :gutter="20">
      <ElCol :md="12" :lg="10">
        <ElCard shadow="never">
          <template #header>账户充值</template>
          <ElDescriptions :column="1" border class="mb-5">
            <ElDescriptionsItem label="当前账号">{{ userInfo.user }}</ElDescriptionsItem>
            <ElDescriptionsItem label="当前余额">
              <span class="text-theme font-bold text-lg">¥ {{ userInfo.money?.toFixed(2) }}</span>
            </ElDescriptionsItem>
          </ElDescriptions>
          <ElForm :model="form" label-width="90px">
            <ElFormItem label="充值金额">
              <ElInput v-model="form.money" placeholder="请输入充值金额" type="number">
                <template #prepend>¥</template>
              </ElInput>
            </ElFormItem>
            <ElFormItem label="支付方式">
              <ElRadioGroup v-model="form.type">
                <ElRadio value="alipay" v-if="conf.is_alipay === '1'">支付宝</ElRadio>
                <ElRadio value="wxpay" v-if="conf.is_wxpay === '1'">微信支付</ElRadio>
                <ElRadio value="qqpay" v-if="conf.is_qqpay === '1'">QQ支付</ElRadio>
              </ElRadioGroup>
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" :loading="loading" @click="handlePay">去支付</ElButton>
            </ElFormItem>
          </ElForm>
          <ElAlert type="warning" :closable="false">
            支付完成后如未自动跳转，请手动刷新页面
          </ElAlert>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getUserInfo, createPay, getAdminConfig } from '@/api/wk'

  defineOptions({ name: 'WkPay' })

  const userInfo = ref<Partial<Api.Auth.UserInfo>>({})
  const conf = ref<Record<string, string>>({})
  const loading = ref(false)
  const form = ref({ money: '', type: 'alipay' })

  const handlePay = async () => {
    if (!form.value.money || parseFloat(form.value.money) <= 0) {
      ElMessage.warning('请输入充值金额')
      return
    }
    loading.value = true
    try {
      const res = await createPay(form.value)
      window.open(res.pay_url, '_blank')
    } finally {
      loading.value = false
    }
  }

  onMounted(async () => {
    [userInfo.value, conf.value] = await Promise.all([getUserInfo(), getAdminConfig()])
  })
</script>

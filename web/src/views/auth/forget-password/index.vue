<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">忘记密码</h3>
          <p class="sub-title">输入账号，重置邮件将发送到绑定邮箱</p>

          <div class="mt-5">
            <span class="input-label">账号</span>
            <ElInput class="custom-height" placeholder="请输入账号" v-model.trim="username" />
          </div>

          <div style="margin-top: 15px">
            <ElButton class="w-full custom-height" type="primary" @click="handleSend" :loading="loading" v-ripple>
              发送重置邮件
            </ElButton>
          </div>

          <div style="margin-top: 15px">
            <ElButton class="w-full custom-height" plain @click="toLogin">返回登录</ElButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'ForgetPassword' })

  const router = useRouter()
  const username = ref('')
  const loading = ref(false)

  const handleSend = async () => {
    if (!username.value) {
      ElMessage.warning('请输入账号')
      return
    }
    loading.value = true
    try {
      await request.post({ url: '/api/auth/forget-password', params: { account: username.value } })
      ElMessage.success('重置邮件已发送，请查收邮箱')
      router.push({ name: 'Login' })
    } catch (e: any) {
      ElMessage.error(e?.message || '发送失败')
    } finally {
      loading.value = false
    }
  }

  const toLogin = () => router.push({ name: 'Login' })
</script>

<style scoped>
  @import '../login/style.css';
</style>

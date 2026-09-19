<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">重置密码</h3>
          <p class="sub-title">请输入新密码</p>

          <div class="mt-5">
            <span class="input-label">新密码</span>
            <ElInput class="custom-height" type="password" show-password placeholder="至少6位" v-model.trim="newPass" />
          </div>
          <div class="mt-4">
            <span class="input-label">确认密码</span>
            <ElInput class="custom-height" type="password" show-password placeholder="再次输入新密码" v-model.trim="confirmPass" />
          </div>

          <div style="margin-top: 15px">
            <ElButton class="w-full custom-height" type="primary" @click="handleReset" :loading="loading" v-ripple>
              确认重置
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
  defineOptions({ name: 'ResetPassword' })

  const route = useRoute()
  const router = useRouter()
  const newPass = ref('')
  const confirmPass = ref('')
  const loading = ref(false)

  const token = route.query.token as string

  const handleReset = async () => {
    if (!newPass.value || !confirmPass.value) {
      ElMessage.warning('请输入完整')
      return
    }
    if (newPass.value !== confirmPass.value) {
      ElMessage.warning('两次密码不一致')
      return
    }
    if (newPass.value.length < 6) {
      ElMessage.warning('密码至少6位')
      return
    }
    loading.value = true
    try {
      await request.post({ url: '/api/auth/reset-password', params: { token, new_pass: newPass.value } })
      ElMessage.success('密码重置成功，请重新登录')
      router.push({ name: 'Login' })
    } catch (e: any) {
      ElMessage.error(e?.message || '重置失败')
    } finally {
      loading.value = false
    }
  }

  const toLogin = () => router.push({ name: 'Login' })
</script>

<style scoped>
  @import '../login/style.css';
</style>

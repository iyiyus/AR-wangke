<!-- 注册页面 -->
<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">{{ $t('register.title') }}</h3>
          <p class="sub-title">{{ $t('register.subTitle') }}</p>
          <ElForm
            class="mt-7.5"
            ref="formRef"
            :model="formData"
            :rules="rules"
            label-position="top"
            :key="formKey"
          >
            <ElFormItem prop="username">
              <div class="qq-row">
                <ElInput
                  class="custom-height"
                  v-model.trim="formData.username"
                  placeholder="请输入QQ号"
                  @input="formData.username = formData.username.replace(/\D/g, '')"
                />
                <img v-if="qqAvatar" :src="qqAvatar" class="qq-avatar" />
              </div>
            </ElFormItem>

            <ElFormItem prop="email">
              <ElInput
                class="custom-height"
                v-model.trim="formData.email"
                placeholder="输入QQ号"
                @input="formData.email = formData.email.replace(/\D/g, '')"
              >
                <template #suffix>@qq.com</template>
              </ElInput>
            </ElFormItem>

            <ElFormItem prop="password">
              <ElInput
                class="custom-height"
                v-model.trim="formData.password"
                :placeholder="$t('register.placeholder.password')"
                type="password"
                autocomplete="off"
                show-password
              />
            </ElFormItem>

            <ElFormItem prop="confirmPassword">
              <ElInput
                class="custom-height"
                v-model.trim="formData.confirmPassword"
                :placeholder="$t('register.placeholder.confirmPassword')"
                type="password"
                autocomplete="off"
                show-password
              />
            </ElFormItem>

            <ElFormItem prop="yqm">
              <ElInput
                class="custom-height"
                v-model.trim="formData.yqm"
                placeholder="邀请码"
              />
            </ElFormItem>

            <ElFormItem prop="agreement">
              <ElCheckbox v-model="formData.agreement">
                {{ $t('register.agreeText') }}
                <span
                  style="color: var(--theme-color); cursor: pointer"
                  @click="privacyVisible = true"
                  >{{ $t('register.privacyPolicy') }}</span
                >
              </ElCheckbox>
            </ElFormItem>

            <div style="margin-top: 15px">
              <ElButton
                class="w-full custom-height"
                type="primary"
                @click="register"
                :loading="loading"
                v-ripple
              >
                {{ $t('register.submitBtnText') }}
              </ElButton>
            </div>

            <div class="mt-5 text-sm text-g-600">
              <span>{{ $t('register.hasAccount') }}</span>
              <RouterLink class="text-theme" :to="{ name: 'Login' }">{{
                $t('register.toLogin')
              }}</RouterLink>
            </div>
          </ElForm>
        </div>
      </div>
    </div>

    <ElDialog v-model="privacyVisible" title="隐私政策" width="600px" :close-on-click-modal="false">
      <div class="privacy-content">
        <h3>一、信息收集</h3>
        <p>我们收集您在注册和使用本服务时提供的：QQ号、昵称、密码、邮箱、订单记录等信息。</p>
        <h3>二、信息使用</h3>
        <p>信息仅用于提供网课代刷、实习打卡、盖章病历等服务，账号安全验证，订单处理和售后。</p>
        <h3>三、信息保护</h3>
        <p>数据传输采用加密方式，敏感信息不对外公开，仅授权人员可访问必要数据。</p>
        <h3>四、信息共享</h3>
        <p>我们不会向任何第三方出售或泄露您的个人信息。为完成订单必须向源台传递的必要信息除外。</p>
        <h3>五、账号安全</h3>
        <p>请妥善保管账号密码，因个人原因导致的信息泄露由用户自行承担。如有问题请联系管理员。</p>
      </div>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import type { FormInstance, FormRules } from 'element-plus'
  import { register as registerApi } from '@/api/wk'

  defineOptions({ name: 'Register' })

  interface RegisterForm {
    username: string
    email: string
    password: string
    confirmPassword: string
    yqm: string
    agreement: boolean
  }

  const USERNAME_MIN_LENGTH = 3
  const USERNAME_MAX_LENGTH = 20
  const PASSWORD_MIN_LENGTH = 6
  const REDIRECT_DELAY = 1000

  const { t, locale } = useI18n()
  const router = useRouter()
  const formRef = ref<FormInstance>()

  const loading = ref(false)
  const formKey = ref(0)
  const privacyVisible = ref(false)

  watch(locale, () => { formKey.value++ })

  const formData = reactive<RegisterForm>({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    yqm: '',
    agreement: false
  })

  const qqAvatar = computed(() => {
    if (/^\d{5,11}$/.test(formData.username)) {
      return `http://q.qlogo.cn/headimg_dl?dst_uin=${formData.username}&spec=100&img_type=jpg`
    }
    return ''
  })

  const validatePassword = (_rule: any, value: string, callback: (error?: Error) => void) => {
    if (!value) { callback(new Error(t('register.placeholder.password'))); return }
    if (formData.confirmPassword) formRef.value?.validateField('confirmPassword')
    callback()
  }

  const validateConfirmPassword = (_rule: any, value: string, callback: (error?: Error) => void) => {
    if (!value) { callback(new Error(t('register.rule.confirmPasswordRequired'))); return }
    if (value !== formData.password) { callback(new Error(t('register.rule.passwordMismatch'))); return }
    callback()
  }

  const validateAgreement = (_rule: any, value: boolean, callback: (error?: Error) => void) => {
    if (!value) { callback(new Error(t('register.rule.agreementRequired'))); return }
    callback()
  }

  const rules = computed<FormRules<RegisterForm>>(() => ({
    username: [
      { required: true, message: '请输入QQ号', trigger: 'blur' },
      { pattern: /^\d{5,11}$/, message: '请输入5-11位数字QQ号', trigger: 'blur' }
    ],
    email: [
      { required: true, message: '请输入QQ号', trigger: 'blur' },
      { pattern: /^\d{5,11}$/, message: '请输入5-11位数字', trigger: 'blur' }
    ],
    password: [
      { required: true, validator: validatePassword, trigger: 'blur' },
      { min: PASSWORD_MIN_LENGTH, message: t('register.rule.passwordLength'), trigger: 'blur' }
    ],
    confirmPassword: [{ required: true, validator: validateConfirmPassword, trigger: 'blur' }],
    yqm: [{ required: true, message: '请输入邀请码', trigger: 'blur' }],
    agreement: [{ validator: validateAgreement, trigger: 'change' }]
  }))

  const register = async () => {
    if (!formRef.value) return
    try {
      await formRef.value.validate()
      loading.value = true
      await registerApi({
        name: formData.username,
        account: formData.username,
        pass: formData.password,
        email: formData.email + '@qq.com',
        yqm: formData.yqm || undefined
      })
      ElMessage.success('注册成功')
      setTimeout(() => router.push({ name: 'Login' }), REDIRECT_DELAY)
    } catch (error) {
      console.error('注册失败:', error)
    } finally {
      loading.value = false
    }
  }

  const toLogin = () => {
    router.push({ name: 'Login' })
  }
</script>

<style scoped>
  @import '../login/style.css';
  .qq-row { display: flex; align-items: center; gap: 12px; width: 100%; }
  .qq-row :deep(.el-input) { flex: 1; }
  .qq-avatar { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; flex-shrink: 0; }
  .privacy-content h3 { font-size: 15px; margin: 16px 0 8px; }
  .privacy-content p { font-size: 13px; color: #666; line-height: 1.8; margin: 0; }
</style>

<template>
  <div class="p-5 max-w-2xl">
    <!-- 修改头像 -->
    <div class="art-card-sm mb-5 p-5 flex items-center gap-5">
      <img :src="avatar" class="w-16 h-16 rounded-full object-cover border border-g-300" />
      <div>
        <ElUpload action="/api/upload" name="file" :show-file-list="false" :headers="uploadHeaders" accept="image/*" :on-success="onAvatarSuccess">
          <ElButton type="primary">上传新头像</ElButton>
        </ElUpload>
        <p class="text-xs text-g-500 mt-2">支持 JPG/PNG，建议正方形图片</p>
      </div>
    </div>

    <!-- 修改邮箱 -->
    <div class="art-card-sm mb-5">
      <h1 class="p-4 text-xl font-normal border-b border-g-300">修改邮箱</h1>
      <ElForm :model="emailForm" class="p-5" label-width="100px" label-position="top">
        <ElFormItem label="当前邮箱" prop="oldEmail">
          <ElInput v-model="emailForm.oldEmail" disabled />
        </ElFormItem>
        <ElFormItem label="新邮箱" prop="newEmail">
          <ElInput v-model="emailForm.newEmail" placeholder="请输入新邮箱" />
        </ElFormItem>
        <div class="flex justify-end">
          <ElButton type="primary" :loading="emailLoading" @click="saveEmail">保存</ElButton>
        </div>
      </ElForm>
    </div>

    <!-- 修改密码 -->
    <div class="art-card-sm">
      <h1 class="p-4 text-xl font-normal border-b border-g-300">修改密码</h1>
      <ElForm :model="pwdForm" class="p-5" label-width="100px" label-position="top">
        <ElFormItem label="当前密码" prop="password">
          <ElInput v-model="pwdForm.password" type="password" show-password />
        </ElFormItem>
        <ElFormItem label="新密码" prop="newPassword">
          <ElInput v-model="pwdForm.newPassword" type="password" show-password />
        </ElFormItem>
        <ElFormItem label="确认新密码" prop="confirmPassword">
          <ElInput v-model="pwdForm.confirmPassword" type="password" show-password />
        </ElFormItem>
        <div class="flex justify-end">
          <ElButton type="primary" :loading="pwdLoading" @click="savePwd">保存</ElButton>
        </div>
      </ElForm>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  import request from '@/utils/http'

  defineOptions({ name: 'UserCenter' })

  const userStore = useUserStore()

  const avatar = ref('https://api.dicebear.com/7.x/avataaars/svg?seed=default')
  const uploadHeaders = { Authorization: 'Bearer ' + (localStorage.getItem('sys-v3.0.1-user') ? JSON.parse(localStorage.getItem('sys-v3.0.1-user')).accessToken : '') }
  const emailForm = reactive({ oldEmail: '', newEmail: '' })
  const emailLoading = ref(false)

  const pwdForm = reactive({ password: '', newPassword: '', confirmPassword: '' })
  const pwdLoading = ref(false)

  onMounted(async () => {
    try {
      const info: any = await userStore.getUserInfo
      if (info?.email) emailForm.oldEmail = info.email
      if (info?.faceimg) avatar.value = info.faceimg
    } catch {}
  })

  const onAvatarSuccess = async (res: any) => {
    const url = res?.data?.url
    if (!url) return
    try {
      await request.post({ url: '/api/user/avatar', params: { avatar: url } })
      avatar.value = url
      ElMessage.success('头像更新成功')
    } catch (e: any) {
      ElMessage.error(e?.message || '更新失败')
    }
  }

  const saveEmail = async () => {
    if (!emailForm.newEmail) return ElMessage.warning('请输入新邮箱')
    emailLoading.value = true
    try {
      await request.post({ url: '/api/user/email', params: { email: emailForm.newEmail } })
      ElMessage.success('邮箱修改成功')
      emailForm.oldEmail = emailForm.newEmail
      emailForm.newEmail = ''
    } catch (e: any) {
      ElMessage.error(e?.message || '修改失败')
    } finally {
      emailLoading.value = false
    }
  }

  const savePwd = async () => {
    if (!pwdForm.password || !pwdForm.newPassword) return ElMessage.warning('请填写完整')
    if (pwdForm.newPassword !== pwdForm.confirmPassword) return ElMessage.warning('两次密码不一致')
    pwdLoading.value = true
    try {
      await request.post({
        url: '/api/user/passwd',
        params: { oldpass: pwdForm.password, newpass: pwdForm.newPassword }
      })
      ElMessage.success('密码修改成功')
      pwdForm.password = ''
      pwdForm.newPassword = ''
      pwdForm.confirmPassword = ''
    } catch (e: any) {
      ElMessage.error(e?.message || '修改失败')
    } finally {
      pwdLoading.value = false
    }
  }
</script>

<template>
  <div class="passwd-page">
    <ElCard shadow="never">
      <template #header>修改密码</template>
      <ElForm :model="form" :rules="rules" ref="formRef" label-width="100px" class="max-w-md">
        <ElFormItem label="原密码" prop="old_pass">
          <ElInput v-model="form.old_pass" type="password" show-password />
        </ElFormItem>
        <ElFormItem label="新密码" prop="new_pass">
          <ElInput v-model="form.new_pass" type="password" show-password />
        </ElFormItem>
        <ElFormItem label="确认密码" prop="confirm_pass">
          <ElInput v-model="form.confirm_pass" type="password" show-password />
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" :loading="submitting" @click="handleSubmit">确认修改</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive } from 'vue'
  import { ElMessage, type FormInstance } from 'element-plus'
  import { changePassword } from '@/api/wk'
  import { useUserStore } from '@/store/modules/user'

  defineOptions({ name: 'WkPasswd' })

  const formRef = ref<FormInstance>()
  const submitting = ref(false)
  const form = reactive({ old_pass: '', new_pass: '', confirm_pass: '' })

  const rules = {
    old_pass: [{ required: true, message: '请输入原密码' }],
    new_pass: [{ required: true, message: '请输入新密码' }, { min: 6, message: '至少6位' }],
    confirm_pass: [
      { required: true, message: '请确认密码' },
      {
        validator: (_r: any, v: string, cb: any) => {
          if (v !== form.new_pass) cb(new Error('两次密码不一致'))
          else cb()
        }
      }
    ]
  }

  const handleSubmit = async () => {
    await formRef.value?.validate()
    submitting.value = true
    try {
      await changePassword({ old_pass: form.old_pass, new_pass: form.new_pass })
      ElMessage.success('修改成功，请重新登录')
      setTimeout(() => useUserStore().logOut(), 1000)
    } finally { submitting.value = false }
  }
</script>

<template>
  <div class="admin-config-page">
    <ElCard shadow="never" v-loading="loading">
      <template #header>
        <span class="font-medium">系统设置</span>
      </template>
      <ElTabs>
        <ElTabPane label="网站配置">
          <ElForm :model="conf" label-width="130px" class="max-w-2xl">
            <ElFormItem label="站点名称"><ElInput v-model="conf.sitename" /></ElFormItem>
            <ElFormItem label="站点关键词"><ElInput v-model="conf.keywords" /></ElFormItem>
            <ElFormItem label="站点描述"><ElInput v-model="conf.description" /></ElFormItem>
            <ElFormItem label="LOGO">
              <div class="flex items-center gap-3">
                <ElInput v-model="conf.logo" placeholder="上传后自动填充" class="flex-1" />
                <ElUpload
                  action="/api/upload"
                  name="file"
                  :show-file-list="false"
                  :headers="uploadHeaders"
                  accept="image/*"
                  :on-success="(res:any) => { if (res?.data?.url) conf.logo = res.data.url }"
                >
                  <ElButton type="primary">上传LOGO</ElButton>
                </ElUpload>
                <img v-if="conf.logo" :src="conf.logo" class="h-10 w-10 object-contain border border-g-300 rounded" />
              </div>
            </ElFormItem>
            <ElFormItem label="Favicon">
              <div class="flex items-center gap-3">
                <ElInput v-model="conf.favicon" placeholder="上传后自动填充" class="flex-1" />
                <ElUpload
                  action="/api/upload"
                  name="file"
                  :show-file-list="false"
                  :headers="uploadHeaders"
                  accept="image/*"
                  :on-success="(res:any) => { if (res?.data?.url) conf.favicon = res.data.url }"
                >
                  <ElButton type="primary">上传Favicon</ElButton>
                </ElUpload>
                <img v-if="conf.favicon" :src="conf.favicon" class="h-10 w-10 object-contain border border-g-300 rounded" />
              </div>
            </ElFormItem>
            <ElFormItem label="站点公告">
              <ElInput v-model="conf.notice" type="textarea" :rows="4" placeholder="支持HTML" />
            </ElFormItem>
            <ElFormItem label="弹窗公告">
              <ElInput v-model="conf.tcgonggao" type="textarea" :rows="3" />
            </ElFormItem>
            <ElFormItem label="开启水印">
              <ElSwitch v-model="conf.sykg" active-value="1" inactive-value="0" />
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="邮箱配置">
          <ElForm :model="conf" label-width="130px" class="max-w-2xl">
            <ElAlert type="info" :closable="false" class="!mb-4"
              title="用于忘记密码发送重置邮件。以QQ邮箱为例：SMTP服务器 smtp.qq.com，端口 465，密码填授权码（不是登录密码）。" />
            <ElFormItem label="SMTP服务器"><ElInput v-model="conf.smtp_host" placeholder="如 smtp.qq.com" /></ElFormItem>
            <ElFormItem label="SMTP端口"><ElInput v-model="conf.smtp_port" placeholder="如 465" /></ElFormItem>
            <ElFormItem label="发信邮箱"><ElInput v-model="conf.smtp_user" placeholder="如 xxx@qq.com" /></ElFormItem>
            <ElFormItem label="邮箱密码/授权码"><ElInput v-model="conf.smtp_pass" type="password" show-password /></ElFormItem>
            <ElFormItem label="发件人名称"><ElInput v-model="conf.smtp_from" placeholder="留空则用发信邮箱" /></ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="代理配置">
          <ElForm :model="conf" label-width="160px" class="max-w-2xl">
            <ElFormItem label="代理开通价格"><ElInput v-model="conf.user_ktmoney" /></ElFormItem>
            <ElFormItem label="开启上级迁移">
              <ElSwitch v-model="conf.sjqykg" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="允许邀请码注册">
              <ElSwitch v-model="conf.user_yqzc" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="允许后台开户">
              <ElSwitch v-model="conf.user_htkh" active-value="1" inactive-value="0" />
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="支付配置">
          <ElForm :model="conf" label-width="130px" class="max-w-2xl">
            <ElFormItem label="易支付API"><ElInput v-model="conf.epay_api" /></ElFormItem>
            <ElFormItem label="商户ID"><ElInput v-model="conf.epay_pid" /></ElFormItem>
            <ElFormItem label="商户KEY"><ElInput v-model="conf.epay_key" /></ElFormItem>
            <ElFormItem label="最低充值金额"><ElInput v-model="conf.zdpay" /></ElFormItem>
            <ElFormItem label="开启支付宝">
              <ElSwitch v-model="conf.is_alipay" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="开启微信支付">
              <ElSwitch v-model="conf.is_wxpay" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="开启QQ支付">
              <ElSwitch v-model="conf.is_qqpay" active-value="1" inactive-value="0" />
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="聚合登录">
          <ElForm :model="conf" label-width="160px" class="max-w-2xl">
            <ElFormItem label="开启快捷登录">
              <ElSwitch v-model="conf.login_kg" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="彩虹聚合登录地址"><ElInput v-model="conf.login_apiurl" /></ElFormItem>
            <ElFormItem label="应用ID"><ElInput v-model="conf.login_appid" /></ElFormItem>
            <ElFormItem label="应用KEY"><ElInput v-model="conf.login_appkey" /></ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="实习助手对接">
          <ElForm :model="conf" label-width="170px" class="max-w-2xl">
            <ElAlert type="info" :closable="false" class="!mb-4"
              title="盖章/病历对接（taowa）。填写源台地址、UID、KEY 与平台 ID 后即可下单测试。" />
            <ElFormItem label="对接平台地址"><ElInput v-model="conf.taowa_url" placeholder="https://www.sxzsjk.top" /></ElFormItem>
            <ElFormItem label="源台 UID"><ElInput v-model="conf.taowa_uid" placeholder="你的UID" /></ElFormItem>
            <ElFormItem label="源台 KEY"><ElInput v-model="conf.taowa_key" placeholder="你的KEY" /></ElFormItem>
            <ElFormItem label="盖章价格倍数"><ElInput v-model="conf.taowa_jg" placeholder="2" /></ElFormItem>
            <ElFormItem label="盖章分类 ID"><ElInput v-model="conf.taowa_fenlei" placeholder="填本系统分类管理里的ID" /></ElFormItem>
            <ElFormItem label="盖章平台 ID"><ElInput v-model="conf.taowa_hid" placeholder="填货源管理里的ID（如1）" /></ElFormItem>
            <ElFormItem label="病历价格倍数"><ElInput v-model="conf.taowa_bljg" placeholder="2" /></ElFormItem>
            <ElFormItem label="病历分类 ID"><ElInput v-model="conf.taowa_blfenlei" placeholder="填本系统分类管理里的ID" /></ElFormItem>
            <ElFormItem label="病历平台 ID"><ElInput v-model="conf.taowa_blhid" placeholder="填货源管理里的ID（如1）" /></ElFormItem>
            <ElFormItem label="源台价格乘数"><ElInput v-model="conf.taowa_danjia" placeholder="1.05" /></ElFormItem>
            <ElFormItem label="源台改价乘数"><ElInput v-model="conf.taowa_newdanjia" placeholder="1.08" /></ElFormItem>
            <ElFormItem label=" ">
              <ElButton type="primary" :loading="syncing" @click="handleSync">一键同步分类和商品</ElButton>
              <span class="text-g-500 text-sm ml-3">从源台拉取公司/病历模板，自动创建本地分类和商品，并回填上方ID</span>
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <ElTabPane label="实习打卡对接">
          <ElForm :model="conf" label-width="170px" class="max-w-2xl">
            <ElAlert type="info" :closable="false" class="!mb-4"
              title="实习打卡对接（daka / TaiShan）。填写源台账号与 token 后即可下单测试。" />
            <ElFormItem label="源台账号"><ElInput v-model="conf.daka_admin" placeholder="您的TaiShan账号" /></ElFormItem>
            <ElFormItem label="源台 Token"><ElInput v-model="conf.daka_token" placeholder="您的token" /></ElFormItem>
            <ElFormItem label="版本号"><ElInput v-model="conf.daka_ts_version" placeholder="260901" /></ElFormItem>
            <ElFormItem label="源台地址(JSON数组)">
              <ElInput v-model="conf.daka_url_list" type="textarea" :rows="3"
                placeholder='["http://location.copilotai.top:4007/copilot/","http://82.156.247.209:4007/copilot/","http://location.tspost.top:4007/copilot/"]' />
            </ElFormItem>
            <ElFormItem label="删单是否退款">
              <ElSelect v-model="conf.daka_del_return" style="width: 120px">
                <ElOption label="是" value="true" />
                <ElOption label="否" value="false" />
              </ElSelect>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
      </ElTabs>

      <div class="mt-4">
        <ElButton type="primary" :loading="saving" @click="handleSave">保存设置</ElButton>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getAdminConfig, saveAdminConfig, syncTaowa } from '@/api/wk'

  defineOptions({ name: 'WkAdminConfig' })

  const loading = ref(false)
  const uploadHeaders = { Authorization: 'Bearer ' + (localStorage.getItem('sys-v3.0.1-user') ? JSON.parse(localStorage.getItem('sys-v3.0.1-user')).accessToken : '') }
  const onLogoSuccess = (res: any) => {
    if (res?.data?.url) conf.logo = res.data.url
    else if (res?.code === 0) conf.logo = res.data.url
  }
  const saving = ref(false)
  const syncing = ref(false)
  const conf = ref<Record<string, string>>({})

  const handleSave = async () => {
    saving.value = true
    try {
      await saveAdminConfig(conf.value)
      ElMessage.success('保存成功')
    } finally { saving.value = false }
  }

  const handleSync = async () => {
    syncing.value = true
    try {
      const res: any = await syncTaowa()
      const d = res?.data || {}
      ElMessage.success(`同步完成：盖章分类ID=${d.gz_fenlei_id||'-'}，新增商品${d.gz_count||0}个；病历分类ID=${d.bl_fenlei_id||'-'}，新增商品${d.bl_count||0}个`)
      conf.value = await getAdminConfig()
    } catch (e: any) {
      ElMessage.error(e?.message || '同步失败')
    } finally { syncing.value = false }
  }

  onMounted(async () => {
    loading.value = true
    try { conf.value = await getAdminConfig() }
    finally { loading.value = false }
  })
</script>

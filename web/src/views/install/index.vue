<!-- 安装向导（复刻 UDID 安装系统逻辑） -->
<template>
  <div class="install-page">
    <div class="install-container">
      <div class="install-header">
        <span class="site-name">网课代刷管理系统</span>
      </div>

      <el-steps :active="currentStep" finish-status="success" class="install-steps">
        <el-step title="许可协议" />
        <el-step title="环境检测" />
        <el-step title="系统配置" />
        <el-step title="注册服务" />
        <el-step title="安装完成" />
      </el-steps>

      <div class="art-card install-card">

        <!-- Step 0: 协议 -->
        <template v-if="currentStep === 0">
          <h3 class="step-title">用户使用协议</h3>
          <p class="step-desc">在开始安装之前，请阅读并接受以下条款</p>
          <div class="agreement-box">
            <p><b>1. 法律合规</b><br>本程序仅供学习交流使用，不得用于违法用途。</p>
            <p><b>2. 免责声明</b><br>本软件按"原样"提供，开发者不对因使用本软件导致的任何损失负责。</p>
            <p><b>3. 授权说明</b><br>严禁用于诈骗、博彩、色情等违法违规内容的传播。</p>
            <p><b>4. 对接配置</b><br>实习打卡源台账号/Token、盖章/病历平台参数需由管理员在系统设置中自行填写，本程序仅负责展示与下单流程。</p>
          </div>
          <el-checkbox v-model="agreed" class="mt-4">我已阅读并同意上述协议</el-checkbox>
          <div class="step-footer">
            <el-button type="primary" :disabled="!agreed" @click="currentStep++">下一步</el-button>
          </div>
        </template>

        <!-- Step 1: 环境检测 -->
        <template v-if="currentStep === 1">
          <h3 class="step-title">环境兼容性检测</h3>
          <p class="step-desc">检测系统运行环境是否满足要求</p>
          <div class="check-list">
            <div class="check-row" v-for="item in envChecks" :key="item.name">
              <span class="text-sm">{{ item.name }}</span>
              <el-tag type="success" size="small">✓ {{ item.status }}</el-tag>
            </div>
            <div class="check-row">
              <span class="text-sm">数据库连接</span>
              <el-tag :type="envStatus.database ? 'success' : 'danger'" size="small">
                {{ envStatus.database ? '✓ 可连接' : '✗ 未连接' }}
              </el-tag>
            </div>
            <div class="check-row">
              <span class="text-sm">目录写入权限</span>
              <el-tag :type="envStatus.writable ? 'success' : 'danger'" size="small">
                {{ envStatus.writable ? '✓ 可写入' : '✗ 不可写' }}
              </el-tag>
            </div>
            <div class="check-row">
              <span class="text-sm">操作系统</span>
              <el-tag type="info" size="small">{{ envStatus.os || '未知' }}</el-tag>
            </div>
          </div>
          <div class="step-footer">
            <el-button @click="currentStep--">上一步</el-button>
            <el-button type="primary" @click="currentStep++">下一步</el-button>
          </div>
        </template>

        <!-- Step 2: 系统配置 -->
        <template v-if="currentStep === 2">
          <h3 class="step-title">系统参数配置</h3>
          <p class="step-desc">请填写数据库连接信息和系统参数</p>
          <el-form :model="form" label-position="top" ref="formRef" :rules="rules">
            <el-tabs v-model="configTab">
              <el-tab-pane label="数据库配置" name="db">
                <el-row :gutter="16">
                  <el-col :span="16">
                    <el-form-item label="主机地址" prop="dbHost">
                      <el-input v-model="form.dbHost" placeholder="127.0.0.1" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="端口">
                      <el-input v-model="form.dbPort" placeholder="3306" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="用户名" prop="dbUser">
                      <el-input v-model="form.dbUser" placeholder="root" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="密码">
                      <el-input v-model="form.dbPassword" type="password" show-password placeholder="数据库密码" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="24">
                    <el-form-item label="数据库名" prop="dbName">
                      <el-input v-model="form.dbName" placeholder="如 wk">
                        <template #append>
                          <el-button :loading="testingDB" @click="testDB">测试连接</el-button>
                        </template>
                      </el-input>
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-tab-pane>

              <el-tab-pane label="应用配置" name="app">
                <el-row :gutter="16">
                  <el-col :span="12">
                    <el-form-item label="服务端口号">
                      <el-input v-model="form.serverPort" placeholder="如 8080" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="systemd 服务名">
                      <el-input v-model="form.projectName" placeholder="如 wk-go" />
                      <div class="form-tip">用于 systemctl restart &lt;服务名&gt;</div>
                    </el-form-item>
                  </el-col>
                  <el-col :span="24">
                    <el-form-item label="网站名称">
                      <el-input v-model="form.siteName" placeholder="如 网课代刷管理系统" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-tab-pane>

              <el-tab-pane label="管理员账号" name="admin">
                <el-row :gutter="16">
                  <el-col :span="12">
                    <el-form-item label="账号" prop="adminUser">
                      <el-input v-model="form.adminUser" placeholder="如 admin" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="密码">
                      <el-input v-model="form.adminPwd" type="password" show-password placeholder="默认 admin123" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-tab-pane>
            </el-tabs>
          </el-form>
          <div class="step-footer">
            <el-button @click="currentStep--">上一步</el-button>
            <el-button type="primary" @click="goToService">下一步</el-button>
          </div>
        </template>

        <!-- Step 3: 注册 systemd 服务 -->
        <template v-if="currentStep === 3">
          <h3 class="step-title">注册 systemd 服务</h3>
          <p class="step-desc">将程序注册为系统服务，实现开机自启和自动重启</p>

          <div class="service-info">
            <div class="info-row"><span>服务名称</span><b>{{ form.projectName }}</b></div>
            <div class="info-row"><span>程序路径</span><b>自动识别</b></div>
            <div class="info-row"><span>重启策略</span><b>Restart=always（崩溃自动拉起）</b></div>
          </div>

          <div v-if="systemdStatus === 'idle'" class="systemd-tip">
            点击下方按钮，系统将自动完成 systemd 服务注册、启用和启动（Linux 部署环境）
          </div>
          <div v-else-if="systemdStatus === 'success'" class="systemd-result">
            ✓ systemd 服务 "{{ systemdResult?.serviceName }}" 注册成功，程序路径：{{ systemdResult?.exePath }}
          </div>
          <div v-else-if="systemdStatus === 'error'" class="systemd-error">
            ✗ 注册失败：{{ systemdError }}，可手动执行 install-service.sh 或跳过
          </div>

          <div class="step-footer">
            <el-button @click="currentStep--">上一步</el-button>
            <el-button v-if="systemdStatus !== 'success'" :loading="registeringSystemd" @click="registerSystemd">
              注册 systemd 服务
            </el-button>
            <el-button v-if="systemdStatus === 'success'" type="primary" :loading="installing" @click="doInstall">
              立即部署
            </el-button>
            <el-button v-if="systemdStatus === 'error'" type="primary" :loading="installing" plain @click="doInstall">
              跳过并部署
            </el-button>
          </div>
        </template>

        <!-- Step 4: 完成 -->
        <template v-if="currentStep === 4">
          <div class="done-wrap">
            <div class="done-icon">✓</div>
            <h3 class="step-title">安装完成</h3>
            <p class="step-desc">系统已成功部署，请保存以下信息</p>
            <div class="result-box">
              <div class="res-row"><span>管理账号</span><b>{{ result.adminUser }}</b></div>
              <div class="res-row"><span>管理密码</span><b>{{ result.adminPwd }}</b></div>
              <div class="res-row"><span>服务端口</span><b>{{ result.serverPort }}</b></div>
              <div class="res-row"><span>网站名称</span><b>{{ result.siteName }}</b></div>
            </div>

            <div class="restart-box" :class="restartStatus">
              <template v-if="restartStatus === 'restarting'">
                <el-icon class="is-loading"><Loading /></el-icon>
                <span>正在重启服务，请稍候... ({{ countdown }}s)</span>
              </template>
              <template v-else-if="restartStatus === 'done'">
                <span>✓ 服务已重启完成</span>
              </template>
              <template v-else-if="restartStatus === 'timeout'">
                <span>⚠ 重启超时，请手动执行 systemctl restart {{ form.projectName }}</span>
              </template>
            </div>

            <el-button v-if="restartStatus === 'done' || restartStatus === 'timeout'"
              type="primary" class="mt-4 w-full" @click="goAdmin">
              进入管理面板
            </el-button>
          </div>
        </template>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { Loading } from '@element-plus/icons-vue'
  import { getInstallStatus, checkInstallDB, registerInstallSystemd, runInstall } from '@/api/install'
  import type { FormInstance } from 'element-plus'

  defineOptions({ name: 'Install' })

  const currentStep = ref(0)
  const agreed = ref(false)
  const testingDB = ref(false)
  const installing = ref(false)
  const configTab = ref('db')
  const formRef = ref<FormInstance>()
  const result = ref<any>({})
  const envStatus = ref<any>({})

  // 环境检测（进入第 2 步时请求）
  watch(currentStep, async (val) => {
    if (val === 1) {
      try {
        const data: any = await getInstallStatus()
        envStatus.value = data || {}
      } catch {}
    }
  })

  const envChecks = [
    { name: 'Go 运行环境', status: '正常' },
    { name: 'MySQL 数据库驱动', status: '已就绪' },
    { name: '文件存储权限', status: '可写入' },
    { name: '网络服务', status: '运行中' }
  ]

  const form = ref({
    dbHost: '127.0.0.1', dbPort: '3306', dbUser: '', dbPassword: '', dbName: '',
    serverPort: '8080', projectName: 'wk-go', siteName: '网课代刷管理系统',
    adminUser: '', adminPwd: ''
  })

  const rules = {
    dbHost: [{ required: true, message: '请输入数据库地址', trigger: 'blur' }],
    dbUser: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    dbName: [{ required: true, message: '请输入数据库名', trigger: 'blur' }],
    adminUser: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }]
  }

  async function testDB() {
    testingDB.value = true
    try {
      await checkInstallDB({
        host: form.value.dbHost,
        port: form.value.dbPort,
        user: form.value.dbUser,
        password: form.value.dbPassword,
        dbName: form.value.dbName
      })
      ElMessage.success('数据库连接成功')
    } catch (e: any) {
      ElMessage.error(e?.message || '数据库连接失败')
    } finally {
      testingDB.value = false
    }
  }

  async function goToService() {
    if (!formRef.value) return
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) {
      configTab.value = 'db'
      return
    }
    currentStep.value = 3
  }

  async function doInstall() {
    installing.value = true
    try {
      const data: any = await runInstall({ ...form.value })
      result.value = data
      currentStep.value = 4
    } catch (e: any) {
      ElMessage.error(e?.message || '安装失败')
    } finally {
      installing.value = false
    }
  }

  // systemd 注册
  const registeringSystemd = ref(false)
  const systemdStatus = ref<'idle' | 'success' | 'error'>('idle')
  const systemdResult = ref<any>(null)
  const systemdError = ref('')

  async function registerSystemd() {
    registeringSystemd.value = true
    try {
      const data: any = await registerInstallSystemd({ serviceName: form.value.projectName })
      systemdResult.value = data
      systemdStatus.value = 'success'
      ElMessage.success('systemd 服务注册成功')
    } catch (e: any) {
      systemdStatus.value = 'error'
      systemdError.value = e?.message || '未知错误'
    } finally {
      registeringSystemd.value = false
    }
  }

  // 重启轮询（进入第 5 步时触发）
  const restartStatus = ref<'restarting' | 'done' | 'timeout'>('restarting')
  const countdown = ref(60)
  let countdownTimer: ReturnType<typeof setInterval> | null = null
  let pollTimer: ReturnType<typeof setInterval> | null = null

  watch(currentStep, (val) => {
    if (val === 4) startRestartPolling()
  })

  function startRestartPolling() {
    restartStatus.value = 'restarting'
    countdown.value = 60
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(countdownTimer!)
        if (restartStatus.value === 'restarting') {
          restartStatus.value = 'timeout'
          clearInterval(pollTimer!)
        }
      }
    }, 1000)
    let wasDown = false
    pollTimer = setInterval(async () => {
      try {
        const data: any = await getInstallStatus()
        // 服务已安装且可访问，直接判断完成
        if (data?.installed) {
          restartStatus.value = 'done'
          clearInterval(pollTimer!)
          clearInterval(countdownTimer!)
          return
        }
        if (wasDown) {
          restartStatus.value = 'done'
          clearInterval(pollTimer!)
          clearInterval(countdownTimer!)
        }
      } catch {
        wasDown = true
      }
    }, 2000)
  }

  onUnmounted(() => {
    if (countdownTimer) clearInterval(countdownTimer)
    if (pollTimer) clearInterval(pollTimer)
  })

  function goAdmin() {
    window.location.hash = '#/'
    window.location.reload()
  }

  onMounted(async () => {
    try {
      const data: any = await getInstallStatus()
      if (data?.installed) {
        window.location.hash = '#/'
        window.location.reload()
      }
    } catch {}
  })
</script>

<style scoped>
  .install-page {
    min-height: 100vh;
    background: var(--default-bg-color, #f4f6f8);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
  }
  .install-container { width: 100%; max-width: 680px; }
  .install-header { display: flex; align-items: center; gap: 10px; margin-bottom: 28px; justify-content: center; }
  .site-name { font-size: 20px; font-weight: 700; color: var(--art-gray-900, #1f2937); }
  .install-steps { margin-bottom: 24px; }
  .install-card { padding: 32px 36px; }
  .step-title { font-size: 20px; font-weight: 700; margin-bottom: 6px; color: var(--art-gray-900, #1f2937); }
  .step-desc { font-size: 14px; color: var(--art-gray-500, #6b7280); margin-bottom: 20px; }
  .form-tip { font-size: 12px; color: var(--art-gray-400, #9ca3af); margin-top: 4px; }
  .agreement-box {
    background: var(--art-gray-100, #f3f4f6);
    border: 1px solid var(--el-border-color);
    border-radius: 8px;
    padding: 16px 20px;
    max-height: 200px;
    overflow-y: auto;
    font-size: 14px;
    line-height: 1.8;
    color: var(--art-gray-700, #374151);
  }
  .agreement-box p { margin: 0 0 10px; }
  .agreement-box b { color: var(--art-gray-900, #1f2937); }
  .check-list { margin-bottom: 8px; }
  .check-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 13px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  .check-row:last-child { border: none; }
  .step-footer { display: flex; justify-content: flex-end; gap: 10px; margin-top: 24px; }
  .service-info {
    background: var(--art-gray-100, #f3f4f6);
    border-radius: 8px;
    padding: 14px 20px;
  }
  .info-row {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    border-bottom: 1px dashed var(--el-border-color-lighter);
    font-size: 14px;
    color: var(--art-gray-500, #6b7280);
  }
  .info-row:last-child { border: none; }
  .info-row b { color: var(--art-gray-900, #1f2937); }
  .systemd-tip {
    margin-top: 16px;
    padding: 12px 16px;
    background: var(--el-color-primary-light-9);
    border-radius: 6px;
    font-size: 13px;
    color: var(--el-color-primary);
  }
  .systemd-result {
    margin-top: 16px;
    padding: 12px 16px;
    background: var(--el-color-success-light-9);
    border-radius: 6px;
    font-size: 13px;
    color: var(--el-color-success);
  }
  .systemd-error {
    margin-top: 16px;
    padding: 12px 16px;
    background: var(--el-color-danger-light-9);
    border-radius: 6px;
    font-size: 13px;
    color: var(--el-color-danger);
  }
  .done-wrap { text-align: center; padding: 10px 0; }
  .done-icon {
    width: 64px;
    height: 64px;
    line-height: 64px;
    margin: 0 auto 16px;
    border-radius: 50%;
    background: var(--el-color-success);
    color: #fff;
    font-size: 28px;
  }
  .result-box {
    text-align: left;
    background: var(--art-gray-100, #f3f4f6);
    border-radius: 8px;
    padding: 14px 24px;
    margin-bottom: 20px;
  }
  .res-row {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    font-size: 14px;
    color: var(--art-gray-500, #6b7280);
  }
  .res-row b { color: var(--art-gray-900, #1f2937); }
  .restart-box {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 12px;
    border-radius: 6px;
    font-size: 14px;
    margin-bottom: 8px;
  }
  .restart-box.restarting { background: var(--el-color-warning-light-9); color: var(--el-color-warning); }
  .restart-box.done { background: var(--el-color-success-light-9); color: var(--el-color-success); }
  .restart-box.timeout { background: var(--el-color-danger-light-9); color: var(--el-color-danger); }
  .w-full { width: 100%; }
  .mt-4 { margin-top: 16px; }
</style>

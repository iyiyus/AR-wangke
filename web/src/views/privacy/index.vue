<template>
  <div class="privacy-page">
    <div class="privacy-box">
      <h1>隐私政策</h1>
      <p class="update-time">更新日期：2026年9月</p>

      <h2>一、信息收集</h2>
      <p>我们收集您在注册和使用本服务时提供的以下信息：</p>
      <ul>
        <li>账号信息：QQ号、昵称、密码</li>
        <li>联系信息：邮箱地址（用于密码找回）</li>
        <li>使用记录：订单记录、操作日志</li>
      </ul>

      <h2>二、信息使用</h2>
      <p>您的信息仅用于以下用途：</p>
      <ul>
        <li>提供网课代刷、实习打卡、盖章病历等服务</li>
        <li>账号安全验证和密码找回</li>
        <li>订单处理和售后服务</li>
        <li>系统安全防护</li>
      </ul>

      <h2>三、信息保护</h2>
      <p>我们采取以下措施保护您的信息安全：</p>
      <ul>
        <li>数据传输采用加密方式</li>
        <li>敏感信息不对外公开</li>
        <li>仅授权人员可访问必要数据</li>
      </ul>

      <h2>四、信息共享</h2>
      <p>我们不会向任何第三方出售、出租或泄露您的个人信息。以下情况除外：</p>
      <ul>
        <li>经您明确同意</li>
        <li>法律法规要求</li>
        <li>为完成订单必须向源台平台传递的必要信息</li>
      </ul>

      <h2>五、账号安全</h2>
      <p>您有责任妥善保管账号密码，因个人原因导致的账号信息泄露由用户自行承担。如发现账号异常，请立即联系管理员。</p>

      <h2>六、联系我们</h2>
      <p>如有任何隐私相关问题，请通过系统工单联系我们。</p>

      <div class="back-row">
        <button class="back-btn" @click="goBack">返回</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'PrivacyPolicy' })
  const goBack = () => history.back()
</script>

<style scoped>
  .privacy-page { min-height: 100vh; background: #f5f7fa; padding: 40px 20px; }
  .privacy-box { max-width: 800px; margin: 0 auto; background: #fff; border-radius: 12px; padding: 40px; }
  h1 { font-size: 24px; margin: 0 0 8px; }
  .update-time { color: #999; font-size: 13px; margin-bottom: 30px; }
  h2 { font-size: 17px; margin: 24px 0 10px; color: #1f2937; }
  p, li { font-size: 14px; color: #4b5563; line-height: 1.8; }
  ul { padding-left: 20px; }
  .back-row { margin-top: 40px; text-align: center; }
  .back-btn { padding: 10px 40px; border-radius: 8px; border: none; background: #2563eb; color: #fff; cursor: pointer; }
</style>

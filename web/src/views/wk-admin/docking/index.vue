<template>
  <div>
    <!-- 对外 API 地址 -->
    <div class="art-card p-5 mb-5">
      <div class="art-card-header mb-4">
        <div class="title">
          <h4>对外 API 地址</h4>
          <p>第三方系统对接本平台使用</p>
        </div>
      </div>
      <ElDescriptions :column="1" border>
        <ElDescriptionsItem label="查询余额">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=getmoney</span>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="查课接口">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=get</span>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="下单接口">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=add</span>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="查单接口">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=chadan</span>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="补刷接口">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=budan</span>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="平台列表">
          <ElTag>POST</ElTag>
          <span class="ml-2">{{ baseUrl }}/open/api?act=getclass</span>
        </ElDescriptionsItem>
      </ElDescriptions>
    </div>

    <!-- 接口参数说明 -->
    <ElRow :gutter="20">
      <ElCol :md="12">
        <div class="art-card p-5 mb-5">
          <div class="art-card-header mb-4">
            <div class="title">
              <h4>查课接口参数</h4>
              <p>act=get</p>
            </div>
          </div>
          <ElTable :data="getParams" size="default" border>
            <ElTableColumn prop="name" label="参数" width="120" />
            <ElTableColumn prop="required" label="必填" width="60" />
            <ElTableColumn prop="desc" label="说明" />
          </ElTable>
        </div>
      </ElCol>
      <ElCol :md="12">
        <div class="art-card p-5 mb-5">
          <div class="art-card-header mb-4">
            <div class="title">
              <h4>下单接口参数</h4>
              <p>act=add</p>
            </div>
          </div>
          <ElTable :data="addParams" size="default" border>
            <ElTableColumn prop="name" label="参数" width="120" />
            <ElTableColumn prop="required" label="必填" width="60" />
            <ElTableColumn prop="desc" label="说明" />
          </ElTable>
        </div>
      </ElCol>
    </ElRow>

    <!-- 响应示例 -->
    <div class="art-card p-5 mb-5">
      <div class="art-card-header mb-4">
        <div class="title">
          <h4>响应示例</h4>
          <p>统一 JSON 格式</p>
        </div>
      </div>
      <ElTabs>
        <ElTabPane label="成功响应">
          <pre class="bg-g-100 p-3 rounded text-sm">{{ successExample }}</pre>
        </ElTabPane>
        <ElTabPane label="失败响应">
          <pre class="bg-g-100 p-3 rounded text-sm">{{ errorExample }}</pre>
        </ElTabPane>
      </ElTabs>
    </div>

    <!-- 状态码 -->
    <div class="art-card p-5 mb-5">
      <div class="art-card-header mb-4">
        <div class="title">
          <h4>状态码说明</h4>
        </div>
      </div>
      <ElTable :data="codeList" size="default" border>
        <ElTableColumn prop="code" label="code" width="100" />
        <ElTableColumn prop="msg" label="说明" />
      </ElTable>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue'

  defineOptions({ name: 'WkAdminDocking' })

  const baseUrl = computed(() => `${window.location.protocol}//${window.location.host}`)

  const getParams = [
    { name: 'uid', required: '是', desc: '用户 UID' },
    { name: 'key', required: '是', desc: '对接密钥' },
    { name: 'platform', required: '是', desc: '平台 cid（getclass 接口获取）' },
    { name: 'school', required: '是', desc: '学校名称' },
    { name: 'user', required: '是', desc: '学生账号' },
    { name: 'pass', required: '是', desc: '学生密码' }
  ]

  const addParams = [
    { name: 'uid', required: '是', desc: '用户 UID' },
    { name: 'key', required: '是', desc: '对接密钥' },
    { name: 'platform', required: '是', desc: '平台 cid' },
    { name: 'school', required: '是', desc: '学校名称' },
    { name: 'user', required: '是', desc: '学生账号' },
    { name: 'pass', required: '是', desc: '学生密码' },
    { name: 'kcid', required: '否', desc: '课程 ID（多个用逗号分隔）' },
    { name: 'kcname', required: '是', desc: '课程名称（多个用逗号分隔）' }
  ]

  const successExample = `{
  "code": 1,
  "msg": "查询成功",
  "data": [
    { "id": "12345", "name": "高等数学" }
  ]
}`

  const errorExample = `{
  "code": -1,
  "msg": "余额不足"
}`

  const codeList = [
    { code: '1 / 0', msg: '成功' },
    { code: '0', msg: '参数为空' },
    { code: '-1', msg: '业务失败（余额不足、账号错误等）' },
    { code: '-2', msg: '密钥错误或平台已下架' }
  ]
</script>

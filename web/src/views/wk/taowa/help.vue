<template>
  <div class="taowa-help-page">
    <div class="art-card p-5">
      <div class="flex-cb mb-4">
        <h4 class="font-bold">帮助文档</h4>
        <ElButton size="small" @click="open">新窗口打开</ElButton>
      </div>
      <div class="help-frame">
        <iframe :src="url" class="w-full h-full border-0 rounded-md" loading="lazy"></iframe>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'WkTaowaHelp' })
  const url = 'https://www.sxzsgy.top/bz.php'
  const open = () => window.open(url, '_blank')
</script>

<style scoped>
  .help-frame {
    position: relative;
    width: 100%;
    height: 70vh;
    background: #f7f8fa;
    border-radius: 8px;
    overflow: hidden;
  }
</style>

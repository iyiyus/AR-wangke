<template>
  <ElRow :gutter="20" class="flex">
    <ElCol v-for="(item, index) in dataList" :key="index" :sm="12" :md="6" :lg="6">
      <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
        <span class="text-g-700 text-sm">{{ item.des }}</span>
        <ArtCountTo
          v-if="item.isMoney"
          class="text-[26px] font-medium mt-2"
          :prefix="'¥ '"
          :decimals="2"
          :target="item.num"
          :duration="1300"
        />
        <ArtCountTo
          v-else
          class="text-[26px] font-medium mt-2"
          :target="item.num"
          :duration="1300"
        />
        <div class="flex-c mt-1">
          <span class="text-xs text-g-600">{{ item.subTitle }}</span>
          <span class="ml-1 text-xs font-semibold" :class="item.changeColor">
            {{ item.change }}
          </span>
        </div>
        <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
          <ArtSvgIcon :icon="item.icon" class="text-xl text-theme" />
        </div>
      </div>
    </ElCol>
  </ElRow>
</template>

<script setup lang="ts">
  import { reactive, onMounted } from 'vue'
  import { getUserInfo } from '@/api/wk'

  interface CardDataItem {
    des: string
    icon: string
    num: number
    change: string
    subTitle: string
    changeColor: string
    isMoney?: boolean
  }

  const dataList = reactive<CardDataItem[]>([
    { des: '账户余额', icon: 'ri:wallet-3-line', num: 0, change: '总充值', subTitle: '', changeColor: 'text-g-600', isMoney: true },
    { des: '查课次数', icon: 'ri:search-eye-line', num: 0, change: '0.0%', subTitle: '下单率', changeColor: 'text-success' },
    { des: '订单数量', icon: 'ri:list-check-2', num: 0, change: '单', subTitle: '累计', changeColor: 'text-g-600' },
    { des: '下单折扣', icon: 'ri:vip-crown-line', num: 1, change: 'UID', subTitle: '', changeColor: 'text-g-600' }
  ])

  onMounted(async () => {
    const info: any = await getUserInfo()
    dataList[0].num = Number(info.money) || 0
    dataList[0].subTitle = `¥ ${info.zcz || 0}`
    dataList[0].change = ''
    dataList[1].num = info.ck || 0
    dataList[1].change = `${(info.lv || 0).toFixed(1)}%`
    dataList[2].num = info.dd || 0
    dataList[2].change = ''
    dataList[2].subTitle = '订单'
    dataList[3].num = Number(info.addprice) || 1
    dataList[3].change = String(info.uid || '-')
  })
</script>

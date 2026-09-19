import { AppRouteRecord } from '@/types/router'

/**
 * 网课平台菜单（原一级菜单已取消，子项全部提升为一级菜单）
 */
export const wkRoutes: AppRouteRecord[] = [
  {
    path: '/wk/dashboard',
    name: 'WkDashboard',
    component: '/wk/dashboard',
    meta: { title: '用户中心', icon: 'ri:home-smile-2-line', fixedTab: true }
  },
  {
    path: '/wk/order',
    name: 'WkOrder',
    component: '/wk/order',
    meta: { title: '任务列表', icon: 'ri:list-check-2' }
  },
  {
    path: '/wk/submit',
    name: 'WkSubmit',
    component: '/wk/submit',
    meta: { title: '提交任务', icon: 'ri:add-circle-line' }
  },
  {
    path: '/wk/pay',
    name: 'WkPay',
    component: '/wk/pay',
    meta: { title: '账户充值', icon: 'ri:wallet-3-line' }
  },
  {
    path: '/wk/paylist',
    name: 'WkPayList',
    component: '/wk/paylist',
    meta: { title: '充值记录', icon: 'ri:bill-line' }
  },
  {
    path: '/wk/log',
    name: 'WkLog',
    component: '/wk/log',
    meta: { title: '消费记录', icon: 'ri:file-list-3-line' }
  },
  {
    path: '/wk/gongdan',
    name: 'WkGongDan',
    component: '/wk/gongdan',
    meta: { title: '工单系统', icon: 'ri:customer-service-2-line' }
  },
  {
    path: '/wk/agents',
    name: 'WkAgents',
    component: '/wk/agents',
    meta: { title: '代理管理', icon: 'ri:team-line' }
  },
  {
    path: '/wk/docking',
    name: 'WkDocking',
    component: '/wk/docking',
    meta: { title: '平台对接', icon: 'ri:plug-line' }
  },
  {
    path: '/wk/myprice',
    name: 'WkMyPrice',
    component: '/wk/myprice',
    meta: { title: '我的价格', icon: 'ri:price-tag-3-line' }
  },
  {
    path: '/wk/charge',
    name: 'WkCharge',
    component: '/wk/charge',
    meta: { title: '联系上级', icon: 'ri:contacts-line' }
  },
  {
    path: '/wk/taowa/add',
    name: 'WkTaowaAdd',
    component: '/wk/taowa/add',
    meta: { title: '实习提交', icon: 'ri:add-circle-line' }
  },
  {
    path: '/wk/taowa/list',
    name: 'WkTaowaList',
    component: '/wk/taowa/list',
    meta: { title: '盖章任务', icon: 'ri:list-check-2' }
  },
  {
    path: '/wk/taowa/add-bl',
    name: 'WkTaowaAddBl',
    component: '/wk/taowa/add-bl',
    meta: { title: '病历提交', icon: 'ri:file-add-line' }
  },
  {
    path: '/wk/taowa/list-bl',
    name: 'WkTaowaListBl',
    component: '/wk/taowa/list-bl',
    meta: { title: '病历任务', icon: 'ri:file-list-line' }
  },
  {
    path: '/wk/taowa/notice',
    name: 'WkTaowaNotice',
    component: '/wk/taowa/notice',
    meta: { title: '源台公告', icon: 'ri:megaphone-line' }
  },
  {
    path: '/wk/taowa/help',
    name: 'WkTaowaHelp',
    component: '/wk/taowa/help',
    meta: { title: '帮助文档', icon: 'ri:question-line' }
  },
  {
    path: '/wk/daka',
    name: 'WkDaka',
    component: '/wk/daka/index',
    meta: { title: '实习打卡', icon: 'ri:calendar-check-line' }
  }
]

export const wkAdminRoutes: AppRouteRecord = {
  name: 'WkAdmin',
  path: '/wk-admin',
  component: '/index/index',
  meta: {
    title: '管理后台',
    icon: 'ri:settings-3-line',
    roles: ['R_SUPER']
  },
  children: [
    {
      path: 'users',
      name: 'WkAdminUsers',
      component: '/wk-admin/users',
      meta: { title: '代理管理', icon: 'ri:team-line' }
    },
    {
      path: 'classes',
      name: 'WkAdminClasses',
      component: '/wk-admin/classes',
      meta: { title: '商品管理', icon: 'ri:apps-2-line' }
    },
    {
      path: 'orders',
      name: 'WkAdminOrders',
      component: '/wk-admin/orders',
      meta: { title: '订单管理', icon: 'ri:order-play-line' }
    },
    {
      path: 'dengji',
      name: 'WkAdminDengji',
      component: '/wk-admin/dengji',
      meta: { title: '等级管理', icon: 'ri:vip-crown-line' }
    },
    {
      path: 'myprice',
      name: 'WkAdminMyPrice',
      component: '/wk-admin/myprice',
      meta: { title: '密价管理', icon: 'ri:price-tag-2-line' }
    },
    {
      path: 'gongdan',
      name: 'WkAdminGongDan',
      component: '/wk-admin/gongdan',
      meta: { title: '工单管理', icon: 'ri:customer-service-line' }
    },
    {
      path: 'config',
      name: 'WkAdminConfig',
      component: '/wk-admin/config',
      meta: { title: '系统设置', icon: 'ri:settings-4-line' }
    },
    {
      path: 'huoyuan',
      name: 'WkAdminHuoYuan',
      component: '/wk-admin/huoyuan',
      meta: { title: '货源管理', icon: 'ri:database-2-line' }
    },
    {
      path: 'fenlei',
      name: 'WkAdminFenLei',
      component: '/wk-admin/fenlei',
      meta: { title: '分类管理', icon: 'ri:folder-line' }
    },
    {
      path: 'log',
      name: 'WkAdminLog',
      component: '/wk-admin/log',
      meta: { title: '系统日志', icon: 'ri:file-list-line' }
    },
    {
      path: 'data',
      name: 'WkAdminData',
      component: '/wk-admin/data',
      meta: { title: '数据统计', icon: 'ri:bar-chart-2-line' }
    },
    {
      path: 'docking',
      name: 'WkAdminDocking',
      component: '/wk-admin/docking',
      meta: { title: '平台对接', icon: 'ri:plug-line' }
    }
  ]
}

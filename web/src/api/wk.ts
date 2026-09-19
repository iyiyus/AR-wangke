import request from '@/utils/http'

// ===== 实习打卡 daka（源台协议 code=0 成功，绕过统一拦截器） =====
import axios from 'axios'
import { useUserStore } from '@/store/modules/user'

export const dakaApi = (act: string, params: Record<string, any> = {}) => {
  const { accessToken } = useUserStore()
  return axios
    .post(`/api/daka/api?act=${act}`, params, {
      headers: { Authorization: `Bearer ${accessToken}` },
      timeout: 120000
    })
    .then((r) => r.data)
}

// ===== 盖章/病历 taowa =====
export const getTaowaCompanies = () =>
  request.get<any>({ url: '/api/taowa/companies' })

export const getTaowaTemplates = () =>
  request.get<any>({ url: '/api/taowa/templates' })

export const getTaowaSpec = (params: { companyName: string; userinfo?: string }) =>
  request.post<any>({ url: '/api/taowa/spec', params })

export const getTaowaDelivery = () =>
  request.get<any>({ url: '/api/taowa/delivery' })

export const getTaowaNotice = () =>
  request.get<any>({ url: '/api/taowa/notice' })

export const addTaowaOrder = (params: any) =>
  request.post<any>({ url: '/api/taowa/order/add', params })

export const addTaowaOrderBL = (params: any) =>
  request.post<any>({ url: '/api/taowa/order/add-bl', params })

export const getTaowaOrderList = (params: Record<string, any>) =>
  request.get<any>({ url: '/api/taowa/order/list', params })

export const cancelTaowaOrder = (oid: number) =>
  request.post({ url: '/api/taowa/order/cancel', params: { oid } })

export const taowaOrderRemark = (oid: number, remark: string) =>
  request.post({ url: '/api/taowa/order/remark', params: { oid, remark } })

export const taowaOrderGxx = (oid: number, fileName: string) =>
  request.post({ url: '/api/taowa/order/gxx', params: { oid, fileName } })

export const taowaOrderStatus = (oid: number, status: string) =>
  request.post({ url: '/api/taowa/order/status', params: { oid, status } })

export const taowaOrderTk = (oid: number) =>
  request.post({ url: '/api/taowa/order/tk', params: { oid } })

export const taowaOrderBjtk = (oid: number) =>
  request.post({ url: '/api/taowa/order/bjtk', params: { oid } })

export const taowaOrderDel = (oid: number) =>
  request.post({ url: '/api/taowa/order/del', params: { oid } })

export const taowaDdlog = (sfdh: string) =>
  request.get<any>({ url: '/api/taowa/order/ddlog', params: { sfdh } })


// ===== 认证 =====
export const login = (params: Api.Auth.LoginParams) =>
  request.post<Api.Auth.LoginResponse>({ url: '/api/auth/login', params })

export const register = (params: { name: string; account: string; pass: string; yqm?: string }) =>
  request.post({ url: '/api/auth/register', params })

export const getUserInfo = () =>
  request.get<Api.Auth.UserInfo>({ url: '/api/user/info' })

export const changePassword = (params: { old_pass: string; new_pass: string }) =>
  request.post({ url: '/api/user/passwd', params })

export const openAPIKey = () =>
  request.post<{ key: string }>({ url: '/api/user/api-key/open', params: {} })

export const resetAPIKey = () =>
  request.post<{ key: string }>({ url: '/api/user/api-key/reset', params: {} })

export const updateNotice = (notice: string) =>
  request.post({ url: '/api/user/notice', params: { notice } })

export const setInvite = (params: { yqm: string; yqprice: string }) =>
  request.post({ url: '/api/user/invite', params })

// ===== 平台 =====
export const getClassAll = () =>
  request.get<Api.Class.Item[]>({ url: '/api/class/all' })

export const getClassList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/class/list', params })

export const addClass = (params: Partial<Api.Class.Item>) =>
  request.post({ url: '/api/admin/class', params })

export const updateClass = (id: number, params: Partial<Api.Class.Item>) =>
  request.put({ url: `/api/admin/class/${id}`, params })

export const deleteClass = (id: number) =>
  request.del({ url: `/api/admin/class/${id}` })

// ===== 分类 =====
export const getFenLeiList = () =>
  request.get<Api.Class.Item[]>({ url: '/api/fenlei/list' })

export const addFenLei = (params: any) =>
  request.post({ url: '/api/admin/fenlei', params })

export const updateFenLei = (id: number, params: any) =>
  request.put({ url: `/api/admin/fenlei/${id}`, params })

export const deleteFenLei = (id: number) =>
  request.del({ url: `/api/admin/fenlei/${id}` })

// ===== 订单 =====
export const getOrderList = (params: Api.Order.QueryParams) =>
  request.get<any>({ url: '/api/order/list', params })

export const addOrder = (params: Api.Order.AddParams) =>
  request.post({ url: '/api/order/add', params })

export const cancelOrder = (oid: number) =>
  request.post({ url: '/api/order/cancel', params: { oid } })

export const restartOrder = (oid: number) =>
  request.post({ url: '/api/order/restart', params: { oid } })

export const exportOrder = (params: Record<string, any>) =>
  request.get<string>({ url: '/api/order/export', params })

export const batchOrderStatus = (oids: number[], status: string) =>
  request.post({ url: '/api/admin/order/status', params: { oids, status } })

export const batchOrderDock = (oids: number[], dock: string) =>
  request.post({ url: '/api/admin/order/dock', params: { oids, dock } })

export const refundOrder = (oids: number[]) =>
  request.post({ url: '/api/admin/order/refund', params: { oids } })

// ===== 充值 =====
export const createPay = (params: { money: string; type: string }) =>
  request.post<{ pay_url: string; out_trade_no: string }>({ url: '/api/pay/create', params })

export const getPayList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/pay/list', params })

// ===== 工单 =====
export const getGongDanList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/gongdan/list', params })

export const addGongDan = (params: { oid?: number; region: string; title?: string; content: string }) =>
  request.post({ url: '/api/gongdan/add', params })

export const replyGongDan = (params: { gid: number; answer: string; state?: string }) =>
  request.post({ url: '/api/admin/gongdan/reply', params })

// ===== 密价 =====
export const getMyPriceList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/myprice/list', params })

export const setMyPrice = (params: any) =>
  request.post({ url: '/api/admin/myprice', params })

export const deleteMyPrice = (id: number) =>
  request.del({ url: `/api/admin/myprice/${id}` })

// ===== 管理员 =====
export const getAdminStats = () =>
  request.get<any>({ url: '/api/admin/stats' })

export const getAdminUsers = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/admin/users', params })

export const adminAddUser = (params: any) =>
  request.post({ url: '/api/admin/user/add', params })

export const adminRecharge = (uid: number, money: number) =>
  request.post({ url: '/api/admin/user/recharge', params: { uid, money } })

export const adminBanUser = (uid: number, active: string) =>
  request.post({ url: '/api/admin/user/ban', params: { uid, active } })

export const adminResetPass = (uid: number, new_pass: string) =>
  request.post({ url: '/api/admin/user/reset-pass', params: { uid, new_pass } })

export const adminSetLevel = (uid: number, addprice: number) =>
  request.post({ url: '/api/admin/user/level', params: { uid, addprice } })

export const adminOpenKey = (uid: number) =>
  request.post<{ key: string }>({ url: '/api/admin/user/api-key', params: { uid } })

export const adminSetYQM = (uid: number, yqm: string) =>
  request.post({ url: '/api/admin/user/yqm', params: { uid, yqm } })

export const getAdminConfig = () =>
  request.get<Record<string, string>>({ url: '/api/admin/config' })

export const saveAdminConfig = (params: Record<string, string>) =>
  request.post({ url: '/api/admin/config', params })

export const syncTaowa = () =>
  request.post<any>({ url: '/api/admin/taowa/sync' })

export const getDengjiList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/dengji/list', params })

export const addDengji = (params: Partial<Api.Dengji.Item>) =>
  request.post({ url: '/api/admin/dengji', params })

export const updateDengji = (id: number, params: Partial<Api.Dengji.Item>) =>
  request.put({ url: `/api/admin/dengji/${id}`, params })

export const deleteDengji = (id: number) =>
  request.del({ url: `/api/admin/dengji/${id}` })

export const getAdminLog = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/admin/log', params })


// ===== 货源管理 =====
export const getHuoYuanList = (params?: Record<string, any>) =>
  request.get<any>({ url: '/api/admin/huoyuan', params })

export const addHuoYuan = (params: any) =>
  request.post({ url: '/api/admin/huoyuan', params })

export const updateHuoYuan = (id: number, params: any) =>
  request.put({ url: `/api/admin/huoyuan/${id}`, params })

export const deleteHuoYuan = (id: number) =>
  request.del({ url: `/api/admin/huoyuan/${id}` })


// ===== 分类管理 =====
export const addFenLeiAdmin = (params: any) =>
  request.post({ url: '/api/admin/fenlei', params })

export const updateFenLeiAdmin = (id: number, params: any) =>
  request.put({ url: `/api/admin/fenlei/${id}`, params })

export const deleteFenLeiAdmin = (id: number) =>
  request.del({ url: `/api/admin/fenlei/${id}` })


// ===== 数据统计（管理员）=====
export const getAdminDashboard = () =>
  request.get<any>({ url: '/api/admin/dashboard' })


// ===== 平台对接 =====
export const queryCourse = (params: { cid: number; school: string; user: string; pass: string }) =>
  request.post<any>({ url: '/api/order/query', params })

export const manualDock = (oid: number) =>
  request.post({ url: '/api/admin/order/manual-dock', params: { oid } })

export const syncProgress = (oid: number) =>
  request.post<any[]>({ url: '/api/admin/order/sync-progress', params: { oid } })


// ===== 批量导入平台 =====
export const fetchRemoteClass = (hid: number) =>
  request.get<{ cid: string; name: string; price: string; noun: string; fenlei: string; content: string }[]>(
    { url: '/api/admin/class/remote', params: { hid } })

export const batchImportClass = (params: {
  hid: number
  fenlei: string
  price: string
  use_remote: boolean
  markup_mode: string
  price_markup: string
  items: { cid: string; name: string; price: string; noun: string; fenlei: string; content: string }[]
}) => request.post<{ imported: number; total: number }>({ url: '/api/admin/class/batch-import', params })


// ===== 批量删除 =====
export const batchDeleteClass = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/class/batch-delete', params: { ids } })
export const batchDeleteFenLei = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/fenlei/batch-delete', params: { ids } })
export const batchDeleteDengji = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/dengji/batch-delete', params: { ids } })
export const batchDeleteHuoYuan = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/huoyuan/batch-delete', params: { ids } })
export const batchDeleteMyPrice = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/myprice/batch-delete', params: { ids } })
export const batchDeleteGongDan = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/gongdan/batch-delete', params: { ids } })
export const batchDeleteUsers = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/users/batch-delete', params: { ids } })
export const batchDeleteOrder = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/order/batch-delete', params: { ids } })
export const batchDeleteLog = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/log/batch-delete', params: { ids } })
export const batchDeletePay = (ids: number[]) =>
  request.post<{ deleted: number }>({ url: '/api/admin/pay/batch-delete', params: { ids } })


// ===== 代理管理（下级） =====
export const getProxyUsers = (params?: Record<string, any>) => request.get<any>({ url: '/api/proxy/user/list', params })
export const proxyAddUser = (params: { name: string; user: string; pass: string; addprice: number }) => request.post({ url: '/api/proxy/user/add', params })
export const proxyRecharge = (uid: number, money: number) => request.post({ url: '/api/proxy/user/recharge', params: { uid, money } })
export const proxySetPrice = (uid: number, addprice: number) => request.post({ url: '/api/proxy/user/price', params: { uid, addprice } })
export const proxyResetPass = (uid: number) => request.post({ url: '/api/proxy/user/reset-pass', params: { uid } })
export const proxyBanUser = (uid: number, active: string) => request.post({ url: '/api/proxy/user/ban', params: { uid, active } })
export const proxySetYQM = (uid: number, yqm: string) => request.post({ url: '/api/proxy/user/yqm', params: { uid, yqm } })
export const proxyOpenKey = (uid: number) => request.post<{ key: string }>({ url: '/api/proxy/user/api-key', params: { uid } })
export const toAnswerGongDan = (gid: number, content: string) => request.post({ url: '/api/gongdan/toanswer', params: { gid, content } })

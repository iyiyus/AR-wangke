import request from '@/utils/http'

// ===== 安装系统（复刻 UDID 安装引导） =====
export const getInstallStatus = () =>
  request.get<any>({ url: '/api/install/status' })

export const checkInstallDB = (params: Record<string, any>) =>
  request.post<any>({ url: '/api/install/check-db', params })

export const registerInstallSystemd = (params: Record<string, any>) =>
  request.post<any>({ url: '/api/install/register-systemd', params })

export const runInstall = (params: Record<string, any>) =>
  request.post<any>({ url: '/api/install/run', params })

export const restartWkService = (serviceName: string) =>
  request.post({ url: '/api/install/restart', params: { serviceName } })

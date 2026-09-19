import { AppRouteRecord } from '@/types/router'
import { wkRoutes, wkAdminRoutes } from './wk'
import { systemRoutes } from './system'

/**
 * 导出所有模块化路由
 */
export const routeModules: AppRouteRecord[] = [
  ...wkRoutes,
  wkAdminRoutes,
  systemRoutes
]

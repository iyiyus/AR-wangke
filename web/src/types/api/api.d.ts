/**
 * API 接口类型定义模块
 *
 * 提供所有后端接口的类型定义
 *
 * ## 主要功能
 *
 * - 通用类型（分页参数、响应结构等）
 * - 认证类型（登录、用户信息等）
 * - 系统管理类型（用户、角色等）
 * - 全局命名空间声明
 *
 * ## 使用场景
 *
 * - API 请求参数类型约束
 * - API 响应数据类型定义
 * - 接口文档类型同步
 *
 * ## 注意事项
 *
 * - 在 .vue 文件使用需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 * - 使用全局命名空间，无需导入即可使用
 *
 * ## 使用方式
 *
 * ```typescript
 * const params: Api.Auth.LoginParams = { userName: 'admin', password: '123456' }
 * const response: Api.Auth.UserInfo = await fetchUserInfo()
 * ```
 *
 * @module types/api/api
 * @author Art Design Pro Team
 */

declare namespace Api {
  /** 通用类型 */
  namespace Common {
    interface PageResult<T = any> {
      list: T[]
      total: number
      page: number
      pageSize: number
    }
    type CommonSearchParams = {
      page?: number
      pageSize?: number
    }
    /** 分页参数 */
    interface PaginationParams {
      current: number
      size: number
      total: number
    }
    type EnableStatus = '1' | '2'
    interface PaginatedResponse<T = any> {
      records: T[]
      current: number
      size: number
      total: number
    }
  }

  /** 认证类型 */
  namespace Auth {
    interface LoginParams {
      account: string
      pass: string
      auth_pass?: string
    }
    interface LoginResponse {
      token: string
      user: UserInfo
    }
    interface UserInfo {
      uid: number
      uuid: number
      user: string
      name: string
      money: number
      zcz: string
      ck: number
      dd: number
      lv: number
      addprice: number
      key: string
      yqm: string
      yqprice: string
      notice: string
      active: string
      faceimg: string
      nickname: string
    }
  }

  /** 网课平台 */
  namespace Class {
    interface Item {
      cid: number
      sort: number
      name: string
      getnoun: string
      noun: string
      price: string
      queryplat: string
      docking: string
      yunsuan: string
      content: string
      addtime: string
      status: number
      fenlei: string
      my_price?: string
      my_mode?: number
    }
  }

  /** 订单 */
  namespace Order {
    interface Item {
      oid: number
      uid: number
      cid: number
      ptname: string
      school: string
      user: string
      pass: string
      kcid: string
      kcname: string
      fees: string
      status: string
      dockstatus: string
      process: string
      remarks: string
      addtime: string
      bsnum: string
    }
    interface AddParams {
      cid: number
      school: string
      user: string
      pass: string
      kcid?: string
      kcname: string
    }
    interface QueryParams {
      current?: number
      size?: number
      page?: number
      pageSize?: number
      oid?: string
      uid?: string
      qq?: string
      cid?: string
      status_text?: string
      dock?: string
      [key: string]: any
    }
  }

  /** 用户（管理员视角） */
  namespace User {
    interface Item {
      uid: number
      uuid: number
      user: string
      name: string
      money: number
      addprice: number
      key: string
      yqm: string
      yqprice: string
      active: string
      endtime: string
    }
  }

  /** 等级 */
  namespace Dengji {
    interface Item {
      id: number
      sort: string
      name: string
      rate: number
      money: number
      status: string
    }
  }

  /** 工单 */
  namespace GongDan {
    interface Item {
      gid: number
      uid: number
      oid: number
      region: string
      title: string
      content: string
      answer: string
      state: string
      addtime: string
    }
  }

  /** 充值记录 */
  namespace Pay {
    interface Item {
      oid: number
      out_trade_no: string
      type: string
      uid: number
      money: string
      status: number
      addtime: string
    }
  }

  /** 系统管理类型（保留兼容） */
  namespace SystemManage {
    type UserList = Api.Common.PaginatedResponse<UserListItem>
    interface UserListItem {
      id: number
      userName: string
      status: string
    }
    type UserSearchParams = Partial<Common.CommonSearchParams>
    type RoleList = Api.Common.PaginatedResponse<RoleListItem>
    interface RoleListItem {
      roleId: number
      roleName: string
    }
    type RoleSearchParams = Partial<Common.CommonSearchParams>
  }
}

import service from '@/utils/request'

export interface AuthorityListParams {
  page?: number
  pageSize?: number
}

export interface Authority {
  authorityId: string | number
  authorityName: string
  parentId: string | number
  defaultRouterString?: string
  defaultRouter?: string
  dataAuthorityId?: Array<string | number>
  children?: Authority[]
  menus?: unknown[]
}

export interface CreateAuthorityParams {
  authorityId: string | number
  authorityName: string
  parentId: string | number
  defaultRouter?: string
}

export interface CopyAuthorityParams {
  authority: CreateAuthorityParams & {
    dataAuthorityId?: Array<string | number>
  }
  oldAuthorityId: number | string
}

// @Router /api/v1/authority/list [post]
export const getAuthorityList = (data?: AuthorityListParams) => {
  return service({
    url: '/api/v1/authority/list',
    method: 'post',
    data
  })
}

// @Summary 删除角色
export const deleteAuthority = (data: { authorityId: number }) => {
  return service({
    url: '/api/v1/authority/delete',
    method: 'post',
    data: data
  })
}

// @Summary 创建角色
export const createAuthority = (data: CreateAuthorityParams) => {
  return service({
    url: '/api/v1/authority/add',
    method: 'post',
    data: data
  })
}

// @Tags authority
// @Summary 拷贝角色
export const copyAuthority = (data: CopyAuthorityParams) => {
  return service({
    url: '/api/v1/authority/copy',
    method: 'post',
    data: data
  })
}

// @Summary 设置角色资源权限
export const setDataAuthority = (data: {
  authorityId: number
  dataAuthorityId: any[]
}) => {
  return service({
    url: '/api/v1/authority/data/authority/update',
    method: 'post',
    data: data
  })
}

// @Summary 修改角色
export const updateAuthority = (data: Authority) => {
  return service({
    url: '/api/v1/authority/update',
    method: 'post',
    data
  })
}

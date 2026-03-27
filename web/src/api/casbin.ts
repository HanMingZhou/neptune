import service from '@/utils/request'

export interface CasbinInfo {
    path: string
    method: string
}

export interface UpdateCasbinData {
    authorityId: string | number
    casbinInfos: CasbinInfo[]
}

// @Tags authority
// @Summary 更改角色api权限
export const UpdateCasbin = (data: UpdateCasbinData) => {
    return service({
        url: '/api/v1/casbin/update',
        method: 'post',
        data: data
    })
}

// @Tags casbin
// @Summary 获取权限列表
export const getPolicyPathByAuthorityId = (data: { authorityId: string | number }) => {
    return service({
        url: '/api/v1/casbin/policy/list',
        method: 'post',
        data
    })
}

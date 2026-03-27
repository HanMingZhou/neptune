import service from '@/utils/request'

export interface ApiData {
    path: string
    description: string
    apiGroup: string
    method: string
}

export interface ApiListParams {
    page: number
    pageSize: number
    path?: string
    description?: string
    apiGroup?: string
    method?: string
    orderKey?: string
    desc?: boolean
}

export interface GetByIdParams {
    id: number | string
}

export interface IdsReq {
    ids: number[] | string[]
}

// @Tags api
// @Summary 分页获取角色列表
export const getApiList = (data: ApiListParams) => {
    return service({
        url: '/api/v1/api/list',
        method: 'post',
        data
    })
}

export const createApi = (data: ApiData) => {
    return service({
        url: '/api/v1/api/add',
        method: 'post',
        data
    })
}

export const getApiById = (data: GetByIdParams) => {
    return service({
        url: '/api/v1/api/get',
        method: 'post',
        data
    })
}

export const updateApi = (data: ApiData) => {
    return service({
        url: '/api/v1/api/update',
        method: 'post',
        data: data
    })
}

export const getAllApis = (data?: any) => {
    return service({
        url: '/api/v1/api/all',
        method: 'post',
        data: data
    })
}

export const deleteApi = (data: ApiData) => {
    return service({
        url: '/api/v1/api/delete',
        method: 'post',
        data: data
    })
}

export const deleteApisByIds = (data: IdsReq) => {
    return service({
        url: '/api/v1/api/delete/multi',
        method: 'post',
        data
    })
}

export const freshCasbin = () => {
    return service({
        url: '/api/v1/api/casbin/fresh',
        method: 'get'
    })
}

export const syncApi = () => {
    return service({
        url: '/api/v1/api/sync',
        method: 'get'
    })
}

export const getApiGroups = () => {
    return service({
        url: '/api/v1/api/group/list',
        method: 'get'
    })
}

export const ignoreApi = (data: any) => {
    return service({
        url: '/api/v1/api/ignore',
        method: 'post',
        data
    })
}

export const enterSyncApi = (data: any) => {
    return service({
        url: '/api/v1/api/sync/enter',
        method: 'post',
        data
    })
}

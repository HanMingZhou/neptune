import service from '@/utils/request'

export interface SysParams {
    ID?: number
    configName?: string
    configKey?: string
    configValue?: string
    configType?: string
    desc?: string
}

export interface SysParamsSearch {
    page: number
    pageSize: number
    configName?: string
    configKey?: string
    configType?: string
}

// @Tags SysParams
// @Summary 创建参数
export const createSysParams = (data: SysParams) => {
    return service({
        url: '/api/v1/params/add',
        method: 'post',
        data
    })
}

// @Tags SysParams
// @Summary 删除参数
export const deleteSysParams = (params: { ID: number }) => {
    return service({
        url: '/api/v1/params/delete',
        method: 'post',
        data: params
    })
}

// @Tags SysParams
// @Summary 批量删除参数
export const deleteSysParamsByIds = (params: { ids: number[] }) => {
    return service({
        url: '/api/v1/params/delete/multi',
        method: 'post',
        data: params
    })
}

// @Tags SysParams
// @Summary 更新参数
export const updateSysParams = (data: SysParams) => {
    return service({
        url: '/api/v1/params/update',
        method: 'post',
        data
    })
}

// @Tags SysParams
// @Summary 用id查询参数
export const findSysParams = (params: { ID: number }) => {
    return service({
        url: '/api/v1/params/get',
        method: 'get',
        params
    })
}

// @Tags SysParams
// @Summary 分页获取参数列表
export const getSysParamsList = (params: SysParamsSearch) => {
    return service({
        url: '/api/v1/params/list',
        method: 'get',
        params
    })
}

// @Tags SysParams
// @Summary 不需要鉴权的参数接口
export const getSysParam = (params: { configKey: string }) => {
    return service({
        url: '/api/v1/params/get/key',
        method: 'get',
        params
    })
}

import service from '@/utils/request'

export interface SysError {
    ID?: number
    // Add properties
}

export interface SysErrorSearch {
    page: number
    pageSize: number
    // Add filter
}

// @Tags SysError
// @Summary 创建错误日志
export const createSysError = (data: any) => {
    return service({
        url: '/api/v1/error/add',
        method: 'post',
        data
    })
}

// @Tags SysError
// @Summary 删除错误日志
export const deleteSysError = (params: { ID: number }) => {
    return service({
        url: '/api/v1/error/delete',
        method: 'post',
        data: params
    })
}

// @Tags SysError
// @Summary 批量删除错误日志
export const deleteSysErrorByIds = (params: { ids: number[] }) => {
    return service({
        url: '/api/v1/error/delete/multi',
        method: 'post',
        data: params
    })
}

// @Tags SysError
// @Summary 更新错误日志
export const updateSysError = (data: any) => {
    return service({
        url: '/api/v1/error/update',
        method: 'post',
        data
    })
}

// @Tags SysError
// @Summary 用id查询错误日志
export const findSysError = (params: { ID: number }) => {
    return service({
        url: '/api/v1/error/get',
        method: 'get',
        params
    })
}

// @Tags SysError
// @Summary 分页获取错误日志列表
export const getSysErrorList = (params: SysErrorSearch) => {
    return service({
        url: '/api/v1/error/list',
        method: 'get',
        params
    })
}


// @Tags SysError
// @Summary 触发错误处理（异步）
export const getSysErrorSolution = (params: { id: string | number }) => {
    return service({
        url: '/api/v1/error/solution/get',
        method: 'get',
        params
    })
}

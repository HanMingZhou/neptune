import service from '@/utils/request'

export interface SysVersion {
    ID?: number
    // Add properties
    [key: string]: any
}

// @Tags SysVersion
// @Summary 删除版本管理
export const deleteSysVersion = (params: { ID: number }) => {
    return service({
        url: '/api/v1/version/delete',
        method: 'post',
        data: params
    })
}

// @Tags SysVersion
// @Summary 批量删除版本管理
export const deleteSysVersionByIds = (params: { ids: number[] }) => {
    return service({
        url: '/api/v1/version/delete/multi',
        method: 'post',
        data: params
    })
}

// @Tags SysVersion
// @Summary 用id查询版本管理
export const findSysVersion = (params: { ID: number }) => {
    return service({
        url: '/api/v1/version/get',
        method: 'get',
        params
    })
}

// @Tags SysVersion
// @Summary 分页获取版本管理列表
export const getSysVersionList = (params: { page: number; pageSize: number }) => {
    return service({
        url: '/api/v1/version/list',
        method: 'get',
        params
    })
}

// @Tags SysVersion
// @Summary 导出版本数据
export const exportVersion = (data: any) => {
    return service({
        url: '/api/v1/version/export',
        method: 'post',
        data
    })
}

// @Tags SysVersion
// @Summary 下载版本JSON数据
export const downloadVersionJson = (params: { ID: string | number }) => {
    return service({
        url: '/api/v1/version/download/json',
        method: 'get',
        params,
        responseType: 'blob'
    })
}

// @Tags SysVersion
// @Summary 导入版本数据
export const importVersion = (data: any) => {
    return service({
        url: '/api/v1/version/import',
        method: 'post',
        data
    })
}

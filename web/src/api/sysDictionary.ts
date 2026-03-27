import service from '@/utils/request'

export interface SysDictionary {
    ID?: number
    name?: string
    type?: string
    status?: boolean
    desc?: string
    sysDictionaryDetails?: any[]
}

export interface SysDictionarySearch {
    page: number
    pageSize: number
    name?: string
    type?: string
    status?: boolean
    desc?: string
}

// @Tags SysDictionary
// @Summary 创建SysDictionary
export const createSysDictionary = (data: SysDictionary) => {
    return service({
        url: '/api/v1/dictionary/add',
        method: 'post',
        data
    })
}

// @Tags SysDictionary
// @Summary 删除SysDictionary
export const deleteSysDictionary = (data: { ID: number }) => {
    return service({
        url: '/api/v1/dictionary/delete',
        method: 'post',
        data
    })
}

// @Tags SysDictionary
// @Summary 更新SysDictionary
export const updateSysDictionary = (data: SysDictionary) => {
    return service({
        url: '/api/v1/dictionary/update',
        method: 'post',
        data
    })
}

// @Tags SysDictionary
// @Summary 用id查询SysDictionary
export const findSysDictionary = (params: { ID: number | string }) => {
    return service({
        url: '/api/v1/dictionary/get',
        method: 'get',
        params
    })
}

// @Tags SysDictionary
// @Summary 分页获取SysDictionary列表
export const getSysDictionaryList = (params: SysDictionarySearch) => {
    return service({
        url: '/api/v1/dictionary/list',
        method: 'get',
        params
    })
}

// @Tags SysDictionary
// @Summary 导出字典JSON（包含字典详情）
export const exportSysDictionary = (params?: { ID: number | string }) => {
    return service({
        url: '/api/v1/dictionary/export',
        method: 'get',
        params
    })
}

// @Tags SysDictionary
// @Summary 导入字典JSON（包含字典详情）
export const importSysDictionary = (data: any) => {
    return service({
        url: '/api/v1/dictionary/import',
        method: 'post',
        data
    })
}

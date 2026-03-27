import service from '@/utils/request'

export interface SysDictionaryDetail {
    ID?: number
    label?: string
    value?: number | string
    status?: boolean
    sort?: number
    sysDictionaryID?: number
}

export interface SysDictionaryDetailSearch {
    page: number
    pageSize: number
    label?: string
    value?: number | string
    status?: boolean
    sysDictionaryID?: number | string
}

// @Tags SysDictionaryDetail
// @Summary 创建SysDictionaryDetail
export const createSysDictionaryDetail = (data: SysDictionaryDetail) => {
    return service({
        url: '/api/v1/dictionary/detail/add',
        method: 'post',
        data
    })
}

// @Tags SysDictionaryDetail
// @Summary 删除SysDictionaryDetail
export const deleteSysDictionaryDetail = (data: { ID: number }) => {
    return service({
        url: '/api/v1/dictionary/detail/delete',
        method: 'post',
        data
    })
}

// @Tags SysDictionaryDetail
// @Summary 更新SysDictionaryDetail
export const updateSysDictionaryDetail = (data: SysDictionaryDetail) => {
    return service({
        url: '/api/v1/dictionary/detail/update',
        method: 'post',
        data
    })
}

// @Tags SysDictionaryDetail
// @Summary 用id查询SysDictionaryDetail
export const findSysDictionaryDetail = (params: { ID: number | string }) => {
    return service({
        url: '/api/v1/dictionary/detail/get',
        method: 'get',
        params
    })
}

// @Tags SysDictionaryDetail
// @Summary 分页获取SysDictionaryDetail列表
export const getSysDictionaryDetailList = (params: SysDictionaryDetailSearch) => {
    return service({
        url: '/api/v1/dictionary/detail/list',
        method: 'get',
        params
    })
}

// @Tags SysDictionaryDetail
// @Summary 获取层级字典详情树形结构（根据字典ID）
export const getDictionaryTreeList = (params: { sysDictionaryID: string | number }) => {
    return service({
        url: '/api/v1/dictionary/detail/tree/list',
        method: 'get',
        params
    })
}

// @Tags SysDictionaryDetail
// @Summary 获取层级字典详情树形结构（根据字典类型）
export const getDictionaryTreeListByType = (params: { dictType: string }) => {
    return service({
        url: '/api/v1/dictionary/detail/tree/type/list',
        method: 'get',
        params
    })
}

// @Tags SysDictionaryDetail
// @Summary 根据父级ID获取字典详情
export const getDictionaryDetailsByParent = (params: { parentID: string | number; includeChildren?: boolean }) => {
    return service({
        url: '/api/v1/dictionary/detail/parent/get',
        method: 'get',
        params
    })
}

// @Tags SysDictionaryDetail
// @Summary 获取字典详情的完整路径
export const getDictionaryPath = (params: { ID: string | number }) => {
    return service({
        url: '/api/v1/dictionary/detail/path/get',
        method: 'get',
        params
    })
}

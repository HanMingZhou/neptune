import service from '@/utils/request'

export interface SysOperationRecordSearch {
    page: number
    pageSize: number
    path?: string
    method?: string
    status?: number
}

// @Tags SysOperationRecord
// @Summary 删除SysOperationRecord
export const deleteSysOperationRecord = (data: { ID: number }) => {
    return service({
        url: '/api/v1/operation/record/delete',
        method: 'post',
        data
    })
}

// @Tags SysOperationRecord
// @Summary 删除SysOperationRecord
export const deleteSysOperationRecordByIds = (data: { ids: number[] }) => {
    return service({
        url: '/api/v1/operation/record/delete/multi',
        method: 'post',
        data
    })
}

// @Tags SysOperationRecord
// @Summary 分页获取SysOperationRecord列表
export const getSysOperationRecordList = (params: SysOperationRecordSearch) => {
    return service({
        url: '/api/v1/operation/record/list',
        method: 'get',
        params
    })
}

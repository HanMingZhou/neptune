import service from '@/utils/request'

export interface ExportTemplate {
    ID?: number
    name?: string
    tableName?: string
    templateID?: string
    // Add others
    [key: string]: any
}

export interface ExportTemplateSearch {
    page: number
    pageSize: number
    name?: string
    tableName?: string
    templateID?: string
}


// @Tags SysExportTemplate
// @Summary 创建导出模板
export const createSysExportTemplate = (data: ExportTemplate) => {
    return service({
        url: '/api/v1/export/template/add',
        method: 'post',
        data
    })
}

// @Tags SysExportTemplate
// @Summary 删除导出模板
export const deleteSysExportTemplate = (data: { ID: number }) => {
    return service({
        url: '/api/v1/export/template/delete',
        method: 'post',
        data
    })
}

// @Tags SysExportTemplate
// @Summary 批量删除导出模板
export const deleteSysExportTemplateByIds = (data: { ids: number[] }) => {
    return service({
        url: '/api/v1/export/template/delete/multi',
        method: 'post',
        data
    })
}

// @Tags SysExportTemplate
// @Summary 更新导出模板
export const updateSysExportTemplate = (data: ExportTemplate) => {
    return service({
        url: '/api/v1/export/template/update',
        method: 'post',
        data
    })
}

// @Tags SysExportTemplate
// @Summary 用id查询导出模板
export const findSysExportTemplate = (params: { ID: number }) => {
    return service({
        url: '/api/v1/export/template/get',
        method: 'get',
        params
    })
}

// @Tags SysExportTemplate
// @Summary 分页获取导出模板列表
export const getSysExportTemplateList = (params: ExportTemplateSearch) => {
    return service({
        url: '/api/v1/export/template/list',
        method: 'get',
        params
    })
}


// ExportExcel 导出表格token
// @Tags SysExportTemplate
// @Summary 导出表格
export const exportExcel = (params: { templateID: string;[key: string]: any }) => {
    return service({
        url: '/api/v1/export/template/export',
        method: 'get',
        params
    })
}

// ExportTemplate 导出表格模板
// @Tags SysExportTemplate
// @Summary 导出表格模板
export const exportTemplate = (params: { templateID: string }) => {
    return service({
        url: '/api/v1/export/template/download',
        method: 'get',
        params
    })
}

// PreviewSQL 预览最终生成的SQL
// @Tags SysExportTemplate
// @Summary 预览最终生成的SQL
export const previewSQL = (params: { templateID: string;[key: string]: any }) => {
    return service({
        url: '/api/v1/export/template/sql/preview',
        method: 'get',
        params
    })
}

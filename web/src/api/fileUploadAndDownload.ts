import service from '@/utils/request'

export interface FileListParams {
    page: number
    pageSize: number
    keyword?: string
}

// @Tags FileUploadAndDownload
// @Summary 分页文件列表
export const getFileList = (data: FileListParams) => {
    return service({
        url: '/api/v1/file/list',
        method: 'post',
        data
    })
}

// @Tags FileUploadAndDownload
// @Summary 删除文件
export const deleteFile = (data: { id: number; key: string }) => {
    return service({
        url: '/api/v1/file/delete',
        method: 'post',
        data
    })
}

/**
 * 编辑文件名或者备注
 * @param data
 * @returns {*}
 */
export const editFileName = (data: { id: number; name: string; tag?: string }) => {
    return service({
        url: '/api/v1/file/update',
        method: 'post',
        data
    })
}

/**
 * 导入URL
 * @param data
 * @returns {*}
 */
export const importURL = (data: { url: string; name?: string }) => {
    return service({
        url: '/api/v1/file/import',
        method: 'post',
        data
    })
}


// 上传文件 暂时用于头像上传
export const uploadFile = (data: FormData) => {
    return service({
        url: "/api/v1/file/upload",
        method: "post",
        data,
    });
};

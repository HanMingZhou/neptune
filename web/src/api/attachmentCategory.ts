import service from '@/utils/request'

export interface AttachmentCategory {
    ID: number
    categoryName: string
}

// 分类列表
export const getCategoryList = () => {
    return service({
        url: '/api/v1/attachment/category/list',
        method: 'get',
    })
}

// 添加/编辑分类
export const addCategory = (data: AttachmentCategory) => {
    return service({
        url: '/api/v1/attachment/category/add',
        method: 'post',
        data
    })
}

// 删除分类
export const deleteCategory = (data: { ID: number }) => {
    return service({
        url: '/api/v1/attachment/category/delete',
        method: 'post',
        data
    })
}

import service from '@/utils/request'

// 获取集群列表
export const getClusterList = (data: { keyword?: string; status?: number }) => {
    return service({
        url: '/api/v1/cms/cluster/list',
        method: 'post',
        data
    })
}

// 创建集群
export const createCluster = (data: any) => {
    return service({
        url: '/api/v1/cms/cluster/add',
        method: 'post',
        data
    })
}

// 更新集群
export const updateCluster = (data: any) => {
    return service({
        url: '/api/v1/cms/cluster/update',
        method: 'post',
        data
    })
}

// 删除集群
export const deleteCluster = (data: { id: number }) => {
    return service({
        url: '/api/v1/cms/cluster/delete',
        method: 'post',
        data
    })
}

import service from '@/utils/request'

// 获取集群节点列表（含资源信息）
export const getClusterNodes = (data: { clusterId: string }) => {
    return service({
        url: '/api/v1/cms/node/list',
        method: 'post',
        data
    })
}

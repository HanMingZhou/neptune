import service from '@/utils/request'
import type { ResourceId } from '@/types/consoleResource'
import type { CmsClusterForm } from '@/types/superAdmin'

type ClusterPayload = Partial<CmsClusterForm> &
  Pick<CmsClusterForm, 'name' | 'status'>

// 获取集群列表
export const getClusterList = (data: {
  keyword?: string
  page?: number
  pageSize?: number
  status?: number
}) => {
  return service({
    url: '/api/v1/cms/cluster/list',
    method: 'post',
    data
  })
}

// 创建集群
export const createCluster = (data: ClusterPayload) => {
  return service({
    url: '/api/v1/cms/cluster/add',
    method: 'post',
    data
  })
}

// 更新集群
export const updateCluster = (data: ClusterPayload) => {
  return service({
    url: '/api/v1/cms/cluster/update',
    method: 'post',
    data
  })
}

// 删除集群
export const deleteCluster = (data: { id: ResourceId }) => {
  return service({
    url: '/api/v1/cms/cluster/delete',
    method: 'post',
    data
  })
}

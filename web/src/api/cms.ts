import service from '@/utils/request'
import type { ResourceId } from '@/types/consoleResource'
import type {
  CmsBatchCreateComputeProductPayload,
  CmsNodeRow,
  CmsProductNodeCandidatesParams,
  CmsProductCatalogParams,
  CmsProductForm,
  CmsProductPricePayload
} from '@/types/superAdmin'

export interface CMSProduct {
  id?: number
  name?: string
  // Add other properties as needed
  [key: string]: any
}

// 获取CMS产品列表
export const getCMSProductList = (params: CmsProductCatalogParams) => {
  return service({
    url: '/api/v1/cms/product/list',
    method: 'get',
    params
  })
}

// 获取产品详情
export const getCMSProductDetail = (id: ResourceId) => {
  return service({
    url: '/api/v1/cms/product/detail',
    method: 'get',
    params: { id }
  })
}

// 创建产品
export const createCMSProduct = (data: CmsProductForm) => {
  return service({
    url: '/api/v1/cms/product/add',
    method: 'post',
    data
  })
}

// 更新产品
export const updateCMSProduct = (data: CmsProductForm) => {
  return service({
    url: '/api/v1/cms/product/update',
    method: 'post',
    data
  })
}

// 更新价格
export const updateCMSProductPrice = (data: CmsProductPricePayload) => {
  return service({
    url: '/api/v1/cms/product/price/update',
    method: 'post',
    data
  })
}

// 删除产品
export const deleteCMSProduct = (data: { id: ResourceId }) => {
  return service({
    url: '/api/v1/cms/product/delete',
    method: 'post',
    data
  })
}

// 获取集群列表
export const getCMSClusterList = () => {
  return service({
    url: '/api/v1/cms/product/cluster/list',
    method: 'get'
  })
}

// 获取地区列表
export const getCMSAreaList = () => {
  return service({
    url: '/api/v1/cms/product/area/list',
    method: 'get'
  })
}

// 获取集群下的K8s节点列表（含资源信息）
export const getCMSClusterNodes = (data: {
  clusterId: ResourceId
  cpu?: number
  memory?: number
  gpuCount?: number
}) => {
  return service({
    url: '/api/v1/cms/product/node/list',
    method: 'get',
    params: data
  })
}

export const getCMSProductNodeCandidates = (
  params: CmsProductNodeCandidatesParams
) => {
  return service({
    url: '/api/v1/cms/product/node/list',
    method: 'get',
    params
  })
}

export const createCMSComputeProductsBatch = (
  data: CmsBatchCreateComputeProductPayload
) => {
  return service({
    url: '/api/v1/cms/product/batch/add',
    method: 'post',
    data
  })
}

// 节点管理 - 获取列表
export const getCMSNodeList = (data: {
  clusterId: ResourceId
  keyword?: string
  page?: number
  pageSize?: number
}) => {
  return service({
    url: '/api/v1/cms/node/list',
    method: 'post',
    data
  })
}

// 节点管理 - 恢复调度
export const uncordonCMSNode = (data: {
  clusterId: ResourceId
  nodeName: CmsNodeRow['nodeName']
}) => {
  return service({
    url: '/api/v1/cms/node/uncordon',
    method: 'post',
    data
  })
}

// 节点管理 - 驱逐节点
export const drainCMSNode = (data: {
  clusterId: ResourceId
  nodeName: CmsNodeRow['nodeName']
}) => {
  return service({
    url: '/api/v1/cms/node/drain',
    method: 'post',
    data
  })
}

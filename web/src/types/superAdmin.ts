import type { ResourceId } from './consoleResource'

export interface LabelValueOption {
  label: string
  value: string
}

export type MethodTone = '' | 'success' | 'warning' | 'danger'

export interface ApiListItem {
  ID: ResourceId
  path: string
  apiGroup: string
  method: string
  description: string
  [key: string]: unknown
}

export interface ApiForm {
  path: string
  apiGroup: string
  method: string
  description: string
}

export interface ApiSyncData {
  newApis: ApiListItem[]
  deleteApis: ApiListItem[]
  ignoreApis: ApiListItem[]
}

export interface ApiGroupListData {
  groups?: string[]
  apiGroupMap?: Record<string, string>
}

export interface ApiListData {
  list?: ApiListItem[]
  total?: number
  page?: number
  pageSize?: number
}

export interface ApiMethodOption extends LabelValueOption {
  type: MethodTone
}

export interface AuthorityTreeNode {
  authorityId: ResourceId
  authorityName: string
  parentId: ResourceId
  dataAuthorityId?: ResourceId[]
  children?: AuthorityTreeNode[]
  menus?: unknown[]
  [key: string]: unknown
}

export interface AuthorityListData {
  list?: AuthorityTreeNode[]
  total?: number
  page?: number
  pageSize?: number
}

export interface AuthorityForm {
  authorityId: ResourceId | ''
  authorityName: string
  parentId: ResourceId
}

export interface AuthorityOption {
  authorityId: ResourceId
  authorityName: string
  disabled?: boolean
  children?: AuthorityOption[]
}

export type AuthorityDialogType = 'add' | 'edit' | 'copy'

export interface AuthorityCopyPayload {
  authority: {
    authorityId: number
    authorityName: string
    dataAuthorityId: ResourceId[]
    parentId: number
  }
  oldAuthorityId: number
}

export interface MenuMetaForm {
  activeName: string
  title: string
  icon: string
  defaultMenu: boolean
  closeTab: boolean
  keepAlive: boolean
  transitionType: string
}

export interface MenuParameter {
  ID?: number
  type: string
  key: string
  value: string
  [key: string]: unknown
}

export interface MenuButton {
  ID?: number
  name: string
  desc: string
  [key: string]: unknown
}

export interface MenuTreeNode {
  ID: number
  path?: string
  name?: string
  hidden?: boolean
  parentId?: ResourceId
  component?: string
  sort?: number
  meta?: Partial<MenuMetaForm>
  children?: MenuTreeNode[]
  parameters?: MenuParameter[]
  menuBtn?: MenuButton[]
  [key: string]: unknown
}

export interface MenuForm {
  ID: number
  path: string
  name: string
  hidden: boolean
  parentId: ResourceId
  component: string
  sort?: number
  meta: MenuMetaForm
  parameters: MenuParameter[]
  menuBtn: MenuButton[]
}

export interface MenuOption {
  title?: string
  ID: number
  disabled?: boolean
  children?: MenuOption[]
}

export interface MenuListData {
  list?: MenuTreeNode[]
  total?: number
  page?: number
  pageSize?: number
}

export interface OperationRecordUser {
  userName?: string
  nickName?: string
  [key: string]: unknown
}

export interface OperationRecordItem {
  ID: number
  status?: number | string
  ip?: string
  method?: string
  path?: string
  body?: string
  resp?: string
  CreatedAt?: string | number
  user?: OperationRecordUser
  [key: string]: unknown
}

export interface OperationRecordListData {
  list?: OperationRecordItem[]
  total?: number
  page?: number
  pageSize?: number
}

export interface OperationRecordSearchInfo {
  method: string
  path: string
  status: string
}

export interface UserAuthority {
  authorityId: ResourceId
  authorityName?: string
  children?: UserAuthority[]
  disabled?: boolean
  [key: string]: unknown
}

export interface UserSearchInfo {
  username: string
  nickName: string
  phone: string
  email: string
}

export interface UserForm {
  ID: number
  id?: number
  userName: string
  password: string
  nickName: string
  headerImg: string
  authorityId: ResourceId | ''
  authorityIds: ResourceId[]
  authorities?: UserAuthority[]
  enable: number
  phone: string
  email: string
  [key: string]: unknown
}

export interface UserRow extends UserForm {
  _authorityDirty: boolean
}

export interface UserResetPasswordForm {
  ID: ResourceId | ''
  userName: string
  nickName: string
  password: string
}

export interface UserListData {
  list?: UserForm[]
  total?: number
  page?: number
  pageSize?: number
}

export interface MenuIconOption {
  key: string
  label: string
}

export interface CmsClusterRow {
  id: number
  name: string
  area?: string
  description?: string
  kubeConfig?: string
  apiServer?: string
  status?: number
  harborAddr?: string
  storageClass?: string
  internalIp?: string
  nodeCount?: number
  createdAt?: string
  [key: string]: unknown
}

export interface CmsClusterForm {
  id: number | null
  name: string
  area: string
  description: string
  kubeconfig: string
  apiServer: string
  status: number
  harborAddr: string
  storageClass: string
}

export interface CmsClusterListData {
  list?: CmsClusterRow[]
  total?: number
  page?: number
  pageSize?: number
}

export interface CmsClusterOption {
  id: ResourceId
  name: string
  area?: string
  harborAddr?: string
  [key: string]: unknown
}

export interface CmsNodeRow {
  nodeName: string
  internalIp?: string
  clusterName?: string
  nodeRole?: string
  area?: string
  schedulable?: boolean
  cpu?: number
  memory?: number
  cpuModel?: string
  cpuAvailable?: number
  cpuAllocatable?: number
  memoryAvailable?: number
  memoryAllocatable?: number
  gpuModel?: string
  gpuCount?: number
  gpuAvailable?: number
  gpuMemory?: number
  driverVersion?: string
  cudaVersion?: string
  vGpuNumber?: number
  vGpuCount?: number
  vGpuMemory?: number
  vGpuCores?: number
  [key: string]: unknown
}

export interface CmsExistingComputeProductSummary {
  id: number
  name: string
  description?: string
  resourceType: CmsProductResourceType
  cpu: number
  memory: number
  gpuModel?: string
  gpuCount: number
  gpuMemory: number
  vGpuNumber?: number
  vGpuMemory: number
  vGpuCores: number
  priceHourly: number
  priceDaily: number
  priceWeekly: number
  priceMonthly: number
  status: number
  maxInstances: number
  usedCapacity: number
  available: number
}

export interface CmsProductNodeCandidate extends CmsNodeRow {
  existingComputeProducts?: CmsExistingComputeProductSummary[]
  canCreateComputeProduct?: boolean
  compatible?: boolean
  disableReason?: string
}

export interface CmsNodeListData {
  nodes?: CmsNodeRow[]
  total?: number
  page?: number
  pageSize?: number
}

export type CmsProductType = 1 | 2
export type CmsProductResourceType = 'cpu' | 'gpu' | 'vgpu'
export type CmsProductFilterResourceType = CmsProductResourceType | ''
export type CmsComputePriceType = 1 | 2 | 3 | 4
export type CmsCatalogPriceType = CmsComputePriceType | 5
export type CmsProductPriceField = CmsCatalogPriceType

export interface CmsProductPriceItem {
  priceType: CmsCatalogPriceType
  price: number
}

export interface CmsProductRow {
  id: number | null
  productType: CmsProductType
  name: string
  description: string
  clusterId: ResourceId | null
  clusterName?: string
  area: string
  nodeName: string
  nodeIp?: string
  nodeType: string
  cpuModel: string
  cpu: number
  memory: number
  gpuModel: string
  gpuCount: number
  gpuMemory: number
  vGpuCount: number
  vGpuNumber?: number
  vGpuMemory: number
  vGpuCores: number
  priceHourly: number
  priceDaily: number
  priceWeekly: number
  priceMonthly: number
  driverVersion?: string
  cudaVersion?: string
  status: number
  sortOrder?: number
  available?: number
  maxInstances: number
  usedCapacity?: number
  version?: number
  prices?: CmsProductPriceItem[]
  storageClass: string
  storagePriceGb: number
  [key: string]: unknown
}

export type CmsProductForm = CmsProductRow

export interface CmsProductPriceForm {
  id: number | null
  priceHourly: number
  priceDaily: number
  priceWeekly: number
  priceMonthly: number
}

export interface CmsProductPricePayload {
  id: number | null
  prices: CmsProductPriceItem[]
}

export interface CmsProductListData {
  list?: CmsProductRow[]
  total?: number
}

export interface CmsProductCatalogParams {
  page: number
  pageSize: number
  productType: CmsProductType
  clusterId?: ResourceId
  area?: string
  resourceType?: CmsProductResourceType
  gpuModel?: string
  availableMin?: number
  availableMax?: number
  maxInstancesMin?: number
  maxInstancesMax?: number
  usedCapacityMin?: number
  usedCapacityMax?: number
  priceType?: CmsCatalogPriceType
  priceMin?: number
  priceMax?: number
  keyword?: string
}

export interface CmsNodeSelectionState {
  resourceType: CmsProductResourceType
  fields: Pick<
    CmsProductForm,
    | 'nodeName'
    | 'area'
    | 'cpu'
    | 'memory'
    | 'cpuModel'
    | 'gpuModel'
    | 'gpuCount'
    | 'gpuMemory'
    | 'driverVersion'
    | 'cudaVersion'
    | 'vGpuNumber'
    | 'vGpuCount'
    | 'vGpuMemory'
    | 'vGpuCores'
  >
  suggestedName: string
}

export interface CmsProductNodeCandidatesParams {
  clusterId: ResourceId
  resourceType?: CmsProductResourceType
  cpu?: number
  memory?: number
  gpuCount?: number
  gpuMemory?: number
  vGpuNumber?: number
  vGpuMemory?: number
  vGpuCores?: number
  excludeProductId?: ResourceId
}

export interface CmsBatchCreateComputeProductPayload {
  productType: CmsProductType
  nodeNames: string[]
  name: string
  description: string
  clusterId: ResourceId
  area: string
  nodeType: string
  cpuModel: string
  cpu: number
  memory: number
  gpuModel: string
  gpuCount: number
  gpuMemory: number
  vGpuNumber: number
  vGpuMemory: number
  vGpuCores: number
  prices: CmsProductPriceItem[]
  driverVersion: string
  cudaVersion: string
  systemDisk?: number
  dataDisk?: number
  status: number
  sortOrder?: number
  maxInstances?: number
}

export interface CmsBatchCreateComputeProductResult {
  createdIds?: number[]
  createdCount?: number
  createdNodes?: string[]
  skippedNodes?: string[]
}

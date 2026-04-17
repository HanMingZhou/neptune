import type {
  ConsoleProduct,
  ConsoleVolume,
  PageListData,
  ResourceId
} from './consoleResource'

export interface StorageUsageReference {
  type?: string
  name?: string
  [key: string]: unknown
}

export interface StorageListItem extends Omit<ConsoleVolume, 'size'> {
  id: ResourceId
  name: string
  pvcName?: string
  size?: string | number
  requestedSize?: string | number
  resizePending?: boolean
  productName?: string
  status?: string
  area?: string
  createdAt?: string | number
  usedBy?: StorageUsageReference[]
}

export interface StorageClusterOption {
  id: ResourceId
  name?: string
  area?: string
  harborAddr?: string
  [key: string]: unknown
}

export interface StorageCreateForm {
  name: string
  size: number
  area: string
  clusterId: ResourceId | ''
  productId: ResourceId | ''
}

export interface StorageExpandForm {
  id: ResourceId
  currentSize: string
  minSize: number
  newSize: number
}

export interface StorageEditForm {
  id: ResourceId
  name: string
}

export interface StorageAreaListData {
  clusters?: StorageClusterOption[]
}

export interface StorageProductOption extends ConsoleProduct {
  name?: string
}

export type StorageListData = PageListData<StorageListItem>
export type StorageProductListData = PageListData<StorageProductOption>

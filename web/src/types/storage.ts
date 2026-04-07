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
  size?: string | number
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

export interface StorageAreaListData {
  clusters?: StorageClusterOption[]
}

export interface StorageProductOption extends ConsoleProduct {
  name?: string
}

export type StorageListData = PageListData<StorageListItem>
export type StorageProductListData = PageListData<StorageProductOption>

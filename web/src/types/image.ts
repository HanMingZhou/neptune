import type { ConsoleImage, PageListData, ResourceId } from './consoleResource'

export type ImageFilterValue = '' | number

export interface ImageListItem extends ConsoleImage {
  type?: number
  usageType?: number
  image?: string
  area?: string
  size?: string
  imagePath?: string
  createTime?: string
}

export interface ImageForm {
  id: ResourceId | null
  name: string
  type: number
  usageType: number
  imageAddr: string
  area: string
  size: string
  imagePath: string
}

export interface ImageMutationPayload extends Omit<ImageForm, 'id'> {
  id?: ResourceId | null
}

export type ImageListData = PageListData<ImageListItem>

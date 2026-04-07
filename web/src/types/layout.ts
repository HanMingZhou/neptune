export type LayoutGroupKey = 'admin' | 'compute' | 'management' | 'resources'

export interface LayoutNavItem {
  key: string
  titleKey: string
  title: string
  icon: string
  routeName?: string
  children?: LayoutNavItem[]
}

export interface LayoutNavGroup {
  title: string
  items: LayoutNavItem[]
}

export interface LayoutNavGroupBucket extends LayoutNavGroup {
  sort: number
}

export type LayoutTranslationValue = string | LayoutTranslationDictionary

export interface LayoutTranslationDictionary {
  [key: string]: LayoutTranslationValue
}

export type LayoutTranslationRegistry = Record<
  string,
  LayoutTranslationDictionary
>

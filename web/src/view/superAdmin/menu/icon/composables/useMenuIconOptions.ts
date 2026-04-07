import { computed } from 'vue'
import config from '@/core/config'
import { elementPlusIconNames } from '@/generated/elementPlusIconNames'
import type { MenuIconOption } from '@/types/superAdmin'

const createElementIconOptions = (): MenuIconOption[] =>
  elementPlusIconNames
    .map((iconName) => ({
      key: iconName,
      label: iconName
    }))
    .sort((left, right) => left.key.localeCompare(right.key))

const baseIconOptions = createElementIconOptions()

export function useMenuIconOptions() {
  const iconOptions = computed<MenuIconOption[]>(() => {
    const mergedOptions = [...baseIconOptions, ...(config.logs || [])]
    const optionMap = new Map<string, MenuIconOption>()

    mergedOptions.forEach((option) => {
      if (!option?.key || optionMap.has(option.key)) {
        return
      }

      optionMap.set(option.key, {
        key: option.key,
        label: option.label || option.key
      })
    })

    return Array.from(optionMap.values())
  })

  return {
    iconOptions
  }
}

import { computed } from 'vue'
import * as ElIconModules from '@element-plus/icons-vue'
import config from '@/core/config'

const toKebabCase = (value = '') =>
  value
    .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .toLowerCase()

const createElementIconOptions = () =>
  Object.keys(ElIconModules)
    .map((iconName) => {
      const normalizedName = toKebabCase(iconName)
      return {
        key: normalizedName,
        label: normalizedName
      }
    })
    .sort((left, right) => left.key.localeCompare(right.key))

const baseIconOptions = createElementIconOptions()

export function useMenuIconOptions() {
  const iconOptions = computed(() => {
    const mergedOptions = [...baseIconOptions, ...(config.logs || [])]
    const optionMap = new Map()

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

<template>
  <span
    v-if="name"
    class="inline-flex items-center justify-center leading-none"
  >
    <el-icon v-if="resolvedIconComponent">
      <component :is="resolvedIconComponent" />
    </el-icon>
    <span
      v-else-if="shouldRenderMaterialIcon"
      class="material-icons leading-none"
    >
      {{ name }}
    </span>
  </span>
</template>

<script setup lang="ts">
import { computed, getCurrentInstance, ref, watch, type Component } from 'vue'

type ResolvedIconComponent = string | Component
type ElementPlusIconModule = Record<string, Component>

const props = withDefaults(
  defineProps<{
    name?: string
  }>(),
  {
    name: ''
  }
)

const instance = getCurrentInstance()
const resolvedIconComponent = ref<ResolvedIconComponent | null>(null)
const resolvingElementIcon = ref(false)

let elementPlusIconModulesPromise: Promise<ElementPlusIconModule> | null = null
const elementPlusIconCache = new Map<string, Component>()

const toPascalCase = (value = '') =>
  value
    .split('-')
    .filter(Boolean)
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join('')

const getIconCandidates = (value = ''): string[] => {
  const trimmedValue = value.trim()

  if (!trimmedValue) {
    return []
  }

  return Array.from(new Set([trimmedValue, toPascalCase(trimmedValue)]))
}

const resolveRegisteredComponent = (
  iconName: string
): ResolvedIconComponent | null => {
  const registeredComponents = instance?.appContext.components || {}

  for (const candidate of getIconCandidates(iconName)) {
    if (registeredComponents[candidate]) {
      return candidate
    }
  }

  return null
}

const loadElementPlusIcon = async (
  iconName: string
): Promise<Component | null> => {
  const candidates = getIconCandidates(iconName)

  for (const candidate of candidates) {
    const cachedIcon = elementPlusIconCache.get(candidate)
    if (cachedIcon) {
      return cachedIcon
    }
  }

  if (!elementPlusIconModulesPromise) {
    elementPlusIconModulesPromise =
      import('@element-plus/icons-vue') as Promise<ElementPlusIconModule>
  }

  const elementPlusIcons = await elementPlusIconModulesPromise

  for (const candidate of candidates) {
    const iconComponent = elementPlusIcons[candidate]
    if (iconComponent) {
      elementPlusIconCache.set(candidate, iconComponent)
      return iconComponent
    }
  }

  return null
}

watch(
  () => props.name,
  async (iconName) => {
    const normalizedName = iconName?.trim() || ''

    if (!normalizedName) {
      resolvedIconComponent.value = null
      resolvingElementIcon.value = false
      return
    }

    const registeredComponent = resolveRegisteredComponent(normalizedName)
    if (registeredComponent) {
      resolvedIconComponent.value = registeredComponent
      resolvingElementIcon.value = false
      return
    }

    resolvingElementIcon.value = true

    try {
      const elementPlusIcon = await loadElementPlusIcon(normalizedName)

      if ((props.name?.trim() || '') !== normalizedName) {
        return
      }

      resolvedIconComponent.value = elementPlusIcon
    } catch (error) {
      if ((props.name?.trim() || '') !== normalizedName) {
        return
      }

      resolvedIconComponent.value = null
      console.error(`Failed to resolve icon "${normalizedName}"`, error)
    } finally {
      if ((props.name?.trim() || '') === normalizedName) {
        resolvingElementIcon.value = false
      }
    }
  },
  { immediate: true }
)

const shouldRenderMaterialIcon = computed(() => {
  const iconName = props.name?.trim()
  return (
    Boolean(iconName) &&
    !resolvedIconComponent.value &&
    !resolvingElementIcon.value
  )
})
</script>

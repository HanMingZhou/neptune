<template>
  <span
    v-if="name"
    class="inline-flex items-center justify-center leading-none"
  >
    <el-icon v-if="resolvedIconComponent">
      <component :is="resolvedIconComponent" />
    </el-icon>
    <span
      v-else
      class="material-icons leading-none"
    >
      {{ name }}
    </span>
  </span>
</template>

<script setup>
import { computed, getCurrentInstance } from 'vue'

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
})

const instance = getCurrentInstance()

const toPascalCase = (value = '') =>
  value
    .split('-')
    .filter(Boolean)
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join('')

const resolvedIconComponent = computed(() => {
  const iconName = props.name?.trim()
  if (!iconName) {
    return null
  }

  const registeredComponents = instance?.appContext.components || {}
  const candidates = [
    iconName,
    toPascalCase(iconName)
  ]

  for (const candidate of candidates) {
    if (registeredComponents[candidate]) {
      return candidate
    }
  }

  return null
})
</script>

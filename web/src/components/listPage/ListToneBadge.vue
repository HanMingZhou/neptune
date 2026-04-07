<template>
  <span :class="resolvedClass" class="console-badge">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ListTone = 'neutral' | 'info' | 'success' | 'warning' | 'danger'

const props = withDefaults(
  defineProps<{
    label?: string
    tone?: ListTone
  }>(),
  {
    label: '',
    tone: 'neutral'
  }
)

const toneClassMap: Record<ListTone, string> = {
  neutral: 'console-badge--neutral',
  info: 'console-badge--info',
  success: 'console-badge--success',
  warning: 'console-badge--warning',
  danger: 'console-badge--danger'
}

const resolvedClass = computed(
  () => toneClassMap[props.tone] || toneClassMap.neutral
)
</script>

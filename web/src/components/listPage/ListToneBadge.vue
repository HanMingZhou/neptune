<template>
  <span :class="resolvedClass" class="console-badge">
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  tone: {
    type: String,
    default: 'neutral'
  }
})

const toneClassMap = {
  neutral: 'console-badge--neutral',
  info: 'console-badge--info',
  success: 'console-badge--success',
  warning: 'console-badge--warning',
  danger: 'console-badge--danger'
}

const resolvedClass = computed(() => toneClassMap[props.tone] || toneClassMap.neutral)
</script>

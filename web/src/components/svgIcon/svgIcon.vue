<template>
  <svg :class="svgClass" aria-hidden="true">
    <use :xlink:href="iconName" />
  </svg>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  iconClass: {
    type: String,
    default: ''
  },
  localIcon: {
    type: String,
    default: ''
  },
  className: {
    type: String,
    default: ''
  }
})

const iconName = computed(() => {
  if (props.localIcon) {
    return `#icon-${props.localIcon}`
  }
  return `#icon-${props.iconClass}`
})
const svgClass = computed(() => {
  if (props.className) {
    return `svg-icon ${props.className}`
  }
  return 'svg-icon'
})
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}
</style>

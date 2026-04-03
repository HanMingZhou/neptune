<template>
  <div class="console-resource-status">
    <span :class="statusClass" class="console-badge text-[11px] flex items-center gap-1.5 border border-transparent">
      <span v-if="showPulse" :class="pulseClass" class="size-1.5 rounded-full"></span>
      {{ label }}
    </span>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },
  pulseClass: {
    type: String,
    default: ''
  },
  showPulse: {
    type: Boolean,
    default: false
  },
  statusClass: {
    type: [String, Array, Object],
    default: ''
  }
})
</script>

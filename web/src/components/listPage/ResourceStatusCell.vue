<template>
  <div class="console-resource-status">
    <span
      :class="statusClass"
      class="console-badge text-[11px] flex items-center gap-1.5 border border-transparent"
    >
      <span
        v-if="showPulse"
        :class="pulseClass"
        class="size-1.5 rounded-full"
      ></span>
      {{ label }}
    </span>
  </div>
</template>

<script setup lang="ts">
type StatusClass = string | string[] | Record<string, boolean>

withDefaults(
  defineProps<{
    label?: string
    pulseClass?: string
    showPulse?: boolean
    statusClass?: StatusClass
  }>(),
  {
    label: '',
    pulseClass: '',
    showPulse: false,
    statusClass: ''
  }
)
</script>

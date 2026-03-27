<template>
  <LogStreamPanel
    :logs="logs"
    :logs-connected="logsConnected"
    :logs-loading="logsLoading"
    :set-logs-ref="setLogsRef"
    @clear="$emit('clear')"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
  />
</template>

<script setup>
import LogStreamPanel from '@/components/detailPage/LogStreamPanel.vue'

defineProps({
  logs: {
    type: String,
    default: ''
  },
  logsConnected: {
    type: Boolean,
    default: false
  },
  logsLoading: {
    type: Boolean,
    default: false
  },
  setLogsRef: {
    type: Function,
    required: true
  }
})

defineEmits(['clear', 'connect', 'disconnect'])
</script>

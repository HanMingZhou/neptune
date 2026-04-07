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

<script setup lang="ts">
import LogStreamPanel from '@/components/detailPage/LogStreamPanel.vue'

withDefaults(
  defineProps<{
    logs?: string
    logsConnected?: boolean
    logsLoading?: boolean
    setLogsRef: (element: HTMLElement | null) => void
  }>(),
  {
    logs: '',
    logsConnected: false,
    logsLoading: false
  }
)

defineEmits<{
  clear: []
  connect: []
  disconnect: []
}>()
</script>

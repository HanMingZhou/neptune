<template>
  <TerminalShell
    :body-style="{ overscrollBehavior: 'contain' }"
    :can-connect="isInstanceRunning"
    :connected-label="'Gateway: 22'"
    :disabled-status-text="notebook.status !== 'Running' ? `${t('instanceStatus')}: ${notebook.status}` : ''"
    :set-terminal-ref="setTerminalRef"
    :terminal-connected="terminalConnected"
    :terminal-title="`root@${notebook.instanceName || ''}-0:~`"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
    @fit="$emit('fit')"
  />
</template>

<script setup>
import { inject } from 'vue'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

defineProps({
  isInstanceRunning: {
    type: Boolean,
    default: false
  },
  notebook: {
    type: Object,
    default: () => ({})
  },
  setTerminalRef: {
    type: Function,
    required: true
  },
  terminalConnected: {
    type: Boolean,
    default: false
  }
})

defineEmits(['connect', 'disconnect', 'fit'])

const t = inject('t', (key) => key)
</script>

<template>
  <TerminalShell
    :body-style="{ overscrollBehavior: 'contain' }"
    :can-connect="isInstanceRunning"
    :connected-label="'Gateway: 22'"
    :disabled-status-text="
      !isInstanceRunning
        ? `${t('instanceStatus')}: ${notebook.status || '-'}`
        : ''
    "
    :set-terminal-ref="setTerminalRef"
    :terminal-connected="terminalConnected"
    :terminal-title="`root@${notebook.instanceName || ''}-0:~`"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
    @fit="$emit('fit')"
  />
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ConsoleNotebook, Translator } from '@/types/consoleResource'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

withDefaults(
  defineProps<{
    isInstanceRunning?: boolean
    notebook?: ConsoleNotebook
    setTerminalRef: (element: HTMLElement | null) => void
    terminalConnected?: boolean
  }>(),
  {
    isInstanceRunning: false,
    notebook: () => ({}),
    terminalConnected: false
  }
)

defineEmits<{
  connect: []
  disconnect: []
  fit: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

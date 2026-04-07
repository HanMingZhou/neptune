<template>
  <TerminalShell
    :body-style="{ overscrollBehavior: 'contain' }"
    :can-connect="canConnectTerminal"
    :connected-label="t('connected')"
    :disabled-status-text="
      !canConnectTerminal ? t('inference.serviceNotReady') : ''
    "
    :set-terminal-ref="setTerminalRef"
    :shortcuts-class="'flex items-center gap-4 text-xs text-slate-400'"
    :terminal-connected="terminalConnected"
    :terminal-title="`${t('inference.terminalTitle')} - ${service.instanceName || ''}`"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
    @fit="$emit('fit')"
  >
    <template #controls-prefix>
      <select
        v-if="pods.length > 1"
        :value="terminalPod"
        class="px-3 py-1.5 border border-border-light dark:border-border-dark rounded-lg text-sm bg-white dark:bg-zinc-800"
        @change="$emit('update:terminal-pod', $event.target.value)"
      >
        <option v-for="pod in pods" :key="pod.name" :value="pod.name">
          {{ pod.name }}
        </option>
      </select>
    </template>
  </TerminalShell>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type {
  ConsoleInferenceService,
  ConsolePod,
  Translator
} from '@/types/consoleResource'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

withDefaults(
  defineProps<{
    canConnectTerminal?: boolean
    pods?: ConsolePod[]
    service?: ConsoleInferenceService
    setTerminalRef: (element: HTMLElement | null) => void
    terminalConnected?: boolean
    terminalPod?: string
  }>(),
  {
    canConnectTerminal: false,
    pods: () => [],
    service: () => ({ id: '' }),
    terminalConnected: false,
    terminalPod: ''
  }
)

defineEmits<{
  connect: []
  disconnect: []
  fit: []
  'update:terminal-pod': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

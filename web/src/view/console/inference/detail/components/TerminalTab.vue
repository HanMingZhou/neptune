<template>
  <TerminalShell
    :body-style="{ overscrollBehavior: 'contain' }"
    :can-connect="canConnectTerminal"
    :connected-label="t('connected')"
    :disabled-status-text="!canConnectTerminal ? t('inference.serviceNotReady') : ''"
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
          <option v-for="pod in pods" :key="pod.name" :value="pod.name">{{ pod.name }}</option>
        </select>
    </template>
  </TerminalShell>
</template>

<script setup>
import { inject } from 'vue'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

defineProps({
  canConnectTerminal: {
    type: Boolean,
    default: false
  },
  pods: {
    type: Array,
    default: () => []
  },
  service: {
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
  },
  terminalPod: {
    type: String,
    default: ''
  }
})

defineEmits(['connect', 'disconnect', 'fit', 'update:terminal-pod'])

const t = inject('t', (key) => key)
</script>

<template>
  <LogStreamPanel
    :logs="logs"
    :logs-connected="logsConnected"
    :logs-loading="logsLoading"
    :set-logs-ref="setLogsRef"
    @clear="$emit('clear')"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
  >
    <template #controls-prefix>
        <select
          v-if="pods.length > 1"
          :disabled="logsConnected"
          :value="selectedPod"
          class="px-3 py-1.5 border border-border-light dark:border-border-dark rounded-lg text-sm bg-white dark:bg-zinc-800"
          @change="$emit('update:selected-pod', $event.target.value)"
        >
          <option v-for="pod in pods" :key="pod.name" :value="pod.name">{{ pod.name }}</option>
        </select>
    </template>
    <template #content>
        <pre v-if="logs" class="text-slate-300 text-sm font-mono whitespace-pre-wrap break-all leading-6">{{ logs }}</pre>
        <div v-else class="text-slate-500 text-sm text-center py-8">{{ t('connectLogStream') }}</div>
    </template>
  </LogStreamPanel>
</template>

<script setup>
import { inject } from 'vue'
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
  pods: {
    type: Array,
    default: () => []
  },
  selectedPod: {
    type: String,
    default: ''
  },
  setLogsRef: {
    type: Function,
    required: true
  }
})

defineEmits(['clear', 'connect', 'disconnect', 'update:selected-pod'])

const t = inject('t', (key) => key)
</script>

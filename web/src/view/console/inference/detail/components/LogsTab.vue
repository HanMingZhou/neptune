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
        <option v-for="pod in pods" :key="pod.name" :value="pod.name">
          {{ pod.name }}
        </option>
      </select>
    </template>
    <template #content>
      <pre
        v-if="logs"
        class="text-slate-300 text-sm font-mono whitespace-pre-wrap break-all leading-6"
        >{{ logs }}</pre
      >
      <div v-else class="text-slate-500 text-sm text-center py-8">
        {{ t('connectLogStream') }}
      </div>
    </template>
  </LogStreamPanel>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ConsolePod, Translator } from '@/types/consoleResource'
import LogStreamPanel from '@/components/detailPage/LogStreamPanel.vue'

withDefaults(
  defineProps<{
    logs?: string
    logsConnected?: boolean
    logsLoading?: boolean
    pods?: ConsolePod[]
    selectedPod?: string
    setLogsRef: (element: HTMLElement | null) => void
  }>(),
  {
    logs: '',
    logsConnected: false,
    logsLoading: false,
    pods: () => [],
    selectedPod: ''
  }
)

defineEmits<{
  clear: []
  connect: []
  disconnect: []
  'update:selected-pod': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

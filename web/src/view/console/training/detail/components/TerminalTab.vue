<template>
  <TerminalShell
    :can-connect="isRunning"
    :connected-label="t('connected')"
    :disabled-status-text="
      !isRunning ? `${t('instanceStatus')}: ${statusLabel}` : ''
    "
    :set-terminal-ref="setTerminalRef"
    :shortcuts-class="'hidden md:flex items-center gap-4 text-slate-400 text-xs'"
    :terminal-connected="terminalConnected"
    :terminal-title="`root@${jobName || ''}-${terminalTaskName}:~`"
    @connect="$emit('connect')"
    @disconnect="$emit('disconnect')"
    @fit="$emit('fit')"
  >
    <template #controls-prefix>
      <el-select
        :model-value="terminalTaskName"
        :placeholder="t('pleaseSelect')"
        style="width: 150px"
        @update:model-value="$emit('update:terminal-task-name', $event)"
      >
        <el-option
          v-if="frameworkType !== 'STANDALONE'"
          :label="masterTaskLabel"
          :value="masterTaskName"
        />
        <el-option
          v-if="workerCount > 0 || frameworkType === 'STANDALONE'"
          label="Worker"
          value="worker"
        />
      </el-select>
    </template>
  </TerminalShell>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

withDefaults(
  defineProps<{
    frameworkType?: string
    isRunning?: boolean
    jobName?: string
    masterTaskLabel?: string
    masterTaskName?: string
    setTerminalRef: (element: HTMLElement | null) => void
    statusLabel?: string
    terminalConnected?: boolean
    terminalTaskName?: string
    workerCount?: number
  }>(),
  {
    frameworkType: '',
    isRunning: false,
    jobName: '',
    masterTaskLabel: '',
    masterTaskName: '',
    statusLabel: '',
    terminalConnected: false,
    terminalTaskName: '',
    workerCount: 0
  }
)

defineEmits<{
  connect: []
  disconnect: []
  fit: []
  'update:terminal-task-name': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

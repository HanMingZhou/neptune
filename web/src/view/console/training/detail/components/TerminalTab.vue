<template>
  <TerminalShell
    :can-connect="isRunning"
    :connected-label="t('connected')"
    :disabled-status-text="!isRunning ? `${t('instanceStatus')}: ${statusLabel}` : ''"
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
          style="width: 150px;"
          @update:model-value="$emit('update:terminal-task-name', $event)"
        >
          <el-option v-if="frameworkType !== 'STANDALONE'" :label="masterTaskLabel" :value="masterTaskName" />
          <el-option v-if="workerCount > 0 || frameworkType === 'STANDALONE'" label="Worker" value="worker" />
        </el-select>
    </template>
  </TerminalShell>
</template>

<script setup>
import { inject } from 'vue'
import TerminalShell from '@/components/detailPage/TerminalShell.vue'

defineProps({
  frameworkType: {
    type: String,
    default: ''
  },
  isRunning: {
    type: Boolean,
    default: false
  },
  jobName: {
    type: String,
    default: ''
  },
  masterTaskLabel: {
    type: String,
    default: ''
  },
  masterTaskName: {
    type: String,
    default: ''
  },
  setTerminalRef: {
    type: Function,
    required: true
  },
  statusLabel: {
    type: String,
    default: ''
  },
  terminalConnected: {
    type: Boolean,
    default: false
  },
  terminalTaskName: {
    type: String,
    default: ''
  },
  workerCount: {
    type: Number,
    default: 0
  }
})

defineEmits(['connect', 'disconnect', 'fit', 'update:terminal-task-name'])

const t = inject('t', (key) => key)
</script>

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
      <el-select
        v-if="pods.length > 1"
        :disabled="terminalConnected"
        :model-value="terminalPod"
        class="inference-pod-select gva-custom-select"
        popper-class="inference-pod-select__popper"
        size="default"
        @update:model-value="$emit('update:terminal-pod', $event)"
      >
        <el-option v-for="pod in pods" :key="pod.name" :label="pod.name" :value="pod.name">
          <div class="inference-pod-option">
            <span class="inference-pod-option__name">{{ pod.name }}</span>
            <span class="inference-pod-option__status">{{ pod.status || 'Unknown' }}</span>
          </div>
        </el-option>
      </el-select>
    </template>
  </TerminalShell>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import { ElOption, ElSelect } from 'element-plus'
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

<style scoped>
.inference-pod-select {
  width: clamp(180px, 22vw, 220px);
  max-width: 100%;
  flex: 0 1 220px;
}

.inference-pod-select :deep(.el-select__wrapper) {
  max-width: 100%;
}

.inference-pod-select :deep(.el-select__selected-item),
.inference-pod-select :deep(.el-select__placeholder),
.inference-pod-select :deep(.el-input__inner) {
  display: block;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.inference-pod-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-width: 0;
  width: 100%;
}

.inference-pod-option__name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgb(15 23 42);
  font-size: 0.875rem;
  font-weight: 600;
}

.inference-pod-option__status {
  flex-shrink: 0;
  padding: 2px 8px;
  border-radius: 9999px;
  background: rgb(241 245 249);
  color: rgb(71 85 105);
  font-size: 0.75rem;
  line-height: 1.25rem;
}

.dark .inference-pod-option__name {
  color: rgb(226 232 240);
}

.dark .inference-pod-option__status {
  background: rgb(39 39 42);
  color: rgb(212 212 216);
}
</style>

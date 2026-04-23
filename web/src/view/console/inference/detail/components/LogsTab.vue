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
      <el-select
        v-if="pods.length > 1"
        :disabled="logsConnected"
        :model-value="selectedPod"
        class="inference-pod-select gva-custom-select"
        popper-class="inference-pod-select__popper"
        size="default"
        @update:model-value="$emit('update:selected-pod', $event)"
      >
        <el-option v-for="pod in pods" :key="pod.name" :label="pod.name" :value="pod.name">
          <div class="inference-pod-option">
            <span class="inference-pod-option__name">{{ pod.name }}</span>
            <span class="inference-pod-option__status">{{ pod.status || 'Unknown' }}</span>
          </div>
        </el-option>
      </el-select>
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
import { ElOption, ElSelect } from 'element-plus'
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

<template>
  <DetailHeaderShell
    :get-status-class="getStatusClass"
    :get-status-label="getStatusLabel"
    :status="service.status"
    :title="service.displayName || t('inference.detail')"
    @back="$emit('back')"
  >
    <template #actions>
      <button
        v-if="isStopped"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--primary"
        @click="$emit('start')"
      >
        {{ t('start') }}
      </button>
      <button
        v-if="isRunningOrPending"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--warning"
        @click="$emit('stop')"
      >
        {{ t('stop') }}
      </button>
      <button
        v-if="isTerminalState"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--danger"
        @click="$emit('delete')"
      >
        {{ t('delete') }}
      </button>
    </template>
  </DetailHeaderShell>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type {
  ConsoleInferenceService,
  Translator
} from '@/types/consoleResource'
import DetailHeaderShell from '@/components/detailPage/DetailHeaderShell.vue'

type StatusResolver = (status?: string) => string

withDefaults(
  defineProps<{
    actionLoading?: boolean
    getStatusClass: StatusResolver
    getStatusLabel: StatusResolver
    isRunningOrPending?: boolean
    isStopped?: boolean
    isTerminalState?: boolean
    service?: ConsoleInferenceService
  }>(),
  {
    actionLoading: false,
    isRunningOrPending: false,
    isStopped: false,
    isTerminalState: false,
    service: () => ({ id: '' })
  }
)

defineEmits<{
  back: []
  delete: []
  start: []
  stop: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

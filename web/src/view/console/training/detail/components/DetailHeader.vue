<template>
  <DetailHeaderShell
    :get-status-class="getStatusClass"
    :get-status-label="getStatusLabel"
    :status="job.status"
    :title="job.displayName || job.jobName || t('trainingJobDetail')"
    @back="$emit('back')"
  >
    <template #actions>
      <button
        v-if="job.enableTensorboard"
        :disabled="!job.tensorboardUrl"
        class="detail-header-action detail-header-action--warning"
        @click="$emit('open-tensorboard')"
      >
        <span class="material-icons text-lg">analytics</span>
        {{ t('openTensorboard') }}
      </button>
      <button
        v-if="isRunning"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--neutral"
        @click="$emit('stop')"
      >
        {{ t('stop') }}
      </button>
      <button
        v-if="isTerminal"
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
import type { ConsoleTrainingJob, Translator } from '@/types/consoleResource'
import DetailHeaderShell from '@/components/detailPage/DetailHeaderShell.vue'

type StatusResolver = (status?: string) => string

withDefaults(
  defineProps<{
    actionLoading?: boolean
    getStatusClass: StatusResolver
    getStatusLabel: StatusResolver
    isRunning?: boolean
    isTerminal?: boolean
    job?: ConsoleTrainingJob
  }>(),
  {
    actionLoading: false,
    isRunning: false,
    isTerminal: false,
    job: () => ({ id: '' })
  }
)

defineEmits<{
  back: []
  delete: []
  'open-tensorboard': []
  stop: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

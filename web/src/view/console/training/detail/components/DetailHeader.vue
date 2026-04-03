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

<script setup>
import { inject } from 'vue'
import DetailHeaderShell from '@/components/detailPage/DetailHeaderShell.vue'

defineProps({
  actionLoading: {
    type: Boolean,
    default: false
  },
  getStatusClass: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  isRunning: {
    type: Boolean,
    default: false
  },
  isTerminal: {
    type: Boolean,
    default: false
  },
  job: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['back', 'delete', 'open-tensorboard', 'stop'])

const t = inject('t', (key) => key)
</script>

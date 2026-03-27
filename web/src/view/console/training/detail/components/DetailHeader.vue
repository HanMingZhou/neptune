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
        class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg text-sm font-bold transition-all flex items-center gap-2 shadow-sm hover:shadow-md active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:grayscale"
        @click="$emit('open-tensorboard')"
      >
        <span class="material-icons text-lg">analytics</span>
        {{ t('openTensorboard') }}
      </button>
      <button
        v-if="isRunning"
        :disabled="actionLoading"
        class="px-4 py-2 bg-slate-500 hover:bg-slate-600 text-white rounded-lg text-sm font-bold transition-colors disabled:opacity-50"
        @click="$emit('stop')"
      >
        {{ t('stop') }}
      </button>
      <button
        v-if="isTerminal"
        :disabled="actionLoading"
        class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-bold transition-colors disabled:opacity-50"
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

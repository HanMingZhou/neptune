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
        class="px-4 py-2 bg-emerald-500 hover:bg-emerald-600 text-white rounded-lg text-sm font-bold transition-colors disabled:opacity-50"
        @click="$emit('start')"
      >
        {{ t('start') }}
      </button>
      <button
        v-if="isRunningOrPending"
        :disabled="actionLoading"
        class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg text-sm font-bold transition-colors disabled:opacity-50"
        @click="$emit('stop')"
      >
        {{ t('stop') }}
      </button>
      <button
        v-if="isTerminalState"
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
  isRunningOrPending: {
    type: Boolean,
    default: false
  },
  isStopped: {
    type: Boolean,
    default: false
  },
  isTerminalState: {
    type: Boolean,
    default: false
  },
  service: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['back', 'delete', 'start', 'stop'])

const t = inject('t', (key) => key)
</script>

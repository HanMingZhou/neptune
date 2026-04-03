<template>
  <DetailHeaderShell
    :get-status-class="getStatusClass"
    :get-status-label="getStatusLabel"
    :status="notebook.status"
    :title="notebook.displayName || t('containerInstanceDetail')"
    @back="$emit('back')"
  >
    <template #actions>
      <button
        v-if="normalizedStatus === 'STOPPED'"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--primary"
        @click="$emit('start')"
      >
        {{ t('start') }}
      </button>
      <button
        v-if="normalizedStatus === 'RUNNING'"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--warning"
        @click="$emit('stop')"
      >
        {{ t('stop') }}
      </button>
      <button
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
import { computed, inject } from 'vue'
import DetailHeaderShell from '@/components/detailPage/DetailHeaderShell.vue'

defineEmits(['back', 'delete', 'start', 'stop'])

const t = inject('t', (key) => key)
const props = defineProps({
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
  notebook: {
    type: Object,
    default: () => ({})
  }
})

const normalizedStatus = computed(() => `${props.notebook.status || ''}`.trim().toUpperCase())
</script>

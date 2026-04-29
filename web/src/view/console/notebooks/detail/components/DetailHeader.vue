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
        v-if="normalizedStatus === 'STOPPED' || normalizedStatus === 'CLOSED'"
        :disabled="actionLoading"
        class="detail-header-action detail-header-action--secondary"
        @click="$emit('edit')"
      >
        {{ t('edit') }}
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

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { ConsoleNotebook, Translator } from '@/types/consoleResource'
import DetailHeaderShell from '@/components/detailPage/DetailHeaderShell.vue'

type StatusResolver = (status?: string) => string

defineEmits<{
  back: []
  delete: []
  edit: []
  start: []
  stop: []
}>()

const t = inject<Translator>('t', (key: string) => key)
const props = withDefaults(
  defineProps<{
    actionLoading?: boolean
    getStatusClass: StatusResolver
    getStatusLabel: StatusResolver
    notebook?: ConsoleNotebook
  }>(),
  {
    actionLoading: false,
    notebook: () => ({})
  }
)

const normalizedStatus = computed(() =>
  `${props.notebook.status || ''}`.trim().toUpperCase()
)
</script>

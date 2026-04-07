<template>
  <BaseTableToolbar
    :breadcrumbs="[t('compute'), t('storage')]"
    :description="t('storageDesc')"
    :loading="loading"
    :title="`${t('storage')}${t('admin')}`"
    @refresh="emit('refresh', $event)"
  >
    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
        {{ t('create') }}{{ t('storage') }}
      </button>
    </template>

    <StorageHintBanner />
  </BaseTableToolbar>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import StorageHintBanner from './StorageHintBanner.vue'

withDefaults(
  defineProps<{
    loading?: boolean
  }>(),
  {
    loading: false
  }
)

const emit = defineEmits<{
  create: []
  refresh: [silent: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

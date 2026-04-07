<template>
  <BaseTableToolbar
    :breadcrumbs="[t('admin'), t('apiManage')]"
    :description="t('apiManageDesc')"
    :loading="loading"
    :title="t('apiManage')"
    @refresh="emit('refresh-table', $event)"
  >
    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--secondary"
        @click="emit('fresh-cache')"
      >
        <span class="material-icons text-[18px]">refresh</span>
        {{ t('refreshCache') }}
      </button>
      <button
        class="list-toolbar-button list-toolbar-button--secondary"
        @click="emit('sync')"
      >
        <span class="material-icons text-[18px]">sync</span>
        {{ t('syncApi') }}
      </button>
      <button
        class="list-toolbar-button list-toolbar-button--danger"
        :disabled="!hasSelection"
        @click="emit('delete-selected')"
      >
        <span class="material-icons text-[18px]">delete</span>
        {{ t('delete') }}
      </button>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('create')"
      >
        <span class="material-icons text-[18px]">add</span>
        {{ t('addApi') }}
      </button>
    </template>
  </BaseTableToolbar>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    hasSelection?: boolean
    loading?: boolean
  }>(),
  {
    hasSelection: false,
    loading: false
  }
)

const emit = defineEmits<{
  create: []
  'delete-selected': []
  'fresh-cache': []
  'refresh-table': [silent: boolean]
  sync: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

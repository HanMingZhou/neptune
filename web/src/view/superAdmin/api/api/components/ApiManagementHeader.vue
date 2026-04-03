<template>
  <PageIntro
    :breadcrumbs="[t('admin'), t('apiManage')]"
    :description="t('apiManageDesc')"
    :title="t('apiManage')"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh-table', $event)" />
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
  </PageIntro>
</template>

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'

defineProps({
  hasSelection: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['create', 'delete-selected', 'fresh-cache', 'refresh-table', 'sync'])
const t = inject('t', (key) => key)
</script>

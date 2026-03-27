<template>
  <PageIntro
    :breadcrumbs="[t('admin'), t('apiManage')]"
    :description="t('apiManageDesc')"
    :title="t('apiManage')"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh-table')" />
      <button
        class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm font-medium hover:bg-slate-50"
        @click="emit('fresh-cache')"
      >
        <span class="material-icons text-[18px]">refresh</span>
        {{ t('refreshCache') }}
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm font-medium hover:bg-slate-50"
        @click="emit('sync')"
      >
        <span class="material-icons text-[18px]">sync</span>
        {{ t('syncApi') }}
      </button>
      <button
        class="flex items-center gap-2 px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-medium transition-colors"
        :disabled="!hasSelection"
        :class="{ 'opacity-50 cursor-not-allowed': !hasSelection }"
        @click="emit('delete-selected')"
      >
        <span class="material-icons text-[18px]">delete</span>
        {{ t('delete') }}
      </button>
      <button
        class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
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

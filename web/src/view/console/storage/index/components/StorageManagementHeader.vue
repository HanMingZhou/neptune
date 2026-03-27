<template>
  <PageIntro
    :breadcrumbs="[t('compute'), t('storage')]"
    :description="t('storageDesc')"
    :title="`${t('storage')}${t('admin')}`"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh')" />
      <button
        class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
        {{ t('create') }}{{ t('storage') }}
      </button>
    </template>
  </PageIntro>

  <StorageHintBanner />
</template>

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import StorageHintBanner from './StorageHintBanner.vue'

defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['create', 'refresh'])
const t = inject('t', (key) => key)
</script>

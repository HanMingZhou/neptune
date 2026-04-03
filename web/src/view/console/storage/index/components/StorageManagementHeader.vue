<template>
  <PageIntro
    :breadcrumbs="[t('compute'), t('storage')]"
    :description="t('storageDesc')"
    :title="`${t('storage')}${t('admin')}`"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh', $event)" />
      <button
        class="list-toolbar-button list-toolbar-button--primary"
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

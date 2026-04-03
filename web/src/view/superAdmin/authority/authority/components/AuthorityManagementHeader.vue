<template>
  <PageIntro
    :breadcrumbs="[t('admin'), t('roles')]"
    :description="t('roleManageDesc')"
    :title="t('roles')"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh', $event)" />
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
        {{ t('addRole') }}
      </button>
    </template>
  </PageIntro>

  <AdminAuthTipBanner />
</template>

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import AdminAuthTipBanner from '@/components/listPage/AdminAuthTipBanner.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'

defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['create', 'refresh'])
const t = inject('t', (key) => key)
</script>

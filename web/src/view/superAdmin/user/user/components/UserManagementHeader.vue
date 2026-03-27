<template>
  <PageIntro
    :breadcrumbs="[t('admin'), t('users')]"
    :description="t('userManageDesc')"
    :title="t('users')"
  >
    <template #actions>
      <RefreshButton :loading="loading" @refresh="emit('refresh')" />
      <button
        class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
        {{ t('addUser') }}
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

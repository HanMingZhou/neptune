<template>
  <BaseTableToolbar
    :breadcrumbs="[t('admin'), t('roles')]"
    :description="t('roleManageDesc')"
    :loading="loading"
    :title="t('roles')"
    @refresh="emit('refresh', $event)"
  >
    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('create')"
      >
        <span class="material-icons">add</span>
        {{ t('addRole') }}
      </button>
    </template>

    <AdminAuthTipBanner />
  </BaseTableToolbar>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import AdminAuthTipBanner from '@/components/listPage/AdminAuthTipBanner.vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'

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

<template>
  <BaseTableToolbar
    :breadcrumbs="[t('admin'), t('productManage')]"
    :description="t('productManageDesc')"
    :show-refresh="false"
    :title="t('productManage')"
  />
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)
</script>

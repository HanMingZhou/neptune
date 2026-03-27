<template>
  <PageIntro
    :breadcrumbs="[t('admin'), t('productManage')]"
    :description="t('productManageDesc')"
    :title="t('productManage')"
  />
</template>

<script setup>
import { inject } from 'vue'
import PageIntro from '@/components/listPage/PageIntro.vue'

const t = inject('t', (key) => key)
</script>

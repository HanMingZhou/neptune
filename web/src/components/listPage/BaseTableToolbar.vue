<template>
  <section class="space-y-4">
    <PageIntro
      :breadcrumbs="breadcrumbs"
      :description="description"
      :title="title"
    >
      <template #actions>
        <RefreshButton
          v-if="showRefresh"
          :loading="loading"
          @refresh="emit('refresh', $event)"
        />
        <slot name="actions" />
      </template>
    </PageIntro>

    <slot />
  </section>
</template>

<script setup lang="ts">
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from './PageIntro.vue'

withDefaults(
  defineProps<{
    breadcrumbs?: string[]
    description?: string
    loading?: boolean
    showRefresh?: boolean
    title?: string
  }>(),
  {
    breadcrumbs: () => [],
    description: '',
    loading: false,
    showRefresh: true,
    title: ''
  }
)

const emit = defineEmits<{
  refresh: [silent: boolean]
}>()
</script>

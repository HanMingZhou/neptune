<template>
  <div class="flex gap-3">
    <RefreshButton :loading="loading" @refresh="$emit('refresh', $event)" />
    <button
      class="flex items-center gap-2 rounded-lg bg-primary px-5 py-2.5 text-sm font-bold text-white shadow-lg shadow-primary/20 transition-all hover:bg-primary-hover"
      @click="$emit('create')"
    >
      <span class="material-icons text-[20px]">add</span>
      {{ t('createInstance') }}
    </button>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'

defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['create', 'refresh'])

const t = inject('t', (key) => key)
</script>

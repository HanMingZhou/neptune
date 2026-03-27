<template>
  <div class="bg-surface-light dark:bg-surface-dark border-b border-border-light dark:border-border-dark px-6 py-4 flex justify-between items-center">
    <div class="flex items-center gap-4">
      <button class="flex items-center gap-1 text-sm text-slate-500 hover:text-primary transition-colors" @click="$emit('back')">
        <span class="material-icons text-lg">arrow_back</span>
        {{ t('back') }}
      </button>
      <h2 class="text-xl font-bold">{{ title }}</h2>
      <span :class="['px-3 py-1 rounded-full text-xs font-bold', getStatusClass(status)]">
        {{ getStatusLabel(status) }}
      </span>
    </div>
    <div class="flex gap-3">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  getStatusClass: {
    type: Function,
    required: true
  },
  getStatusLabel: {
    type: Function,
    required: true
  },
  status: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
})

defineEmits(['back'])

const t = inject('t', (key) => key)
</script>

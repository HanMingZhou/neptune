<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t(titleKey) }}
    </h3>
    <div class="flex gap-3 flex-wrap">
      <button
        v-for="item in payTypes"
        :key="item.value"
        :class="[
          'px-5 py-2 rounded-lg text-sm font-medium border transition-all',
          payType === item.value
            ? 'bg-primary text-white border-primary shadow-lg shadow-primary/20'
            : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
        ]"
        @click="$emit('update:payType', item.value)"
      >
        {{ t(item.labelKey) }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  payType: {
    type: Number,
    required: true
  },
  payTypes: {
    type: Array,
    default: () => []
  },
  titleKey: {
    type: String,
    default: 'orderMethod'
  }
})

defineEmits(['update:payType'])

const t = inject('t', (key) => key)
</script>

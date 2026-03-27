<template>
  <div class="p-4 bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl flex items-center gap-4">
    <div class="relative flex-1 max-w-sm">
      <input
        :value="searchName"
        type="text"
        :placeholder="t('searchKeyPlaceholder')"
        class="w-full h-[32px] pl-9 pr-4 bg-gray-50 dark:bg-zinc-900 border border-gray-200 dark:border-border-dark rounded-md text-xs focus:ring-0 focus:border-primary focus:outline-none text-slate-700 dark:text-slate-300"
        @input="emit('update:searchName', $event.target.value)"
        @keyup.enter="emit('search')"
      />
      <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
    </div>
    <button
      class="bg-primary hover:bg-primary-hover text-white px-4 py-1.5 rounded text-xs font-bold transition-all"
      @click="emit('search')"
    >
      {{ t('search') }}
    </button>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  searchName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['search', 'update:searchName'])
const t = inject('t', (key) => key)
</script>

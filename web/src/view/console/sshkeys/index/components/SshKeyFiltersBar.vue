<template>
  <div class="console-filter-card px-5 py-4">
    <div class="list-filter-bar">
    <div class="list-filter-field min-w-[240px] max-w-sm">
      <input
        :value="searchName"
        type="text"
        :placeholder="t('searchKeyPlaceholder')"
        class="list-search-input !w-full"
        @input="emit('update:searchName', $event.target.value)"
        @keyup.enter="emit('search')"
      />
      <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
    </div>
    <button
      class="list-toolbar-button list-toolbar-button--primary"
      @click="emit('search')"
    >
      {{ t('search') }}
    </button>
    </div>
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

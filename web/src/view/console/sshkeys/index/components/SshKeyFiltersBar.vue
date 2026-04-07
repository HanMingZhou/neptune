<template>
  <BaseFilterBar>
    <template #default>
      <div class="list-filter-field min-w-[240px] max-w-sm">
        <input
          :value="searchName"
          type="text"
          :placeholder="t('searchKeyPlaceholder')"
          class="list-search-input !w-full"
          @input="handleInput"
          @keyup.enter="emit('search')"
        />
        <span
          class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]"
          >search</span
        >
      </div>
    </template>

    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('search')"
      >
        {{ t('search') }}
      </button>
    </template>
  </BaseFilterBar>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseFilterBar from '@/components/listPage/BaseFilterBar.vue'
import type { Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    searchName?: string
  }>(),
  {
    searchName: ''
  }
)

const emit = defineEmits<{
  search: []
  'update:searchName': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)

const handleInput = (event: Event): void => {
  emit('update:searchName', (event.target as HTMLInputElement).value)
}
</script>

<template>
  <BaseFilterBar wrap-main>
    <template #default>
      <div class="list-search-field max-w-[320px]">
        <span
          class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]"
          >search</span
        >
        <input
          :value="searchName"
          type="text"
          :placeholder="t('searchStoragePlaceholder')"
          class="list-search-input !w-full"
          @input="handleSearchNameInput"
          @keyup.enter="emit('refresh')"
        />
      </div>

      <el-select
        :model-value="searchStatus"
        :placeholder="`${t('status')}: ${t('all')}`"
        clearable
        class="list-filter-select !w-[168px]"
        size="small"
        @update:model-value="handleSearchStatusChange"
      >
        <el-option :label="t('Creating')" value="Creating" />
        <el-option :label="t('Bound')" value="Bound" />
        <el-option :label="t('Error')" value="Error" />
      </el-select>
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
    searchStatus?: string
  }>(),
  {
    searchName: '',
    searchStatus: ''
  }
)

const emit = defineEmits<{
  refresh: []
  'update:search-name': [value: string]
  'update:search-status': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)

const handleSearchNameInput = (event: Event): void => {
  emit('update:search-name', (event.target as HTMLInputElement).value)
}

const handleSearchStatusChange = (value?: string): void => {
  emit('update:search-status', value ?? '')
}
</script>

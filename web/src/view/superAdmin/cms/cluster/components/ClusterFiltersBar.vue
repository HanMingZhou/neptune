<template>
  <BaseFilterBar>
    <template #default>
      <el-select
        v-model="statusModel"
        :placeholder="t('status')"
        clearable
        class="list-filter-select !w-[200px]"
      >
        <el-option :label="t('allStatus')" :value="undefined" />
        <el-option :label="t('enable')" :value="1" />
        <el-option :label="t('disable')" :value="0" />
      </el-select>

      <div class="list-filter-field max-w-[280px]">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
          >search</span
        >
        <input
          v-model="keywordModel"
          type="text"
          :placeholder="t('clusterSearchDesc')"
          class="list-search-input !w-full"
          @keyup.enter="emit('search')"
        />
      </div>
    </template>

    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('search')"
      >
        <span class="material-icons text-[18px]">search</span>
        {{ t('searchQuery') }}
      </button>
      <button
        class="list-toolbar-button list-toolbar-button--secondary"
        @click="emit('reset')"
      >
        <span class="material-icons text-[18px]">autorenew</span>
        {{ t('reset') }}
      </button>
    </template>
  </BaseFilterBar>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseFilterBar from '@/components/listPage/BaseFilterBar.vue'
import type { Translator } from '@/types/consoleResource'

const props = withDefaults(
  defineProps<{
    filterKeyword?: string
    filterStatus?: number
  }>(),
  {
    filterKeyword: '',
    filterStatus: undefined
  }
)

const emit = defineEmits<{
  reset: []
  search: []
  'update:filter-keyword': [value: string]
  'update:filter-status': [value: number | undefined]
}>()

const t = inject<Translator>('t', (key: string) => key)

const keywordModel = computed({
  get: () => props.filterKeyword,
  set: (value: string) => emit('update:filter-keyword', value)
})

const statusModel = computed({
  get: () => props.filterStatus,
  set: (value?: number) => emit('update:filter-status', value)
})
</script>

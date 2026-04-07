<template>
  <BaseFilterBar>
    <template #default>
      <div class="relative flex-1 max-w-[160px]">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
          >category</span
        >
        <el-select
          v-model="filterTypeModel"
          :placeholder="t('imageType')"
          clearable
          class="!w-full list-filter-select gva-custom-select"
        >
          <el-option :label="t('all')" value="" />
          <el-option :label="t('systemImage')" :value="1" />
          <el-option :label="t('customImage')" :value="2" />
        </el-select>
      </div>

      <div class="relative flex-1 max-w-[160px]">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
          >workspaces</span
        >
        <el-select
          v-model="filterUsageTypeModel"
          :placeholder="t('imageUsageType')"
          clearable
          class="!w-full list-filter-select gva-custom-select"
        >
          <el-option :label="t('all')" value="" />
          <el-option :label="t('usageNotebook')" :value="1" />
          <el-option :label="t('usageTrain')" :value="2" />
          <el-option :label="t('usageInfer')" :value="3" />
        </el-select>
      </div>

      <div class="relative flex-1 max-w-[240px]">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
          >search</span
        >
        <input
          v-model="filterKeywordModel"
          type="text"
          :placeholder="t('imageSearchDesc')"
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
import type { ImageFilterValue } from '@/types/image'

const props = withDefaults(
  defineProps<{
    filterKeyword?: string
    filterType?: ImageFilterValue
    filterUsageType?: ImageFilterValue
  }>(),
  {
    filterKeyword: '',
    filterType: '',
    filterUsageType: ''
  }
)

const emit = defineEmits<{
  reset: []
  search: []
  'update:filter-keyword': [value: string]
  'update:filter-type': [value: ImageFilterValue]
  'update:filter-usage-type': [value: ImageFilterValue]
}>()

const t = inject<Translator>('t', (key: string) => key)

const filterKeywordModel = computed({
  get: () => props.filterKeyword,
  set: (value: string) => emit('update:filter-keyword', value)
})

const filterTypeModel = computed({
  get: () => props.filterType,
  set: (value?: ImageFilterValue) => emit('update:filter-type', value ?? '')
})

const filterUsageTypeModel = computed({
  get: () => props.filterUsageType,
  set: (value?: ImageFilterValue) =>
    emit('update:filter-usage-type', value ?? '')
})
</script>

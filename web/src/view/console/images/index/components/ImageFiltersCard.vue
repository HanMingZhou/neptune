<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="p-4 border-b border-border-light dark:border-border-dark flex flex-wrap gap-4 items-center">
      <div class="relative flex-1 max-w-[160px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">category</span>
        <el-select v-model="filterTypeModel" :placeholder="t('imageType')" clearable class="!w-full gva-custom-select">
          <el-option :label="t('all')" value="" />
          <el-option :label="t('systemImage')" :value="1" />
          <el-option :label="t('customImage')" :value="2" />
        </el-select>
      </div>

      <div class="relative flex-1 max-w-[160px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">workspaces</span>
        <el-select v-model="filterUsageTypeModel" :placeholder="t('imageUsageType')" clearable class="!w-full gva-custom-select">
          <el-option :label="t('all')" value="" />
          <el-option :label="t('usageNotebook')" :value="1" />
          <el-option :label="t('usageTrain')" :value="2" />
          <el-option :label="t('usageInfer')" :value="3" />
        </el-select>
      </div>

      <div class="relative flex-1 max-w-[240px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
        <input
          v-model="filterKeywordModel"
          type="text"
          :placeholder="t('imageSearchDesc')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @keyup.enter="$emit('search')"
        />
      </div>

      <div class="flex gap-2">
        <button
          class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium transition-all"
          @click="$emit('search')"
        >
          <span class="material-icons text-[18px]">search</span>
          {{ t('searchQuery') }}
        </button>
        <button
          class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark hover:bg-slate-50 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium transition-all"
          @click="$emit('reset')"
        >
          <span class="material-icons text-[18px]">autorenew</span>
          {{ t('reset') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  filterKeyword: {
    type: String,
    default: ''
  },
  filterType: {
    type: [String, Number],
    default: ''
  },
  filterUsageType: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['reset', 'search', 'update:filter-keyword', 'update:filter-type', 'update:filter-usage-type'])
const t = inject('t', (key) => key)

const filterKeywordModel = computed({
  get: () => props.filterKeyword,
  set: (value) => emit('update:filter-keyword', value)
})

const filterTypeModel = computed({
  get: () => props.filterType,
  set: (value) => emit('update:filter-type', value)
})

const filterUsageTypeModel = computed({
  get: () => props.filterUsageType,
  set: (value) => emit('update:filter-usage-type', value)
})
</script>

<style scoped>
.gva-custom-select :deep(.el-input__wrapper) {
  @apply bg-slate-50 dark:bg-zinc-800 border-none shadow-none pl-10;
}
</style>

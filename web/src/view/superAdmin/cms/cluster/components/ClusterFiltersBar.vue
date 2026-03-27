<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="p-4 border-b border-border-light dark:border-border-dark flex flex-wrap gap-4 items-center">
      <el-select
        v-model="statusModel"
        :placeholder="t('status')"
        clearable
        class="!w-[200px]"
      >
        <el-option :label="t('allStatus')" :value="undefined" />
        <el-option :label="t('enable')" :value="1" />
        <el-option :label="t('disable')" :value="0" />
      </el-select>

      <div class="relative flex-1 max-w-[280px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
        <input
          v-model="keywordModel"
          type="text"
          :placeholder="t('clusterSearchDesc')"
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
  filterStatus: {
    type: Number,
    default: undefined
  }
})

const emit = defineEmits(['reset', 'search', 'update:filter-keyword', 'update:filter-status'])
const t = inject('t', (key) => key)

const keywordModel = computed({
  get: () => props.filterKeyword,
  set: (value) => emit('update:filter-keyword', value)
})

const statusModel = computed({
  get: () => props.filterStatus,
  set: (value) => emit('update:filter-status', value)
})
</script>

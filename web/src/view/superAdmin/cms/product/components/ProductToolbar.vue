<template>
  <div class="space-y-4">
    <div class="flex gap-2">
      <button
        class="rounded-lg px-4 py-2 text-sm font-bold transition-all"
        :class="activeTab === 1 ? 'bg-primary text-white shadow-lg shadow-primary/20' : 'bg-slate-100 text-slate-600 hover:bg-slate-200 dark:bg-zinc-800 dark:text-slate-400 dark:hover:bg-zinc-700'"
        @click="$emit('update:active-tab', 1)"
      >
        <span class="material-icons mr-1 align-middle text-[18px]">memory</span>
        {{ t('computeProduct') }}
      </button>
      <button
        class="rounded-lg px-4 py-2 text-sm font-bold transition-all"
        :class="activeTab === 2 ? 'bg-primary text-white shadow-lg shadow-primary/20' : 'bg-slate-100 text-slate-600 hover:bg-slate-200 dark:bg-zinc-800 dark:text-slate-400 dark:hover:bg-zinc-700'"
        @click="$emit('update:active-tab', 2)"
      >
        <span class="material-icons mr-1 align-middle text-[18px]">folder_open</span>
        {{ t('storageProduct') }}
      </button>
    </div>

    <div class="overflow-hidden rounded-xl border border-border-light bg-surface-light shadow-sm dark:border-border-dark dark:bg-surface-dark">
      <div class="flex flex-wrap items-center gap-4 border-b border-border-light p-4 dark:border-border-dark">
        <div class="relative max-w-[180px] flex-1">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">hub</span>
          <el-select
            v-model="clusterModel"
            :placeholder="t('cluster')"
            clearable
            class="!w-full gva-custom-select"
            @change="$emit('search')"
          >
            <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
          </el-select>
        </div>

        <div class="relative max-w-[180px] flex-1">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">place</span>
          <el-select
            v-model="areaModel"
            :placeholder="t('area')"
            clearable
            class="!w-full gva-custom-select"
            @change="$emit('search')"
          >
            <el-option v-for="area in areas" :key="area" :label="area" :value="area" />
          </el-select>
        </div>

        <div class="relative max-w-[240px] flex-1">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">search</span>
          <input
            v-model="keywordModel"
            type="text"
            :placeholder="t('searchProductDesc')"
            class="w-full rounded-lg border-none bg-slate-50 py-2 pl-10 pr-4 text-sm outline-none focus:ring-1 focus:ring-primary dark:bg-zinc-800"
            @keyup.enter="$emit('search')"
          />
        </div>

        <div class="flex gap-2">
          <button
            class="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-all hover:bg-primary-hover"
            @click="$emit('search')"
          >
            <span class="material-icons text-[18px]">search</span>
            {{ t('searchQuery') }}
          </button>
          <button
            class="flex items-center gap-2 rounded-lg border border-border-light bg-white px-4 py-2 text-sm font-medium transition-all hover:bg-slate-50 dark:border-border-dark dark:bg-zinc-800 dark:hover:bg-zinc-700"
            @click="$emit('reset')"
          >
            <span class="material-icons text-[18px]">autorenew</span>
            {{ t('reset') }}
          </button>
        </div>

        <div class="ml-auto flex gap-3">
          <RefreshButton :loading="loading" @refresh="$emit('refresh', $event)" />
          <button
            class="flex items-center gap-2 rounded-lg bg-primary px-5 py-2.5 text-sm font-bold text-white shadow-lg shadow-primary/20 transition-all hover:bg-primary-hover"
            @click="$emit('create')"
          >
            <span class="material-icons text-[20px]">add</span>
            {{ activeTab === 1 ? t('newComputeProduct') : t('newStorageProduct') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'

const props = defineProps({
  activeTab: {
    type: Number,
    default: 1
  },
  areas: {
    type: Array,
    default: () => []
  },
  clusters: {
    type: Array,
    default: () => []
  },
  filterArea: {
    type: [String, Number],
    default: ''
  },
  filterClusterId: {
    type: [String, Number],
    default: ''
  },
  filterKeyword: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'create',
  'refresh',
  'reset',
  'search',
  'update:active-tab',
  'update:filter-area',
  'update:filter-cluster-id',
  'update:filter-keyword'
])

const t = inject('t', (key) => key)

const clusterModel = computed({
  get: () => props.filterClusterId,
  set: (value) => emit('update:filter-cluster-id', value)
})

const areaModel = computed({
  get: () => props.filterArea,
  set: (value) => emit('update:filter-area', value)
})

const keywordModel = computed({
  get: () => props.filterKeyword,
  set: (value) => emit('update:filter-keyword', value)
})
</script>

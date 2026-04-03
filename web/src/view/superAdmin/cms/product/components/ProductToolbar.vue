<template>
  <div class="console-filter-card overflow-hidden">
    <div class="border-b border-border-light px-5 py-4 dark:border-border-dark">
      <div class="flex flex-wrap items-center gap-2 min-w-0">
        <button
          class="list-toolbar-button min-w-0 max-w-full"
          :class="activeTab === 1 ? 'list-toolbar-button--primary' : 'list-toolbar-button--secondary'"
          @click="$emit('update:active-tab', 1)"
        >
          <span class="material-icons mr-1 align-middle text-[18px]">memory</span>
          <span class="truncate">{{ t('computeProduct') }}</span>
        </button>
        <button
          class="list-toolbar-button min-w-0 max-w-full"
          :class="activeTab === 2 ? 'list-toolbar-button--primary' : 'list-toolbar-button--secondary'"
          @click="$emit('update:active-tab', 2)"
        >
          <span class="material-icons mr-1 align-middle text-[18px]">folder_open</span>
          <span class="truncate">{{ t('storageProduct') }}</span>
        </button>

        <div class="list-toolbar-actions list-toolbar-actions--push product-toolbar-actions">
          <RefreshButton :loading="loading" @refresh="$emit('refresh', $event)" />
          <button
            class="list-toolbar-button list-toolbar-button--primary min-w-0 max-w-full"
            @click="$emit('create')"
          >
            <span class="material-icons text-[20px]">add</span>
            <span class="hidden xl:inline">{{ activeTab === 1 ? t('newComputeProduct') : t('newStorageProduct') }}</span>
            <span class="xl:hidden">{{ t('create') }}</span>
          </button>
        </div>
      </div>
    </div>

    <div class="px-5 py-4">
      <div class="list-filter-bar">
        <div class="list-toolbar-main">
          <div class="list-filter-field list-filter-field--compact max-w-[180px]">
            <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">hub</span>
            <el-select
              v-model="clusterModel"
              :placeholder="t('cluster')"
              clearable
              class="!w-full list-filter-select gva-custom-select"
              @change="$emit('search')"
            >
              <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
            </el-select>
          </div>

          <div class="list-filter-field list-filter-field--compact max-w-[180px]">
            <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">place</span>
            <el-select
              v-model="areaModel"
              :placeholder="t('area')"
              clearable
              class="!w-full list-filter-select gva-custom-select"
              @change="$emit('search')"
            >
              <el-option v-for="area in areas" :key="area" :label="area" :value="area" />
            </el-select>
          </div>

          <div class="list-search-field max-w-[280px]">
            <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-[20px] text-slate-400">search</span>
            <input
              v-model="keywordModel"
              type="text"
              :placeholder="t('searchProductDesc')"
              class="list-search-input !w-full"
              @keyup.enter="$emit('search')"
            />
          </div>
        </div>

        <div class="list-toolbar-actions">
          <button
            class="list-toolbar-button list-toolbar-button--primary"
            @click="$emit('search')"
          >
            <span class="material-icons text-[18px]">search</span>
            {{ t('searchQuery') }}
          </button>
          <button
            class="list-toolbar-button list-toolbar-button--secondary"
            @click="$emit('reset')"
          >
            <span class="material-icons text-[18px]">autorenew</span>
            {{ t('reset') }}
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

<style scoped>
.product-toolbar-actions {
  justify-content: flex-end;
}
</style>

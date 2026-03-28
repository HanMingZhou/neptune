<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="list-filter-bar border-b border-border-light p-4 dark:border-border-dark">
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
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
        <input
          v-model="keywordModel"
          type="text"
          :placeholder="t('clusterSearchDesc')"
          class="list-search-input !w-full"
          @keyup.enter="$emit('search')"
        />
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

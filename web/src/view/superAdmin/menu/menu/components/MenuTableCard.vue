<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="p-4 border-b border-border-light dark:border-border-dark flex flex-wrap gap-4 items-center">
      <div class="relative flex-1 max-w-[240px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
        <input
          v-model="searchKeywordModel"
          type="text"
          :placeholder="t('searchMenuDesc')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-900 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @keyup.enter="$emit('search')"
        />
      </div>

      <div class="flex gap-2">
        <button
          class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium"
          @click="$emit('search')"
        >
          <span class="material-icons text-[18px]">search</span>
          {{ t('searchQuery') }}
        </button>
        <button
          class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark hover:bg-slate-50 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium"
          @click="$emit('reset')"
        >
          <span class="material-icons text-[18px]">refresh</span>
          {{ t('reset') }}
        </button>
      </div>
    </div>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="w-full">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-3 py-3 whitespace-nowrap">{{ t('id') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('displayName') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('menuIcon') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('routeName') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('path') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('status') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('parentMenu') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('sort') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('componentPath') }}</th>
            <th class="px-3 py-3 whitespace-nowrap text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <template v-for="row in items" :key="row.ID">
            <tr class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors">
              <td class="px-3 py-3 text-sm font-mono text-slate-500 whitespace-nowrap">{{ row.ID }}</td>
              <td class="px-3 py-3 font-bold text-primary text-sm whitespace-nowrap">{{ row.meta?.title }}</td>
              <td class="px-3 py-3 whitespace-nowrap">
                <div v-if="row.meta?.icon" class="flex items-center gap-1.5 text-sm text-slate-600">
                  <el-icon><component :is="row.meta.icon" /></el-icon>
                  <span>{{ row.meta.icon }}</span>
                </div>
              </td>
              <td class="px-3 py-3 text-sm text-slate-600 dark:text-slate-400 font-mono whitespace-nowrap">{{ row.name }}</td>
              <td class="px-3 py-3 text-sm text-slate-600 dark:text-slate-400 font-mono whitespace-nowrap">{{ row.path }}</td>
              <td class="px-3 py-3 whitespace-nowrap">
                <span v-if="row.hidden" class="px-2 py-0.5 rounded-full text-xs font-bold bg-slate-500/10 text-slate-500">{{ t('hidden') }}</span>
                <span v-else class="px-2 py-0.5 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-500">{{ t('visible') }}</span>
              </td>
              <td class="px-3 py-3 text-sm text-slate-500 whitespace-nowrap">{{ row.parentId }}</td>
              <td class="px-3 py-3 text-sm text-slate-500 whitespace-nowrap">{{ row.sort }}</td>
              <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap max-w-[250px] truncate" :title="row.component">{{ row.component }}</td>
              <td class="px-3 py-3 text-center whitespace-nowrap">
                <div class="flex justify-center gap-2 items-center">
                  <button
                    class="bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                    :title="t('addChildMenu')"
                    @click="$emit('add', row.ID)"
                  >
                    <span class="material-icons text-[14px]">add</span>
                    <span class="hidden xl:inline">{{ t('add') }}</span>
                  </button>
                  <button
                    class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                    :title="t('edit')"
                    @click="$emit('edit', row.ID)"
                  >
                    <span class="material-icons text-[14px]">edit</span>
                    <span class="hidden xl:inline">{{ t('edit') }}</span>
                  </button>
                  <button
                    class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                    :title="t('delete')"
                    @click="$emit('delete', row.ID)"
                  >
                    <span class="material-icons text-[14px]">delete</span>
                    <span class="hidden xl:inline">{{ t('delete') }}</span>
                  </button>
                </div>
              </td>
            </tr>

            <template v-if="row.children && row.children.length">
              <tr
                v-for="child in row.children"
                :key="child.ID"
                class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors bg-slate-50/50 dark:bg-zinc-900/30"
              >
                <td class="px-3 py-3 text-sm font-mono text-slate-400 pl-8 whitespace-nowrap">
                  <span class="text-slate-300 dark:text-zinc-600 mr-1">└</span>
                  {{ child.ID }}
                </td>
                <td class="px-3 py-3 text-sm text-slate-600 dark:text-slate-400 whitespace-nowrap">{{ child.meta?.title }}</td>
                <td class="px-3 py-3 whitespace-nowrap">
                  <div v-if="child.meta?.icon" class="flex items-center gap-1.5 text-sm text-slate-500">
                    <el-icon><component :is="child.meta.icon" /></el-icon>
                    <span>{{ child.meta.icon }}</span>
                  </div>
                </td>
                <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap">{{ child.name }}</td>
                <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap">{{ child.path }}</td>
                <td class="px-3 py-3 whitespace-nowrap">
                  <span v-if="child.hidden" class="px-2 py-0.5 rounded-full text-xs font-bold bg-slate-500/10 text-slate-500">{{ t('hidden') }}</span>
                  <span v-else class="px-2 py-0.5 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-500">{{ t('visible') }}</span>
                </td>
                <td class="px-3 py-3 text-xs text-slate-500 whitespace-nowrap">{{ child.parentId }}</td>
                <td class="px-3 py-3 text-xs text-slate-500 whitespace-nowrap">{{ child.sort }}</td>
                <td class="px-3 py-3 text-xs text-slate-400 font-mono whitespace-nowrap max-w-[250px] truncate" :title="child.component">{{ child.component }}</td>
                <td class="px-3 py-3 text-center whitespace-nowrap">
                  <div class="flex justify-center gap-2 items-center">
                    <button
                      class="bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                      :title="t('addChildMenu')"
                      @click="$emit('add', child.ID)"
                    >
                      <span class="material-icons text-[14px]">add</span>
                      <span class="hidden xl:inline">{{ t('add') }}</span>
                    </button>
                    <button
                      class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                      :title="t('edit')"
                      @click="$emit('edit', child.ID)"
                    >
                      <span class="material-icons text-[14px]">edit</span>
                      <span class="hidden xl:inline">{{ t('edit') }}</span>
                    </button>
                    <button
                      class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-0.5"
                      :title="t('delete')"
                      @click="$emit('delete', child.ID)"
                    >
                      <span class="material-icons text-[14px]">delete</span>
                      <span class="hidden xl:inline">{{ t('delete') }}</span>
                    </button>
                  </div>
                </td>
              </tr>
            </template>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  searchKeyword: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['add', 'delete', 'edit', 'reset', 'search', 'update:search-keyword'])
const t = inject('t', (key) => key)

const searchKeywordModel = computed({
  get: () => props.searchKeyword,
  set: (value) => emit('update:search-keyword', value)
})
</script>

<template>
  <div
    v-loading="loading"
    class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm"
  >
    <div class="p-4 border-b border-border-light dark:border-border-dark flex flex-wrap gap-4 items-center">
      <div class="relative flex-1 max-w-[240px]">
        <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
        <input
          :value="searchKeyword"
          type="text"
          :placeholder="t('searchRoleDesc')"
          class="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-zinc-900 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none"
          @input="emit('update:searchKeyword', $event.target.value)"
          @keyup.enter="emit('search')"
        />
      </div>
      <div class="flex gap-2">
        <button
          class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium"
          @click="emit('search')"
        >
          <span class="material-icons text-[18px]">search</span>
          {{ t('searchQuery') }}
        </button>
        <button
          class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border border-border-light dark:border-border-dark hover:bg-slate-50 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium"
          @click="resetSearch"
        >
          <span class="material-icons text-[18px]">refresh</span>
          {{ t('reset') }}
        </button>
      </div>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('roleId') }}</th>
            <th class="px-6 py-4">{{ t('roleName') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <template v-for="row in items" :key="row.authorityId">
            <tr class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors">
              <td class="px-6 py-4 text-sm font-mono text-slate-500">{{ row.authorityId }}</td>
              <td class="px-6 py-4 font-bold text-primary text-sm">{{ row.authorityName }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex justify-center gap-2 items-center flex-wrap">
                  <button
                    class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('configure', row)"
                  >
                    <span class="material-icons text-[14px]">settings</span>
                    {{ t('permission') }}
                  </button>
                  <button
                    class="bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('add-child', row.authorityId)"
                  >
                    <span class="material-icons text-[14px]">add</span>
                    {{ t('add') }}
                  </button>
                  <button
                    class="bg-purple-500/10 hover:bg-purple-500/20 text-purple-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('copy', row)"
                  >
                    <span class="material-icons text-[14px]">content_copy</span>
                    {{ t('copy') }}
                  </button>
                  <button
                    class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('edit', row)"
                  >
                    <span class="material-icons text-[14px]">edit</span>
                    {{ t('edit') }}
                  </button>
                  <button
                    class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('delete', row)"
                  >
                    <span class="material-icons text-[14px]">delete</span>
                    {{ t('delete') }}
                  </button>
                </div>
              </td>
            </tr>

            <tr
              v-for="child in row.children || []"
              :key="child.authorityId"
              class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors bg-slate-50/50 dark:bg-zinc-900/30"
            >
              <td class="px-6 py-4 text-sm font-mono text-slate-400 pl-12">
                <span class="text-slate-300 dark:text-zinc-600 mr-2">└</span>
                {{ child.authorityId }}
              </td>
              <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{{ child.authorityName }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex justify-center gap-2 items-center flex-wrap">
                  <button
                    class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('configure', child)"
                  >
                    <span class="material-icons text-[14px]">settings</span>
                    {{ t('permission') }}
                  </button>
                  <button
                    class="bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('add-child', child.authorityId)"
                  >
                    <span class="material-icons text-[14px]">add</span>
                    {{ t('add') }}
                  </button>
                  <button
                    class="bg-purple-500/10 hover:bg-purple-500/20 text-purple-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('copy', child)"
                  >
                    <span class="material-icons text-[14px]">content_copy</span>
                    {{ t('copy') }}
                  </button>
                  <button
                    class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('edit', child)"
                  >
                    <span class="material-icons text-[14px]">edit</span>
                    {{ t('edit') }}
                  </button>
                  <button
                    class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                    @click="emit('delete', child)"
                  >
                    <span class="material-icons text-[14px]">delete</span>
                    {{ t('delete') }}
                  </button>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
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

const emit = defineEmits([
  'add-child',
  'configure',
  'copy',
  'delete',
  'edit',
  'reset',
  'search',
  'update:searchKeyword'
])

const t = inject('t', (key) => key)

const resetSearch = () => {
  emit('update:searchKeyword', '')
  emit('reset')
}
</script>

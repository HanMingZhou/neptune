<template>
  <TableCard>
    <template #toolbar>
      <div class="list-filter-bar">
        <div class="list-filter-field max-w-[240px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
          <input
            v-model="searchKeywordModel"
            type="text"
            :placeholder="t('searchMenuDesc')"
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
            <span class="material-icons text-[18px]">refresh</span>
            {{ t('reset') }}
          </button>
        </div>
      </div>
    </template>

    <div class="overflow-x-auto" v-loading="loading">
      <table class="console-table console-table--compact w-full">
        <thead>
          <tr>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('id') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('displayName') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('menuIcon') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('routeName') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('path') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('status') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('parentMenu') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('sort') }}</th>
            <th class="px-3 py-3 whitespace-nowrap">{{ t('componentPath') }}</th>
            <th class="console-actions-header px-3 py-3 whitespace-nowrap">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <template v-for="row in items" :key="row.ID">
            <tr class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors">
              <td class="px-3 py-3 text-sm font-mono text-slate-500 whitespace-nowrap">{{ row.ID }}</td>
              <td class="px-3 py-3 font-bold text-primary text-sm whitespace-nowrap">{{ row.meta?.title }}</td>
              <td class="px-3 py-3 whitespace-nowrap">
                <div v-if="row.meta?.icon" class="flex items-center gap-1.5 text-sm text-slate-600">
                  <AppIcon :name="row.meta.icon" />
                  <span>{{ row.meta.icon }}</span>
                </div>
              </td>
              <td class="px-3 py-3 text-sm text-slate-600 dark:text-slate-400 font-mono whitespace-nowrap">{{ row.name }}</td>
              <td class="px-3 py-3 text-sm text-slate-600 dark:text-slate-400 font-mono whitespace-nowrap">{{ row.path }}</td>
              <td class="px-3 py-3 whitespace-nowrap">
                <ListToneBadge
                  :label="row.hidden ? t('hidden') : t('visible')"
                  :tone="row.hidden ? 'neutral' : 'success'"
                />
              </td>
              <td class="px-3 py-3 text-sm text-slate-500 whitespace-nowrap">{{ row.parentId }}</td>
              <td class="px-3 py-3 text-sm text-slate-500 whitespace-nowrap">{{ row.sort }}</td>
              <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap max-w-[250px] truncate" :title="row.component">{{ row.component }}</td>
              <td class="console-actions-cell px-3 py-3 whitespace-nowrap">
                <div class="list-row-actions">
                  <button
                    class="list-row-button list-row-button--neutral"
                    :title="t('addChildMenu')"
                    @click="$emit('add', row.ID)"
                  >
                    {{ t('add') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--info"
                    :title="t('edit')"
                    @click="$emit('edit', row.ID)"
                  >
                    {{ t('edit') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--danger"
                    :title="t('delete')"
                    @click="$emit('delete', row.ID)"
                  >
                    {{ t('delete') }}
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
                <td class="px-3 py-3 text-sm is-primary whitespace-nowrap">{{ child.meta?.title }}</td>
                <td class="px-3 py-3 whitespace-nowrap">
                  <div v-if="child.meta?.icon" class="flex items-center gap-1.5 text-sm text-slate-500">
                    <AppIcon :name="child.meta.icon" />
                    <span>{{ child.meta.icon }}</span>
                  </div>
                </td>
                <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap">{{ child.name }}</td>
                <td class="px-3 py-3 text-xs text-slate-500 font-mono whitespace-nowrap">{{ child.path }}</td>
                <td class="px-3 py-3 whitespace-nowrap">
                  <ListToneBadge
                    :label="child.hidden ? t('hidden') : t('visible')"
                    :tone="child.hidden ? 'neutral' : 'success'"
                  />
                </td>
                <td class="px-3 py-3 text-xs text-slate-500 whitespace-nowrap">{{ child.parentId }}</td>
                <td class="px-3 py-3 text-xs text-slate-500 whitespace-nowrap">{{ child.sort }}</td>
                <td class="px-3 py-3 text-xs text-slate-400 font-mono whitespace-nowrap max-w-[250px] truncate" :title="child.component">{{ child.component }}</td>
                <td class="console-actions-cell px-3 py-3 whitespace-nowrap">
                  <div class="list-row-actions">
                    <button
                      class="list-row-button list-row-button--neutral"
                      :title="t('addChildMenu')"
                      @click="$emit('add', child.ID)"
                    >
                      {{ t('add') }}
                    </button>
                    <button
                      class="list-row-button list-row-button--info"
                      :title="t('edit')"
                      @click="$emit('edit', child.ID)"
                    >
                      {{ t('edit') }}
                    </button>
                    <button
                      class="list-row-button list-row-button--danger"
                      :title="t('delete')"
                      @click="$emit('delete', child.ID)"
                    >
                      {{ t('delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </template>
          </template>
        </tbody>
      </table>
    </div>
  </TableCard>
</template>

<script setup>
import { computed, inject } from 'vue'
import AppIcon from '@/components/AppIcon.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'

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

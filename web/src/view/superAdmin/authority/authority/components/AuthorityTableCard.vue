<template>
  <TableCard v-loading="loading">
    <template #toolbar>
      <div class="list-filter-bar">
        <div class="list-filter-field max-w-[240px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
          <input
            :value="searchKeyword"
            type="text"
            :placeholder="t('searchRoleDesc')"
            class="list-search-input !w-full"
            @input="emit('update:searchKeyword', $event.target.value)"
            @keyup.enter="emit('search')"
          />
        </div>
        <div class="list-toolbar-actions">
          <button
            class="list-toolbar-button list-toolbar-button--primary"
            @click="emit('search')"
          >
            <span class="material-icons text-[18px]">search</span>
            {{ t('searchQuery') }}
          </button>
          <button
            class="list-toolbar-button list-toolbar-button--secondary"
            @click="resetSearch"
          >
            <span class="material-icons text-[18px]">refresh</span>
            {{ t('reset') }}
          </button>
        </div>
      </div>
    </template>
    <div class="overflow-x-auto">
      <table class="console-table console-table--compact w-full min-w-[840px]">
        <thead>
          <tr>
            <th>{{ t('roleId') }}</th>
            <th>{{ t('roleName') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <template v-for="row in items" :key="row.authorityId">
            <tr class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors">
              <td class="is-code is-secondary">{{ row.authorityId }}</td>
              <td><span class="is-primary">{{ row.authorityName }}</span></td>
              <td class="console-actions-cell">
                <div class="list-row-actions">
                  <button
                    class="list-row-button list-row-button--info"
                    @click="emit('configure', row)"
                  >
                    <span class="material-icons text-[14px]">settings</span>
                    {{ t('permission') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--neutral"
                    @click="emit('add-child', row.authorityId)"
                  >
                    <span class="material-icons text-[14px]">add</span>
                    {{ t('add') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--neutral"
                    @click="emit('copy', row)"
                  >
                    <span class="material-icons text-[14px]">content_copy</span>
                    {{ t('copy') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--info"
                    @click="emit('edit', row)"
                  >
                    <span class="material-icons text-[14px]">edit</span>
                    {{ t('edit') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--danger"
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
              <td class="is-code is-secondary pl-12">
                <span class="text-slate-300 dark:text-zinc-600 mr-2">└</span>
                {{ child.authorityId }}
              </td>
              <td><span class="is-secondary">{{ child.authorityName }}</span></td>
              <td class="console-actions-cell">
                <div class="list-row-actions">
                  <button
                    class="list-row-button list-row-button--info"
                    @click="emit('configure', child)"
                  >
                    <span class="material-icons text-[14px]">settings</span>
                    {{ t('permission') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--neutral"
                    @click="emit('add-child', child.authorityId)"
                  >
                    <span class="material-icons text-[14px]">add</span>
                    {{ t('add') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--neutral"
                    @click="emit('copy', child)"
                  >
                    <span class="material-icons text-[14px]">content_copy</span>
                    {{ t('copy') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--info"
                    @click="emit('edit', child)"
                  >
                    <span class="material-icons text-[14px]">edit</span>
                    {{ t('edit') }}
                  </button>
                  <button
                    class="list-row-button list-row-button--danger"
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
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import TableCard from '@/components/listPage/TableCard.vue'

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

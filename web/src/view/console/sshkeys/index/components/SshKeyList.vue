<template>
  <TableCard>
    <div class="overflow-x-auto">
      <table class="console-table">
        <thead>
          <tr>
            <th>{{ t('name') }}</th>
            <th>{{ t('fingerprint') }}</th>
            <th class="text-center">{{ t('status') }}</th>
            <th>{{ t('createdAt') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <div class="h-5 w-5 animate-spin rounded-full border-b-2 border-primary"></div>
                {{ t('loading') }}
              </div>
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">
              <div class="flex flex-col items-center gap-3">
                <span class="material-icons text-5xl text-slate-300">vpn_key</span>
                <div class="space-y-1">
                  <p class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ t('noSshKeyData') }}</p>
                  <p class="text-sm text-slate-500">{{ t('sshKeyManageDesc') }}</p>
                </div>
                <button class="list-toolbar-button list-toolbar-button--primary" @click="emit('create')">
                  <span class="material-icons text-[18px]">add</span>
                  {{ t('newSshKey') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-for="key in items" v-else :key="key.id">
            <td>
              <div class="flex min-w-0 flex-col gap-1">
                <span class="is-primary text-sm">{{ key.name }}</span>
                <span class="text-xs text-slate-500">{{ t('sshKey') }}</span>
              </div>
            </td>
            <td class="max-w-[420px]">
              <code class="detail-inline-chip is-code block max-w-full truncate">{{ key.fingerprint }}</code>
            </td>
            <td class="text-center">
              <span
                :class="key.isDefault ? 'console-badge console-badge--success' : 'console-badge console-badge--neutral'"
              >
                {{ key.isDefault ? t('isDefault') : t('notSet') }}
              </span>
            </td>
            <td class="text-sm text-slate-500">{{ key.createdAt }}</td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  v-if="!key.isDefault"
                  class="list-row-button list-row-button--info"
                  @click="emit('set-default', key)"
                >
                  {{ t('setDefault') }}
                </button>
                <button class="list-row-button list-row-button--danger" @click="emit('delete', key)">
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        :total="items.length"
        :total-text="t('totalRecords', { total: items.length })"
        :show-pagination="false"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['create', 'delete', 'set-default'])
const t = inject('t', (key) => key)
</script>

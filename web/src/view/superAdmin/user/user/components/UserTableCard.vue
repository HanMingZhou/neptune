<template>
  <TableCard>
    <template #toolbar>
      <div class="list-filter-bar">
        <div class="list-filter-field list-filter-field--compact max-w-[180px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">search</span>
          <input
            v-model="searchInfo.username"
            type="text"
            :placeholder="t('username')"
            class="list-search-input !w-full"
            @keyup.enter="$emit('search')"
          />
        </div>
        <div class="list-filter-field list-filter-field--compact max-w-[180px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">person</span>
          <input
            v-model="searchInfo.nickname"
            type="text"
            :placeholder="t('nickname')"
            class="list-search-input !w-full"
            @keyup.enter="$emit('search')"
          />
        </div>
        <div class="list-filter-field list-filter-field--compact max-w-[180px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">phone</span>
          <input
            v-model="searchInfo.phone"
            type="text"
            :placeholder="t('phone')"
            class="list-search-input !w-full"
            @keyup.enter="$emit('search')"
          />
        </div>
        <div class="list-filter-field list-filter-field--compact max-w-[180px]">
          <span class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]">email</span>
          <input
            v-model="searchInfo.email"
            type="text"
            :placeholder="t('email')"
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
      <table class="console-table console-table--compact w-full min-w-[1100px]">
        <thead>
          <tr>
            <th>{{ t('id') }}</th>
            <th>{{ t('username') }}</th>
            <th>{{ t('nickname') }}</th>
            <th>{{ t('phone') }}</th>
            <th>{{ t('email') }}</th>
            <th>{{ t('userRole') }}</th>
            <th>{{ t('enable') }}</th>
            <th class="console-actions-header">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.id || row.ID"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="is-secondary is-code">{{ row.id || row.ID }}</td>
            <td><span class="is-primary">{{ row.userName }}</span></td>
            <td class="is-secondary">{{ row.nickName }}</td>
            <td class="is-secondary">{{ row.phone || '-' }}</td>
            <td class="is-secondary">{{ row.email || '-' }}</td>
            <td>
              <div class="flex justify-center">
                <el-cascader
                  v-model="row.authorityIds"
                  :options="authOptions"
                  :show-all-levels="false"
                  collapse-tags
                  :props="cascaderProps"
                  :clearable="false"
                  size="small"
                  @change="$emit('authority-dirty', row)"
                  @visible-change="(flag) => $emit('change-authority', { row, flag })"
                />
              </div>
            </td>
            <td>
              <div class="flex justify-center">
                <el-switch
                  v-model="row.enable"
                  inline-prompt
                  :active-value="1"
                  :inactive-value="2"
                  @change="(value) => $emit('switch-enable', { row, value })"
                />
              </div>
            </td>
            <td class="console-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
                >
                  <span class="material-icons text-[14px]">delete</span>
                  {{ t('delete') }}
                </button>
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('edit', row)"
                >
                  <span class="material-icons text-[14px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--warning"
                  @click="$emit('reset-password', row)"
                >
                  <span class="material-icons text-[14px]">key</span>
                  {{ t('reset') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[10, 30, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup>
import { inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import TableCard from '@/components/listPage/TableCard.vue'

defineProps({
  authOptions: {
    type: Array,
    default: () => []
  },
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  searchInfo: {
    type: Object,
    required: true
  },
  total: {
    type: Number,
    default: 0
  }
})

defineEmits([
  'authority-dirty',
  'change-authority',
  'delete',
  'edit',
  'page-change',
  'reset',
  'reset-password',
  'search',
  'size-change',
  'switch-enable'
])

const t = inject('t', (key) => key)

const cascaderProps = {
  multiple: true,
  checkStrictly: true,
  label: 'authorityName',
  value: 'authorityId',
  disabled: 'disabled',
  emitPath: false
}
</script>

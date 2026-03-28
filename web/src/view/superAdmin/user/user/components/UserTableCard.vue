<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl overflow-hidden shadow-sm">
    <div class="list-filter-bar border-b border-border-light p-4 dark:border-border-dark">
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

    <div class="overflow-x-auto" v-loading="loading">
      <table class="w-full min-w-[1100px]">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('id') }}</th>
            <th class="px-6 py-4">{{ t('username') }}</th>
            <th class="px-6 py-4">{{ t('nickname') }}</th>
            <th class="px-6 py-4">{{ t('phone') }}</th>
            <th class="px-6 py-4">{{ t('email') }}</th>
            <th class="px-6 py-4">{{ t('userRole') }}</th>
            <th class="px-6 py-4">{{ t('enable') }}</th>
            <th class="px-6 py-4 text-center">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.id || row.ID"
            class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors"
          >
            <td class="px-6 py-4 text-center text-sm font-mono text-slate-500">{{ row.id || row.ID }}</td>
            <td class="px-6 py-4 text-center font-bold text-primary hover:underline cursor-pointer text-sm">{{ row.userName }}</td>
            <td class="px-6 py-4 text-center text-sm text-slate-600 dark:text-slate-400">{{ row.nickName }}</td>
            <td class="px-6 py-4 text-center text-sm text-slate-600 dark:text-slate-400">{{ row.phone || '-' }}</td>
            <td class="px-6 py-4 text-center text-sm text-slate-600 dark:text-slate-400">{{ row.email || '-' }}</td>
            <td class="px-6 py-4 text-center">
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
            <td class="px-6 py-4 text-center">
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
            <td class="px-6 py-4 text-center">
              <div class="flex justify-center gap-2 items-center">
                <button
                  class="bg-red-500/10 hover:bg-red-500/20 text-red-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                  @click="$emit('delete', row)"
                >
                  <span class="material-icons text-[14px]">delete</span>
                  {{ t('delete') }}
                </button>
                <button
                  class="bg-primary/10 hover:bg-primary/20 text-primary px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
                  @click="$emit('edit', row)"
                >
                  <span class="material-icons text-[14px]">edit</span>
                  {{ t('edit') }}
                </button>
                <button
                  class="bg-amber-500/10 hover:bg-amber-500/20 text-amber-600 px-2 py-1 rounded-sm text-xs font-bold transition-colors flex items-center gap-1"
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

    <div class="px-6 py-4 bg-slate-50 dark:bg-zinc-800/30 flex items-center justify-between border-t border-border-light dark:border-border-dark">
      <span class="text-xs text-slate-500">{{ t('totalRecords', { total }) }}</span>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        background
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

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

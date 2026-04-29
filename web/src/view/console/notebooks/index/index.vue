<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('compute'), t('notebooks')]"
      :description="t('notebookDesc')"
      :loading="loading"
      :title="`${t('notebooks')}${t('admin')}`"
      @refresh="fetchNotebooks"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="goToCreate"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('rentInstance') }}
        </button>
      </template>
    </BaseTableToolbar>

    <NotebookTableCard
      class="min-h-0 flex-1"
      :btn-loading="btnLoading"
      :get-status-style="getStatusStyle"
      :get-status-text="getStatusText"
      :loading="loading"
      :notebooks="notebooks"
      :page="page"
      :page-size="pageSize"
      :search-query="searchQuery"
      :status-filter="statusFilter"
      :total="total"
      @copy="copyText"
      @create="goToCreate"
      @delete="handleDelete"
      @detail="goToDetail"
      @edit="goToEdit"
      @page-change="handlePageChange"
      @refresh="fetchNotebooks"
      @search-change="searchQuery = $event"
      @show-key-settings="showKeySettings"
      @show-ssh="showSSHInfo"
      @size-change="handleSizeChange"
      @start="handleStart"
      @status-change="statusFilter = $event"
      @stop="handleStop"
    />

    <SshDialog
      :current-notebook="currentNotebook"
      :show-password="showPassword"
      :visible="showSSHDialog"
      @copy="copyText"
      @update:show-password="showPassword = $event"
      @update:visible="showSSHDialog = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'
import NotebookTableCard from './components/NotebookTableCard.vue'
import SshDialog from './components/SshDialog.vue'
import { useNotebookList } from './composables/useNotebookList'

const t = inject<Translator>('t', (key: string) => key)

const {
  btnLoading,
  copyText,
  currentNotebook,
  fetchNotebooks,
  getStatusStyle,
  getStatusText,
  goToCreate,
  goToDetail,
  goToEdit,
  handleDelete,
  handlePageChange,
  handleSizeChange,
  handleStart,
  handleStop,
  loading,
  notebooks,
  page,
  pageSize,
  searchQuery,
  showKeySettings,
  showPassword,
  showSSHDialog,
  showSSHInfo,
  statusFilter,
  total
} = useNotebookList()
</script>

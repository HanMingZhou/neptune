<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('compute'), t('notebooks')]"
      :description="t('notebookDesc')"
      :title="`${t('notebooks')}${t('admin')}`"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchNotebooks" />
        <button class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all" @click="goToCreate">
          <span class="material-icons text-[20px]">add</span>
          {{ t('rentInstance') }}
        </button>
      </template>
    </PageIntro>

    <NotebookTableCard
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

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import NotebookTableCard from './components/NotebookTableCard.vue'
import SshDialog from './components/SshDialog.vue'
import { useNotebookList } from './composables/useNotebookList'

const t = inject('t', (key) => key)

const {
  btnLoading,
  copyText,
  currentNotebook,
  fetchNotebooks,
  getStatusStyle,
  getStatusText,
  goToCreate,
  goToDetail,
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

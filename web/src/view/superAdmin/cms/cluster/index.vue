<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('admin'), t('clusterManage')]"
      :description="t('clusterManageDesc')"
      :loading="loading"
      :title="t('clusterManage')"
      @refresh="fetchClusters"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('clusterAdd') }}
        </button>
      </template>
    </BaseTableToolbar>

    <ManagementListShell class="min-h-0 flex-1">
      <template #filters>
        <ClusterFiltersBar
          :filter-keyword="filterKeyword"
          :filter-status="filterStatus"
          @reset="handleResetFilters"
          @search="handleSearch"
          @update:filter-keyword="filterKeyword = $event"
          @update:filter-status="filterStatus = $event"
        />
      </template>

      <ClusterTableCard
        class="min-h-0 flex-1"
        :items="clusters"
        :loading="loading"
        :page="page"
        :page-size="pageSize"
        :total="total"
        @delete="handleDelete"
        @edit="openEditDialog"
        @page-change="handleCurrentChange"
        @size-change="handleSizeChange"
        @view-kubeconfig="viewKubeConfig"
      />
    </ManagementListShell>

    <ClusterEditorDialog
      v-model="showDialog"
      :dialog-title="dialogTitle"
      :form="form"
      :is-edit="isEdit"
      :rules="formRules"
      :submitting="submitting"
      @close="closeDialog"
      @submit="handleSubmit"
    />

    <KubeConfigDialog
      v-model="showKubeConfigDialog"
      :content="viewingKubeConfig"
      @close="closeKubeConfigDialog"
      @copy="copyKubeConfig"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import ClusterEditorDialog from './components/ClusterEditorDialog.vue'
import ClusterFiltersBar from './components/ClusterFiltersBar.vue'
import ClusterTableCard from './components/ClusterTableCard.vue'
import KubeConfigDialog from './components/KubeConfigDialog.vue'
import { useClusterManagementPage } from './composables/useClusterManagementPage'

const t = inject('t', (key) => key)

const {
  closeDialog,
  closeKubeConfigDialog,
  clusters,
  copyKubeConfig,
  dialogTitle,
  fetchClusters,
  filterKeyword,
  filterStatus,
  form,
  formRules,
  handleCurrentChange,
  handleDelete,
  handleResetFilters,
  handleSearch,
  handleSizeChange,
  handleSubmit,
  initialize,
  isEdit,
  loading,
  openCreateDialog,
  openEditDialog,
  page,
  pageSize,
  showDialog,
  showKubeConfigDialog,
  submitting,
  total,
  viewKubeConfig,
  viewingKubeConfig
} = useClusterManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

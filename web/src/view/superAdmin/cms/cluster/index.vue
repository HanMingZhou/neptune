<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('admin'), t('clusterManage')]"
      :description="t('clusterManageDesc')"
      :title="t('clusterManage')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="fetchClusters" />
        <button
          class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('clusterAdd') }}
        </button>
      </template>
    </PageIntro>

    <ClusterFiltersBar
      :filter-keyword="filterKeyword"
      :filter-status="filterStatus"
      @reset="handleResetFilters"
      @search="fetchClusters"
      @update:filter-keyword="filterKeyword = $event"
      @update:filter-status="filterStatus = $event"
    />

    <ClusterTableCard
      :items="clusters"
      :loading="loading"
      @delete="handleDelete"
      @edit="openEditDialog"
      @view-kubeconfig="viewKubeConfig"
    />

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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
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
  handleDelete,
  handleResetFilters,
  handleSubmit,
  initialize,
  isEdit,
  loading,
  openCreateDialog,
  openEditDialog,
  showDialog,
  showKubeConfigDialog,
  submitting,
  viewKubeConfig,
  viewingKubeConfig
} = useClusterManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

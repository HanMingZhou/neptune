<template>
  <div class="console-page-container space-y-6">
    <StorageManagementHeader
      :loading="loading"
      @create="openCreateDialog"
      @refresh="fetchList"
    />

    <ManagementListShell>
      <template #filters>
        <StorageFiltersBar
          :search-name="searchName"
          :search-status="searchStatus"
          @refresh="fetchList"
          @update:search-name="searchName = $event"
          @update:search-status="searchStatus = $event"
        />
      </template>

      <StorageTableCard
        v-model:page="pageInfo.page"
        :btn-loading="btnLoading"
        :items="volumeList"
        :loading="loading"
        :total="total"
        @create="openCreateDialog"
        @delete="handleDelete"
        @expand="openExpandDialog"
        @refresh="fetchList"
      />
    </ManagementListShell>

    <StorageDialogsHost
      :cluster-options="clusterOptions"
      :create-form="createForm"
      :creating="creating"
      :expand-form="expandForm"
      :expanding="expanding"
      :show-create-dialog="showCreateDialog"
      :show-expand-dialog="showExpandDialog"
      :storage-products="storageProducts"
      @cluster-change="onClusterChange"
      @submit-create="handleCreate"
      @submit-expand="handleExpand"
      @update:create-visible="showCreateDialog = $event"
      @update:expand-visible="showExpandDialog = $event"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import StorageDialogsHost from './components/StorageDialogsHost.vue'
import StorageFiltersBar from './components/StorageFiltersBar.vue'
import StorageManagementHeader from './components/StorageManagementHeader.vue'
import StorageTableCard from './components/StorageTableCard.vue'
import { useStorageList } from './composables/useStorageList'

const t = inject('t', (key) => key)

const {
  btnLoading,
  clusterOptions,
  createForm,
  creating,
  expanding,
  expandForm,
  fetchAreas,
  fetchList,
  handleCreate,
  handleDelete,
  handleExpand,
  loading,
  onClusterChange,
  openCreateDialog,
  openExpandDialog,
  pageInfo,
  searchName,
  searchStatus,
  showCreateDialog,
  showExpandDialog,
  storageProducts,
  total,
  volumeList
} = useStorageList({ t })

onMounted(() => {
  fetchList()
  fetchAreas()
})
</script>

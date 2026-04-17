<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <StorageManagementHeader
      :loading="loading"
      @create="openCreateDialog"
      @refresh="fetchList"
    />

    <ManagementListShell class="min-h-0 flex-1">
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
        class="min-h-0 flex-1"
        v-model:page="pageInfo.page"
        v-model:page-size="pageInfo.pageSize"
        :btn-loading="btnLoading"
        :items="volumeList"
        :loading="loading"
        :page-size="pageInfo.pageSize"
        :total="total"
        @create="openCreateDialog"
        @delete="handleDelete"
        @edit="openEditDialog"
        @expand="openExpandDialog"
        @refresh="fetchList"
      />
    </ManagementListShell>

    <StorageDialogsHost
      :cluster-options="clusterOptions"
      :create-form="createForm"
      :creating="creating"
      :edit-form="editForm"
      :editing="editing"
      :expand-form="expandForm"
      :expanding="expanding"
      :show-create-dialog="showCreateDialog"
      :show-edit-dialog="showEditDialog"
      :show-expand-dialog="showExpandDialog"
      :storage-products="storageProducts"
      @cluster-change="onClusterChange"
      @submit-create="handleCreate"
      @submit-edit="handleEdit"
      @submit-expand="handleExpand"
      @update:create-visible="showCreateDialog = $event"
      @update:edit-visible="showEditDialog = $event"
      @update:expand-visible="showExpandDialog = $event"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import type { Translator } from '@/types/consoleResource'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import StorageDialogsHost from './components/StorageDialogsHost.vue'
import StorageFiltersBar from './components/StorageFiltersBar.vue'
import StorageManagementHeader from './components/StorageManagementHeader.vue'
import StorageTableCard from './components/StorageTableCard.vue'
import { useStorageList } from './composables/useStorageList'

const t = inject<Translator>('t', (key: string) => key)

const {
  btnLoading,
  clusterOptions,
  createForm,
  creating,
  editForm,
  editing,
  expanding,
  expandForm,
  fetchAreas,
  fetchList,
  handleCreate,
  handleDelete,
  handleEdit,
  handleExpand,
  loading,
  onClusterChange,
  openCreateDialog,
  openEditDialog,
  openExpandDialog,
  pageInfo,
  searchName,
  searchStatus,
  showCreateDialog,
  showEditDialog,
  showExpandDialog,
  storageProducts,
  total,
  volumeList
} = useStorageList({ t })

onMounted(() => {
  void fetchList()
  void fetchAreas()
})
</script>

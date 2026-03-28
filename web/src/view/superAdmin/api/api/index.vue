<template>
  <div class="console-page-container space-y-6">
    <ApiManagementHeader
      :has-selection="hasSelection"
      :loading="loading"
      @create="openCreateDialog"
      @delete-selected="onDeleteSelected"
      @fresh-cache="onFresh"
      @refresh-table="getTableData"
      @sync="onSync"
    />

    <ApiTableSection
      :all-selected="isAllSelected"
      :api-group-options="apiGroupOptions"
      :get-method-class="getMethodClass"
      :items="tableData"
      :loading="loading"
      :method-options="methodOptions"
      :page="page"
      :page-size="pageSize"
      :search-api-group="searchApiGroup"
      :search-description="searchDescription"
      :search-method="searchMethod"
      :search-path="searchPath"
      :selected-ids="selectedApiIds"
      :total="total"
      @delete="deleteSingleApi"
      @edit="openEditDialog"
      @page-change="handleCurrentChange"
      @reset="onReset"
      @search="onSubmit"
      @size-change="handleSizeChange"
      @toggle-select="toggleSelect"
      @toggle-select-all="toggleSelectAll"
      @update:search-api-group="searchApiGroup = $event"
      @update:search-description="searchDescription = $event"
      @update:search-method="searchMethod = $event"
      @update:search-path="searchPath = $event"
    />

    <ApiDialogsHost
      :api-completion-loading="apiCompletionLoading"
      :api-group-options="apiGroupOptions"
      :dialog-form-visible="dialogFormVisible"
      :dialog-title="dialogTitle"
      :form="form"
      :method-options="methodOptions"
      :rules="rules"
      :sync-api-data="syncApiData"
      :sync-api-flag="syncApiFlag"
      :syncing="syncing"
      @add-one="addApiFromSync"
      @ai-completion="apiCompletion"
      @close-dialog="closeDialog"
      @close-sync="closeSyncDialog"
      @ignore="ignoreApiEntry($event.row, $event.flag)"
      @submit-dialog="submitDialog"
      @submit-sync="submitSync"
      @update-dialog-visible="dialogFormVisible = $event"
      @update-sync-visible="syncApiFlag = $event"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import ApiDialogsHost from './components/ApiDialogsHost.vue'
import ApiManagementHeader from './components/ApiManagementHeader.vue'
import ApiTableSection from './components/ApiTableSection.vue'
import { useApiManagementPage } from './composables/useApiManagementPage'

const t = inject('t', (key) => key)

const {
  addApiFromSync,
  apiCompletion,
  apiCompletionLoading,
  apiGroupOptions,
  closeDialog,
  closeSyncDialog,
  deleteSingleApi,
  dialogFormVisible,
  dialogTitle,
  form,
  getMethodClass,
  getTableData,
  handleCurrentChange,
  handleSizeChange,
  hasSelection,
  ignoreApiEntry,
  initialize,
  isAllSelected,
  loading,
  methodOptions,
  onDeleteSelected,
  onFresh,
  onReset,
  onSubmit,
  onSync,
  openCreateDialog,
  openEditDialog,
  page,
  pageSize,
  rules,
  searchApiGroup,
  searchDescription,
  searchMethod,
  searchPath,
  selectedApiIds,
  submitDialog,
  submitSync,
  syncing,
  syncApiData,
  syncApiFlag,
  tableData,
  toggleSelect,
  toggleSelectAll,
  total
} = useApiManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('admin'), t('operationRecord')]"
      :description="t('operationRecordDesc')"
      :loading="loading"
      :title="t('operationRecord')"
      @refresh="getTableData"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--danger"
          :disabled="!hasSelection"
          @click="deleteSelectedRecords"
        >
          <span class="material-icons text-[18px]">delete</span>
          {{ t('batchDelete') }}
        </button>
      </template>
    </BaseTableToolbar>

    <OperationRecordTableCard
      class="min-h-0 flex-1"
      :all-selected="isAllSelected"
      :items="tableData"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-info="searchInfo"
      :selected-ids="selectedRecordIds"
      :total="total"
      @delete="deleteRecord"
      @page-change="handleCurrentChange"
      @reset="onReset"
      @search="onSubmit"
      @size-change="handleSizeChange"
      @toggle-select="toggleSelect"
      @toggle-select-all="toggleSelectAll"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import OperationRecordTableCard from './components/OperationRecordTableCard.vue'
import { useOperationRecordPage } from './composables/useOperationRecordPage'

const t = inject('t', (key) => key)

const {
  deleteRecord,
  deleteSelectedRecords,
  getTableData,
  handleCurrentChange,
  handleSizeChange,
  hasSelection,
  initialize,
  isAllSelected,
  loading,
  onReset,
  onSubmit,
  page,
  pageSize,
  searchInfo,
  selectedRecordIds,
  tableData,
  toggleSelect,
  toggleSelectAll,
  total
} = useOperationRecordPage({ t })

onMounted(() => {
  void initialize()
})
</script>

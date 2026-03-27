<template>
  <div class="max-w-[1400px] mx-auto space-y-6">
    <PageIntro
      :breadcrumbs="[t('admin'), t('operationRecord')]"
      :description="t('operationRecordDesc')"
      :title="t('operationRecord')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="getTableData" />
        <button
          class="flex items-center gap-2 px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-medium transition-colors"
          :disabled="!hasSelection"
          :class="{ 'opacity-50 cursor-not-allowed': !hasSelection }"
          @click="deleteSelectedRecords"
        >
          <span class="material-icons text-[18px]">delete</span>
          {{ t('batchDelete') }}
        </button>
      </template>
    </PageIntro>

    <OperationRecordTableCard
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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
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
  initialize()
})
</script>

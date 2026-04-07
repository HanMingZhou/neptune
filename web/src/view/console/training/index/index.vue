<template>
  <div class="console-page-container space-y-6">
    <BaseTableToolbar
      :breadcrumbs="[t('compute'), t('training')]"
      :description="t('trainingDesc')"
      :loading="loading"
      :title="`${t('training')}${t('admin')}`"
      @refresh="handleRefresh"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="goToCreate"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('createTraining') }}
        </button>
      </template>
    </BaseTableToolbar>

    <TrainingTableCard
      :btn-loading="btnLoading"
      :filter-status="filterStatus"
      :get-framework-label="getFrameworkLabel"
      :get-framework-style="getFrameworkStyle"
      :get-status-label="getStatusLabel"
      :get-status-style="getStatusStyle"
      :jobs="jobs"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-query="searchQuery"
      :total="total"
      @copy="copyText"
      @create="goToCreate"
      @delete="handleDelete"
      @detail="goToDetail"
      @logs="viewLogs"
      @open-tensorboard="openTensorboard"
      @page-change="handlePageChange"
      @refresh="handleRefresh"
      @search-change="searchQuery = $event"
      @size-change="handleSizeChange"
      @status-change="filterStatus = $event"
      @stop="handleStop"
    />
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'
import TrainingTableCard from './components/TrainingTableCard.vue'
import { useTrainingList } from './composables/useTrainingList'

const t = inject<Translator>('t', (key: string) => key)

const {
  btnLoading,
  copyText,
  filterStatus,
  getFrameworkLabel,
  getFrameworkStyle,
  getStatusLabel,
  getStatusStyle,
  goToCreate,
  goToDetail,
  handleDelete,
  handlePageChange,
  handleRefresh,
  handleSizeChange,
  handleStop,
  jobs,
  loading,
  openTensorboard,
  page,
  pageSize,
  searchQuery,
  total,
  viewLogs
} = useTrainingList()
</script>

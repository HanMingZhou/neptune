<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('compute'), t('training')]"
      :description="t('trainingDesc')"
      :title="`${t('training')}${t('admin')}`"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="handleRefresh" />
        <button class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all" @click="goToCreate">
          <span class="material-icons text-[20px]">add</span>
          {{ t('createTraining') }}
        </button>
      </template>
    </PageIntro>

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

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import TrainingTableCard from './components/TrainingTableCard.vue'
import { useTrainingList } from './composables/useTrainingList'

const t = inject('t', (key) => key)

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

<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('compute'), t('inference')]"
      :description="t('inferenceDesc')"
      :title="t('inference')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="handleRefresh" />
        <button class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all" @click="goToCreate">
          <span class="material-icons text-[20px]">add</span>
          {{ t('createInference') }}
        </button>
      </template>
    </PageIntro>

    <InferenceTableCard
      :btn-loading="btnLoading"
      :filter-framework="filterFramework"
      :filter-status="filterStatus"
      :get-deploy-type-style="getDeployTypeStyle"
      :get-framework-style="getFrameworkStyle"
      :get-status-label="getStatusLabel"
      :get-status-style="getStatusStyle"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-query="searchQuery"
      :services="services"
      :total="total"
      @copy="copyText"
      @create="goToCreate"
      @delete="handleDelete"
      @detail="goToDetail"
      @framework-change="filterFramework = $event"
      @logs="viewLogs"
      @page-change="handlePageChange"
      @refresh="handleRefresh"
      @search-change="searchQuery = $event"
      @size-change="handleSizeChange"
      @start="handleStart"
      @status-change="filterStatus = $event"
      @stop="handleStop"
    />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import InferenceTableCard from './components/InferenceTableCard.vue'
import { useInferenceList } from './composables/useInferenceList'

const t = inject('t', (key) => key)

const {
  btnLoading,
  copyText,
  filterFramework,
  filterStatus,
  getDeployTypeStyle,
  getFrameworkStyle,
  getStatusLabel,
  getStatusStyle,
  goToCreate,
  goToDetail,
  handleDelete,
  handlePageChange,
  handleRefresh,
  handleSizeChange,
  handleStart,
  handleStop,
  loading,
  page,
  pageSize,
  searchQuery,
  services,
  total,
  viewLogs
} = useInferenceList()
</script>

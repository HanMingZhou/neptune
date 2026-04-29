<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('compute'), t('inference')]"
      :description="t('inferenceDesc')"
      :loading="loading"
      :title="t('inference')"
      @refresh="handleRefresh"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="goToCreate"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('createInference') }}
        </button>
      </template>
    </BaseTableToolbar>

    <InferenceTableCard
      class="min-h-0 flex-1"
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
      @edit="goToEdit"
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

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import type { Translator } from '@/types/consoleResource'
import InferenceTableCard from './components/InferenceTableCard.vue'
import { useInferenceList } from './composables/useInferenceList'

const t = inject<Translator>('t', (key: string) => key)

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
  goToEdit,
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

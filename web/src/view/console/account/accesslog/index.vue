<template>
  <div class="console-page-container space-y-6">
    <BaseTableToolbar
      :description="t('operationRecordDesc')"
      :loading="loading"
      :title="t('audit.title')"
      @refresh="handleRefresh"
    />

    <ActiveSessionsSection :sessions="activeSessions" @kill="handleKill" />

    <AccessLogTable
      v-model:filter-status="filterStatus"
      v-model:page="page"
      v-model:page-size="pageSize"
      v-model:search-ip="searchIp"
      :loading="loading"
      :records="loginLogs"
      :total="total"
      @page-change="handlePageChange"
      @search="handleSearch"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import ActiveSessionsSection from './components/ActiveSessionsSection.vue'
import AccessLogTable from './components/AccessLogTable.vue'
import { useAccessLog } from './composables/useAccessLog'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)

const {
  activeSessions,
  fetchAll,
  filterStatus,
  handleKill,
  handlePageChange,
  handleRefresh,
  handleSearch,
  handleSizeChange,
  loading,
  loginLogs,
  page,
  pageSize,
  searchIp,
  total
} = useAccessLog({ t })

onMounted(() => {
  void fetchAll()
})
</script>

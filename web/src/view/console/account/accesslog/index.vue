<template>
  <div class="console-page-container space-y-6">
    <PageIntro :title="t('audit.title')" :description="t('operationRecordDesc')">
      <template #actions>
        <RefreshButton :loading="loading" @refresh="handleRefresh" />
      </template>
    </PageIntro>

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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import ActiveSessionsSection from './components/ActiveSessionsSection.vue'
import AccessLogTable from './components/AccessLogTable.vue'
import { useAccessLog } from './composables/useAccessLog'

const t = inject('t', (key) => key)

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
  fetchAll()
})
</script>

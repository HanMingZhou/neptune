<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('admin'), t('nodeManage')]"
      :description="t('nodeManageDesc')"
      :loading="loading"
      :title="t('nodeManage')"
      @refresh="refreshData"
    />

    <ManagementListShell class="min-h-0 flex-1">
      <template #filters>
        <NodeFiltersBar
          :clusters="clusters"
          :current-cluster-area="currentClusterArea"
          :filter-cluster-id="filterClusterId"
          :filter-keyword="filterKeyword"
          @reset="handleResetFilters"
          @search="handleSearch"
          @update:filter-cluster-id="handleClusterChange"
          @update:filter-keyword="filterKeyword = $event"
        />
      </template>

      <NodeTableCard
        class="min-h-0 flex-1"
        :items="nodes"
        :loading="loading"
        :page="page"
        :page-size="pageSize"
        :total="total"
        @drain="handleDrain"
        @page-change="handleCurrentChange"
        @size-change="handleSizeChange"
        @uncordon="handleUncordon"
      />
    </ManagementListShell>
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import NodeFiltersBar from './components/NodeFiltersBar.vue'
import NodeTableCard from './components/NodeTableCard.vue'
import { useNodeManagementPage } from './composables/useNodeManagementPage'

const t = inject('t', (key) => key)

const {
  clusters,
  currentClusterArea,
  filterClusterId,
  filterKeyword,
  handleClusterChange,
  handleCurrentChange,
  handleDrain,
  handleResetFilters,
  handleSearch,
  handleSizeChange,
  handleUncordon,
  initialize,
  loading,
  nodes,
  page,
  pageSize,
  refreshData,
  total
} = useNodeManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

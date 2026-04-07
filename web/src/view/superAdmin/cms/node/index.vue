<template>
  <div class="console-page-container space-y-6">
    <BaseTableToolbar
      :breadcrumbs="[t('admin'), t('nodeManage')]"
      :description="t('nodeManageDesc')"
      :loading="loading"
      :title="t('nodeManage')"
      @refresh="refreshData"
    />

    <ManagementListShell>
      <template #filters>
        <NodeFiltersBar
          :clusters="clusters"
          :current-cluster-area="currentClusterArea"
          :filter-cluster-id="filterClusterId"
          :filter-keyword="filterKeyword"
          @reset="handleResetFilters"
          @search="fetchNodes"
          @update:filter-cluster-id="handleClusterChange"
          @update:filter-keyword="filterKeyword = $event"
        />
      </template>

      <NodeTableCard
        :items="nodes"
        :loading="loading"
        @drain="handleDrain"
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
  fetchNodes,
  filterClusterId,
  filterKeyword,
  handleClusterChange,
  handleDrain,
  handleResetFilters,
  handleUncordon,
  initialize,
  loading,
  nodes,
  refreshData
} = useNodeManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

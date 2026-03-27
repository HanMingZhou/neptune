<template>
  <div class="max-w-[1400px] mx-auto space-y-6">
    <PageIntro
      :breadcrumbs="[t('admin'), t('nodeManage')]"
      :description="t('nodeManageDesc')"
      :title="t('nodeManage')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="refreshData" />
      </template>
    </PageIntro>

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

    <NodeTableCard
      :items="nodes"
      :loading="loading"
      @drain="handleDrain"
      @uncordon="handleUncordon"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
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
  initialize()
})
</script>

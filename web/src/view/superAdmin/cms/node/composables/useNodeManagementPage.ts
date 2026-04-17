import { computed, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  drainCMSNode,
  getCMSClusterList,
  getCMSNodeList,
  uncordonCMSNode
} from '@/api/cms'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsNodeListData,
  CmsNodeRow
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

interface UseNodeManagementPageOptions {
  t?: Translator
}

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

const isMasterNode = (row: CmsNodeRow): boolean => row.nodeRole === 'master'

export function useNodeManagementPage({
  t
}: UseNodeManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const page = ref(1)
  const pageSize = ref(15)
  const loading = ref(false)
  const nodes = ref<CmsNodeRow[]>([])
  const clusters = ref<CmsClusterOption[]>([])
  const filterClusterId = ref<ResourceId | undefined>(undefined)
  const filterKeyword = ref('')
  const total = ref(0)

  const currentClusterArea = computed(() => {
    const currentCluster = clusters.value.find(
      (item) => String(item.id) === String(filterClusterId.value)
    )
    return currentCluster?.area || '-'
  })

  const ensureValidClusterSelection = (): void => {
    if (clusters.value.length === 0) {
      filterClusterId.value = undefined
      return
    }

    const hasCurrentCluster = clusters.value.some(
      (item) => String(item.id) === String(filterClusterId.value)
    )
    if (!hasCurrentCluster) {
      filterClusterId.value = clusters.value[0].id
    }
  }

  const fetchClusters = async (): Promise<void> => {
    try {
      const res = (await getCMSClusterList()) as ApiResponse<CmsClusterOption[]>
      if (res.code === 0) {
        clusters.value = res.data ?? []
        ensureValidClusterSelection()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch CMS clusters:', error)
      ElMessage.error(translate('failed'))
    }
  }

  const fetchNodes = async (silent = false): Promise<void> => {
    if (!filterClusterId.value) {
      nodes.value = []
      total.value = 0
      return
    }

    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getCMSNodeList({
        clusterId: filterClusterId.value,
        keyword: filterKeyword.value || undefined,
        page: page.value,
        pageSize: pageSize.value
      })) as ApiResponse<CmsNodeListData>

      if (res.code === 0) {
        nodes.value = res.data?.nodes ?? []
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? page.value
        pageSize.value = res.data?.pageSize ?? pageSize.value
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch CMS nodes:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const refreshData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      await fetchClusters()

      if (filterClusterId.value) {
        const res = (await getCMSNodeList({
          clusterId: filterClusterId.value,
          keyword: filterKeyword.value || undefined,
          page: page.value,
          pageSize: pageSize.value
        })) as ApiResponse<CmsNodeListData>

        if (res.code === 0) {
          nodes.value = res.data?.nodes ?? []
          total.value = res.data?.total ?? 0
          page.value = res.data?.page ?? page.value
          pageSize.value = res.data?.pageSize ?? pageSize.value
        } else {
          ElMessage.error(res.msg || translate('failed'))
        }
      } else {
        nodes.value = []
        total.value = 0
      }
    } catch (error: unknown) {
      console.error('Failed to refresh node data:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await refreshData()
  }

  const handleSearch = (): void => {
    page.value = 1
    void fetchNodes()
  }

  const handleClusterChange = async (
    clusterId: ResourceId | undefined
  ): Promise<void> => {
    page.value = 1
    filterClusterId.value = clusterId
    await fetchNodes()
  }

  const handleResetFilters = (): void => {
    page.value = 1
    filterKeyword.value = ''
    void fetchNodes()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void fetchNodes()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void fetchNodes()
  }

  const handleUncordon = async (row: CmsNodeRow): Promise<void> => {
    if (!filterClusterId.value) {
      return
    }

    try {
      await ElMessageBox.confirm(
        translate('confirmUncordon', { name: row.nodeName }),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await uncordonCMSNode({
        clusterId: filterClusterId.value,
        nodeName: row.nodeName
      })

      if (res.code === 0) {
        ElMessage.success(translate('success'))
        await fetchNodes()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const handleDrain = async (row: CmsNodeRow): Promise<void> => {
    if (!filterClusterId.value) {
      return
    }

    if (isMasterNode(row)) {
      ElMessage.warning(translate('masterNodeEvictDisabled'))
      return
    }

    try {
      await ElMessageBox.confirm(
        translate('confirmDrain', { name: row.nodeName }),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await drainCMSNode({
        clusterId: filterClusterId.value,
        nodeName: row.nodeName
      })

      if (res.code === 0) {
        ElMessage.success(translate('success'))
        await fetchNodes()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('failed'))
      }
    }
  }

  return {
    clusters,
    currentClusterArea,
    fetchNodes,
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
  }
}

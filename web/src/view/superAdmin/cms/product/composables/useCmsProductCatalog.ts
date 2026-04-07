import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getCMSAreaList, getCMSClusterList, getCMSProductList } from '@/api/cms'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsProductListData,
  CmsProductRow,
  CmsProductType
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

interface UseCmsProductCatalogOptions {
  t: Translator
}

export const useCmsProductCatalog = ({ t }: UseCmsProductCatalogOptions) => {
  const loading = ref(false)
  const products = ref<CmsProductRow[]>([])
  const clusters = ref<CmsClusterOption[]>([])
  const areas = ref<string[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const activeTab = ref<CmsProductType>(1)
  const filterClusterId = ref<ResourceId | ''>('')
  const filterArea = ref('')
  const filterKeyword = ref('')

  const fetchProducts = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getCMSProductList({
        page: currentPage.value,
        pageSize: pageSize.value,
        productType: activeTab.value,
        clusterId: filterClusterId.value || undefined,
        area: filterArea.value || undefined,
        keyword: filterKeyword.value || undefined
      })) as ApiResponse<CmsProductListData>

      if (res.code === 0) {
        products.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error: unknown) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const fetchClusters = async (): Promise<void> => {
    try {
      const res = (await getCMSClusterList()) as ApiResponse<CmsClusterOption[]>
      if (res.code === 0) {
        clusters.value = res.data ?? []
      }
    } catch (error: unknown) {
      console.error(error)
    }
  }

  const fetchAreas = async (): Promise<void> => {
    try {
      const res = (await getCMSAreaList()) as ApiResponse<string[]>
      if (res.code === 0) {
        areas.value = res.data ?? []
      }
    } catch (error: unknown) {
      console.error(error)
    }
  }

  const initialize = async (): Promise<void> => {
    await Promise.all([fetchClusters(), fetchAreas()])
    await fetchProducts()
  }

  const setActiveTab = (value: CmsProductType): void => {
    if (activeTab.value === value) {
      return
    }

    activeTab.value = value
    currentPage.value = 1
    void fetchProducts()
  }

  const handleSearch = (): void => {
    currentPage.value = 1
    void fetchProducts()
  }

  const handleResetFilters = (): void => {
    filterClusterId.value = ''
    filterArea.value = ''
    filterKeyword.value = ''
    currentPage.value = 1
    void fetchProducts()
  }

  const handlePageChange = (value: number): void => {
    currentPage.value = value
    void fetchProducts()
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    currentPage.value = 1
    void fetchProducts()
  }

  return {
    activeTab,
    areas,
    clusters,
    currentPage,
    fetchAreas,
    fetchClusters,
    fetchProducts,
    filterArea,
    filterClusterId,
    filterKeyword,
    handlePageChange,
    handleResetFilters,
    handleSearch,
    handleSizeChange,
    initialize,
    loading,
    pageSize,
    products,
    setActiveTab,
    total
  }
}

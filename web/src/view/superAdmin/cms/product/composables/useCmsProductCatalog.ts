import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getCMSAreaList, getCMSClusterList, getCMSProductList } from '@/api/cms'
import { getProductFilters } from '@/api/product'
import type {
  FilterOption,
  ProductFilterData,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsProductListData,
  CmsProductFilterResourceType,
  CmsProductPriceField,
  CmsProductRow,
  CmsProductType
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'

interface UseCmsProductCatalogOptions {
  t: Translator
}

export const useCmsProductCatalog = ({ t }: UseCmsProductCatalogOptions) => {
  const getDefaultPriceField = (
    productType: CmsProductType
  ): CmsProductPriceField =>
    productType === 2 ? 5 : 1

  const loading = ref(false)
  const products = ref<CmsProductRow[]>([])
  const clusters = ref<CmsClusterOption[]>([])
  const areas = ref<string[]>([])
  const gpuModels = ref<FilterOption[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(15)
  const activeTab = ref<CmsProductType>(1)
  const filterClusterId = ref<ResourceId | ''>('')
  const filterArea = ref('')
  const filterResourceType = ref<CmsProductFilterResourceType>('')
  const filterGpuModel = ref('')
  const filterAvailableMin = ref<number | undefined>(undefined)
  const filterAvailableMax = ref<number | undefined>(undefined)
  const filterMaxInstancesMin = ref<number | undefined>(undefined)
  const filterMaxInstancesMax = ref<number | undefined>(undefined)
  const filterUsedCapacityMin = ref<number | undefined>(undefined)
  const filterUsedCapacityMax = ref<number | undefined>(undefined)
  const filterPriceField = ref<CmsProductPriceField>(
    getDefaultPriceField(activeTab.value)
  )
  const filterPriceMin = ref<number | undefined>(undefined)
  const filterPriceMax = ref<number | undefined>(undefined)
  const filterKeyword = ref('')
  let searchTimer: ReturnType<typeof setTimeout> | null = null

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
        resourceType:
          activeTab.value === 1 ? filterResourceType.value || undefined : undefined,
        gpuModel:
          activeTab.value === 1 ? filterGpuModel.value || undefined : undefined,
        availableMin: filterAvailableMin.value,
        availableMax: filterAvailableMax.value,
        maxInstancesMin: filterMaxInstancesMin.value,
        maxInstancesMax: filterMaxInstancesMax.value,
        usedCapacityMin: filterUsedCapacityMin.value,
        usedCapacityMax: filterUsedCapacityMax.value,
        priceType: filterPriceField.value,
        priceMin: filterPriceMin.value,
        priceMax: filterPriceMax.value,
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
      ElMessage.error(getErrorMessage(error, t('failed')))
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

  const fetchFilterOptions = async (): Promise<void> => {
    if (activeTab.value !== 1) {
      gpuModels.value = []
      return
    }

    try {
      const res = (await getProductFilters({
        productType: activeTab.value
      })) as ApiResponse<ProductFilterData>

      if (res.code === 0) {
        gpuModels.value = res.data?.gpuModels ?? []
      }
    } catch (error: unknown) {
      console.error(error)
    }
  }

  const initialize = async (): Promise<void> => {
    await Promise.all([fetchClusters(), fetchAreas(), fetchFilterOptions()])
    await fetchProducts()
  }

  const setActiveTab = (value: CmsProductType): void => {
    if (activeTab.value === value) {
      return
    }

    activeTab.value = value
    filterPriceField.value = getDefaultPriceField(value)
    if (value === 2) {
      filterResourceType.value = ''
      filterGpuModel.value = ''
    }
    void fetchFilterOptions()
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
    filterResourceType.value = ''
    filterGpuModel.value = ''
    filterAvailableMin.value = undefined
    filterAvailableMax.value = undefined
    filterMaxInstancesMin.value = undefined
    filterMaxInstancesMax.value = undefined
    filterUsedCapacityMin.value = undefined
    filterUsedCapacityMax.value = undefined
    filterPriceField.value = getDefaultPriceField(activeTab.value)
    filterPriceMin.value = undefined
    filterPriceMax.value = undefined
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

  const scheduleSearch = (): void => {
    if (searchTimer) {
      clearTimeout(searchTimer)
    }

    searchTimer = setTimeout(() => {
      currentPage.value = 1
      void fetchProducts(true)
    }, 250)
  }

  watch(
    [
      filterClusterId,
      filterArea,
      filterResourceType,
      filterGpuModel,
      filterAvailableMin,
      filterAvailableMax,
      filterMaxInstancesMin,
      filterMaxInstancesMax,
      filterUsedCapacityMin,
      filterUsedCapacityMax,
      filterPriceField,
      filterPriceMin,
      filterPriceMax
    ],
    () => {
      scheduleSearch()
    }
  )

  watch(filterResourceType, (value) => {
    if (value === 'cpu' && filterGpuModel.value) {
      filterGpuModel.value = ''
    }
  })

  watch(filterKeyword, () => {
    scheduleSearch()
  })

  return {
    activeTab,
    areas,
    clusters,
    currentPage,
    fetchAreas,
    fetchClusters,
    fetchFilterOptions,
    fetchProducts,
    filterArea,
    filterAvailableMax,
    filterAvailableMin,
    filterClusterId,
    filterResourceType,
    filterGpuModel,
    gpuModels,
    filterKeyword,
    filterMaxInstancesMax,
    filterMaxInstancesMin,
    filterPriceField,
    filterPriceMax,
    filterPriceMin,
    filterUsedCapacityMax,
    filterUsedCapacityMin,
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


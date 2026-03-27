import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getCMSAreaList,
  getCMSClusterList,
  getCMSProductList
} from '@/api/cms'

export const useCmsProductCatalog = ({ t }) => {
  const loading = ref(false)
  const products = ref([])
  const clusters = ref([])
  const areas = ref([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const activeTab = ref(1)
  const filterClusterId = ref('')
  const filterArea = ref('')
  const filterKeyword = ref('')

  const fetchProducts = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getCMSProductList({
        page: currentPage.value,
        pageSize: pageSize.value,
        productType: activeTab.value,
        clusterId: filterClusterId.value || undefined,
        area: filterArea.value || undefined,
        keyword: filterKeyword.value || undefined
      })

      if (res.code === 0) {
        products.value = res.data?.list || []
        total.value = res.data?.total || 0
      } else {
        ElMessage.error(res.msg || t('failed'))
      }
    } catch (error) {
      console.error(error)
      ElMessage.error(t('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const fetchClusters = async () => {
    try {
      const res = await getCMSClusterList()
      if (res.code === 0) {
        clusters.value = res.data || []
      }
    } catch (error) {
      console.error(error)
    }
  }

  const fetchAreas = async () => {
    try {
      const res = await getCMSAreaList()
      if (res.code === 0) {
        areas.value = res.data || []
      }
    } catch (error) {
      console.error(error)
    }
  }

  const initialize = async () => {
    await Promise.all([
      fetchClusters(),
      fetchAreas()
    ])
    await fetchProducts()
  }

  const setActiveTab = (value) => {
    if (activeTab.value === value) {
      return
    }

    activeTab.value = value
    currentPage.value = 1
    fetchProducts()
  }

  const handleSearch = () => {
    currentPage.value = 1
    fetchProducts()
  }

  const handleResetFilters = () => {
    filterClusterId.value = ''
    filterArea.value = ''
    filterKeyword.value = ''
    currentPage.value = 1
    fetchProducts()
  }

  const handlePageChange = (value) => {
    currentPage.value = value
    fetchProducts()
  }

  const handleSizeChange = (value) => {
    pageSize.value = value
    currentPage.value = 1
    fetchProducts()
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

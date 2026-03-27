import { reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createVolume, deleteVolume, expandVolume, getAreaList, getVolumeList } from '@/api/volume'
import { getProductList } from '@/api/product'

export function useStorageList({ t }) {
  const translate = t || ((key) => key)

  const volumeList = ref([])
  const total = ref(0)
  const loading = ref(false)
  const clusterOptions = ref([])
  const searchName = ref('')
  const searchStatus = ref('')
  const storageProducts = ref([])

  const pageInfo = reactive({
    page: 1,
    pageSize: 10
  })

  const showCreateDialog = ref(false)
  const creating = ref(false)
  const createForm = reactive({
    name: '',
    size: 50,
    area: '',
    clusterId: '',
    productId: ''
  })

  const showExpandDialog = ref(false)
  const expanding = ref(false)
  const expandForm = reactive({
    id: 0,
    currentSize: '',
    minSize: 0,
    newSize: 0
  })

  const btnLoading = reactive({})

  const fetchStorageProducts = async (clusterId) => {
    try {
      const res = await getProductList({
        productType: 2,
        clusterId: clusterId || undefined,
        status: 1,
        pageSize: 100
      })

      if (res.code === 0) {
        storageProducts.value = res.data?.list || []

        const hasCurrentProduct = storageProducts.value.some((item) => item.id === createForm.productId)
        if (storageProducts.value.length > 0 && !hasCurrentProduct) {
          createForm.productId = storageProducts.value[0].id
        }
      }
    } catch (error) {
      console.error('获取存储产品失败', error)
    }
  }

  const fetchList = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getVolumeList({
        ...pageInfo,
        name: searchName.value || undefined,
        status: searchStatus.value || undefined
      })

      if (res.code === 0) {
        volumeList.value = res.data?.list || []
        total.value = res.data?.total || 0
      }
    } catch (error) {
      console.error('获取存储列表失败', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  watch([searchName, searchStatus], () => {
    pageInfo.page = 1
    fetchList()
  })

  const fetchAreas = async () => {
    try {
      const res = await getAreaList()
      if (res.code === 0) {
        clusterOptions.value = res.data?.clusters || []
        if (clusterOptions.value.length > 0 && !createForm.clusterId) {
          const firstCluster = clusterOptions.value[0]
          createForm.clusterId = firstCluster.id
          createForm.area = firstCluster.area
          fetchStorageProducts(firstCluster.id)
        }
      }
    } catch (error) {
      console.error('获取集群列表失败', error)
    }
  }

  const onClusterChange = (clusterId) => {
    const selectedCluster = clusterOptions.value.find((item) => item.id === clusterId)
    if (selectedCluster) {
      createForm.area = selectedCluster.area
    }

    createForm.productId = ''
    fetchStorageProducts(clusterId)
  }

  const openCreateDialog = () => {
    showCreateDialog.value = true
  }

  const closeCreateDialog = () => {
    showCreateDialog.value = false
  }

  const resetCreateForm = () => {
    createForm.name = ''
    createForm.size = 50
  }

  const handleCreate = async () => {
    if (!createForm.name || !createForm.clusterId || !createForm.productId) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    creating.value = true
    try {
      const res = await createVolume(createForm)
      if (res.code === 0) {
        ElMessage.success(translate('createSuccess'))
        closeCreateDialog()
        resetCreateForm()
        fetchList()
        return
      }

      ElMessage.error(res.msg || translate('createFailed'))
    } catch (error) {
      ElMessage.error(translate('createFailed'))
    } finally {
      creating.value = false
    }
  }

  const openExpandDialog = (row) => {
    expandForm.id = row.id
    expandForm.currentSize = row.size

    const sizeNum = parseInt(row.size, 10) || 50
    expandForm.minSize = sizeNum + 10
    expandForm.newSize = sizeNum + 50

    showExpandDialog.value = true
  }

  const closeExpandDialog = () => {
    showExpandDialog.value = false
  }

  const handleExpand = async () => {
    if (expandForm.newSize <= expandForm.minSize - 10) {
      ElMessage.warning(translate('capacityError'))
      return
    }

    expanding.value = true
    try {
      const res = await expandVolume({
        id: expandForm.id,
        size: expandForm.newSize
      })

      if (res.code === 0) {
        ElMessage.success(translate('expandSuccess'))
        closeExpandDialog()
        fetchList()
        return
      }

      ElMessage.error(res.msg || translate('expandFailed'))
    } catch (error) {
      ElMessage.error(translate('expandFailed'))
    } finally {
      expanding.value = false
    }
  }

  const handleDelete = async (row) => {
    if (btnLoading[row.id]) {
      return
    }

    if (row.usedBy && row.usedBy.length > 0) {
      const users = row.usedBy.map((item) => `${item.type}: ${item.name}`).join(', ')
      ElMessage.error(translate('volumeInUse', { users }))
      return
    }

    try {
      btnLoading[row.id] = true
      await ElMessageBox.confirm(translate('confirmDeleteVolume', { name: row.name }), translate('tip'), {
        type: 'warning',
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel')
      })

      const res = await deleteVolume({ id: row.id })
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        fetchList()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      // cancelled
    } finally {
      btnLoading[row.id] = false
    }
  }

  return {
    btnLoading,
    clusterOptions,
    closeCreateDialog,
    closeExpandDialog,
    createForm,
    creating,
    expanding,
    expandForm,
    fetchAreas,
    fetchList,
    handleCreate,
    handleDelete,
    handleExpand,
    loading,
    onClusterChange,
    openCreateDialog,
    openExpandDialog,
    pageInfo,
    searchName,
    searchStatus,
    showCreateDialog,
    showExpandDialog,
    storageProducts,
    total,
    volumeList
  }
}

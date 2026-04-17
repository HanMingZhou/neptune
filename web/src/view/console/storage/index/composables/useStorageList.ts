import { reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createVolume,
  deleteVolume,
  expandVolume,
  getAreaList,
  getVolumeList,
  updateVolume
} from '@/api/volume'
import { getProductList } from '@/api/product'
import type { Translator } from '@/types/consoleResource'
import type { ApiResponse } from '@/utils/request'
import type {
  StorageAreaListData,
  StorageCreateForm,
  StorageEditForm,
  StorageExpandForm,
  StorageListData,
  StorageListItem,
  StorageClusterOption,
  StorageProductListData,
  StorageProductOption
} from '@/types/storage'

interface UseStorageListOptions {
  t?: Translator
}

interface StoragePageInfo {
  page: number
  pageSize: number
}

export function useStorageList({ t }: UseStorageListOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const volumeList = ref<StorageListItem[]>([])
  const total = ref(0)
  const loading = ref(false)
  const clusterOptions = ref<StorageClusterOption[]>([])
  const searchName = ref('')
  const searchStatus = ref('')
  const storageProducts = ref<StorageProductOption[]>([])

  const pageInfo = reactive<StoragePageInfo>({
    page: 1,
    pageSize: 15
  })

  const showCreateDialog = ref(false)
  const creating = ref(false)
  const createForm = reactive<StorageCreateForm>({
    name: '',
    size: 50,
    area: '',
    clusterId: '',
    productId: ''
  })

  const showExpandDialog = ref(false)
  const expanding = ref(false)
  const expandForm = reactive<StorageExpandForm>({
    id: 0,
    currentSize: '',
    minSize: 0,
    newSize: 0
  })
  const showEditDialog = ref(false)
  const editing = ref(false)
  const editForm = reactive<StorageEditForm>({
    id: 0,
    name: ''
  })

  const btnLoading = reactive<Record<string, boolean>>({})

  const fetchStorageProducts = async (
    clusterId?: StorageCreateForm['clusterId']
  ): Promise<void> => {
    try {
      const res = (await getProductList({
        productType: 2,
        clusterId: clusterId || undefined,
        status: 1,
        pageSize: 150
      })) as ApiResponse<StorageProductListData>

      if (res.code === 0) {
        storageProducts.value = res.data?.list ?? []

        const hasCurrentProduct = storageProducts.value.some(
          (item) => item.id === createForm.productId
        )
        if (storageProducts.value.length > 0 && !hasCurrentProduct) {
          createForm.productId = storageProducts.value[0].id
        }
      }
    } catch (error: unknown) {
      console.error('获取存储产品失败', error)
    }
  }

  const fetchList = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getVolumeList({
        ...pageInfo,
        name: searchName.value || undefined,
        status: searchStatus.value || undefined
      })) as ApiResponse<StorageListData>

      if (res.code === 0) {
        volumeList.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
      }
    } catch (error: unknown) {
      console.error('获取存储列表失败', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  watch([searchName, searchStatus], () => {
    pageInfo.page = 1
    void fetchList()
  })

  const fetchAreas = async (): Promise<void> => {
    try {
      const res = (await getAreaList()) as ApiResponse<StorageAreaListData>
      if (res.code === 0) {
        clusterOptions.value = res.data?.clusters ?? []
        if (clusterOptions.value.length > 0 && !createForm.clusterId) {
          const firstCluster = clusterOptions.value[0]
          createForm.clusterId = firstCluster.id
          createForm.area = firstCluster.area ?? ''
          void fetchStorageProducts(firstCluster.id)
        }
      }
    } catch (error: unknown) {
      console.error('获取集群列表失败', error)
    }
  }

  const onClusterChange = (clusterId: StorageCreateForm['clusterId']): void => {
    const selectedCluster = clusterOptions.value.find(
      (item) => item.id === clusterId
    )
    if (selectedCluster) {
      createForm.area = selectedCluster.area ?? ''
    }

    createForm.productId = ''
    void fetchStorageProducts(clusterId)
  }

  const openCreateDialog = (): void => {
    showCreateDialog.value = true
  }

  const closeCreateDialog = (): void => {
    showCreateDialog.value = false
  }

  const resetCreateForm = (): void => {
    createForm.name = ''
    createForm.size = 50
  }

  const handleCreate = async (): Promise<void> => {
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
        void fetchList()
        return
      }

      ElMessage.error(res.msg || translate('createFailed'))
    } catch {
      ElMessage.error(translate('createFailed'))
    } finally {
      creating.value = false
    }
  }

  const openExpandDialog = (row: StorageListItem): void => {
    expandForm.id = row.id
    expandForm.currentSize = row.size ? String(row.size) : ''

    const sizeNum = parseInt(String(row.size ?? ''), 10) || 50
    const requestedSizeNum =
      parseInt(String(row.requestedSize ?? row.size ?? ''), 10) || sizeNum
    expandForm.minSize = requestedSizeNum + 10
    expandForm.newSize = requestedSizeNum + 50

    showExpandDialog.value = true
  }

  const closeExpandDialog = (): void => {
    showExpandDialog.value = false
  }

  const openEditDialog = (row: StorageListItem): void => {
    editForm.id = row.id
    editForm.name = row.name
    showEditDialog.value = true
  }

  const closeEditDialog = (): void => {
    showEditDialog.value = false
  }

  const handleExpand = async (): Promise<void> => {
    if (expandForm.newSize < expandForm.minSize) {
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
        ElMessage.success(
          typeof res.msg === 'string' && res.msg
            ? res.msg
            : translate('expandSuccess')
        )
        closeExpandDialog()
        void fetchList()
        return
      }

      ElMessage.error(res.msg || translate('expandFailed'))
    } catch {
      ElMessage.error(translate('expandFailed'))
    } finally {
      expanding.value = false
    }
  }

  const handleEdit = async (): Promise<void> => {
    if (!editForm.name.trim()) {
      ElMessage.warning(translate('inputName'))
      return
    }

    editing.value = true
    try {
      const res = await updateVolume({
        id: editForm.id,
        name: editForm.name.trim()
      })
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        closeEditDialog()
        void fetchList()
        return
      }

      ElMessage.error(res.msg || translate('changeFailed'))
    } catch {
      ElMessage.error(translate('changeFailed'))
    } finally {
      editing.value = false
    }
  }

  const handleDelete = async (row: StorageListItem): Promise<void> => {
    const rowId = String(row.id)

    if (btnLoading[rowId]) {
      return
    }

    if (row.usedBy && row.usedBy.length > 0) {
      const users = row.usedBy
        .map((item) => `${item.type}: ${item.name}`)
        .join(', ')
      ElMessage.error(translate('volumeInUse', { users }))
      return
    }

    try {
      btnLoading[rowId] = true
      await ElMessageBox.confirm(
        translate('confirmDeleteVolume', { name: row.name }),
        translate('tip'),
        {
          type: 'warning',
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel')
        }
      )

      const res = await deleteVolume({ id: row.id })
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        void fetchList()
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch {
      // cancelled
    } finally {
      btnLoading[rowId] = false
    }
  }

  return {
    btnLoading,
    clusterOptions,
    closeCreateDialog,
    closeEditDialog,
    closeExpandDialog,
    createForm,
    creating,
    editForm,
    editing,
    expanding,
    expandForm,
    fetchAreas,
    fetchList,
    handleCreate,
    handleDelete,
    handleEdit,
    handleExpand,
    loading,
    onClusterChange,
    openCreateDialog,
    openEditDialog,
    openExpandDialog,
    pageInfo,
    searchName,
    searchStatus,
    showCreateDialog,
    showEditDialog,
    showExpandDialog,
    storageProducts,
    total,
    volumeList
  }
}


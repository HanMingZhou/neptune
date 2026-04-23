import { computed, reactive, ref, watch } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createImage,
  deleteImage,
  getImageList,
  updateImage
} from '@/api/image'
import { getAreaList } from '@/api/volume'
import type { Translator } from '@/types/consoleResource'
import type { CmsClusterOption } from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'
import { composeImageAddr } from '@/utils/imageRegistry'
import type {
  ImageFilterValue,
  ImageForm,
  ImageListData,
  ImageListItem,
  ImageMutationPayload
} from '@/types/image'
import type { StorageAreaListData } from '@/types/storage'

interface UseImageManagementOptions {
  t?: Translator
}

const createDefaultForm = (): ImageForm => ({
  id: null,
  clusterId: '',
  name: '',
  type: 1,
  usageType: 1,
  imageAddr: '',
  area: '',
  size: '',
  imagePath: ''
})

export function useImageManagement({ t }: UseImageManagementOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const loading = ref(false)
  const submitting = ref(false)
  const imageAddrManuallyEdited = ref(false)
  const images = ref<ImageListItem[]>([])
  const total = ref(0)
  const clusterOptions = ref<CmsClusterOption[]>([])
  const currentPage = ref(1)
  const pageSize = ref(15)
  const filterKeyword = ref('')
  const filterType = ref<ImageFilterValue>('')
  const filterUsageType = ref<ImageFilterValue>('')
  const showDialog = ref(false)
  const isEdit = ref(false)

  const form = reactive<ImageForm>(createDefaultForm())

  const formRules: FormRules<ImageForm> = {
    clusterId: [
      { required: true, message: translate('selectCluster'), trigger: 'change' }
    ],
    name: [
      { required: true, message: translate('inputName'), trigger: 'blur' }
    ],
    imagePath: [
      { required: true, message: translate('fillAllFields'), trigger: 'blur' }
    ],
    usageType: [
      { required: true, message: translate('pleaseSelect'), trigger: 'change' }
    ]
  }

  const selectedCluster = computed<CmsClusterOption | null>(
    () =>
      clusterOptions.value.find((item) => item.id === form.clusterId) || null
  )
  const generatedImageAddr = computed(() =>
    composeImageAddr(selectedCluster.value?.harborAddr, form.imagePath)
  )

  const resetForm = (): void => {
    Object.assign(form, createDefaultForm())
    imageAddrManuallyEdited.value = false
  }

  const fetchImages = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getImageList({
        page: currentPage.value,
        pageSize: pageSize.value,
        name: filterKeyword.value || undefined,
        type: filterType.value || undefined,
        usageType: filterUsageType.value || undefined
      })) as ApiResponse<ImageListData>

      if (res.code === 0) {
        images.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const fetchClusters = async (): Promise<void> => {
    try {
      const res = (await getAreaList()) as ApiResponse<StorageAreaListData>
      if (res.code === 0) {
        clusterOptions.value = (res.data?.clusters ?? []) as CmsClusterOption[]
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('failed')))
    }
  }

  const handleRefresh = (silent = false): Promise<void> => fetchImages(silent)

  const handleSearch = (): void => {
    currentPage.value = 1
    void fetchImages()
  }

  const handleReset = (): void => {
    filterKeyword.value = ''
    filterType.value = ''
    filterUsageType.value = ''
    currentPage.value = 1
    void fetchImages()
  }

  const handlePageChange = (): void => {
    void fetchImages()
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    currentPage.value = 1
    void fetchImages()
  }

  const openCreateDialog = (): void => {
    resetForm()
    isEdit.value = false
    showDialog.value = true
  }

  const openEditDialog = (row: ImageListItem): void => {
    isEdit.value = true
    form.id = row.id
    form.clusterId = row.clusterId ?? ''
    form.name = row.name
    form.type = row.type ?? 1
    form.usageType = row.usageType ?? 1
    form.imageAddr = row.image || ''
    form.area = row.area || ''
    form.size = row.size || ''
    form.imagePath = row.imagePath || ''
    imageAddrManuallyEdited.value =
      Boolean(form.imageAddr) && form.imageAddr !== generatedImageAddr.value
    showDialog.value = true
  }

  const handleDialogClosed = (): void => {
    if (!showDialog.value) {
      resetForm()
      isEdit.value = false
    }
  }

  const submitImage = async (): Promise<void> => {
    submitting.value = true

    try {
      const data: ImageMutationPayload = { ...form }
      if (!isEdit.value) {
        delete data.id
      }

      const res = isEdit.value
        ? await updateImage(data)
        : await createImage(data)
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        showDialog.value = false
        fetchImages()
        return
      }

      ElMessage.error(getErrorMessage(res, translate('failed')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      submitting.value = false
    }
  }

  const handleImageAddrInput = (value: string): void => {
    form.imageAddr = value
    imageAddrManuallyEdited.value = value !== generatedImageAddr.value
  }

  watch(
    [
      () => form.clusterId,
      () => form.imagePath,
      selectedCluster,
      generatedImageAddr
    ],
    () => {
      if (!form.clusterId) {
        form.area = ''
        if (!imageAddrManuallyEdited.value) {
          form.imageAddr = ''
        }
        return
      }

      if (!selectedCluster.value) {
        return
      }

      form.area = selectedCluster.value.area ?? ''
      if (!imageAddrManuallyEdited.value || !form.imageAddr) {
        form.imageAddr = generatedImageAddr.value
      }
    },
    { immediate: true }
  )

  const handleDelete = (row: ImageListItem): void => {
    ElMessageBox.confirm(
      translate('confirmDelete', { name: row.name }),
      translate('tip'),
      {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      }
    )
      .then(async () => {
        try {
          const res = await deleteImage({ id: row.id })
          if (res.code === 0) {
            ElMessage.success(translate('success'))
            fetchImages()
            return
          }

          ElMessage.error(getErrorMessage(res, translate('failed')))
        } catch (error) {
          ElMessage.error(getErrorMessage(error, translate('failed')))
        }
      })
      .catch(() => {})
  }

  return {
    currentPage,
    clusterOptions,
    fetchImages,
    fetchClusters,
    filterKeyword,
    filterType,
    filterUsageType,
    form,
    formRules,
    generatedImageAddr,
    handleDelete,
    handleDialogClosed,
    handleImageAddrInput,
    handlePageChange,
    handleRefresh,
    handleReset,
    handleSearch,
    handleSizeChange,
    images,
    isEdit,
    loading,
    openCreateDialog,
    openEditDialog,
    pageSize,
    showDialog,
    submitImage,
    submitting,
    total
  }
}


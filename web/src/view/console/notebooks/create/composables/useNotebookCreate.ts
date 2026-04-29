import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { RouteLocationNormalizedLoaded, Router } from 'vue-router'
import { getAggregateProductList, getProductFilters } from '@/api/product'
import { getImageList } from '@/api/image'
import { createNotebook, getNotebookDetail, updateNotebook } from '@/api/notebook'
import { getSSHKeyList } from '@/api/sshkey'
import { getVolumeList } from '@/api/volume'
import type { ApiResponse } from '@/utils/request'
import {
  DEFAULT_VOLUME_MOUNT_PATH,
  NOTEBOOK_IMAGE_TABS,
  NOTEBOOK_PAY_TYPES
} from '../constants'
import {
  getSubmitErrorMessage,
  isResourceNameErrorMessage,
  validateK8sResourceName,
  validateTensorBoardPath
} from '@/utils/resourceValidators'
import {
  buildGpuResourceFilterOptions,
  findGpuResourceFilterOption
} from '@/utils/gpuFilters'
import { decorateConsoleImages } from '@/utils/imageRegistry'
import type {
  ConsoleImage,
  ConsoleNotebookDetail,
  ConsoleProduct,
  ConsoleSshKey,
  ConsoleVolume,
  FilterOption,
  ListData,
  ProductFilterData,
  ResourceId,
  Translator
} from '@/types/consoleResource'

export type NotebookImageTab = (typeof NOTEBOOK_IMAGE_TABS)[number]['value']
export type NotebookPayType = (typeof NOTEBOOK_PAY_TYPES)[number]['value']

export interface NotebookFilters {
  area: string
  gpuModel: string
  cpuModel: string
}

export interface NotebookFieldErrors {
  instanceName: string
  tensorboardLogPath: string
}

interface NotebookCreatePayload {
  displayName: string
  productId: ResourceId
  imageId: ResourceId
  payType: NotebookPayType
  tensorBoard: boolean
  tensorBoardPath: string
  sshKeyId: ResourceId | 0
  enableSshPassword: boolean
  volumeMounts: Array<{
    pvcId: ResourceId
    mountsPath: string
  }>
}

interface NotebookUpdatePayload {
  id: ResourceId
  displayName: string
  productId: ResourceId
  imageId: ResourceId
  chargeType: NotebookPayType
}

interface UseNotebookCreateOptions {
  route: Pick<RouteLocationNormalizedLoaded, 'query'>
  t?: Translator
  router: Pick<Router, 'go' | 'push'>
}

export const useNotebookCreate = ({ route, t, router }: UseNotebookCreateOptions) => {
  const translate: Translator = (key, params) =>
    typeof t === 'function' ? t(key, params) : key

  const queryId = computed(() => {
    const value = route.query.id
    return Array.isArray(value) ? value[0] || '' : value || ''
  })
  const notebookId = computed<ResourceId | null>(() => {
    const value = Number(queryId.value)
    return Number.isFinite(value) && value > 0 ? value : null
  })
  const isEditMode = computed(() => {
    const value = route.query.mode
    const normalized = `${Array.isArray(value) ? value[0] || '' : value || ''}`
      .trim()
      .toLowerCase()
    return normalized === 'edit' && Boolean(notebookId.value)
  })
  const submitLoading = ref(false)

  const products = ref<ConsoleProduct[]>([])
  const selectedProduct = ref<ConsoleProduct | null>(null)
  const isSelectableProduct = (
    product: ConsoleProduct | null | undefined
  ): product is ConsoleProduct => (product?.available ?? 0) > 0
  const payType = ref<NotebookPayType>(1)
  const filters = ref<NotebookFilters>({ area: '', gpuModel: '', cpuModel: '' })
  const areas = ref<string[]>([])
  const gpuModels = ref<FilterOption[]>([])
  const cpuModels = ref<FilterOption[]>([])

  const images = ref<ConsoleImage[]>([])
  const selectedImage = ref<ResourceId | ''>('')
  const activeTab = ref<NotebookImageTab>('base')

  const instanceName = ref('')
  const enableTensorboard = ref(false)
  const tensorboardLogPath = ref('')
  const sshKeys = ref<ConsoleSshKey[]>([])
  const selectedSshKey = ref<ResourceId | null>(null)
  const enableSshPassword = ref(false)

  const availableVolumes = ref<ConsoleVolume[]>([])
  const selectedVolumeId = ref<ResourceId | null>(null)
  const volumeMountPath = ref(DEFAULT_VOLUME_MOUNT_PATH)
  const fieldErrors = ref<NotebookFieldErrors>({
    instanceName: '',
    tensorboardLogPath: ''
  })

  const payTypes = NOTEBOOK_PAY_TYPES

  const pageTitleKey = computed(() =>
    isEditMode.value ? 'edit' : 'createInstance'
  )

  const submitLabelKey = computed(() =>
    isEditMode.value ? 'save' : 'createNow'
  )

  const selectedImageRecord = computed(
    () => images.value.find((item) => item.id === selectedImage.value) || null
  )

  const imageTabs = computed(() =>
    NOTEBOOK_IMAGE_TABS.map((tab) => ({
      ...tab,
      label: translate(tab.labelKey)
    }))
  )

  const selectedVolumeName = computed(() => {
    const volume = availableVolumes.value.find(
      (item) => item.id === selectedVolumeId.value
    )
    if (!volume) return ''
    return volume.size ? `${volume.name} (${volume.size})` : volume.name
  })

  const priceUnitText = computed(() => {
    switch (payType.value) {
      case 1:
        return translate('unitHour')
      case 2:
        return translate('unitDay')
      case 3:
        return translate('unitWeek')
      case 4:
        return translate('unitMonth')
      default:
        return translate('unitHour')
    }
  })

  const imageDescription = computed(() => {
    switch (activeTab.value) {
      case 'base':
        return translate('baseImageDesc')
      case 'community':
        return translate('communityImageDesc')
      case 'my':
        return translate('myImageDesc')
      default:
        return ''
    }
  })

  const filteredImages = computed(() => images.value)

  const getUnitPrice = (product: ConsoleProduct | null | undefined): number => {
    if (!product) return 0

    switch (payType.value) {
      case 1:
        return product.priceHourly ?? 0
      case 2:
        return product.priceDaily ?? 0
      case 3:
        return product.priceWeekly ?? 0
      case 4:
        return product.priceMonthly ?? 0
      default:
        return 0
    }
  }

  const formatPrice = (product: ConsoleProduct | null | undefined): string => {
    const price = getUnitPrice(product)
    return payType.value === 1 ? price.toFixed(4) : price.toFixed(2)
  }

  const totalPrice = computed(() => {
    if (!selectedProduct.value) return '0.00'

    const total = getUnitPrice(selectedProduct.value)
    return payType.value === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const canCreate = computed(() =>
    Boolean(
      isSelectableProduct(selectedProduct.value) &&
        selectedImage.value &&
        instanceName.value.trim() &&
        (isEditMode.value ||
          !enableTensorboard.value ||
          (selectedVolumeId.value &&
            volumeMountPath.value.trim().startsWith('/') &&
            tensorboardLogPath.value.trim().startsWith(volumeMountPath.value.trim())))
    )
  )

  const updateFieldError = <K extends keyof NotebookFieldErrors>(
    field: K,
    message = ''
  ): void => {
    fieldErrors.value = {
      ...fieldErrors.value,
      [field]: message
    }
  }

  const updateInstanceName = (value: string): void => {
    instanceName.value = value
    if (fieldErrors.value.instanceName) {
      updateFieldError('instanceName')
    }
  }

  const updateTensorboardLogPath = (value: string): void => {
    tensorboardLogPath.value = value
    if (fieldErrors.value.tensorboardLogPath) {
      updateFieldError('tensorboardLogPath')
    }
  }

  const validateInstanceNameField = (): boolean => {
    const error = validateK8sResourceName(instanceName.value, {
      t: translate,
      fieldKey: 'instanceName',
      example: 'my-notebook'
    })

    updateFieldError('instanceName', error || '')
    return !error
  }

  const validateTensorboardLogPathField = (): boolean => {
    const error = !isEditMode.value && enableTensorboard.value
      ? validateTensorBoardPath(tensorboardLogPath.value, translate)
      : null

    updateFieldError('tensorboardLogPath', error || '')
    return !error
  }

  const validateCreateForm = (): string => {
    const isInstanceNameValid = validateInstanceNameField()
    const isTensorboardPathValid = validateTensorboardLogPathField()

    if (!isInstanceNameValid) {
      return fieldErrors.value.instanceName
    }

    if (!selectedProduct.value || !selectedImage.value) {
      return translate('fillAllFields')
    }
    if (!isSelectableProduct(selectedProduct.value)) {
      return translate('fillAllFields')
    }

    if (!isTensorboardPathValid) {
      return fieldErrors.value.tensorboardLogPath
    }
    if (!isEditMode.value && enableTensorboard.value && !selectedVolumeId.value) {
      return '启用 TensorBoard 时必须选择持久化数据卷'
    }
    if (!isEditMode.value && enableTensorboard.value && !volumeMountPath.value.trim()) {
      return '启用 TensorBoard 时必须填写数据卷挂载路径'
    }
    if (
      !isEditMode.value &&
      enableTensorboard.value &&
      !volumeMountPath.value.trim().startsWith('/')
    ) {
      return '数据卷挂载路径必须是绝对路径'
    }
    if (
      !isEditMode.value &&
      enableTensorboard.value &&
      !tensorboardLogPath.value.trim().startsWith(volumeMountPath.value.trim())
    ) {
      return 'TensorBoard 日志路径必须位于所选数据卷挂载路径下'
    }

    return ''
  }

  const syncSelectedProduct = (list: ConsoleProduct[]): void => {
    if (!list.length) {
      selectedProduct.value = null
      availableVolumes.value = []
      selectedVolumeId.value = null
      return
    }

    const firstSelectable = list.find(isSelectableProduct)

    if (!selectedProduct.value) {
      if (firstSelectable) {
        selectProduct(firstSelectable)
      } else {
        selectedProduct.value = null
      }
      return
    }

    const nextSelected = list.find(
      (item) => item.id === selectedProduct.value?.id && isSelectableProduct(item)
    )
    if (!nextSelected) {
      if (firstSelectable) {
        selectProduct(firstSelectable)
      } else {
        selectedProduct.value = null
      }
      return
    }

    selectedProduct.value = nextSelected
  }

  const loadProducts = async (clusterId?: ResourceId | null): Promise<void> => {
    try {
      const params: {
        page: number
        pageSize: number
        productType: number
        area?: string
        gpuModel?: string
        gpuResourceType?: 'gpu' | 'vgpu'
        vGpuNumber?: number
        vGpuMemory?: number
        vGpuCores?: number
        cpuModel?: string
        clusterId?: ResourceId
      } = { page: 1, pageSize: 100, productType: 1 }

      if (filters.value.area) params.area = filters.value.area
      if (filters.value.gpuModel) {
        const selectedGpuFilter = findGpuResourceFilterOption(
          gpuModels.value,
          filters.value.gpuModel
        )

        if (selectedGpuFilter) {
          params.gpuModel =
            selectedGpuFilter.gpuModel || selectedGpuFilter.model || ''
          params.gpuResourceType = selectedGpuFilter.resourceType

          if (selectedGpuFilter.resourceType === 'vgpu') {
            params.vGpuNumber = selectedGpuFilter.vGpuNumber
            params.vGpuMemory = selectedGpuFilter.vGpuMemory
            params.vGpuCores = selectedGpuFilter.vGpuCores
          }
        } else {
          params.gpuModel = filters.value.gpuModel
        }
      }
      if (filters.value.cpuModel) params.cpuModel = filters.value.cpuModel
      if (clusterId || selectedImageRecord.value?.clusterId) {
        params.clusterId = (clusterId ||
          selectedImageRecord.value?.clusterId) as ResourceId
      }

      const res = (await getAggregateProductList(params)) as ApiResponse<
        ListData<ConsoleProduct>
      >
      if (res.code === 0) {
        const list = res.data?.list || []
        products.value = list
        syncSelectedProduct(list)
      }
    } catch (error) {
      console.error('加载产品失败', error)
      ElMessage.error(getSubmitErrorMessage(error, translate('error')))
    }
  }

  const loadFilters = async (): Promise<void> => {
    try {
      const res = (await getProductFilters({
        productType: 1
      })) as ApiResponse<ProductFilterData>
      if (res.code === 0) {
        areas.value = res.data?.areas || []
        gpuModels.value = buildGpuResourceFilterOptions(res.data, translate)
        cpuModels.value = res.data?.cpuModels || []
      }
    } catch (error) {
      console.error('加载筛选条件失败', error)
    }
  }

  const loadImages = async (clusterId?: ResourceId | null): Promise<void> => {
    images.value = await fetchImagesByTab(activeTab.value, clusterId)
    if (!images.value.some((item) => item.id === selectedImage.value)) {
      selectedImage.value = ''
    }
  }

  const fetchImagesByTab = async (
    tab: NotebookImageTab,
    clusterId?: ResourceId | null
  ): Promise<ConsoleImage[]> => {
    const params: {
      page: number
      pageSize: number
      usageType: number
      type?: number
      clusterId?: ResourceId
    } = { page: 1, pageSize: 100, usageType: 1 }

    if (tab === 'base') params.type = 1
    if (tab === 'my') params.type = 2
    if (clusterId || selectedProduct.value?.clusterId) {
      params.clusterId = (clusterId || selectedProduct.value?.clusterId) as ResourceId
    }

    const res = (await getImageList(params)) as ApiResponse<ListData<ConsoleImage>>
    if (res.code !== 0) {
      return []
    }

    return decorateConsoleImages(res.data?.list || [])
  }

  const resolveEditImageTab = async (
    imageId: ResourceId | undefined,
    imageType: number | undefined,
    clusterId?: ResourceId | null
  ): Promise<NotebookImageTab> => {
    if (imageType === 1) return 'base'
    if (imageType === 2) return 'my'
    if (!imageId) return 'base'

    const tabs: NotebookImageTab[] = ['base', 'community', 'my']
    for (const tab of tabs) {
      const list = await fetchImagesByTab(tab, clusterId)
      if (list.some((item) => item.id === imageId)) {
        return tab
      }
    }

    return 'base'
  }

  const loadSSHKeys = async (): Promise<void> => {
    try {
      const res = (await getSSHKeyList({
        page: 1,
        pageSize: 100
      })) as ApiResponse<ListData<ConsoleSshKey>>
      if (res.code === 0) {
        sshKeys.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载SSH密钥失败', error)
    }
  }

  const loadVolumes = async (): Promise<void> => {
    if (!selectedProduct.value) return

    try {
      const res = (await getVolumeList({
        page: 1,
        pageSize: 100,
        clusterId: selectedProduct.value.clusterId
      })) as ApiResponse<ListData<ConsoleVolume>>

      if (res.code === 0) {
        availableVolumes.value = res.data?.list || []

        if (
          selectedVolumeId.value &&
          !availableVolumes.value.some(
            (item) => item.id === selectedVolumeId.value
          )
        ) {
          selectedVolumeId.value = null
        }
      }
    } catch (error) {
      console.error('加载数据盘失败', error)
    }
  }

  const selectProduct = (product: ConsoleProduct): void => {
    if (!isSelectableProduct(product)) return

    selectedProduct.value = product
    selectedVolumeId.value = null
    volumeMountPath.value = DEFAULT_VOLUME_MOUNT_PATH
    void loadImages(product.clusterId)
    void loadVolumes()
  }

  const changeFilter = <K extends keyof NotebookFilters>(
    key: K,
    value: NotebookFilters[K]
  ): void => {
    filters.value = {
      ...filters.value,
      [key]: value
    }
    void loadProducts()
  }

  const changeImageTab = (tab: NotebookImageTab): void => {
    activeTab.value = tab
    selectedImage.value = ''
    void loadProducts()
    void loadImages()
  }

  const onVolumeChange = (value: ResourceId | null): void => {
    selectedVolumeId.value = value
    if (value && !volumeMountPath.value) {
      volumeMountPath.value = DEFAULT_VOLUME_MOUNT_PATH
    }
  }

  const goBack = (): void => {
    router.go(-1)
  }

  const loadNotebookDetail = async (): Promise<ConsoleNotebookDetail | null> => {
    if (!notebookId.value) return null

    try {
      const res = (await getNotebookDetail({
        id: notebookId.value
      })) as ApiResponse<ConsoleNotebookDetail>
      if (res.code === 0) {
        return res.data || {}
      }
    } catch (error) {
      console.error('加载笔记本详情失败', error)
      ElMessage.error(getSubmitErrorMessage(error, translate('error')))
    }

    return null
  }

  const prefillEditForm = async (): Promise<void> => {
    const notebook = await loadNotebookDetail()
    if (!notebook) return

    instanceName.value = `${notebook.displayName || notebook.instanceName || ''}`
    payType.value = (Number(notebook.payType) || 1) as NotebookPayType
    selectedSshKey.value = (notebook.sshKeyId as ResourceId | undefined) || null
    enableSshPassword.value = Boolean(notebook.enableSshPassword)
    enableTensorboard.value = Boolean(notebook.enableTensorboard)
    tensorboardLogPath.value = `${notebook.tensorboardLogPath || ''}`

    const clusterId = (notebook.clusterId as ResourceId | undefined) || null
    await loadProducts(clusterId)

    const matchedProduct = products.value.find(
      (item) => item.id === notebook.productId
    )
    if (matchedProduct) {
      selectedProduct.value = matchedProduct
    }

    activeTab.value = await resolveEditImageTab(
      notebook.imageId as ResourceId | undefined,
      Number(notebook.imageType) || undefined,
      clusterId
    )
    await loadImages(clusterId)

    if (
      notebook.imageId &&
      images.value.some((item) => item.id === notebook.imageId)
    ) {
      selectedImage.value = notebook.imageId as ResourceId
    }
  }

  const handleSubmit = async (): Promise<void> => {
    const validationMessage = validateCreateForm()
    if (validationMessage) {
      ElMessage.warning(validationMessage)
      return
    }

    const product = selectedProduct.value
    const imageId = selectedImage.value
    if (!product || !imageId || !isSelectableProduct(product)) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    try {
      submitLoading.value = true

      if (isEditMode.value) {
        const params: NotebookUpdatePayload = {
          id: notebookId.value as ResourceId,
          displayName: instanceName.value.trim(),
          productId: product.id,
          imageId,
          chargeType: payType.value
        }

        const res = await updateNotebook(params)
        if (res.code === 0) {
          ElMessage.success(translate('success'))
          await router
            .push({ name: 'notebookDetail', query: { id: notebookId.value } })
            .catch(() => {
              return router.push({
                path: '/layout/notebooks/detail',
                query: { id: notebookId.value }
              })
            })
          return
        }

        const submitMessage = res.msg || translate('operationFailed')
        ElMessage.error(submitMessage)
        return
      }

      const volumeMounts: NotebookCreatePayload['volumeMounts'] = []
      if (selectedVolumeId.value) {
        volumeMounts.push({
          pvcId: selectedVolumeId.value,
          mountsPath: volumeMountPath.value || '/data/volume-1'
        })
      }

      const params: NotebookCreatePayload = {
        displayName: instanceName.value.trim(),
        productId: product.id,
        imageId,
        payType: payType.value,
        tensorBoard: enableTensorboard.value,
        tensorBoardPath: tensorboardLogPath.value,
        sshKeyId: selectedSshKey.value || 0,
        enableSshPassword: enableSshPassword.value,
        volumeMounts
      }

      const res = await createNotebook(params)
      if (res.code === 0) {
        ElMessage.success(translate('createSuccess'))
        router.go(-1)
        return
      }

      const submitMessage = res.msg || translate('createFailed')
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('instanceName', submitMessage)
      }
      ElMessage.error(submitMessage)
    } catch (error) {
      console.error('创建失败', error)
      const submitMessage = getSubmitErrorMessage(
        error,
        translate(isEditMode.value ? 'operationFailed' : 'createFailed')
      )
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('instanceName', submitMessage)
      }
      ElMessage.error(submitMessage)
    } finally {
      submitLoading.value = false
    }
  }

  onMounted(() => {
    void (async () => {
      await loadFilters()
      await loadSSHKeys()
      if (isEditMode.value) {
        await prefillEditForm()
        return
      }
      await loadProducts()
      await loadImages(selectedProduct.value?.clusterId)
    })()
  })

  watch(
    () => selectedProduct.value?.id ?? null,
    (value) => {
      if (!value && !selectedImageRecord.value?.clusterId) {
        void loadImages()
      }
    }
  )

  return {
    activeTab,
    areas,
    availableVolumes,
    canCreate,
    changeFilter,
    changeImageTab,
    cpuModels,
    enableSshPassword,
    enableTensorboard,
    fieldErrors,
    filteredImages,
    filters,
    formatPrice,
    goBack,
    gpuModels,
    handleSubmit,
    imageDescription,
    imageTabs,
    instanceName,
    isEditMode,
    updateInstanceName,
    updateTensorboardLogPath,
    onVolumeChange,
    pageTitleKey,
    payType,
    payTypes,
    priceUnitText,
    products,
    selectProduct,
    selectedImage,
    selectedProduct,
    selectedSshKey,
    selectedVolumeId,
    selectedVolumeName,
    sshKeys,
    submitLabelKey,
    submitLoading,
    tensorboardLogPath,
    totalPrice,
    validateInstanceNameField,
    validateTensorboardLogPathField,
    volumeMountPath
  }
}

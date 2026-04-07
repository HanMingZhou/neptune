import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { Router } from 'vue-router'
import { getProductFilters, getProductList } from '@/api/product'
import { getImageList } from '@/api/image'
import { createNotebook } from '@/api/notebook'
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
import type {
  ConsoleImage,
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

interface UseNotebookCreateOptions {
  t?: Translator
  router: Pick<Router, 'go'>
}

export const useNotebookCreate = ({ t, router }: UseNotebookCreateOptions) => {
  const translate: Translator = (key, params) =>
    typeof t === 'function' ? t(key, params) : key

  const products = ref<ConsoleProduct[]>([])
  const selectedProduct = ref<ConsoleProduct | null>(null)
  const gpuCount = ref(1)
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

    const count = selectedProduct.value.gpuCount > 0 ? gpuCount.value : 1
    const total = getUnitPrice(selectedProduct.value) * count
    return payType.value === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const canCreate = computed(() =>
    Boolean(
      selectedProduct.value && selectedImage.value && instanceName.value.trim()
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
    const error = enableTensorboard.value
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

    if (!isTensorboardPathValid) {
      return fieldErrors.value.tensorboardLogPath
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

    if (!selectedProduct.value) {
      selectProduct(list[0])
      return
    }

    const nextSelected = list.find(
      (item) => item.id === selectedProduct.value?.id
    )
    if (!nextSelected) {
      selectProduct(list[0])
      return
    }

    selectedProduct.value = nextSelected
  }

  const loadProducts = async (): Promise<void> => {
    try {
      const params: {
        page: number
        pageSize: number
        productType: number
        area?: string
        gpuModel?: string
        cpuModel?: string
      } = { page: 1, pageSize: 100, productType: 1 }

      if (filters.value.area) params.area = filters.value.area
      if (filters.value.gpuModel) params.gpuModel = filters.value.gpuModel
      if (filters.value.cpuModel) params.cpuModel = filters.value.cpuModel

      const res = (await getProductList(params)) as ApiResponse<
        ListData<ConsoleProduct>
      >
      if (res.code === 0) {
        const list = res.data?.list || []
        products.value = list
        syncSelectedProduct(list)
      }
    } catch (error) {
      console.error('加载产品失败', error)
      ElMessage.error(translate('error'))
    }
  }

  const loadFilters = async (): Promise<void> => {
    try {
      const res = (await getProductFilters({
        productType: 1
      })) as ApiResponse<ProductFilterData>
      if (res.code === 0) {
        areas.value = res.data?.areas || []
        gpuModels.value = res.data?.gpuModels || []
        cpuModels.value = res.data?.cpuModels || []
      }
    } catch (error) {
      console.error('加载筛选条件失败', error)
    }
  }

  const loadImages = async (): Promise<void> => {
    try {
      const params: {
        page: number
        pageSize: number
        usageType: number
        type?: number
      } = { page: 1, pageSize: 100, usageType: 1 }

      if (activeTab.value === 'base') params.type = 1
      if (activeTab.value === 'my') params.type = 2

      const res = (await getImageList(params)) as ApiResponse<
        ListData<ConsoleImage>
      >
      if (res.code === 0) {
        images.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载镜像失败', error)
    }
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
    selectedProduct.value = product
    gpuCount.value = 1
    selectedVolumeId.value = null
    volumeMountPath.value = DEFAULT_VOLUME_MOUNT_PATH
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

  const handleCreate = async (): Promise<void> => {
    const validationMessage = validateCreateForm()
    if (validationMessage) {
      ElMessage.warning(validationMessage)
      return
    }

    const product = selectedProduct.value
    const imageId = selectedImage.value
    if (!product || !imageId) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    try {
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
        translate('createFailed')
      )
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('instanceName', submitMessage)
      }
      ElMessage.error(submitMessage)
    }
  }

  onMounted(() => {
    void loadProducts()
    void loadFilters()
    void loadImages()
    void loadSSHKeys()
  })

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
    gpuCount,
    gpuModels,
    handleCreate,
    imageDescription,
    imageTabs,
    instanceName,
    updateInstanceName,
    updateTensorboardLogPath,
    onVolumeChange,
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
    tensorboardLogPath,
    totalPrice,
    validateInstanceNameField,
    validateTensorboardLogPathField,
    volumeMountPath
  }
}

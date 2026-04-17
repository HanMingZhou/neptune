import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { Router } from 'vue-router'
import { createTrainingJob } from '@/api/training'
import { getImageList } from '@/api/image'
import { getAggregateProductList, getProductFilters } from '@/api/product'
import { getVolumeList } from '@/api/volume'
import type { ApiResponse } from '@/utils/request'
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
import { getVGpuNumber } from '@/utils/vgpu'
import type {
  ConsoleImage,
  ConsoleProduct,
  ConsoleVolume,
  ListData,
  ProductFilterData,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import {
  TRAINING_FRAMEWORK_TYPES,
  TRAINING_IMAGE_TABS,
  TRAINING_PAY_TYPES
} from '../constants'

export type TrainingFrameworkType =
  (typeof TRAINING_FRAMEWORK_TYPES)[number]['value']
export type TrainingImageTab = (typeof TRAINING_IMAGE_TABS)[number]['value']
export type TrainingPayType = (typeof TRAINING_PAY_TYPES)[number]['value']
export type ScheduleStrategy = 'BALANCED' | 'STRICT'

type TrainingFormKey = keyof TrainingCreateForm

export type UpdateTrainingFieldPayload<
  K extends TrainingFormKey = TrainingFormKey
> = {
  key: K
  value: TrainingCreateForm[K]
}

export interface TrainingMount {
  mountType: string
  pvcId: ResourceId | null
  mountPath: string
  readOnly: boolean
}

export interface TrainingEnv {
  name: string
  value: string
}

export interface TrainingCreateForm {
  name: string
  frameworkType: TrainingFrameworkType
  imageId: ResourceId | null
  startupCommand: string
  workerCount: number
  scheduleStrategy: ScheduleStrategy
  gpuPerWorker: number
  gpuType: string
  useSHM: boolean
  shmSize: string
  enableTensorboard: boolean
  tensorboardLogPath: string
  payType: TrainingPayType
  mounts: TrainingMount[]
  envs: TrainingEnv[]
}

export interface TrainingFilters {
  area: string
  gpuModel: string
  cpuModel: string
}

export interface TrainingFieldErrors {
  name: string
  tensorboardLogPath: string
}

export interface TrainingTotalResources {
  gpu: number
  nodes: number
}

interface TrainingCreatePayload {
  name: string
  frameworkType: TrainingFrameworkType
  imageId: ResourceId
  startupCommand: string
  clusterId: ResourceId
  resourceId: ResourceId
  templateProductId: ResourceId
  instanceCount: number
  workerCount: number
  scheduleStrategy: ScheduleStrategy
  enableTensorboard: boolean
  tensorboardLogPath: string
  payType: TrainingPayType
  mounts: Array<{
    mountType: string
    pvcId: ResourceId
    mountPath: string
    readOnly: boolean
  }>
  envs: TrainingEnv[]
}

interface UseTrainingCreateOptions {
  t?: Translator
  router: Pick<Router, 'go'>
}

const MULTI_WORKER_FRAMEWORKS: TrainingFrameworkType[] = ['PYTORCH_DDP', 'MPI']

export const useTrainingCreate = ({ t, router }: UseTrainingCreateOptions) => {
  const translate: Translator = (key, params) =>
    typeof t === 'function' ? t(key, params) : key

  const form = reactive<TrainingCreateForm>({
    name: '',
    frameworkType: 'STANDALONE',
    imageId: null,
    startupCommand: '',
    workerCount: 2,
    scheduleStrategy: 'BALANCED',
    gpuPerWorker: 1,
    gpuType: 'nvidia.com/gpu',
    useSHM: true,
    shmSize: '4Gi',
    enableTensorboard: false,
    tensorboardLogPath: 'logs',
    payType: 1,
    mounts: [],
    envs: []
  })

  const filters = reactive<TrainingFilters>({
    area: '',
    gpuModel: '',
    cpuModel: ''
  })

  const images = ref<ConsoleImage[]>([])
  const pvcs = ref<ConsoleVolume[]>([])
  const products = ref<ConsoleProduct[]>([])
  const selectedProduct = ref<ConsoleProduct | null>(null)
  const areas = ref<string[]>([])
  const gpuModels = ref<NonNullable<ProductFilterData['gpuModels']>>([])
  const cpuModels = ref<NonNullable<ProductFilterData['cpuModels']>>([])
  const creating = ref(false)
  const fieldErrors = reactive<TrainingFieldErrors>({
    name: '',
    tensorboardLogPath: ''
  })

  const activeTab = ref<TrainingImageTab>('base')
  const payTypes = TRAINING_PAY_TYPES

  const frameworkTypes = computed(() =>
    TRAINING_FRAMEWORK_TYPES.map((item) => ({
      ...item,
      hint: translate(item.hintKey)
    }))
  )

  const imageTabs = computed(() =>
    TRAINING_IMAGE_TABS.map((tab) => ({
      ...tab,
      label: translate(tab.labelKey)
    }))
  )

  const selectedImage = computed(
    () => images.value.find((item) => item.id === form.imageId) || null
  )
  const filteredImages = computed(() => images.value)
  const showWorkerCount = computed(() =>
    MULTI_WORKER_FRAMEWORKS.includes(form.frameworkType)
  )

  const totalResources = computed<TrainingTotalResources | null>(() => {
    if (!selectedProduct.value) return null

    const nodes = showWorkerCount.value ? form.workerCount : 1
    const gpuPerNode =
      selectedProduct.value.gpuCount > 0
        ? form.gpuPerWorker
        : getVGpuNumber(selectedProduct.value)
    const gpu = gpuPerNode * nodes
    return { gpu, nodes }
  })

  const priceUnitText = computed(() => {
    switch (form.payType) {
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

  const getUnitPrice = (product: ConsoleProduct | null | undefined): number => {
    if (!product) return 0

    switch (form.payType) {
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
    return form.payType === 1 ? price.toFixed(4) : price.toFixed(2)
  }

  const totalPrice = computed(() => {
    if (!selectedProduct.value) return '0.00'

    const nodes = showWorkerCount.value ? form.workerCount : 1
    const total = getUnitPrice(selectedProduct.value) * nodes
    return form.payType === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const availableCapacity = computed(
    () => selectedProduct.value?.available ?? 0
  )

  const totalAllowedNodes = computed(() => {
    if (!selectedProduct.value) return 0
    if (!showWorkerCount.value) {
      return availableCapacity.value
    }

    const strictMax = selectedProduct.value.strictMax ?? 0
    const balancedMax = selectedProduct.value.balancedMax ?? 0
    if (form.scheduleStrategy === 'STRICT') {
      return strictMax > 0 ? strictMax : availableCapacity.value
    }
    return balancedMax > 0 ? balancedMax : availableCapacity.value
  })

  const maxWorkerCount = computed(() => {
    if (!selectedProduct.value) return 100
    return Math.max(2, totalAllowedNodes.value)
  })

  const canCreate = computed(() => {
    if (
      !form.name.trim() ||
      !form.imageId ||
      !form.startupCommand.trim() ||
      !selectedProduct.value
    ) {
      return false
    }

    const requiredNodes = showWorkerCount.value ? form.workerCount : 1
    if (requiredNodes > totalAllowedNodes.value) return false
    if (showWorkerCount.value && form.workerCount < 2) return false
    return true
  })

  const loadImages = async (clusterId?: ResourceId | null): Promise<void> => {
    try {
      const params: {
        page: number
        pageSize: number
        usageType: number
        type?: number
        clusterId?: ResourceId
      } = { page: 1, pageSize: 100, usageType: 2 }

      if (activeTab.value === 'base') params.type = 1
      if (activeTab.value === 'my') params.type = 2
      if (clusterId || selectedProduct.value?.clusterId) {
        params.clusterId = (clusterId ||
          selectedProduct.value?.clusterId) as ResourceId
      }

      const res = (await getImageList(params)) as ApiResponse<
        ListData<ConsoleImage>
      >
      if (res.code === 0) {
        images.value = decorateConsoleImages(res.data?.list || [])
        if (!images.value.some((item) => item.id === form.imageId)) {
          form.imageId = null
        }
      }
    } catch (error) {
      console.error('加载镜像失败', error)
    }
  }

  const changeImageTab = (tab: TrainingImageTab): void => {
    activeTab.value = tab
    form.imageId = null
    void loadProducts()
    void loadImages()
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
      } = {
        page: 1,
        pageSize: 100,
        productType: 1
      }

      if (filters.area) params.area = filters.area
      if (filters.gpuModel) {
        const selectedGpuFilter = findGpuResourceFilterOption(
          gpuModels.value,
          filters.gpuModel
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
          params.gpuModel = filters.gpuModel
        }
      }
      if (filters.cpuModel) params.cpuModel = filters.cpuModel
      if (clusterId || selectedImage.value?.clusterId) {
        params.clusterId = (clusterId ||
          selectedImage.value?.clusterId) as ResourceId
      }

      const res = (await getAggregateProductList(params)) as ApiResponse<
        ListData<ConsoleProduct>
      >
      if (res.code === 0) {
        products.value = res.data?.list || []
        if (
          selectedProduct.value &&
          !products.value.find((item) => item.id === selectedProduct.value?.id)
        ) {
          selectedProduct.value = null
        }
      }
    } catch (error) {
      console.error('加载产品失败', error)
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

  const loadPvcs = async (): Promise<void> => {
    try {
      const res = (await getVolumeList({
        page: 1,
        pageSize: 100
      })) as ApiResponse<ListData<ConsoleVolume>>
      if (res.code === 0) {
        pvcs.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载 PVC 失败', error)
    }
  }

  const changeFilter = <K extends keyof TrainingFilters>(
    key: K,
    value: TrainingFilters[K]
  ): void => {
    filters[key] = value
    void loadProducts()
  }

  const selectProduct = (product: ConsoleProduct): void => {
    selectedProduct.value = product
    if (product.gpuCount > 0) {
      form.gpuPerWorker = 1
      form.gpuType = product.gpuModel || ''
    } else {
      form.gpuPerWorker = 0
      form.gpuType = ''
    }
    void loadImages(product.clusterId)
  }

  const decreaseWorker = (): void => {
    if (form.workerCount > 2) form.workerCount -= 1
  }

  const increaseWorker = (): void => {
    if (form.workerCount < maxWorkerCount.value) form.workerCount += 1
  }

  const addMount = (): void => {
    form.mounts.push({
      mountType: 'DATASET',
      pvcId: null,
      mountPath: '',
      readOnly: false
    })
  }

  const removeMount = (index: number): void => {
    form.mounts.splice(index, 1)
  }

  const addEnv = (): void => {
    form.envs.push({ name: '', value: '' })
  }

  const removeEnv = (index: number): void => {
    form.envs.splice(index, 1)
  }

  const updateFieldError = <K extends keyof TrainingFieldErrors>(
    field: K,
    message = ''
  ): void => {
    fieldErrors[field] = message
  }

  const updateField = <K extends TrainingFormKey>({
    key,
    value
  }: UpdateTrainingFieldPayload<K>): void => {
    form[key] = value

    if (key === 'name' && fieldErrors.name) {
      updateFieldError('name')
    }

    if (key === 'tensorboardLogPath' && fieldErrors.tensorboardLogPath) {
      updateFieldError('tensorboardLogPath')
    }

    if (
      key === 'enableTensorboard' &&
      !value &&
      fieldErrors.tensorboardLogPath
    ) {
      updateFieldError('tensorboardLogPath')
    }
  }

  const validateNameField = (): boolean => {
    const error = validateK8sResourceName(form.name, {
      t: translate,
      fieldKey: 'jobName',
      example: 'my-training'
    })

    updateFieldError('name', error || '')
    return !error
  }

  const validateTensorboardLogPathField = (): boolean => {
    const error = form.enableTensorboard
      ? validateTensorBoardPath(form.tensorboardLogPath, translate)
      : null

    updateFieldError('tensorboardLogPath', error || '')
    return !error
  }

  const validateCreateForm = (): string => {
    const isNameValid = validateNameField()
    const isTensorboardPathValid = validateTensorboardLogPathField()

    if (!isNameValid) {
      return fieldErrors.name
    }

    if (
      !form.imageId ||
      !form.startupCommand.trim() ||
      !selectedProduct.value
    ) {
      return translate('fillAllFields')
    }

    if (!isTensorboardPathValid) {
      return fieldErrors.tensorboardLogPath
    }

    const requiredNodes = showWorkerCount.value ? form.workerCount : 1
    if (
      requiredNodes > totalAllowedNodes.value ||
      (showWorkerCount.value && form.workerCount < 2)
    ) {
      return translate('fillAllFields')
    }

    return ''
  }

  const insertMpiExample = (): void => {
    const totalSlots = form.workerCount * form.gpuPerWorker
    const np = totalSlots > 0 ? totalSlots : form.workerCount
    form.startupCommand = `mpirun -np ${np} -H ${'${MPI_HOSTS}'} --allow-run-as-root --bind-to none -map-by slot -x NCCL_DEBUG=INFO python train.py`
  }

  const insertPytorchExample = (): void => {
    const nprocPerNode = form.gpuPerWorker > 0 ? form.gpuPerWorker : 1
    form.startupCommand = `torchrun --nnodes=${'${WORLD_SIZE}'} --node_rank=${'${RANK}'} --master_addr=${'${MASTER_ADDR}'} --master_port=${'${MASTER_PORT}'} --nproc_per_node=${nprocPerNode} train.py`
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
    const imageId = form.imageId
    if (!product || !imageId) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    creating.value = true
    try {
      const params: TrainingCreatePayload = {
        name: form.name.trim(),
        frameworkType: form.frameworkType,
        imageId,
        startupCommand: form.startupCommand,
        clusterId: product.clusterId,
        resourceId: product.id,
        templateProductId: product.templateProductId || product.id,
        instanceCount: showWorkerCount.value ? form.workerCount : 1,
        workerCount: showWorkerCount.value ? form.workerCount : 0,
        scheduleStrategy: showWorkerCount.value
          ? form.scheduleStrategy
          : 'BALANCED',
        enableTensorboard: form.enableTensorboard,
        tensorboardLogPath: form.tensorboardLogPath,
        payType: form.payType,
        mounts: form.mounts
          .filter((mount) => mount.pvcId && mount.mountPath)
          .map((mount) => ({
            mountType: mount.mountType,
            pvcId: mount.pvcId as ResourceId,
            mountPath: mount.mountPath,
            readOnly: mount.readOnly
          })),
        envs: form.envs.filter((env) => env.name && env.value)
      }

      const res = await createTrainingJob(params)
      if (res.code === 0) {
        ElMessage.success(translate('createSuccess'))
        router.go(-1)
        return
      }

      const submitMessage = res.msg || translate('createFailed')
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('name', submitMessage)
      }
      ElMessage.error(submitMessage)
    } catch (error) {
      console.error('创建失败', error)
      const submitMessage = getSubmitErrorMessage(
        error,
        translate('createFailed')
      )
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('name', submitMessage)
      }
      ElMessage.error(submitMessage)
    } finally {
      creating.value = false
    }
  }

  onMounted(() => {
    void loadImages()
    void loadPvcs()
    void loadFilters()
    void loadProducts()
  })

  watch(
    () => form.imageId,
    (value) => {
      if (!value || !selectedImage.value?.clusterId) {
        return
      }

      void loadProducts(selectedImage.value.clusterId)
    }
  )

  watch(
    () => selectedProduct.value?.id ?? null,
    (value) => {
      if (!value && !selectedImage.value?.clusterId) {
        void loadImages()
      }
    }
  )

  return {
    activeTab,
    addEnv,
    addMount,
    areas,
    availableCapacity,
    canCreate,
    changeFilter,
    changeImageTab,
    cpuModels,
    creating,
    decreaseWorker,
    fieldErrors,
    filteredImages,
    filters,
    form,
    formatPrice,
    frameworkTypes,
    goBack,
    gpuModels,
    handleCreate,
    imageTabs,
    increaseWorker,
    insertMpiExample,
    insertPytorchExample,
    maxWorkerCount,
    payTypes,
    priceUnitText,
    products,
    pvcs,
    removeEnv,
    removeMount,
    selectProduct,
    selectedProduct,
    showWorkerCount,
    totalAllowedNodes,
    totalPrice,
    totalResources,
    updateField,
    validateNameField,
    validateTensorboardLogPathField
  }
}

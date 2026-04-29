import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { Router } from 'vue-router'
import { useRoute } from 'vue-router'
import {
  createInferenceService,
  getInferenceServiceDetail,
  updateInferenceService
} from '@/api/inference'
import { getImageList } from '@/api/image'
import { getAggregateProductList, getProductFilters } from '@/api/product'
import { getVolumeList } from '@/api/volume'
import type { ApiResponse } from '@/utils/request'
import {
  getSubmitErrorMessage,
  isResourceNameErrorMessage,
  validateK8sResourceName
} from '@/utils/resourceValidators'
import {
  buildGpuResourceFilterOptions,
  findGpuResourceFilterOption
} from '@/utils/gpuFilters'
import { decorateConsoleImages } from '@/utils/imageRegistry'
import type {
  ConsoleImage,
  ConsoleInferenceDetail,
  ConsoleProduct,
  ConsoleVolume,
  FilterOption,
  ListData,
  ProductFilterData,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import {
  INFERENCE_AUTH_TYPES,
  INFERENCE_IMAGE_TABS,
  INFERENCE_PAY_TYPES
} from '../constants'

export type InferenceAuthType = (typeof INFERENCE_AUTH_TYPES)[number]['value']
export type InferenceImageTab = (typeof INFERENCE_IMAGE_TABS)[number]['value']
export type InferencePayType = (typeof INFERENCE_PAY_TYPES)[number]['value']
export type InferenceFramework = '' | 'VLLM' | 'SGLANG'
export type InferenceDeployType = 'STANDALONE' | 'DISTRIBUTED'
export type InferenceScheduleStrategy = 'BALANCED' | 'STRICT'

type InferenceFormKey = keyof InferenceCreateForm

export type UpdateInferenceFieldPayload<
  K extends InferenceFormKey = InferenceFormKey
> = {
  key: K
  value: InferenceCreateForm[K]
}

export interface InferenceMount {
  pvcId: ResourceId | null
  mountPath: string
  subPath: string
  readOnly: boolean
}

export interface InferenceEnv {
  name: string
  value: string
}

export interface InferenceCreateForm {
  displayName: string
  framework: InferenceFramework
  deployType: InferenceDeployType
  modelPvcId: ResourceId | ''
  modelMountPath: string
  imageId: ResourceId | ''
  productId: ResourceId | ''
  payType: InferencePayType
  workerCount: number
  scheduleStrategy: InferenceScheduleStrategy
  autoRestart: boolean
  maxRestarts: number
  servicePort: number
  authType: InferenceAuthType
  command: string
  args: string
  mounts: InferenceMount[]
  envs: InferenceEnv[]
}

export interface InferenceFilters {
  area: string
  gpuModel: string
}

export interface InferenceFieldErrors {
  displayName: string
}

interface InferenceCreatePayload {
  id?: ResourceId
  name: string
  framework?: Exclude<InferenceFramework, ''>
  deployType: InferenceDeployType
  modelPvcId: ResourceId
  modelMountPath: string
  imageId: ResourceId
  productId: ResourceId
  templateProductId: ResourceId
  instanceCount: number
  scheduleStrategy: InferenceScheduleStrategy
  payType: InferencePayType
  servicePort: number
  authType: InferenceAuthType
  command: string
  args: string[]
  mounts: Array<{
    pvcId: ResourceId
    mountPath: string
    subPath?: string
    readOnly: boolean
  }>
  envs: InferenceEnv[]
  workerCount?: number
  autoRestart?: boolean
  maxRestarts?: number
}

interface UseInferenceCreateOptions {
  t?: Translator
  router: Pick<Router, 'push'>
}

export const useInferenceCreate = ({
  t,
  router
}: UseInferenceCreateOptions) => {
  const route = useRoute()
  const translate: Translator = (key, params) =>
    typeof t === 'function' ? t(key, params) : key

  const editServiceId = computed(() => Number(route.query.id) || 0)
  const isEditMode = computed(() => editServiceId.value > 0)
  const editLoading = ref(false)
  const pageTitleKey = computed(() =>
    isEditMode.value ? 'inference.editTitle' : 'inference.createTitle'
  )
  const submitLabelKey = computed(() =>
    isEditMode.value ? 'save' : 'inference.createService'
  )

  const frameworks = computed(() => [
    { label: 'vLLM', value: 'VLLM', icon: 'auto_awesome' },
    { label: 'SGLang', value: 'SGLANG', icon: 'psychology' }
  ])

  const deployTypes = computed(() => [
    { label: translate('standalone'), value: 'STANDALONE' },
    { label: translate('distributed'), value: 'DISTRIBUTED' }
  ])

  const authTypes = INFERENCE_AUTH_TYPES
  const payTypes = INFERENCE_PAY_TYPES

  const form = reactive<InferenceCreateForm>({
    displayName: '',
    framework: '',
    deployType: 'STANDALONE',
    modelPvcId: '',
    modelMountPath: '/model',
    imageId: '',
    productId: '',
    payType: 1,
    workerCount: 2,
    scheduleStrategy: 'BALANCED',
    autoRestart: false,
    maxRestarts: 3,
    servicePort: 30000,
    authType: 1,
    command: '',
    args: '',
    mounts: [],
    envs: []
  })

  const filters = reactive<InferenceFilters>({
    area: '',
    gpuModel: ''
  })

  const images = ref<ConsoleImage[]>([])
  const pvcs = ref<ConsoleVolume[]>([])
  const products = ref<ConsoleProduct[]>([])
  const areas = ref<string[]>([])
  const gpuModelsList = ref<FilterOption[]>([])
  const activeTab = ref<InferenceImageTab>('base')
  const loading = ref(false)
  const fieldErrors = reactive<InferenceFieldErrors>({
    displayName: ''
  })

  const imageTabs = computed(() =>
    INFERENCE_IMAGE_TABS.map((tab) => ({
      ...tab,
      label: translate(tab.labelKey)
    }))
  )

  const selectedImage = computed(
    () => images.value.find((item) => item.id === form.imageId) || null
  )

  const imageOptions = computed(() =>
    images.value.map((image) => ({
      label: image.name,
      value: image.id,
      desc: image.description
    }))
  )

  const hasCommand = computed(() => Boolean(form.command.trim()))
  const frameworkRequired = computed(() => form.deployType === 'DISTRIBUTED')
  const selectedProduct = computed(
    () => products.value.find((item) => item.id === form.productId) || null
  )
  const isSelectableProduct = (
    product: ConsoleProduct | null | undefined
  ): product is ConsoleProduct =>
    Boolean(product && (isEditMode.value || (product.available ?? 0) > 0))
  const maxDistributedCount = computed(() => {
    const product = selectedProduct.value
    if (!product) return 0
    if (form.scheduleStrategy === 'STRICT') {
      return product.strictMax && product.strictMax > 0
        ? product.strictMax
        : (product.available ?? 0)
    }
    return product.balancedMax && product.balancedMax > 0
      ? product.balancedMax
      : (product.available ?? 0)
  })

  const canCreate = computed(() => {
    if (
      !form.displayName.trim() ||
      !form.modelPvcId ||
      !form.imageId ||
      !form.productId
    )
      return false
    if (!isSelectableProduct(selectedProduct.value)) return false
    if (!hasCommand.value) return false
    if (frameworkRequired.value && !form.framework) return false
    if (form.deployType === 'DISTRIBUTED' && form.workerCount < 2) return false
    if (
      form.deployType === 'DISTRIBUTED' &&
      form.workerCount > maxDistributedCount.value
    )
      return false
    return true
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

  const totalPrice = computed(() => {
    const product = selectedProduct.value
    if (!product) return '0.00'

    let total = getUnitPrice(product)
    if (form.deployType === 'DISTRIBUTED') {
      total *= form.workerCount
    }

    return form.payType === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const formatPrice = (product: ConsoleProduct | null | undefined): string =>
    form.payType === 1
      ? getUnitPrice(product).toFixed(4)
      : getUnitPrice(product).toFixed(2)

  const onDeployTypeChange = (value: InferenceDeployType): void => {
    form.deployType = value
    if (value === 'STANDALONE') {
      form.scheduleStrategy = 'BALANCED'
      return
    }

    if (!form.framework) {
      form.framework = 'VLLM'
    }
  }

  const loadImages = async (clusterId?: ResourceId | null): Promise<void> => {
    try {
      const params: {
        page: number
        pageSize: number
        usageType: number
        type?: number
        clusterId?: ResourceId
      } = { page: 1, pageSize: 100, usageType: 3 }

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
          form.imageId = ''
        }
      }
    } catch (error) {
      console.error('加载镜像失败', error)
    }
  }

  const changeImageTab = (tab: InferenceImageTab): void => {
    activeTab.value = tab
    form.imageId = ''
    void loadProducts()
    void loadImages()
  }

  const loadProducts = async (
    clusterId?: ResourceId | null,
    silent = false
  ): Promise<void> => {
    if (!silent) loading.value = true

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
        clusterId?: ResourceId
      } = { page: 1, pageSize: 100, productType: 1 }

      if (filters.area) params.area = filters.area
      if (filters.gpuModel) {
        const selectedGpuFilter = findGpuResourceFilterOption(
          gpuModelsList.value,
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
      if (clusterId || selectedImage.value?.clusterId) {
        params.clusterId = (clusterId ||
          selectedImage.value?.clusterId) as ResourceId
      }

      const res = (await getAggregateProductList(params)) as ApiResponse<
        ListData<ConsoleProduct>
      >
      if (res.code === 0) {
        products.value = res.data?.list || []
        const nextSelected = products.value.find(
          (item) => item.id === form.productId && isSelectableProduct(item)
        )
        if (nextSelected) {
          form.productId = nextSelected.id
          return
        }

        form.productId = products.value.find(isSelectableProduct)?.id ?? ''
      }
    } catch (error) {
      console.error(error)
    } finally {
      if (!silent) loading.value = false
    }
  }

  const loadFilters = async (): Promise<void> => {
    try {
      const res = (await getProductFilters({
        productType: 1
      })) as ApiResponse<ProductFilterData>
      if (res.code === 0) {
        areas.value = res.data?.areas || []
        gpuModelsList.value = buildGpuResourceFilterOptions(res.data, translate)
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

  const changeFilter = <K extends keyof InferenceFilters>(
    key: K,
    value: InferenceFilters[K]
  ): void => {
    filters[key] = value
    void loadProducts()
  }

  const addMount = (): void => {
    form.mounts.push({
      pvcId: null,
      mountPath: '',
      subPath: '',
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

  const updateFieldError = (
    field: keyof InferenceFieldErrors,
    message = ''
  ): void => {
    fieldErrors[field] = message
  }

  const updateField = <K extends InferenceFormKey>({
    key,
    value
  }: UpdateInferenceFieldPayload<K>): void => {
    form[key] = value

    if (key === 'displayName' && fieldErrors.displayName) {
      updateFieldError('displayName')
    }
  }

  const validateDisplayNameField = (): boolean => {
    const error = validateK8sResourceName(form.displayName, {
      t: translate,
      fieldKey: 'name',
      example: 'my-service'
    })

    updateFieldError('displayName', error || '')
    return !error
  }

  const validateCreateForm = (): string => {
    const isDisplayNameValid = validateDisplayNameField()
    if (!isDisplayNameValid) {
      return fieldErrors.displayName
    }

    if (
      !form.modelPvcId ||
      !form.imageId ||
      !form.productId ||
      !hasCommand.value
    ) {
      return translate('fillAllFields')
    }
    if (!isSelectableProduct(selectedProduct.value)) {
      return translate('fillAllFields')
    }

    if (frameworkRequired.value && !form.framework) {
      return translate('fillAllFields')
    }

    if (
      form.deployType === 'DISTRIBUTED' &&
      (form.workerCount < 2 || form.workerCount > maxDistributedCount.value)
    ) {
      return translate('fillAllFields')
    }

    return ''
  }

  const handleCancel = (): void => {
    router.push({ name: 'inference' })
  }

  const handleSubmit = async (): Promise<void> => {
    const validationMessage = validateCreateForm()
    if (validationMessage) {
      ElMessage.warning(validationMessage)
      return
    }

    const modelPvcId = form.modelPvcId
    const imageId = form.imageId
    const productId = form.productId
    if (!modelPvcId || !imageId || !productId || !isSelectableProduct(selectedProduct.value)) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    loading.value = true
    try {
      const params: InferenceCreatePayload = {
        id: isEditMode.value ? editServiceId.value : undefined,
        name: form.displayName.trim(),
        framework: form.framework || undefined,
        deployType: form.deployType,
        modelPvcId,
        modelMountPath: form.modelMountPath || '/model',
        imageId,
        productId,
        templateProductId: productId,
        instanceCount: form.deployType === 'DISTRIBUTED' ? form.workerCount : 1,
        scheduleStrategy:
          form.deployType === 'DISTRIBUTED'
            ? form.scheduleStrategy
            : 'BALANCED',
        payType: form.payType,
        servicePort: form.servicePort,
        authType: form.authType,
        command: form.command.trim(),
        args: form.args.trim()
          ? form.args
              .trim()
              .split(/\n+/)
              .map((item) => item.trim())
              .filter(Boolean)
          : [],
        mounts: form.mounts
          .filter((mount) => mount.pvcId && mount.mountPath)
          .map((mount) => ({
            pvcId: mount.pvcId as ResourceId,
            mountPath: mount.mountPath,
            subPath: mount.subPath || undefined,
            readOnly: mount.readOnly
          })),
        envs: form.envs
          .filter((env) => env.name && env.value)
          .map((env) => ({
            name: env.name,
            value: env.value
          }))
      }

      if (form.deployType === 'DISTRIBUTED') {
        params.workerCount = form.workerCount
        params.autoRestart = form.autoRestart
        params.maxRestarts = form.autoRestart ? form.maxRestarts : 0
      }

      const res = isEditMode.value
        ? await updateInferenceService(params)
        : await createInferenceService(params)
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        router.push({ name: 'inference' })
        return
      }

      const submitMessage =
        res.msg || translate(isEditMode.value ? 'error' : 'createFailed')
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('displayName', submitMessage)
      }
      ElMessage.error(submitMessage)
    } catch (error) {
      console.error(isEditMode.value ? '更新失败' : '创建失败', error)
      const submitMessage = getSubmitErrorMessage(
        error,
        translate(isEditMode.value ? 'error' : 'createFailed')
      )
      if (isResourceNameErrorMessage(submitMessage)) {
        updateFieldError('displayName', submitMessage)
      }
      ElMessage.error(submitMessage)
    } finally {
      loading.value = false
    }
  }

  const applyInferenceDetail = async (
    detail: ConsoleInferenceDetail
  ): Promise<void> => {
    if (`${detail.status || ''}`.toUpperCase() !== 'STOPPED') {
      ElMessage.warning(translate('onlyStoppedCanEditInference'))
      router.push({ name: 'inference' })
      return
    }

    form.displayName = detail.displayName || detail.instanceName || ''
    form.framework = ((detail.framework || '') as InferenceFramework) || ''
    form.deployType =
      (detail.deployType as InferenceDeployType) || 'STANDALONE'
    form.modelPvcId = (detail.modelPvcId as ResourceId) || ''
    form.modelMountPath = detail.modelMountPath || '/model'
    form.imageId = (detail.imageId as ResourceId) || ''
    form.productId = (detail.productId as ResourceId) || ''
    form.payType = (Number(detail.payType) || 1) as InferencePayType
    form.workerCount = Math.max(
      2,
      Number(detail.instanceCount ?? detail.workerCount) || 2
    )
    form.scheduleStrategy =
      (detail.scheduleStrategy as InferenceScheduleStrategy) || 'BALANCED'
    form.autoRestart = Boolean(detail.autoRestart)
    form.maxRestarts = Number(detail.maxRestarts) || 3
    form.servicePort = Number(detail.servicePort) || 30000
    form.authType = (Number(detail.authType) || 1) as InferenceAuthType
    form.command = detail.command || ''
    form.args = Array.isArray(detail.args)
      ? detail.args.join('\n')
      : `${detail.args ?? ''}`
    form.mounts = (detail.mounts || [])
      .filter((mount) => mount.pvcId)
      .map((mount) => ({
        pvcId: (mount.pvcId as ResourceId) || null,
        mountPath: mount.mountPath || '',
        subPath: mount.subPath || '',
        readOnly: Boolean(mount.readOnly)
      }))
    form.envs = (detail.envs || [])
      .filter((env) => env.name)
      .map((env) => ({ name: env.name || '', value: `${env.value ?? ''}` }))

    await loadProducts(detail.clusterId || null, true)
    await loadImages(detail.clusterId || null)
  }

  const loadEditDetail = async (): Promise<void> => {
    if (!isEditMode.value) return
    editLoading.value = true
    try {
      const res = (await getInferenceServiceDetail({
        id: editServiceId.value
      })) as ApiResponse<ConsoleInferenceDetail>
      if (res.code === 0 && res.data) {
        await applyInferenceDetail(res.data)
      } else {
        ElMessage.error(res.msg || translate('error'))
        router.push({ name: 'inference' })
      }
    } catch (error) {
      console.error('加载推理服务详情失败', error)
      ElMessage.error(getSubmitErrorMessage(error, translate('error')))
      router.push({ name: 'inference' })
    } finally {
      editLoading.value = false
    }
  }

  onMounted(async () => {
    void loadPvcs()
    void loadFilters()
    if (isEditMode.value) {
      await loadEditDetail()
      return
    }
    await loadProducts()
    await loadImages(selectedProduct.value?.clusterId)
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
    () => form.productId,
    (value) => {
      if (!value) {
        if (!selectedImage.value?.clusterId) {
          void loadImages()
        }
        return
      }

      if (!selectedProduct.value?.clusterId) {
        return
      }

      void loadImages(selectedProduct.value.clusterId)
    }
  )

  return {
    activeTab,
    addEnv,
    addMount,
    areas,
    authTypes,
    canCreate,
    changeFilter,
    changeImageTab,
    deployTypes,
    editLoading,
    fieldErrors,
    filters,
    form,
    formatPrice,
    frameworkRequired,
    frameworks,
    gpuModelsList,
    handleCancel,
    handleSubmit,
    imageOptions,
    imageTabs,
    loading,
    onDeployTypeChange,
    pageTitleKey,
    payTypes,
    priceUnitText,
    products,
    pvcs,
    removeEnv,
    removeMount,
    submitLabelKey,
    totalPrice,
    updateField,
    validateDisplayNameField
  }
}

import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createInferenceService } from '@/api/inference'
import { getImageList } from '@/api/image'
import { getProductFilters, getProductList } from '@/api/product'
import { getVolumeList } from '@/api/volume'
import {
  INFERENCE_AUTH_TYPES,
  INFERENCE_IMAGE_TABS,
  INFERENCE_PAY_TYPES
} from '../constants'

export const useInferenceCreate = ({ t, router }) => {
  const translate = (key) => (typeof t === 'function' ? t(key) : key)

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

  const form = reactive({
    displayName: '',
    framework: '',
    deployType: 'STANDALONE',
    modelPvcId: '',
    modelMountPath: '/model',
    imageId: '',
    productId: '',
    payType: 1,
    workerCount: 2,
    autoRestart: false,
    maxRestarts: 3,
    servicePort: 30000,
    authType: 1,
    command: '',
    args: '',
    mounts: [],
    envs: []
  })

  const filters = reactive({
    area: '',
    gpuModel: ''
  })

  const images = ref([])
  const pvcs = ref([])
  const products = ref([])
  const areas = ref([])
  const gpuModelsList = ref([])
  const activeTab = ref('base')
  const loading = ref(false)

  const imageTabs = computed(() =>
    INFERENCE_IMAGE_TABS.map((tab) => ({
      ...tab,
      label: translate(tab.labelKey)
    }))
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

  const canCreate = computed(() => {
    if (!form.displayName || !form.modelPvcId || !form.imageId || !form.productId) return false
    if (!hasCommand.value) return false
    if (frameworkRequired.value && !form.framework) return false
    if (form.deployType === 'DISTRIBUTED' && form.workerCount < 2) return false
    return true
  })

  const priceUnitText = computed(() => {
    const units = {
      1: translate('unitHour'),
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }
    return units[form.payType] || translate('unitHour')
  })

  const getUnitPrice = (product) => {
    if (!product) return 0

    switch (form.payType) {
      case 1:
        return product.priceHourly || 0
      case 2:
        return product.priceDaily || 0
      case 3:
        return product.priceWeekly || 0
      case 4:
        return product.priceMonthly || 0
      default:
        return 0
    }
  }

  const totalPrice = computed(() => {
    if (!form.productId) return '0.00'
    const product = products.value.find((item) => item.id === form.productId)
    if (!product) return '0.00'

    let total = getUnitPrice(product)
    if (form.deployType === 'DISTRIBUTED') {
      total *= form.workerCount
    }

    return form.payType === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const formatPrice = (product) =>
    form.payType === 1 ? getUnitPrice(product).toFixed(4) : getUnitPrice(product).toFixed(2)

  const onDeployTypeChange = (value) => {
    form.deployType = value
    if (value === 'STANDALONE') {
      form.framework = ''
      return
    }

    if (!form.framework) {
      form.framework = 'VLLM'
    }
  }

  const loadImages = async () => {
    try {
      const params = { page: 1, pageSize: 100, usageType: 3 }
      if (activeTab.value === 'base') params.type = 1
      if (activeTab.value === 'my') params.type = 2

      const res = await getImageList(params)
      if (res.code === 0) {
        images.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载镜像失败', error)
    }
  }

  const changeImageTab = (tab) => {
    activeTab.value = tab
    form.imageId = ''
    loadImages()
  }

  const loadProducts = async (silent = false) => {
    if (!silent) loading.value = true

    try {
      const params = { page: 1, pageSize: 100, productType: 1 }
      if (filters.area) params.area = filters.area
      if (filters.gpuModel) params.gpuModel = filters.gpuModel

      const res = await getProductList(params)
      if (res.code === 0) {
        products.value = res.data?.list || []
        if (products.value.length > 0 && !form.productId) {
          form.productId = products.value[0].id
        }
      }
    } catch (error) {
      console.error(error)
    } finally {
      if (!silent) loading.value = false
    }
  }

  const loadFilters = async () => {
    try {
      const res = await getProductFilters({ productType: 1 })
      if (res.code === 0) {
        areas.value = res.data?.areas || []
        gpuModelsList.value = res.data?.gpuModels || []
      }
    } catch (error) {
      console.error('加载筛选条件失败', error)
    }
  }

  const loadPvcs = async () => {
    try {
      const res = await getVolumeList({ page: 1, pageSize: 100 })
      if (res.code === 0) {
        pvcs.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载 PVC 失败', error)
    }
  }

  const changeFilter = (key, value) => {
    filters[key] = value
    loadProducts()
  }

  const addMount = () => {
    form.mounts.push({ pvcId: null, mountPath: '', subPath: '', readOnly: false })
  }

  const removeMount = (index) => {
    form.mounts.splice(index, 1)
  }

  const addEnv = () => {
    form.envs.push({ name: '', value: '' })
  }

  const removeEnv = (index) => {
    form.envs.splice(index, 1)
  }

  const handleCancel = () => {
    router.push({ name: 'inference' })
  }

  const handleSubmit = async () => {
    if (!canCreate.value) {
      ElMessage.warning('请填写完整信息')
      return
    }

    loading.value = true
    try {
      const params = {
        name: form.displayName,
        framework: form.framework || undefined,
        deployType: form.deployType,
        modelPvcId: form.modelPvcId,
        modelMountPath: form.modelMountPath || '/model',
        imageId: form.imageId,
        productId: form.productId,
        payType: form.payType,
        servicePort: form.servicePort,
        authType: form.authType,
        command: form.command.trim(),
        args: form.args.trim()
          ? form.args.trim().split(/\n+/).map((item) => item.trim()).filter(Boolean)
          : [],
        mounts: form.mounts
          .filter((mount) => mount.pvcId && mount.mountPath)
          .map((mount) => ({
            pvcId: mount.pvcId,
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

      const res = await createInferenceService(params)
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        router.push({ name: 'inference' })
        return
      }

      ElMessage.error(res.msg || translate('error'))
    } catch (error) {
      console.error('创建失败', error)
      ElMessage.error(translate('error'))
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    loadImages()
    loadPvcs()
    loadFilters()
    loadProducts()
  })

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
    payTypes,
    priceUnitText,
    products,
    pvcs,
    removeEnv,
    removeMount,
    totalPrice
  }
}

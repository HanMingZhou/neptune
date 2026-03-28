import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createTrainingJob } from '@/api/training'
import { getImageList } from '@/api/image'
import { getAggregateProductList, getProductFilters } from '@/api/product'
import { getVolumeList } from '@/api/volume'
import {
  TRAINING_FRAMEWORK_TYPES,
  TRAINING_IMAGE_TABS,
  TRAINING_PAY_TYPES
} from '../constants'

export const useTrainingCreate = ({ t, router }) => {
  const translate = (key, params) => (typeof t === 'function' ? t(key, params) : key)

  const form = reactive({
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

  const filters = reactive({
    area: '',
    gpuModel: '',
    cpuModel: ''
  })

  const images = ref([])
  const pvcs = ref([])
  const products = ref([])
  const selectedProduct = ref(null)
  const areas = ref([])
  const gpuModels = ref([])
  const cpuModels = ref([])
  const creating = ref(false)

  const activeTab = ref('base')
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

  const filteredImages = computed(() => images.value)
  const showWorkerCount = computed(() => ['PYTORCH_DDP', 'MPI'].includes(form.frameworkType))

  const totalResources = computed(() => {
    if (!selectedProduct.value) return null
    const nodes = showWorkerCount.value ? form.workerCount : 1
    const gpu = (selectedProduct.value.gpuCount > 0 ? form.gpuPerWorker : 0) * nodes
    return { gpu, nodes }
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

  const formatPrice = (product) =>
    form.payType === 1 ? getUnitPrice(product).toFixed(4) : getUnitPrice(product).toFixed(2)

  const totalPrice = computed(() => {
    if (!selectedProduct.value) return '0.00'

    const nodes = showWorkerCount.value ? form.workerCount : 1
    const total = getUnitPrice(selectedProduct.value) * nodes
    return form.payType === 1 ? total.toFixed(4) : total.toFixed(2)
  })

  const availableCapacity = computed(() => selectedProduct.value?.available || 0)

  const totalAllowedNodes = computed(() => {
    if (!selectedProduct.value) return 0
    if (!showWorkerCount.value) {
      return availableCapacity.value
    }

    const strictMax = selectedProduct.value.strictMax
    const balancedMax = selectedProduct.value.balancedMax
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
    if (!form.name || !form.imageId || !form.startupCommand || !selectedProduct.value) {
      return false
    }

    const requiredNodes = showWorkerCount.value ? form.workerCount : 1
    if (requiredNodes > totalAllowedNodes.value) return false
    if (showWorkerCount.value && form.workerCount < 2) return false
    return true
  })

  const loadImages = async () => {
    try {
      const params = { page: 1, pageSize: 100, usageType: 2 }
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
    form.imageId = null
    loadImages()
  }

  const loadProducts = async () => {
    try {
      const params = {
        page: 1,
        pageSize: 100,
        productType: 1
      }

      if (filters.area) params.area = filters.area
      if (filters.gpuModel) params.gpuModel = filters.gpuModel
      if (filters.cpuModel) params.cpuModel = filters.cpuModel

      const res = await getAggregateProductList(params)
      if (res.code === 0) {
        products.value = res.data?.list || []
        if (selectedProduct.value && !products.value.find((item) => item.id === selectedProduct.value.id)) {
          selectedProduct.value = null
        }
      }
    } catch (error) {
      console.error('加载产品失败', error)
    }
  }

  const loadFilters = async () => {
    try {
      const res = await getProductFilters({ productType: 1 })
      if (res.code === 0) {
        areas.value = res.data?.areas || []
        gpuModels.value = res.data?.gpuModels || []
        cpuModels.value = res.data?.cpuModels || []
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

  const selectProduct = (product) => {
    selectedProduct.value = product
    if (product.gpuCount > 0) {
      form.gpuPerWorker = 1
      form.gpuType = product.gpuModel
    } else {
      form.gpuPerWorker = 0
      form.gpuType = ''
    }
  }

  const decreaseWorker = () => {
    if (form.workerCount > 2) form.workerCount -= 1
  }

  const increaseWorker = () => {
    if (form.workerCount < maxWorkerCount.value) form.workerCount += 1
  }

  const addMount = () => {
    form.mounts.push({ mountType: 'DATASET', pvcId: null, mountPath: '', readOnly: false })
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

  const insertMpiExample = () => {
    const totalSlots = form.workerCount * form.gpuPerWorker
    const np = totalSlots > 0 ? totalSlots : form.workerCount
    form.startupCommand = `mpirun -np ${np} -H \${MPI_HOSTS} --allow-run-as-root --bind-to none -map-by slot -x NCCL_DEBUG=INFO python train.py`
  }

  const insertPytorchExample = () => {
    const nprocPerNode = form.gpuPerWorker > 0 ? form.gpuPerWorker : 1
    form.startupCommand = `torchrun --nnodes=\${WORLD_SIZE} --node_rank=\${RANK} --master_addr=\${MASTER_ADDR} --master_port=\${MASTER_PORT} --nproc_per_node=${nprocPerNode} train.py`
  }

  const goBack = () => router.go(-1)

  const handleCreate = async () => {
    if (!canCreate.value) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    creating.value = true
    try {
      const params = {
        name: form.name,
        frameworkType: form.frameworkType,
        imageId: form.imageId,
        startupCommand: form.startupCommand,
        clusterId: selectedProduct.value.clusterId,
        resourceId: selectedProduct.value.id,
        templateProductId: selectedProduct.value.templateProductId || selectedProduct.value.id,
        instanceCount: showWorkerCount.value ? form.workerCount : 1,
        workerCount: showWorkerCount.value ? form.workerCount : 0,
        scheduleStrategy: showWorkerCount.value ? form.scheduleStrategy : 'BALANCED',
        enableTensorboard: form.enableTensorboard,
        tensorboardLogPath: form.tensorboardLogPath,
        payType: form.payType,
        mounts: form.mounts
          .filter((mount) => mount.pvcId && mount.mountPath)
          .map((mount) => ({
            mountType: mount.mountType,
            pvcId: mount.pvcId,
            mountPath: mount.mountPath,
            readOnly: mount.readOnly
          })),
        envs: form.envs.filter((env) => env.name && env.value)
      }

      const res = await createTrainingJob(params)
      if (res.code === 0) {
        ElMessage.success(translate('createSuccess'))
        router.go(-1)
      }
    } catch (error) {
      console.error('创建失败', error)
    } finally {
      creating.value = false
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
    availableCapacity,
    canCreate,
    changeFilter,
    changeImageTab,
    cpuModels,
    creating,
    decreaseWorker,
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
    totalResources
  }
}

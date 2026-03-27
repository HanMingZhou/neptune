import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getProductFilters, getProductList } from '@/api/product'
import { getImageList } from '@/api/image'
import { createNotebook } from '@/api/notebook'
import { getSSHKeyList } from '@/api/sshkey'
import { getVolumeList } from '@/api/volume'
import {
  DEFAULT_VOLUME_MOUNT_PATH,
  NOTEBOOK_IMAGE_TABS,
  NOTEBOOK_PAY_TYPES
} from '../constants'
import { validateTensorBoardPath } from '../validators'

export const useNotebookCreate = ({ t, router }) => {
  const translate = (key) => (typeof t === 'function' ? t(key) : key)

  const products = ref([])
  const selectedProduct = ref(null)
  const gpuCount = ref(1)
  const payType = ref(1)
  const filters = ref({ area: '', gpuModel: '', cpuModel: '' })
  const areas = ref([])
  const gpuModels = ref([])
  const cpuModels = ref([])

  const images = ref([])
  const selectedImage = ref('')
  const activeTab = ref('base')

  const instanceName = ref('')
  const enableTensorboard = ref(false)
  const tensorboardLogPath = ref('')
  const sshKeys = ref([])
  const selectedSshKey = ref(null)
  const enableSshPassword = ref(false)

  const availableVolumes = ref([])
  const selectedVolumeId = ref(null)
  const volumeMountPath = ref(DEFAULT_VOLUME_MOUNT_PATH)

  const payTypes = NOTEBOOK_PAY_TYPES

  const imageTabs = computed(() =>
    NOTEBOOK_IMAGE_TABS.map((tab) => ({
      ...tab,
      label: translate(tab.labelKey)
    }))
  )

  const selectedVolumeName = computed(() => {
    const volume = availableVolumes.value.find((item) => item.id === selectedVolumeId.value)
    return volume ? `${volume.name} (${volume.size})` : ''
  })

  const priceUnitText = computed(() => {
    const units = {
      1: translate('unitHour'),
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }
    return units[payType.value] || translate('unitHour')
  })

  const imageDescription = computed(() => {
    const descriptions = {
      base: translate('baseImageDesc'),
      community: translate('communityImageDesc'),
      my: translate('myImageDesc')
    }
    return descriptions[activeTab.value] || ''
  })

  const filteredImages = computed(() => images.value)

  const getUnitPrice = (product) => {
    if (!product) return 0

    switch (payType.value) {
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

  const formatPrice = (product) => {
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
    Boolean(selectedProduct.value && selectedImage.value && instanceName.value.trim())
  )

  const syncSelectedProduct = (list) => {
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

    const nextSelected = list.find((item) => item.id === selectedProduct.value.id)
    if (!nextSelected) {
      selectProduct(list[0])
      return
    }

    selectedProduct.value = nextSelected
  }

  const loadProducts = async () => {
    try {
      const params = { page: 1, pageSize: 100, productType: 1 }

      if (filters.value.area) params.area = filters.value.area
      if (filters.value.gpuModel) params.gpuModel = filters.value.gpuModel
      if (filters.value.cpuModel) params.cpuModel = filters.value.cpuModel

      const res = await getProductList(params)
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

  const loadImages = async () => {
    try {
      const params = { page: 1, pageSize: 100, usageType: 1 }
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

  const loadSSHKeys = async () => {
    try {
      const res = await getSSHKeyList({ page: 1, pageSize: 100 })
      if (res.code === 0) {
        sshKeys.value = res.data?.list || []
      }
    } catch (error) {
      console.error('加载SSH密钥失败', error)
    }
  }

  const loadVolumes = async () => {
    if (!selectedProduct.value) return

    try {
      const res = await getVolumeList({
        page: 1,
        pageSize: 100,
        clusterId: selectedProduct.value.clusterId
      })

      if (res.code === 0) {
        availableVolumes.value = res.data?.list || []

        if (
          selectedVolumeId.value &&
          !availableVolumes.value.some((item) => item.id === selectedVolumeId.value)
        ) {
          selectedVolumeId.value = null
        }
      }
    } catch (error) {
      console.error('加载数据盘失败', error)
    }
  }

  const selectProduct = (product) => {
    selectedProduct.value = product
    gpuCount.value = 1
    selectedVolumeId.value = null
    volumeMountPath.value = DEFAULT_VOLUME_MOUNT_PATH
    loadVolumes()
  }

  const changeFilter = (key, value) => {
    filters.value = {
      ...filters.value,
      [key]: value
    }
    loadProducts()
  }

  const changeImageTab = (tab) => {
    activeTab.value = tab
    selectedImage.value = ''
    loadImages()
  }

  const onVolumeChange = (value) => {
    selectedVolumeId.value = value
    if (value && !volumeMountPath.value) {
      volumeMountPath.value = DEFAULT_VOLUME_MOUNT_PATH
    }
  }

  const goBack = () => router.go(-1)

  const handleCreate = async () => {
    if (!canCreate.value) {
      ElMessage.warning(translate('fillAllFields'))
      return
    }

    if (enableTensorboard.value && tensorboardLogPath.value) {
      const pathError = validateTensorBoardPath(tensorboardLogPath.value, translate)
      if (pathError) {
        ElMessage.error(pathError)
        return
      }
    }

    try {
      const volumeMounts = []
      if (selectedVolumeId.value) {
        volumeMounts.push({
          pvcId: selectedVolumeId.value,
          mountsPath: volumeMountPath.value || '/data/volume-1'
        })
      }

      const params = {
        displayName: instanceName.value,
        productId: selectedProduct.value.id,
        imageId: selectedImage.value,
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

      ElMessage.error(res.msg || translate('createFailed'))
    } catch (error) {
      console.error('创建失败', error)
      ElMessage.error(translate('createFailed'))
    }
  }

  onMounted(() => {
    loadProducts()
    loadFilters()
    loadImages()
    loadSSHKeys()
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
    volumeMountPath
  }
}

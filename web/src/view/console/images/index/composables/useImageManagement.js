import { reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createImage, deleteImage, getImageList, updateImage } from '@/api/image'

export function useImageManagement({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const submitting = ref(false)
  const images = ref([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = 10
  const filterKeyword = ref('')
  const filterType = ref('')
  const filterUsageType = ref('')
  const showDialog = ref(false)
  const isEdit = ref(false)

  const form = reactive({
    id: null,
    name: '',
    type: 1,
    usageType: 1,
    imageAddr: '',
    area: '',
    size: '',
    imagePath: ''
  })

  const formRules = {
    name: [{ required: true, message: () => translate('inputName'), trigger: 'blur' }],
    usageType: [{ required: true, message: () => translate('pleaseSelect'), trigger: 'change' }]
  }

  const resetForm = () => {
    form.id = null
    form.name = ''
    form.type = 1
    form.usageType = 1
    form.imageAddr = ''
    form.area = ''
    form.size = ''
    form.imagePath = ''
  }

  const fetchImages = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getImageList({
        page: currentPage.value,
        pageSize,
        name: filterKeyword.value || undefined,
        type: filterType.value || undefined,
        usageType: filterUsageType.value || undefined
      })

      if (res.code === 0) {
        images.value = res.data?.list || []
        total.value = res.data?.total || 0
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error) {
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const handleRefresh = (silent = false) => fetchImages(silent)

  const handleSearch = () => {
    currentPage.value = 1
    fetchImages()
  }

  const handleReset = () => {
    filterKeyword.value = ''
    filterType.value = ''
    filterUsageType.value = ''
    currentPage.value = 1
    fetchImages()
  }

  const handlePageChange = () => {
    fetchImages()
  }

  const openCreateDialog = () => {
    resetForm()
    isEdit.value = false
    showDialog.value = true
  }

  const openEditDialog = (row) => {
    isEdit.value = true
    form.id = row.id
    form.name = row.name
    form.type = row.type
    form.usageType = row.usageType
    form.imageAddr = row.image || ''
    form.area = row.area || ''
    form.size = row.size || ''
    form.imagePath = row.imagePath || ''
    showDialog.value = true
  }

  const handleDialogClosed = () => {
    if (!showDialog.value) {
      resetForm()
      isEdit.value = false
    }
  }

  const submitImage = async () => {
    submitting.value = true

    try {
      const data = { ...form }
      if (!isEdit.value) {
        delete data.id
      }

      const res = isEdit.value ? await updateImage(data) : await createImage(data)
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        showDialog.value = false
        fetchImages()
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error) {
      ElMessage.error(translate('failed'))
    } finally {
      submitting.value = false
    }
  }

  const handleDelete = (row) => {
    ElMessageBox.confirm(translate('confirmDelete', { name: row.name }), translate('tip'), {
      confirmButtonText: translate('confirm'),
      cancelButtonText: translate('cancel'),
      type: 'danger'
    })
      .then(async () => {
        try {
          const res = await deleteImage({ id: row.id })
          if (res.code === 0) {
            ElMessage.success(translate('success'))
            fetchImages()
            return
          }

          ElMessage.error(res.msg || translate('failed'))
        } catch (error) {
          ElMessage.error(translate('failed'))
        }
      })
      .catch(() => {})
  }

  return {
    currentPage,
    fetchImages,
    filterKeyword,
    filterType,
    filterUsageType,
    form,
    formRules,
    handleDelete,
    handleDialogClosed,
    handlePageChange,
    handleRefresh,
    handleReset,
    handleSearch,
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

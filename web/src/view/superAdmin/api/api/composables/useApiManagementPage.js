import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createApi,
  deleteApi,
  deleteApisByIds,
  enterSyncApi,
  freshCasbin,
  getApiById,
  getApiGroups,
  getApiList,
  ignoreApi,
  syncApi,
  updateApi
} from '@/api/api'

const createDefaultForm = () => ({
  path: '',
  apiGroup: '',
  method: '',
  description: ''
})

const createDefaultSyncData = () => ({
  newApis: [],
  deleteApis: [],
  ignoreApis: []
})

export function useApiManagementPage({ t }) {
  const translate = t || ((key) => key)

  const apiCompletionLoading = ref(false)
  const apiGroupMap = ref({})
  const apiGroupOptions = ref([])
  const dialogFormVisible = ref(false)
  const form = reactive(createDefaultForm())
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(10)
  const searchApiGroup = ref('')
  const searchDescription = ref('')
  const searchMethod = ref('')
  const searchPath = ref('')
  const selectedApis = ref([])
  const syncing = ref(false)
  const syncApiData = ref(createDefaultSyncData())
  const syncApiFlag = ref(false)
  const tableData = ref([])
  const total = ref(0)
  const type = ref('addApi')

  const methodOptions = computed(() => [
    { value: 'POST', label: translate('methodCreate'), type: 'success' },
    { value: 'GET', label: translate('methodRead'), type: '' },
    { value: 'PUT', label: translate('methodUpdate'), type: 'warning' },
    { value: 'DELETE', label: translate('methodDelete'), type: 'danger' }
  ])

  const rules = computed(() => ({
    path: [{ required: true, message: translate('inputApiPath'), trigger: 'blur' }],
    apiGroup: [{ required: true, message: translate('inputApiGroup'), trigger: 'blur' }],
    method: [{ required: true, message: translate('selectMethod'), trigger: 'blur' }],
    description: [{ required: true, message: translate('inputApiDesc'), trigger: 'blur' }]
  }))

  const dialogTitle = computed(() => type.value === 'addApi' ? translate('addApi') : translate('editApi'))
  const hasSelection = computed(() => selectedApis.value.length > 0)
  const isAllSelected = computed(() => {
    if (tableData.value.length === 0) {
      return false
    }

    return tableData.value.every((row) => selectedApis.value.some((item) => item.ID === row.ID))
  })
  const selectedApiIds = computed(() => selectedApis.value.map((item) => item.ID))

  const getMethodClass = (method) => {
    const classes = {
      POST: 'bg-emerald-500/10 text-emerald-500',
      GET: 'bg-primary/10 text-primary',
      PUT: 'bg-amber-500/10 text-amber-500',
      DELETE: 'bg-red-500/10 text-red-500'
    }

    return classes[method] || 'bg-slate-500/10 text-slate-500'
  }

  const resetForm = () => {
    Object.assign(form, createDefaultForm())
  }

  const getGroup = async () => {
    const res = await getApiGroups()
    if (res.code === 0) {
      apiGroupOptions.value = (res.data?.groups || []).map((item) => ({ label: item, value: item }))
      apiGroupMap.value = res.data?.apiGroupMap || {}
    }
  }

  const getTableData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getApiList({
        page: page.value,
        pageSize: pageSize.value,
        path: searchPath.value || undefined,
        description: searchDescription.value || undefined,
        apiGroup: searchApiGroup.value || undefined,
        method: searchMethod.value || undefined
      })

      if (res.code === 0) {
        tableData.value = res.data?.list || []
        total.value = res.data?.total || 0
        page.value = res.data?.page || 1
        pageSize.value = res.data?.pageSize || 10
      }
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await Promise.all([
      getTableData(),
      getGroup()
    ])
  }

  const isSelected = (row) => selectedApis.value.some((item) => item.ID === row.ID)

  const toggleSelect = (row) => {
    const index = selectedApis.value.findIndex((item) => item.ID === row.ID)
    if (index > -1) {
      selectedApis.value.splice(index, 1)
    } else {
      selectedApis.value.push(row)
    }
  }

  const toggleSelectAll = (checked) => {
    if (checked) {
      selectedApis.value = [...tableData.value]
      return
    }

    selectedApis.value = []
  }

  const onReset = () => {
    searchPath.value = ''
    searchDescription.value = ''
    searchApiGroup.value = ''
    searchMethod.value = ''
    getTableData()
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (value) => {
    pageSize.value = value
    getTableData()
  }

  const handleCurrentChange = (value) => {
    page.value = value
    getTableData()
  }

  const onDeleteSelected = async () => {
    try {
      await ElMessageBox.confirm(
        translate('deleteConfirm'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const ids = selectedApis.value.map((item) => item.ID)
      const res = await deleteApisByIds({ ids })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: translate('success') })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value -= 1
        }
        selectedApis.value = []
        await getTableData()
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const onFresh = async () => {
    try {
      await ElMessageBox.confirm(
        translate('confirmFreshCache'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await freshCasbin()
      if (res.code === 0) {
        ElMessage({ type: 'success', message: translate('success') })
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const closeSyncDialog = () => {
    syncApiFlag.value = false
  }

  const onSync = async () => {
    const res = await syncApi()
    if (res.code === 0) {
      res.data?.newApis?.forEach((item) => {
        item.apiGroup = apiGroupMap.value[item.path.split('/')[1]]
      })
      syncApiData.value = {
        newApis: res.data?.newApis || [],
        deleteApis: res.data?.deleteApis || [],
        ignoreApis: res.data?.ignoreApis || []
      }
      syncApiFlag.value = true
    }
  }

  const ignoreApiEntry = async (row, flag) => {
    const res = await ignoreApi({ path: row.path, method: row.method, flag })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: res.msg })
      if (flag) {
        syncApiData.value.newApis = syncApiData.value.newApis.filter(
          (item) => !(item.path === row.path && item.method === row.method)
        )
        syncApiData.value.ignoreApis.push(row)
      } else {
        syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(
          (item) => !(item.path === row.path && item.method === row.method)
        )
        syncApiData.value.newApis.push(row)
      }
    }
  }

  const addApiFromSync = async (row) => {
    if (!row.apiGroup) {
      ElMessage({ type: 'error', message: translate('selectApiGroupFirst') })
      return
    }

    if (!row.description) {
      ElMessage({ type: 'error', message: translate('inputApiDescFirst') })
      return
    }

    const res = await createApi(row)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: translate('apiAddSuccess'), showClose: true })
      syncApiData.value.newApis = syncApiData.value.newApis.filter(
        (item) => !(item.path === row.path && item.method === row.method)
      )
    }

    await Promise.all([
      getTableData(),
      getGroup()
    ])
  }

  const submitSync = async () => {
    if (syncApiData.value.newApis.some((item) => !item.apiGroup || !item.description)) {
      ElMessage({ type: 'error', message: translate('apiMissingInfo') })
      return
    }

    syncing.value = true
    const res = await enterSyncApi(syncApiData.value)
    syncing.value = false

    if (res.code === 0) {
      ElMessage({ type: 'success', message: res.msg })
      syncApiFlag.value = false
      await getTableData()
    }
  }

  const apiCompletion = async () => {
    ElMessage({ type: 'warning', message: translate('aiCompletionFailed') })
  }

  const openCreateDialog = () => {
    resetForm()
    type.value = 'addApi'
    dialogFormVisible.value = true
  }

  const closeDialog = () => {
    resetForm()
    dialogFormVisible.value = false
  }

  const openEditDialog = async (row) => {
    const res = await getApiById({ id: row.ID })
    if (res.code === 0 && res.data?.api) {
      resetForm()
      Object.assign(form, res.data.api)
      type.value = 'edit'
      dialogFormVisible.value = true
    }
  }

  const submitDialog = async () => {
    const api = type.value === 'addApi' ? createApi : updateApi
    const res = await api({ ...form })

    if (res.code === 0) {
      ElMessage({ type: 'success', message: translate('success'), showClose: true })
    }

    await getTableData()
    if (type.value === 'addApi') {
      await getGroup()
    }
    closeDialog()
  }

  const deleteSingleApi = async (row) => {
    try {
      await ElMessageBox.confirm(
        translate('confirmDeleteApi'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await deleteApi(row)
      if (res.code === 0) {
        ElMessage({ type: 'success', message: translate('success') })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await Promise.all([
          getTableData(),
          getGroup()
        ])
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  return {
    apiCompletion,
    apiCompletionLoading,
    apiGroupOptions,
    closeDialog,
    deleteSingleApi,
    closeSyncDialog,
    dialogFormVisible,
    dialogTitle,
    form,
    getMethodClass,
    getTableData,
    handleCurrentChange,
    handleSizeChange,
    hasSelection,
    ignoreApiEntry,
    initialize,
    isAllSelected,
    isSelected,
    loading,
    methodOptions,
    onDeleteSelected,
    onFresh,
    onReset,
    onSubmit,
    onSync,
    openCreateDialog,
    openEditDialog,
    page,
    pageSize,
    rules,
    searchApiGroup,
    searchDescription,
    searchMethod,
    searchPath,
    selectedApiIds,
    submitDialog,
    submitSync,
    syncing,
    syncApiData,
    syncApiFlag,
    tableData,
    toggleSelect,
    toggleSelectAll,
    total,
    addApiFromSync
  }
}


import { computed, reactive, ref } from 'vue'
import type { FormRules } from 'element-plus'
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
import type { Translator } from '@/types/consoleResource'
import type {
  ApiForm,
  ApiGroupListData,
  ApiListData,
  ApiListItem,
  ApiMethodOption,
  ApiSyncData,
  LabelValueOption
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'

interface UseApiManagementPageOptions {
  t?: Translator
}

type ApiDialogType = 'addApi' | 'edit'

const createDefaultForm = (): ApiForm => ({
  path: '',
  apiGroup: '',
  method: '',
  description: ''
})

const createDefaultSyncData = (): ApiSyncData => ({
  newApis: [],
  deleteApis: [],
  ignoreApis: []
})

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useApiManagementPage({ t }: UseApiManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const apiCompletionLoading = ref(false)
  const apiGroupMap = ref<Record<string, string>>({})
  const apiGroupOptions = ref<LabelValueOption[]>([])
  const dialogFormVisible = ref(false)
  const form = reactive<ApiForm>(createDefaultForm())
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(15)
  const searchApiGroup = ref('')
  const searchDescription = ref('')
  const searchMethod = ref('')
  const searchPath = ref('')
  const selectedApis = ref<ApiListItem[]>([])
  const syncing = ref(false)
  const syncApiData = ref<ApiSyncData>(createDefaultSyncData())
  const syncApiFlag = ref(false)
  const tableData = ref<ApiListItem[]>([])
  const total = ref(0)
  const type = ref<ApiDialogType>('addApi')

  const methodOptions = computed<ApiMethodOption[]>(() => [
    { value: 'POST', label: translate('methodCreate'), type: 'success' },
    { value: 'GET', label: translate('methodRead'), type: '' },
    { value: 'PUT', label: translate('methodUpdate'), type: 'warning' },
    { value: 'DELETE', label: translate('methodDelete'), type: 'danger' }
  ])

  const rules = computed<FormRules<ApiForm>>(() => ({
    path: [
      { required: true, message: translate('inputApiPath'), trigger: 'blur' }
    ],
    apiGroup: [
      { required: true, message: translate('inputApiGroup'), trigger: 'blur' }
    ],
    method: [
      { required: true, message: translate('selectMethod'), trigger: 'blur' }
    ],
    description: [
      { required: true, message: translate('inputApiDesc'), trigger: 'blur' }
    ]
  }))

  const dialogTitle = computed(() =>
    type.value === 'addApi' ? translate('addApi') : translate('editApi')
  )
  const hasSelection = computed(() => selectedApis.value.length > 0)
  const isAllSelected = computed(() => {
    if (tableData.value.length === 0) {
      return false
    }

    return tableData.value.every((row) =>
      selectedApis.value.some((item) => item.ID === row.ID)
    )
  })
  const selectedApiIds = computed<Array<ApiListItem['ID']>>(() =>
    selectedApis.value.map((item) => item.ID)
  )

  const getMethodClass = (method: string): string => {
    const classes: Record<string, string> = {
      POST: 'bg-emerald-500/10 text-emerald-500',
      GET: 'bg-primary/10 text-primary',
      PUT: 'bg-amber-500/10 text-amber-500',
      DELETE: 'bg-red-500/10 text-red-500'
    }

    return classes[method] || 'bg-slate-500/10 text-slate-500'
  }

  const resetForm = (): void => {
    Object.assign(form, createDefaultForm())
  }

  const getGroup = async (): Promise<void> => {
    const res = (await getApiGroups()) as ApiResponse<ApiGroupListData>
    if (res.code === 0) {
      apiGroupOptions.value = (res.data?.groups ?? []).map((item) => ({
        label: item,
        value: item
      }))
      apiGroupMap.value = res.data?.apiGroupMap ?? {}
    }
  }

  const getTableData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getApiList({
        page: page.value,
        pageSize: pageSize.value,
        path: searchPath.value || undefined,
        description: searchDescription.value || undefined,
        apiGroup: searchApiGroup.value || undefined,
        method: searchMethod.value || undefined
      })) as ApiResponse<ApiListData>

      if (res.code === 0) {
        tableData.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? 1
        pageSize.value = res.data?.pageSize ?? 10
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!silent) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await Promise.all([getTableData(), getGroup()])
  }

  const isSelected = (row: ApiListItem): boolean =>
    selectedApis.value.some((item) => item.ID === row.ID)

  const toggleSelect = (row: ApiListItem): void => {
    const index = selectedApis.value.findIndex((item) => item.ID === row.ID)
    if (index > -1) {
      selectedApis.value.splice(index, 1)
    } else {
      selectedApis.value.push(row)
    }
  }

  const toggleSelectAll = (checked: boolean): void => {
    if (checked) {
      selectedApis.value = [...tableData.value]
      return
    }

    selectedApis.value = []
  }

  const onReset = (): void => {
    searchPath.value = ''
    searchDescription.value = ''
    searchApiGroup.value = ''
    searchMethod.value = ''
    void getTableData()
  }

  const onSubmit = (): void => {
    page.value = 1
    void getTableData()
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    void getTableData()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void getTableData()
  }

  const onDeleteSelected = async (): Promise<void> => {
    try {
      await ElMessageBox.confirm(translate('deleteConfirm'), translate('tip'), {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      })

      const ids = selectedApis.value.map((item) => item.ID)
      const res = await deleteApisByIds({ ids })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: translate('success') })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value -= 1
        }
        selectedApis.value = []
        await getTableData()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    }
  }

  const onFresh = async (): Promise<void> => {
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
        page.value = 1
        selectedApis.value = []
        await getTableData()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    }
  }

  const closeSyncDialog = (): void => {
    syncApiFlag.value = false
  }

  const onSync = async (): Promise<void> => {
    const res = (await syncApi()) as ApiResponse<Partial<ApiSyncData>>
    if (res.code === 0) {
      const nextNewApis = (res.data?.newApis ?? []).map((item) => ({
        ...item,
        apiGroup:
          item.apiGroup || apiGroupMap.value[item.path.split('/')[1]] || ''
      }))

      syncApiData.value = {
        newApis: nextNewApis,
        deleteApis: res.data?.deleteApis ?? [],
        ignoreApis: res.data?.ignoreApis ?? []
      }
      syncApiFlag.value = true
    }
  }

  const ignoreApiEntry = async (
    row: ApiListItem,
    flag: boolean
  ): Promise<void> => {
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

  const addApiFromSync = async (row: ApiListItem): Promise<void> => {
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
      ElMessage({
        type: 'success',
        message: translate('apiAddSuccess'),
        showClose: true
      })
      syncApiData.value.newApis = syncApiData.value.newApis.filter(
        (item) => !(item.path === row.path && item.method === row.method)
      )
    }

    await Promise.all([getTableData(), getGroup()])
  }

  const submitSync = async (): Promise<void> => {
    if (
      syncApiData.value.newApis.some(
        (item) => !item.apiGroup || !item.description
      )
    ) {
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

  const apiCompletion = async (): Promise<void> => {
    ElMessage({ type: 'warning', message: translate('aiCompletionFailed') })
  }

  const openCreateDialog = (): void => {
    resetForm()
    type.value = 'addApi'
    dialogFormVisible.value = true
  }

  const closeDialog = (): void => {
    resetForm()
    dialogFormVisible.value = false
  }

  const openEditDialog = async (row: ApiListItem): Promise<void> => {
    const res = (await getApiById({ id: row.ID })) as ApiResponse<{
      api?: ApiListItem
    }>
    if (res.code === 0 && res.data?.api) {
      resetForm()
      Object.assign(form, res.data.api)
      type.value = 'edit'
      dialogFormVisible.value = true
    }
  }

  const submitDialog = async (): Promise<void> => {
    const api = type.value === 'addApi' ? createApi : updateApi
    const res = await api({ ...form })

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate('success'),
        showClose: true
      })
    }

    await getTableData()
    if (type.value === 'addApi') {
      await getGroup()
    }
    closeDialog()
  }

  const deleteSingleApi = async (row: ApiListItem): Promise<void> => {
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
        await Promise.all([getTableData(), getGroup()])
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    }
  }

  return {
    addApiFromSync,
    apiCompletion,
    apiCompletionLoading,
    apiGroupOptions,
    closeDialog,
    closeSyncDialog,
    deleteSingleApi,
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
    total
  }
}


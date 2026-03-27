import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteSysOperationRecord,
  deleteSysOperationRecordByIds,
  getSysOperationRecordList
} from '@/api/sysOperationRecord'

const createDefaultSearchInfo = () => ({
  method: '',
  path: '',
  status: ''
})

const normalizeStatus = (status) => {
  if (status === '' || status === null || status === undefined) {
    return undefined
  }

  const numericStatus = Number(status)
  return Number.isFinite(numericStatus) ? numericStatus : status
}

export function useOperationRecordPage({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(10)
  const searchInfo = reactive(createDefaultSearchInfo())
  const selectedRecords = ref([])
  const tableData = ref([])
  const total = ref(0)

  const hasSelection = computed(() => selectedRecords.value.length > 0)
  const isAllSelected = computed(() => {
    if (tableData.value.length === 0) {
      return false
    }

    return tableData.value.every((row) => selectedRecords.value.some((item) => item.ID === row.ID))
  })
  const selectedRecordIds = computed(() => selectedRecords.value.map((item) => item.ID))

  const buildSearchParams = () => ({
    method: searchInfo.method || undefined,
    path: searchInfo.path || undefined,
    status: normalizeStatus(searchInfo.status)
  })

  const clearSelection = () => {
    selectedRecords.value = []
  }

  const getTableData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getSysOperationRecordList({
        page: page.value,
        pageSize: pageSize.value,
        ...buildSearchParams()
      })

      if (res.code === 0) {
        tableData.value = res.data?.list || []
        total.value = res.data?.total || 0
        page.value = res.data?.page || 1
        pageSize.value = res.data?.pageSize || 10
        clearSelection()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      console.error('Failed to fetch operation records:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await getTableData()
  }

  const onReset = () => {
    page.value = 1
    Object.assign(searchInfo, createDefaultSearchInfo())
    getTableData()
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (value) => {
    page.value = 1
    pageSize.value = value
    getTableData()
  }

  const handleCurrentChange = (value) => {
    page.value = value
    getTableData()
  }

  const toggleSelect = (row) => {
    const index = selectedRecords.value.findIndex((item) => item.ID === row.ID)
    if (index > -1) {
      selectedRecords.value.splice(index, 1)
      return
    }

    selectedRecords.value.push(row)
  }

  const toggleSelectAll = (checked) => {
    selectedRecords.value = checked ? [...tableData.value] : []
  }

  const deleteSelectedRecords = async () => {
    try {
      await ElMessageBox.confirm(
        translate('confirmDeleteSelected'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const ids = selectedRecords.value.map((item) => item.ID)
      const res = await deleteSysOperationRecordByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: translate('success')
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value -= 1
        }
        await getTableData()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const deleteRecord = async (row) => {
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

      const res = await deleteSysOperationRecord({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: translate('success')
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await getTableData()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(translate('failed'))
      }
    }
  }

  return {
    deleteRecord,
    deleteSelectedRecords,
    getTableData,
    handleCurrentChange,
    handleSizeChange,
    hasSelection,
    initialize,
    isAllSelected,
    loading,
    onReset,
    onSubmit,
    page,
    pageSize,
    searchInfo,
    selectedRecordIds,
    tableData,
    toggleSelect,
    toggleSelectAll,
    total
  }
}

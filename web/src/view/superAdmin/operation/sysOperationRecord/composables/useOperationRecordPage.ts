import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  deleteSysOperationRecord,
  deleteSysOperationRecordByIds,
  getSysOperationRecordList
} from '@/api/sysOperationRecord'
import type { Translator } from '@/types/consoleResource'
import type {
  OperationRecordItem,
  OperationRecordListData,
  OperationRecordSearchInfo
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'

interface UseOperationRecordPageOptions {
  t?: Translator
}

const createDefaultSearchInfo = (): OperationRecordSearchInfo => ({
  method: '',
  path: '',
  status: ''
})

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

const normalizeStatus = (
  status: OperationRecordSearchInfo['status']
): number | undefined => {
  if (status === '') {
    return undefined
  }

  const numericStatus = Number(status)
  return Number.isFinite(numericStatus) ? numericStatus : undefined
}

export function useOperationRecordPage({
  t
}: UseOperationRecordPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(15)
  const searchInfo = reactive<OperationRecordSearchInfo>(
    createDefaultSearchInfo()
  )
  const selectedRecords = ref<OperationRecordItem[]>([])
  const tableData = ref<OperationRecordItem[]>([])
  const total = ref(0)

  const hasSelection = computed(() => selectedRecords.value.length > 0)
  const isAllSelected = computed(() => {
    if (tableData.value.length === 0) {
      return false
    }

    return tableData.value.every((row) =>
      selectedRecords.value.some((item) => item.ID === row.ID)
    )
  })
  const selectedRecordIds = computed<Array<number>>(() =>
    selectedRecords.value.map((item) => item.ID)
  )

  const buildSearchParams = () => ({
    method: searchInfo.method || undefined,
    path: searchInfo.path || undefined,
    status: normalizeStatus(searchInfo.status)
  })

  const clearSelection = (): void => {
    selectedRecords.value = []
  }

  const getTableData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getSysOperationRecordList({
        page: page.value,
        pageSize: pageSize.value,
        ...buildSearchParams()
      })) as ApiResponse<OperationRecordListData>

      if (res.code === 0) {
        tableData.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? 1
        pageSize.value = res.data?.pageSize ?? 10
        clearSelection()
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch operation records:', error)
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await getTableData()
  }

  const onReset = (): void => {
    page.value = 1
    Object.assign(searchInfo, createDefaultSearchInfo())
    void getTableData()
  }

  const onSubmit = (): void => {
    page.value = 1
    void getTableData()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void getTableData()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void getTableData()
  }

  const toggleSelect = (row: OperationRecordItem): void => {
    const index = selectedRecords.value.findIndex((item) => item.ID === row.ID)
    if (index > -1) {
      selectedRecords.value.splice(index, 1)
      return
    }

    selectedRecords.value.push(row)
  }

  const toggleSelectAll = (checked: boolean): void => {
    selectedRecords.value = checked ? [...tableData.value] : []
  }

  const deleteSelectedRecords = async (): Promise<void> => {
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
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
      }
    }
  }

  const deleteRecord = async (row: OperationRecordItem): Promise<void> => {
    try {
      await ElMessageBox.confirm(translate('deleteConfirm'), translate('tip'), {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      })

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
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(getErrorMessage(error, translate('failed')))
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


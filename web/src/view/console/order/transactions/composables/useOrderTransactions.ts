import { computed, reactive, ref, watch } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, type FormRules } from 'element-plus'
import {
  getOrderOverview,
  getTransactionList,
  rechargeBalance
} from '@/api/order'
import { useBtnAuth } from '@/utils/btnAuth'
import { useUserStore } from '@/pinia/modules/user'
import type { ApiResponse } from '@/utils/request'
import type { PageListData, Translator } from '@/types/consoleResource'
import type {
  OrderOverviewData,
  RechargeForm,
  TransactionListItem,
  TransactionRecord
} from '@/types/order'

const SUPER_ADMIN_AUTHORITY_ID = 888

type TransactionType = 'Recharge' | 'Refund' | 'Consumption' | 'Adjustment'
type ValidateCallback = (error?: Error) => void

interface UseOrderTransactionsOptions {
  t?: Translator
}

export function useOrderTransactions({ t }: UseOrderTransactionsOptions = {}) {
  const translate: Translator = t || ((key: string) => key)
  const btnAuth = useBtnAuth()
  const userStore = useUserStore()

  const loading = ref(false)
  const balance = ref(0)
  const searchQuery = ref('')
  const filterType = ref('all')
  const dateRange = ref<string[]>([])
  const transactions = ref<TransactionRecord[]>([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const rechargeDialogVisible = ref(false)
  const rechargeSubmitting = ref(false)
  const rechargeForm = reactive<RechargeForm>({
    amount: 100,
    method: 4,
    remark: ''
  })

  const rechargeRules = reactive<FormRules<RechargeForm>>({
    amount: [
      {
        validator: (
          _rule: unknown,
          value: number,
          callback: ValidateCallback
        ) => {
          if (value === undefined || value === null || Number(value) <= 0) {
            callback(new Error(translate('order.rechargeAmountRequired')))
            return
          }
          callback()
        },
        trigger: 'blur'
      }
    ]
  })

  const canSystemRecharge = computed(() => {
    const authorityId = Number(
      userStore.userInfo?.authorityId ||
        userStore.userInfo?.authority?.authorityId ||
        0
    )
    return (
      Boolean((btnAuth as Record<string, unknown>).recharge_system) ||
      authorityId === SUPER_ADMIN_AUTHORITY_ID
    )
  })

  const rechargeDisabledReason = computed(() =>
    canSystemRecharge.value ? '' : translate('order.personalRechargeComingSoon')
  )

  const getMethodText = (method?: number): string => {
    const methodMap: Record<number, string> = {
      1: translate('order.alipayRecharge'),
      2: translate('order.wechatRecharge'),
      3: translate('order.balancePay'),
      4: translate('order.platformRecharge')
    }
    return methodMap[method || 0] || '-'
  }

  const getResourceTypeText = (type?: number): string => {
    const map: Record<number, string> = {
      1: translate('notebook'),
      2: translate('training'),
      3: translate('inference'),
      4: translate('storage')
    }
    return map[type || 0] || ''
  }

  const getTypeStyle = (type: TransactionType): string => {
    const map: Record<TransactionType, string> = {
      Recharge: 'order-tone-chip--recharge',
      Refund: 'order-tone-chip--warning',
      Consumption: 'order-tone-chip--consume',
      Adjustment: 'order-tone-chip--consume'
    }
    return map[type]
  }

  const getStatusStyle = (status?: number): string => {
    const map: Record<number, string> = {
      1: 'order-tone-chip--recharge',
      2: 'order-tone-chip--consume',
      3: 'order-tone-chip--warning',
      4: 'order-tone-chip--consume'
    }
    return map[status || 0] || 'order-tone-chip--consume'
  }

  const getStatusText = (status?: number): string => {
    const map: Record<number, string> = {
      1: 'inUse',
      2: 'stopped',
      3: 'settled',
      4: 'expired'
    }
    return translate(map[status || 0]) || translate('unknown')
  }

  const formatDate = (date?: string | number): string =>
    dayjs(date).format('YYYY-MM-DD HH:mm')

  const resolveTransactionStatus = (
    item: TransactionListItem
  ): Pick<TransactionRecord, 'statusStyle' | 'statusText'> => {
    if (
      item.type === 1 ||
      item.type === 3 ||
      item.type === 4 ||
      !item.orderStatus ||
      Number(item.orderStatus) <= 0
    ) {
      return {
        statusStyle: 'order-tone-chip--recharge',
        statusText: translate('success')
      }
    }

    return {
      statusStyle: getStatusStyle(item.orderStatus),
      statusText: getStatusText(item.orderStatus)
    }
  }

  const getTransactionType = (type?: number): TransactionType =>
    type === 1
      ? 'Recharge'
      : type === 2
        ? 'Consumption'
        : type === 3
          ? 'Refund'
          : 'Adjustment'

  const fetchData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const [overviewRes, txRes] = await Promise.all([
        getOrderOverview() as Promise<ApiResponse<OrderOverviewData>>,
        getTransactionList({
          page: page.value,
          pageSize: pageSize.value,
          keyword: searchQuery.value,
          type:
            filterType.value === 'all'
              ? 0
              : filterType.value === 'Recharge'
                ? 1
                : 2,
          ...(dateRange.value.length === 2
            ? {
                startTime: `${dateRange.value[0]} 00:00:00`,
                endTime: `${dateRange.value[1]} 23:59:59`
              }
            : {})
        }) as Promise<ApiResponse<PageListData<TransactionListItem>>>
      ])

      if (overviewRes.code === 0) {
        balance.value = overviewRes.data?.balance ?? 0
      }

      if (txRes.code === 0) {
        transactions.value = (txRes.data?.list || []).map((item) => {
          const type = getTransactionType(item.type)
          const rawAmount = item.amount ?? 0

          return {
            ...resolveTransactionStatus(item),
            id: item.transactionNo || `TXN-${item.id}`,
            time: formatDate(item.createdAt),
            type,
            typeLabel: translate(
              type === 'Recharge'
                ? 'recharge'
                : type === 'Consumption'
                  ? 'consumption'
                  : type === 'Refund'
                    ? 'refund'
                    : 'adjustment'
            ),
            typeStyle: getTypeStyle(type),
            amount: `${rawAmount > 0 ? '+' : ''}¥${Math.abs(rawAmount).toFixed(6)}`,
            method: getMethodText(item.method),
            rawAmount,
            orderStatus: item.orderStatus,
            remark: typeof item.remark === 'string' ? item.remark : '',
            resourceName:
              typeof item.resourceName === 'string' ? item.resourceName : '',
            resourceTypeText: getResourceTypeText(item.resourceType),
            orderNo:
              typeof item.orderNo === 'string'
                ? item.orderNo
                : item.type === 1
                  ? item.transactionNo
                  : '',
            balanceAfter: item.balanceAfter ?? 0
          }
        })
        total.value = txRes.data?.total || 0
      }
    } catch (error) {
      console.error('获取交易列表失败:', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const handleSearch = (): void => {
    page.value = 1
    void fetchData()
  }

  const handleDateChange = (): void => {
    page.value = 1
    void fetchData()
  }

  const handlePageChange = (value: number): void => {
    page.value = value
    void fetchData()
  }

  const handleSizeChange = (value: number): void => {
    pageSize.value = value
    page.value = 1
    void fetchData()
  }

  const resetRechargeForm = (): void => {
    rechargeForm.amount = 100
    rechargeForm.method = 4
    rechargeForm.remark = ''
  }

  const openRechargeDialog = (): void => {
    if (!canSystemRecharge.value) {
      ElMessage.info(
        rechargeDisabledReason.value || translate('order.comingSoon')
      )
      return
    }
    resetRechargeForm()
    rechargeDialogVisible.value = true
  }

  const closeRechargeDialog = (): void => {
    rechargeDialogVisible.value = false
    resetRechargeForm()
  }

  const submitRecharge = async (): Promise<void> => {
    if (!canSystemRecharge.value) {
      ElMessage.error(translate('order.rechargePermissionDenied'))
      return
    }

    rechargeSubmitting.value = true
    try {
      const res = await rechargeBalance({
        amount: rechargeForm.amount,
        method: rechargeForm.method,
        remark: rechargeForm.remark
      })

      if (res.code !== 0) {
        ElMessage.error(res.msg || translate('order.rechargeFailed'))
        return
      }

      ElMessage.success(translate('order.rechargeSuccess'))
      closeRechargeDialog()
      await fetchData()
    } catch (error) {
      console.error('充值失败:', error)
    } finally {
      rechargeSubmitting.value = false
    }
  }

  watch(filterType, () => {
    page.value = 1
    void fetchData()
  })

  return {
    balance,
    canSystemRecharge,
    closeRechargeDialog,
    dateRange,
    fetchData,
    filterType,
    handleDateChange,
    handlePageChange,
    handleSearch,
    handleSizeChange,
    loading,
    openRechargeDialog,
    page,
    pageSize,
    rechargeDialogVisible,
    rechargeDisabledReason,
    rechargeForm,
    rechargeRules,
    rechargeSubmitting,
    searchQuery,
    submitRecharge,
    total,
    transactions
  }
}

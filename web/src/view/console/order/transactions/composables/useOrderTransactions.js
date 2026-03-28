import { computed, reactive, ref, watch } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { getOrderOverview, getTransactionList, rechargeBalance } from '@/api/order'
import { useBtnAuth } from '@/utils/btnAuth'
import { useUserStore } from '@/pinia/modules/user'

const SUPER_ADMIN_AUTHORITY_ID = 888

export function useOrderTransactions({ t }) {
  const translate = t || ((key) => key)
  const btnAuth = useBtnAuth()
  const userStore = useUserStore()

  const loading = ref(false)
  const balance = ref(0)
  const searchQuery = ref('')
  const filterType = ref('all')
  const dateRange = ref([])
  const transactions = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const rechargeDialogVisible = ref(false)
  const rechargeSubmitting = ref(false)
  const rechargeForm = reactive({
    amount: 100,
    method: 4,
    remark: ''
  })

  const rechargeRules = {
    amount: [
      {
        validator: (_, value, callback) => {
          if (value === undefined || value === null || Number(value) <= 0) {
            callback(new Error(translate('order.rechargeAmountRequired')))
            return
          }
          callback()
        },
        trigger: 'blur'
      }
    ]
  }

  const canSystemRecharge = computed(() => {
    const authorityId = Number(userStore.userInfo?.authorityId || userStore.userInfo?.authority?.authorityId || 0)
    return Boolean(btnAuth.recharge_system) || authorityId === SUPER_ADMIN_AUTHORITY_ID
  })

  const rechargeDisabledReason = computed(() => (
    canSystemRecharge.value ? '' : translate('order.personalRechargeComingSoon')
  ))

  const getMethodText = (method) => {
    const methodMap = {
      1: translate('order.alipayRecharge'),
      2: translate('order.wechatRecharge'),
      3: translate('order.balancePay'),
      4: translate('order.platformRecharge')
    }
    return methodMap[method] || '-'
  }

  const getResourceTypeText = (type) => {
    const map = { 1: translate('notebook'), 2: translate('training'), 3: translate('inference'), 4: translate('storage') }
    return map[type] || ''
  }

  const getTypeStyle = (type) => {
    const map = {
      Recharge: 'bg-emerald-500/10 text-emerald-500',
      Refund: 'bg-amber-500/10 text-amber-500',
      Consumption: 'bg-slate-500/10 text-slate-500',
      Adjustment: 'bg-blue-500/10 text-blue-500'
    }
    return map[type] || 'bg-slate-500/10 text-slate-500'
  }

  const getStatusStyle = (status) => {
    const map = {
      1: 'bg-emerald-500/10 text-emerald-500',
      2: 'bg-slate-500/10 text-slate-500',
      3: 'bg-blue-500/10 text-blue-500',
      4: 'bg-slate-400/10 text-slate-400'
    }
    return map[status] || 'bg-slate-500/10 text-slate-500'
  }

  const getStatusText = (status) => {
    const map = {
      1: 'inUse',
      2: 'stopped',
      3: 'settled',
      4: 'expired'
    }
    return translate(map[status]) || translate('unknown')
  }

  const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

  const resolveTransactionStatus = (item) => {
    if (item.type === 1 || item.type === 3 || item.type === 4 || !item.orderStatus || Number(item.orderStatus) <= 0) {
      return {
        statusStyle: 'bg-emerald-500/10 text-emerald-500',
        statusText: translate('success')
      }
    }

    return {
      statusStyle: getStatusStyle(item.orderStatus),
      statusText: getStatusText(item.orderStatus)
    }
  }

  const fetchData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const [overviewRes, txRes] = await Promise.all([
        getOrderOverview(),
        getTransactionList({
          page: page.value,
          pageSize: pageSize.value,
          keyword: searchQuery.value,
          type: filterType.value === 'all' ? 0 : (filterType.value === 'Recharge' ? 1 : 2),
          ...(dateRange.value?.length === 2
            ? {
                startTime: `${dateRange.value[0]} 00:00:00`,
                endTime: `${dateRange.value[1]} 23:59:59`
              }
            : {})
        })
      ])

      if (overviewRes.code === 0) {
        balance.value = overviewRes.data?.balance ?? 0
      }

      if (txRes.code === 0) {
        transactions.value = (txRes.data.list || []).map((item) => ({
          ...resolveTransactionStatus(item),
          id: item.transactionNo || `TXN-${item.id}`,
          time: formatDate(item.createdAt),
          type: item.type === 1 ? 'Recharge' : (item.type === 2 ? 'Consumption' : (item.type === 3 ? 'Refund' : 'Adjustment')),
          typeLabel: translate((item.type === 1 ? 'recharge' : (item.type === 2 ? 'consumption' : item.type === 3 ? 'refund' : 'adjustment'))),
          typeStyle: getTypeStyle(item.type === 1 ? 'Recharge' : (item.type === 2 ? 'Consumption' : (item.type === 3 ? 'Refund' : 'Adjustment'))),
          amount: `${item.amount > 0 ? '+' : ''}¥${Math.abs(item.amount).toFixed(6)}`,
          method: getMethodText(item.method),
          rawAmount: item.amount,
          orderStatus: item.orderStatus,
          remark: item.remark,
          resourceName: item.resourceName || '',
          resourceTypeText: getResourceTypeText(item.resourceType),
          orderNo: item.orderNo || (item.type === 1 ? item.transactionNo : ''),
          balanceAfter: item.balanceAfter ?? 0
        }))
        total.value = txRes.data.total
      }
    } catch (error) {
      console.error('获取交易列表失败:', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const handleSearch = () => {
    page.value = 1
    fetchData()
  }

  const handleDateChange = () => {
    page.value = 1
    fetchData()
  }

  const handlePageChange = (value) => {
    page.value = value
    fetchData()
  }

  const handleSizeChange = (value) => {
    pageSize.value = value
    page.value = 1
    fetchData()
  }

  const resetRechargeForm = () => {
    rechargeForm.amount = 100
    rechargeForm.method = 4
    rechargeForm.remark = ''
  }

  const openRechargeDialog = () => {
    if (!canSystemRecharge.value) {
      ElMessage.info(rechargeDisabledReason.value || translate('order.comingSoon'))
      return
    }
    resetRechargeForm()
    rechargeDialogVisible.value = true
  }

  const closeRechargeDialog = () => {
    rechargeDialogVisible.value = false
    resetRechargeForm()
  }

  const submitRecharge = async () => {
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
    fetchData()
  })

  return {
    balance,
    dateRange,
    fetchData,
    filterType,
    handleDateChange,
    handlePageChange,
    handleSearch,
    handleSizeChange,
    loading,
    page,
    pageSize,
    searchQuery,
    total,
    transactions,
    canSystemRecharge,
    rechargeDialogVisible,
    rechargeDisabledReason,
    rechargeForm,
    rechargeRules,
    rechargeSubmitting,
    openRechargeDialog,
    closeRechargeDialog,
    submitRecharge
  }
}

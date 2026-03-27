import { ref, watch } from 'vue'
import dayjs from 'dayjs'
import { getOrderOverview, getTransactionList } from '@/api/order'

export function useOrderTransactions({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const balance = ref(0)
  const searchQuery = ref('')
  const filterType = ref('all')
  const dateRange = ref([])
  const transactions = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const getMethodText = (method) => {
    const methodMap = {
      1: '支付宝',
      2: '微信',
      3: '余额',
      4: '系统'
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
          id: item.transactionNo || `TXN-${item.id}`,
          time: formatDate(item.createdAt),
          type: item.type === 1 ? 'Recharge' : (item.type === 2 ? 'Consumption' : (item.type === 3 ? 'Refund' : 'Adjustment')),
          typeLabel: translate((item.type === 1 ? 'recharge' : (item.type === 2 ? 'consumption' : item.type === 3 ? 'refund' : 'adjustment'))),
          typeStyle: getTypeStyle(item.type === 1 ? 'Recharge' : (item.type === 2 ? 'Consumption' : (item.type === 3 ? 'Refund' : 'Adjustment'))),
          amount: `${item.amount > 0 ? '+' : ''}¥${Math.abs(item.amount).toFixed(6)}`,
          method: getMethodText(item.method),
          rawAmount: item.amount,
          orderStatus: item.orderStatus,
          statusStyle: item.orderStatus !== undefined ? getStatusStyle(item.orderStatus) : '',
          statusText: item.orderStatus !== undefined ? getStatusText(item.orderStatus) : '',
          remark: item.remark,
          resourceName: item.resourceName || '',
          resourceTypeText: getResourceTypeText(item.resourceType),
          orderNo: item.orderNo || '',
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
    transactions
  }
}

import { computed, ref } from 'vue'
import { getOrderOverview, getUsageList } from '@/api/order'

export function useOrderUsage({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const searchQuery = ref('')
  const chargeFilter = ref('all')
  const dateRange = ref([])
  const totalExpenditure = ref(0)
  const mainItem = ref({ name: translate('noData'), percent: 0 })
  const usageItems = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const getResourceTypeText = (type) => {
    const map = { 1: translate('notebook'), 2: translate('training'), 3: translate('inference'), 4: translate('storage') }
    return map[type] || '-'
  }

  const getChargeTypeStyle = (source) => (
    source === '按量计费' ? 'bg-blue-500/10 text-blue-600' : 'bg-purple-500/10 text-purple-600'
  )

  const formatUsage = (item) => {
    if (item.source === '按量计费') {
      if (item.duration > 0) {
        if (item.duration < 3600) {
          const minutes = Math.floor(item.duration / 60)
          const seconds = item.duration % 60
          let text = ''
          if (minutes > 0) text += `${minutes}${translate('unitMinute')}`
          if (seconds > 0) text += `${seconds}${translate('unitSecond')}`
          return text || `0${translate('unitSecond')}`
        }

        const hours = (item.duration / 3600).toFixed(2)
        return `${hours} ${translate('unitHour')}`
      }

      return translate('order.calculating')
    }

    const unitMap = {
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }

    return `${item.quantity} ${unitMap[item.chargeType] || ''}`
  }

  const getUnitText = (item) => {
    const unitMap = {
      1: translate('unitHour'),
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }

    return unitMap[item.chargeType] || translate('unitHour')
  }

  const getStatusStyle = (item) => {
    if (item.source === '按量计费') {
      const map = {
        1: 'bg-emerald-500/10 text-emerald-500',
        2: 'bg-slate-500/10 text-slate-500',
        3: 'bg-blue-500/10 text-blue-500'
      }
      return map[item.status] || 'bg-slate-500/10 text-slate-500'
    }

    return item.status === 1
      ? 'bg-emerald-500/10 text-emerald-500'
      : 'bg-amber-500/10 text-amber-500'
  }

  const filteredUsageItems = computed(() => {
    let items = usageItems.value

    if (chargeFilter.value !== 'all') {
      items = items.filter((item) => item.source === chargeFilter.value)
    }

    if (searchQuery.value) {
      const keyword = searchQuery.value.toLowerCase()
      items = items.filter((item) => getResourceTypeText(item.resourceType).toLowerCase().includes(keyword))
    }

    return items
  })

  const fetchData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const params = {
        page: page.value,
        pageSize: pageSize.value
      }

      if (dateRange.value && dateRange.value.length === 2) {
        params.startTime = `${dateRange.value[0]} 00:00:00`
        params.endTime = `${dateRange.value[1]} 23:59:59`
      }

      const [overviewRes, listRes] = await Promise.all([
        getOrderOverview(),
        getUsageList(params)
      ])

      if (overviewRes.code === 0) {
        const data = overviewRes.data
        totalExpenditure.value = data.mtdEstimate
        if (data.monthlyRanking && data.monthlyRanking.length > 0) {
          mainItem.value = {
            name: data.monthlyRanking[0].name,
            percent: data.monthlyRanking[0].percent.toFixed(1)
          }
        }
      }

      if (listRes.code === 0) {
        usageItems.value = listRes.data.list || []
        total.value = listRes.data.total || 0
      }
    } catch (error) {
      console.error('获取使用详情失败:', error)
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

  return {
    chargeFilter,
    dateRange,
    fetchData,
    filteredUsageItems,
    getChargeTypeStyle,
    getResourceTypeText,
    getStatusStyle,
    getUnitText,
    formatUsage,
    handleDateChange,
    handlePageChange,
    handleSearch,
    handleSizeChange,
    loading,
    mainItem,
    page,
    pageSize,
    searchQuery,
    total,
    totalExpenditure
  }
}

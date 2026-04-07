import { computed, ref } from 'vue'
import { getOrderOverview, getUsageList } from '@/api/order'
import type { ApiResponse } from '@/utils/request'
import type { PageListData, Translator } from '@/types/consoleResource'
import type {
  OrderOverviewData,
  UsageListItem,
  UsageMainItem
} from '@/types/order'

interface UseOrderUsageOptions {
  t?: Translator
}

export function useOrderUsage({ t }: UseOrderUsageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const loading = ref(false)
  const searchQuery = ref('')
  const chargeFilter = ref('all')
  const dateRange = ref<string[]>([])
  const totalExpenditure = ref(0)
  const mainItem = ref<UsageMainItem>({ name: translate('noData'), percent: 0 })
  const usageItems = ref<UsageListItem[]>([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const getResourceTypeText = (type?: number): string => {
    const map: Record<number, string> = {
      1: translate('notebook'),
      2: translate('training'),
      3: translate('inference'),
      4: translate('storage')
    }
    return map[type || 0] || '-'
  }

  const getChargeTypeStyle = (source?: string): string =>
    source === '按量计费'
      ? 'order-tone-chip--consume'
      : 'order-tone-chip--warning'

  const formatUsage = (item: UsageListItem): string => {
    if (item.source === '按量计费') {
      if ((item.duration || 0) > 0) {
        if ((item.duration || 0) < 3600) {
          const minutes = Math.floor((item.duration || 0) / 60)
          const seconds = (item.duration || 0) % 60
          let text = ''
          if (minutes > 0) text += `${minutes}${translate('unitMinute')}`
          if (seconds > 0) text += `${seconds}${translate('unitSecond')}`
          return text || `0${translate('unitSecond')}`
        }

        const hours = ((item.duration || 0) / 3600).toFixed(2)
        return `${hours} ${translate('unitHour')}`
      }

      return translate('order.calculating')
    }

    const unitMap: Record<number, string> = {
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }

    return `${item.quantity || 0} ${unitMap[item.chargeType || 0] || ''}`
  }

  const getUnitText = (item: UsageListItem): string => {
    const unitMap: Record<number, string> = {
      1: translate('unitHour'),
      2: translate('unitDay'),
      3: translate('unitWeek'),
      4: translate('unitMonth')
    }

    return unitMap[item.chargeType || 0] || translate('unitHour')
  }

  const getStatusStyle = (item: UsageListItem): string => {
    if (item.source === '按量计费') {
      const map: Record<number, string> = {
        1: 'order-tone-chip--recharge',
        2: 'order-tone-chip--consume',
        3: 'order-tone-chip--warning'
      }
      return map[item.status || 0] || 'order-tone-chip--consume'
    }

    return item.status === 1
      ? 'order-tone-chip--recharge'
      : 'order-tone-chip--warning'
  }

  const filteredUsageItems = computed<UsageListItem[]>(() => {
    let items = usageItems.value

    if (chargeFilter.value !== 'all') {
      items = items.filter((item) => item.source === chargeFilter.value)
    }

    if (searchQuery.value) {
      const keyword = searchQuery.value.toLowerCase()
      items = items.filter((item) =>
        getResourceTypeText(item.resourceType).toLowerCase().includes(keyword)
      )
    }

    return items
  })

  const fetchData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const params: {
        page: number
        pageSize: number
        startTime?: string
        endTime?: string
      } = {
        page: page.value,
        pageSize: pageSize.value
      }

      if (dateRange.value.length === 2) {
        params.startTime = `${dateRange.value[0]} 00:00:00`
        params.endTime = `${dateRange.value[1]} 23:59:59`
      }

      const [overviewRes, listRes] = await Promise.all([
        getOrderOverview() as Promise<ApiResponse<OrderOverviewData>>,
        getUsageList(params) as Promise<
          ApiResponse<PageListData<UsageListItem>>
        >
      ])

      if (overviewRes.code === 0) {
        const data = overviewRes.data || {}
        totalExpenditure.value = data.mtdEstimate ?? 0
        if (data.monthlyRanking && data.monthlyRanking.length > 0) {
          mainItem.value = {
            name: data.monthlyRanking[0].name || translate('noData'),
            percent: Number((data.monthlyRanking[0].percent || 0).toFixed(1))
          }
        }
      }

      if (listRes.code === 0) {
        usageItems.value = listRes.data?.list || []
        total.value = listRes.data?.total || 0
      }
    } catch (error) {
      console.error('获取使用详情失败:', error)
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

  return {
    chargeFilter,
    dateRange,
    fetchData,
    filteredUsageItems,
    formatUsage,
    getChargeTypeStyle,
    getResourceTypeText,
    getStatusStyle,
    getUnitText,
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

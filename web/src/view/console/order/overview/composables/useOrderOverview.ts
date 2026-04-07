import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getOrderOverview } from '@/api/order'
import type { ApiResponse } from '@/utils/request'
import type { Translator } from '@/types/consoleResource'
import type {
  ConsumptionTrendPoint,
  OrderOverviewData,
  OrderOverviewMetric,
  OrderSettlementMetric,
  OverviewRankingItem
} from '@/types/order'

interface UseOrderOverviewOptions {
  t?: Translator
}

export function useOrderOverview({ t }: UseOrderOverviewOptions = {}) {
  const translate: Translator = t || ((key: string) => key)
  const router = useRouter()

  const loading = ref(false)
  const balance = ref(0)
  const mtdEstimate = ref<OrderOverviewMetric>({
    amount: 0,
    change: 0,
    lastMonth: 0
  })
  const lastMonthSettlement = ref<OrderSettlementMetric>({
    amount: 0,
    period: ''
  })
  const trendPeriod = ref('7')
  const rankingData = ref<OverviewRankingItem[]>([])
  const chartData = ref<ConsumptionTrendPoint[]>([])

  const fetchData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getOrderOverview()) as ApiResponse<OrderOverviewData>
      if (res.code === 0) {
        const data = res.data || {}
        balance.value = data.balance ?? 0
        mtdEstimate.value = {
          amount: data.mtdEstimate ?? 0,
          change: 0,
          lastMonth: data.lastMonthSettlement ?? 0
        }
        lastMonthSettlement.value = {
          amount: data.lastMonthSettlement ?? 0,
          period: translate('order.lastMonthSettlementDesc')
        }
        rankingData.value = data.monthlyRanking || []
        chartData.value = data.consumptionTrend || []
      }
    } catch (error) {
      console.error('获取财务概览失败:', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const handleRecharge = (): void => {
    // TODO: Implement recharge logic.
  }

  const goToTransactions = (): void => {
    void router.push({ name: 'transactions' })
  }

  const goToUsage = (): void => {
    void router.push({ name: 'order-usage' })
  }

  return {
    balance,
    chartData,
    fetchData,
    goToTransactions,
    goToUsage,
    handleRecharge,
    lastMonthSettlement,
    loading,
    mtdEstimate,
    rankingData,
    trendPeriod
  }
}

import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getOrderOverview } from '@/api/order'

export function useOrderOverview({ t }) {
  const translate = t || ((key) => key)
  const router = useRouter()

  const loading = ref(false)
  const balance = ref(0)
  const mtdEstimate = ref({
    amount: 0,
    change: 0,
    lastMonth: 0
  })
  const lastMonthSettlement = ref({
    amount: 0,
    period: ''
  })
  const trendPeriod = ref('7')
  const rankingData = ref([])
  const chartData = ref([])

  const fetchData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getOrderOverview()
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

  const handleRecharge = () => {
    // TODO: Implement recharge logic.
  }

  const goToTransactions = () => {
    router.push({ name: 'transactions' })
  }

  const goToUsage = () => {
    router.push({ name: 'order-usage' })
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

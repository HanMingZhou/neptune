<template>
  <div class="rounded-2xl border border-slate-200 bg-white p-8 dark:border-border-dark dark:bg-surface-dark">
    <div class="mb-8 flex items-center justify-between">
      <h3 class="text-sm font-bold uppercase tracking-widest text-slate-400">{{ t('consumptionTrend') }} (CNY)</h3>
      <el-select v-model="trendPeriodModel" size="small" class="!w-32">
        <el-option :label="t('last7Days')" value="7" />
        <el-option :label="t('last30Days')" value="30" />
      </el-select>
    </div>
    <div ref="chartRef" class="h-64"></div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, onUnmounted, ref, watch } from 'vue'
import * as echarts from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps({
  chartData: {
    type: Array,
    default: () => []
  },
  trendPeriod: {
    type: String,
    default: '7'
  }
})

const emit = defineEmits(['update:trend-period'])
const t = inject('t', (key) => key)

const chartRef = ref(null)
let chartInstance = null

const trendPeriodModel = computed({
  get: () => props.trendPeriod,
  set: (value) => emit('update:trend-period', value)
})

const visibleChartData = computed(() => {
  const limit = Number(props.trendPeriod)
  if (!Number.isFinite(limit) || limit <= 0) {
    return props.chartData
  }

  return props.chartData.slice(-limit)
})

const updateChart = () => {
  if (!chartInstance) {
    return
  }

  chartInstance.setOption({
    xAxis: {
      data: visibleChartData.value.map((item) => item.date)
    },
    series: [
      {
        data: visibleChartData.value.map((item) => item.amount)
      }
    ]
  })
}

const initChart = () => {
  if (!chartRef.value) {
    return
  }

  chartInstance = echarts.init(chartRef.value)
  chartInstance.setOption({
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: visibleChartData.value.map((item) => item.date),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: {
        color: '#94a3b8',
        fontSize: 10,
        fontWeight: 'bold'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: {
        lineStyle: {
          color: '#e2e8f0',
          type: 'dashed',
          opacity: 0.5
        }
      },
      axisLabel: {
        color: '#94a3b8',
        fontSize: 10,
        fontWeight: 'bold'
      }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#1d1d1d',
      borderWidth: 0,
      borderRadius: 12,
      textStyle: { color: '#fff' },
      formatter: (params) => {
        const item = params[0]
        return `<div class="p-2">
          <div class="text-xs text-slate-400 mb-1">${item.name}</div>
          <div class="text-sm font-bold">¥${item.value}</div>
        </div>`
      }
    },
    series: [
      {
        data: visibleChartData.value.map((item) => item.amount),
        type: 'bar',
        barWidth: 40,
        itemStyle: {
          borderRadius: [4, 4, 0, 0],
          color: (params) => {
            return params.dataIndex === visibleChartData.value.length - 1 ? '#165DFF' : 'rgba(22, 93, 255, 0.2)'
          }
        },
        emphasis: {
          itemStyle: {
            color: '#165DFF'
          }
        }
      }
    ]
  })
}

const handleResize = () => {
  chartInstance?.resize()
}

watch(visibleChartData, () => {
  updateChart()
}, { deep: true })

onMounted(() => {
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
  chartInstance = null
})
</script>

// 本组件参考 arco-pro 的实现
// https://github.com/arco-design/arco-design-pro-vue/blob/main/arco-design-pro-vite/src/hooks/chart-option.ts

import { computed, type ComputedRef } from 'vue'
import { useAppStore } from '@/pinia'
import type { EChartsOption } from 'echarts'

export default function useChartOption(sourceOption: (isDark: boolean) => EChartsOption): { chartOption: ComputedRef<EChartsOption> } {
    const appStore = useAppStore()
    const isDark = computed(() => {
        return appStore.isDark
    })
    const chartOption = computed(() => {
        return sourceOption(isDark.value)
    })
    return {
        chartOption
    }
}

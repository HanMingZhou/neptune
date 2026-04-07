import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardData } from '@/api/dashboard'
import type { Translator } from '@/types/consoleResource'
import type { ApiResponse } from '@/utils/request'
import type {
  DashboardData,
  DashboardQuickEntry,
  DashboardQuickEntryKey,
  DashboardRecentInstance,
  DashboardStatCard,
  DashboardStats,
  DashboardUsageTrendViewModel
} from '@/types/dashboard'

interface UseDashboardPageOptions {
  t?: Translator
}

const createDefaultStats = (): DashboardStats => ({
  runningNotebooks: 0,
  runningTraining: 0,
  runningInference: 0,
  totalNotebooks: 0,
  totalTraining: 0,
  totalInference: 0,
  totalGpuInUse: 0,
  storageUsed: 0,
  storageVolumeCount: 0
})

export function useDashboardPage({ t }: UseDashboardPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)
  const router = useRouter()

  const loading = ref(false)
  const stats = ref<DashboardStats>(createDefaultStats())
  const recentInstances = ref<DashboardRecentInstance[]>([])
  const usageTrends = ref<DashboardUsageTrendViewModel[]>([])

  const statCards = computed<DashboardStatCard[]>(() => {
    const currentStats = stats.value
    const totalRunning =
      currentStats.runningNotebooks +
      currentStats.runningTraining +
      currentStats.runningInference

    return [
      {
        key: 'runningInstances',
        titleKey: 'runningInstances',
        value: totalRunning,
        sub: `${currentStats.runningNotebooks} ${translate('notebook')} / ${currentStats.runningTraining} ${translate('training')} / ${currentStats.runningInference} ${translate('inference')}`,
        icon: 'computer',
        hoverClass: 'hover:border-emerald-500',
        iconContainerClass: 'bg-emerald-500/10',
        iconClass: 'text-emerald-500'
      },
      {
        key: 'gpuInUse',
        titleKey: 'gpuInUse',
        value: currentStats.totalGpuInUse,
        sub: translate('gpuInUseDesc'),
        icon: 'developer_board',
        hoverClass: 'hover:border-blue-500',
        iconContainerClass: 'bg-blue-500/10',
        iconClass: 'text-blue-500'
      },
      {
        key: 'storageUsed',
        titleKey: 'storageUsage',
        value: `${currentStats.storageUsed} GB`,
        sub: `${currentStats.storageVolumeCount} ${translate('volumes')}`,
        icon: 'folder',
        hoverClass: 'hover:border-purple-500',
        iconContainerClass: 'bg-purple-500/10',
        iconClass: 'text-purple-500'
      },
      {
        key: 'totalInstances',
        titleKey: 'totalInstances',
        value:
          currentStats.totalNotebooks +
          currentStats.totalTraining +
          currentStats.totalInference,
        sub: `${currentStats.totalNotebooks} ${translate('notebook')} / ${currentStats.totalTraining} ${translate('training')} / ${currentStats.totalInference} ${translate('inference')}`,
        icon: 'dns',
        hoverClass: 'hover:border-amber-500',
        iconContainerClass: 'bg-amber-500/10',
        iconClass: 'text-amber-500'
      }
    ]
  })

  const quickEntries = computed<DashboardQuickEntry[]>(() => [
    {
      key: 'notebook',
      icon: 'add_to_photos',
      labelKey: 'newNotebook'
    },
    {
      key: 'training',
      icon: 'model_training',
      labelKey: 'newTraining'
    },
    {
      key: 'storage',
      icon: 'folder',
      labelKey: 'storage'
    },
    {
      key: 'sshkeys',
      icon: 'vpn_key',
      labelKey: 'sshkeys'
    }
  ])

  const fetchData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getDashboardData({
        days: 7
      })) as ApiResponse<DashboardData>
      if (res.code === 0 && res.data) {
        const data = res.data
        stats.value = {
          ...createDefaultStats(),
          ...(data.stats ?? {})
        }
        recentInstances.value = data.recentInstances ?? []

        const trends = data.usageTrends ?? []
        const maxTasks = Math.max(
          ...trends.map((item) => item.runningTasks ?? 0),
          1
        )
        usageTrends.value = trends.map((item) => ({
          ...item,
          dateLabel: item.date ? item.date.slice(5) : '',
          barHeight: ((item.runningTasks ?? 0) / maxTasks) * 200
        }))
      }
    } catch (error: unknown) {
      console.error('Fetch dashboard data failed', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const getStatusClass = (status?: string): string => {
    const map: Record<string, string> = {
      Running: 'bg-emerald-500/10 text-emerald-500',
      RUNNING: 'bg-emerald-500/10 text-emerald-500',
      Stopped: 'bg-slate-500/10 text-slate-500',
      STOPPED: 'bg-slate-500/10 text-slate-500',
      Pending: 'bg-amber-500/10 text-amber-500',
      PENDING: 'bg-amber-500/10 text-amber-500',
      Failed: 'bg-red-500/10 text-red-500',
      FAILED: 'bg-red-500/10 text-red-500',
      Creating: 'bg-blue-500/10 text-blue-500',
      CREATING: 'bg-blue-500/10 text-blue-500',
      SUBMITTED: 'bg-blue-500/10 text-blue-500',
      SUCCEEDED: 'bg-emerald-500/10 text-emerald-500',
      KILLED: 'bg-slate-500/10 text-slate-500'
    }

    return map[status ?? ''] || 'bg-slate-500/10 text-slate-500'
  }

  const getTypeClass = (type?: string): string => {
    const map: Record<string, string> = {
      notebook: 'bg-blue-500/10 text-blue-600',
      training: 'bg-orange-500/10 text-orange-600',
      inference: 'bg-violet-500/10 text-violet-600'
    }

    return map[type ?? ''] || 'bg-slate-500/10 text-slate-500'
  }

  const goToDetail = (item: DashboardRecentInstance): void => {
    if (item.type === 'notebook') {
      void router.push({ name: 'notebookDetail', query: { id: item.id } })
      return
    }

    if (item.type === 'training') {
      void router.push({ name: 'trainingDetail', query: { id: item.id } })
      return
    }

    if (item.type === 'inference') {
      void router.push({ name: 'inferenceDetail', query: { id: item.id } })
    }
  }

  const goToCreate = (): void => {
    void router.push({ name: 'notebookCreate' })
  }

  const goToTrainingCreate = (): void => {
    void router.push({ name: 'trainingCreate' })
  }

  const goToNotebooks = (): void => {
    void router.push({ name: 'notebooks' })
  }

  const goToStorage = (): void => {
    void router.push({ name: 'storage' })
  }

  const goToSSHKeys = (): void => {
    void router.push({ name: 'sshkeys' })
  }

  const handleQuickEntry = (key: DashboardQuickEntryKey): void => {
    const actionMap: Record<DashboardQuickEntryKey, () => void> = {
      notebook: goToCreate,
      training: goToTrainingCreate,
      storage: goToStorage,
      sshkeys: goToSSHKeys
    }

    actionMap[key]?.()
  }

  return {
    fetchData,
    getStatusClass,
    getTypeClass,
    goToCreate,
    goToDetail,
    goToNotebooks,
    handleQuickEntry,
    loading,
    quickEntries,
    recentInstances,
    statCards,
    usageTrends
  }
}

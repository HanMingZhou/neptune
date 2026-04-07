import type { ResourceId } from './consoleResource'

export interface DashboardStats {
  runningNotebooks: number
  runningTraining: number
  runningInference: number
  totalNotebooks: number
  totalTraining: number
  totalInference: number
  totalGpuInUse: number
  storageUsed: number
  storageVolumeCount: number
}

export interface DashboardRecentInstance {
  id: ResourceId
  name?: string
  type?: string
  status?: string
  gpu?: number
  createdAt?: string
  [key: string]: unknown
}

export interface DashboardUsageTrend {
  date?: string
  runningTasks?: number
  [key: string]: unknown
}

export interface DashboardUsageTrendViewModel extends DashboardUsageTrend {
  dateLabel: string
  barHeight: number
}

export interface DashboardData {
  stats?: Partial<DashboardStats>
  recentInstances?: DashboardRecentInstance[]
  usageTrends?: DashboardUsageTrend[]
  [key: string]: unknown
}

export interface DashboardStatCard {
  key: string
  titleKey: string
  value: number | string
  sub: string
  icon: string
  hoverClass: string
  iconContainerClass: string
  iconClass: string
}

export type DashboardQuickEntryKey =
  | 'notebook'
  | 'training'
  | 'storage'
  | 'sshkeys'

export interface DashboardQuickEntry {
  key: DashboardQuickEntryKey
  icon: string
  labelKey: string
}

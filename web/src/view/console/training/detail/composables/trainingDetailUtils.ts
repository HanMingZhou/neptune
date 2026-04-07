import type {
  ConsolePod,
  ConsoleTrainingDetail,
  DetailTab,
  Translator
} from '@/types/consoleResource'

type TrainingJobLike = Partial<ConsoleTrainingDetail>

export const createTrainingDetailTabs = (t: Translator): DetailTab[] => [
  { key: 'overview', label: t('overview'), icon: 'description' },
  { key: 'logs', label: t('logs'), icon: 'article' },
  { key: 'terminal', label: t('terminal'), icon: 'terminal' },
  { key: 'pods', label: t('instanceList'), icon: 'dns' }
]

export const processTrainingLogData = (data = ''): string => {
  let processed = String(data).replace(/\r(?!\n)/g, '\n')
  processed = processed.replace(/\n{3,}/g, '\n\n')
  return processed
}

export const getDefaultTrainingTaskName = (
  job: TrainingJobLike = {}
): string => {
  if (job.frameworkType === 'MPI') {
    return 'mpimaster'
  }

  if (job.frameworkType === 'STANDALONE') {
    return 'worker'
  }

  return 'master'
}

export const getTrainingMasterTaskName = (job: TrainingJobLike = {}): string =>
  job.frameworkType === 'MPI' ? 'mpimaster' : 'master'

export const getTrainingMasterTaskLabel = (
  job: TrainingJobLike = {}
): string => (job.frameworkType === 'MPI' ? 'MPI Master' : 'Master')

export const getTrainingPodCount = (
  job: TrainingJobLike = {},
  taskName = ''
): number => {
  if (taskName === 'master' || taskName === 'mpimaster') {
    return 1
  }

  const count = Number(job.workerCount) || 1

  if (job.frameworkType !== 'MPI') {
    return Math.max(0, count - 1) || 1
  }

  return count
}

export const resolveTrainingLogTaskName = (
  job: TrainingJobLike = {},
  pod: Partial<ConsolePod> = {}
): string => {
  const podName = pod?.name || ''

  if (podName.includes('worker')) {
    return 'worker'
  }

  return getTrainingMasterTaskName(job)
}

export const getTrainingFrameworkLabel = (type?: string): string => {
  const map: Record<string, string> = {
    PYTORCH_DDP: 'PyTorch DDP',
    STANDALONE: 'StandAlone',
    MPI: 'MPI'
  }

  return map[type || ''] || type || '-'
}

export const getTrainingPayTypeLabel = (
  t: Translator,
  type?: number | string
): string => {
  const map: Record<number, string> = {
    1: 'payHourly',
    2: 'payDaily',
    3: 'payWeekly',
    4: 'payMonthly'
  }

  const normalizedType = Number(type)
  return t(map[normalizedType] || 'payHourly')
}

export const getTrainingStatusLabel = (
  t: Translator,
  status?: string
): string => t(status || '') || status || '-'

export const getTrainingStatusClass = (status?: string): string => {
  const map: Record<string, string> = {
    RUNNING: 'bg-emerald-500/10 text-emerald-500',
    PENDING: 'bg-amber-500/10 text-amber-500',
    CREATING: 'bg-amber-500/10 text-amber-500',
    SUCCEEDED: 'bg-blue-500/10 text-blue-500',
    FAILED: 'bg-red-500/10 text-red-500',
    KILLED: 'bg-slate-500/10 text-slate-500'
  }

  return map[status || ''] || 'bg-slate-500/10 text-slate-500'
}

export const getTrainingPodStatusClass = (status?: string): string => {
  const normalizedStatus = status?.toLowerCase()
  const map: Record<string, string> = {
    running: 'bg-emerald-500/10 text-emerald-500',
    pending: 'bg-amber-500/10 text-amber-500',
    failed: 'bg-red-500/10 text-red-500'
  }

  return map[normalizedStatus || ''] || 'bg-slate-500/10 text-slate-500'
}

export const formatTrainingTime = (
  time?: string | number | null
): string | null => {
  if (!time) {
    return null
  }

  return new Date(time).toLocaleString()
}

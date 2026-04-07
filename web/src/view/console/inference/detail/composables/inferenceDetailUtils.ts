import type {
  ConsoleInferenceDetail,
  DetailTab,
  Translator
} from '@/types/consoleResource'

type InferenceServiceLike = Partial<ConsoleInferenceDetail>

export const createInferenceDetailTabs = (t: Translator): DetailTab[] => [
  { key: 'overview', label: t('overview'), icon: 'description' },
  { key: 'logs', label: t('logs'), icon: 'article' },
  { key: 'terminal', label: t('terminal'), icon: 'terminal' },
  { key: 'pods', label: t('instanceList'), icon: 'dns' }
]

export const buildInferenceApiEndpoint = (
  service: InferenceServiceLike = {}
): string => {
  if (!service.accessUrl) {
    return ''
  }

  return `${service.accessUrl}/v1`
}

export const buildInferenceCurlExample = (
  service: InferenceServiceLike = {},
  apiEndpoint = ''
): string => {
  if (!apiEndpoint) {
    return ''
  }

  const authHeader =
    service.authType === 2
      ? '-H "API-Key: <YOUR_API_KEY>"'
      : '-H "Authorization: Bearer <JWT_TOKEN>"'

  return `curl ${apiEndpoint}/chat/completions \\
  ${authHeader} \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${service.modelPath || 'model'}",
    "messages": [{"role": "user", "content": "Hello"}],
    "max_tokens": 512
  }'`
}

export const processInferenceLogData = (data = ''): string => {
  let processed = String(data).replace(/\r(?!\n)/g, '\n')
  processed = processed.replace(/\n{3,}/g, '\n\n')
  return processed
}

export const formatInferenceCommand = (
  service: InferenceServiceLike = {}
): string => {
  const parts: string[] = []
  const command =
    typeof service.command === 'string' ? service.command.trim() : ''
  const toArray = (value: string[] | string | undefined): string[] =>
    Array.isArray(value)
      ? value
      : typeof value === 'string' && value
        ? [value]
        : []
  const args = toArray(service.args)
  const extraArgs = toArray(service.extraArgs)

  if (command) {
    parts.push(command)
  }

  if (args.length) {
    parts.push(args.join(' '))
  }

  if (extraArgs.length) {
    parts.push(extraArgs.join(' '))
  }

  return parts.join(' \\\n  ')
}

export const formatInferenceTime = (time?: string | number | null): string => {
  if (!time) {
    return '-'
  }

  return new Date(time).toLocaleString()
}

export const getInferenceStatusLabel = (
  t: Translator,
  status?: string
): string => t(status || '') || status || '-'

export const getInferenceStatusClass = (status?: string): string => {
  if (status === 'RUNNING') return 'bg-emerald-500/10 text-emerald-500'
  if (status === 'PENDING' || status === 'CREATING')
    return 'bg-amber-500/10 text-amber-500'
  if (status === 'FAILED') return 'bg-red-500/10 text-red-500'
  if (status === 'STOPPED') return 'bg-slate-500/10 text-slate-500'
  if (status === 'DELETING' || status === 'RESTARTING')
    return 'bg-blue-500/10 text-blue-500'
  return 'bg-slate-500/10 text-slate-500'
}

export const getInferencePodStatusClass = (status?: string): string => {
  const normalizedStatus = status?.toLowerCase()

  if (normalizedStatus === 'running')
    return 'bg-emerald-500/10 text-emerald-500'
  if (normalizedStatus === 'pending') return 'bg-amber-500/10 text-amber-500'
  if (normalizedStatus === 'failed') return 'bg-red-500/10 text-red-500'
  return 'bg-slate-500/10 text-slate-500'
}

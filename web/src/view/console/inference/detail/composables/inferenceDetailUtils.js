export const createInferenceDetailTabs = (t) => [
  { key: 'overview', label: t('overview'), icon: 'description' },
  { key: 'logs', label: t('logs'), icon: 'article' },
  { key: 'terminal', label: t('terminal'), icon: 'terminal' },
  { key: 'pods', label: t('instanceList'), icon: 'dns' }
]

export const buildInferenceApiEndpoint = (service = {}) => {
  if (!service.accessUrl) {
    return ''
  }

  return `${service.accessUrl}/v1`
}

export const buildInferenceCurlExample = (service = {}, apiEndpoint = '') => {
  if (!apiEndpoint) {
    return ''
  }

  const authHeader = service.authType === 2
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

export const processInferenceLogData = (data = '') => {
  let processed = String(data).replace(/\r(?!\n)/g, '\n')
  processed = processed.replace(/\n{3,}/g, '\n\n')
  return processed
}

export const formatInferenceCommand = (service = {}) => {
  const parts = []
  const command = typeof service.command === 'string' ? service.command.trim() : ''
  const toArray = (value) => Array.isArray(value) ? value : (typeof value === 'string' && value ? [value] : [])
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

export const formatInferenceTime = (time) => {
  if (!time) {
    return '-'
  }

  return new Date(time).toLocaleString()
}

export const getInferenceStatusLabel = (t, status) => t(status) || status

export const getInferenceStatusClass = (status) => {
  if (status === 'RUNNING') return 'bg-emerald-500/10 text-emerald-500'
  if (status === 'PENDING' || status === 'CREATING') return 'bg-amber-500/10 text-amber-500'
  if (status === 'FAILED') return 'bg-red-500/10 text-red-500'
  if (status === 'STOPPED') return 'bg-slate-500/10 text-slate-500'
  if (status === 'DELETING' || status === 'RESTARTING') return 'bg-blue-500/10 text-blue-500'
  return 'bg-slate-500/10 text-slate-500'
}

export const getInferencePodStatusClass = (status) => {
  const normalizedStatus = status?.toLowerCase()

  if (normalizedStatus === 'running') return 'bg-emerald-500/10 text-emerald-500'
  if (normalizedStatus === 'pending') return 'bg-amber-500/10 text-amber-500'
  if (normalizedStatus === 'failed') return 'bg-red-500/10 text-red-500'
  return 'bg-slate-500/10 text-slate-500'
}

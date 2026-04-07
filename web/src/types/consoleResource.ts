export type ResourceId = number | string
export type TranslationParams = Record<string, string | number>
export type Translator = (key: string, params?: TranslationParams) => string

export interface ListData<T> {
  list?: T[]
}

export interface PageListData<T> extends ListData<T> {
  total?: number
}

export interface FilterOption {
  model: string
  available?: number
}

export interface ProductFilterData {
  areas?: string[]
  gpuModels?: FilterOption[]
  cpuModels?: FilterOption[]
}

export interface ConsoleProduct {
  id: ResourceId
  clusterId: ResourceId
  clusterName?: string
  templateProductId?: ResourceId | null
  name?: string
  nodeType?: string
  area?: string
  cpu?: number
  cpuModel?: string
  memory?: number
  gpuCount: number
  gpuModel?: string
  gpuMemory?: number
  vGpuCount?: number
  vGpuMemory?: number
  vGpuCores?: number
  available?: number
  strictMax?: number
  balancedMax?: number
  systemDisk?: number
  cudaVersion?: string
  priceHourly?: number
  priceDaily?: number
  priceWeekly?: number
  priceMonthly?: number
  [key: string]: unknown
}

export interface ConsoleImage {
  id: ResourceId
  name: string
  description?: string
  [key: string]: unknown
}

export interface ConsoleVolume {
  id: ResourceId
  name: string
  size?: string
  clusterId?: ResourceId
  [key: string]: unknown
}

export interface ConsoleSshKey {
  id: ResourceId
  name?: string
  title?: string
  publicKey?: string
  isDefault?: boolean
  [key: string]: unknown
}

export interface ConsoleTrainingJob {
  id: ResourceId
  displayName?: string
  jobName?: string
  frameworkType?: string
  status?: string
  totalGpuCount?: number
  workerCount?: number
  createdAt?: string | number
  duration?: string
  enableTensorboard?: boolean
  tensorboardUrl?: string
  [key: string]: unknown
}

export interface ConsoleNotebook {
  id?: ResourceId
  displayName?: string
  instanceName?: string
  status?: string
  cpu?: number | string
  memory?: number | string
  gpuCount?: number
  gpu?: number
  gpuModel?: string
  jupyterUrl?: string
  enableTensorboard?: boolean
  tensorboardUrl?: string
  sshKeyCommand?: string
  sshCommand?: string
  sshPassword?: string
  [key: string]: unknown
}

export interface ConsoleInferenceService {
  id: ResourceId
  displayName?: string
  instanceName?: string
  framework?: string
  status?: string
  gpu?: number
  deployType?: string
  createdAt?: string | number
  [key: string]: unknown
}

export interface DetailTab {
  key: string
  label: string
  icon: string
}

export interface ConsolePod {
  name: string
  status?: string
  hostIP?: string
  podIP?: string
  [key: string]: unknown
}

export interface ConsoleEnvVar {
  name?: string
  value?: string
  [key: string]: unknown
}

export interface ConsoleTrainingMount {
  pvcName?: string
  mountPath?: string
  [key: string]: unknown
}

export interface ConsoleTrainingDetail extends ConsoleTrainingJob {
  cpu?: number | string
  memory?: number | string
  gpuModel?: string
  gpuType?: string
  imageName?: string
  image?: string
  workerGpu?: number
  clusterName?: string
  area?: string
  payType?: number
  price?: number
  startupCommand?: string
  startedAt?: string | number
  finishedAt?: string | number
  tensorboardLogPath?: string
  mounts?: ConsoleTrainingMount[]
  envs?: ConsoleEnvVar[]
  errorMsg?: string
}

export interface ConsoleNotebookVolumeMount {
  pvcName?: string
  mountsPath?: string
  mountPath?: string
  [key: string]: unknown
}

export interface ConsoleNotebookDetail extends ConsoleNotebook {
  price?: number | string
  payType?: number
  createdAt?: string | number
  imageName?: string
  volumeMounts?: ConsoleNotebookVolumeMount[] | string
  tensorboardLogPath?: string
}

export interface ConsoleInferenceMount {
  pvcName?: string
  mountPath?: string
  subPath?: string
  readOnly?: boolean
  [key: string]: unknown
}

export interface ConsoleInferenceDetail extends ConsoleInferenceService {
  cpu?: number | string
  memory?: number | string
  gpuModel?: string
  imageName?: string
  accessUrl?: string
  authType?: number
  modelMountPath?: string
  modelPath?: string
  servicePort?: number | string
  maxTokens?: number | string
  maxConcurrency?: number | string
  workerCount?: number
  autoRestart?: boolean
  restartCount?: number
  maxRestarts?: number
  command?: string
  args?: string[] | string
  extraArgs?: string[] | string
  mounts?: ConsoleInferenceMount[]
  envs?: ConsoleEnvVar[]
  errorMsg?: string
  startedAt?: string | number
}

export interface InferenceApiKey {
  id: ResourceId
  name?: string
  apiKey?: string
  createdAt?: string | number
  [key: string]: unknown
}

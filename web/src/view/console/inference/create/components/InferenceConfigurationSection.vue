<template>
  <div class="space-y-6">
    <EngineConfigSection
      :deploy-types="deployTypes"
      :field-errors="fieldErrors"
      :form="form"
      :framework-required="frameworkRequired"
      :frameworks="frameworks"
      :pvcs="pvcs"
      @deploy-type-change="onDeployTypeChange"
      @update:field="emit('update:field', $event)"
      @validate:display-name="emit('validate:display-name')"
    />

    <ImagePickerSection
      :active-tab="activeTab"
      :description-class="'text-sm text-slate-500 mb-4 italic'"
      :description-text="imageDescription"
      :empty-state-class="'col-span-full text-center py-8 text-slate-400 font-mono'"
      :items="imageOptions"
      :selected-value="form.imageId"
      :show-hot-tag="true"
      :tabs="imageTabs"
      :uppercase-label="true"
      @change-tab="changeImageTab"
      @update:selected-value="
        emit('update:field', { key: 'imageId', value: $event })
      "
    />

    <ServiceAdvancedSection
      :auth-types="authTypes"
      :form="form"
      :pvcs="pvcs"
      @add-env="addEnv"
      @add-mount="addMount"
      @remove-env="removeEnv"
      @remove-mount="removeMount"
      @update:env="emit('update:env', $event)"
      @update:field="emit('update:field', $event)"
      @update:mount="emit('update:mount', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type {
  ConsoleVolume,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import EngineConfigSection from './EngineConfigSection.vue'
import ServiceAdvancedSection from './ServiceAdvancedSection.vue'
import type {
  InferenceAuthType,
  InferenceCreateForm,
  InferenceDeployType,
  InferenceEnv,
  InferenceFieldErrors,
  InferenceFramework,
  InferenceImageTab,
  InferenceMount,
  UpdateInferenceFieldPayload
} from '../composables/useInferenceCreate'

interface AuthTypeOption {
  label: string
  value: InferenceAuthType
}

interface DeployTypeOption {
  label: string
  value: InferenceDeployType
}

interface FrameworkOption {
  icon: string
  label: string
  value: Exclude<InferenceFramework, ''>
}

interface ImageOption {
  desc?: string
  label: string
  value: ResourceId
}

interface ImageTabOption {
  value: InferenceImageTab
  label: string
}

interface InferenceEnvUpdatePayload {
  index: number
  key: keyof InferenceEnv
  value: InferenceEnv[keyof InferenceEnv]
}

interface InferenceMountUpdatePayload {
  index: number
  key: keyof InferenceMount
  value: InferenceMount[keyof InferenceMount]
}

const props = withDefaults(
  defineProps<{
    activeTab: InferenceImageTab
    addEnv: () => void
    addMount: () => void
    authTypes?: AuthTypeOption[]
    changeImageTab: (tab: InferenceImageTab) => void
    deployTypes?: DeployTypeOption[]
    fieldErrors?: InferenceFieldErrors
    form: InferenceCreateForm
    frameworkRequired?: boolean
    frameworks?: FrameworkOption[]
    imageOptions?: ImageOption[]
    imageTabs?: ImageTabOption[]
    onDeployTypeChange: (value: InferenceDeployType) => void
    pvcs?: ConsoleVolume[]
    removeEnv: (index: number) => void
    removeMount: (index: number) => void
  }>(),
  {
    authTypes: () => [],
    deployTypes: () => [],
    fieldErrors: () => ({
      displayName: ''
    }),
    frameworkRequired: false,
    frameworks: () => [],
    imageOptions: () => [],
    imageTabs: () => [],
    pvcs: () => []
  }
)

const emit = defineEmits<{
  'update:env': [payload: InferenceEnvUpdatePayload]
  'update:field': [payload: UpdateInferenceFieldPayload]
  'update:mount': [payload: InferenceMountUpdatePayload]
  'validate:display-name': []
}>()

const t = inject<Translator>('t', (key: string) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

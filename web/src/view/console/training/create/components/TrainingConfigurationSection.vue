<template>
  <div class="space-y-6">
    <ImagePickerSection
      :active-tab="activeTab"
      :description-class="'text-xs text-slate-500 mb-4 italic'"
      :description-text="imageDescription"
      :items="filteredImages"
      :label-key="'name'"
      :selected-value="form.imageId"
      :tabs="imageTabs"
      :value-key="'id'"
      @change-tab="changeImageTab"
      @update:selected-value="form.imageId = $event"
    />

    <JobConfigSection
      :available-capacity="availableCapacity"
      :field-errors="fieldErrors"
      :form="form"
      :framework-types="frameworkTypes"
      :max-worker-count="maxWorkerCount"
      :pvcs="pvcs"
      :selected-product="selectedProduct"
      :show-worker-count="showWorkerCount"
      @add-env="addEnv"
      @add-mount="addMount"
      @decrease-worker="decreaseWorker"
      @increase-worker="increaseWorker"
      @insert-mpi-example="insertMpiExample"
      @insert-pytorch-example="insertPytorchExample"
      @remove-env="removeEnv"
      @remove-mount="removeMount"
      @update:env="form.envs[$event.index][$event.key] = $event.value"
      @update:field="updateField($event)"
      @update:mount="form.mounts[$event.index][$event.key] = $event.value"
      @update:schedule-strategy="form.scheduleStrategy = $event || 'BALANCED'"
      @update:workerCount="
        form.workerCount = Math.max(2, Math.min($event || 2, maxWorkerCount))
      "
      @validate:name="emit('validate:name')"
      @validate:tensorboard-log-path="emit('validate:tensorboard-log-path')"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type {
  ConsoleImage,
  ConsoleProduct,
  ConsoleVolume,
  Translator
} from '@/types/consoleResource'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import JobConfigSection from './JobConfigSection.vue'
import type {
  TrainingCreateForm,
  TrainingFieldErrors,
  TrainingFrameworkType,
  TrainingImageTab,
  UpdateTrainingFieldPayload
} from '../composables/useTrainingCreate'

interface TrainingFrameworkOption {
  value: TrainingFrameworkType
  label: string
  hint?: string
}

interface TrainingImageTabOption {
  value: TrainingImageTab
  label: string
}

const props = withDefaults(
  defineProps<{
    activeTab: TrainingImageTab
    availableCapacity?: number
    addEnv: () => void
    addMount: () => void
    changeImageTab: (tab: TrainingImageTab) => void
    decreaseWorker: () => void
    filteredImages?: ConsoleImage[]
    fieldErrors?: TrainingFieldErrors
    form: TrainingCreateForm
    frameworkTypes?: TrainingFrameworkOption[]
    imageTabs?: TrainingImageTabOption[]
    increaseWorker: () => void
    insertMpiExample: () => void
    insertPytorchExample: () => void
    maxWorkerCount?: number
    pvcs?: ConsoleVolume[]
    removeEnv: (index: number) => void
    removeMount: (index: number) => void
    selectedProduct?: ConsoleProduct | null
    showWorkerCount?: boolean
    updateField: (payload: UpdateTrainingFieldPayload) => void
  }>(),
  {
    availableCapacity: 0,
    filteredImages: () => [],
    fieldErrors: () => ({
      name: '',
      tensorboardLogPath: ''
    }),
    frameworkTypes: () => [],
    imageTabs: () => [],
    maxWorkerCount: 2,
    pvcs: () => [],
    selectedProduct: null,
    showWorkerCount: false
  }
)

const emit = defineEmits<{
  'validate:name': []
  'validate:tensorboard-log-path': []
}>()

const t = inject<Translator>('t', (key: string) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

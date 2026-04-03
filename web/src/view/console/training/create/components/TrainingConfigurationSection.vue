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
      @update:workerCount="form.workerCount = Math.max(2, Math.min($event || 2, maxWorkerCount))"
      @validate:name="emit('validate:name')"
      @validate:tensorboard-log-path="emit('validate:tensorboard-log-path')"
    />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import JobConfigSection from './JobConfigSection.vue'

const props = defineProps({
  activeTab: {
    type: String,
    required: true
  },
  availableCapacity: {
    type: Number,
    default: 0
  },
  addEnv: {
    type: Function,
    required: true
  },
  addMount: {
    type: Function,
    required: true
  },
  changeImageTab: {
    type: Function,
    required: true
  },
  decreaseWorker: {
    type: Function,
    required: true
  },
  filteredImages: {
    type: Array,
    default: () => []
  },
  fieldErrors: {
    type: Object,
    default: () => ({})
  },
  form: {
    type: Object,
    required: true
  },
  frameworkTypes: {
    type: Array,
    default: () => []
  },
  imageTabs: {
    type: Array,
    default: () => []
  },
  insertMpiExample: {
    type: Function,
    required: true
  },
  insertPytorchExample: {
    type: Function,
    required: true
  },
  increaseWorker: {
    type: Function,
    required: true
  },
  maxWorkerCount: {
    type: Number,
    default: 2
  },
  pvcs: {
    type: Array,
    default: () => []
  },
  removeEnv: {
    type: Function,
    required: true
  },
  removeMount: {
    type: Function,
    required: true
  },
  selectedProduct: {
    type: Object,
    default: null
  },
  showWorkerCount: {
    type: Boolean,
    default: false
  },
  updateField: {
    type: Function,
    required: true
  }
})

const emit = defineEmits(['validate:name', 'validate:tensorboard-log-path'])

const t = inject('t', (key) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

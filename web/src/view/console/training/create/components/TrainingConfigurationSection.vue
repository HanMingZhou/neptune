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
      :form="form"
      :framework-types="frameworkTypes"
      :pvcs="pvcs"
      @add-env="addEnv"
      @add-mount="addMount"
      @insert-mpi-example="insertMpiExample"
      @insert-pytorch-example="insertPytorchExample"
      @remove-env="removeEnv"
      @remove-mount="removeMount"
      @update:env="form.envs[$event.index][$event.key] = $event.value"
      @update:field="form[$event.key] = $event.value"
      @update:mount="form.mounts[$event.index][$event.key] = $event.value"
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
  filteredImages: {
    type: Array,
    default: () => []
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
  }
})

const t = inject('t', (key) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

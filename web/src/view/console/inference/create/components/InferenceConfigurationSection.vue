<template>
  <div class="space-y-6">
    <EngineConfigSection
      :deploy-types="deployTypes"
      :form="form"
      :framework-required="frameworkRequired"
      :frameworks="frameworks"
      :pvcs="pvcs"
      @deploy-type-change="onDeployTypeChange"
      @update:field="emit('update:field', $event)"
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
      @update:selected-value="emit('update:field', { key: 'imageId', value: $event })"
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

<script setup>
import { computed, inject } from 'vue'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import EngineConfigSection from './EngineConfigSection.vue'
import ServiceAdvancedSection from './ServiceAdvancedSection.vue'

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
  authTypes: {
    type: Array,
    default: () => []
  },
  changeImageTab: {
    type: Function,
    required: true
  },
  deployTypes: {
    type: Array,
    default: () => []
  },
  form: {
    type: Object,
    required: true
  },
  frameworkRequired: {
    type: Boolean,
    default: false
  },
  frameworks: {
    type: Array,
    default: () => []
  },
  imageOptions: {
    type: Array,
    default: () => []
  },
  imageTabs: {
    type: Array,
    default: () => []
  },
  onDeployTypeChange: {
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

const emit = defineEmits(['update:env', 'update:field', 'update:mount'])

const t = inject('t', (key) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

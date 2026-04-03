<template>
  <div class="space-y-6">
    <ImagePickerSection
      :active-tab="activeTab"
      :description-text="imageDescription"
      :items="filteredImages"
      :label-key="'name'"
      :selected-value="selectedImage"
      :tabs="imageTabs"
      :value-key="'id'"
      @change-tab="changeImageTab"
      @update:selected-value="emit('update:selected-image', $event)"
    />

    <OtherConfigSection
      :enable-ssh-password="enableSshPassword"
      :enable-tensorboard="enableTensorboard"
      :instance-name="instanceName"
      :instance-name-error="fieldErrors.instanceName"
      :selected-ssh-key="selectedSshKey"
      :ssh-keys="sshKeys"
      :tensorboard-log-path="tensorboardLogPath"
      :tensorboard-path-error="fieldErrors.tensorboardLogPath"
      @update:enable-ssh-password="emit('update:enable-ssh-password', $event)"
      @update:enable-tensorboard="emit('update:enable-tensorboard', $event)"
      @update:instance-name="emit('update:instance-name', $event)"
      @update:selected-ssh-key="emit('update:selected-ssh-key', $event)"
      @update:tensorboard-log-path="emit('update:tensorboard-log-path', $event)"
      @validate:instance-name="emit('validate:instance-name')"
      @validate:tensorboard-log-path="emit('validate:tensorboard-log-path')"
    />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import OtherConfigSection from './OtherConfigSection.vue'

const props = defineProps({
  activeTab: {
    type: String,
    required: true
  },
  changeImageTab: {
    type: Function,
    required: true
  },
  enableSshPassword: {
    type: Boolean,
    default: false
  },
  enableTensorboard: {
    type: Boolean,
    default: false
  },
  filteredImages: {
    type: Array,
    default: () => []
  },
  fieldErrors: {
    type: Object,
    default: () => ({})
  },
  imageTabs: {
    type: Array,
    default: () => []
  },
  instanceName: {
    type: String,
    default: ''
  },
  selectedImage: {
    type: [Number, String],
    default: null
  },
  selectedSshKey: {
    type: [Number, String],
    default: null
  },
  sshKeys: {
    type: Array,
    default: () => []
  },
  tensorboardLogPath: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'update:enable-ssh-password',
  'update:enable-tensorboard',
  'update:instance-name',
  'update:selected-image',
  'update:selected-ssh-key',
  'update:tensorboard-log-path',
  'validate:instance-name',
  'validate:tensorboard-log-path'
])

const t = inject('t', (key) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

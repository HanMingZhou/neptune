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

<script setup lang="ts">
import { computed, inject } from 'vue'
import type {
  ConsoleImage,
  ConsoleSshKey,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import ImagePickerSection from '@/components/createPage/ImagePickerSection.vue'
import OtherConfigSection from './OtherConfigSection.vue'
import type {
  NotebookFieldErrors,
  NotebookImageTab
} from '../composables/useNotebookCreate'

interface NotebookImageTabOption {
  value: NotebookImageTab
  label: string
}

type SelectableId = ResourceId | '' | null

const props = withDefaults(
  defineProps<{
    activeTab: NotebookImageTab
    changeImageTab: (tab: NotebookImageTab) => void
    enableSshPassword?: boolean
    enableTensorboard?: boolean
    fieldErrors?: NotebookFieldErrors
    filteredImages?: ConsoleImage[]
    imageTabs?: NotebookImageTabOption[]
    instanceName?: string
    selectedImage?: SelectableId
    selectedSshKey?: ResourceId | null
    sshKeys?: ConsoleSshKey[]
    tensorboardLogPath?: string
  }>(),
  {
    enableSshPassword: false,
    enableTensorboard: false,
    fieldErrors: () => ({
      instanceName: '',
      tensorboardLogPath: ''
    }),
    filteredImages: () => [],
    imageTabs: () => [],
    instanceName: '',
    selectedImage: null,
    selectedSshKey: null,
    sshKeys: () => [],
    tensorboardLogPath: ''
  }
)

const emit = defineEmits<{
  'update:enable-ssh-password': [value: boolean]
  'update:enable-tensorboard': [value: boolean]
  'update:instance-name': [value: string]
  'update:selected-image': [value: ResourceId]
  'update:selected-ssh-key': [value: ResourceId | null]
  'update:tensorboard-log-path': [value: string]
  'validate:instance-name': []
  'validate:tensorboard-log-path': []
}>()

const t = inject<Translator>('t', (key: string) => key)
const imageDescription = computed(() => t(`${props.activeTab}ImageDesc`))
</script>

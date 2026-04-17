<template>
  <div
    class="console-create-card console-create-card--section w-full"
  >
    <h3 class="console-create-card__title mb-4">
      <span class="console-create-card__title-marker"></span>
      {{ t('mountExistingData') }}
    </h3>
    <el-select
      :model-value="selectedVolumeId"
      :placeholder="t('selectDataDisk')"
      clearable
      class="w-full"
      @update:model-value="$emit('update:selectedVolumeId', $event)"
      @change="$emit('volume-change', $event)"
    >
      <template #empty>
        <div class="p-4 text-center text-slate-400">{{ t('noData') }}</div>
      </template>
      <el-option
        v-for="volume in availableVolumes"
        :key="volume.id"
        :label="`${volume.name} (${volume.size}Gi) - ${volume.type === 1 ? 'Dataset' : 'Model'}`"
        :value="volume.id"
      />
    </el-select>
    <div v-if="selectedVolumeId" class="mt-4">
      <label class="block text-sm text-slate-500 mb-2">{{
        t('mountPath')
      }}</label>
      <input
        :value="volumeMountPath"
        :placeholder="t('enterMountPath')"
        class="create-form-input"
        type="text"
        @input="$emit('update:volumeMountPath', $event.target.value)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type {
  ConsoleVolume,
  ResourceId,
  Translator
} from '@/types/consoleResource'

withDefaults(
  defineProps<{
    availableVolumes?: ConsoleVolume[]
    selectedVolumeId?: ResourceId | null
    volumeMountPath?: string
  }>(),
  {
    availableVolumes: () => [],
    selectedVolumeId: null,
    volumeMountPath: ''
  }
)

defineEmits<{
  'update:selectedVolumeId': [value: ResourceId | null]
  'update:volumeMountPath': [value: string]
  'volume-change': [value: ResourceId | null]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

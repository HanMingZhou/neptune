<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
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
      <label class="block text-sm text-slate-500 mb-2">{{ t('mountPath') }}</label>
      <input
        :value="volumeMountPath"
        :placeholder="t('enterMountPath')"
        class="w-full px-4 py-2 border border-border-light dark:border-border-dark rounded-lg text-sm bg-white dark:bg-zinc-800 focus:ring-1 focus:ring-primary outline-none"
        type="text"
        @input="$emit('update:volumeMountPath', $event.target.value)"
      />
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  availableVolumes: {
    type: Array,
    default: () => []
  },
  selectedVolumeId: {
    type: [Number, String],
    default: null
  },
  volumeMountPath: {
    type: String,
    default: ''
  }
})

defineEmits(['update:selectedVolumeId', 'update:volumeMountPath', 'volume-change'])

const t = inject('t', (key) => key)
</script>

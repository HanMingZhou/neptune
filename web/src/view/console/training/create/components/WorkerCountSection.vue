<template>
  <div v-if="showWorkerCount" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('workerNodes') }}
    </h3>
    <div>
      <label class="block text-sm text-slate-500 mb-2">
        {{ t('workerCount') }} ({{ t('workerCountHint') }})<span class="text-red-500">*</span>
      </label>
      <div class="flex items-center">
        <button
          :disabled="workerCount <= 2"
          class="w-9 h-9 border border-border-light dark:border-border-dark rounded-l-lg hover:bg-slate-50 dark:hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="$emit('decrease-worker')"
        >
          -
        </button>
        <input
          :max="maxWorkerCount"
          :value="workerCount"
          class="w-16 h-9 text-center border-y border-border-light dark:border-border-dark bg-white dark:bg-zinc-800 outline-none"
          min="2"
          type="number"
          @input="$emit('update:workerCount', Number($event.target.value))"
        />
        <button
          :disabled="workerCount >= maxWorkerCount"
          class="w-9 h-9 border border-border-light dark:border-border-dark rounded-r-lg hover:bg-slate-50 dark:hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="$emit('increase-worker')"
        >
          +
        </button>
      </div>
      <p v-if="selectedProduct" class="mt-2 text-sm text-slate-500">
        {{ t('availableCapacity', { count: availableCapacity, unit: selectedProduct.gpuCount > 0 ? 'GPU' : t('instance') }) }}
      </p>
      <div v-if="frameworkType === 'MPI'" class="mt-2 flex items-center gap-2 text-sm text-amber-600 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg">
        <span class="material-icons text-base">info</span>
        {{ t('mpiModeHint') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  availableCapacity: {
    type: Number,
    default: 0
  },
  frameworkType: {
    type: String,
    required: true
  },
  maxWorkerCount: {
    type: Number,
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
  workerCount: {
    type: Number,
    required: true
  }
})

defineEmits(['decrease-worker', 'increase-worker', 'update:workerCount'])

const t = inject('t', (key) => key)
</script>

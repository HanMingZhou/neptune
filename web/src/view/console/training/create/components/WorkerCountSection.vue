<template>
  <div
    v-if="showWorkerCount"
    :class="embedded
      ? 'rounded-xl border border-blue-200/70 bg-blue-50/70 p-4 dark:border-blue-800/60 dark:bg-blue-950/20'
      : 'bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6'"
  >
    <h3 v-if="!embedded" class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('workerNodes') }}
    </h3>
    <div v-else class="mb-4 flex items-center gap-2 text-sm font-semibold text-slate-700 dark:text-slate-200">
      <span class="material-icons text-[18px] text-primary">tune</span>
      {{ t('workerNodes') }}
    </div>
    <div>
      <label class="block text-sm text-slate-500 mb-2">
        {{ t('workerCount') }} ({{ t('workerCountHint') }})<span class="text-red-500">*</span>
      </label>
      <div class="create-form-stepper">
        <button
          :disabled="workerCount <= 2"
          class="create-form-stepper__button"
          type="button"
          @click="$emit('decrease-worker')"
        >
          -
        </button>
        <input
          :max="maxWorkerCount"
          :value="workerCount"
          class="create-form-stepper__input"
          min="2"
          type="number"
          @input="$emit('update:workerCount', Number($event.target.value))"
        />
        <button
          :disabled="workerCount >= maxWorkerCount"
          class="create-form-stepper__button"
          type="button"
          @click="$emit('increase-worker')"
        >
          +
        </button>
      </div>
      <div class="mt-4">
        <label class="block text-sm text-slate-500 mb-2">调度策略</label>
        <el-select
          :model-value="scheduleStrategy"
          class="w-full"
          @update:model-value="$emit('update:schedule-strategy', $event)"
        >
          <el-option label="智能均衡（默认）" value="BALANCED" />
          <el-option label="严格分布式（节点不足则失败）" value="STRICT" />
        </el-select>
      </div>
      <p v-if="!selectedProduct" class="mt-2 text-sm text-slate-500">
        请先选择资源规格，再根据可用节点数调整实例数量
      </p>
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
  embedded: {
    type: Boolean,
    default: false
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
  scheduleStrategy: {
    type: String,
    default: 'BALANCED'
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

defineEmits(['decrease-worker', 'increase-worker', 'update:schedule-strategy', 'update:workerCount'])

const t = inject('t', (key) => key)
</script>

<template>
  <div class="fixed bottom-0 left-0 right-0 bg-surface-light/95 dark:bg-surface-dark/95 backdrop-blur border-t border-border-light dark:border-border-dark py-4 z-50">
    <div :class="['max-w-[1400px] mx-auto px-6 flex justify-end items-center', actionGapClass]">
      <div class="flex items-baseline gap-2">
        <span :class="priceLabelClass">{{ t(priceLabelKey) }}:</span>
        <span :class="priceValueClass">¥{{ totalPrice }}</span>
        <span :class="unitLabelClass">/{{ priceUnitText }}</span>
      </div>
      <div class="flex gap-3">
        <button
          class="px-6 py-2.5 border border-border-light dark:border-border-dark rounded-lg text-sm font-bold hover:bg-slate-50 dark:hover:bg-zinc-800 transition-colors"
          @click="$emit('back')"
        >
          {{ t(cancelLabelKey) }}
        </button>
        <button
          :class="[
            submitBaseClass,
            canSubmit && !loading ? enabledSubmitClass : disabledSubmitClass
          ]"
          :disabled="!canSubmit || loading"
          @click="$emit('submit')"
        >
          <span v-if="showSpinner && loading" class="material-icons animate-spin text-sm mr-1">autorenew</span>
          {{ loading ? t(loadingLabelKey) : t(submitLabelKey) }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  actionGapClass: {
    type: String,
    default: 'gap-6'
  },
  canSubmit: {
    type: Boolean,
    default: false
  },
  cancelLabelKey: {
    type: String,
    default: 'cancel'
  },
  disabledSubmitClass: {
    type: String,
    default: 'bg-slate-200 dark:bg-zinc-700 text-slate-400 cursor-not-allowed'
  },
  enabledSubmitClass: {
    type: String,
    default: 'bg-primary hover:bg-primary-hover text-white shadow-lg shadow-primary/20'
  },
  loading: {
    type: Boolean,
    default: false
  },
  loadingLabelKey: {
    type: String,
    default: 'submitting'
  },
  priceLabelClass: {
    type: String,
    default: 'text-sm text-slate-500'
  },
  priceLabelKey: {
    type: String,
    default: 'totalPrice'
  },
  priceUnitText: {
    type: String,
    required: true
  },
  priceValueClass: {
    type: String,
    default: 'text-2xl font-bold text-red-500'
  },
  showSpinner: {
    type: Boolean,
    default: false
  },
  submitBaseClass: {
    type: String,
    default: 'px-6 py-2.5 rounded-lg text-sm font-bold transition-all'
  },
  submitLabelKey: {
    type: String,
    default: 'createNow'
  },
  totalPrice: {
    type: String,
    required: true
  },
  unitLabelClass: {
    type: String,
    default: 'text-sm text-slate-400'
  }
})

defineEmits(['back', 'submit'])

const t = inject('t', (key) => key)
</script>

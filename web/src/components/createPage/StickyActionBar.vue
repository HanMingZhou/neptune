<template>
  <div
    class="fixed bottom-0 left-0 right-0 border-t border-border-light bg-surface-light py-4 z-50 dark:border-border-dark dark:bg-background-dark"
  >
    <div
      :class="[
        'console-page-container flex justify-end items-center',
        actionGapClass
      ]"
    >
      <div class="flex items-baseline gap-2">
        <span :class="priceLabelClass">{{ t(priceLabelKey) }}:</span>
        <span :class="priceValueClass">¥{{ totalPrice }}</span>
        <span :class="unitLabelClass">/{{ priceUnitText }}</span>
      </div>
      <div class="flex gap-3">
        <button
          class="px-6 py-2.5 border border-border-light bg-surface-light text-sm font-bold transition-colors hover:bg-slate-50 dark:border-border-dark dark:bg-surface-dark dark:hover:bg-zinc-800"
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
          <span
            v-if="showSpinner && loading"
            class="material-icons animate-spin text-sm mr-1"
            >autorenew</span
          >
          {{ loading ? t(loadingLabelKey) : t(submitLabelKey) }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    actionGapClass?: string
    canSubmit?: boolean
    cancelLabelKey?: string
    disabledSubmitClass?: string
    enabledSubmitClass?: string
    loading?: boolean
    loadingLabelKey?: string
    priceLabelClass?: string
    priceLabelKey?: string
    priceUnitText: string
    priceValueClass?: string
    showSpinner?: boolean
    submitBaseClass?: string
    submitLabelKey?: string
    totalPrice: string | number
    unitLabelClass?: string
  }>(),
  {
    actionGapClass: 'gap-6',
    canSubmit: false,
    cancelLabelKey: 'cancel',
    disabledSubmitClass:
      'bg-slate-200 dark:bg-zinc-700 text-slate-400 cursor-not-allowed',
    enabledSubmitClass: 'bg-primary hover:bg-primary-hover text-white',
    loading: false,
    loadingLabelKey: 'submitting',
    priceLabelClass: 'text-sm text-slate-500',
    priceLabelKey: 'totalPrice',
    priceValueClass: 'text-2xl font-bold text-red-500',
    showSpinner: false,
    submitBaseClass: 'px-6 py-2.5 text-sm font-bold transition-all',
    submitLabelKey: 'createNow',
    unitLabelClass: 'text-sm text-slate-400'
  }
)

defineEmits<{
  back: []
  submit: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

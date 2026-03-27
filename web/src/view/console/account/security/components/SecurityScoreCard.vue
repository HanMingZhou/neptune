<template>
  <div class="p-6 bg-primary/10 border border-primary/20 rounded-xl">
    <div class="flex items-center gap-3 mb-4">
      <span class="material-icons text-primary">verified_user</span>
      <h3 class="font-bold text-sm">
        {{ t('security.scoreTitle') }}：{{ levelLabel }}
      </h3>
    </div>
    <div class="w-full bg-slate-200 dark:bg-zinc-800 h-1.5 rounded-full overflow-hidden mb-4">
      <div class="bg-primary h-full transition-all duration-500" :style="{ width: `${score}%` }"></div>
    </div>
    <p class="text-[11px] text-slate-600 dark:text-slate-400 leading-relaxed">
      {{ t('security.scoreExceed') }}
    </p>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  score: {
    type: Number,
    default: 0
  }
})

const t = inject('t', (key) => key)

const levelLabel = computed(() => {
  if (props.score >= 80) {
    return t('security.excellent')
  }

  if (props.score >= 60) {
    return t('security.good')
  }

  return t('security.risk')
})
</script>

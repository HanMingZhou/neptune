<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-5 shadow-sm">
    <h3 class="text-xs font-black uppercase text-slate-400 mb-4">{{ t('security.akTitle') }}</h3>
    <p class="text-xs text-slate-500 mb-4">{{ t('security.akDesc') }}</p>
    <div class="space-y-3">
      <div class="p-3 bg-slate-50 dark:bg-zinc-800 rounded-lg flex items-center justify-between border border-border-light dark:border-border-dark">
        <div>
          <p class="text-[10px] font-bold text-slate-400">Access Key ID</p>
          <p class="text-xs font-mono">{{ accessKeyId || t('security.notGenerated') }}</p>
        </div>
        <button
          class="material-icons text-slate-400 text-sm hover:text-primary disabled:opacity-40"
          :disabled="!accessKeyId"
          @click="$emit('copy')"
        >
          {{ t('copy') }}
        </button>
      </div>
      <button
        :disabled="loading"
        class="w-full py-2 bg-slate-900 text-white dark:bg-white dark:text-slate-900 text-xs font-bold rounded-lg hover:opacity-90 transition-opacity flex items-center justify-center gap-2"
        @click="$emit('generate')"
      >
        <div
          v-if="loading"
          class="size-3 border-2 border-slate-400 border-t-white dark:border-slate-300 dark:border-t-slate-900 rounded-full animate-spin"
        ></div>
        {{ t('security.manageAk') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  accessKeyId: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['generate', 'copy'])

const t = inject('t', (key) => key)
</script>

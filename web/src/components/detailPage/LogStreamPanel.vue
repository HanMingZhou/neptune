<template>
  <div class="h-[calc(100vh-220px)] flex flex-col">
    <div class="console-detail-card rounded-xl overflow-hidden flex flex-col h-full">
      <div class="px-4 py-3 border-b border-border-light dark:border-border-dark" :class="toolbarClass">
        <slot name="controls-prefix" />
        <button
          v-if="!logsConnected"
          :disabled="logsLoading"
          class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-bold flex items-center gap-2 disabled:opacity-50"
          @click="$emit('connect')"
        >
          <span class="material-icons text-lg">link</span>
          {{ logsLoading ? t('connecting') : t('connectLogStream') }}
        </button>
        <button
          v-else
          class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white rounded-lg text-sm font-bold flex items-center gap-2"
          @click="$emit('disconnect')"
        >
          <span class="material-icons text-lg">link_off</span>
          {{ t('disconnectLogStream') }}
        </button>
        <span v-if="logsConnected" class="flex items-center gap-2 text-sm text-emerald-500">
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          {{ t('streaming') }}
        </span>
        <button :class="clearButtonClass" @click="$emit('clear')">
          <span class="material-icons">delete_outline</span>
        </button>
        <slot name="controls-suffix" />
      </div>
      <div :ref="setLogsRef" class="flex-1 bg-zinc-900 p-4 overflow-y-auto custom-scrollbar">
        <slot name="content">
          <pre class="text-slate-300 text-sm font-mono whitespace-pre-wrap break-all leading-6">{{ logs || t('connectLogStream') }}</pre>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  clearButtonClass: {
    type: String,
    default: 'ml-auto text-slate-400 hover:text-slate-600 dark:hover:text-slate-300'
  },
  logs: {
    type: String,
    default: ''
  },
  logsConnected: {
    type: Boolean,
    default: false
  },
  logsLoading: {
    type: Boolean,
    default: false
  },
  setLogsRef: {
    type: Function,
    required: true
  },
  toolbarClass: {
    type: String,
    default: 'flex items-center gap-4'
  }
})

defineEmits(['clear', 'connect', 'disconnect'])

const t = inject('t', (key) => key)
</script>

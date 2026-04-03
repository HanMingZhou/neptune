<template>
  <div class="h-[calc(100vh-220px)] flex flex-col gap-4">
    <div class="console-detail-card rounded-xl px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <slot name="controls-prefix" />
        <button
          v-if="!terminalConnected"
          :disabled="!canConnect"
          class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-bold flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm hover:shadow-md"
          @click="$emit('connect')"
        >
          <span class="material-icons text-lg">terminal</span>
          {{ t('connectTerminal') }}
        </button>
        <button
          v-else
          class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg text-sm font-bold flex items-center gap-2 transition-all shadow-sm hover:shadow-md"
          @click="$emit('disconnect')"
        >
          <span class="material-icons text-lg">link_off</span>
          {{ t('disconnect') }}
        </button>
        <span v-if="terminalConnected" class="flex items-center gap-2 text-sm text-emerald-500 font-medium bg-emerald-500/10 px-3 py-1 rounded-full border border-emerald-500/20">
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          {{ connectedLabel }}
        </span>
        <span v-else-if="!canConnect && disabledStatusText" class="text-xs text-amber-500 bg-amber-500/10 px-3 py-1 rounded-full border border-amber-500/20">
          {{ disabledStatusText }}
        </span>
      </div>

      <div v-if="showShortcuts" :class="shortcutsClass">
        <div class="flex items-center gap-1">
          <kbd class="px-1.5 py-0.5 rounded border border-border-light dark:border-border-dark bg-slate-100 dark:bg-zinc-800 font-mono text-[10px]">Ctrl+C</kbd>
          <span>{{ t('copy') }}</span>
        </div>
        <div class="flex items-center gap-1">
          <kbd class="px-1.5 py-0.5 rounded border border-border-light dark:border-border-dark bg-slate-100 dark:bg-zinc-800 font-mono text-[10px]">Ctrl+V</kbd>
          <span>{{ t('paste') }}</span>
        </div>
      </div>
    </div>

    <div class="flex-1 rounded-xl overflow-hidden flex flex-col border border-border-light dark:border-zinc-700/50 bg-[#1e1e1e]">
      <div class="bg-[#2d2d2d] px-4 py-2 flex items-center justify-between select-none">
        <div class="flex gap-2 group">
          <button class="w-3 h-3 rounded-full bg-[#ff5f56] hover:bg-[#ff5f56]/80 flex items-center justify-center text-[8px] text-black/0 group-hover:text-black/50 transition-colors" title="Disconnect" @click="$emit('disconnect')">×</button>
          <span class="w-3 h-3 rounded-full bg-[#ffbd2e] hover:bg-[#ffbd2e]/80"></span>
          <button class="w-3 h-3 rounded-full bg-[#27c93f] hover:bg-[#27c93f]/80 flex items-center justify-center text-[8px] text-black/0 group-hover:text-black/50 transition-colors" title="Fit Window" @click="$emit('fit')">↔</button>
        </div>
        <div class="flex-1 text-center text-xs text-slate-400 font-mono flex items-center justify-center gap-2">
          <span class="material-icons text-[14px]">dns</span>
          {{ terminalTitle }}
        </div>
        <div class="w-14"></div>
      </div>

      <div :ref="setTerminalRef" class="flex-1 p-1 relative overflow-hidden bg-[#1e1e1e]" :style="bodyStyle">
        <div v-if="!terminalConnected" class="absolute inset-0 z-10 flex items-center justify-center bg-[#1e1e1e]/80 backdrop-blur-sm transition-all duration-500">
          <div class="text-center transform transition-transform hover:scale-105 duration-300">
            <div class="w-20 h-20 bg-zinc-800 rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-xl border border-zinc-700">
              <span class="material-icons text-4xl text-slate-400">terminal</span>
            </div>
            <h3 class="text-slate-200 font-bold text-lg mb-2">{{ t('terminalDisconnected') }}</h3>
            <p class="text-slate-500 text-sm max-w-xs mx-auto mb-6">{{ t('terminalDesc') }}</p>
            <button
              v-if="canConnect"
              class="px-6 py-2 bg-primary hover:bg-primary-hover text-white rounded-full text-sm font-bold shadow-lg shadow-primary/20 transition-all hover:shadow-primary/40 hover:-translate-y-0.5 active:translate-y-0"
              @click="$emit('connect')"
            >
              {{ t('connectTerminal') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  bodyStyle: {
    type: [String, Object],
    default: ''
  },
  canConnect: {
    type: Boolean,
    default: false
  },
  connectedLabel: {
    type: String,
    default: ''
  },
  disabledStatusText: {
    type: String,
    default: ''
  },
  setTerminalRef: {
    type: Function,
    required: true
  },
  shortcutsClass: {
    type: String,
    default: 'flex items-center gap-4'
  },
  showShortcuts: {
    type: Boolean,
    default: true
  },
  terminalConnected: {
    type: Boolean,
    default: false
  },
  terminalTitle: {
    type: String,
    default: ''
  }
})

defineEmits(['connect', 'disconnect', 'fit'])

const t = inject('t', (key) => key)
</script>

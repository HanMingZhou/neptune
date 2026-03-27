<template>
  <section class="space-y-4">
    <h2 class="text-xs font-black uppercase tracking-widest text-slate-400">{{ t('audit.currentSessions') }}</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="session in sessions"
        :key="session.id"
        :class="[
          'p-5 border rounded-xl bg-white dark:bg-surface-dark transition-all',
          session.current ? 'border-primary ring-1 ring-primary/20' : 'border-gray-200 dark:border-border-dark'
        ]"
      >
        <div class="flex justify-between items-start mb-3">
          <div class="flex items-center gap-3">
            <div class="size-10 rounded-lg bg-slate-100 dark:bg-zinc-800 flex items-center justify-center text-slate-500">
              <span class="material-icons">
                {{ session.device && session.device.includes('iPhone') ? 'smartphone' : 'laptop' }}
              </span>
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="text-sm font-bold text-slate-800 dark:text-white">{{ session.device || '-' }}</span>
                <span
                  v-if="session.current"
                  class="bg-emerald-500/10 text-emerald-500 px-1.5 py-0.5 rounded text-[9px] font-black uppercase"
                >
                  {{ t('audit.currentDevice') }}
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-0.5">{{ session.location || '-' }} · {{ session.ip }}</p>
            </div>
          </div>
          <button
            v-if="!session.current"
            class="text-[10px] font-black uppercase text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 px-2 py-1 rounded transition-colors"
            @click="$emit('kill', session.id)"
          >
            {{ t('audit.killSession') }}
          </button>
        </div>
        <div class="flex items-center justify-between pt-3 border-t border-gray-100 dark:border-border-dark text-[10px] text-slate-400">
          <span>{{ t('audit.loginTime') }}: {{ session.time }}</span>
          <span class="flex items-center gap-1">
            <span class="size-1.5 rounded-full bg-emerald-500"></span>
            {{ t('audit.sessionActive') }}
          </span>
        </div>
      </div>
    </div>
    <div v-if="sessions.length === 0" class="text-center text-sm text-slate-400 py-6">{{ t('noData') }}</div>
  </section>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  sessions: {
    type: Array,
    default: () => []
  }
})

defineEmits(['kill'])

const t = inject('t', (key) => key)
</script>

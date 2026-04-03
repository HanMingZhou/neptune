<template>
  <div class="console-detail-card rounded-xl overflow-hidden">
    <div class="px-6 py-4 border-b border-border-light dark:border-border-dark flex items-center justify-between">
      <h3 class="font-bold">{{ t('instanceList') }}</h3>
      <button
        :disabled="podsLoading"
        class="px-3 py-1.5 border border-border-light dark:border-border-dark rounded-lg text-sm hover:bg-slate-50 dark:hover:bg-zinc-800 flex items-center gap-1"
        @click="$emit('refresh')"
      >
        <span class="material-icons text-lg" :class="{ 'animate-spin': podsLoading }">refresh</span>
        {{ t('refresh') }}
      </button>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-zinc-800/50 border-b border-border-light dark:border-border-dark text-slate-500 text-xs font-bold uppercase tracking-wider">
            <th class="px-6 py-4">{{ t('name') }}</th>
            <th class="px-6 py-4">{{ t('status') }}</th>
            <th class="px-6 py-4">{{ t('nodeIp') }}</th>
            <th class="px-6 py-4">{{ t('instanceIp') }}</th>
            <th class="px-6 py-4">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr v-for="pod in pods" :key="pod.name" class="hover:bg-slate-50 dark:hover:bg-zinc-800/40 transition-colors">
            <td class="px-6 py-4 text-sm font-mono">{{ pod.name }}</td>
            <td class="px-6 py-4">
              <span :class="['px-2 py-1 rounded-full text-xs font-bold', getPodStatusClass(pod.status)]">{{ pod.status }}</span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ pod.hostIP || '-' }}</td>
            <td class="px-6 py-4 text-sm text-slate-500">{{ pod.podIP || '-' }}</td>
            <td class="px-6 py-4">
              <button class="text-primary hover:underline text-sm" @click="$emit('view-logs', pod)">{{ t('logs') }}</button>
            </td>
          </tr>
          <tr v-if="pods.length === 0">
            <td colspan="5" class="px-6 py-12 text-center text-slate-400">{{ t('noData') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  getPodStatusClass: {
    type: Function,
    required: true
  },
  pods: {
    type: Array,
    default: () => []
  },
  podsLoading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['refresh', 'view-logs'])

const t = inject('t', (key) => key)
</script>

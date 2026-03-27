<template>
  <div class="overflow-hidden rounded-xl border border-border-light bg-surface-light text-sm shadow-sm dark:border-border-dark dark:bg-surface-dark">
    <div class="flex items-center justify-between border-b border-border-light px-6 py-4 dark:border-border-dark">
      <h3 class="text-sm font-bold uppercase tracking-widest text-slate-400">{{ t('recentInstances') }}</h3>
      <button class="text-xs font-bold text-primary hover:underline" @click="$emit('view-all')">{{ t('viewAll') }}</button>
    </div>
    <table class="w-full text-left">
      <thead>
        <tr class="bg-slate-50 text-[10px] font-black uppercase tracking-widest text-slate-500 dark:bg-zinc-800/50">
          <th class="px-6 py-4">{{ t('name') }}</th>
          <th class="px-6 py-4">{{ t('type') }}</th>
          <th class="px-6 py-4">{{ t('status') }}</th>
          <th class="px-6 py-4">{{ t('gpu') }}</th>
          <th class="px-6 py-4">{{ t('createdAt') }}</th>
          <th class="px-6 py-4 text-right">{{ t('actions') }}</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-border-light dark:divide-border-dark">
        <tr v-if="items.length === 0">
          <td colspan="6" class="px-6 py-8 text-center text-slate-400">{{ t('noData') }}</td>
        </tr>
        <tr
          v-for="item in items"
          :key="`${item.type}-${item.id}`"
          class="cursor-pointer transition-colors hover:bg-slate-50 dark:hover:bg-zinc-800"
          @click="$emit('open-detail', item)"
        >
          <td class="px-6 py-4">
            <span class="text-[13px] font-bold text-primary">{{ item.name }}</span>
          </td>
          <td class="px-6 py-4">
            <span class="rounded px-2 py-0.5 text-[10px] font-bold uppercase" :class="getTypeClass(item.type)">
              {{ t(item.type) }}
            </span>
          </td>
          <td class="px-6 py-4">
            <span class="rounded-full px-2 py-0.5 text-[10px] font-black uppercase" :class="getStatusClass(item.status)">
              <span
                class="mr-1 inline-block size-1.5 rounded-full bg-current"
                :class="item.status === 'Running' || item.status === 'RUNNING' ? 'animate-pulse' : ''"
              ></span>
              {{ t(item.status) }}
            </span>
          </td>
          <td class="px-6 py-4">
            <span v-if="item.gpu > 0" class="rounded bg-primary/10 px-2 py-0.5 text-[11px] font-mono font-bold text-primary">
              {{ item.gpu }} GPU
            </span>
            <span v-else class="text-[11px] text-slate-400">CPU</span>
          </td>
          <td class="px-6 py-4 text-[12px] text-slate-400">{{ item.createdAt }}</td>
          <td class="px-6 py-4 text-right">
            <button class="material-icons text-[18px] text-slate-400 transition-colors hover:text-primary">open_in_new</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  getStatusClass: {
    type: Function,
    required: true
  },
  getTypeClass: {
    type: Function,
    required: true
  },
  items: {
    type: Array,
    default: () => []
  }
})

defineEmits(['open-detail', 'view-all'])

const t = inject('t', (key) => key)
</script>

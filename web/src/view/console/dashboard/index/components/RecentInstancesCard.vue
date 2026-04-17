<template>
  <div
    class="overflow-hidden rounded-xl border border-border-light bg-surface-light text-sm shadow-sm dark:border-border-dark dark:bg-surface-dark"
  >
    <div
      class="flex items-center justify-between border-b border-border-light px-6 py-4 dark:border-border-dark"
    >
      <h3
        class="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-300"
      >
        {{ t('recentInstances') }}
      </h3>
      <button
        class="text-xs font-bold text-primary hover:underline"
        @click="$emit('view-all')"
      >
        {{ t('viewAll') }}
      </button>
    </div>
    <table class="w-full text-left">
      <thead>
        <tr
          class="bg-slate-50 text-xs font-black uppercase tracking-widest text-slate-500 dark:bg-zinc-800/50"
        >
          <th class="px-6 py-3">{{ t('name') }}</th>
          <th class="px-6 py-3">{{ t('type') }}</th>
          <th class="px-6 py-3">{{ t('status') }}</th>
          <th class="px-6 py-3">{{ t('gpu') }}</th>
          <th class="px-6 py-3">{{ t('createdAt') }}</th>
          <th class="px-6 py-3 text-right">{{ t('actions') }}</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-border-light dark:divide-border-dark">
        <tr v-if="items.length === 0">
          <td colspan="6" class="px-6 py-8 text-center text-sm text-slate-400">
            {{ t('noData') }}
          </td>
        </tr>
        <tr
          v-for="item in items"
          :key="`${item.type}-${item.id}`"
          class="recent-instance-row cursor-pointer transition-colors"
          :class="getRowClass(item.type)"
          @click="$emit('open-detail', item)"
        >
          <td class="px-6 py-3">
            <span class="text-[13px] font-bold text-primary">{{
              item.name
            }}</span>
          </td>
          <td class="px-6 py-3">
            <span
              class="recent-chip recent-chip--type px-2 py-0.5 text-[10px] font-bold uppercase"
              :class="getTypeClass(item.type)"
            >
              {{ t(item.type || '') }}
            </span>
          </td>
          <td class="px-6 py-3">
            <span
              class="recent-chip recent-chip--status px-2 py-0.5 text-[10px] font-black uppercase"
              :class="getStatusClass(item.status)"
            >
              <span
                class="mr-1 inline-block size-1.5 rounded-full bg-current"
                :class="
                  item.status === 'Running' || item.status === 'RUNNING'
                    ? 'animate-pulse'
                    : ''
                "
              ></span>
              {{ t(item.status || '') }}
            </span>
          </td>
          <td class="px-6 py-3">
            <span
              v-if="(item.gpu ?? 0) > 0"
              class="recent-chip recent-chip--gpu px-2 py-0.5 text-[11px] font-mono font-bold"
              :class="getGpuClass(item.gpu)"
            >
              {{ item.gpu ?? 0 }} GPU
            </span>
            <span
              v-else
              class="recent-chip recent-chip--gpu recent-chip--cpu px-2 py-0.5 text-[11px] font-mono font-bold"
            >
              CPU
            </span>
          </td>
          <td class="px-6 py-3 text-[12px] text-slate-400">
            {{ item.createdAt }}
          </td>
          <td class="px-6 py-3 text-right">
            <button
              class="material-icons text-[18px] text-slate-400 transition-colors hover:text-primary"
            >
              open_in_new
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import type { DashboardRecentInstance } from '@/types/dashboard'

withDefaults(
  defineProps<{
    getStatusClass: (status?: string) => string
    getTypeClass: (type?: string) => string
    items?: DashboardRecentInstance[]
  }>(),
  {
    items: () => []
  }
)

defineEmits<{
  'open-detail': [item: DashboardRecentInstance]
  'view-all': []
}>()

const t = inject<Translator>('t', (key: string) => key)

const getRowClass = (type?: string): string => {
  const map: Record<string, string> = {
    notebook: 'recent-instance-row--notebook',
    training: 'recent-instance-row--training',
    inference: 'recent-instance-row--inference'
  }

  return map[type ?? ''] || ''
}

const getGpuClass = (gpu?: number): string =>
  (gpu ?? 0) > 0 ? 'recent-chip--gpu-active' : 'recent-chip--cpu'
</script>

<style scoped>
.recent-chip {
  display: inline-flex;
  align-items: center;
  border-width: 1px;
  border-style: solid;
  line-height: 1.2;
  white-space: nowrap;
}

.recent-chip--type {
  border-radius: 9999px;
}

.recent-chip--status {
  border-radius: 9999px;
}

.recent-chip--gpu {
  border-radius: 0.5rem;
}

.recent-chip--type-notebook {
  border-color: rgb(96 165 250);
  background: rgb(239 246 255);
  color: rgb(29 78 216);
}

.recent-chip--type-training {
  border-color: rgb(251 146 60);
  background: rgb(255 247 237);
  color: rgb(194 65 12);
}

.recent-chip--type-inference {
  border-color: rgb(167 139 250);
  background: rgb(245 243 255);
  color: rgb(109 40 217);
}

.recent-chip--type-default {
  border-color: rgb(203 213 225);
  background: rgb(248 250 252);
  color: rgb(71 85 105);
}

.recent-chip--status-running {
  border-color: rgb(110 231 183);
  background: rgb(236 253 245);
  color: rgb(4 120 87);
}

.recent-chip--status-stopped {
  border-color: rgb(203 213 225);
  background: rgb(248 250 252);
  color: rgb(71 85 105);
}

.recent-chip--status-pending {
  border-color: rgb(252 211 77);
  background: rgb(255 251 235);
  color: rgb(180 83 9);
}

.recent-chip--status-failed {
  border-color: rgb(252 165 165);
  background: rgb(254 242 242);
  color: rgb(185 28 28);
}

.recent-chip--status-creating {
  border-color: rgb(147 197 253);
  background: rgb(239 246 255);
  color: rgb(29 78 216);
}

.recent-chip--status-default {
  border-color: rgb(203 213 225);
  background: rgb(248 250 252);
  color: rgb(71 85 105);
}

.recent-chip--gpu-active {
  border-color: rgb(125 211 252);
  background: rgb(240 249 255);
  color: rgb(3 105 161);
}

.recent-chip--cpu {
  border-color: rgb(203 213 225);
  background: rgb(241 245 249);
  color: rgb(100 116 139);
}

.recent-instance-row {
  box-shadow: inset 0 0 0 0 transparent;
}

.recent-instance-row:hover {
  background: rgb(248 250 252);
}

.recent-instance-row--notebook {
  box-shadow: inset 3px 0 0 rgb(59 130 246);
}

.recent-instance-row--training {
  box-shadow: inset 3px 0 0 rgb(249 115 22);
}

.recent-instance-row--inference {
  box-shadow: inset 3px 0 0 rgb(139 92 246);
}

.recent-instance-row--notebook td:first-child {
  background: linear-gradient(90deg, rgb(59 130 246 / 0.08), transparent 42%);
}

.recent-instance-row--training td:first-child {
  background: linear-gradient(90deg, rgb(249 115 22 / 0.08), transparent 42%);
}

.recent-instance-row--inference td:first-child {
  background: linear-gradient(90deg, rgb(139 92 246 / 0.08), transparent 42%);
}

.dark .recent-instance-row:hover {
  background: rgb(39 39 42);
}

.dark .recent-instance-row--notebook td:first-child {
  background: linear-gradient(90deg, rgb(59 130 246 / 0.12), transparent 46%);
}

.dark .recent-instance-row--training td:first-child {
  background: linear-gradient(90deg, rgb(249 115 22 / 0.12), transparent 46%);
}

.dark .recent-instance-row--inference td:first-child {
  background: linear-gradient(90deg, rgb(139 92 246 / 0.14), transparent 46%);
}

.dark .recent-chip--type-notebook {
  border-color: rgb(59 130 246 / 0.7);
  background: rgb(30 41 59);
  color: rgb(191 219 254);
}

.dark .recent-chip--type-training {
  border-color: rgb(249 115 22 / 0.78);
  background: rgb(67 20 7);
  color: rgb(254 215 170);
}

.dark .recent-chip--type-inference {
  border-color: rgb(139 92 246 / 0.78);
  background: rgb(46 16 101);
  color: rgb(221 214 254);
}

.dark .recent-chip--type-default,
.dark .recent-chip--status-stopped,
.dark .recent-chip--status-default,
.dark .recent-chip--cpu {
  border-color: rgb(82 82 91);
  background: rgb(39 39 42);
  color: rgb(203 213 225);
}

.dark .recent-chip--status-running {
  border-color: rgb(16 185 129 / 0.78);
  background: rgb(6 44 31);
  color: rgb(167 243 208);
}

.dark .recent-chip--status-pending {
  border-color: rgb(245 158 11 / 0.76);
  background: rgb(69 26 3);
  color: rgb(253 230 138);
}

.dark .recent-chip--status-failed {
  border-color: rgb(239 68 68 / 0.76);
  background: rgb(69 10 10);
  color: rgb(252 165 165);
}

.dark .recent-chip--status-creating {
  border-color: rgb(59 130 246 / 0.75);
  background: rgb(30 41 59);
  color: rgb(191 219 254);
}

.dark .recent-chip--gpu-active {
  border-color: rgb(56 189 248 / 0.76);
  background: rgb(8 47 73);
  color: rgb(186 230 253);
}
</style>

<template>
  <aside class="w-64 h-full border-r border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark flex flex-col shrink-0 overflow-hidden">
    <div class="p-6 pb-2 mb-2 sticky top-0 bg-surface-light dark:bg-surface-dark z-20">
      <div class="flex items-center gap-3 cursor-pointer group transition-all" @click="router.push({ path: '/' })">
        <div class="relative">
          <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-xl blur opacity-25 group-hover:opacity-50 transition duration-700"></div>
          <div class="relative w-11 h-11 rounded-xl overflow-hidden border border-white/10 shadow-2xl bg-slate-950">
            <img :src="logoUrl" class="w-full h-full object-cover scale-[1.7] group-hover:scale-[1.9] transition duration-1000" alt="Logo" />
          </div>
        </div>
        <div class="flex flex-col">
          <div class="flex items-center">
            <span class="text-xl font-[950] tracking-tighter text-slate-900 dark:text-white uppercase transition-colors group-hover:text-primary">
              机器<span class="text-blue-600 dark:text-blue-400">学习</span>
            </span>
          </div>
          <div class="flex items-center gap-1.5 -mt-1">
            <span class="text-[9px] font-black tracking-[4px] text-slate-400 dark:text-slate-500 uppercase opacity-70">Platform</span>
            <div class="h-[1px] w-8 bg-slate-200 dark:bg-slate-800 transition-all group-hover:w-12"></div>
          </div>
        </div>
      </div>
    </div>

    <nav class="flex-1 overflow-y-auto custom-scrollbar px-3 py-4 space-y-6">
      <div v-for="group in dynamicNavigation" :key="group.title" class="space-y-1">
        <p class="text-[10px] font-black text-slate-400 dark:text-slate-500 uppercase tracking-[2px] px-3 mb-2">
          {{ t(group.title) }}
        </p>
        <div v-for="item in group.items" :key="item.key">
          <NavMenuItem :item="item" :level="0" />
        </div>
      </div>
    </nav>

    <div class="shrink-0 border-t border-border-light bg-surface-light p-4 dark:border-border-dark dark:bg-surface-dark">
      <button
        class="w-full flex items-center justify-between gap-3 border border-border-light px-3 py-2.5 text-xs font-bold uppercase tracking-wider text-slate-500 transition-colors hover:bg-slate-50 dark:border-border-dark dark:bg-surface-dark dark:text-slate-400 dark:hover:bg-zinc-800"
        @click="emit('toggle-theme')"
      >
        <span class="flex items-center gap-3">
          <span class="material-icons text-[18px]">{{ isDark ? 'dark_mode' : 'light_mode' }}</span>
          {{ t('theme') }}
        </span>
        <span class="text-[10px] tracking-[0.16em] text-slate-400 dark:text-slate-500">
          {{ isDark ? 'DARK' : 'LIGHT' }}
        </span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { inject } from 'vue'
import { useRouter } from 'vue-router'
import logoUrl from '@/assets/logo.png'
import NavMenuItem from '@/components/NavMenuItem.vue'

defineProps({
  dynamicNavigation: {
    type: Array,
    default: () => []
  },
  isDark: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['toggle-theme'])
const router = useRouter()
const t = inject('t', (key) => key)
</script>

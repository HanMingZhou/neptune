<template>
  <header
    class="h-16 border-b border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark px-8 flex items-center justify-between shrink-0 z-10"
  >
    <div class="flex-1 max-w-lg">
      <div class="relative group">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
          >search</span
        >
        <input
          type="text"
          :placeholder="t('menuSearchPlaceholder')"
          class="w-full pl-10 pr-4 py-2 bg-slate-100 dark:bg-zinc-800 border-none rounded-lg text-sm focus:ring-1 focus:ring-primary transition-all outline-none"
        />
      </div>
    </div>
    <div class="flex items-center gap-4">
      <button
        class="px-3 py-1.5 border border-border-light dark:border-border-dark rounded-md text-[10px] font-black hover:bg-slate-100 dark:hover:bg-zinc-800 flex items-center gap-2 transition-all shadow-sm"
        @click="emit('toggle-lang')"
      >
        <span class="material-icons text-[16px]">translate</span>
        {{ lang === 'en' ? 'EN' : 'CN' }}
      </button>
      <button
        class="material-icons text-slate-400 hover:text-primary transition-colors"
      >
        notifications
      </button>

      <LayoutUserPopover
        :user-balance="userBalance"
        :user-info="userInfo"
        @change-user-auth="emit('change-user-auth', $event)"
        @logout="emit('logout')"
      />
    </div>
  </header>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import LayoutUserPopover from './LayoutUserPopover.vue'

interface LayoutHeaderUserInfo {
  nickName?: string
  userName?: string
  authorityId?: ResourceId
  authority?: {
    authorityId?: ResourceId
    authorityName?: string
  }
  authorities?: Array<{
    authorityId: ResourceId
    authorityName?: string
  }>
}

withDefaults(
  defineProps<{
    lang?: string
    userBalance?: number
    userInfo?: LayoutHeaderUserInfo
  }>(),
  {
    lang: 'zh',
    userBalance: 0,
    userInfo: () => ({})
  }
)

const emit = defineEmits<{
  'change-user-auth': [authorityId: ResourceId]
  logout: []
  'toggle-lang': []
}>()
const t = inject<Translator>('t', (key: string) => key)
</script>

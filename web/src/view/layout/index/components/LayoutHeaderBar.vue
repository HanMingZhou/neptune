<template>
  <header
    class="h-16 border-b border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark px-8 flex items-center justify-between shrink-0 z-10"
  >
    <div class="relative flex-1 max-w-xl">
      <CommandMenu
        :groups="dynamicNavigation"
        :placeholder="t('menuSearchPlaceholder')"
      />
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
import type { LayoutNavGroup } from '@/types/layout'
import CommandMenu from '@/components/commandMenu/index.vue'
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
    dynamicNavigation?: LayoutNavGroup[]
    lang?: string
    userBalance?: number
    userInfo?: LayoutHeaderUserInfo
  }>(),
  {
    dynamicNavigation: () => [],
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

<template>
  <div class="relative">
    <div
      class="flex items-center gap-3 pl-4 border-l border-border-light dark:border-border-dark cursor-pointer"
      @click="showUserPopover = !showUserPopover"
    >
      <div class="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-xs shadow-sm border border-primary/10">
        {{ avatarText }}
      </div>
      <div class="hidden lg:block text-left">
        <p class="text-sm font-bold leading-none">{{ displayName }}</p>
        <p class="text-[9px] text-slate-500 font-bold mt-1 uppercase tracking-tighter">{{ currentAuthorityName }}</p>
      </div>
      <span class="material-icons text-slate-400 text-[16px] transition-transform" :class="{ 'rotate-180': showUserPopover }">expand_more</span>
    </div>

    <div v-if="showUserPopover" class="fixed inset-0 z-40" @click="showUserPopover = false"></div>

    <transition name="popover-fade">
      <div v-if="showUserPopover" class="absolute right-0 top-full mt-2 w-80 bg-white dark:bg-zinc-900 border border-border-light dark:border-border-dark rounded-xl shadow-2xl z-50 overflow-hidden">
        <div class="p-4 bg-gradient-to-r from-primary/5 to-blue-500/5 border-b border-border-light dark:border-border-dark">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-sm border border-primary/10">
              {{ avatarText }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-bold truncate">{{ displayName }}</p>
              <p class="text-[11px] text-slate-400 font-mono truncate">{{ username }}</p>
            </div>
          </div>
        </div>

        <div class="mx-3 mt-3 p-3 bg-slate-50 dark:bg-zinc-800 rounded-lg border border-border-light dark:border-border-dark">
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-3">
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{{ t('dashboard.balance') }}</p>
              <span class="shrink-0 material-icons text-primary/30 text-2xl">account_balance_wallet</span>
            </div>
            <p class="break-all text-lg leading-tight font-black tabular-nums text-primary">
              {{ fullBalanceText }}
            </p>
          </div>
        </div>

        <div class="p-2 mt-1">
          <template v-if="alternativeAuthorities.length > 0">
            <p class="px-3 py-1.5 text-[10px] font-bold text-slate-400 uppercase tracking-wider">{{ t('switchRole') }}</p>
            <button
              v-for="auth in alternativeAuthorities"
              :key="auth.authorityId"
              class="w-full flex items-center gap-3 px-3 py-2 text-sm rounded-lg hover:bg-slate-50 dark:hover:bg-zinc-800 transition-colors"
              @click="handleChangeAuthority(auth.authorityId)"
            >
              <span class="material-icons text-[18px] text-slate-400">swap_horiz</span>
              <span>{{ auth.authorityName }}</span>
            </button>
            <div class="mx-3 my-1 border-t border-border-light dark:border-border-dark"></div>
          </template>

          <button
            class="w-full flex items-center gap-3 px-3 py-2 text-sm rounded-lg hover:bg-slate-50 dark:hover:bg-zinc-800 transition-colors"
            @click="goToProfile"
          >
            <span class="material-icons text-[18px] text-slate-400">person</span>
            <span>{{ t('person') }}</span>
          </button>

          <div class="mx-3 my-1 border-t border-border-light dark:border-border-dark"></div>

          <button
            class="w-full flex items-center gap-3 px-3 py-2 text-sm rounded-lg hover:bg-red-50 dark:hover:bg-red-500/10 text-red-500 transition-colors"
            @click="handleLogout"
          >
            <span class="material-icons text-[18px]">logout</span>
            <span>{{ t('logout') }}</span>
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  userBalance: {
    type: Number,
    default: 0
  },
  userInfo: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['change-user-auth', 'logout'])
const router = useRouter()
const t = inject('t', (key) => key)
const showUserPopover = ref(false)

const avatarText = computed(() => props.userInfo?.nickName?.charAt(0)?.toUpperCase() || 'U')
const displayName = computed(() => props.userInfo?.nickName || 'Admin_User')
const username = computed(() => props.userInfo?.userName || '-')
const currentAuthorityName = computed(() => props.userInfo?.authority?.authorityName || 'System Root')
const currentAuthorityId = computed(() => props.userInfo?.authorityId ?? props.userInfo?.authority?.authorityId)
const alternativeAuthorities = computed(() =>
  (props.userInfo?.authorities || []).filter((item) => item.authorityId !== currentAuthorityId.value)
)
const fullBalanceText = computed(() =>
  `¥${Number(props.userBalance || 0).toLocaleString(undefined, {
    maximumFractionDigits: 6,
    minimumFractionDigits: 2
  })}`
)

const closePopover = () => {
  showUserPopover.value = false
}

const handleChangeAuthority = (authorityId) => {
  emit('change-user-auth', authorityId)
  closePopover()
}

const goToProfile = () => {
  router.push({ name: 'person' })
  closePopover()
}

const handleLogout = () => {
  emit('logout')
  closePopover()
}
</script>

<style scoped>
.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>

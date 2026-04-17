<template>
  <div ref="rootRef" class="command-menu-shell relative w-full">
    <div
      class="command-menu-trigger group relative overflow-hidden rounded-[22px] border border-slate-200/70 bg-white/80 shadow-[0_18px_44px_-30px_rgba(15,23,42,0.55)] transition-all duration-300 focus-within:border-sky-400/50 focus-within:shadow-[0_24px_70px_-32px_rgba(14,165,233,0.45)] dark:border-zinc-800/90 dark:bg-zinc-950/85 dark:focus-within:border-sky-400/45"
    >
      <div class="pointer-events-none absolute inset-0 overflow-hidden">
        <div
          class="absolute inset-y-1 left-1.5 w-24 rounded-[18px] bg-[linear-gradient(135deg,rgba(14,165,233,0.18),rgba(59,130,246,0.04)_58%,transparent)] blur-[1px]"
        ></div>
        <div
          class="absolute -left-8 top-[-18px] h-20 w-28 rounded-full bg-sky-400/15 blur-2xl transition-opacity duration-300 group-focus-within:opacity-100 dark:bg-sky-500/10"
        ></div>
        <div
          class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(59,130,246,0.14),transparent_42%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.08),transparent_38%)] opacity-90"
        ></div>
      </div>

      <span
        class="pointer-events-none absolute left-3 top-1/2 z-10 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-2xl border border-white/70 bg-white/90 text-slate-500 shadow-sm transition-all duration-300 group-focus-within:border-sky-200 group-focus-within:text-sky-600 dark:border-white/8 dark:bg-zinc-900/90 dark:text-slate-300 dark:group-focus-within:border-sky-500/30 dark:group-focus-within:text-sky-300"
      >
        <span class="material-icons text-[18px] leading-none">search</span>
      </span>

      <input
        ref="inputRef"
        v-model="query"
        type="text"
        :placeholder="placeholder"
        class="relative z-10 w-full rounded-[22px] bg-transparent py-3 pl-14 pr-18 text-sm font-medium text-slate-700 outline-none placeholder:text-slate-400 dark:text-slate-100 dark:placeholder:text-slate-500"
        @focus="openMenu"
        @keydown.down.prevent="moveSelection(1)"
        @keydown.up.prevent="moveSelection(-1)"
        @keydown.enter.prevent="selectActiveResult"
        @keydown.esc.prevent="closeMenu"
      />

      <button
        type="button"
        class="absolute right-2.5 top-1/2 z-10 inline-flex -translate-y-1/2 items-center gap-1 rounded-full border border-slate-200/80 bg-white/90 px-2.5 py-1 text-[10px] font-black tracking-[0.24em] text-slate-400 transition-all duration-300 hover:border-sky-200 hover:text-sky-600 dark:border-zinc-700/80 dark:bg-zinc-900/90 dark:text-slate-500 dark:hover:border-sky-500/30 dark:hover:text-sky-300"
        @click="focusInput"
      >
        {{ shortcutLabel }}
      </button>
    </div>

    <transition name="command-menu-panel">
      <div
        v-if="isOpen"
        class="command-menu-surface absolute left-0 right-0 top-[calc(100%+0.95rem)] z-40 overflow-hidden rounded-[28px] border border-slate-200/75 bg-white/92 shadow-[0_32px_95px_-36px_rgba(15,23,42,0.55)] backdrop-blur-xl dark:border-zinc-800/90 dark:bg-zinc-950/92"
      >
        <div class="pointer-events-none absolute inset-0 overflow-hidden">
          <div
            class="absolute inset-x-0 top-0 h-24 bg-[linear-gradient(180deg,rgba(59,130,246,0.16),rgba(255,255,255,0))] dark:bg-[linear-gradient(180deg,rgba(14,165,233,0.16),rgba(0,0,0,0))]"
          ></div>
          <div
            class="absolute left-6 top-6 h-32 w-32 rounded-full bg-sky-400/10 blur-3xl dark:bg-sky-500/10"
          ></div>
          <div
            class="absolute -right-12 top-12 h-40 w-40 rounded-full bg-blue-500/10 blur-3xl dark:bg-blue-500/10"
          ></div>
          <div class="command-menu-grid absolute inset-x-0 bottom-0 top-0 opacity-40"></div>
        </div>

        <div
          class="relative flex items-start justify-between gap-4 border-b border-slate-200/80 px-5 py-4 dark:border-zinc-800/90"
        >
          <div class="min-w-0">
            <div
              class="inline-flex items-center gap-2 rounded-full border border-sky-100 bg-white/80 px-3 py-1 text-[10px] font-black uppercase tracking-[0.28em] text-sky-700 shadow-sm dark:border-sky-500/20 dark:bg-zinc-900/85 dark:text-sky-300"
            >
              <span class="h-2 w-2 rounded-full bg-sky-500 shadow-[0_0_0_4px_rgba(14,165,233,0.12)]"></span>
              {{ t('menus') }}
            </div>
            <p
              class="mt-3 max-w-[28rem] text-sm font-semibold tracking-tight text-slate-700 dark:text-slate-100"
            >
              {{ query ? t('searchMenuDesc') : t('menuSearchPlaceholder') }}
            </p>
            <p class="mt-1 text-xs text-slate-400 dark:text-slate-500">
              {{ currentRouteName || '/' }}
            </p>
          </div>
          <div
            class="shrink-0 rounded-[20px] border border-slate-200/80 bg-white/80 px-3 py-2 text-right shadow-sm dark:border-zinc-800/90 dark:bg-zinc-900/80"
          >
            <div
              class="text-[10px] font-black uppercase tracking-[0.24em] text-slate-400 dark:text-slate-500"
            >
              {{ query ? t('searchQuery') : t('menus') }}
            </div>
            <div class="mt-1 text-lg font-black tracking-tight text-slate-800 dark:text-white">
              {{ resultCountLabel }}
            </div>
          </div>
        </div>

        <div
          v-if="filteredEntries.length"
          class="relative max-h-[460px] overflow-y-auto px-3 pb-3 pt-2"
        >
          <button
            v-for="(entry, index) in filteredEntries"
            :key="entry.key"
            type="button"
            :ref="(element) => setItemRef(element, index)"
            class="group relative mb-2 flex w-full items-center gap-4 overflow-hidden rounded-[22px] border px-4 py-3.5 text-left transition-all duration-300 last:mb-0"
            :class="
              index === activeIndex
                ? 'border-sky-200/80 bg-[linear-gradient(135deg,rgba(14,165,233,0.12),rgba(255,255,255,0.96)_48%)] text-slate-900 shadow-[0_24px_44px_-34px_rgba(14,165,233,0.65)] dark:border-sky-500/25 dark:bg-[linear-gradient(135deg,rgba(14,165,233,0.18),rgba(9,9,11,0.95)_52%)] dark:text-white'
                : 'border-slate-200/75 bg-white/72 text-slate-600 hover:border-slate-300 hover:bg-white/90 hover:shadow-[0_18px_30px_-28px_rgba(15,23,42,0.45)] dark:border-zinc-800/85 dark:bg-zinc-900/68 dark:text-slate-300 dark:hover:border-zinc-700 dark:hover:bg-zinc-900/92'
            "
            @mouseenter="activeIndex = index"
            @mousedown.prevent="navigateToEntry(entry)"
          >
            <div
              class="pointer-events-none absolute inset-0 bg-[linear-gradient(120deg,rgba(59,130,246,0.07),transparent_40%,transparent_60%,rgba(14,165,233,0.05))] opacity-0 transition-opacity duration-300 group-hover:opacity-100"
              :class="index === activeIndex ? 'opacity-100' : ''"
            ></div>

            <div
              class="relative z-10 flex h-12 w-12 shrink-0 items-center justify-center rounded-[18px] border border-slate-200/80 bg-slate-50/90 text-slate-500 shadow-sm transition-all duration-300 dark:border-zinc-800 dark:bg-zinc-900 dark:text-slate-300"
              :class="
                index === activeIndex
                  ? 'border-sky-200 bg-sky-50 text-sky-600 shadow-[0_12px_24px_-18px_rgba(14,165,233,0.7)] dark:border-sky-500/30 dark:bg-sky-500/10 dark:text-sky-300'
                  : ''
              "
            >
              <AppIcon :name="entry.icon" class="text-[18px]" />
            </div>

            <div class="relative z-10 min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span
                  class="inline-flex items-center rounded-full border border-slate-200/80 bg-slate-100/90 px-2.5 py-1 text-[10px] font-black uppercase tracking-[0.22em] text-slate-500 dark:border-zinc-800/90 dark:bg-zinc-900 dark:text-slate-400"
                  :class="
                    index === activeIndex
                      ? 'border-sky-200 bg-sky-50 text-sky-700 dark:border-sky-500/25 dark:bg-sky-500/10 dark:text-sky-300'
                      : ''
                  "
                >
                  {{ entry.groupTitle }}
                </span>
                <span
                  v-if="entry.routeName === currentRouteName"
                  class="inline-flex h-2.5 w-2.5 rounded-full bg-emerald-500 shadow-[0_0_0_5px_rgba(16,185,129,0.12)]"
                ></span>
              </div>

              <div class="mt-2 flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <p class="truncate text-[15px] font-black tracking-tight">
                    {{ entry.title }}
                  </p>
                  <div
                    class="mt-2 inline-flex max-w-full items-center gap-2 rounded-full border border-slate-200/75 bg-slate-50/80 px-2.5 py-1 text-[11px] font-medium text-slate-500 dark:border-zinc-800 dark:bg-zinc-900/90 dark:text-slate-400"
                    :class="
                      index === activeIndex
                        ? 'border-sky-100 bg-white/80 text-slate-600 dark:border-sky-500/20 dark:bg-zinc-900 dark:text-slate-300'
                        : ''
                    "
                  >
                    <span class="material-icons text-[14px]">alt_route</span>
                    <span class="truncate font-mono">{{ entry.fullPath }}</span>
                  </div>
                </div>

                <span
                  class="hidden shrink-0 rounded-full border border-slate-200/80 bg-white/75 px-2.5 py-1 font-mono text-[10px] font-bold uppercase tracking-[0.18em] text-slate-400 shadow-sm sm:inline-flex dark:border-zinc-800 dark:bg-zinc-900 dark:text-slate-500"
                  :class="
                    index === activeIndex
                      ? 'border-sky-100 text-sky-600 dark:border-sky-500/20 dark:text-sky-300'
                      : ''
                  "
                >
                  {{ entry.routeName }}
                </span>
              </div>

              <div
                class="mt-3 flex flex-wrap items-center gap-1.5 text-[11px] text-slate-400 dark:text-slate-500"
              >
                <span
                  v-for="parent in entry.parents.slice(0, 3)"
                  :key="`${entry.key}-${parent}`"
                  class="inline-flex items-center rounded-full border border-slate-200/70 bg-white/70 px-2 py-1 font-medium dark:border-zinc-800 dark:bg-zinc-900"
                >
                  {{ parent }}
                </span>
                <span
                  v-if="entry.parents.length > 3"
                  class="inline-flex items-center rounded-full border border-slate-200/70 bg-white/70 px-2 py-1 font-medium dark:border-zinc-800 dark:bg-zinc-900"
                >
                  +{{ entry.parents.length - 3 }}
                </span>
                <span
                  v-if="!entry.parents.length"
                  class="inline-flex items-center rounded-full border border-dashed border-slate-200/80 px-2 py-1 font-medium dark:border-zinc-800"
                >
                  {{ formatMeta(entry) }}
                </span>
              </div>
            </div>

            <div
              class="relative z-10 flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-transparent bg-transparent transition-all duration-300"
              :class="
                index === activeIndex
                  ? 'border-sky-100 bg-white/90 text-sky-600 shadow-sm dark:border-sky-500/20 dark:bg-zinc-900 dark:text-sky-300'
                  : 'text-slate-300 group-hover:border-slate-200 group-hover:bg-white/80 group-hover:text-slate-500 dark:text-slate-600 dark:group-hover:border-zinc-800 dark:group-hover:bg-zinc-900 dark:group-hover:text-slate-300'
              "
            >
              <span class="material-icons text-[18px]">north_east</span>
            </div>
          </button>
        </div>

        <div
          v-else
          class="relative flex flex-col items-center justify-center gap-3 px-6 py-12 text-center"
        >
          <div
            class="flex h-16 w-16 items-center justify-center rounded-full border border-slate-200/80 bg-white/80 shadow-[0_24px_45px_-30px_rgba(15,23,42,0.55)] dark:border-zinc-800 dark:bg-zinc-900"
          >
            <span
              class="material-icons text-[30px] text-slate-300 dark:text-slate-600"
              >search_off</span
            >
          </div>
          <span
            class="text-[10px] font-black uppercase tracking-[0.28em] text-slate-400 dark:text-slate-500"
          >
            {{ t('menus') }}
          </span>
          <p class="text-base font-black tracking-tight text-slate-700 dark:text-slate-200">
            {{ t('noData') }}
          </p>
          <p class="max-w-[18rem] text-xs leading-5 text-slate-400 dark:text-slate-500">
            {{ t('searchMenuDesc') }}
          </p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useRoute, useRouter } from 'vue-router'
import AppIcon from '@/components/AppIcon.vue'
import type { Translator } from '@/types/consoleResource'
import type { LayoutNavGroup, LayoutNavItem } from '@/types/layout'

interface CommandEntry {
  key: string
  title: string
  routeName: string
  icon: string
  groupTitle: string
  parents: string[]
  fullPath: string
  searchText: string
}

const props = withDefaults(
  defineProps<{
    groups?: LayoutNavGroup[]
    placeholder?: string
  }>(),
  {
    groups: () => [],
    placeholder: ''
  }
)

const t = inject<Translator>('t', (key: string) => key)
const route = useRoute()
const router = useRouter()

const rootRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const itemRefs = ref<HTMLElement[]>([])
const query = ref('')
const isOpen = ref(false)
const activeIndex = ref(0)

const currentRouteName = computed(() =>
  typeof route.name === 'string' ? route.name : ''
)

const shortcutLabel = computed(() => {
  if (typeof navigator === 'undefined') {
    return 'CTRL K'
  }

  return /mac/i.test(navigator.platform) ? 'CMD K' : 'CTRL K'
})

const resolveTitle = (item: Pick<LayoutNavItem, 'title' | 'titleKey'>): string => {
  if (!item.titleKey) {
    return item.title
  }

  const translated = t(item.titleKey)
  return translated !== item.titleKey ? translated : item.title
}

const flattenEntries = (
  items: LayoutNavItem[],
  groupTitle: string,
  parents: string[] = []
): CommandEntry[] =>
  items.flatMap((item) => {
    const title = resolveTitle(item)
    const nextParents = item.routeName ? parents : [...parents, title]
    const childrenEntries = item.children?.length
      ? flattenEntries(item.children, groupTitle, nextParents)
      : []

    if (!item.routeName) {
      return childrenEntries
    }

    const fullPath = router.resolve({ name: item.routeName }).fullPath
    const searchText = [
      title,
      item.title,
      item.titleKey,
      item.routeName,
      fullPath,
      groupTitle,
      ...parents
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase()

    return [
      {
        key: item.key,
        title,
        routeName: item.routeName,
        icon: item.icon,
        groupTitle,
        parents,
        fullPath,
        searchText
      },
      ...childrenEntries
    ]
  })

const allEntries = computed<CommandEntry[]>(() =>
  props.groups.flatMap((group) => {
    const groupTitle = t(group.title)
    return flattenEntries(group.items || [], groupTitle)
  })
)

const normalizedQuery = computed(() => query.value.trim().toLowerCase())

const getEntryScore = (entry: CommandEntry, keyword: string): number => {
  if (!keyword) {
    return 0
  }

  const title = entry.title.toLowerCase()
  const routeName = entry.routeName.toLowerCase()
  const fullPath = entry.fullPath.toLowerCase()
  const breadcrumbs = entry.parents.join(' ').toLowerCase()
  let score = 0

  if (title === keyword) {
    score += 400
  }
  if (title.startsWith(keyword)) {
    score += 280
  } else if (title.includes(keyword)) {
    score += 180
  }

  if (routeName === keyword) {
    score += 220
  }
  if (routeName.startsWith(keyword)) {
    score += 140
  } else if (routeName.includes(keyword)) {
    score += 80
  }

  if (fullPath.includes(keyword)) {
    score += 60
  }
  if (breadcrumbs.includes(keyword)) {
    score += 30
  }

  return score - Math.min(title.length, 32)
}

const filteredEntries = computed<CommandEntry[]>(() => {
  const keyword = normalizedQuery.value

  if (!keyword) {
    return allEntries.value
  }

  return [...allEntries.value]
    .filter((entry) => entry.searchText.includes(keyword))
    .sort((left, right) => getEntryScore(right, keyword) - getEntryScore(left, keyword))
})

const resultCountLabel = computed(() => `${filteredEntries.value.length}`)

const setActiveResult = (): void => {
  if (!filteredEntries.value.length) {
    activeIndex.value = 0
    return
  }

  if (!normalizedQuery.value) {
    const currentIndex = filteredEntries.value.findIndex(
      (entry) => entry.routeName === currentRouteName.value
    )
    activeIndex.value = currentIndex >= 0 ? currentIndex : 0
    return
  }

  activeIndex.value = Math.min(activeIndex.value, filteredEntries.value.length - 1)
}

const setItemRef = (element: Element | null, index: number): void => {
  if (!(element instanceof HTMLElement)) {
    return
  }

  itemRefs.value[index] = element
}

const scrollActiveResultIntoView = (): void => {
  const activeElement = itemRefs.value[activeIndex.value]
  activeElement?.scrollIntoView({
    block: 'nearest'
  })
}

const openMenu = (): void => {
  isOpen.value = true
  setActiveResult()
}

const closeMenu = (): void => {
  isOpen.value = false
}

const focusInput = async (): Promise<void> => {
  openMenu()
  await nextTick()
  inputRef.value?.focus()
  inputRef.value?.select()
}

const moveSelection = (step: number): void => {
  openMenu()

  if (!filteredEntries.value.length) {
    return
  }

  const nextIndex =
    (activeIndex.value + step + filteredEntries.value.length) %
    filteredEntries.value.length
  activeIndex.value = nextIndex
}

const navigateToEntry = async (entry: CommandEntry): Promise<void> => {
  query.value = ''
  closeMenu()

  if (entry.routeName === currentRouteName.value) {
    return
  }

  await router.push({ name: entry.routeName }).catch((error: unknown) => {
    console.error('Failed to navigate from command menu:', error)
  })
}

const selectActiveResult = (): void => {
  const entry = filteredEntries.value[activeIndex.value]
  if (!entry) {
    return
  }

  void navigateToEntry(entry)
}

const formatMeta = (entry: CommandEntry): string => {
  const parts = [entry.groupTitle, ...entry.parents].filter(Boolean)
  return parts.join(' / ')
}

const handleGlobalKeydown = (event: KeyboardEvent): void => {
  const target = event.target as HTMLElement | null
  const isTypingTarget =
    target?.tagName === 'INPUT' ||
    target?.tagName === 'TEXTAREA' ||
    target?.isContentEditable

  if (isTypingTarget) {
    return
  }

  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    void focusInput()
  }
}

watch(normalizedQuery, () => {
  activeIndex.value = 0
})

watch(filteredEntries, () => {
  itemRefs.value = []
  setActiveResult()
})

watch(activeIndex, () => {
  void nextTick(() => {
    scrollActiveResultIntoView()
  })
})

watch(
  () => route.fullPath,
  () => {
    query.value = ''
    closeMenu()
  }
)

onClickOutside(rootRef, () => {
  closeMenu()
})

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<style scoped>
.command-menu-shell {
  isolation: isolate;
}

.command-menu-trigger::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  opacity: 0.75;
  pointer-events: none;
}

.dark .command-menu-trigger::after {
  border-color: rgba(255, 255, 255, 0.05);
  opacity: 1;
}

.command-menu-surface::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 28px;
  border: 1px solid rgba(255, 255, 255, 0.55);
  pointer-events: none;
}

.dark .command-menu-surface::before {
  border-color: rgba(255, 255, 255, 0.04);
}

.command-menu-grid {
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.06) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.06) 1px, transparent 1px);
  background-position: center;
  background-size: 28px 28px;
}

.command-menu-panel-enter-active,
.command-menu-panel-leave-active {
  transition:
    opacity 0.22s ease,
    transform 0.22s ease,
    filter 0.22s ease;
}

.command-menu-panel-enter-from,
.command-menu-panel-leave-to {
  opacity: 0;
  filter: blur(6px);
  transform: translateY(-10px) scale(0.985);
}
</style>

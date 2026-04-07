<template>
  <div
    class="console-create-card bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6 space-y-5"
  >
    <h3 class="text-base font-bold flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t(titleKey) }}
    </h3>

    <div v-for="group in groups" :key="group.key">
      <label class="block text-sm text-slate-500 mb-3">{{
        t(group.labelKey)
      }}</label>
      <div class="flex gap-3 flex-wrap">
        <button
          :class="chipClass(filters[group.key] === '')"
          @click="$emit('change', { key: group.key, value: '' })"
        >
          {{ t(allLabelKey) }}
        </button>
        <button
          v-for="option in group.options"
          :key="getOptionValue(group, option)"
          :class="
            chipClass(filters[group.key] === getOptionValue(group, option))
          "
          @click="
            $emit('change', {
              key: group.key,
              value: getOptionValue(group, option)
            })
          "
        >
          {{ getOptionLabel(group, option) }}
          <span
            v-if="
              group.optionCountKey &&
              getOptionCount(group, option) !== undefined
            "
            class="text-xs opacity-70"
          >
            ({{ getOptionCount(group, option) }})
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'

type FilterValue = string | number
type FilterOption = FilterValue | Record<string, unknown>

interface FilterChipGroup {
  key: string
  labelKey: string
  options: FilterOption[]
  optionValueKey?: string
  optionLabelKey?: string
  optionCountKey?: string
}

withDefaults(
  defineProps<{
    allLabelKey?: string
    filters: Record<string, FilterValue>
    groups?: FilterChipGroup[]
    titleKey?: string
  }>(),
  {
    allLabelKey: 'all',
    groups: () => [],
    titleKey: 'resourceFilter'
  }
)

defineEmits<{
  change: [payload: { key: string; value: FilterValue | '' }]
}>()

const t = inject<Translator>('t', (key: string) => key)

const chipClass = (active: boolean) => [
  'px-4 py-2 rounded-lg text-sm border transition-all',
  active
    ? 'bg-primary text-white border-primary'
    : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
]

const getOptionValue = (
  group: FilterChipGroup,
  option: FilterOption
): FilterValue => {
  if (group.optionValueKey && option && typeof option === 'object') {
    return option[group.optionValueKey] as FilterValue
  }
  return option as FilterValue
}

const getOptionLabel = (
  group: FilterChipGroup,
  option: FilterOption
): string | number => {
  if (group.optionLabelKey && option && typeof option === 'object') {
    return option[group.optionLabelKey] as string | number
  }
  return option as string | number
}

const getOptionCount = (
  group: FilterChipGroup,
  option: FilterOption
): number | string | undefined => {
  if (!group.optionCountKey || !option || typeof option !== 'object') {
    return undefined
  }

  return option[group.optionCountKey] as number | string | undefined
}
</script>

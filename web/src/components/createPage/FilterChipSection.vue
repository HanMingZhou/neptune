<template>
  <div
    class="console-create-card console-create-card--section space-y-5"
  >
    <h3 class="console-create-card__title">
      <span class="console-create-card__title-marker"></span>
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
          <span class="flex flex-col items-center gap-2 text-center">
            <span class="leading-5">
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
            </span>
            <span
              v-if="getOptionMetaFields(group, option).length"
              class="flex flex-wrap justify-center gap-1.5"
            >
              <span
                v-for="field in getOptionMetaFields(group, option)"
                :key="field.key"
                :class="fieldChipClass(filters[group.key] === getOptionValue(group, option))"
              >
                <span class="opacity-75">{{ field.label }}</span>
                <span class="font-semibold">{{ field.value }}</span>
              </span>
            </span>
            <span
              v-else-if="getOptionMeta(group, option)"
              class="mt-1 text-center text-[11px] leading-4 opacity-75"
            >
              {{ getOptionMeta(group, option) }}
            </span>
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
  optionMetaKey?: string
  optionMetaFieldsKey?: string
  optionValueKey?: string
  optionLabelKey?: string
  optionCountKey?: string
}

interface OptionMetaField {
  key: string
  label: string
  value: string
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
  'inline-flex items-center justify-center px-4 py-2 rounded-lg text-sm border transition-all',
  active
    ? 'bg-primary text-white border-primary'
    : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
]

const fieldChipClass = (active: boolean) => [
  'inline-flex items-center gap-1 rounded-md px-2 py-1 text-[11px] leading-none',
  active
    ? 'bg-white/15 text-white ring-1 ring-white/20'
    : 'bg-slate-50 text-slate-600 ring-1 ring-slate-200'
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

const getOptionMeta = (
  group: FilterChipGroup,
  option: FilterOption
): string | undefined => {
  if (!group.optionMetaKey || !option || typeof option !== 'object') {
    return undefined
  }

  return option[group.optionMetaKey] as string | undefined
}

const getOptionMetaFields = (
  group: FilterChipGroup,
  option: FilterOption
): OptionMetaField[] => {
  if (!group.optionMetaFieldsKey || !option || typeof option !== 'object') {
    return []
  }

  const fields = option[group.optionMetaFieldsKey] as OptionMetaField[] | undefined
  return Array.isArray(fields) ? fields : []
}
</script>

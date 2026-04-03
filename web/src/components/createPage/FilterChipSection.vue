<template>
  <div class="console-create-card bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6 space-y-5">
    <h3 class="text-base font-bold flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t(titleKey) }}
    </h3>

    <div v-for="group in groups" :key="group.key">
      <label class="block text-sm text-slate-500 mb-3">{{ t(group.labelKey) }}</label>
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
          :class="chipClass(filters[group.key] === getOptionValue(group, option))"
          @click="$emit('change', { key: group.key, value: getOptionValue(group, option) })"
        >
          {{ getOptionLabel(group, option) }}
          <span v-if="group.optionCountKey && getOptionCount(group, option) !== undefined" class="text-xs opacity-70">
            ({{ getOptionCount(group, option) }})
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const props = defineProps({
  allLabelKey: {
    type: String,
    default: 'all'
  },
  filters: {
    type: Object,
    required: true
  },
  groups: {
    type: Array,
    default: () => []
  },
  titleKey: {
    type: String,
    default: 'resourceFilter'
  }
})

defineEmits(['change'])

const t = inject('t', (key) => key)

const chipClass = (active) => [
  'px-4 py-2 rounded-lg text-sm border transition-all',
  active
    ? 'bg-primary text-white border-primary'
    : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
]

const getOptionValue = (group, option) => {
  if (group.optionValueKey && option && typeof option === 'object') {
    return option[group.optionValueKey]
  }
  return option
}

const getOptionLabel = (group, option) => {
  if (group.optionLabelKey && option && typeof option === 'object') {
    return option[group.optionLabelKey]
  }
  return option
}

const getOptionCount = (group, option) => {
  if (!group.optionCountKey || !option || typeof option !== 'object') return undefined
  return option[group.optionCountKey]
}
</script>

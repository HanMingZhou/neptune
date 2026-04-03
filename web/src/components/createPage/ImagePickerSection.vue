<template>
  <div class="console-create-card bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t(titleKey) }}
    </h3>
    <div class="flex gap-3 mb-3">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        :class="[
          'px-4 py-2 rounded-lg text-sm border transition-all relative',
          activeTab === tab.value
            ? 'bg-primary text-white border-primary'
            : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
        ]"
        @click="$emit('change-tab', tab.value)"
      >
        {{ tab.label }}
        <span v-if="showHotTag && tab.value === hotTabValue" class="absolute -top-1 -right-1 bg-red-500 text-white text-[10px] px-1.5 rounded">
          {{ t(hotLabelKey) }}
        </span>
      </button>
    </div>
    <p :class="descriptionClass">{{ descriptionText }}</p>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div
        v-for="item in items"
        :key="getItemValue(item)"
        :class="[
          'border rounded-xl p-4 cursor-pointer transition-all',
          selectedValue === getItemValue(item)
            ? 'border-primary bg-primary/5'
            : 'border-border-light dark:border-border-dark hover:border-primary hover:-translate-y-1'
        ]"
        @click="$emit('update:selectedValue', getItemValue(item))"
      >
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2 overflow-hidden mr-2">
            <span class="text-xl flex-shrink-0">📦</span>
            <el-tooltip :content="getItemLabel(item)" placement="top">
              <span :class="['font-bold text-sm truncate', uppercaseLabel ? 'uppercase' : '']">{{ getItemLabel(item) }}</span>
            </el-tooltip>
          </div>
          <span v-if="selectedValue === getItemValue(item)" class="text-primary font-bold flex-shrink-0">✓</span>
        </div>
        <p class="text-xs text-slate-500 line-clamp-2">{{ getItemDescription(item) || t('noData') }}</p>
      </div>
      <div v-if="items.length === 0" :class="emptyStateClass">
        {{ t('noData') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const props = defineProps({
  activeTab: {
    type: String,
    required: true
  },
  descriptionClass: {
    type: String,
    default: 'text-sm text-slate-500 mb-4'
  },
  descriptionKey: {
    type: String,
    default: 'description'
  },
  descriptionText: {
    type: String,
    default: ''
  },
  emptyStateClass: {
    type: String,
    default: 'col-span-full text-center py-8 text-slate-400'
  },
  hotLabelKey: {
    type: String,
    default: 'hot'
  },
  hotTabValue: {
    type: String,
    default: 'base'
  },
  items: {
    type: Array,
    default: () => []
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  selectedValue: {
    type: [Number, String],
    default: ''
  },
  showHotTag: {
    type: Boolean,
    default: true
  },
  tabs: {
    type: Array,
    default: () => []
  },
  titleKey: {
    type: String,
    default: 'selectImage'
  },
  uppercaseLabel: {
    type: Boolean,
    default: false
  },
  valueKey: {
    type: String,
    default: 'value'
  }
})

defineEmits(['change-tab', 'update:selectedValue'])

const t = inject('t', (key) => key)

const getItemValue = (item) => item?.[props.valueKey]
const getItemLabel = (item) => item?.[props.labelKey]
const getItemDescription = (item) => item?.[props.descriptionKey]
</script>

<template>
  <div
    class="console-create-card bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6"
  >
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
        <span
          v-if="showHotTag && tab.value === hotTabValue"
          class="absolute -top-1 -right-1 bg-red-500 text-white text-[10px] px-1.5 rounded"
        >
          {{ t(hotLabelKey) }}
        </span>
      </button>
    </div>
    <p :class="descriptionClass">{{ descriptionText }}</p>
    <div
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
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
              <span
                :class="[
                  'font-bold text-sm truncate',
                  uppercaseLabel ? 'uppercase' : ''
                ]"
                >{{ getItemLabel(item) }}</span
              >
            </el-tooltip>
          </div>
          <span
            v-if="selectedValue === getItemValue(item)"
            class="text-primary font-bold flex-shrink-0"
            >✓</span
          >
        </div>
        <p class="text-xs text-slate-500 line-clamp-2">
          {{ getItemDescription(item) || t('noData') }}
        </p>
      </div>
      <div v-if="items.length === 0" :class="emptyStateClass">
        {{ t('noData') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { ResourceId, Translator } from '@/types/consoleResource'

interface ImagePickerTab {
  value: string
  label: string
}

type PickerItem = Record<string, unknown>

const props = withDefaults(
  defineProps<{
    activeTab: string
    descriptionClass?: string
    descriptionKey?: string
    descriptionText?: string
    emptyStateClass?: string
    hotLabelKey?: string
    hotTabValue?: string
    items?: PickerItem[]
    labelKey?: string
    selectedValue?: ResourceId | ''
    showHotTag?: boolean
    tabs?: ImagePickerTab[]
    titleKey?: string
    uppercaseLabel?: boolean
    valueKey?: string
  }>(),
  {
    descriptionClass: 'text-sm text-slate-500 mb-4',
    descriptionKey: 'description',
    descriptionText: '',
    emptyStateClass: 'col-span-full text-center py-8 text-slate-400',
    hotLabelKey: 'hot',
    hotTabValue: 'base',
    items: () => [],
    labelKey: 'label',
    selectedValue: '',
    showHotTag: true,
    tabs: () => [],
    titleKey: 'selectImage',
    uppercaseLabel: false,
    valueKey: 'value'
  }
)

defineEmits<{
  'change-tab': [value: string]
  'update:selectedValue': [value: ResourceId]
}>()

const t = inject<Translator>('t', (key: string) => key)

const getItemValue = (item: PickerItem): ResourceId | undefined =>
  item?.[props.valueKey] as ResourceId | undefined
const getItemLabel = (item: PickerItem): string | undefined =>
  item?.[props.labelKey] as string | undefined
const getItemDescription = (item: PickerItem): string | undefined =>
  item?.[props.descriptionKey] as string | undefined
</script>

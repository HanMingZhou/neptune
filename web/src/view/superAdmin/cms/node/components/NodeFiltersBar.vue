<template>
  <BaseFilterBar>
    <template #default>
      <el-select
        v-model="clusterModel"
        :placeholder="t('cluster')"
        class="list-filter-select !w-[240px]"
      >
        <el-option
          v-for="cluster in clusters"
          :key="cluster.id"
          :label="`${cluster.name} (${cluster.area || '-'})`"
          :value="cluster.id"
        />
      </el-select>

      <div
        v-if="clusterModel"
        class="flex items-center gap-2 px-3 py-2 bg-slate-50 dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg"
      >
        <span class="material-icons text-[18px] text-slate-400">public</span>
        <span class="text-xs font-bold text-slate-600 dark:text-slate-400"
          >{{ t('area') }}:</span
        >
        <span class="text-xs font-black text-primary">{{
          currentClusterArea
        }}</span>
      </div>

      <div class="list-filter-field max-w-[280px]">
        <span
          class="material-icons absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
          >search</span
        >
        <input
          v-model="keywordModel"
          type="text"
          :placeholder="t('searchNodeDesc')"
          class="list-search-input !w-full"
          @keyup.enter="emit('search')"
        />
      </div>
    </template>

    <template #actions>
      <button
        class="list-toolbar-button list-toolbar-button--primary"
        @click="emit('search')"
      >
        <span class="material-icons text-[18px]">search</span>
        {{ t('searchQuery') }}
      </button>
      <button
        class="list-toolbar-button list-toolbar-button--secondary"
        @click="emit('reset')"
      >
        <span class="material-icons text-[18px]">autorenew</span>
        {{ t('reset') }}
      </button>
    </template>
  </BaseFilterBar>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseFilterBar from '@/components/listPage/BaseFilterBar.vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type { CmsClusterOption } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    clusters?: CmsClusterOption[]
    currentClusterArea?: string
    filterClusterId?: ResourceId | undefined
    filterKeyword?: string
  }>(),
  {
    clusters: () => [],
    currentClusterArea: '-',
    filterClusterId: undefined,
    filterKeyword: ''
  }
)

const emit = defineEmits<{
  reset: []
  search: []
  'update:filter-cluster-id': [value: ResourceId | undefined]
  'update:filter-keyword': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)

const clusterModel = computed({
  get: () => props.filterClusterId,
  set: (value?: ResourceId) => emit('update:filter-cluster-id', value)
})

const keywordModel = computed({
  get: () => props.filterKeyword,
  set: (value: string) => emit('update:filter-keyword', value)
})
</script>

<template>
  <div class="console-filter-card overflow-hidden">
    <div class="border-b border-border-light px-5 py-4 dark:border-border-dark">
      <div class="product-toolbar-top">
        <div class="product-toolbar-tabs">
          <button
            class="list-toolbar-button min-w-0 max-w-full"
            :class="
              activeTab === 1
                ? 'list-toolbar-button--primary'
                : 'list-toolbar-button--secondary'
            "
            @click="emit('update:active-tab', 1)"
          >
            <span class="material-icons mr-1 align-middle text-[18px]">memory</span>
            <span class="truncate">{{ t('computeProduct') }}</span>
          </button>
          <button
            class="list-toolbar-button min-w-0 max-w-full"
            :class="
              activeTab === 2
                ? 'list-toolbar-button--primary'
                : 'list-toolbar-button--secondary'
            "
            @click="emit('update:active-tab', 2)"
          >
            <span class="material-icons mr-1 align-middle text-[18px]">
              folder_open
            </span>
            <span class="truncate">{{ t('storageProduct') }}</span>
          </button>
        </div>

        <div class="list-toolbar-actions product-toolbar-actions">
          <RefreshButton :loading="loading" @refresh="emit('refresh', $event)" />
          <button
            class="list-toolbar-button list-toolbar-button--primary min-w-0 max-w-full"
            @click="emit('create')"
          >
            <span class="material-icons text-[20px]">add</span>
            <span class="hidden xl:inline">
              {{ activeTab === 1 ? t('newComputeProduct') : t('newStorageProduct') }}
            </span>
            <span class="xl:hidden">{{ t('create') }}</span>
          </button>
        </div>
      </div>
    </div>

    <div class="product-filter-shell">
      <div class="product-filter-row">
        <div class="product-filter-field">
          <label class="product-filter-label">{{ t('cluster') }}</label>
          <el-select
            v-model="clusterModel"
            :placeholder="t('selectCluster')"
            clearable
            class="!w-full gva-custom-select"
          >
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
        </div>

        <div class="product-filter-field">
          <label class="product-filter-label">{{ t('area') }}</label>
          <el-select
            v-model="areaModel"
            :placeholder="t('inference.selectArea')"
            clearable
            class="!w-full gva-custom-select"
          >
            <el-option
              v-for="area in areas"
              :key="area"
              :label="area"
              :value="area"
            />
          </el-select>
        </div>

        <div class="product-filter-field product-filter-field--keyword">
          <label class="product-filter-label">{{ t('searchQuery') }}</label>
          <el-input
            v-model="keywordModel"
            :placeholder="t('searchProductDesc')"
            clearable
            @keyup.enter="emit('search')"
          >
            <template #prefix>
              <span class="material-icons text-[18px] text-slate-400">search</span>
            </template>
          </el-input>
        </div>

        <div class="product-filter-actions">
          <button
            class="list-toolbar-button list-toolbar-button--secondary"
            :class="{
              'product-advanced-toggle--active':
                showAdvanced || activeAdvancedFilterCount > 0
            }"
            @click="showAdvanced = !showAdvanced"
          >
            <span class="material-icons text-[18px]">
              {{ showAdvanced ? 'expand_less' : 'tune' }}
            </span>
            {{ showAdvanced ? t('collapseFilters') : t('advancedSearch') }}
            <span
              v-if="activeAdvancedFilterCount > 0"
              class="product-filter-count"
            >
              {{ activeAdvancedFilterCount }}
            </span>
          </button>
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
        </div>
      </div>

      <div v-show="showAdvanced" class="product-advanced-panel">
        <div class="product-advanced-header">
          <div class="product-advanced-title">
            <span class="material-icons text-[16px]">tune</span>
            <span>{{ t('advancedSearch') }}</span>
          </div>
          <span v-if="activeAdvancedFilterCount > 0" class="product-advanced-meta">
            {{ activeAdvancedFilterCount }} {{ t('searchQuery') }}
          </span>
        </div>

        <div class="product-advanced-grid">
          <div
            v-if="activeTab === 1"
            class="product-filter-field product-filter-field--panel"
          >
            <label class="product-filter-label">{{ t('resourceType') }}</label>
            <div class="product-resource-type-group">
              <button
                v-for="option in resourceTypeOptions"
                :key="option.value || 'all'"
                type="button"
                class="product-resource-type-button"
                :class="[
                  `product-resource-type-button--${option.tone}`,
                  {
                    'is-active': resourceTypeModel === option.value
                  }
                ]"
                @click="resourceTypeModel = option.value"
              >
                {{ option.label }}
              </button>
            </div>
          </div>

          <div
            v-if="activeTab === 1 && resourceTypeModel !== 'cpu'"
            class="product-filter-field product-filter-field--panel"
          >
            <label class="product-filter-label">{{ t('gpuModel') }}</label>
            <el-select
              v-model="gpuModelModel"
              filterable
              :placeholder="t('selectGpuModel')"
              clearable
              class="!w-full gva-custom-select"
            >
              <el-option
                v-for="gpu in gpuModels"
                :key="gpu.model"
                :label="gpu.model"
                :value="gpu.model"
              />
            </el-select>
          </div>

          <div class="product-filter-field product-filter-field--panel">
            <label class="product-filter-label">{{ t('priceField') }}</label>
            <el-select
              v-model="priceFieldModel"
              class="!w-full gva-custom-select"
            >
              <el-option
                v-for="option in priceFieldOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </div>

          <div class="product-filter-field product-filter-field--range">
            <label class="product-filter-label">{{ t('remainingInventory') }}</label>
            <div class="product-range-row">
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('minValue') }}</span>
                <el-input-number
                  v-model="availableMinModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
              <span class="product-range-divider">-</span>
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('maxValue') }}</span>
                <el-input-number
                  v-model="availableMaxModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
            </div>
          </div>

          <div class="product-filter-field product-filter-field--range">
            <label class="product-filter-label">{{ t('usedInstances') }}</label>
            <div class="product-range-row">
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('minValue') }}</span>
                <el-input-number
                  v-model="usedCapacityMinModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
              <span class="product-range-divider">-</span>
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('maxValue') }}</span>
                <el-input-number
                  v-model="usedCapacityMaxModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
            </div>
          </div>

          <div class="product-filter-field product-filter-field--range">
            <label class="product-filter-label">{{ t('maxInstances') }}</label>
            <div class="product-range-row">
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('minValue') }}</span>
                <el-input-number
                  v-model="maxInstancesMinModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
              <span class="product-range-divider">-</span>
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('maxValue') }}</span>
                <el-input-number
                  v-model="maxInstancesMaxModel"
                  :controls="false"
                  :min="0"
                  class="!w-full"
                />
              </div>
            </div>
          </div>

          <div class="product-filter-field product-filter-field--range">
            <label class="product-filter-label">{{ t('price') }}</label>
            <div class="product-range-row">
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('minValue') }}</span>
                <el-input-number
                  v-model="priceMinModel"
                  :controls="false"
                  :min="0"
                  :precision="activeTab === 2 ? 4 : 2"
                  class="!w-full"
                />
              </div>
              <span class="product-range-divider">-</span>
              <div class="product-range-item">
                <span class="product-range-hint">{{ t('maxValue') }}</span>
                <el-input-number
                  v-model="priceMaxModel"
                  :controls="false"
                  :min="0"
                  :precision="activeTab === 2 ? 4 : 2"
                  class="!w-full"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import type {
  FilterOption,
  ResourceId,
  Translator
} from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsProductFilterResourceType,
  CmsProductPriceField,
  CmsProductType
} from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    activeTab?: CmsProductType
    areas?: string[]
    clusters?: CmsClusterOption[]
    filterArea?: string
    filterAvailableMax?: number
    filterAvailableMin?: number
    filterClusterId?: ResourceId | ''
    filterResourceType?: CmsProductFilterResourceType
    filterGpuModel?: string
    filterKeyword?: string
    filterMaxInstancesMax?: number
    filterMaxInstancesMin?: number
    filterPriceField?: CmsProductPriceField
    filterPriceMax?: number
    filterPriceMin?: number
    filterUsedCapacityMax?: number
    filterUsedCapacityMin?: number
    gpuModels?: FilterOption[]
    loading?: boolean
  }>(),
  {
    activeTab: 1,
    areas: () => [],
    clusters: () => [],
    filterArea: '',
    filterAvailableMax: undefined,
    filterAvailableMin: undefined,
    filterClusterId: '',
    filterResourceType: '',
    filterGpuModel: '',
    filterKeyword: '',
    filterMaxInstancesMax: undefined,
    filterMaxInstancesMin: undefined,
    filterPriceField: 1,
    filterPriceMax: undefined,
    filterPriceMin: undefined,
    filterUsedCapacityMax: undefined,
    filterUsedCapacityMin: undefined,
    gpuModels: () => [],
    loading: false
  }
)

const emit = defineEmits<{
  create: []
  refresh: [silent: boolean]
  reset: []
  search: []
  'update:active-tab': [value: CmsProductType]
  'update:filter-area': [value: string]
  'update:filter-available-max': [value?: number]
  'update:filter-available-min': [value?: number]
  'update:filter-cluster-id': [value: ResourceId | '']
  'update:filter-resource-type': [value: CmsProductFilterResourceType]
  'update:filter-gpu-model': [value: string]
  'update:filter-keyword': [value: string]
  'update:filter-max-instances-max': [value?: number]
  'update:filter-max-instances-min': [value?: number]
  'update:filter-price-field': [value: CmsProductPriceField]
  'update:filter-price-max': [value?: number]
  'update:filter-price-min': [value?: number]
  'update:filter-used-capacity-max': [value?: number]
  'update:filter-used-capacity-min': [value?: number]
}>()

const t = inject<Translator>('t', (key: string) => key)
const showAdvanced = ref(false)

const activeAdvancedFilterCount = computed(() => {
  const filters = [
    props.activeTab === 1 && !!props.filterResourceType,
    props.activeTab === 1 && !!props.filterGpuModel,
    props.filterAvailableMin !== undefined,
    props.filterAvailableMax !== undefined,
    props.filterUsedCapacityMin !== undefined,
    props.filterUsedCapacityMax !== undefined,
    props.filterMaxInstancesMin !== undefined,
    props.filterMaxInstancesMax !== undefined,
    props.filterPriceMin !== undefined,
    props.filterPriceMax !== undefined,
    props.activeTab === 1 ? props.filterPriceField !== 1 : props.filterPriceField !== 5
  ]

  return filters.filter(Boolean).length
})

const clusterModel = computed({
  get: () => props.filterClusterId,
  set: (value: ResourceId | '' | undefined) =>
    emit('update:filter-cluster-id', value ?? '')
})

const areaModel = computed({
  get: () => props.filterArea,
  set: (value?: string) => emit('update:filter-area', value ?? '')
})

const gpuModelModel = computed({
  get: () => props.filterGpuModel,
  set: (value?: string) => emit('update:filter-gpu-model', value ?? '')
})

const resourceTypeModel = computed({
  get: () => props.filterResourceType,
  set: (value: CmsProductFilterResourceType) => {
    emit('update:filter-resource-type', value)
    if (value === 'cpu' && props.filterGpuModel) {
      emit('update:filter-gpu-model', '')
    }
  }
})

const keywordModel = computed({
  get: () => props.filterKeyword,
  set: (value: string) => emit('update:filter-keyword', value)
})

const availableMinModel = computed({
  get: () => props.filterAvailableMin,
  set: (value?: number) => emit('update:filter-available-min', value)
})

const availableMaxModel = computed({
  get: () => props.filterAvailableMax,
  set: (value?: number) => emit('update:filter-available-max', value)
})

const usedCapacityMinModel = computed({
  get: () => props.filterUsedCapacityMin,
  set: (value?: number) => emit('update:filter-used-capacity-min', value)
})

const usedCapacityMaxModel = computed({
  get: () => props.filterUsedCapacityMax,
  set: (value?: number) => emit('update:filter-used-capacity-max', value)
})

const maxInstancesMinModel = computed({
  get: () => props.filterMaxInstancesMin,
  set: (value?: number) => emit('update:filter-max-instances-min', value)
})

const maxInstancesMaxModel = computed({
  get: () => props.filterMaxInstancesMax,
  set: (value?: number) => emit('update:filter-max-instances-max', value)
})

const priceFieldModel = computed({
  get: () => props.filterPriceField,
  set: (value?: CmsProductPriceField) =>
    emit('update:filter-price-field', value ?? 1)
})

const priceMinModel = computed({
  get: () => props.filterPriceMin,
  set: (value?: number) => emit('update:filter-price-min', value)
})

const priceMaxModel = computed({
  get: () => props.filterPriceMax,
  set: (value?: number) => emit('update:filter-price-max', value)
})

const priceFieldOptions = computed(() =>
  props.activeTab === 2
    ? [
        {
          label: t('storagePriceGb'),
          value: 5 as CmsProductPriceField
        }
      ]
    : [
        { label: t('priceHourly'), value: 1 as CmsProductPriceField },
        { label: t('priceDaily'), value: 2 as CmsProductPriceField },
        { label: t('priceWeekly'), value: 3 as CmsProductPriceField },
        {
          label: t('priceMonthly'),
          value: 4 as CmsProductPriceField
        }
      ]
)

const resourceTypeOptions = computed<
  Array<{
    label: string
    tone: 'neutral' | 'cpu' | 'gpu' | 'vgpu'
    value: CmsProductFilterResourceType
  }>
>(() => [
  {
    label: t('all'),
    tone: 'neutral',
    value: ''
  },
  {
    label: t('cpuOnly'),
    tone: 'cpu',
    value: 'cpu'
  },
  {
    label: t('gpu'),
    tone: 'gpu',
    value: 'gpu'
  },
  {
    label: 'vGPU',
    tone: 'vgpu',
    value: 'vgpu'
  }
])
</script>

<style scoped>
.product-toolbar-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  flex-wrap: wrap;
}

.product-toolbar-tabs {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px;
  border: 1px solid rgb(226 232 240);
  border-radius: 18px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(248 250 252));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.85);
}

.product-toolbar-actions {
  justify-content: flex-end;
  gap: 10px;
}

.product-filter-shell {
  padding: 18px 20px 20px;
  background:
    linear-gradient(180deg, rgb(248 250 252 / 0.96), rgb(255 255 255)),
    radial-gradient(circle at top right, rgb(191 219 254 / 0.45), transparent 34%);
}

.product-filter-row {
  display: grid;
  grid-template-columns: 220px 180px minmax(320px, 1fr) auto;
  gap: 14px;
  align-items: end;
}

.product-filter-field {
  min-width: 0;
}

.product-filter-field--keyword {
  min-width: 260px;
}

.product-filter-label {
  display: inline-flex;
  margin-bottom: 8px;
  color: rgb(100 116 139);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  line-height: 1;
  text-transform: uppercase;
}

.product-filter-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  flex-wrap: wrap;
}

.product-filter-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 999px;
  background: rgb(15 23 42);
  color: rgb(255 255 255);
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

.product-advanced-toggle--active {
  border-color: rgb(148 163 184) !important;
  background: rgb(241 245 249) !important;
  color: rgb(15 23 42) !important;
}

.product-advanced-panel {
  margin-top: 16px;
  padding-top: 18px;
  border-top: 1px solid rgb(226 232 240);
}

.product-advanced-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.product-advanced-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: rgb(30 41 59);
  font-size: 13px;
  font-weight: 700;
  line-height: 1;
}

.product-advanced-meta {
  color: rgb(100 116 139);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.product-advanced-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.product-filter-field--panel,
.product-filter-field--range {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 96px;
  padding: 12px 12px 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 16px;
  background: rgb(255 255 255 / 0.7);
}

.product-resource-type-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.product-resource-type-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 34px;
  padding: 0 12px;
  border: 1px solid rgb(203 213 225);
  border-radius: 999px;
  background: rgb(255 255 255 / 0.92);
  color: rgb(71 85 105);
  font-size: 12px;
  font-weight: 700;
  line-height: 1;
  transition:
    border-color 0.18s ease,
    background-color 0.18s ease,
    color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.product-resource-type-button:hover {
  transform: translateY(-1px);
  border-color: rgb(148 163 184);
  color: rgb(15 23 42);
}

.product-resource-type-button.is-active {
  box-shadow: 0 10px 24px rgb(15 23 42 / 0.08);
}

.product-resource-type-button--neutral.is-active {
  border-color: rgb(100 116 139);
  background: linear-gradient(180deg, rgb(255 255 255), rgb(241 245 249));
  color: rgb(15 23 42);
}

.product-resource-type-button--cpu.is-active {
  border-color: rgb(148 163 184);
  background: linear-gradient(180deg, rgb(248 250 252), rgb(226 232 240));
  color: rgb(51 65 85);
}

.product-resource-type-button--gpu.is-active {
  border-color: rgb(245 158 11);
  background: linear-gradient(180deg, rgb(255 251 235), rgb(254 243 199));
  color: rgb(180 83 9);
}

.product-resource-type-button--vgpu.is-active {
  border-color: rgb(6 182 212);
  background: linear-gradient(180deg, rgb(236 254 255), rgb(207 250 254));
  color: rgb(14 116 144);
}

.product-range-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: 10px;
  align-items: center;
}

.product-range-item {
  min-width: 0;
}

.product-range-hint {
  display: inline-flex;
  margin-bottom: 6px;
  color: rgb(148 163 184);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.product-range-divider {
  color: rgb(148 163 184);
  font-size: 16px;
  font-weight: 700;
  line-height: 1;
  transform: translateY(10px);
}

.product-filter-field :deep(.el-select__wrapper),
.product-filter-field :deep(.el-input__wrapper),
.product-filter-field :deep(.el-input-number) {
  border-radius: 14px;
  background: rgb(255 255 255 / 0.88);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.92),
    0 1px 2px rgb(15 23 42 / 0.04);
}

.product-filter-field :deep(.el-input-number .el-input__wrapper) {
  box-shadow: none;
}

.product-filter-field :deep(.el-select__wrapper.is-focused),
.product-filter-field :deep(.el-input__wrapper.is-focus),
.product-filter-field :deep(.el-input-number.is-focus) {
  box-shadow:
    0 0 0 1px rgb(59 130 246 / 0.24),
    0 10px 24px rgb(59 130 246 / 0.08);
}

.dark .product-filter-shell {
  background:
    linear-gradient(180deg, rgb(24 24 27 / 0.98), rgb(15 23 42 / 0.94)),
    radial-gradient(circle at top right, rgb(30 64 175 / 0.28), transparent 36%);
}

.dark .product-toolbar-tabs {
  border-color: rgb(51 65 85);
  background: linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.04);
}

.dark .product-filter-label {
  color: rgb(148 163 184);
}

.dark .product-filter-count {
  background: rgb(96 165 250);
  color: rgb(15 23 42);
}

.dark .product-advanced-toggle--active {
  border-color: rgb(71 85 105) !important;
  background: rgb(39 39 42) !important;
  color: rgb(226 232 240) !important;
}

.dark .product-advanced-panel {
  border-top-color: rgb(51 65 85);
}

.dark .product-advanced-title {
  color: rgb(226 232 240);
}

.dark .product-advanced-meta {
  color: rgb(148 163 184);
}

.dark .product-filter-field--panel,
.dark .product-filter-field--range {
  border-color: rgb(63 63 70);
  background: rgb(24 24 27 / 0.62);
}

.dark .product-resource-type-button {
  border-color: rgb(71 85 105);
  background: rgb(24 24 27 / 0.84);
  color: rgb(203 213 225);
}

.dark .product-resource-type-button:hover {
  border-color: rgb(100 116 139);
  color: white;
}

.dark .product-resource-type-button.is-active {
  box-shadow: 0 12px 24px rgb(2 6 23 / 0.24);
}

.dark .product-resource-type-button--neutral.is-active {
  border-color: rgb(148 163 184);
  background: linear-gradient(180deg, rgb(51 65 85), rgb(30 41 59));
  color: rgb(241 245 249);
}

.dark .product-resource-type-button--cpu.is-active {
  border-color: rgb(100 116 139);
  background: linear-gradient(180deg, rgb(71 85 105), rgb(51 65 85));
  color: rgb(226 232 240);
}

.dark .product-resource-type-button--gpu.is-active {
  border-color: rgb(245 158 11);
  background: linear-gradient(180deg, rgb(146 64 14), rgb(120 53 15));
  color: rgb(255 237 213);
}

.dark .product-resource-type-button--vgpu.is-active {
  border-color: rgb(34 211 238);
  background: linear-gradient(180deg, rgb(14 116 144), rgb(8 145 178));
  color: rgb(207 250 254);
}

.dark .product-range-hint,
.dark .product-range-divider {
  color: rgb(113 113 122);
}

.dark .product-filter-field :deep(.el-select__wrapper),
.dark .product-filter-field :deep(.el-input__wrapper),
.dark .product-filter-field :deep(.el-input-number) {
  background: rgb(24 24 27 / 0.86);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.03),
    0 1px 2px rgb(0 0 0 / 0.24);
}

@media (max-width: 1440px) {
  .product-filter-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .product-filter-actions {
    grid-column: 1 / -1;
    justify-content: flex-start;
  }

  .product-advanced-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .product-toolbar-top,
  .product-filter-row,
  .product-advanced-grid {
    grid-template-columns: 1fr;
  }

  .product-toolbar-tabs,
  .product-filter-field--keyword {
    min-width: 0;
  }

  .product-toolbar-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .product-range-row {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .product-range-divider {
    display: none;
  }
}
</style>

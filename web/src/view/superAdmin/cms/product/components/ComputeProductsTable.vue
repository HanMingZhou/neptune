<template>
  <TableCard :page-size="pageSize">
    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table
        class="console-table console-table--compact product-table w-full min-w-[1980px]"
        v-loading="loading"
      >
        <colgroup>
          <col style="width: 96px" />
          <col style="width: 200px" />
          <col style="width: 120px" />
          <col style="width: 140px" />
          <col style="width: 170px" />
          <col style="width: 150px" />
          <col style="width: 120px" />
          <col style="width: 280px" />
          <col style="width: 180px" />
          <col style="width: 140px" />
          <col style="width: 220px" />
          <col style="width: 320px" />
          <col style="width: 228px" />
        </colgroup>
        <thead>
          <tr>
            <th>{{ t('id') }}</th>
            <th class="product-column">{{ t('productName') }}</th>
            <th class="text-center">{{ t('status') }}</th>
            <th class="cluster-column">{{ t('clusterName') }}</th>
            <th class="node-column">{{ t('nodeName') }}</th>
            <th class="node-ip-column">{{ t('nodeIp') }}</th>
            <th>{{ t('area') }}</th>
            <th class="spec-column">{{ t('spec') }}</th>
            <th class="resource-column">{{ t('resourceConfig') }}</th>
            <th class="cuda-column">{{ t('cudaVersion') }}</th>
            <th class="text-center inventory-column">{{ t('inventory') }}</th>
            <th class="price-column">{{ t('prices') }}({{ t('priceHourly') }}/{{ t('priceDaily') }}/{{ t('priceWeekly') }}/{{ t('priceMonthly') }})</th>
            <th class="console-actions-header product-actions-column">
              {{ t('actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-light dark:divide-border-dark">
          <tr
            v-for="row in items"
            :key="row.id"
            class="transition-colors hover:bg-slate-50 dark:hover:bg-zinc-800/40"
          >
            <td class="is-code is-secondary">{{ row.id }}</td>
            <td class="product-cell">
              <span class="is-primary product-name">{{ row.name }}</span>
            </td>
            <td class="text-center">
              <ListToneBadge
                :label="row.status === 1 ? t('onShelf') : t('offShelf')"
                :tone="row.status === 1 ? 'success' : 'danger'"
              />
            </td>
            <td class="cluster-cell">
              <span class="meta-pill">{{ row.clusterName || '-' }}</span>
            </td>
            <td class="node-cell">
              <span class="meta-pill meta-pill--code">{{ row.nodeName }}</span>
            </td>
            <td class="node-ip-cell">
              <span class="meta-pill meta-pill--code">
                {{ getNodeIpText(row) }}
              </span>
            </td>
            <td class="is-secondary">{{ row.area }}</td>
            <td class="spec-cell">
              <div class="gpu-config-panel" :title="getSpecSummaryText(row)">
                <div
                  v-if="hasGpuSpec(row) || hasVGpuSpec(row)"
                  class="gpu-config-stack"
                >
                  <div class="gpu-config-title-row">
                    <div class="gpu-config-title">
                      {{ row.gpuModel || t('gpu') }}
                    </div>
                    <span v-if="hasVGpuSpec(row)" class="gpu-config-type-badge"
                      >vGPU</span
                    >
                  </div>
                  <div class="gpu-config-metric-grid">
                    <span
                      v-for="entry in getSpecEntries(row)"
                      :key="entry.key"
                      class="gpu-config-metric"
                    >
                      <span class="gpu-config-metric__label">{{
                        entry.label
                      }}</span>
                      <span class="gpu-config-metric__value">{{
                        entry.value
                      }}</span>
                    </span>
                  </div>
                </div>
                <div v-else class="gpu-config-empty">CPU ONLY</div>
              </div>
            </td>
            <td class="resource-cell">
              <div class="resource-surface">
                <div class="resource-chip-group">
                  <span
                    v-for="entry in getResourceConfigEntries(row)"
                    :key="entry.key"
                    class="resource-chip"
                  >
                    <span class="resource-chip__label">{{ entry.label }}</span>
                    <span class="resource-chip__value">{{ entry.value }}</span>
                  </span>
                </div>
              </div>
            </td>
            <td class="cuda-cell">
              <span class="meta-pill meta-pill--code">
                {{ getCudaVersionText(row) }}
              </span>
            </td>
            <td class="text-center inventory-cell">
              <div class="metric-cell-shell">
                <div
                  class="inline-metrics inline-metrics--surface inventory-metrics"
                >
                  <span class="inline-metric">
                    <span class="inline-metric__label">{{
                      t('usedInstances')
                    }}</span>
                    <span
                      class="inline-metric__value inline-metric__value--info"
                    >
                      {{ row.usedCapacity ?? 0 }}
                    </span>
                  </span>
                  <span class="inline-metrics__divider">|</span>
                  <span class="inline-metric">
                    <span class="inline-metric__label">{{
                      t('maxInstances')
                    }}</span>
                    <span class="inline-metric__value">
                      {{ row.maxInstances || 0 }}
                    </span>
                  </span>
                  <span class="inline-metrics__divider">|</span>
                  <span class="inline-metric">
                    <span class="inline-metric__label">{{
                      t('remainingInventory')
                    }}</span>
                    <span
                      class="inline-metric__value inline-metric__value--success"
                    >
                      {{ row.available ?? 0 }}
                    </span>
                  </span>
                </div>
              </div>
            </td>
            <td class="price-cell">
              <div class="metric-cell-shell">
                <div
                  class="inline-metrics inline-metrics--surface price-metrics-grid"
                  :title="getPriceSummaryText(row)"
                >
                  <span class="inline-metric price-metric">
                    <span class="inline-metric__label">H</span>
                    <span
                      class="inline-metric__value inline-metric__value--price"
                      :title="formatPriceValue(row.priceHourly)"
                    >
                      {{ formatPriceValue(row.priceHourly) }}
                    </span>
                  </span>
                  <span class="inline-metric price-metric">
                    <span class="inline-metric__label">D</span>
                    <span
                      class="inline-metric__value inline-metric__value--price"
                      :title="formatPriceValue(row.priceDaily)"
                    >
                      {{ formatPriceValue(row.priceDaily) }}
                    </span>
                  </span>
                  <span class="inline-metric price-metric">
                    <span class="inline-metric__label">W</span>
                    <span
                      class="inline-metric__value inline-metric__value--price"
                      :title="formatPriceValue(row.priceWeekly)"
                    >
                      {{ formatPriceValue(row.priceWeekly) }}
                    </span>
                  </span>
                  <span class="inline-metric price-metric">
                    <span class="inline-metric__label">M</span>
                    <span
                      class="inline-metric__value inline-metric__value--price"
                      :title="formatPriceValue(row.priceMonthly)"
                    >
                      {{ formatPriceValue(row.priceMonthly) }}
                    </span>
                  </span>
                </div>
              </div>
            </td>
            <td class="console-actions-cell product-actions-cell">
              <div class="list-row-actions">
                <button
                  class="list-row-button list-row-button--info"
                  @click="$emit('edit', row)"
                >
                  {{ t('edit') }}
                </button>
                <button
                  class="list-row-button list-row-button--warning"
                  @click="$emit('adjust-price', row)"
                >
                  {{ t('prices') }}
                </button>
                <button
                  class="list-row-button list-row-button--danger"
                  @click="$emit('delete', row)"
                >
                  {{ t('delete') }}
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="items.length === 0">
            <td
              colspan="13"
              class="px-6 py-12 text-center text-sm text-slate-400"
            >
              {{ t('noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <ListPaginationBar
        v-model:current-page="pageModel"
        v-model:page-size="pageSizeModel"
        :total="total"
        :total-text="t('totalRecords', { total })"
        :page-sizes="[15, 20, 50, 100]"
        layout="sizes, prev, pager, next, jumper"
        @current-change="$emit('page-change', $event)"
        @size-change="$emit('size-change', $event)"
      />
    </template>
  </TableCard>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import ListPaginationBar from '@/components/listPage/ListPaginationBar.vue'
import ListToneBadge from '@/components/listPage/ListToneBadge.vue'
import TableCard from '@/components/listPage/TableCard.vue'
import type { Translator } from '@/types/consoleResource'
import type { CmsProductRow } from '@/types/superAdmin'
import {
  buildGpuFieldEntries,
  buildVGpuFieldEntries,
  formatVGpuSpec,
  hasGpuSpec,
  hasVGpuSpec
} from '@/utils/vgpu'

const props = withDefaults(
  defineProps<{
    items?: CmsProductRow[]
    loading?: boolean
    page?: number
    pageSize?: number
    total?: number
  }>(),
  {
    items: () => [],
    loading: false,
    page: 1,
    pageSize: 15,
    total: 0
  }
)

const emit = defineEmits<{
  'adjust-price': [row: CmsProductRow]
  delete: [row: CmsProductRow]
  edit: [row: CmsProductRow]
  'page-change': [page: number]
  'size-change': [pageSize: number]
  'update:page': [value: number]
  'update:page-size': [value: number]
}>()
const t = inject<Translator>('t', (key: string) => key)

const pageModel = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:page-size', value)
})

type SpecEntry = {
  key: string
  label: string
  value: string
}

const getTranslatedSpecLabel = (key: string, fallback: string): string => {
  const label = t(key)
  return label && label !== key ? label : fallback
}

const getSpecLabel = (entry: SpecEntry): string => {
  if (entry.key === 'count') {
    return getTranslatedSpecLabel('gpuCountShort', entry.label)
  }
  if (entry.key === 'memory') {
    return getTranslatedSpecLabel('gpuMemoryShort', entry.label)
  }
  if (entry.key === 'cores') {
    return getTranslatedSpecLabel('vGpuCoresShort', entry.label)
  }
  return entry.label
}

const getGpuSpecDisplayText = (row: CmsProductRow): string => {
  const gpuModel = row.gpuModel?.trim() || 'GPU'

  if ((row.gpuCount || 0) > 0) {
    return row.gpuMemory
      ? `${gpuModel} × ${row.gpuCount || 0} · ${row.gpuMemory}GB`
      : `${gpuModel} × ${row.gpuCount || 0}`
  }

  if (hasVGpuSpec(row)) {
    const vGpuLabel = formatVGpuSpec(row)
    return row.gpuModel ? `${gpuModel} vGPU ${vGpuLabel}` : `vGPU ${vGpuLabel}`
  }

  return 'CPU ONLY'
}

const getCpuSpecText = (row: CmsProductRow): string => `${row.cpu}C`

const getMemorySpecText = (row: CmsProductRow): string => `${row.memory}GB`

const getGpuFieldEntries = (row: CmsProductRow): SpecEntry[] =>
  buildGpuFieldEntries(row, t).map((entry) => ({
    key: entry.key,
    label: entry.label,
    value: entry.value
  }))

const getVGpuFieldEntries = (row: CmsProductRow): SpecEntry[] =>
  buildVGpuFieldEntries(row, t).map((entry) => ({
    key: entry.key,
    label: entry.label,
    value: entry.value
  }))

const getSpecEntries = (row: CmsProductRow): SpecEntry[] =>
  (hasVGpuSpec(row) ? getVGpuFieldEntries(row) : getGpuFieldEntries(row)).map(
    (entry) => ({
      ...entry,
      label: getSpecLabel(entry)
    })
  )

const getResourceConfigEntries = (row: CmsProductRow): SpecEntry[] => [
  {
    key: 'cpu',
    label: 'CPU',
    value: getCpuSpecText(row)
  },
  {
    key: 'memory',
    label: t('memory'),
    value: getMemorySpecText(row)
  }
]

const getNodeIpText = (row: CmsProductRow): string => row.nodeIp?.trim() || '-'

const getCudaVersionText = (row: CmsProductRow): string =>
  row.cudaVersion?.trim() || '-'

const getSpecSummaryText = (row: CmsProductRow): string =>
  `${getGpuSpecDisplayText(row)} | CPU ${getCpuSpecText(row)} | Memory ${getMemorySpecText(row)}`

const formatPriceValue = (price?: number | null): string =>
  `¥${(price ?? 0).toFixed(2)}`

const getPriceSummaryText = (row: CmsProductRow): string =>
  `H ${formatPriceValue(row.priceHourly)} | D ${formatPriceValue(row.priceDaily)} | W ${formatPriceValue(row.priceWeekly)} | M ${formatPriceValue(row.priceMonthly)}`
</script>

<style scoped>
.product-column {
  min-width: 180px;
}

.cluster-column {
  min-width: 140px;
}

.node-column {
  min-width: 170px;
}

.node-ip-column {
  min-width: 150px;
}

.spec-column,
.spec-cell {
  width: 270px;
  min-width: 240px;
}

.resource-column,
.resource-cell {
  width: 200px;
  min-width: 200px;
}

.cuda-column {
  min-width: 140px;
}

.inventory-column {
  min-width: 220px;
}

.price-column {
  width: 320px;
  min-width: 320px;
}

th.product-actions-column,
td.product-actions-cell {
  width: 228px !important;
  min-width: 228px !important;
  max-width: 228px !important;
  box-sizing: border-box;
  overflow: hidden;
}

.product-table {
  table-layout: fixed;
}

.product-cell,
.cluster-cell,
.node-cell,
.node-ip-cell,
.spec-cell,
.resource-cell,
.cuda-cell,
.inventory-cell,
.price-cell {
  vertical-align: middle;
}

.inventory-cell,
.price-cell {
  overflow: hidden;
}

.metric-cell-shell {
  display: flex;
  width: 100%;
  max-width: 100%;
  justify-content: center;
  overflow: hidden;
}

.product-name {
  display: block;
  overflow: hidden;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta-pill {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  padding: 0 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 999px;
  background: rgb(248 250 252);
  color: rgb(71 85 105);
  font-size: 12px;
  font-weight: 600;
  line-height: 1;
}

.meta-pill--code {
  font-family:
    'JetBrains Mono', 'Fira Code', Consolas, 'SFMono-Regular', Monaco, monospace;
}

.spec-cell {
  white-space: normal;
}

.gpu-config-panel {
  display: flex;
  width: 100%;
  min-height: 36px;
  align-items: center;
  justify-content: center;
  padding: 7px 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 16px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(248 250 252));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.86);
}

.gpu-config-stack {
  display: flex;
  width: 100%;
  min-width: 0;
  flex-direction: column;
  gap: 5px;
}

.gpu-config-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.gpu-config-title {
  min-width: 0;
  overflow: hidden;
  color: var(--color-primary);
  font-size: 11px;
  font-weight: 700;
  line-height: 1.1;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gpu-config-type-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  border-radius: 999px;
  background: rgb(236 254 255);
  padding: 0.125rem 0.375rem;
  color: rgb(14 116 144);
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

.gpu-config-metric-grid {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(0, 1fr);
  align-items: center;
  width: 100%;
  gap: 8px;
}

.gpu-config-metric {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-width: 0;
  min-height: 24px;
  padding: 0 8px;
  border: 1px solid rgb(226 232 240);
  border-radius: 10px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(241 245 249));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.92);
}

.gpu-config-metric__label {
  color: rgb(148 163 184);
  font-size: 9px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.02em;
  text-transform: none;
  white-space: nowrap;
}

.gpu-config-metric__value {
  color: rgb(51 65 85);
  font-size: 11px;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.gpu-config-empty {
  color: rgb(148 163 184);
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.spec-surface,
.inventory-cell .inline-metrics--surface,
.price-cell .inline-metrics--surface {
  display: inline-flex;
  align-items: center;
  min-height: 30px;
  max-width: 100%;
  min-width: 0;
  padding: 0 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 14px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(248 250 252));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.86);
  overflow: hidden;
}

.spec-surface {
  width: 100%;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  gap: 6px;
  padding: 8px 10px;
  text-align: center;
}

.resource-surface {
  display: flex;
  width: 100%;
  min-height: 36px;
  align-items: center;
  justify-content: center;
  padding: 7px 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 16px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(248 250 252));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.86);
}

.spec-surface--gpu {
  border-color: rgb(253 230 138);
  background:
    radial-gradient(circle at top right, rgb(255 247 237), transparent 42%),
    linear-gradient(180deg, rgb(255 255 255), rgb(255 251 235));
}

.spec-surface--vgpu {
  border-color: rgb(165 243 252);
  background:
    radial-gradient(circle at top right, rgb(240 253 250), transparent 42%),
    linear-gradient(180deg, rgb(255 255 255), rgb(236 254 255));
}

.spec-header {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 6px;
  width: 100%;
  min-width: 0;
}

.spec-kind-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  height: 22px;
  padding: 0 8px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  white-space: nowrap;
}

.spec-kind-badge--gpu {
  background: linear-gradient(135deg, rgb(245 158 11), rgb(251 191 36));
  color: rgb(120 53 15);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.35);
}

.spec-kind-badge--vgpu {
  background: linear-gradient(135deg, rgb(34 211 238), rgb(45 212 191));
  color: rgb(22 78 99);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.35);
}

.spec-kind-badge--cpu {
  background: linear-gradient(135deg, rgb(226 232 240), rgb(203 213 225));
  color: rgb(71 85 105);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.5);
}

.spec-model {
  min-width: 0;
  overflow: hidden;
  font-size: 12px;
  font-weight: 800;
  line-height: 1.2;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.spec-model--gpu {
  color: rgb(217 119 6);
}

.spec-model--vgpu {
  color: rgb(8 145 178);
}

.spec-model--cpu {
  color: rgb(51 65 85);
}

.spec-chip-group {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 5px;
  width: 100%;
}

.spec-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-height: 22px;
  padding: 0 7px;
  border-radius: 9px;
  border: 1px solid rgb(226 232 240);
  background: rgb(255 255 255 / 0.88);
  white-space: nowrap;
}

.spec-chip--primary.spec-chip--gpu {
  border-color: rgb(252 211 77);
  background: rgb(255 251 235);
}

.spec-chip--primary.spec-chip--vgpu {
  border-color: rgb(153 246 228);
  background: rgb(240 253 250);
}

.spec-chip__label {
  color: rgb(100 116 139);
  font-size: 9px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.spec-chip__value {
  color: rgb(30 41 59);
  font-size: 11px;
  font-weight: 800;
  line-height: 1;
}

.resource-chip-group {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(0, 1fr);
  align-items: center;
  width: 100%;
  gap: 8px;
}

.resource-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-width: 0;
  min-height: 24px;
  padding: 0 8px;
  border: 1px solid rgb(226 232 240);
  border-radius: 10px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(241 245 249));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.92);
  white-space: nowrap;
}

.resource-chip__label {
  color: rgb(148 163 184);
  font-size: 9px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.02em;
  text-transform: none;
  white-space: nowrap;
}

.resource-chip__value {
  color: rgb(51 65 85);
  font-size: 11px;
  font-weight: 700;
  line-height: 1;
}

.inline-metrics {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-wrap: nowrap;
  width: 100%;
  min-width: 0;
}

.inline-metrics--center {
  justify-content: center;
}

.inline-metrics--surface {
  width: auto;
  max-width: 100%;
  justify-content: center;
  gap: 5px;
}

.inventory-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  align-items: stretch;
  width: 100%;
  padding-top: 8px;
  padding-bottom: 8px;
  gap: 8px;
}

.inventory-metrics .inline-metric {
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  white-space: normal;
  text-align: center;
}

.inventory-metrics .inline-metric__label {
  font-size: 10px;
  line-height: 1.2;
  text-transform: none;
}

.inventory-metrics .inline-metric__value {
  flex: 0 0 auto;
  overflow: visible;
  font-size: 12px;
  line-height: 1.2;
}

.inventory-metrics .inline-metrics__divider {
  display: none;
}

.inline-metric {
  display: inline-flex;
  align-items: baseline;
  min-width: 0;
  gap: 3px;
  flex: 0 1 auto;
  white-space: nowrap;
}

.inline-metric__label {
  color: rgb(100 116 139);
  font-size: 11px;
  font-weight: 600;
  line-height: 1.15;
  text-transform: uppercase;
  flex: 0 0 auto;
}

.inline-metric__value {
  color: rgb(51 65 85);
  min-width: 0;
  overflow: hidden;
  font-size: 11px;
  font-weight: 700;
  line-height: 1.15;
  text-overflow: ellipsis;
}

.inline-metric__value--info {
  color: var(--row-button-info-text);
}

.inline-metric__value--success {
  color: var(--row-button-success-text);
}

.inline-metric__value--price {
  color: var(--row-button-danger-text);
}

.inline-metrics__divider {
  color: rgb(203 213 225);
  font-weight: 700;
  line-height: 1.15;
  flex: 0 0 auto;
}

.price-cell .inline-metrics--surface {
  justify-content: stretch;
}

.price-metrics-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  align-items: center;
  width: 100%;
  padding-top: 8px;
  padding-bottom: 8px;
  gap: 8px 12px;
}

.price-metrics-grid .inline-metrics__divider {
  display: none;
}

.price-metric {
  min-width: 0;
  justify-content: space-between;
  gap: 8px;
}

.price-metric .inline-metric__label {
  flex: 0 0 auto;
}

.price-metric .inline-metric__value {
  display: block;
  min-width: 0;
  overflow: hidden;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price-cell .inline-metric__value {
  letter-spacing: -0.01em;
}

.product-actions-cell {
  background: rgb(255 255 255);
  background-clip: padding-box;
}

.product-actions-cell .list-row-actions {
  width: 100%;
  justify-content: center;
}

tbody tr:hover > .product-actions-cell {
  background: var(--list-table-hover-bg);
}

.dark .inline-metric__label {
  color: rgb(148 163 184);
}

.dark .meta-pill {
  border-color: rgb(63 63 70);
  background: rgb(39 39 42);
  color: rgb(212 212 216);
}

.dark .product-actions-cell {
  background: rgb(24 24 27);
}

.dark .spec-surface,
.dark .resource-surface,
.dark .gpu-config-panel,
.dark .inventory-cell .inline-metrics--surface,
.dark .price-cell .inline-metrics--surface {
  border-color: rgb(63 63 70);
  background: linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.03);
}

.dark .inline-metric__value {
  color: rgb(226 232 240);
}

.dark .spec-surface--gpu {
  border-color: rgb(146 64 14);
  background:
    radial-gradient(
      circle at top right,
      rgb(120 53 15 / 0.35),
      transparent 42%
    ),
    linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
}

.dark .spec-surface--vgpu {
  border-color: rgb(14 116 144);
  background:
    radial-gradient(
      circle at top right,
      rgb(8 145 178 / 0.26),
      transparent 42%
    ),
    linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
}

.dark .spec-kind-badge--gpu {
  background: linear-gradient(135deg, rgb(180 83 9), rgb(217 119 6));
  color: rgb(255 237 213);
}

.dark .spec-kind-badge--vgpu {
  background: linear-gradient(135deg, rgb(14 116 144), rgb(13 148 136));
  color: rgb(207 250 254);
}

.dark .spec-kind-badge--cpu {
  background: linear-gradient(135deg, rgb(71 85 105), rgb(100 116 139));
  color: rgb(226 232 240);
}

.dark .gpu-config-title {
  color: rgb(96 165 250);
}

.dark .gpu-config-type-badge {
  background: rgb(8 145 178 / 0.18);
  color: rgb(165 243 252);
}

.dark .gpu-config-metric {
  border-color: rgb(63 63 70);
  background: linear-gradient(180deg, rgb(39 39 42 / 0.96), rgb(24 24 27));
}

.dark .gpu-config-metric__label {
  color: rgb(113 113 122);
}

.dark .gpu-config-metric__value {
  color: rgb(226 232 240);
}

.dark .gpu-config-empty {
  color: rgb(148 163 184);
}

.dark .spec-model {
  color: rgb(226 232 240);
}

.dark .spec-model--gpu {
  color: rgb(251 191 36);
}

.dark .spec-model--vgpu {
  color: rgb(103 232 249);
}

.dark .spec-chip-group--secondary {
  border-top-color: rgb(71 85 105);
}

.dark .spec-chip {
  border-color: rgb(63 63 70);
  background: rgb(24 24 27 / 0.82);
}

.dark .spec-chip--primary.spec-chip--gpu {
  border-color: rgb(146 64 14);
  background: rgb(120 53 15 / 0.3);
}

.dark .spec-chip--primary.spec-chip--vgpu {
  border-color: rgb(14 116 144);
  background: rgb(8 145 178 / 0.22);
}

.dark .spec-chip__label {
  color: rgb(148 163 184);
}

.dark .spec-chip__value {
  color: rgb(226 232 240);
}

.dark .resource-chip {
  border-color: rgb(63 63 70);
  background: linear-gradient(180deg, rgb(39 39 42 / 0.96), rgb(24 24 27));
}

.dark .resource-chip__label {
  color: rgb(148 163 184);
}

.dark .resource-chip__value {
  color: rgb(226 232 240);
}

.dark .inline-metric__value--gpu {
  color: rgb(147 197 253);
}

.dark .inline-metric__value--info {
  color: rgb(191 219 254);
}

.dark .inline-metric__value--success {
  color: rgb(110 231 183);
}

.dark .inline-metric__value--price {
  color: rgb(252 165 165);
}

.dark .inline-metrics__divider {
  color: rgb(71 85 105);
}
</style>

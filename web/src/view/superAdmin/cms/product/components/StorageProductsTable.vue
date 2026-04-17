<template>
  <TableCard :page-size="pageSize">
    <div
      class="console-table-scroll console-table-scroll--fill overflow-x-auto"
    >
      <table
        class="console-table console-table--compact product-table w-full min-w-[1680px]"
        v-loading="loading"
      >
        <colgroup>
          <col style="width: 96px" />
          <col style="width: 180px" />
          <col style="width: 140px" />
          <col style="width: 180px" />
          <col style="width: 120px" />
          <col style="width: 180px" />
          <col style="width: 220px" />
          <col style="width: 120px" />
          <col style="width: 220px" />
          <col style="width: 228px" />
        </colgroup>
        <thead>
          <tr>
            <th>{{ t('id') }}</th>
            <th class="product-column">{{ t('productName') }}</th>
            <th class="cluster-column">{{ t('clusterName') }}</th>
            <th class="storage-column">{{ t('storageClass') }}</th>
            <th>{{ t('area') }}</th>
            <th class="price-column">{{ t('prices') }}</th>
            <th class="text-center inventory-column">{{ t('inventory') }}</th>
            <th class="text-center">{{ t('status') }}</th>
            <th class="desc-column">{{ t('paramDesc') }}</th>
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
            <td class="cluster-cell">
              <span class="meta-pill">{{ row.clusterName || '-' }}</span>
            </td>
            <td class="storage-cell">
              <span class="meta-pill meta-pill--code">{{ row.storageClass }}</span>
            </td>
            <td class="is-secondary">{{ row.area }}</td>
            <td class="price-cell">
              <span class="price-surface">
                ¥{{ row.storagePriceGb?.toFixed(4) || '0' }}/GB/日
              </span>
            </td>
            <td class="text-center inventory-cell">
              <div class="metric-cell-shell">
                <div class="inline-metrics inline-metrics--surface inventory-metrics">
                  <span class="inline-metric">
                    <span class="inline-metric__label">{{ t('usedInstances') }}</span>
                    <span class="inline-metric__value inline-metric__value--info">
                      {{ row.usedCapacity ?? 0 }}
                    </span>
                  </span>
                  <span class="inline-metrics__divider">|</span>
                  <span class="inline-metric">
                    <span class="inline-metric__label">{{ t('maxInstances') }}</span>
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
            <td class="text-center">
              <ListToneBadge
                :label="row.status === 1 ? t('onShelf') : t('offShelf')"
                :tone="row.status === 1 ? 'success' : 'neutral'"
              />
            </td>
            <td class="desc-cell">
              <span class="desc-text" :title="row.description">
                {{ row.description || '-' }}
              </span>
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
              colspan="10"
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
</script>

<style scoped>
.product-column {
  min-width: 180px;
}

.cluster-column {
  min-width: 140px;
}

.storage-column {
  min-width: 180px;
}

.price-column {
  width: 180px;
  min-width: 180px;
}

.desc-column {
  min-width: 220px;
}

.inventory-column {
  min-width: 220px;
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

.price-surface {
  display: inline-flex;
  align-items: center;
  max-width: 100%;
  min-height: 30px;
  padding: 0 10px;
  border: 1px solid rgb(226 232 240);
  border-radius: 14px;
  background: linear-gradient(180deg, rgb(255 255 255), rgb(248 250 252));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.86);
  color: var(--row-button-danger-text);
  font-size: 12px;
  font-weight: 700;
  line-height: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.inventory-cell {
  overflow: hidden;
}

.metric-cell-shell {
  display: flex;
  width: 100%;
  max-width: 100%;
  justify-content: center;
  overflow: hidden;
}

.inventory-cell .inline-metrics--surface {
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

.inline-metrics {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-wrap: nowrap;
  width: 100%;
  min-width: 0;
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

.inline-metrics__divider {
  color: rgb(203 213 225);
  font-weight: 700;
  line-height: 1.15;
  flex: 0 0 auto;
}

.desc-text {
  display: block;
  max-width: 100%;
  overflow: hidden;
  color: rgb(100 116 139);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.dark .meta-pill {
  border-color: rgb(63 63 70);
  background: rgb(39 39 42);
  color: rgb(212 212 216);
}

.dark .product-actions-cell {
  background: rgb(24 24 27);
}

.dark .price-surface {
  border-color: rgb(63 63 70);
  background: linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.03);
  color: rgb(252 165 165);
}

.dark .inventory-cell .inline-metrics--surface {
  border-color: rgb(63 63 70);
  background: linear-gradient(180deg, rgb(39 39 42), rgb(24 24 27));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 0.03);
}

.dark .inline-metric__label {
  color: rgb(148 163 184);
}

.dark .inline-metric__value {
  color: rgb(226 232 240);
}

.dark .desc-text {
  color: rgb(148 163 184);
}
</style>


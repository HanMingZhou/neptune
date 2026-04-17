<template>
  <div class="flex items-center gap-4">
    <RefreshButton :loading="loading" @refresh="$emit('refresh', $event)" />
    <div
      class="transactions-balance-card"
    >
      <div
        class="transactions-balance-card__label"
      >
        {{ t('dashboard.balance') }} (CNY)
      </div>
      <div
        class="transactions-balance-card__value"
      >
        ¥{{
          balance.toLocaleString(undefined, {
            minimumFractionDigits: 2,
            maximumFractionDigits: 6
          })
        }}
      </div>
    </div>
    <el-tooltip
      :disabled="canRecharge || !disabledReason"
      :content="disabledReason"
      placement="bottom"
    >
      <span>
        <button
          :disabled="!canRecharge"
          class="list-toolbar-button list-toolbar-button--primary"
          @click="$emit('recharge')"
        >
          <span class="material-icons">add_circle</span>
          {{ t('dashboard.recharge') }}
        </button>
      </span>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Translator } from '@/types/consoleResource'
import RefreshButton from '@/components/RefreshButton/index.vue'

withDefaults(
  defineProps<{
    balance?: number
    loading?: boolean
    canRecharge?: boolean
    disabledReason?: string
  }>(),
  {
    balance: 0,
    loading: false,
    canRecharge: false,
    disabledReason: ''
  }
)

defineEmits<{
  refresh: [silent: boolean]
  recharge: []
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

<style scoped>
.transactions-balance-card {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  height: var(--toolbar-button-height);
  min-height: var(--toolbar-button-height);
  padding: 0 0.85rem;
  border: 1px solid var(--toolbar-button-secondary-border);
  border-radius: var(--toolbar-button-radius);
  background: var(--toolbar-button-secondary-bg);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.transactions-balance-card__label {
  color: var(--toolbar-button-secondary-text);
  font-size: 0.625rem;
  font-weight: 800;
  line-height: 1;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  white-space: nowrap;
}

.transactions-balance-card__value {
  font-variant-numeric: tabular-nums;
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--color-primary);
  line-height: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .transactions-balance-card {
  border-color: var(--toolbar-button-secondary-border);
  background: var(--toolbar-button-secondary-bg);
}

@media (max-width: 768px) {
  .transactions-balance-card {
    gap: 0.45rem;
    padding: 0 0.65rem;
  }

  .transactions-balance-card__label {
    display: none;
  }
}
</style>

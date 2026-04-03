<template>
  <div class="flex items-center gap-4 h-14">
    <RefreshButton :loading="loading" @refresh="$emit('refresh', $event)" />
    <div class="flex flex-col justify-center text-right px-6 h-full border border-slate-200 dark:border-border-dark rounded-xl bg-white dark:bg-surface-dark shadow-sm">
      <p class="text-[10px] font-black uppercase text-slate-400 leading-none mb-1">{{ t('dashboard.balance') }} (CNY)</p>
      <p class="text-2xl font-black text-primary leading-none">
        ¥{{ balance.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 6 }) }}
      </p>
    </div>
    <el-tooltip :disabled="canRecharge || !disabledReason" :content="disabledReason" placement="bottom">
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

<script setup>
import { inject } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'

defineProps({
  balance: {
    type: Number,
    default: 0
  },
  loading: {
    type: Boolean,
    default: false
  },
  canRecharge: {
    type: Boolean,
    default: false
  },
  disabledReason: {
    type: String,
    default: ''
  }
})

defineEmits(['refresh', 'recharge'])

const t = inject('t', (key) => key)
</script>

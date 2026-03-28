<template>
  <div class="flex items-center gap-4 h-14">
    <RefreshButton class="h-full flex items-center justify-center aspect-square" :loading="loading" @refresh="$emit('refresh', $event)" />
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
          class="px-6 py-3 rounded-xl font-bold text-sm flex items-center gap-2 transition-all"
          :class="canRecharge
            ? 'bg-primary hover:bg-primary-hover text-white shadow-lg shadow-primary/20'
            : 'bg-slate-300 dark:bg-zinc-600 text-white cursor-not-allowed opacity-60'"
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

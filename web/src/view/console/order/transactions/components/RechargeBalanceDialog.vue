<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('order.rechargeDialogTitle')"
    width="560px"
    class="custom-dialog"
    @close="handleClose"
  >
    <div class="space-y-5 pt-2">
      <section class="space-y-3">
        <div class="text-sm font-bold text-slate-700 dark:text-slate-200">{{ t('order.rechargeMethod') }}</div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <button
            type="button"
            class="rounded-2xl border px-4 py-4 text-left transition-all"
            :class="form.method === 4
              ? 'border-primary bg-primary/5 shadow-sm shadow-primary/10'
              : 'border-slate-200 dark:border-border-dark bg-white dark:bg-surface-dark hover:border-primary/40'"
            @click="form.method = 4"
          >
            <div class="flex items-center justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">{{ t('order.platformRecharge') }}</div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ t('order.platformRechargeDesc') }}</div>
              </div>
              <span v-if="form.method === 4" class="material-icons text-primary text-lg">check_circle</span>
            </div>
          </button>

          <button
            type="button"
            disabled
            class="rounded-2xl border border-slate-200 dark:border-border-dark bg-slate-50 dark:bg-zinc-900/40 px-4 py-4 text-left opacity-60 cursor-not-allowed"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">{{ t('order.alipayRecharge') }}</div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ t('order.personalRechargeComingSoon') }}</div>
              </div>
              <span class="text-[10px] font-black uppercase tracking-wide text-slate-400">{{ t('order.comingSoon') }}</span>
            </div>
          </button>

          <button
            type="button"
            disabled
            class="rounded-2xl border border-slate-200 dark:border-border-dark bg-slate-50 dark:bg-zinc-900/40 px-4 py-4 text-left opacity-60 cursor-not-allowed"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">{{ t('order.wechatRecharge') }}</div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ t('order.personalRechargeComingSoon') }}</div>
              </div>
              <span class="text-[10px] font-black uppercase tracking-wide text-slate-400">{{ t('order.comingSoon') }}</span>
            </div>
          </button>
        </div>
      </section>

      <section class="rounded-2xl border border-slate-200 dark:border-border-dark bg-slate-50/80 dark:bg-zinc-900/40 p-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <div class="text-sm font-bold text-slate-900 dark:text-white">{{ t('order.rechargeAmount') }}</div>
            <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ t('order.rechargeAmountHint') }}</div>
          </div>
          <div class="text-right">
            <div class="text-[11px] font-bold uppercase tracking-wide text-slate-400">{{ t('dashboard.balance') }}</div>
            <div class="text-sm font-mono font-bold text-slate-700 dark:text-slate-200">{{ t('order.instantArrival') }}</div>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-2">
          <button
            v-for="preset in amountPresets"
            :key="preset"
            type="button"
            class="rounded border px-3 py-1.5 text-xs font-bold transition-all"
            :class="Number(form.amount) === preset
              ? 'border-primary bg-primary text-white shadow-sm shadow-primary/20'
              : 'border-slate-200 dark:border-border-dark bg-white dark:bg-surface-dark text-slate-600 dark:text-slate-300 hover:border-primary/40'"
            @click="form.amount = preset"
          >
            ¥{{ preset }}
          </button>
        </div>
      </section>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('order.rechargeAmount')" prop="amount">
          <div class="recharge-amount-shell">
            <span class="recharge-amount-prefix">¥</span>
            <el-input-number
              v-model="form.amount"
              :min="0.01"
              :precision="2"
              :step="100"
              :controls="false"
              class="recharge-amount-input"
            />
          </div>
          <div class="mt-2 text-xs text-slate-500 dark:text-slate-400">{{ t('order.rechargeAmountFootnote') }}</div>
        </el-form-item>
        <el-form-item :label="t('remark')" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            :placeholder="t('order.rechargeRemarkPlaceholder')"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3 pt-4 border-t border-slate-100 dark:border-border-dark">
        <button
          class="list-toolbar-button list-toolbar-button--secondary"
          @click="dialogVisible = false"
        >
          {{ t('general.cancel') }}
        </button>
        <button
          :disabled="submitting"
          class="list-toolbar-button list-toolbar-button--primary"
          @click="handleSubmit"
        >
          <span v-if="submitting" class="material-icons animate-spin text-sm">autorenew</span>
          {{ t('order.confirmRecharge') }}
        </button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

const amountPresets = [100, 500, 1000, 5000]

const props = defineProps({
  form: {
    type: Object,
    required: true
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  rules: {
    type: Object,
    default: () => ({})
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['closed', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)
const formRef = ref(null)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const handleClose = () => {
  formRef.value?.clearValidate()
  emit('closed')
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('submit')
  } catch (error) {
    // validation error
  }
}
</script>

<style scoped>
.recharge-amount-shell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  min-height: 70px;
  padding: 0.95rem 1.1rem;
  border: 1px solid rgb(226 232 240);
  border-radius: 1rem;
  background: rgba(248, 250, 252, 0.9);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.dark .recharge-amount-shell {
  border-color: rgb(63 63 70);
  background: rgba(24, 24, 27, 0.85);
}

.recharge-amount-prefix {
  font-size: 1.75rem;
  font-weight: 800;
  color: rgb(15 23 42);
}

.dark .recharge-amount-prefix {
  color: rgb(248 250 252);
}

:deep(.recharge-amount-input) {
  width: 100% !important;
}

:deep(.recharge-amount-input .el-input__wrapper) {
  min-height: auto !important;
  padding: 0 !important;
  border: none !important;
  border-radius: 0 !important;
  background: transparent !important;
  box-shadow: none !important;
}

:deep(.recharge-amount-input .el-input__inner) {
  font-size: 1.75rem !important;
  font-weight: 800 !important;
  color: rgb(15 23 42) !important;
}

.dark :deep(.recharge-amount-input .el-input__inner) {
  color: rgb(248 250 252) !important;
}
</style>

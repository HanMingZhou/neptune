<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('general.cancel')"
    :model="form"
    :rules="rules"
    :submit-text="t('order.confirmRecharge')"
    :submitting="submitting"
    :title="t('order.rechargeDialogTitle')"
    class="recharge-balance-dialog"
    width="560px"
    @closed="emit('closed')"
    @submit="emit('submit')"
  >
    <div class="space-y-5 pt-2">
      <section class="space-y-3">
        <div class="text-sm font-bold text-slate-700 dark:text-slate-200">
          {{ t('order.rechargeMethod') }}
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <el-button
            class="recharge-method-card"
            :class="
              form.method === 4
                ? 'recharge-method-card--active'
                : 'recharge-method-card--idle'
            "
            @click="form.method = 4"
          >
            <div class="flex items-center justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">
                  {{ t('order.platformRecharge') }}
                </div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                  {{ t('order.platformRechargeDesc') }}
                </div>
              </div>
                <span
                  v-if="form.method === 4"
                  class="material-icons text-primary text-lg"
                  >check_circle</span
                >
              </div>
          </el-button>

          <el-button
            disabled
            class="recharge-method-card recharge-method-card--disabled"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">
                  {{ t('order.alipayRecharge') }}
                </div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                  {{ t('order.personalRechargeComingSoon') }}
                </div>
              </div>
              <span
                class="text-[10px] font-black uppercase tracking-wide text-slate-400"
                >{{ t('order.comingSoon') }}</span
              >
            </div>
          </el-button>

          <el-button
            disabled
            class="recharge-method-card recharge-method-card--disabled"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-white">
                  {{ t('order.wechatRecharge') }}
                </div>
                <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                  {{ t('order.personalRechargeComingSoon') }}
                </div>
              </div>
              <span
                class="text-[10px] font-black uppercase tracking-wide text-slate-400"
                >{{ t('order.comingSoon') }}</span
              >
            </div>
          </el-button>
        </div>
      </section>

      <section
        class="rounded-2xl border border-slate-200 dark:border-border-dark bg-slate-50/80 dark:bg-zinc-900/40 p-4"
      >
        <div class="flex items-center justify-between gap-3">
          <div>
            <div class="text-sm font-bold text-slate-900 dark:text-white">
              {{ t('order.rechargeAmount') }}
            </div>
            <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">
              {{ t('order.rechargeAmountHint') }}
            </div>
          </div>
          <div class="text-right">
            <div
              class="text-[11px] font-bold uppercase tracking-wide text-slate-400"
            >
              {{ t('dashboard.balance') }}
            </div>
            <div
              class="text-sm font-mono font-bold text-slate-700 dark:text-slate-200"
            >
              {{ t('order.instantArrival') }}
            </div>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-2">
          <el-button
            v-for="preset in amountPresets"
            :key="preset"
            class="recharge-preset-button"
            :class="
              Number(form.amount) === preset
                ? 'recharge-preset-button--active'
                : 'recharge-preset-button--idle'
            "
            @click="form.amount = preset"
          >
            ¥{{ preset }}
          </el-button>
        </div>
      </section>

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
        <div class="mt-2 text-xs text-slate-500 dark:text-slate-400">
          {{ t('order.rechargeAmountFootnote') }}
        </div>
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
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { RechargeForm } from '@/types/order'

const amountPresets = [100, 500, 1000, 5000]

const props = defineProps<{
  form: RechargeForm
  modelValue: boolean
  rules: FormRules<RechargeForm>
  submitting: boolean
}>()

const emit = defineEmits<{
  closed: []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

<style scoped>
.recharge-method-card {
  width: 100%;
  height: auto;
  margin: 0;
  justify-content: flex-start;
  padding: 1rem;
  border-radius: 1rem;
  white-space: normal;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.recharge-method-card :deep(.el-button__text) {
  display: block;
  width: 100%;
}

.recharge-method-card--active {
  border-color: var(--el-color-primary);
  background: rgb(37 99 235 / 0.08);
  box-shadow: 0 6px 18px rgb(37 99 235 / 0.1);
}

.recharge-method-card--idle {
  border-color: rgb(226 232 240);
  background: rgb(255 255 255);
}

.recharge-method-card--idle:hover {
  border-color: rgb(37 99 235 / 0.4);
}

.recharge-method-card--disabled {
  border-color: rgb(226 232 240);
  background: rgb(248 250 252);
  opacity: 0.6;
}

.dark .recharge-method-card--active {
  background: rgb(37 99 235 / 0.12);
}

.dark .recharge-method-card--idle,
.dark .recharge-method-card--disabled {
  border-color: rgb(63 63 70);
}

.dark .recharge-method-card--idle {
  background: rgb(24 24 27);
}

.dark .recharge-method-card--disabled {
  background: rgb(24 24 27 / 0.4);
}

.recharge-preset-button {
  height: auto;
  margin: 0;
  padding: 0.375rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  font-weight: 700;
}

.recharge-preset-button--active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary);
  color: white;
  box-shadow: 0 6px 16px rgb(37 99 235 / 0.2);
}

.recharge-preset-button--idle {
  border-color: rgb(226 232 240);
  background: rgb(255 255 255);
  color: rgb(71 85 105);
}

.recharge-preset-button--idle:hover {
  border-color: rgb(37 99 235 / 0.4);
}

.dark .recharge-preset-button--idle {
  border-color: rgb(63 63 70);
  background: rgb(24 24 27);
  color: rgb(203 213 225);
}

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
</style>

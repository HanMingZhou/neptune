<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('general.cancel')"
    form-class="recharge-balance-dialog__form"
    label-width="108px"
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
    <div class="recharge-balance-dialog__content">
      <section class="recharge-balance-dialog__section">
        <div class="recharge-section-title text-sm font-bold">
          {{ t('order.rechargeMethod') }}
        </div>
        <div class="recharge-method-grid">
          <button
            type="button"
            class="recharge-method-card"
            :class="
              form.method === 4
                ? 'recharge-method-card--active'
                : 'recharge-method-card--idle'
            "
            @click="form.method = 4"
          >
            <div class="recharge-method-card__inner">
              <div class="recharge-method-card__meta">
                <div class="recharge-section-title text-sm font-bold">
                  {{ t('order.platformRecharge') }}
                </div>
                <div class="recharge-section-desc text-xs">
                  {{ t('order.platformRechargeDesc') }}
                </div>
              </div>
              <span
                v-if="form.method === 4"
                class="recharge-method-card__check material-icons"
                >check_circle</span
              >
            </div>
          </button>

          <button
            type="button"
            disabled
            class="recharge-method-card recharge-method-card--disabled"
          >
            <div class="recharge-method-card__inner recharge-method-card__inner--top">
              <div class="recharge-method-card__meta">
                <div class="recharge-section-title text-sm font-bold">
                  {{ t('order.alipayRecharge') }}
                </div>
                <div class="recharge-section-desc text-xs">
                  {{ t('order.personalRechargeComingSoon') }}
                </div>
              </div>
              <span class="recharge-method-card__badge">{{
                t('order.comingSoon')
              }}</span>
            </div>
          </button>

          <button
            type="button"
            disabled
            class="recharge-method-card recharge-method-card--disabled"
          >
            <div class="recharge-method-card__inner recharge-method-card__inner--top">
              <div class="recharge-method-card__meta">
                <div class="recharge-section-title text-sm font-bold">
                  {{ t('order.wechatRecharge') }}
                </div>
                <div class="recharge-section-desc text-xs">
                  {{ t('order.personalRechargeComingSoon') }}
                </div>
              </div>
              <span class="recharge-method-card__badge">{{
                t('order.comingSoon')
              }}</span>
            </div>
          </button>
        </div>
      </section>

      <section class="recharge-amount-panel">
        <div class="recharge-amount-panel__head">
          <div>
            <div class="recharge-section-title text-sm font-bold">
              {{ t('order.rechargeAmount') }}
            </div>
            <div class="recharge-section-desc text-xs">
              {{ t('order.rechargeAmountHint') }}
            </div>
          </div>
          <div class="recharge-amount-panel__meta">
            <div class="recharge-amount-panel__meta-label">
              {{ t('dashboard.balance') }}
            </div>
            <div class="recharge-amount-panel__meta-value">
              {{ t('order.instantArrival') }}
            </div>
          </div>
        </div>

        <div class="recharge-preset-grid">
          <button
            type="button"
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
          </button>
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
        <div class="recharge-amount-footnote">
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

<style>
.recharge-balance-dialog .el-dialog__body {
  padding: 20px 24px 14px;
}

.recharge-balance-dialog__form {
  display: grid;
  gap: 1.1rem;
}

.recharge-balance-dialog .recharge-balance-dialog__content {
  display: grid;
  gap: 1.25rem;
  padding-top: 0.35rem;
}

.recharge-balance-dialog .recharge-balance-dialog__section {
  display: grid;
  gap: 0.75rem;
}

.recharge-balance-dialog .recharge-section-title {
  line-height: 1.3;
  color: rgb(30 41 59);
}

.recharge-balance-dialog .recharge-section-desc {
  line-height: 1.35;
  margin-top: 0.25rem;
  color: rgb(100 116 139);
}

.recharge-balance-dialog .recharge-method-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 0.75rem;
}

@media (min-width: 768px) {
  .recharge-balance-dialog .recharge-method-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

.recharge-balance-dialog .recharge-method-card {
  display: block;
  width: 100%;
  height: auto !important;
  min-height: 88px !important;
  margin: 0;
  justify-content: flex-start !important;
  align-items: stretch !important;
  padding: 1rem !important;
  border: 1px solid rgb(226 232 240);
  border-radius: 1rem;
  white-space: normal;
  text-align: left;
  font: inherit;
  line-height: 1.25 !important;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
}

.recharge-balance-dialog .recharge-method-card__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.recharge-balance-dialog .recharge-method-card__inner--top {
  align-items: flex-start;
}

.recharge-balance-dialog .recharge-method-card__meta {
  min-width: 0;
}

.recharge-balance-dialog .recharge-method-card__check {
  font-size: 1.125rem;
  color: var(--el-color-primary);
}

.recharge-balance-dialog .recharge-method-card__badge {
  font-size: 0.625rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: rgb(148 163 184);
}

.recharge-balance-dialog .recharge-method-card--active {
  border-color: var(--el-color-primary);
  background: rgb(37 99 235 / 0.08);
  box-shadow: 0 6px 18px rgb(37 99 235 / 0.1);
}

.recharge-balance-dialog .recharge-method-card--idle {
  border-color: rgb(226 232 240);
  background: rgb(255 255 255);
}

.recharge-balance-dialog .recharge-method-card:disabled {
  cursor: not-allowed;
}

.recharge-balance-dialog .recharge-method-card--idle:hover {
  border-color: rgb(37 99 235 / 0.4);
}

.recharge-balance-dialog .recharge-method-card:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.recharge-balance-dialog .recharge-method-card--disabled {
  border-color: rgb(226 232 240);
  background: rgb(248 250 252);
  opacity: 0.6;
}

.recharge-balance-dialog .recharge-amount-panel {
  border: 1px solid rgb(226 232 240);
  border-radius: 1rem;
  padding: 1rem;
  background: rgb(248 250 252 / 0.8);
}

.recharge-balance-dialog .recharge-amount-panel__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.recharge-balance-dialog .recharge-amount-panel__meta {
  text-align: right;
}

.recharge-balance-dialog .recharge-amount-panel__meta-label {
  font-size: 0.6875rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: rgb(148 163 184);
}

.recharge-balance-dialog .recharge-amount-panel__meta-value {
  margin-top: 0.125rem;
  font-size: 0.875rem;
  font-weight: 700;
  color: rgb(51 65 85);
}

.recharge-balance-dialog .recharge-preset-grid {
  margin-top: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.recharge-balance-dialog .recharge-preset-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: auto !important;
  min-height: 30px !important;
  margin: 0;
  padding: 0.375rem 0.75rem !important;
  border: 1px solid rgb(226 232 240);
  border-radius: 0.5rem;
  font-size: 0.75rem;
  font-weight: 700;
  text-align: center;
  font: inherit;
  line-height: 1.2 !important;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease;
}

.recharge-balance-dialog .recharge-preset-button--active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary);
  color: white;
  box-shadow: 0 6px 16px rgb(37 99 235 / 0.2);
}

.recharge-balance-dialog .recharge-preset-button--idle {
  border-color: rgb(226 232 240);
  background: rgb(255 255 255);
  color: rgb(71 85 105);
}

.recharge-balance-dialog .recharge-preset-button--idle:hover {
  border-color: rgb(37 99 235 / 0.4);
}

.recharge-balance-dialog .recharge-preset-button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.recharge-balance-dialog .recharge-amount-shell {
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

.recharge-balance-dialog .recharge-amount-prefix {
  font-size: 1.75rem;
  font-weight: 800;
  color: rgb(15 23 42);
}

.recharge-balance-dialog .recharge-amount-input {
  width: 100% !important;
}

.recharge-balance-dialog .recharge-amount-input .el-input__wrapper {
  min-height: auto !important;
  padding: 0 !important;
  border: none !important;
  border-radius: 0 !important;
  background: transparent !important;
  box-shadow: none !important;
}

.recharge-balance-dialog .recharge-amount-input .el-input__inner {
  font-size: 1.75rem !important;
  font-weight: 800 !important;
  color: rgb(15 23 42) !important;
}

.recharge-balance-dialog .recharge-amount-footnote {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: rgb(100 116 139);
}

.dark .recharge-balance-dialog .recharge-section-title {
  color: rgb(248 250 252);
}

.dark .recharge-balance-dialog .recharge-section-desc {
  color: rgb(148 163 184);
}

.dark .recharge-balance-dialog .recharge-method-card--active {
  background: rgb(37 99 235 / 0.12);
}

.dark .recharge-balance-dialog .recharge-method-card--idle,
.dark .recharge-balance-dialog .recharge-method-card--disabled {
  border-color: rgb(63 63 70);
}

.dark .recharge-balance-dialog .recharge-method-card--idle {
  background: rgb(24 24 27);
}

.dark .recharge-balance-dialog .recharge-method-card--disabled {
  background: rgb(24 24 27 / 0.4);
}

.dark .recharge-balance-dialog .recharge-preset-button--idle {
  border-color: rgb(63 63 70);
  background: rgb(24 24 27);
  color: rgb(203 213 225);
}

.dark .recharge-balance-dialog .recharge-amount-panel {
  border-color: rgb(63 63 70);
  background: rgba(24, 24, 27, 0.45);
}

.dark .recharge-balance-dialog .recharge-amount-shell {
  border-color: rgb(63 63 70);
  background: rgba(24, 24, 27, 0.85);
}

.dark .recharge-balance-dialog .recharge-amount-panel__meta-value {
  color: rgb(203 213 225);
}

.dark .recharge-balance-dialog .recharge-amount-prefix,
.dark .recharge-balance-dialog .recharge-amount-input .el-input__inner {
  color: rgb(248 250 252);
}

.dark .recharge-balance-dialog .recharge-amount-footnote {
  color: rgb(148 163 184);
}
</style>

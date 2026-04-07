<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('general.cancel')"
    :model="form"
    :rules="rules"
    :submit-text="t('general.confirm')"
    :submitting="submitting"
    :title="t('order.applyInvoice')"
    class="invoice-apply-dialog"
    label-width="120px"
    width="500px"
    @closed="emit('closed')"
    @submit="emit('submit')"
  >
    <div class="mt-4">
      <el-form-item :label="t('order.amount')" prop="amount">
        <el-input-number
          v-model="form.amount"
          :min="0.01"
          :precision="2"
          :step="100"
          class="w-full"
        />
      </el-form-item>
      <el-form-item :label="`${t('order.invoiceTitle')} ID`" prop="titleId">
        <el-input-number
          v-model="form.titleId"
          :min="1"
          :step="1"
          class="w-full"
        />
        <div class="text-[10px] text-slate-400 mt-1 leading-tight">
          {{ t('order.comingSoon') }}
        </div>
      </el-form-item>
      <el-form-item :label="`${t('order.invoiceAddress')} ID`" prop="addressId">
        <el-input-number
          v-model="form.addressId"
          :min="1"
          :step="1"
          class="w-full"
        />
        <div class="text-[10px] text-slate-400 mt-1 leading-tight">
          {{ t('order.comingSoon') }}
        </div>
      </el-form-item>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { ApplyInvoiceForm } from '@/types/order'

const props = withDefaults(
  defineProps<{
    form: ApplyInvoiceForm
    modelValue?: boolean
    rules?: FormRules<ApplyInvoiceForm>
    submitting?: boolean
  }>(),
  {
    modelValue: false,
    rules: () => ({}),
    submitting: false
  }
)

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

<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('order.applyInvoice')"
    width="500px"
    class="custom-dialog"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" class="mt-4">
      <el-form-item :label="t('order.amount')" prop="amount">
        <el-input-number v-model="form.amount" :min="0.01" :precision="2" :step="100" class="w-full" />
      </el-form-item>
      <el-form-item :label="`${t('order.invoiceTitle')} ID`" prop="titleId">
        <el-input-number v-model="form.titleId" :min="1" :step="1" class="w-full" />
        <div class="text-[10px] text-slate-400 mt-1 leading-tight">{{ t('order.comingSoon') }}</div>
      </el-form-item>
      <el-form-item :label="`${t('order.invoiceAddress')} ID`" prop="addressId">
        <el-input-number v-model="form.addressId" :min="1" :step="1" class="w-full" />
        <div class="text-[10px] text-slate-400 mt-1 leading-tight">{{ t('order.comingSoon') }}</div>
      </el-form-item>
    </el-form>
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
          {{ t('general.confirm') }}
        </button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

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

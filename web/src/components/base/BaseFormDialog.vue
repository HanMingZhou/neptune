<template>
  <BaseDialog
    :model-value="modelValue"
    :title="title"
    :width="width"
    :top="top"
    :align-center="alignCenter"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    :destroy-on-close="destroyOnClose"
    :shell="shell"
    :shell-size="shellSize"
    :show-close="showClose"
    @close="emit('close')"
    @closed="emit('closed')"
    @opened="emit('opened')"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template v-if="$slots.header" #header="{ requestClose }">
      <slot name="header" :request-close="requestClose" />
    </template>

    <el-form
      v-if="modelValue || keepMounted"
      ref="formRef"
      :model="model"
      :rules="rules"
      :label-position="labelPosition"
      :label-width="labelWidth"
      :class="['overlay-form-shell', 'dialog-form-shell', formClass]"
    >
      <slot :form-ref="formRef" :submit-form="submitForm" />
    </el-form>

    <slot name="append" />

    <template #footer="{ requestClose }">
      <slot
        v-if="$slots.footer"
        name="footer"
        :form-ref="formRef"
        :request-close="requestClose"
        :submit-form="submitForm"
      />
      <div v-else class="flex justify-end gap-3">
        <el-button @click="requestClose">{{ cancelText }}</el-button>
        <el-button
          type="primary"
          :disabled="submitDisabled"
          :loading="submitting"
          @click="submitForm"
        >
          {{ submitText }}
        </el-button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import BaseDialog from './BaseDialog.vue'

type DialogShellSize = 'md' | 'lg' | 'xl'
type LabelPosition = 'left' | 'right' | 'top'

const props = withDefaults(
  defineProps<{
    alignCenter?: boolean
    cancelText?: string
    closeOnClickModal?: boolean
    closeOnPressEscape?: boolean
    destroyOnClose?: boolean
    formClass?: string
    keepMounted?: boolean
    labelPosition?: LabelPosition
    labelWidth?: number | string
    model: object
    modelValue?: boolean
    rules?: FormRules
    shell?: boolean
    shellSize?: DialogShellSize
    showClose?: boolean
    submitDisabled?: boolean
    submitText?: string
    submitting?: boolean
    title?: string
    top?: string
    width?: string
  }>(),
  {
    alignCenter: false,
    cancelText: 'Cancel',
    closeOnClickModal: true,
    closeOnPressEscape: true,
    destroyOnClose: false,
    formClass: '',
    keepMounted: false,
    labelPosition: 'left',
    labelWidth: 'auto',
    modelValue: false,
    rules: () => ({}),
    shell: true,
    shellSize: undefined,
    showClose: true,
    submitDisabled: false,
    submitText: 'Submit',
    submitting: false,
    title: '',
    top: undefined,
    width: undefined
  }
)

const emit = defineEmits<{
  close: []
  closed: []
  opened: []
  submit: []
  'update:modelValue': [value: boolean]
}>()

const formRef = ref<FormInstance | null>(null)

const submitForm = async (): Promise<void> => {
  if (!formRef.value) {
    emit('submit')
    return
  }

  try {
    await formRef.value.validate()
    emit('submit')
  } catch {
    // Validation feedback is already shown by Element Plus.
  }
}

watch(
  () => props.modelValue,
  async (value: boolean | undefined) => {
    if (!value) {
      return
    }

    await nextTick()
    formRef.value?.clearValidate()
  }
)

defineExpose({
  formRef,
  submitForm
})
</script>

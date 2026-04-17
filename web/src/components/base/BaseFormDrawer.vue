<template>
  <BaseDrawer
    :model-value="modelValue"
    :title="title"
    :size="size"
    :direction="direction"
    :show-close="showClose"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    :destroy-on-close="destroyOnClose"
    @close="emit('close')"
    @closed="emit('closed')"
    @opened="emit('opened')"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #header="{ requestClose }">
      <slot
        v-if="$slots.header"
        name="header"
        :form-ref="formRef"
        :request-close="requestClose"
        :submit-form="submitForm"
      />
      <div v-else class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ title }}</span>
        <div class="flex gap-3">
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
      </div>
    </template>

    <slot name="prepend" />

    <el-form
      v-show="modelValue || keepMounted"
      ref="formRef"
      :model="model"
      :rules="rules"
      :label-position="labelPosition"
      :label-width="labelWidth"
      :inline="inline"
      :class="['overlay-form-shell', 'drawer-form-shell', formClass]"
    >
      <slot :form-ref="formRef" :submit-form="submitForm" />
    </el-form>

    <slot name="append" />

    <template v-if="$slots.footer" #footer="{ requestClose }">
      <slot
        name="footer"
        :form-ref="formRef"
        :request-close="requestClose"
        :submit-form="submitForm"
      />
    </template>
  </BaseDrawer>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import BaseDrawer from './BaseDrawer.vue'

type LabelPosition = 'left' | 'right' | 'top'
type DrawerDirection = 'rtl' | 'ltr' | 'ttb' | 'btt'

const props = withDefaults(
  defineProps<{
    cancelText?: string
    closeOnClickModal?: boolean
    closeOnPressEscape?: boolean
    destroyOnClose?: boolean
    direction?: DrawerDirection
    formClass?: string
    inline?: boolean
    keepMounted?: boolean
    labelPosition?: LabelPosition
    labelWidth?: number | string
    model: object
    modelValue?: boolean
    rules?: FormRules
    showClose?: boolean
    size?: string | number
    submitDisabled?: boolean
    submitText?: string
    submitting?: boolean
    title?: string
  }>(),
  {
    cancelText: 'Cancel',
    closeOnClickModal: true,
    closeOnPressEscape: true,
    destroyOnClose: false,
    direction: 'rtl',
    formClass: '',
    inline: false,
    keepMounted: false,
    labelPosition: 'left',
    labelWidth: 'auto',
    modelValue: false,
    rules: () => ({}),
    showClose: false,
    size: 480,
    submitDisabled: false,
    submitText: 'Submit',
    submitting: false,
    title: ''
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

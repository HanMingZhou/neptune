<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    :class="dialogClass"
    :model="form"
    :rules="rules"
    :submit-text="t('confirm')"
    :submitting="loading"
    :title="title"
    form-class="py-4"
    label-position="top"
    width="400px"
    @closed="emit('closed')"
    @submit="emit('submit')"
  >
    <el-form-item :label="phoneLabel" prop="phone">
      <el-input v-model="form.phone" :placeholder="phonePlaceholder">
        <template #prefix>
          <el-icon><Phone /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item :label="codeLabel" prop="code">
      <div class="overlay-code-row">
        <el-input
          v-model="form.code"
          :placeholder="codePlaceholder"
          class="overlay-code-row__input"
        >
          <template #prefix>
            <el-icon><Key /></el-icon>
          </template>
        </el-input>
        <el-button
          type="primary"
          :disabled="disableRequestCode"
          class="overlay-code-row__action"
          @click="emit('request-code')"
        >
          {{ requestCodeText }}
        </el-button>
      </div>
    </el-form-item>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { Key, Phone } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { AccountPhoneForm } from '@/types/account'

const props = withDefaults(
  defineProps<{
    codeLabel?: string
    codePlaceholder?: string
    disableRequestCode?: boolean
    dialogClass?: string
    form: AccountPhoneForm
    loading?: boolean
    modelValue?: boolean
    phoneLabel?: string
    phonePlaceholder?: string
    requestCodeText?: string
    rules?: FormRules<AccountPhoneForm>
    title?: string
  }>(),
  {
    codeLabel: '验证码',
    codePlaceholder: '请输入验证码',
    disableRequestCode: false,
    dialogClass: '',
    loading: false,
    modelValue: false,
    phoneLabel: '新手机号',
    phonePlaceholder: '请输入新手机号',
    requestCodeText: '获取验证码',
    rules: () => ({}),
    title: ''
  }
)

const emit = defineEmits<{
  closed: []
  'request-code': []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

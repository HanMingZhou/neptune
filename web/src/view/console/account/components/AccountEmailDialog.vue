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
    <el-form-item :label="emailLabel" prop="email">
      <el-input v-model="form.email" :placeholder="emailPlaceholder">
        <template #prefix>
          <el-icon><Message /></el-icon>
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
import { Key, Message } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { AccountEmailForm } from '@/types/account'
import type { Translator } from '@/types/consoleResource'

const props = withDefaults(
  defineProps<{
    codeLabel?: string
    codePlaceholder?: string
    disableRequestCode?: boolean
    dialogClass?: string
    emailLabel?: string
    emailPlaceholder?: string
    form: AccountEmailForm
    loading?: boolean
    modelValue?: boolean
    requestCodeText?: string
    rules?: FormRules<AccountEmailForm>
    title?: string
  }>(),
  {
    codeLabel: '验证码',
    codePlaceholder: '请输入验证码',
    disableRequestCode: false,
    dialogClass: '',
    emailLabel: '新邮箱',
    emailPlaceholder: '请输入新邮箱',
    loading: false,
    modelValue: false,
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

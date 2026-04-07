<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    class="account-security-dialog"
    :model="formModel"
    :rules="rules"
    title="关闭 MFA 验证"
    width="400px"
    form-class="py-4"
    label-position="top"
    @closed="emit('closed')"
    @submit="emit('confirm')"
  >
    <div class="space-y-4">
      <div
        class="rounded-xl border border-amber-500/20 bg-amber-500/10 p-3 text-sm text-amber-600"
      >
        关闭 MFA 后，登录时将不再需要验证码。请输入当前的 MFA 验证码以确认关闭。
      </div>

      <el-form-item label="验证码" prop="code">
        <el-input
          v-model="formModel.code"
          placeholder="请输入6位验证码"
          maxlength="6"
        >
          <template #prefix>
            <el-icon><Key /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </div>

    <template #footer="{ requestClose, submitForm }">
      <div class="dialog-footer">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button
          type="danger"
          :loading="loading"
          :disabled="formModel.code.length !== 6"
          @click="submitForm"
        >
          确认关闭
        </el-button>
      </div>
    </template>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject, reactive, watch } from 'vue'
import { Key } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'

interface MfaDisableForm {
  code: string
}

const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    loading?: boolean
    code?: string
  }>(),
  {
    modelValue: false,
    loading: false,
    code: ''
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'update:code': [value: string]
  confirm: []
  closed: []
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const formModel = reactive<MfaDisableForm>({
  code: ''
})

const rules: FormRules<MfaDisableForm> = {
  code: [
    { required: true, message: '请输入6位验证码', trigger: 'blur' },
    { min: 6, max: 6, message: '请输入6位验证码', trigger: 'blur' }
  ]
}

watch(
  () => props.code,
  (value: string) => {
    if (value !== formModel.code) {
      formModel.code = value
    }
  },
  { immediate: true }
)

watch(
  () => formModel.code,
  (value: string) => {
    if (value !== props.code) {
      emit('update:code', value)
    }
  }
)
</script>

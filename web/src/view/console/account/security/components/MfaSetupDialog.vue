<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    class="account-security-dialog"
    :model="formModel"
    :rules="rules"
    title="启用虚拟 MFA 设备"
    width="460px"
    form-class="py-4"
    label-position="top"
    @closed="emit('closed')"
    @submit="emit('confirm')"
  >
    <div class="space-y-5">
      <div class="text-sm leading-relaxed text-slate-600 dark:text-slate-400">
        <p class="mb-2">
          请使用 <strong>Google Authenticator</strong>、<strong
            >Microsoft Authenticator</strong
          >
          或其他支持 TOTP 的应用扫描下方二维码：
        </p>
      </div>

      <div class="flex flex-col items-center gap-4">
        <div class="rounded-xl border border-slate-200 bg-white p-3 shadow-sm">
          <img
            v-if="qr"
            :src="qr"
            alt="MFA QR Code"
            class="h-[180px] w-[180px]"
          />
          <div
            v-else
            class="flex h-[180px] w-[180px] items-center justify-center text-sm text-slate-400"
          >
            加载中...
          </div>
        </div>

        <div v-if="secret" class="w-full">
          <p class="mb-1 text-xs text-slate-400">无法扫码？手动输入密钥：</p>
          <div
            class="flex items-center gap-2 rounded-lg border border-border-light bg-slate-50 p-3 dark:border-border-dark dark:bg-zinc-800"
          >
            <code class="flex-1 select-all break-all font-mono text-xs">{{
              secret
            }}</code>
            <el-button
              link
              type="primary"
              class="shrink-0 !h-auto !min-h-0 !p-0 text-xs"
              @click="emit('copy-secret')"
            >
              复制
            </el-button>
          </div>
        </div>
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
          type="primary"
          :loading="loading"
          :disabled="formModel.code.length !== 6"
          @click="submitForm"
        >
          确认绑定
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

interface MfaSetupForm {
  code: string
}

const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    loading?: boolean
    qr?: string
    secret?: string
    code?: string
  }>(),
  {
    modelValue: false,
    loading: false,
    qr: '',
    secret: '',
    code: ''
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'update:code': [value: string]
  confirm: []
  'copy-secret': []
  closed: []
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const formModel = reactive<MfaSetupForm>({
  code: ''
})

const rules: FormRules<MfaSetupForm> = {
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

<template>
  <el-dialog
    v-model="dialogVisible"
    title="启用虚拟 MFA 设备"
    width="460px"
    class="custom-dialog"
    @close="handleClose"
  >
    <div class="py-4 space-y-5">
      <div class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
        <p class="mb-2">
          请使用 <strong>Google Authenticator</strong>、<strong>Microsoft Authenticator</strong>
          或其他支持 TOTP 的应用扫描下方二维码：
        </p>
      </div>
      <div class="flex flex-col items-center gap-4">
        <div class="p-3 bg-white rounded-xl border border-slate-200 shadow-sm">
          <img v-if="qr" :src="qr" alt="MFA QR Code" class="w-[180px] h-[180px]" />
          <div v-else class="w-[180px] h-[180px] flex items-center justify-center text-slate-400 text-sm">
            加载中...
          </div>
        </div>
        <div v-if="secret" class="w-full">
          <p class="text-xs text-slate-400 mb-1">无法扫码？手动输入密钥：</p>
          <div class="flex items-center gap-2 p-3 bg-slate-50 dark:bg-zinc-800 rounded-lg border border-border-light dark:border-border-dark">
            <code class="text-xs font-mono flex-1 break-all select-all">{{ secret }}</code>
            <button class="text-xs text-primary hover:underline flex-shrink-0" @click="$emit('copy-secret')">
              复制
            </button>
          </div>
        </div>
      </div>
      <el-form @submit.prevent="$emit('confirm')">
        <el-form-item label="验证码">
          <el-input v-model="codeModel" placeholder="请输入6位验证码" maxlength="6">
            <template #prefix>
              <el-icon><Key /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" :loading="loading" :disabled="code.length !== 6" @click="$emit('confirm')">
          确认绑定
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject } from 'vue'
import { Key } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  qr: {
    type: String,
    default: ''
  },
  secret: {
    type: String,
    default: ''
  },
  code: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'update:code', 'confirm', 'copy-secret', 'closed'])
const t = inject('t', (key) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const codeModel = computed({
  get: () => props.code,
  set: (value) => emit('update:code', value)
})

const handleClose = () => {
  emit('closed')
}
</script>

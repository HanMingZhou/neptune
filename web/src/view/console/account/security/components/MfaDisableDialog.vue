<template>
  <el-dialog
    v-model="dialogVisible"
    title="关闭 MFA 验证"
    width="400px"
    class="custom-dialog"
    @close="handleClose"
  >
    <div class="py-4 space-y-4">
      <div class="p-3 bg-amber-500/10 border border-amber-500/20 rounded-xl text-amber-600 text-sm">
        关闭 MFA 后，登录时将不再需要验证码。请输入当前的 MFA 验证码以确认关闭。
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
        <el-button type="danger" :loading="loading" :disabled="code.length !== 6" @click="$emit('confirm')">
          确认关闭
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
  code: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'update:code', 'confirm', 'closed'])
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

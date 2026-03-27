<template>
  <el-dialog
    v-model="dialogVisible"
    :title="title"
    width="400px"
    class="custom-dialog"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="py-4"
      @submit.prevent="handleSubmit"
    >
      <el-form-item :label="phoneLabel" prop="phone">
        <el-input v-model="form.phone" :placeholder="phonePlaceholder">
          <template #prefix>
            <el-icon><Phone /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item :label="codeLabel" prop="code">
        <div class="flex gap-4 w-full">
          <el-input v-model="form.code" :placeholder="codePlaceholder" class="flex-1">
            <template #prefix>
              <el-icon><Key /></el-icon>
            </template>
          </el-input>
          <el-button
            type="primary"
            :disabled="disableRequestCode"
            class="w-32"
            @click="emit('request-code')"
          >
            {{ requestCodeText }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ t('confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { Key, Phone } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  form: {
    type: Object,
    required: true
  },
  rules: {
    type: Object,
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  },
  disableRequestCode: {
    type: Boolean,
    default: false
  },
  requestCodeText: {
    type: String,
    default: '获取验证码'
  },
  phoneLabel: {
    type: String,
    default: '新手机号'
  },
  phonePlaceholder: {
    type: String,
    default: '请输入新手机号'
  },
  codeLabel: {
    type: String,
    default: '验证码'
  },
  codePlaceholder: {
    type: String,
    default: '请输入验证码'
  }
})

const emit = defineEmits(['update:modelValue', 'submit', 'request-code', 'closed'])
const t = inject('t', (key) => key)
const formRef = ref(null)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const needsValidation = computed(() => Object.keys(props.rules || {}).length > 0)

const handleSubmit = () => {
  if (!formRef.value || !needsValidation.value) {
    emit('submit')
    return
  }

  formRef.value.validate((valid) => {
    if (valid) {
      emit('submit')
    }
  })
}

const handleClose = () => {
  formRef.value?.clearValidate()
  emit('closed')
}
</script>

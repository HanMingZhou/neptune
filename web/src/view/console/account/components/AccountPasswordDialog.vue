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
      label-width="90px"
      label-position="top"
      class="py-4"
      @submit.prevent="handleSubmit"
    >
      <el-form-item label="原密码" prop="password">
        <el-input v-model="form.password" show-password placeholder="原密码">
          <template #prefix>
            <el-icon :class="iconClass"><Lock /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input v-model="form.newPassword" show-password placeholder="新密码">
          <template #prefix>
            <el-icon :class="iconClass"><Lock /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input v-model="form.confirmPassword" show-password placeholder="确认密码">
          <template #prefix>
            <el-icon :class="iconClass"><Lock /></el-icon>
          </template>
        </el-input>
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
import { Lock } from '@element-plus/icons-vue'

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
  iconClass: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'submit', 'closed'])
const t = inject('t', (key) => key)
const formRef = ref(null)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const handleSubmit = () => {
  if (!formRef.value) {
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

<template>
  <BaseFormDialog
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    form-class="reset-password-dialog__form"
    label-width="96px"
    :model="form"
    :submit-text="t('confirm')"
    :title="t('resetPassword')"
    class="reset-password-dialog"
    width="500px"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <div class="reset-password-dialog__content">
      <el-form-item :label="t('username')">
        <el-input v-model="form.userName" disabled />
      </el-form-item>
      <el-form-item :label="t('nickname')">
        <el-input v-model="form.nickName" disabled />
      </el-form-item>
      <el-form-item :label="t('password')">
        <div class="reset-password-dialog__password-row">
          <el-input
            v-model="form.password"
            class="flex-1"
            :placeholder="t('password')"
            show-password
          />
          <el-button type="primary" @click="$emit('generate-password')">
            {{ t('refresh') }}
          </el-button>
        </div>
      </el-form-item>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { UserResetPasswordForm } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    form: UserResetPasswordForm
    modelValue?: boolean
  }>(),
  {
    modelValue: false
  }
)

const emit = defineEmits<{
  close: []
  'generate-password': []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

<style scoped>
.reset-password-dialog__content {
  display: grid;
  gap: 0.2rem;
}

.reset-password-dialog__content :deep(.el-form-item) {
  margin-bottom: 0.9rem;
}

.reset-password-dialog__content :deep(.el-form-item__label) {
  justify-content: flex-end;
  color: rgb(71 85 105);
  line-height: 1.2rem;
}

.reset-password-dialog__content :deep(.el-form-item__content) {
  min-width: 0;
}

.reset-password-dialog__password-row {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.75rem;
}

.reset-password-dialog__password-row :deep(.el-input) {
  min-width: 0;
}

.reset-password-dialog__password-row :deep(.el-button) {
  height: 42px;
  padding-left: 1rem;
  padding-right: 1rem;
  white-space: nowrap;
}

.dark .reset-password-dialog__content :deep(.el-form-item__label) {
  color: rgb(148 163 184);
}
</style>

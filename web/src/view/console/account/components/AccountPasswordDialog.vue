<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    :class="dialogClass"
    :model="form"
    :rules="rules"
    :shell-size="'md'"
    :submit-text="t('confirm')"
    :submitting="loading"
    :title="title"
    form-class="py-4"
    label-position="top"
    label-width="90px"
    @closed="emit('closed')"
    @submit="emit('submit')"
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
      <el-input
        v-model="form.confirmPassword"
        show-password
        placeholder="确认密码"
      >
        <template #prefix>
          <el-icon :class="iconClass"><Lock /></el-icon>
        </template>
      </el-input>
    </el-form-item>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { Lock } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { AccountPasswordForm } from '@/types/account'

const props = withDefaults(
  defineProps<{
    dialogClass?: string
    form: AccountPasswordForm
    iconClass?: string
    loading?: boolean
    modelValue?: boolean
    rules?: FormRules<AccountPasswordForm>
    title?: string
  }>(),
  {
    dialogClass: '',
    iconClass: '',
    loading: false,
    modelValue: false,
    rules: () => ({}),
    title: ''
  }
)

const emit = defineEmits<{
  closed: []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

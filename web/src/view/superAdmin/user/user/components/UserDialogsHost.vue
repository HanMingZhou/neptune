<template>
  <ResetPasswordDialog
    :model-value="resetPwdDialog"
    :form="resetPwdInfo"
    @close="emit('close-reset-password')"
    @generate-password="emit('generate-password')"
    @submit="emit('submit-reset-password')"
  />

  <UserEditorDrawer
    :model-value="addUserDialog"
    :auth-options="authOptions"
    :dialog-flag="dialogFlag"
    :form="userInfo"
    :rules="rules"
    @close="emit('close-add-user')"
    @submit="emit('submit-user')"
  />
</template>

<script setup lang="ts">
import type { FormRules } from 'element-plus'
import ResetPasswordDialog from './ResetPasswordDialog.vue'
import UserEditorDrawer from './UserEditorDrawer.vue'
import type {
  UserAuthority,
  UserForm,
  UserResetPasswordForm
} from '@/types/superAdmin'

withDefaults(
  defineProps<{
    addUserDialog?: boolean
    authOptions?: UserAuthority[]
    dialogFlag?: 'add' | 'edit'
    resetPwdDialog?: boolean
    resetPwdInfo: UserResetPasswordForm
    rules: FormRules<UserForm>
    userInfo: UserForm
  }>(),
  {
    addUserDialog: false,
    authOptions: () => [],
    dialogFlag: 'add',
    resetPwdDialog: false
  }
)

const emit = defineEmits<{
  'close-add-user': []
  'close-reset-password': []
  'generate-password': []
  'submit-reset-password': []
  'submit-user': []
}>()
</script>

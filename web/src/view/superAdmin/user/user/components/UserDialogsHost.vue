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

<script setup>
import ResetPasswordDialog from './ResetPasswordDialog.vue'
import UserEditorDrawer from './UserEditorDrawer.vue'

defineProps({
  addUserDialog: {
    type: Boolean,
    default: false
  },
  authOptions: {
    type: Array,
    default: () => []
  },
  dialogFlag: {
    type: String,
    default: 'add'
  },
  resetPwdDialog: {
    type: Boolean,
    default: false
  },
  resetPwdInfo: {
    type: Object,
    required: true
  },
  rules: {
    type: Object,
    required: true
  },
  userInfo: {
    type: Object,
    required: true
  }
})

const emit = defineEmits([
  'close-add-user',
  'close-reset-password',
  'generate-password',
  'submit-reset-password',
  'submit-user'
])
</script>

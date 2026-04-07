<template>
  <AccountPasswordDialog
    v-model="passwordDialog.visible"
    dialog-class="account-profile-dialog"
    :form="passwordDialog.form"
    :icon-class="passwordDialog.iconClass"
    :loading="passwordDialog.loading"
    :rules="passwordDialog.rules"
    :title="passwordDialog.title"
    @closed="passwordDialog.clear"
    @submit="passwordDialog.submit"
  />

  <AccountPhoneDialog
    v-model="phoneDialog.visible"
    dialog-class="account-profile-dialog"
    :form="phoneDialog.form"
    :loading="phoneDialog.loading"
    :rules="phoneDialog.rules"
    :disable-request-code="phoneDialog.disableRequestCode"
    :request-code-text="phoneDialog.requestCodeText"
    :title="phoneDialog.title"
    @closed="phoneDialog.close"
    @request-code="phoneDialog.requestCode"
    @submit="phoneDialog.submit"
  />

  <AccountEmailDialog
    v-model="emailDialog.visible"
    dialog-class="account-profile-dialog"
    :form="emailDialog.form"
    :loading="emailDialog.loading"
    :rules="emailDialog.rules"
    :disable-request-code="emailDialog.disableRequestCode"
    :request-code-text="emailDialog.requestCodeText"
    :title="emailDialog.title"
    @closed="emailDialog.close"
    @request-code="emailDialog.requestCode"
    @submit="emailDialog.submit"
  />
</template>

<script setup lang="ts">
import type { FormRules } from 'element-plus'
import type {
  AccountEmailForm,
  AccountPasswordForm,
  AccountPhoneForm
} from '@/types/account'
import AccountEmailDialog from '../../components/AccountEmailDialog.vue'
import AccountPasswordDialog from '../../components/AccountPasswordDialog.vue'
import AccountPhoneDialog from '../../components/AccountPhoneDialog.vue'

interface PasswordDialogController {
  visible: boolean
  form: AccountPasswordForm
  iconClass?: string
  loading: boolean
  rules: FormRules<AccountPasswordForm>
  title: string
  clear: () => void
  submit: () => void | Promise<void>
}

interface PhoneDialogController {
  visible: boolean
  form: AccountPhoneForm
  loading: boolean
  rules: FormRules<AccountPhoneForm>
  disableRequestCode: boolean
  requestCodeText: string
  title: string
  close: () => void
  requestCode: () => void
  submit: () => void | Promise<void>
}

interface EmailDialogController {
  visible: boolean
  form: AccountEmailForm
  loading: boolean
  rules: FormRules<AccountEmailForm>
  disableRequestCode: boolean
  requestCodeText: string
  title: string
  close: () => void
  requestCode: () => void
  submit: () => void | Promise<void>
}

defineProps<{
  emailDialog: EmailDialogController
  passwordDialog: PasswordDialogController
  phoneDialog: PhoneDialogController
}>()
</script>

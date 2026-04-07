<template>
  <AccountPasswordDialog
    v-model="passwordDialog.visible"
    dialog-class="account-security-dialog"
    :form="passwordDialog.form"
    :rules="passwordDialog.rules"
    :loading="passwordDialog.loading"
    :title="passwordDialog.title"
    @closed="passwordDialog.clear"
    @submit="passwordDialog.submit"
  />

  <AccountPhoneDialog
    v-model="phoneDialog.visible"
    dialog-class="account-security-dialog"
    :form="phoneDialog.form"
    :rules="phoneDialog.rules"
    :loading="phoneDialog.loading"
    :disable-request-code="phoneDialog.disableRequestCode"
    :request-code-text="phoneDialog.requestCodeText"
    :title="phoneDialog.title"
    @closed="phoneDialog.close"
    @request-code="phoneDialog.requestCode"
    @submit="phoneDialog.submit"
  />

  <AccountEmailDialog
    v-model="emailDialog.visible"
    dialog-class="account-security-dialog"
    :form="emailDialog.form"
    :rules="emailDialog.rules"
    :loading="emailDialog.loading"
    :disable-request-code="emailDialog.disableRequestCode"
    :request-code-text="emailDialog.requestCodeText"
    :title="emailDialog.title"
    @closed="emailDialog.close"
    @request-code="emailDialog.requestCode"
    @submit="emailDialog.submit"
  />

  <MfaSetupDialog
    v-model="mfaSetupDialog.visible"
    :code="mfaSetupDialog.code"
    :loading="mfaSetupDialog.loading"
    :qr="mfaSetupDialog.qr"
    :secret="mfaSetupDialog.secret"
    @closed="mfaSetupDialog.close"
    @confirm="mfaSetupDialog.confirm"
    @copy-secret="mfaSetupDialog.copySecret"
    @update:code="mfaSetupDialog.code = $event"
  />

  <MfaDisableDialog
    v-model="mfaDisableDialog.visible"
    :code="mfaDisableDialog.code"
    :loading="mfaDisableDialog.loading"
    @closed="mfaDisableDialog.close"
    @confirm="mfaDisableDialog.confirm"
    @update:code="mfaDisableDialog.code = $event"
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
import MfaDisableDialog from './MfaDisableDialog.vue'
import MfaSetupDialog from './MfaSetupDialog.vue'

interface PasswordDialogController {
  visible: boolean
  form: AccountPasswordForm
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

interface MfaSetupDialogController {
  visible: boolean
  code: string
  loading: boolean
  qr: string
  secret: string
  close: () => void
  confirm: () => void | Promise<void>
  copySecret: () => void | Promise<void>
}

interface MfaDisableDialogController {
  visible: boolean
  code: string
  loading: boolean
  close: () => void
  confirm: () => void | Promise<void>
}

defineProps<{
  emailDialog: EmailDialogController
  mfaDisableDialog: MfaDisableDialogController
  mfaSetupDialog: MfaSetupDialogController
  passwordDialog: PasswordDialogController
  phoneDialog: PhoneDialogController
}>()
</script>

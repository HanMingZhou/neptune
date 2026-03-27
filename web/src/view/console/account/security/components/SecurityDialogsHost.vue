<template>
  <AccountPasswordDialog
    v-model="passwordDialog.visible"
    :form="passwordDialog.form"
    :rules="passwordDialog.rules"
    :loading="passwordDialog.loading"
    :title="passwordDialog.title"
    @closed="passwordDialog.clear"
    @submit="passwordDialog.submit"
  />

  <AccountPhoneDialog
    v-model="phoneDialog.visible"
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

<script setup>
import AccountEmailDialog from '../../components/AccountEmailDialog.vue'
import AccountPasswordDialog from '../../components/AccountPasswordDialog.vue'
import AccountPhoneDialog from '../../components/AccountPhoneDialog.vue'
import MfaDisableDialog from './MfaDisableDialog.vue'
import MfaSetupDialog from './MfaSetupDialog.vue'

defineProps({
  emailDialog: {
    type: Object,
    required: true
  },
  mfaDisableDialog: {
    type: Object,
    required: true
  },
  mfaSetupDialog: {
    type: Object,
    required: true
  },
  passwordDialog: {
    type: Object,
    required: true
  },
  phoneDialog: {
    type: Object,
    required: true
  }
})
</script>

<style>
.custom-dialog {
  :deep(.el-dialog) {
    @apply rounded-xl overflow-hidden;
  }

  :deep(.el-dialog__header) {
    @apply m-0 p-6 border-b border-slate-100 dark:border-slate-800;
  }

  :deep(.el-dialog__title) {
    @apply text-lg font-semibold text-slate-800 dark:text-slate-200;
  }

  :deep(.el-dialog__body) {
    @apply p-6 bg-slate-50 dark:bg-slate-900/50;
  }

  :deep(.el-dialog__footer) {
    @apply p-6 border-t border-slate-100 dark:border-slate-800 m-0 bg-white dark:bg-slate-900;
  }

  :deep(.el-form-item__label) {
    @apply pb-2 font-medium text-slate-700 dark:text-slate-300;
  }

  :deep(.el-input__wrapper) {
    @apply px-4 py-2 transition-all duration-200;
  }

  :deep(.el-input__wrapper.is-focus) {
    @apply shadow-[0_0_0_2px_rgba(59,130,246,0.1)] dark:shadow-[0_0_0_2px_rgba(59,130,246,0.2)];
  }
}
</style>

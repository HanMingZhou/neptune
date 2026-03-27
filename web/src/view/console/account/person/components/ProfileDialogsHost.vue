<template>
  <AccountPasswordDialog
    v-model="passwordDialog.visible"
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

<script setup>
import AccountEmailDialog from '../../components/AccountEmailDialog.vue'
import AccountPasswordDialog from '../../components/AccountPasswordDialog.vue'
import AccountPhoneDialog from '../../components/AccountPhoneDialog.vue'

defineProps({
  emailDialog: {
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
  :deep(.el-dialog__header) {
    @apply mb-0 pb-4 border-b border-gray-100 dark:border-gray-700;
  }

  :deep(.el-dialog__footer) {
    @apply mt-0 pt-4 border-t border-gray-100 dark:border-gray-700;
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

<template>
  <div class="max-w-[1000px] mx-auto space-y-8">
    <SecuritySettingsPanel
      :account-info="accountInfo"
      :ak-loading="akLoading"
      :get-status-color="getStatusColor"
      :mask-string="maskString"
      :mfa-loading="mfaLoading"
      @bind-github="bindGithub"
      @change-email="emailDialog.visible = true"
      @change-password="passwordDialog.visible = true"
      @change-phone="phoneDialog.visible = true"
      @copy-access-key="copyAccessKeyId"
      @generate-access-key="handleGenerateAK"
      @mfa-action="handleMfaAction"
      @setup-notification="setupNotification"
    />

    <SecurityDialogsHost
      :email-dialog="emailDialog"
      :mfa-disable-dialog="mfaDisableDialog"
      :mfa-setup-dialog="mfaSetupDialog"
      :password-dialog="passwordDialog"
      :phone-dialog="phoneDialog"
    />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import SecurityDialogsHost from './components/SecurityDialogsHost.vue'
import SecuritySettingsPanel from './components/SecuritySettingsPanel.vue'
import { useAccountSecurity } from './composables/useAccountSecurity'

const t = inject('t', (key) => key)

const {
  accountInfo,
  akLoading,
  bindGithub,
  copyAccessKeyId,
  emailDialog,
  getStatusColor,
  handleGenerateAK,
  handleMfaAction,
  maskString,
  mfaDisableDialog,
  mfaLoading,
  mfaSetupDialog,
  passwordDialog,
  phoneDialog,
  setupNotification
} = useAccountSecurity({ t })
</script>

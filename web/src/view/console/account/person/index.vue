<template>
  <div>
    <ProfileOverviewPanel
      :edit-flag="editFlag"
      :nick-name="nickName"
      :user-info="userStore.userInfo"
      @change-email="emailDialog.visible = true"
      @change-password="passwordDialog.visible = true"
      @change-phone="phoneDialog.visible = true"
      @close-edit="closeEdit"
      @confirm-edit="enterEdit"
      @open-edit="openEdit"
      @update:nick-name="nickName = $event"
    />

    <ProfileDialogsHost
      :email-dialog="emailDialog"
      :password-dialog="passwordDialog"
      :phone-dialog="phoneDialog"
    />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import ProfileDialogsHost from './components/ProfileDialogsHost.vue'
import ProfileOverviewPanel from './components/ProfileOverviewPanel.vue'
import { useAccountProfile } from './composables/useAccountProfile'

const t = inject('t', (key) => key)

const {
  closeEdit,
  editFlag,
  emailDialog,
  enterEdit,
  nickName,
  openEdit,
  passwordDialog,
  phoneDialog,
  userStore
} = useAccountProfile({ t })
</script>

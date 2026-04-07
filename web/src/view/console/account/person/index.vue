<template>
  <div class="console-page-container space-y-6">
    <BaseTableToolbar
      :breadcrumbs="[t('person')]"
      :description="t('profileCenterDesc')"
      :show-refresh="false"
      :title="t('person')"
    />

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

<script setup lang="ts">
import { inject } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import ProfileDialogsHost from './components/ProfileDialogsHost.vue'
import ProfileOverviewPanel from './components/ProfileOverviewPanel.vue'
import { useAccountProfile } from './composables/useAccountProfile'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)

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

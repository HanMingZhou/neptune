<template>
  <div class="max-w-[1400px] mx-auto space-y-6">
    <UserManagementHeader
      :loading="loading"
      @create="openCreateDialog"
      @refresh="getTableData"
    />

    <UserTableCard
      :auth-options="authOptions"
      :items="tableData"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-info="searchInfo"
      :total="total"
      @authority-dirty="markAuthorityDirty"
      @change-authority="changeAuthority"
      @delete="deleteUserFunc"
      @edit="openEditDialog"
      @page-change="handleCurrentChange"
      @reset="onReset"
      @reset-password="openResetPasswordDialog"
      @search="onSubmit"
      @size-change="handleSizeChange"
      @switch-enable="switchEnable"
    />

    <UserDialogsHost
      :add-user-dialog="addUserDialog"
      :auth-options="authOptions"
      :dialog-flag="dialogFlag"
      :reset-pwd-dialog="resetPwdDialog"
      :reset-pwd-info="resetPwdInfo"
      :rules="rules"
      :user-info="userInfo"
      @close-add-user="closeAddUserDialog"
      @close-reset-password="closeResetPwdDialog"
      @generate-password="generateRandomPassword"
      @submit-reset-password="confirmResetPassword"
      @submit-user="submitUserDialog"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import UserDialogsHost from './components/UserDialogsHost.vue'
import UserManagementHeader from './components/UserManagementHeader.vue'
import UserTableCard from './components/UserTableCard.vue'
import { useUserManagementPage } from './composables/useUserManagementPage'

const t = inject('t', (key) => key)

const {
  addUserDialog,
  authOptions,
  changeAuthority,
  closeAddUserDialog,
  closeResetPwdDialog,
  confirmResetPassword,
  deleteUserFunc,
  dialogFlag,
  generateRandomPassword,
  getTableData,
  handleCurrentChange,
  handleSizeChange,
  initialize,
  loading,
  markAuthorityDirty,
  onReset,
  onSubmit,
  openCreateDialog,
  openEditDialog,
  openResetPasswordDialog,
  page,
  pageSize,
  resetPwdDialog,
  resetPwdInfo,
  rules,
  searchInfo,
  submitUserDialog,
  switchEnable,
  tableData,
  total,
  userInfo
} = useUserManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

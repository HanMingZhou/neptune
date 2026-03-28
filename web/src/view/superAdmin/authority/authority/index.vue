<template>
  <div class="console-page-container space-y-6">
    <AuthorityManagementHeader
      :loading="loading"
      @create="addAuthority(0)"
      @refresh="getTableData"
    />

    <AuthorityTableCard
      :items="filteredTableData"
      :loading="loading"
      :search-keyword="searchKeyword"
      @add-child="addAuthority"
      @configure="openDrawer"
      @copy="copyAuthorityFunc"
      @delete="deleteAuth"
      @edit="editAuthority"
      @reset="handleResetSearch"
      @search="handleSearch"
      @update:search-keyword="searchKeyword = $event"
    />

    <AuthorityDialogsHost
      :active-row="activeRow"
      :authority-form-visible="authorityFormVisible"
      :authority-options="authorityOptions"
      :dialog-title="dialogTitle"
      :dialog-type="dialogType"
      :drawer="drawer"
      :form="form"
      :rules="rules"
      :submitting="submitting"
      :table-data="tableData"
      @change-row="changeRow"
      @close-authority-form="closeAuthorityForm"
      @close-drawer="closeDrawer"
      @submit-authority-form="submitAuthorityForm"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import AuthorityDialogsHost from './components/AuthorityDialogsHost.vue'
import AuthorityManagementHeader from './components/AuthorityManagementHeader.vue'
import AuthorityTableCard from './components/AuthorityTableCard.vue'
import { useAuthorityManagementPage } from './composables/useAuthorityManagementPage'

const t = inject('t', (key) => key)

const {
  activeRow,
  addAuthority,
  authorityFormVisible,
  authorityOptions,
  changeRow,
  closeAuthorityForm,
  closeDrawer,
  copyAuthorityFunc,
  deleteAuth,
  dialogTitle,
  dialogType,
  drawer,
  editAuthority,
  filteredTableData,
  form,
  getTableData,
  handleResetSearch,
  handleSearch,
  initialize,
  loading,
  openDrawer,
  rules,
  searchKeyword,
  submitAuthorityForm,
  submitting,
  tableData
} = useAuthorityManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <AuthorityManagementHeader
      :loading="loading"
      @create="addAuthority(0)"
      @refresh="initialize"
    />

    <AuthorityTableCard
      class="min-h-0 flex-1"
      :items="pagedTableData"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-keyword="searchKeyword"
      :total="total"
      @add-child="addAuthority"
      @configure="openDrawer"
      @copy="copyAuthorityFunc"
      @delete="deleteAuth"
      @edit="editAuthority"
      @page-change="handleCurrentChange"
      @reset="handleResetSearch"
      @search="handleSearch"
      @size-change="handleSizeChange"
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

<script setup lang="ts">
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
  form,
  handleCurrentChange,
  handleResetSearch,
  handleSearch,
  handleSizeChange,
  initialize,
  loading,
  openDrawer,
  page,
  pageSize,
  pagedTableData,
  rules,
  searchKeyword,
  submitAuthorityForm,
  submitting,
  tableData,
  total
} = useAuthorityManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

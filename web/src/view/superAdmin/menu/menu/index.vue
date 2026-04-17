<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('admin'), t('menus')]"
      :description="t('manageMenusDesc')"
      :loading="loading"
      :title="t('menus')"
      @refresh="initialize"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog(0)"
        >
          <span class="material-icons">add</span>
          {{ t('addRootMenu') }}
        </button>
      </template>
    </BaseTableToolbar>

    <MenuTableCard
      class="min-h-0 flex-1"
      :items="pagedTableData"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :search-keyword="searchKeyword"
      :total="total"
      @add="openCreateDialog"
      @delete="handleDeleteMenu"
      @edit="openEditDialog"
      @page-change="handleCurrentChange"
      @reset="handleResetSearch"
      @search="handleSearch"
      @size-change="handleSizeChange"
      @update:search-keyword="searchKeyword = $event"
    />

    <MenuEditorDrawer
      v-model="dialogFormVisible"
      v-model:check-flag="checkFlag"
      :dialog-title="dialogTitle"
      :form="form"
      :is-edit="isEdit"
      :menu-options="menuOptions"
      :rules="rules"
      @add-button="addButton"
      @add-parameter="addParameter"
      @close="closeDialog"
      @component-change="fmtComponent"
      @delete-button="deleteButton"
      @delete-parameter="deleteParameter"
      @name-change="changeName"
      @submit="handleSubmitMenu"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import MenuEditorDrawer from './components/MenuEditorDrawer.vue'
import MenuTableCard from './components/MenuTableCard.vue'
import { useMenuManagementPage } from './composables/useMenuManagementPage'

const t = inject('t', (key) => key)

const {
  addButton,
  addParameter,
  changeName,
  checkFlag,
  closeDialog,
  deleteButton,
  deleteParameter,
  dialogFormVisible,
  dialogTitle,
  fmtComponent,
  form,
  handleCurrentChange,
  handleDeleteMenu,
  handleResetSearch,
  handleSearch,
  handleSizeChange,
  handleSubmitMenu,
  initialize,
  isEdit,
  loading,
  menuOptions,
  openCreateDialog,
  openEditDialog,
  page,
  pageSize,
  pagedTableData,
  rules,
  searchKeyword,
  total
} = useMenuManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('admin'), t('menus')]"
      :description="t('manageMenusDesc')"
      :title="t('menus')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="getTableData" />
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog(0)"
        >
          <span class="material-icons">add</span>
          {{ t('addRootMenu') }}
        </button>
      </template>
    </PageIntro>

    <MenuTableCard
      :items="filteredTableData"
      :loading="loading"
      :search-keyword="searchKeyword"
      @add="openCreateDialog"
      @delete="handleDeleteMenu"
      @edit="openEditDialog"
      @reset="handleResetSearch"
      @search="handleSearch"
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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
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
  filteredTableData,
  fmtComponent,
  form,
  getTableData,
  handleDeleteMenu,
  handleResetSearch,
  handleSearch,
  handleSubmitMenu,
  initialize,
  isEdit,
  loading,
  menuOptions,
  openCreateDialog,
  openEditDialog,
  rules,
  searchKeyword
} = useMenuManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

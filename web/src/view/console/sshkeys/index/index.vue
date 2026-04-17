<template>
  <div class="console-page-container flex min-h-full flex-col gap-6">
    <BaseTableToolbar
      :breadcrumbs="[t('management'), t('sshKeyManage')]"
      :description="t('sshKeyManageDesc')"
      :loading="loading"
      :title="t('sshKeyManage')"
      @refresh="loadKeys"
    >
      <template #actions>
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('newSshKey') }}
        </button>
      </template>
    </BaseTableToolbar>

    <ManagementListShell class="min-h-0 flex-1">
      <template #filters>
        <SshKeyFiltersBar
          :search-name="searchName"
          @search="searchKeys"
          @update:search-name="searchName = $event"
        />
      </template>

      <SshKeyList
        class="min-h-0 flex-1"
        :items="keys"
        :loading="loading"
        :page="page"
        :page-size="pageSize"
        :total="total"
        @create="openCreateDialog"
        @delete="handleDelete"
        @page-change="handleCurrentChange"
        @size-change="handleSizeChange"
        @set-default="setDefault"
      />
    </ManagementListShell>

    <SshKeyCreateDialog
      v-model="showCreateDialog"
      :form="createForm"
      :loading="creating"
      :rules="rules"
      @close="closeCreateDialog"
      @submit="handleCreate"
    />
  </div>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import type { Translator } from '@/types/consoleResource'
import BaseTableToolbar from '@/components/listPage/BaseTableToolbar.vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import SshKeyCreateDialog from './components/SshKeyCreateDialog.vue'
import SshKeyFiltersBar from './components/SshKeyFiltersBar.vue'
import SshKeyList from './components/SshKeyList.vue'
import { useSshKeyManagementPage } from './composables/useSshKeyManagementPage'

const t = inject<Translator>('t', (key: string) => key)

const {
  closeCreateDialog,
  createForm,
  creating,
  handleCreate,
  handleDelete,
  handleCurrentChange,
  handleSizeChange,
  initialize,
  keys,
  loading,
  loadKeys,
  openCreateDialog,
  page,
  pageSize,
  rules,
  searchKeys,
  searchName,
  setDefault,
  showCreateDialog,
  total
} = useSshKeyManagementPage({ t })

onMounted(() => {
  void initialize()
})
</script>

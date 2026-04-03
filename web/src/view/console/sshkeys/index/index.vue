<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('management'), t('sshKeyManage')]"
      :description="t('sshKeyManageDesc')"
      :title="t('sshKeyManage')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="loadKeys" />
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('newSshKey') }}
        </button>
      </template>
    </PageIntro>

    <ManagementListShell>
      <template #filters>
        <SshKeyFiltersBar
          :search-name="searchName"
          @search="searchKeys"
          @update:search-name="searchName = $event"
        />
      </template>

      <SshKeyList
        :items="keys"
        :loading="loading"
        @create="openCreateDialog"
        @delete="handleDelete"
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

<script setup>
import { inject, onMounted } from 'vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import SshKeyCreateDialog from './components/SshKeyCreateDialog.vue'
import SshKeyFiltersBar from './components/SshKeyFiltersBar.vue'
import SshKeyList from './components/SshKeyList.vue'
import { useSshKeyManagementPage } from './composables/useSshKeyManagementPage'

const t = inject('t', (key) => key)

const {
  closeCreateDialog,
  createForm,
  creating,
  handleCreate,
  handleDelete,
  initialize,
  keys,
  loading,
  loadKeys,
  openCreateDialog,
  rules,
  searchKeys,
  searchName,
  setDefault,
  showCreateDialog
} = useSshKeyManagementPage({ t })

onMounted(() => {
  initialize()
})
</script>

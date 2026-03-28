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
          class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('newSshKey') }}
        </button>
      </template>
    </PageIntro>

    <SshKeyFiltersBar
      :search-name="searchName"
      @search="searchKeys"
      @update:search-name="searchName = $event"
    />

    <SshKeyList
      :items="keys"
      :loading="loading"
      @create="openCreateDialog"
      @delete="handleDelete"
      @set-default="setDefault"
    />

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

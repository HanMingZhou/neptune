<template>
  <div class="console-page-container space-y-6">
    <PageIntro
      :breadcrumbs="[t('resources'), t('imageManage')]"
      :description="t('imageManageDesc')"
      :title="t('imageManage')"
    >
      <template #actions>
        <RefreshButton :loading="loading" @refresh="handleRefresh" />
        <button
          class="list-toolbar-button list-toolbar-button--primary"
          @click="openCreateDialog"
        >
          <span class="material-icons text-[20px]">add</span>
          {{ t('imageAdd') }}
        </button>
      </template>
    </PageIntro>

    <ManagementListShell>
      <template #filters>
        <ImageFiltersCard
          v-model:filter-keyword="filterKeyword"
          v-model:filter-type="filterType"
          v-model:filter-usage-type="filterUsageType"
          @reset="handleReset"
          @search="handleSearch"
        />
      </template>

      <ImageTableCard
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :items="images"
        :loading="loading"
        :total="total"
        @delete="handleDelete"
        @edit="openEditDialog"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </ManagementListShell>

    <ImageDialog
      v-model="showDialog"
      :form="form"
      :is-edit="isEdit"
      :rules="formRules"
      :submitting="submitting"
      @closed="handleDialogClosed"
      @submit="submitImage"
    />
  </div>
</template>

<script setup>
import { inject, onMounted } from 'vue'
import ManagementListShell from '@/components/listPage/ManagementListShell.vue'
import RefreshButton from '@/components/RefreshButton/index.vue'
import PageIntro from '@/components/listPage/PageIntro.vue'
import ImageDialog from './components/ImageDialog.vue'
import ImageFiltersCard from './components/ImageFiltersCard.vue'
import ImageTableCard from './components/ImageTableCard.vue'
import { useImageManagement } from './composables/useImageManagement'

const t = inject('t', (key) => key)

const {
  currentPage,
  fetchImages,
  filterKeyword,
  filterType,
  filterUsageType,
  form,
  formRules,
  handleDelete,
  handleDialogClosed,
  handlePageChange,
  handleRefresh,
  handleReset,
  handleSearch,
  images,
  isEdit,
  loading,
  openCreateDialog,
  openEditDialog,
  pageSize,
  showDialog,
  submitImage,
  submitting,
  total
} = useImageManagement({ t })

onMounted(() => {
  fetchImages()
})
</script>

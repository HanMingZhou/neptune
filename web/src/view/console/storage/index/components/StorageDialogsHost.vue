<template>
  <CreateStorageDialog
    :model-value="showCreateDialog"
    :cluster-options="clusterOptions"
    :creating="creating"
    :form="createForm"
    :storage-products="storageProducts"
    @cluster-change="emit('cluster-change', $event)"
    @submit="emit('submit-create')"
    @update:model-value="emit('update:create-visible', $event)"
  />

  <ExpandStorageDialog
    :model-value="showExpandDialog"
    :expanding="expanding"
    :form="expandForm"
    @submit="emit('submit-expand')"
    @update:model-value="emit('update:expand-visible', $event)"
  />

  <EditStorageDialog
    :model-value="showEditDialog"
    :editing="editing"
    :form="editForm"
    @submit="emit('submit-edit')"
    @update:model-value="emit('update:edit-visible', $event)"
  />
</template>

<script setup lang="ts">
import type {
  StorageClusterOption,
  StorageCreateForm,
  StorageEditForm,
  StorageExpandForm,
  StorageProductOption
} from '@/types/storage'
import CreateStorageDialog from './CreateStorageDialog.vue'
import EditStorageDialog from './EditStorageDialog.vue'
import ExpandStorageDialog from './ExpandStorageDialog.vue'

withDefaults(
  defineProps<{
    clusterOptions?: StorageClusterOption[]
    createForm: StorageCreateForm
    creating?: boolean
    editForm: StorageEditForm
    editing?: boolean
    expanding?: boolean
    expandForm: StorageExpandForm
    showCreateDialog?: boolean
    showEditDialog?: boolean
    showExpandDialog?: boolean
    storageProducts?: StorageProductOption[]
  }>(),
  {
    clusterOptions: () => [],
    creating: false,
    editing: false,
    expanding: false,
    showCreateDialog: false,
    showEditDialog: false,
    showExpandDialog: false,
    storageProducts: () => []
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: StorageCreateForm['clusterId']]
  'submit-create': []
  'submit-edit': []
  'submit-expand': []
  'update:create-visible': [value: boolean]
  'update:edit-visible': [value: boolean]
  'update:expand-visible': [value: boolean]
}>()
</script>

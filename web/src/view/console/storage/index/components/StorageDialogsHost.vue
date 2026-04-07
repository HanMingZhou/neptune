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
</template>

<script setup lang="ts">
import type {
  StorageClusterOption,
  StorageCreateForm,
  StorageExpandForm,
  StorageProductOption
} from '@/types/storage'
import CreateStorageDialog from './CreateStorageDialog.vue'
import ExpandStorageDialog from './ExpandStorageDialog.vue'

withDefaults(
  defineProps<{
    clusterOptions?: StorageClusterOption[]
    createForm: StorageCreateForm
    creating?: boolean
    expanding?: boolean
    expandForm: StorageExpandForm
    showCreateDialog?: boolean
    showExpandDialog?: boolean
    storageProducts?: StorageProductOption[]
  }>(),
  {
    clusterOptions: () => [],
    creating: false,
    expanding: false,
    showCreateDialog: false,
    showExpandDialog: false,
    storageProducts: () => []
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: StorageCreateForm['clusterId']]
  'submit-create': []
  'submit-expand': []
  'update:create-visible': [value: boolean]
  'update:expand-visible': [value: boolean]
}>()
</script>

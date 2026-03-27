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

<script setup>
import CreateStorageDialog from './CreateStorageDialog.vue'
import ExpandStorageDialog from './ExpandStorageDialog.vue'

defineProps({
  clusterOptions: {
    type: Array,
    default: () => []
  },
  createForm: {
    type: Object,
    required: true
  },
  creating: {
    type: Boolean,
    default: false
  },
  expanding: {
    type: Boolean,
    default: false
  },
  expandForm: {
    type: Object,
    required: true
  },
  showCreateDialog: {
    type: Boolean,
    default: false
  },
  showExpandDialog: {
    type: Boolean,
    default: false
  },
  storageProducts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits([
  'cluster-change',
  'submit-create',
  'submit-expand',
  'update:create-visible',
  'update:expand-visible'
])
</script>

<style scoped>
:deep(.custom-select .el-input__wrapper) {
  background-color: transparent !important;
  box-shadow: 0 0 0 1px var(--el-border-color) inset !important;
}

:deep(.custom-select .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-color-primary) inset !important;
}
</style>

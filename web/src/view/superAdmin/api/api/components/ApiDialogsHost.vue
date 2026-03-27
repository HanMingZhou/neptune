<template>
  <SyncApiDrawer
    :model-value="syncApiFlag"
    :api-completion-loading="apiCompletionLoading"
    :api-group-options="apiGroupOptions"
    :sync-api-data="syncApiData"
    :syncing="syncing"
    @add-one="emit('add-one', $event)"
    @ai-completion="emit('ai-completion')"
    @close="emit('close-sync')"
    @ignore="emit('ignore', $event)"
    @submit="emit('submit-sync')"
    @update:model-value="emit('update-sync-visible', $event)"
  />

  <ApiEditorDrawer
    :model-value="dialogFormVisible"
    :api-group-options="apiGroupOptions"
    :dialog-title="dialogTitle"
    :form="form"
    :method-options="methodOptions"
    :rules="rules"
    @close="emit('close-dialog')"
    @submit="emit('submit-dialog')"
    @update:model-value="emit('update-dialog-visible', $event)"
  />
</template>

<script setup>
import ApiEditorDrawer from './ApiEditorDrawer.vue'
import SyncApiDrawer from './SyncApiDrawer.vue'

defineProps({
  apiCompletionLoading: {
    type: Boolean,
    default: false
  },
  apiGroupOptions: {
    type: Array,
    default: () => []
  },
  dialogFormVisible: {
    type: Boolean,
    default: false
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  form: {
    type: Object,
    required: true
  },
  methodOptions: {
    type: Array,
    default: () => []
  },
  rules: {
    type: Object,
    default: () => ({})
  },
  syncApiData: {
    type: Object,
    required: true
  },
  syncApiFlag: {
    type: Boolean,
    default: false
  },
  syncing: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'add-one',
  'ai-completion',
  'close-dialog',
  'close-sync',
  'ignore',
  'submit-dialog',
  'submit-sync',
  'update-dialog-visible',
  'update-sync-visible'
])
</script>

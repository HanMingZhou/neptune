<template>
  <ApiKeyDialog
    v-if="showApiKeyDialog"
    :api-keys="apiKeys"
    :format-time="formatTime"
    :new-key-name="newKeyName"
    :newly-created-key="newlyCreatedKey"
    @close="emit('close-api-key-dialog')"
    @copy="emit('copy', $event)"
    @create="emit('create-api-key')"
    @delete="emit('delete-api-key', $event)"
    @update:new-key-name="emit('update:new-key-name', $event)"
  />
</template>

<script setup>
import ApiKeyDialog from './ApiKeyDialog.vue'

defineProps({
  apiKeys: {
    type: Array,
    default: () => []
  },
  formatTime: {
    type: Function,
    required: true
  },
  newKeyName: {
    type: String,
    default: ''
  },
  newlyCreatedKey: {
    type: Object,
    default: null
  },
  showApiKeyDialog: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'close-api-key-dialog',
  'copy',
  'create-api-key',
  'delete-api-key',
  'update:new-key-name'
])
</script>

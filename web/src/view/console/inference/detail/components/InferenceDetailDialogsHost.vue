<template>
  <ApiKeyDialog
    v-if="showApiKeyDialog"
    :api-keys="apiKeys"
    :api-keys-loading="apiKeysLoading"
    :creating-api-key="creatingApiKey"
    :deleting-api-key-id="deletingApiKeyId"
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

<script setup lang="ts">
import ApiKeyDialog from './ApiKeyDialog.vue'
import type { InferenceApiKey } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    apiKeys?: InferenceApiKey[]
    apiKeysLoading?: boolean
    creatingApiKey?: boolean
    deletingApiKeyId?: number | null
    formatTime: (value?: string | number) => string
    newKeyName?: string
    newlyCreatedKey?: string
    showApiKeyDialog?: boolean
  }>(),
  {
    apiKeys: () => [],
    apiKeysLoading: false,
    creatingApiKey: false,
    deletingApiKeyId: null,
    newKeyName: '',
    newlyCreatedKey: '',
    showApiKeyDialog: false
  }
)

const emit = defineEmits<{
  'close-api-key-dialog': []
  copy: [text: string]
  'create-api-key': []
  'delete-api-key': [apiKey: InferenceApiKey]
  'update:new-key-name': [value: string]
}>()
</script>

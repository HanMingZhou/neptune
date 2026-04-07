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

<script setup lang="ts">
import type { FormRules } from 'element-plus'
import ApiEditorDrawer from './ApiEditorDrawer.vue'
import SyncApiDrawer from './SyncApiDrawer.vue'
import type {
  ApiForm,
  ApiListItem,
  ApiMethodOption,
  ApiSyncData,
  LabelValueOption
} from '@/types/superAdmin'

interface SyncIgnorePayload {
  row: ApiListItem
  flag: boolean
}

withDefaults(
  defineProps<{
    apiCompletionLoading?: boolean
    apiGroupOptions?: LabelValueOption[]
    dialogFormVisible?: boolean
    dialogTitle?: string
    form: ApiForm
    methodOptions?: ApiMethodOption[]
    rules?: FormRules<ApiForm>
    syncApiData: ApiSyncData
    syncApiFlag?: boolean
    syncing?: boolean
  }>(),
  {
    apiCompletionLoading: false,
    apiGroupOptions: () => [],
    dialogFormVisible: false,
    dialogTitle: '',
    methodOptions: () => [],
    rules: () => ({}),
    syncApiFlag: false,
    syncing: false
  }
)

const emit = defineEmits<{
  'add-one': [row: ApiListItem]
  'ai-completion': []
  'close-dialog': []
  'close-sync': []
  ignore: [payload: SyncIgnorePayload]
  'submit-dialog': []
  'submit-sync': []
  'update-dialog-visible': [value: boolean]
  'update-sync-visible': [value: boolean]
}>()
</script>

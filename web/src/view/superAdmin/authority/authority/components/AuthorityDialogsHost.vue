<template>
  <AuthorityEditorDrawer
    :model-value="authorityFormVisible"
    :authority-options="authorityOptions"
    :dialog-title="dialogTitle"
    :dialog-type="dialogType"
    :form="form"
    :rules="rules"
    :submitting="submitting"
    @close="emit('close-authority-form')"
    @submit="emit('submit-authority-form')"
  />

  <AuthorityPermissionDrawer
    :model-value="drawer"
    :active-row="activeRow"
    :authority="tableData"
    @change-row="(key, value) => emit('change-row', key, value)"
    @close="emit('close-drawer')"
  />
</template>

<script setup lang="ts">
import type { FormRules } from 'element-plus'
import AuthorityEditorDrawer from './AuthorityEditorDrawer.vue'
import AuthorityPermissionDrawer from './AuthorityPermissionDrawer.vue'
import type { ResourceId } from '@/types/consoleResource'
import type {
  AuthorityDialogType,
  AuthorityForm,
  AuthorityOption,
  AuthorityTreeNode
} from '@/types/superAdmin'

interface AuthorityPermissionRef {
  authorityId: ResourceId
}

interface AuthorityPermissionRow extends AuthorityTreeNode {
  dataAuthorityId?: AuthorityPermissionRef[]
  defaultRouter?: string
}

withDefaults(
  defineProps<{
    activeRow?: AuthorityPermissionRow
    authorityFormVisible?: boolean
    authorityOptions?: AuthorityOption[]
    dialogTitle?: string
    dialogType?: AuthorityDialogType
    drawer?: boolean
    form: AuthorityForm
    rules: FormRules<AuthorityForm>
    submitting?: boolean
    tableData?: AuthorityTreeNode[]
  }>(),
  {
    activeRow: () => ({}) as AuthorityPermissionRow,
    authorityFormVisible: false,
    authorityOptions: () => [],
    dialogTitle: '',
    dialogType: 'add',
    drawer: false,
    submitting: false,
    tableData: () => []
  }
)

const emit = defineEmits<{
  'change-row': [key: string, value: unknown]
  'close-authority-form': []
  'close-drawer': []
  'submit-authority-form': []
}>()
</script>

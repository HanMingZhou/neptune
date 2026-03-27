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

<script setup>
import AuthorityEditorDrawer from './AuthorityEditorDrawer.vue'
import AuthorityPermissionDrawer from './AuthorityPermissionDrawer.vue'

defineProps({
  activeRow: {
    type: Object,
    default: () => ({})
  },
  authorityFormVisible: {
    type: Boolean,
    default: false
  },
  authorityOptions: {
    type: Array,
    default: () => []
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  dialogType: {
    type: String,
    default: 'add'
  },
  drawer: {
    type: Boolean,
    default: false
  },
  form: {
    type: Object,
    required: true
  },
  rules: {
    type: Object,
    required: true
  },
  submitting: {
    type: Boolean,
    default: false
  },
  tableData: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['change-row', 'close-authority-form', 'close-drawer', 'submit-authority-form'])
</script>

<template>
  <BaseFormDrawer
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :size="480"
    :submit-text="t('confirm')"
    :submitting="submitting"
    :title="dialogTitle"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <el-form-item :label="t('parentMenu')" prop="parentId">
      <el-cascader
        v-model="form.parentId"
        style="width: 100%"
        :disabled="dialogType === 'add'"
        :options="authorityOptions"
        :props="cascaderProps"
        :show-all-levels="false"
        filterable
        :placeholder="t('selectParentRole')"
      />
    </el-form-item>
    <el-form-item :label="t('roleId')" prop="authorityId">
      <el-input
        v-model="form.authorityId"
        :disabled="dialogType === 'edit'"
        autocomplete="off"
        maxlength="15"
        :placeholder="t('inputRoleId')"
      />
    </el-form-item>
    <el-form-item :label="t('roleName')" prop="authorityName">
      <el-input
        v-model="form.authorityName"
        autocomplete="off"
        :placeholder="t('inputRoleName')"
      />
    </el-form-item>
  </BaseFormDrawer>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDrawer from '@/components/base/BaseFormDrawer.vue'
import type { Translator } from '@/types/consoleResource'
import type {
  AuthorityDialogType,
  AuthorityForm,
  AuthorityOption
} from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    authorityOptions: AuthorityOption[]
    dialogTitle?: string
    dialogType?: AuthorityDialogType
    form: AuthorityForm
    modelValue?: boolean
    rules: FormRules<AuthorityForm>
    submitting?: boolean
  }>(),
  {
    dialogTitle: '',
    dialogType: 'add',
    modelValue: false,
    submitting: false
  }
)

const emit = defineEmits<{
  close: []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const cascaderProps = {
  checkStrictly: true,
  label: 'authorityName',
  value: 'authorityId',
  disabled: 'disabled',
  emitPath: false
}
</script>

<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :submit-text="t('save')"
    :submitting="editing"
    :title="`${t('edit')}${t('storage')}`"
    width="450px"
    form-class="py-2"
    label-position="top"
    @submit="emit('submit')"
  >
    <el-form-item :label="t('name')" prop="name">
      <el-input v-model="form.name" :placeholder="t('inputName')" />
    </el-form-item>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { StorageEditForm } from '@/types/storage'

const props = withDefaults(
  defineProps<{
    editing?: boolean
    form: StorageEditForm
    modelValue?: boolean
  }>(),
  {
    editing: false,
    modelValue: false
  }
)

const emit = defineEmits<{
  submit: []
  'update:modelValue': [value: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const rules: FormRules<StorageEditForm> = {
  name: [{ required: true, message: t('inputName'), trigger: 'blur' }]
}
</script>

<template>
  <BaseFormDrawer
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :size="500"
    :submit-text="t('confirm')"
    :title="dialogTitle"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <template #prepend>
      <div
        class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3 mb-6"
      >
        <span class="material-icons text-amber-500">info</span>
        <span class="text-sm text-amber-700 dark:text-amber-400">{{
          t('menuTip')
        }}</span>
      </div>
    </template>

    <el-form-item :label="t('path')" prop="path">
      <el-input v-model="form.path" :placeholder="t('apiPath')" />
    </el-form-item>
    <el-form-item :label="t('method')" prop="method">
      <el-select
        v-model="form.method"
        :placeholder="t('select')"
        style="width: 100%"
      >
        <el-option
          v-for="item in methodOptions"
          :key="item.value"
          :label="`${item.label}(${item.value})`"
          :value="item.value"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('apiGroup')" prop="apiGroup">
      <el-select
        v-model="form.apiGroup"
        :placeholder="t('apiGroup')"
        allow-create
        filterable
      >
        <el-option
          v-for="item in apiGroupOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('apiDesc')" prop="description">
      <el-input v-model="form.description" :placeholder="t('apiDesc')" />
    </el-form-item>
  </BaseFormDrawer>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDrawer from '@/components/base/BaseFormDrawer.vue'
import type { Translator } from '@/types/consoleResource'
import type {
  ApiForm,
  ApiMethodOption,
  LabelValueOption
} from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    apiGroupOptions?: LabelValueOption[]
    dialogTitle?: string
    form: ApiForm
    methodOptions?: ApiMethodOption[]
    modelValue?: boolean
    rules: FormRules<ApiForm>
  }>(),
  {
    apiGroupOptions: () => [],
    dialogTitle: '',
    methodOptions: () => [],
    modelValue: false
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
</script>

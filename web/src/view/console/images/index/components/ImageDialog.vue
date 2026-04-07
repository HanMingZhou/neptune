<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :align-center="true"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :shell-size="'md'"
    :submit-text="isEdit ? t('save') : t('create')"
    :submitting="submitting"
    :title="isEdit ? t('imageEdit') : t('imageAdd')"
    label-position="top"
    @closed="emit('closed')"
    @submit="emit('submit')"
  >
    <div class="space-y-1">
      <div
        class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2"
      >
        {{ t('basicInfo') }}
      </div>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <el-form-item :label="t('name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('inputName')" />
        </el-form-item>
        <el-form-item :label="t('imageArea')">
          <el-input v-model="form.area" :placeholder="t('inputArea')" />
        </el-form-item>
        <el-form-item :label="t('imageType')">
          <el-select v-model="form.type">
            <el-option :label="t('systemImage')" :value="1" />
            <el-option :label="t('customImage')" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('imageUsageType')" prop="usageType">
          <el-select v-model="form.usageType">
            <el-option :label="t('usageNotebook')" :value="1" />
            <el-option :label="t('usageTrain')" :value="2" />
            <el-option :label="t('usageInfer')" :value="3" />
          </el-select>
        </el-form-item>
      </div>

      <div
        class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark"
      >
        {{ t('resourceConfig') }}
      </div>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <el-form-item :label="t('imageAddr')" class="md:col-span-2">
          <el-input
            v-model="form.imageAddr"
            placeholder="registry.example.com/image:tag"
          />
        </el-form-item>
        <el-form-item :label="t('imagePath')">
          <el-input v-model="form.imagePath" placeholder="/path/to/image" />
        </el-form-item>
        <el-form-item :label="t('imageSize')">
          <el-input v-model="form.size" placeholder="e.g. 2.5GB" />
        </el-form-item>
      </div>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { ImageForm } from '@/types/image'

const props = withDefaults(
  defineProps<{
    form: ImageForm
    isEdit?: boolean
    modelValue?: boolean
    rules?: FormRules<ImageForm>
    submitting?: boolean
  }>(),
  {
    isEdit: false,
    modelValue: false,
    rules: () => ({}),
    submitting: false
  }
)

const emit = defineEmits<{
  closed: []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

<template>
  <BaseFormDialog
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :shell-size="'xl'"
    :submit-text="isEdit ? t('save') : t('create')"
    :submitting="submitting"
    :title="dialogTitle"
    label-position="top"
    top="4vh"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <div class="space-y-1">
      <div
        class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2"
      >
        {{ t('basicInfo') }}
      </div>
      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <el-form-item :label="t('clusterName')" prop="name">
          <el-input v-model="form.name" :placeholder="t('inputName')" />
        </el-form-item>
        <el-form-item :label="t('area')" prop="area">
          <el-input v-model="form.area" :placeholder="t('inputArea')" />
        </el-form-item>
        <el-form-item :label="t('desc')" class="lg:col-span-2">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            :placeholder="t('inputDesc')"
          />
        </el-form-item>
      </div>

      <div
        class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark"
      >
        {{ t('resourceConfig') }}
      </div>
      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <el-form-item :label="t('clusterApiServer')" class="lg:col-span-2">
          <el-input
            v-model="form.apiServer"
            placeholder="https://x.x.x.x:6443"
          />
        </el-form-item>
        <el-form-item :label="t('clusterHarbor')">
          <el-input
            v-model="form.harborAddr"
            placeholder="harbor.example.com"
          />
        </el-form-item>
        <el-form-item :label="t('storageClass')">
          <el-input
            v-model="form.storageClass"
            :placeholder="t('inputStorageClass')"
          />
        </el-form-item>
        <el-form-item :label="t('clusterKubeconfig')" class="lg:col-span-2">
          <el-input
            v-model="form.kubeconfig"
            type="textarea"
            :rows="10"
            placeholder="kubeconfig YAML"
            class="font-mono"
          />
        </el-form-item>
      </div>

      <div
        class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark"
      >
        {{ t('status') }}
      </div>
      <el-form-item>
        <el-switch
          v-model="form.status"
          :active-value="1"
          :inactive-value="0"
          :active-text="t('enable')"
          :inactive-text="t('disable')"
        />
      </el-form-item>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { CmsClusterForm } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    dialogTitle?: string
    form: CmsClusterForm
    isEdit?: boolean
    modelValue?: boolean
    rules: FormRules<CmsClusterForm>
    submitting?: boolean
  }>(),
  {
    dialogTitle: '',
    isEdit: false,
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
</script>

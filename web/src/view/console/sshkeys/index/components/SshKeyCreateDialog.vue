<template>
  <BaseFormDialog
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :shell-size="'lg'"
    :submit-text="t('add')"
    :submitting="loading"
    :title="t('newSshKey')"
    label-position="top"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <el-form-item :label="t('sshKeyName')" prop="name">
      <el-input
        v-model="form.name"
        :placeholder="t('sshKeyNamePlaceholder')"
        maxlength="50"
      />
      <div class="text-xs text-slate-400 mt-1">
        {{ t('sshKeyNamePlaceholder') }}
      </div>
    </el-form-item>

    <el-form-item :label="t('publicKeyContent')" prop="publicKey">
      <el-input
        v-model="form.publicKey"
        type="textarea"
        :rows="8"
        :placeholder="t('publicKeyPlaceholder')"
        class="ssh-key-create-dialog__textarea"
      />
      <div class="text-xs text-slate-400 mt-1">
        {{ t('publicKeyPlaceholder') }}
      </div>
    </el-form-item>

    <template #append>
      <div class="bg-slate-50 dark:bg-zinc-800 rounded-lg p-4 mt-4 text-xs">
        <h4 class="font-bold text-sm mb-2 flex items-center gap-2">
          <span class="material-icons text-[16px]">tips_and_updates</span>
          {{ t('sshKeyHintTitle') }}
        </h4>
        <ol
          class="text-slate-600 dark:text-slate-400 list-decimal list-inside space-y-1"
        >
          <li>{{ t('sshKeyHintStep1') }}</li>
          <li>{{ t('sshKeyHintStep2') }}</li>
          <li>{{ t('sshKeyHintStep3') }}</li>
        </ol>
      </div>
    </template>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { SshKeyCreateForm } from '@/types/sshkey'

const props = withDefaults(
  defineProps<{
    form: SshKeyCreateForm
    loading?: boolean
    modelValue?: boolean
    rules: FormRules<SshKeyCreateForm>
  }>(),
  {
    loading: false,
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

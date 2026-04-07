<template>
  <BaseFormDialog
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :model="form"
    :submit-text="t('confirm')"
    :title="t('resetPassword')"
    width="500px"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <el-form-item :label="t('username')">
      <el-input v-model="form.userName" disabled />
    </el-form-item>
    <el-form-item :label="t('nickname')">
      <el-input v-model="form.nickName" disabled />
    </el-form-item>
    <el-form-item :label="t('password')">
      <div class="flex w-full gap-3">
        <el-input
          v-model="form.password"
          class="flex-1"
          :placeholder="t('password')"
          show-password
        />
        <el-button type="primary" @click="$emit('generate-password')">
          {{ t('refresh') }}
        </el-button>
      </div>
    </el-form-item>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { UserResetPasswordForm } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    form: UserResetPasswordForm
    modelValue?: boolean
  }>(),
  {
    modelValue: false
  }
)

const emit = defineEmits<{
  close: []
  'generate-password': []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    class="storage-expand-dialog"
    :model="form"
    :rules="rules"
    :submit-text="`${t('confirm')}${t('expand')}`"
    :submitting="expanding"
    :title="`${t('expand')}${t('storage')}`"
    width="450px"
    form-class="py-2"
    label-position="top"
    @submit="emit('submit')"
  >
    <div class="space-y-5">
      <div
        class="flex items-center justify-between rounded-xl border border-blue-100 bg-blue-50/50 p-4 dark:border-blue-900/30 dark:bg-blue-900/10"
      >
        <span class="text-sm font-medium text-blue-600 dark:text-blue-400">{{
          t('currentSize')
        }}</span>
        <span
          class="font-mono text-lg font-black text-blue-700 dark:text-blue-300"
          >{{ form.currentSize }}</span
        >
      </div>

      <el-form-item :label="`${t('expandTo')} (GB)`" prop="newSize">
        <el-input-number
          v-model.number="form.newSize"
          :min="form.minSize"
          :max="2000"
          :step="10"
          controls-position="right"
          class="w-full storage-create-dialog__number"
        />
        <p class="mt-2 text-[11px] text-slate-400">
          {{ t('expandHint') || 'Only upward expansion is supported' }}
        </p>
      </el-form-item>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { StorageExpandForm } from '@/types/storage'

const props = withDefaults(
  defineProps<{
    expanding?: boolean
    form: StorageExpandForm
    modelValue?: boolean
  }>(),
  {
    expanding: false,
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

const rules = computed<FormRules<StorageExpandForm>>(() => ({
  newSize: [
    { required: true, message: t('capacityRange'), trigger: 'change' },
    {
      validator: (_rule, value: number, callback) => {
        if (typeof value !== 'number' || value < props.form.minSize) {
          callback(new Error(t('capacityError')))
          return
        }

        if (value > 2000) {
          callback(new Error(t('capacityRange')))
          return
        }

        callback()
      },
      trigger: 'change'
    }
  ]
}))
</script>

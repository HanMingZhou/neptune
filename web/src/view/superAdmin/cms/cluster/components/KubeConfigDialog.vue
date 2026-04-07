<template>
  <BaseDialog
    v-model="visibleModel"
    :title="t('clusterKubeconfig')"
    align-center
    width="700px"
    @close="emit('close')"
  >
    <div
      class="bg-slate-50 dark:bg-zinc-900 border border-border-light dark:border-border-dark rounded-lg p-4 max-h-[500px] overflow-auto"
    >
      <pre
        class="text-xs font-mono text-slate-700 dark:text-slate-300 whitespace-pre-wrap break-all"
        >{{ content || '-' }}</pre
      >
    </div>

    <template #footer="{ requestClose }">
      <div class="flex justify-end gap-3">
        <el-button @click="$emit('copy')">{{ t('copy') }}</el-button>
        <el-button type="primary" @click="requestClose">{{
          t('close')
        }}</el-button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseDialog from '@/components/base/BaseDialog.vue'
import type { Translator } from '@/types/consoleResource'

const props = withDefaults(
  defineProps<{
    content?: string
    modelValue?: boolean
  }>(),
  {
    content: '',
    modelValue: false
  }
)

const emit = defineEmits<{
  close: []
  copy: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

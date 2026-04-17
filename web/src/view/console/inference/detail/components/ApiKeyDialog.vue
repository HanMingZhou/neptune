<template>
  <BaseDialog
    v-model="visibleModel"
    :title="t('inference.keyManage')"
    class="inference-api-key-dialog"
    destroy-on-close
    width="600px"
  >
    <div class="space-y-4">
      <div class="flex flex-col items-stretch gap-3 sm:flex-row">
        <el-input
          class="flex-1"
          :model-value="newKeyName"
          :placeholder="t('inference.inputKeyName')"
          size="large"
          @update:model-value="emit('update:new-key-name', $event)"
        />
        <el-button
          class="shrink-0 sm:self-auto"
          type="primary"
          size="large"
          :disabled="!newKeyName.trim()"
          @click="emit('create')"
        >
          <span class="material-icons mr-1 text-base">add</span>
          {{ t('inference.createKey') }}
        </el-button>
      </div>

      <div
        v-if="newlyCreatedKey"
        class="rounded-lg border border-emerald-200 bg-emerald-50 p-4 dark:border-emerald-800 dark:bg-emerald-900/20"
      >
        <div class="mb-1 text-xs text-emerald-600 dark:text-emerald-400">
          {{ t('inference.newKeyHint') }}
        </div>
        <div class="flex items-center gap-2">
          <code
            class="flex-1 break-all font-mono text-sm text-emerald-700 dark:text-emerald-300"
            >{{ newlyCreatedKey }}</code
          >
          <el-button
            link
            type="primary"
            class="!h-auto !min-h-0 !p-0 text-xs"
            @click="emit('copy', newlyCreatedKey)"
          >
            {{ t('copy') }}
          </el-button>
        </div>
      </div>

      <div v-if="apiKeys.length" class="space-y-3">
        <div
          v-for="key in apiKeys"
          :key="key.id"
          class="flex items-center justify-between rounded-lg border border-border-light p-3 dark:border-border-dark"
        >
          <div class="min-w-0 flex-1">
            <div class="text-sm font-medium">{{ key.name }}</div>
            <div class="mt-1 break-all font-mono text-xs text-slate-400">
              {{ key.apiKey }}
            </div>
            <div class="mt-1 text-xs text-slate-400">
              {{ formatTime(key.createdAt) }}
            </div>
          </div>
          <el-button
            link
            type="danger"
            class="ml-4 shrink-0"
            @click="emit('delete', key)"
          >
            <span class="material-icons mr-1 text-lg">delete</span>
            {{ t('delete') }}
          </el-button>
        </div>
      </div>
      <div v-else class="py-8 text-center text-sm text-slate-400">
        {{ t('inference.noApiKey') }}
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseDialog from '@/components/base/BaseDialog.vue'
import type { InferenceApiKey, Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    apiKeys?: InferenceApiKey[]
    formatTime: (value?: string | number) => string
    newKeyName?: string
    newlyCreatedKey?: string
  }>(),
  {
    apiKeys: () => [],
    newKeyName: '',
    newlyCreatedKey: ''
  }
)

const emit = defineEmits<{
  close: []
  copy: [text: string]
  create: []
  delete: [apiKey: InferenceApiKey]
  'update:new-key-name': [value: string]
}>()

const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => true,
  set: (value: boolean) => {
    if (!value) {
      emit('close')
    }
  }
})
</script>

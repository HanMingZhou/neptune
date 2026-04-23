<template>
  <BaseDialog
    v-model="visibleModel"
    :title="t('inference.keyManage')"
    class="inference-api-key-dialog"
    destroy-on-close
    width="680px"
  >
    <div class="space-y-5">
      <section
        class="rounded-2xl border border-slate-200 bg-gradient-to-br from-slate-50 via-white to-blue-50/70 p-5 shadow-sm dark:border-slate-800 dark:from-slate-950 dark:via-slate-900 dark:to-blue-950/30"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="space-y-1">
            <h3 class="text-base font-semibold text-slate-900 dark:text-slate-100">
              {{ t('inference.createKey') }}
            </h3>
            <p class="text-sm text-slate-500 dark:text-slate-400">
              {{ t('inference.newKeyHint') }}
            </p>
          </div>
          <div class="rounded-full bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-600 dark:text-blue-300">
            API Key
          </div>
        </div>

        <div class="mt-4 flex flex-col items-stretch gap-3 lg:flex-row">
          <el-input
            class="flex-1"
            :model-value="newKeyName"
            :placeholder="t('inference.inputKeyName')"
            size="large"
            @keyup.enter="emit('create')"
            @update:model-value="emit('update:new-key-name', $event)"
          />
          <el-button
            class="shrink-0 lg:min-w-[136px]"
            type="primary"
            size="large"
            :loading="creatingApiKey"
            :disabled="!trimmedKeyName"
            @click="emit('create')"
          >
            <span class="material-icons mr-1 text-base">add</span>
            {{ t('inference.createKey') }}
          </el-button>
        </div>
      </section>

      <section
        v-if="newlyCreatedKey"
        class="rounded-2xl border border-emerald-200 bg-gradient-to-br from-emerald-50 via-white to-emerald-100/70 p-5 shadow-sm dark:border-emerald-800 dark:from-emerald-950/40 dark:via-slate-900 dark:to-emerald-900/30"
      >
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <div class="text-sm font-semibold text-emerald-700 dark:text-emerald-300">
              {{ t('success') }}
            </div>
            <div class="mt-1 text-xs leading-5 text-emerald-700/80 dark:text-emerald-300/80">
              {{ t('inference.newKeyHint') }}
            </div>
          </div>
          <el-button
            type="success"
            plain
            @click="emit('copy', newlyCreatedKey)"
          >
            <span class="material-icons mr-1 text-[18px]">content_copy</span>
            {{ t('copy') }}
          </el-button>
        </div>
        <div class="mt-4 rounded-xl bg-slate-950 px-4 py-3 shadow-inner">
          <code class="block break-all font-mono text-sm leading-6 text-emerald-300">
            {{ newlyCreatedKey }}
          </code>
        </div>
      </section>

      <section
        class="rounded-2xl border border-border-light bg-white/90 p-4 shadow-sm backdrop-blur-sm dark:border-border-dark dark:bg-slate-900/80"
      >
        <div class="mb-4 flex items-center justify-between gap-3">
          <div>
            <h3 class="text-sm font-semibold text-slate-900 dark:text-slate-100">
              {{ t('inference.keyManage') }}
            </h3>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
              {{ apiKeys.length }}
            </p>
          </div>
        </div>

        <div v-loading="apiKeysLoading" class="min-h-[180px]">
          <div v-if="apiKeys.length" class="space-y-3">
            <article
              v-for="key in apiKeys"
              :key="key.id"
              class="rounded-xl border border-slate-200 bg-slate-50/80 p-4 transition-all hover:border-slate-300 hover:shadow-sm dark:border-slate-800 dark:bg-slate-950/40 dark:hover:border-slate-700"
            >
              <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                <div class="min-w-0 flex-1">
                  <div class="flex flex-wrap items-center gap-2">
                    <h4 class="max-w-full truncate text-sm font-semibold text-slate-900 dark:text-slate-100">
                      {{ key.name || '-' }}
                    </h4>
                    <span class="rounded-full bg-slate-200 px-2 py-0.5 text-[11px] font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                      API Key
                    </span>
                  </div>
                  <div class="mt-2 rounded-lg bg-white px-3 py-2 font-mono text-xs text-slate-500 shadow-inner dark:bg-slate-900 dark:text-slate-400">
                    {{ key.apiKey || '-' }}
                  </div>
                  <div class="mt-2 text-xs text-slate-400">
                    {{ formatTime(key.createdAt) }}
                  </div>
                </div>

                <div class="flex shrink-0 items-center justify-end gap-2">
                  <el-button
                    type="danger"
                    plain
                    :loading="deletingApiKeyId === Number(key.id)"
                    @click="emit('delete', key)"
                  >
                    <span class="material-icons mr-1 text-[18px]">delete</span>
                    {{ t('delete') }}
                  </el-button>
                </div>
              </div>
            </article>
          </div>
          <div v-else class="flex min-h-[180px] items-center justify-center rounded-xl border border-dashed border-slate-200 bg-slate-50/60 text-sm text-slate-400 dark:border-slate-800 dark:bg-slate-950/30">
            {{ t('inference.noApiKey') }}
          </div>
        </div>
      </section>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseDialog from '@/components/base/BaseDialog.vue'
import type { InferenceApiKey, Translator } from '@/types/consoleResource'

const props = withDefaults(
  defineProps<{
    apiKeys?: InferenceApiKey[]
    apiKeysLoading?: boolean
    creatingApiKey?: boolean
    deletingApiKeyId?: number | null
    formatTime: (value?: string | number) => string
    newKeyName?: string
    newlyCreatedKey?: string
  }>(),
  {
    apiKeys: () => [],
    apiKeysLoading: false,
    creatingApiKey: false,
    deletingApiKeyId: null,
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

const trimmedKeyName = computed(() => props.newKeyName.trim())

const visibleModel = computed({
  get: () => true,
  set: (value: boolean) => {
    if (!value) {
      emit('close')
    }
  }
})
</script>

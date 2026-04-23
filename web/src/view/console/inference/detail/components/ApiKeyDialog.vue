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

        <div class="api-key-create-row mt-4 flex flex-col items-stretch gap-3 lg:flex-row">
          <label class="api-key-create-input-shell flex-1">
            <span class="material-icons api-key-create-input-icon">key</span>
            <input
              class="api-key-create-native-input"
              :value="newKeyName"
              :placeholder="t('inference.inputKeyName')"
              type="text"
              @input="handleKeyNameInput"
              @keyup.enter="emit('create')"
            />
          </label>
          <button
            class="api-key-create-native-button lg:min-w-[136px]"
            type="button"
            :disabled="!trimmedKeyName || creatingApiKey"
            @click="emit('create')"
          >
            <span :class="['material-icons', creatingApiKey ? 'api-key-create-spin' : '']">{{
              creatingApiKey ? 'progress_activity' : 'add'
            }}</span>
            <span class="api-key-create-native-button__text">{{ t('inference.createKey') }}</span>
          </button>
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

const handleKeyNameInput = (event: Event) => {
  emit('update:new-key-name', (event.target as HTMLInputElement).value)
}

const visibleModel = computed({
  get: () => true,
  set: (value: boolean) => {
    if (!value) {
      emit('close')
    }
  }
})
</script>

<style scoped>
.api-key-create-row {
  align-items: stretch;
}

.api-key-create-input-shell {
  display: flex;
  align-items: center;
  gap: 0.7rem;
  min-width: 0;
  height: 46px;
  min-height: 46px;
  box-sizing: border-box;
  border-radius: 15px;
  border: 1px solid rgb(203 213 225 / 0.95);
  background: linear-gradient(180deg, rgb(255 255 255 / 0.98), rgb(248 250 252 / 0.96));
  padding: 0 1rem;
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.95),
    0 10px 24px -18px rgb(15 23 42 / 0.45);
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    background 0.18s ease,
    transform 0.18s ease;
}

.api-key-create-input-shell:focus-within {
  border-color: rgb(49 88 212);
  background: linear-gradient(180deg, rgb(255 255 255), rgb(241 245 249));
  box-shadow:
    0 0 0 3px rgb(49 88 212 / 0.12),
    0 14px 26px -20px rgb(49 88 212 / 0.45);
}

.api-key-create-input-icon {
  flex-shrink: 0;
  font-size: 18px;
  color: rgb(71 85 105);
}

.api-key-create-native-input {
  flex: 1;
  min-width: 0;
  height: 100%;
  border: none;
  background: transparent;
  padding: 0;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(15 23 42);
  outline: none;
}

.api-key-create-native-input::placeholder {
  color: rgb(148 163 184);
}

.api-key-create-native-button {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  height: 46px;
  min-height: 46px;
  border-radius: 15px;
  border: 1px solid rgb(49 88 212 / 0.96);
  background: linear-gradient(135deg, rgb(49 88 212), rgb(37 99 235));
  padding: 0 1.1rem;
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: rgb(255 255 255);
  white-space: nowrap;
  flex-shrink: 0;
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.18),
    0 14px 28px -18px rgb(37 99 235 / 0.75);
  transition:
    background 0.18s ease,
    border-color 0.18s ease,
    opacity 0.18s ease,
    transform 0.18s ease,
    box-shadow 0.18s ease;
}

.api-key-create-native-button::before {
  content: '';
  position: absolute;
  inset: 1px;
  z-index: 0;
  border-radius: 14px;
  background: linear-gradient(180deg, rgb(255 255 255 / 0.18), rgb(255 255 255 / 0.02));
  opacity: 0.9;
}

.api-key-create-native-button > * {
  position: relative;
  z-index: 1;
}

.api-key-create-native-button:hover:not(:disabled) {
  border-color: rgb(37 69 179);
  background: linear-gradient(135deg, rgb(37 69 179), rgb(29 78 216));
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.2),
    0 16px 30px -18px rgb(29 78 216 / 0.78);
  transform: translateY(-1px);
}

.api-key-create-native-button:active:not(:disabled) {
  transform: translateY(0);
}

.api-key-create-native-button:disabled {
  cursor: not-allowed;
  border-color: rgb(203 213 225);
  background: linear-gradient(180deg, rgb(226 232 240), rgb(226 232 240));
  color: rgb(100 116 139);
  box-shadow: none;
}

.api-key-create-native-button:disabled::before {
  opacity: 0;
}

.api-key-create-native-button .material-icons {
  line-height: 1;
  font-size: 18px;
}

.api-key-create-spin {
  animation: api-key-create-spin 1s linear infinite;
}

.dark .api-key-create-native-input {
  color: rgb(226 232 240);
}

.dark .api-key-create-input-shell {
  border-color: rgb(51 65 85 / 0.96);
  background: linear-gradient(180deg, rgb(15 23 42 / 0.96), rgb(15 23 42 / 0.82));
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.04),
    0 14px 28px -22px rgb(2 6 23 / 0.85);
}

.dark .api-key-create-input-shell:focus-within {
  border-color: rgb(96 165 250);
  background: linear-gradient(180deg, rgb(15 23 42), rgb(15 23 42 / 0.88));
  box-shadow:
    0 0 0 3px rgb(96 165 250 / 0.16),
    0 16px 32px -24px rgb(37 99 235 / 0.7);
}

.dark .api-key-create-input-icon {
  color: rgb(148 163 184);
}

.dark .api-key-create-native-input::placeholder {
  color: rgb(100 116 139);
}

.dark .api-key-create-native-button {
  border-color: rgb(59 130 246);
  background: linear-gradient(135deg, rgb(37 99 235), rgb(29 78 216));
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 0.12),
    0 16px 30px -18px rgb(29 78 216 / 0.8);
}

.dark .api-key-create-native-button:hover:not(:disabled) {
  border-color: rgb(29 78 216);
  background: linear-gradient(135deg, rgb(29 78 216), rgb(30 64 175));
}

.dark .api-key-create-native-button:disabled {
  border-color: rgb(51 65 85);
  background: linear-gradient(180deg, rgb(30 41 59), rgb(30 41 59));
  color: rgb(148 163 184);
}

@keyframes api-key-create-spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}
</style>

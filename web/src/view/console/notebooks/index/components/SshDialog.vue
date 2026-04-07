<template>
  <BaseDialog
    :model-value="visible"
    :title="t('sshLogin')"
    class="rounded-xl overflow-hidden ssh-dialog"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <div class="space-y-6">
      <div v-if="currentNotebook.sshKeyCommand" class="space-y-3">
        <h4 class="font-bold text-sm flex items-center gap-2">
          <span class="material-icons text-primary text-[18px]">vpn_key</span>
          {{ t('sshLogin') }} ({{ t('keyAuth') }})
        </h4>
        <div
          class="bg-zinc-900 rounded-lg p-4 group relative border border-zinc-800"
        >
          <pre
            class="text-emerald-400 text-xs font-mono overflow-x-auto custom-scrollbar"
            >{{ currentNotebook.sshKeyCommand }}</pre
          >
          <el-button
            circle
            text
            class="ssh-dialog__copy-button"
            @click="$emit('copy', currentNotebook.sshKeyCommand)"
          >
            <span class="material-icons text-[14px]">content_copy</span>
          </el-button>
        </div>
      </div>

      <div v-if="currentNotebook.sshCommand" class="space-y-3">
        <h4 class="font-bold text-sm flex items-center gap-2">
          <span class="material-icons text-amber-500 text-[18px]"
            >password</span
          >
          {{ t('sshLogin') }} ({{ t('passwordAuth') }})
        </h4>
        <div
          class="bg-zinc-900 rounded-lg p-4 group relative border border-zinc-800"
        >
          <pre
            class="text-emerald-400 text-xs font-mono overflow-x-auto custom-scrollbar"
            >{{ currentNotebook.sshCommand }}</pre
          >
          <el-button
            circle
            text
            class="ssh-dialog__copy-button"
            @click="$emit('copy', currentNotebook.sshCommand)"
          >
            <span class="material-icons text-[14px]">content_copy</span>
          </el-button>
        </div>
        <div
          class="flex items-center gap-4 bg-gray-50 dark:bg-zinc-900 p-3 rounded-lg border border-gray-100 dark:border-border-dark"
        >
          <span class="text-xs text-slate-500">{{ t('password') }}:</span>
          <span
            class="font-mono bg-white dark:bg-zinc-800 border border-gray-200 dark:border-border-dark px-2 py-1 rounded text-sm text-slate-700 dark:text-slate-300"
          >
            {{ showPassword ? currentNotebook.sshPassword : '••••••••' }}
          </span>
          <el-button
            link
            type="primary"
            class="!h-auto !min-h-0 !px-0 text-xs font-bold"
            @click="$emit('update:show-password', !showPassword)"
          >
            {{ showPassword ? t('hide') : t('show') }}
          </el-button>
          <el-button
            link
            type="primary"
            class="ml-auto !h-auto !min-h-0 text-xs font-bold"
            @click="$emit('copy', currentNotebook.sshPassword)"
          >
            <span class="material-icons text-[14px]">content_copy</span>
            {{ t('copy') }}
          </el-button>
        </div>
      </div>

      <div
        v-if="!currentNotebook.sshKeyCommand && !currentNotebook.sshCommand"
        class="text-center py-12 text-slate-400 bg-gray-50 dark:bg-zinc-900 rounded-xl border border-dashed border-gray-200 dark:border-border-dark"
      >
        <span class="material-icons text-4xl mb-2 opacity-50">block</span>
        <p>{{ t('noData') }}</p>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import BaseDialog from '@/components/base/BaseDialog.vue'
import type { ConsoleNotebook, Translator } from '@/types/consoleResource'

withDefaults(
  defineProps<{
    currentNotebook?: ConsoleNotebook
    showPassword?: boolean
    visible?: boolean
  }>(),
  {
    currentNotebook: () => ({}),
    showPassword: false,
    visible: false
  }
)

defineEmits<{
  copy: [value: string | undefined]
  'update:show-password': [value: boolean]
  'update:visible': [value: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)
</script>

<style scoped>
.ssh-dialog__copy-button {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  color: rgb(255 255 255 / 0.9);
  background: rgb(255 255 255 / 0.12);
  border: 1px solid rgb(255 255 255 / 0.14);
  opacity: 0;
  transition:
    opacity 0.2s ease,
    background-color 0.2s ease;
}

.group:hover .ssh-dialog__copy-button {
  opacity: 1;
}

.ssh-dialog__copy-button:hover,
.ssh-dialog__copy-button:focus-visible {
  color: rgb(255 255 255 / 0.95);
  background: rgb(255 255 255 / 0.18);
}
</style>

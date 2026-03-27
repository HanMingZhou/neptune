<template>
  <el-dialog
    :model-value="visible"
    :title="t('sshLogin')"
    class="rounded-xl overflow-hidden"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <div class="space-y-6">
      <div v-if="currentNotebook.sshKeyCommand" class="space-y-3">
        <h4 class="font-bold text-sm flex items-center gap-2">
          <span class="material-icons text-primary text-[18px]">vpn_key</span>
          {{ t('sshLogin') }} ({{ t('keyAuth') }})
        </h4>
        <div class="bg-zinc-900 rounded-lg p-4 group relative border border-zinc-800">
          <pre class="text-emerald-400 text-xs font-mono overflow-x-auto custom-scrollbar">{{ currentNotebook.sshKeyCommand }}</pre>
          <button class="absolute top-2 right-2 p-1.5 bg-white/10 hover:bg-white/20 rounded text-white opacity-0 group-hover:opacity-100 transition-opacity" @click="$emit('copy', currentNotebook.sshKeyCommand)">
            <span class="material-icons text-[14px]">content_copy</span>
          </button>
        </div>
      </div>

      <div v-if="currentNotebook.sshCommand" class="space-y-3">
        <h4 class="font-bold text-sm flex items-center gap-2">
          <span class="material-icons text-amber-500 text-[18px]">password</span>
          {{ t('sshLogin') }} ({{ t('passwordAuth') }})
        </h4>
        <div class="bg-zinc-900 rounded-lg p-4 group relative border border-zinc-800">
          <pre class="text-emerald-400 text-xs font-mono overflow-x-auto custom-scrollbar">{{ currentNotebook.sshCommand }}</pre>
          <button class="absolute top-2 right-2 p-1.5 bg-white/10 hover:bg-white/20 rounded text-white opacity-0 group-hover:opacity-100 transition-opacity" @click="$emit('copy', currentNotebook.sshCommand)">
            <span class="material-icons text-[14px]">content_copy</span>
          </button>
        </div>
        <div class="flex items-center gap-4 bg-gray-50 dark:bg-zinc-900 p-3 rounded-lg border border-gray-100 dark:border-border-dark">
          <span class="text-xs text-slate-500">{{ t('password') }}:</span>
          <span class="font-mono bg-white dark:bg-zinc-800 border border-gray-200 dark:border-border-dark px-2 py-1 rounded text-sm text-slate-700 dark:text-slate-300">
            {{ showPassword ? currentNotebook.sshPassword : '••••••••' }}
          </span>
          <button class="text-xs font-bold text-primary hover:underline" @click="$emit('update:show-password', !showPassword)">
            {{ showPassword ? t('hide') : t('show') }}
          </button>
          <button class="text-xs font-bold text-slate-500 hover:text-slate-700 ml-auto flex items-center gap-1" @click="$emit('copy', currentNotebook.sshPassword)">
            <span class="material-icons text-[14px]">content_copy</span> {{ t('copy') }}
          </button>
        </div>
      </div>

      <div v-if="!currentNotebook.sshKeyCommand && !currentNotebook.sshCommand" class="text-center py-12 text-slate-400 bg-gray-50 dark:bg-zinc-900 rounded-xl border border-dashed border-gray-200 dark:border-border-dark">
        <span class="material-icons text-4xl mb-2 opacity-50">block</span>
        <p>{{ t('noData') }}</p>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  currentNotebook: {
    type: Object,
    default: () => ({})
  },
  showPassword: {
    type: Boolean,
    default: false
  },
  visible: {
    type: Boolean,
    default: false
  }
})

defineEmits(['copy', 'update:show-password', 'update:visible'])

const t = inject('t', (key) => key)
</script>

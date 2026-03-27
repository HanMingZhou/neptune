<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="$emit('close')">
    <div class="bg-surface-light dark:bg-surface-dark rounded-xl shadow-2xl w-[600px] max-h-[80vh] flex flex-col">
      <div class="px-6 py-4 border-b border-border-light dark:border-border-dark flex items-center justify-between">
        <h3 class="font-bold">{{ t('inference.keyManage') }}</h3>
        <button class="text-slate-400 hover:text-slate-600" @click="$emit('close')">
          <span class="material-icons">close</span>
        </button>
      </div>
      <div class="p-6 flex-1 overflow-y-auto">
        <div class="flex gap-3 mb-4">
          <input
            :value="newKeyName"
            :placeholder="t('inference.inputKeyName')"
            class="flex-1 px-3 py-2 border border-border-light dark:border-border-dark rounded-lg text-sm bg-white dark:bg-zinc-800"
            @input="$emit('update:new-key-name', $event.target.value)"
          />
          <button
            :disabled="!newKeyName.trim()"
            class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-bold disabled:opacity-50 flex items-center gap-1"
            @click="$emit('create')"
          >
            <span class="material-icons text-lg">add</span>
            {{ t('inference.createKey') }}
          </button>
        </div>

        <div v-if="newlyCreatedKey" class="mb-4 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-lg p-4">
          <div class="text-xs text-emerald-600 dark:text-emerald-400 mb-1">{{ t('inference.newKeyHint') }}</div>
          <div class="flex items-center gap-2">
            <code class="flex-1 text-sm font-mono text-emerald-700 dark:text-emerald-300 break-all">{{ newlyCreatedKey }}</code>
            <button class="text-primary text-xs hover:underline" @click="$emit('copy', newlyCreatedKey)">{{ t('copy') }}</button>
          </div>
        </div>

        <div v-if="apiKeys.length" class="space-y-3">
          <div v-for="key in apiKeys" :key="key.id" class="flex items-center justify-between p-3 border border-border-light dark:border-border-dark rounded-lg">
            <div>
              <div class="text-sm font-medium">{{ key.name }}</div>
              <div class="text-xs text-slate-400 font-mono mt-1">{{ key.apiKey }}</div>
              <div class="text-xs text-slate-400 mt-1">{{ formatTime(key.createdAt) }}</div>
            </div>
            <button class="text-red-500 hover:text-red-600" @click="$emit('delete', key)">
              <span class="material-icons text-lg">delete</span>
            </button>
          </div>
        </div>
        <div v-else class="text-sm text-slate-400 text-center py-8">{{ t('inference.noApiKey') }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  apiKeys: {
    type: Array,
    default: () => []
  },
  formatTime: {
    type: Function,
    required: true
  },
  newKeyName: {
    type: String,
    default: ''
  },
  newlyCreatedKey: {
    type: String,
    default: ''
  }
})

defineEmits(['close', 'copy', 'create', 'delete', 'update:new-key-name'])

const t = inject('t', (key) => key)
</script>

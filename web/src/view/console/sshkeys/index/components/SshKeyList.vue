<template>
  <div v-if="loading" class="py-20 text-center text-slate-400">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
    {{ t('loading') }}
  </div>

  <div
    v-else-if="items.length === 0"
    class="bg-surface-light dark:bg-surface-dark border border-dashed border-border-light dark:border-border-dark rounded-xl py-16 text-center"
  >
    <span class="material-icons text-5xl text-slate-300 mb-4">vpn_key</span>
    <h3 class="text-lg font-bold mb-2">{{ t('noSshKeyData') }}</h3>
    <p class="text-slate-500 mb-6">{{ t('sshKeyManageDesc') }}</p>
    <button
      class="bg-primary hover:bg-primary-hover text-white px-5 py-2.5 rounded-lg font-bold text-sm"
      @click="emit('create')"
    >
      {{ t('newSshKey') }}
    </button>
  </div>

  <div v-else class="grid gap-4">
    <div
      v-for="key in items"
      :key="key.id"
      class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-5 hover:border-primary transition-colors"
    >
      <div class="flex justify-between items-start mb-3 gap-4">
        <div class="min-w-0">
          <h3 class="font-bold text-lg flex items-center gap-2 flex-wrap">
            <span>{{ key.name }}</span>
            <span v-if="key.isDefault" class="px-2 py-0.5 bg-primary/10 text-primary text-xs font-bold rounded-full">
              {{ t('isDefault') }}
            </span>
          </h3>
          <p class="font-mono text-sm text-slate-500 break-all">{{ key.fingerprint }}</p>
        </div>
        <div class="flex gap-2 shrink-0">
          <button
            v-if="!key.isDefault"
            class="text-sm font-medium text-primary hover:text-primary-hover"
            @click="emit('set-default', key)"
          >
            {{ t('setDefault') }}
          </button>
          <button
            class="text-sm font-medium text-red-500 hover:text-red-600"
            @click="emit('delete', key)"
          >
            {{ t('delete') }}
          </button>
        </div>
      </div>
      <div class="pt-3 border-t border-border-light dark:border-border-dark flex items-center gap-2 text-sm text-slate-400">
        <span class="material-icons text-[16px]">calendar_today</span>
        {{ t('createdAt') }} {{ key.createdAt }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  items: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['create', 'delete', 'set-default'])
const t = inject('t', (key) => key)
</script>

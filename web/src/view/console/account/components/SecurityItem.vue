<template>
  <div class="flex items-start justify-between p-6 border-b border-border-light dark:border-border-dark last:border-0 hover:bg-slate-50/50 dark:hover:bg-zinc-800/30 transition-colors">
    <div class="flex gap-4">
      <div class="size-10 rounded-lg bg-slate-100 dark:bg-zinc-800 flex items-center justify-center text-slate-500">
        <span class="material-icons">{{ icon }}</span>
      </div>
      <div>
        <h3 class="font-bold text-sm mb-1">{{ title }}</h3>
        <p class="text-xs text-slate-500 leading-relaxed max-w-md">{{ desc }}</p>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <span v-if="statusText" :class="`text-[11px] font-black uppercase px-2 py-0.5 rounded ${statusColor}`">{{ statusText }}</span>
      <button 
        v-if="actionText"
        class="text-primary text-xs font-bold hover:underline disabled:opacity-50 disabled:no-underline flex items-center gap-1" 
        :disabled="loading || disabled" 
        @click="$emit('click')"
      >
        <div v-if="loading" class="size-3 border-2 border-primary/30 border-t-primary rounded-full animate-spin"></div>
        {{ actionText }}
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  icon: String,
  title: String,
  desc: String,
  statusText: String,
  statusColor: {
    type: String,
    default: 'text-slate-500'
  },
  actionText: String,
  loading: Boolean,
  disabled: Boolean
})

defineEmits(['click'])
</script>

<template>
  <div class="console-page-container px-6 pb-4">
    <div class="rounded-xl border border-border-light bg-surface-light px-2 dark:border-border-dark dark:bg-surface-dark md:px-3">
      <div class="flex gap-1 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="[
            'inline-flex shrink-0 items-center gap-1.5 border-b-2 px-3.5 py-3 text-sm transition-colors md:px-4',
            activeTab === tab.key
              ? 'border-primary font-semibold text-primary'
              : 'border-transparent text-slate-500 hover:text-primary'
          ]"
          @click="$emit('update:activeTab', tab.key)"
        >
          <span class="material-icons text-base md:text-lg">{{ tab.icon }}</span>
          {{ tab.label }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  activeTab: {
    type: String,
    required: true
  },
  tabs: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:activeTab'])
</script>

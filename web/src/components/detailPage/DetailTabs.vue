<template>
  <div class="bg-surface-light dark:bg-surface-dark border-b border-border-light dark:border-border-dark px-6 flex gap-1">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      :class="[
        'flex items-center gap-2 px-5 py-3.5 text-sm border-b-2 transition-colors',
        activeTab === tab.key ? 'text-primary border-primary font-bold' : 'text-slate-500 border-transparent hover:text-primary'
      ]"
      @click="$emit('update:activeTab', tab.key)"
    >
      <span class="material-icons text-lg">{{ tab.icon }}</span>
      {{ tab.label }}
    </button>
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

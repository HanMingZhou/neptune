<template>
  <button 
    @click="handleManualRefresh" 
    class="list-toolbar-button list-toolbar-button--secondary list-toolbar-button--icon"
    :disabled="loading"
    :title="t('refresh')"
    type="button"
  >
    <span class="material-icons" :class="{ 'animate-spin': loading }">refresh</span>
  </button>
</template>

<script setup>
import { onMounted, onUnmounted, inject } from 'vue'
import { REFRESH_INTERVAL } from '@/utils/constants'

const props = defineProps({
  loading: Boolean
})

const emit = defineEmits(['refresh'])
const t = inject('t')

let timer = null

const handleManualRefresh = () => {
  emit('refresh', false) // false means NOT silent
}

const startTimer = () => {
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    if (!props.loading) {
      emit('refresh', true) // true means silent
    }
  }, REFRESH_INTERVAL)
}

onMounted(() => {
  startTimer()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>

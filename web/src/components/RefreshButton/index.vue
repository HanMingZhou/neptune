<template>
  <button
    @click="handleManualRefresh"
    class="list-toolbar-button list-toolbar-button--secondary list-toolbar-button--icon"
    :disabled="loading"
    :title="t('refresh')"
    type="button"
  >
    <span class="material-icons" :class="{ 'animate-spin': loading }"
      >refresh</span
    >
  </button>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, inject } from 'vue'
import { REFRESH_INTERVAL } from '@/utils/constants'
import type { Translator } from '@/types/consoleResource'

const props = withDefaults(
  defineProps<{
    loading?: boolean
  }>(),
  {
    loading: false
  }
)

const emit = defineEmits<{
  refresh: [silent: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

let timer: ReturnType<typeof setInterval> | null = null

const handleManualRefresh = (): void => {
  emit('refresh', false) // false means NOT silent
}

const startTimer = (): void => {
  if (timer) {
    clearInterval(timer)
  }
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
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>

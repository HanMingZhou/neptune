<template>
  <el-dialog
    v-model="visibleModel"
    :title="t('clusterKubeconfig')"
    width="700px"
    align-center
    :before-close="handleBeforeClose"
  >
    <div class="bg-slate-50 dark:bg-zinc-900 border border-border-light dark:border-border-dark rounded-lg p-4 max-h-[500px] overflow-auto">
      <pre class="text-xs font-mono text-slate-700 dark:text-slate-300 whitespace-pre-wrap break-all">{{ content || '-' }}</pre>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <el-button @click="$emit('copy')">{{ t('copy') }}</el-button>
        <el-button type="primary" @click="requestClose">{{ t('close') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'copy', 'update:modelValue'])
const t = inject('t', (key) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const requestClose = () => {
  emit('close')
}

const handleBeforeClose = (done) => {
  emit('close')
  done()
}
</script>

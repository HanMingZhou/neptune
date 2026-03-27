<template>
  <el-dialog v-model="dialogVisible" :title="`${t('expand')}${t('storage')}`" width="450px" class="custom-dialog">
    <div class="space-y-5 py-2">
      <div class="p-4 bg-blue-50/50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-900/30 rounded-xl flex items-center justify-between">
        <span class="text-sm text-blue-600 dark:text-blue-400 font-medium">{{ t('currentSize') }}</span>
        <span class="text-lg font-black text-blue-700 dark:text-blue-300 font-mono">{{ form.currentSize }}</span>
      </div>
      <div>
        <label class="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-2">{{ t('expandTo') }} (GB)</label>
        <div class="flex items-center gap-3">
          <input
            v-model.number="form.newSize"
            type="number"
            :min="form.minSize"
            :max="2000"
            class="flex-1 px-4 py-2 bg-slate-50 dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none transition-all"
          />
          <div class="flex gap-1">
            <button
              class="w-10 h-10 flex items-center justify-center bg-slate-100 dark:bg-zinc-800 hover:bg-slate-200 dark:hover:bg-zinc-700 rounded-lg text-slate-500 transition-colors"
              @click="form.newSize = Math.max(form.minSize, form.newSize - 10)"
            >
              <span class="material-icons text-sm">remove</span>
            </button>
            <button
              class="w-10 h-10 flex items-center justify-center bg-slate-100 dark:bg-zinc-800 hover:bg-slate-200 dark:hover:bg-zinc-700 rounded-lg text-slate-500 transition-colors"
              @click="form.newSize = Math.min(2000, form.newSize + 10)"
            >
              <span class="material-icons text-sm">add</span>
            </button>
          </div>
        </div>
        <p class="text-[11px] text-slate-400 mt-2">{{ t('expandHint') || 'Only upward expansion is supported' }}</p>
      </div>
    </div>
    <template #footer>
      <div class="flex justify-end gap-3 px-1">
        <button
          class="px-5 py-2 rounded-lg text-sm font-bold border border-border-light dark:border-border-dark hover:bg-slate-50 dark:hover:bg-zinc-800 transition-all"
          @click="dialogVisible = false"
        >
          {{ t('cancel') }}
        </button>
        <button
          :disabled="expanding"
          class="bg-primary hover:bg-primary-hover text-white px-6 py-2 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all disabled:opacity-50"
          @click="$emit('submit')"
        >
          <span v-if="expanding" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></span>
          {{ t('confirm') }}{{ t('expand') }}
        </button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  expanding: {
    type: Boolean,
    default: false
  },
  form: {
    type: Object,
    required: true
  },
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'update:modelValue'])
const t = inject('t', (key) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

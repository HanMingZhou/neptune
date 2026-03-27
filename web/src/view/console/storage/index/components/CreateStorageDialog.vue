<template>
  <el-dialog v-model="dialogVisible" :title="`${t('create')}${t('storage')}`" width="500px" class="custom-dialog">
    <div class="space-y-5 py-2">
      <div>
        <label class="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-2">{{ t('cluster') }}</label>
        <el-select
          v-model="form.clusterId"
          :placeholder="t('selectCluster')"
          class="w-full custom-select"
          @change="$emit('cluster-change', $event)"
        >
          <el-option
            v-for="item in clusterOptions"
            :key="item.id"
            :label="`${item.area} - ${item.name}`"
            :value="item.id"
          />
        </el-select>
      </div>
      <div>
        <label class="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-2">{{ t('storageProduct') }}</label>
        <el-select v-model="form.productId" :placeholder="t('selectStorageProduct')" class="w-full custom-select">
          <el-option v-for="item in storageProducts" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
      </div>
      <div>
        <label class="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-2">{{ t('name') }}</label>
        <input
          v-model="form.name"
          type="text"
          :placeholder="t('inputName')"
          class="w-full px-4 py-2 bg-slate-50 dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none transition-all"
        />
      </div>
      <div>
        <label class="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-2">{{ t('capacity') }} (GB)</label>
        <div class="flex items-center gap-3">
          <input
            v-model.number="form.size"
            type="number"
            :min="10"
            :max="2000"
            class="flex-1 px-4 py-2 bg-slate-50 dark:bg-zinc-800 border border-border-light dark:border-border-dark rounded-lg text-sm focus:ring-1 focus:ring-primary outline-none transition-all"
          />
          <div class="flex gap-1">
            <button
              class="w-10 h-10 flex items-center justify-center bg-slate-100 dark:bg-zinc-800 hover:bg-slate-200 dark:hover:bg-zinc-700 rounded-lg text-slate-500 transition-colors"
              @click="form.size = Math.max(10, form.size - 10)"
            >
              <span class="material-icons text-sm">remove</span>
            </button>
            <button
              class="w-10 h-10 flex items-center justify-center bg-slate-100 dark:bg-zinc-800 hover:bg-slate-200 dark:hover:bg-zinc-700 rounded-lg text-slate-500 transition-colors"
              @click="form.size = Math.min(2000, form.size + 10)"
            >
              <span class="material-icons text-sm">add</span>
            </button>
          </div>
        </div>
        <p class="text-[11px] text-slate-400 mt-2">{{ t('capacityRange') || 'Range: 10GB - 2000GB' }}</p>
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
          :disabled="creating"
          class="bg-primary hover:bg-primary-hover text-white px-6 py-2 rounded-lg font-bold text-sm shadow-lg shadow-primary/20 flex items-center gap-2 transition-all disabled:opacity-50"
          @click="$emit('submit')"
        >
          <span v-if="creating" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></span>
          {{ t('create') }}
        </button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  clusterOptions: {
    type: Array,
    default: () => []
  },
  creating: {
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
  },
  storageProducts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['cluster-change', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

<template>
  <div class="console-filter-card px-5 py-4">
    <div class="list-filter-bar">
      <div class="list-toolbar-main">
        <div class="list-search-field max-w-[320px]">
          <span class="material-icons absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
          <input
            :value="searchName"
            type="text"
            :placeholder="t('searchStoragePlaceholder')"
            class="list-search-input !w-full"
            @input="emit('update:search-name', $event.target.value)"
            @keyup.enter="emit('refresh')"
          />
        </div>

        <el-select
          :model-value="searchStatus"
          :placeholder="`${t('status')}: ${t('all')}`"
          clearable
          class="list-filter-select !w-[168px]"
          size="small"
          @update:model-value="emit('update:search-status', $event)"
        >
          <el-option :label="t('Creating')" value="Creating" />
          <el-option :label="t('Bound')" value="Bound" />
          <el-option :label="t('Error')" value="Error" />
        </el-select>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  searchName: {
    type: String,
    default: ''
  },
  searchStatus: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['refresh', 'update:search-name', 'update:search-status'])
const t = inject('t', (key) => key)
</script>

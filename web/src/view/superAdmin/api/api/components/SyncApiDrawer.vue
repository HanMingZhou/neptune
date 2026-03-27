<template>
  <el-drawer
    v-model="visibleModel"
    :size="900"
    :before-close="requestClose"
    :show-close="false"
  >
    <template #header>
      <div class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ t('syncRoute') }}</span>
        <div class="flex gap-3">
          <el-button :loading="apiCompletionLoading" @click="requestClose">{{ t('cancel') }}</el-button>
          <el-button type="primary" :loading="syncing || apiCompletionLoading" @click="$emit('submit')">
            {{ t('confirm') }}
          </el-button>
        </div>
      </div>
    </template>

    <div class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3 mb-6">
      <span class="material-icons text-amber-500">info</span>
      <span class="text-sm text-amber-700 dark:text-amber-400">{{ t('syncTip') }}</span>
    </div>

    <div class="mb-6">
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('newRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('newRouteTip') }}</span>
        <button
          class="ml-auto px-3 py-1.5 bg-primary hover:bg-primary-hover text-white rounded text-xs font-medium flex items-center gap-1"
          @click="$emit('ai-completion')"
        >
          {{ t('aiCompletion') }}
        </button>
      </div>

      <el-table v-loading="syncing || apiCompletionLoading" :data="syncApiData.newApis" class="rounded-lg overflow-hidden">
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column :label="t('apiGroup')" min-width="150">
          <template #default="{ row }">
            <el-select v-model="row.apiGroup" :placeholder="t('select')" allow-create filterable>
              <el-option v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="t('apiDesc')" min-width="150">
          <template #default="{ row }">
            <el-input v-model="row.description" :placeholder="t('apiDesc')" />
          </template>
        </el-table-column>
        <el-table-column :label="t('method')" min-width="100">
          <template #default="{ row }">{{ row.method }}</template>
        </el-table-column>
        <el-table-column :label="t('actions')" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$emit('add-one', row)">{{ t('singleAdd') }}</el-button>
            <el-button type="primary" link @click="$emit('ignore', { row, flag: true })">{{ t('ignore') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="mb-6">
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('deletedRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('deletedRouteTip') }}</span>
      </div>

      <el-table :data="syncApiData.deleteApis" class="rounded-lg overflow-hidden">
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column :label="t('apiGroup')" min-width="150" prop="apiGroup" />
        <el-table-column :label="t('apiDesc')" min-width="150" prop="description" />
        <el-table-column :label="t('method')" min-width="100" prop="method" />
      </el-table>
    </div>

    <div>
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('ignoredRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('ignoredRouteTip') }}</span>
      </div>

      <el-table :data="syncApiData.ignoreApis" class="rounded-lg overflow-hidden">
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column :label="t('apiGroup')" min-width="150" prop="apiGroup" />
        <el-table-column :label="t('apiDesc')" min-width="150" prop="description" />
        <el-table-column :label="t('method')" min-width="100" prop="method" />
        <el-table-column :label="t('actions')" min-width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$emit('ignore', { row, flag: false })">{{ t('cancelIgnore') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </el-drawer>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  apiCompletionLoading: {
    type: Boolean,
    default: false
  },
  apiGroupOptions: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  syncing: {
    type: Boolean,
    default: false
  },
  syncApiData: {
    type: Object,
    default: () => ({ newApis: [], deleteApis: [], ignoreApis: [] })
  }
})

const emit = defineEmits(['add-one', 'ai-completion', 'close', 'ignore', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const requestClose = () => {
  emit('close')
}
</script>

<template>
  <BaseDrawer
    v-model="visibleModel"
    :size="900"
    :show-close="false"
    @close="emit('close')"
  >
    <template #header="{ requestClose }">
      <div class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ t('syncRoute') }}</span>
        <div class="flex gap-3">
          <el-button :loading="apiCompletionLoading" @click="requestClose">{{
            t('cancel')
          }}</el-button>
          <el-button
            type="primary"
            :loading="syncing || apiCompletionLoading"
            @click="emit('submit')"
          >
            {{ t('confirm') }}
          </el-button>
        </div>
      </div>
    </template>

    <div
      class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3 mb-6"
    >
      <span class="material-icons text-amber-500">info</span>
      <span class="text-sm text-amber-700 dark:text-amber-400">{{
        t('syncTip')
      }}</span>
    </div>

    <div class="mb-6">
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('newRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('newRouteTip') }}</span>
        <el-button
          class="ml-auto"
          size="small"
          type="primary"
          :loading="apiCompletionLoading"
          @click="emit('ai-completion')"
        >
          {{ t('aiCompletion') }}
        </el-button>
      </div>

      <el-table
        v-loading="syncing || apiCompletionLoading"
        :data="syncApiData.newApis"
        class="rounded-lg overflow-hidden"
      >
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column :label="t('apiGroup')" min-width="150">
          <template #default="{ row }">
            <el-select
              v-model="row.apiGroup"
              :placeholder="t('select')"
              allow-create
              filterable
            >
              <el-option
                v-for="item in apiGroupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
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
            <el-button type="primary" link @click="emit('add-one', row)">{{
              t('singleAdd')
            }}</el-button>
            <el-button
              type="primary"
              link
              @click="emit('ignore', { row, flag: true })"
              >{{ t('ignore') }}</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="mb-6">
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('deletedRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('deletedRouteTip') }}</span>
      </div>

      <el-table
        :data="syncApiData.deleteApis"
        class="rounded-lg overflow-hidden"
      >
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column
          :label="t('apiGroup')"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          :label="t('apiDesc')"
          min-width="150"
          prop="description"
        />
        <el-table-column :label="t('method')" min-width="100" prop="method" />
      </el-table>
    </div>

    <div>
      <div class="flex items-center gap-3 mb-4">
        <h4 class="font-bold">{{ t('ignoredRoute') }}</h4>
        <span class="text-xs text-slate-500">{{ t('ignoredRouteTip') }}</span>
      </div>

      <el-table
        :data="syncApiData.ignoreApis"
        class="rounded-lg overflow-hidden"
      >
        <el-table-column :label="t('apiPath')" min-width="150" prop="path" />
        <el-table-column
          :label="t('apiGroup')"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          :label="t('apiDesc')"
          min-width="150"
          prop="description"
        />
        <el-table-column :label="t('method')" min-width="100" prop="method" />
        <el-table-column :label="t('actions')" min-width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="emit('ignore', { row, flag: false })"
              >{{ t('cancelIgnore') }}</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </div>
  </BaseDrawer>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import BaseDrawer from '@/components/base/BaseDrawer.vue'
import type { Translator } from '@/types/consoleResource'
import type {
  ApiListItem,
  ApiSyncData,
  LabelValueOption
} from '@/types/superAdmin'

interface SyncIgnorePayload {
  row: ApiListItem
  flag: boolean
}

const props = withDefaults(
  defineProps<{
    apiCompletionLoading?: boolean
    apiGroupOptions?: LabelValueOption[]
    modelValue?: boolean
    syncing?: boolean
    syncApiData: ApiSyncData
  }>(),
  {
    apiCompletionLoading: false,
    apiGroupOptions: () => [],
    modelValue: false,
    syncing: false
  }
)

const emit = defineEmits<{
  'add-one': [row: ApiListItem]
  'ai-completion': []
  close: []
  ignore: [payload: SyncIgnorePayload]
  submit: []
  'update:modelValue': [value: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})
</script>

<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    class="storage-create-dialog"
    :model="form"
    :rules="rules"
    shell-size="lg"
    :submit-text="t('create')"
    :submitting="creating"
    :title="`${t('create')}${t('storage')}`"
    form-class="py-2"
    label-position="top"
    @submit="emit('submit')"
  >
    <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
      <el-form-item :label="t('cluster')" prop="clusterId">
        <el-select
          v-model="form.clusterId"
          :placeholder="t('selectCluster')"
          class="w-full"
          @change="emit('cluster-change', $event)"
        >
          <el-option
            v-for="item in clusterOptions"
            :key="item.id"
            :label="`${item.area} - ${item.name}`"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('storageProduct')" prop="productId">
        <el-select
          v-model="form.productId"
          :placeholder="t('selectStorageProduct')"
          class="w-full"
        >
          <el-option
            v-for="item in storageProducts"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('name')" prop="name" class="md:col-span-2">
        <el-input
          v-model="form.name"
          :placeholder="t('inputName')"
          class="w-full"
        />
      </el-form-item>

      <el-form-item
        :label="`${t('capacity')} (GB)`"
        prop="size"
        class="md:col-span-2"
      >
        <el-input-number
          v-model.number="form.size"
          :min="10"
          :max="2000"
          :step="10"
          controls-position="right"
          class="w-full storage-create-dialog__number"
        />
        <p class="mt-2 text-[11px] text-slate-400">
          {{ t('capacityRange') || 'Range: 10GB - 2000GB' }}
        </p>
      </el-form-item>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type {
  StorageClusterOption,
  StorageCreateForm,
  StorageProductOption
} from '@/types/storage'

const props = withDefaults(
  defineProps<{
    clusterOptions?: StorageClusterOption[]
    creating?: boolean
    form: StorageCreateForm
    modelValue?: boolean
    storageProducts?: StorageProductOption[]
  }>(),
  {
    clusterOptions: () => [],
    creating: false,
    modelValue: false,
    storageProducts: () => []
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: StorageCreateForm['clusterId']]
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const rules: FormRules<StorageCreateForm> = {
  clusterId: [
    { required: true, message: t('fillAllFields'), trigger: 'change' }
  ],
  productId: [
    { required: true, message: t('fillAllFields'), trigger: 'change' }
  ],
  name: [{ required: true, message: t('inputName'), trigger: 'blur' }],
  size: [
    { required: true, message: t('capacityRange'), trigger: 'change' },
    {
      validator: (_rule, value: number, callback) => {
        if (typeof value !== 'number' || value < 10 || value > 2000) {
          callback(new Error(t('capacityRange')))
          return
        }

        callback()
      },
      trigger: 'change'
    }
  ]
}
</script>

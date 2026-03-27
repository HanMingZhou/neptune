<template>
  <el-dialog
    v-if="productForm.productType === 1"
    v-model="editorVisibleModel"
    :title="dialogTitle"
    width="880px"
    align-center
    class="product-dialog"
  >
    <el-form ref="productFormRef" :model="productForm" label-width="100px" :rules="computeRules" class="product-form">
      <div v-if="!isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-primary text-[18px]">hub</span>
          <span>{{ t('selectCluster') }}</span>
        </div>
        <el-form-item :label="t('cluster')" prop="clusterId" class="mt-4">
          <el-select v-model="productForm.clusterId" :placeholder="t('selectCluster')" class="w-full" @change="emit('cluster-change', $event)">
            <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
          </el-select>
        </el-form-item>
      </div>

      <div v-if="productForm.clusterId && !isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-violet-500 text-[18px]">dns</span>
          <span>{{ t('nodeList') }}</span>
        </div>
        <el-form-item prop="nodeName" class="!m-0 !h-0 overflow-hidden">
          <el-input v-model="productForm.nodeName" class="hidden" />
        </el-form-item>
        <div class="mt-4 grid max-h-[280px] grid-cols-2 gap-3 overflow-y-auto pr-1" v-loading="loadingNodes">
          <div v-if="clusterNodes.length === 0" class="col-span-2 py-12 text-center text-sm text-slate-400">
            <span class="material-icons mb-2 block text-3xl opacity-30">dns</span>
            {{ t('noNodes') }}
          </div>
          <div
            v-for="node in clusterNodes"
            :key="node.nodeName"
            class="node-card"
            :class="{ 'is-selected': productForm.nodeName === node.nodeName }"
            @click="emit('node-select', node)"
          >
            <div class="mb-2.5 flex items-center justify-between">
              <span class="mr-2 truncate font-mono text-sm font-bold text-slate-800 dark:text-slate-200">{{ node.nodeName }}</span>
              <span
                class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-bold"
                :class="node.schedulable ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'"
              >
                {{ node.schedulable ? t('schedulable') : t('unschedulable') }}
              </span>
            </div>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-slate-500">
              <span class="flex items-center gap-1"><span class="material-icons text-[12px] text-blue-400">memory</span>CPU: {{ node.cpu }}核</span>
              <span class="flex items-center gap-1"><span class="material-icons text-[12px] text-cyan-400">sd_card</span>{{ node.memory }}GB</span>
              <span v-if="node.gpuCount > 0" class="flex items-center gap-1 font-bold text-amber-600"><span class="material-icons text-[12px]">developer_board</span>{{ node.gpuModel || 'GPU' }} × {{ node.gpuCount }}</span>
              <span v-if="node.vGpuNumber > 0" class="flex items-center gap-1 font-bold text-violet-600"><span class="material-icons text-[12px]">grid_view</span>vGPU: {{ node.vGpuNumber }}</span>
              <span v-if="!node.gpuCount && !node.vGpuNumber" class="text-slate-400">CPU Only</span>
            </div>
          </div>
        </div>
      </div>

      <template v-if="isEdit || productForm.nodeName">
        <div class="section-header">
          <span class="material-icons text-blue-500 text-[18px]">inventory_2</span>
          <span>{{ t('productInfo') }}</span>
        </div>
        <div class="form-section">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('productName')" prop="name">
                <el-input v-model="productForm.name" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('area')" prop="area">
                <el-input v-model="productForm.area" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('cpuCores')" prop="cpu">
                <el-input-number v-model="productForm.cpu" :min="1" :max="nodeMaxCpu" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('memoryGb')" prop="memory">
                <el-input-number v-model="productForm.memory" :min="1" :max="nodeMaxMemory" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('cpuModel')">
                <el-input v-model="productForm.cpuModel" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-amber-500 text-[18px]">developer_board</span>
          <span>{{ t('resourceType') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('resourceType')">
            <el-radio-group v-model="resourceTypeModel" class="resource-type-group" @change="emit('resource-type-change', $event)">
              <el-radio-button value="cpu">
                <span class="flex items-center gap-1.5"><span class="material-icons text-[14px]">memory</span>CPU Only</span>
              </el-radio-button>
              <el-radio-button value="gpu">
                <span class="flex items-center gap-1.5"><span class="material-icons text-[14px]">developer_board</span>GPU</span>
              </el-radio-button>
              <el-radio-button value="vgpu">
                <span class="flex items-center gap-1.5"><span class="material-icons text-[14px]">grid_view</span>vGPU</span>
              </el-radio-button>
            </el-radio-group>
          </el-form-item>

          <div v-show="resourceType === 'gpu'">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="t('gpuModel')">
                  <el-input v-model="productForm.gpuModel" class="w-full" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('gpuCount')">
                  <el-input-number v-model="productForm.gpuCount" :min="0" :max="nodeMaxGpuCount" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="t('gpuMemory')">
                  <el-input-number v-model="productForm.gpuMemory" :min="0" :max="nodeMaxGpuMemory" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-show="resourceType === 'vgpu'">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="t('vGpuCount')">
                  <el-input-number v-model="productForm.vGpuCount" :min="0" :max="nodeMaxVGpuCount" class="w-full" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('vGpuMemory')">
                  <el-input-number v-model="productForm.vGpuMemory" :min="0" :max="nodeMaxVGpuMemory" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="t('vGpuCores')">
                  <el-input-number v-model="productForm.vGpuCores" :min="0" :max="nodeMaxVGpuCores" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>

        <div class="section-header">
          <span class="material-icons text-emerald-500 text-[18px]">payments</span>
          <span>{{ t('priceSettingsSingleCard') }}</span>
        </div>
        <div class="form-section">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('priceHourly')" prop="priceHourly">
                <el-input-number v-model="productForm.priceHourly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('priceDaily')">
                <el-input-number v-model="productForm.priceDaily" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('priceWeekly')">
                <el-input-number v-model="productForm.priceWeekly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('priceMonthly')">
                <el-input-number v-model="productForm.priceMonthly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-slate-400 text-[18px]">tune</span>
          <span>{{ t('otherSettings') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('status')">
            <el-switch v-model="productForm.status" :active-value="1" :inactive-value="0" :active-text="t('onShelf')" :inactive-text="t('offShelf')" />
          </el-form-item>
          <el-form-item :label="t('paramDesc')" class="align-start">
            <el-input v-model="productForm.description" type="textarea" :rows="3" :placeholder="t('inputProductDesc')" />
          </el-form-item>
        </div>
      </template>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <button class="btn-cancel" @click="editorVisibleModel = false">{{ t('cancel') }}</button>
        <button class="btn-primary" :disabled="!canSubmit || submitting" @click="submitProduct">
          <span v-if="submitting" class="material-icons mr-1 animate-spin text-sm">autorenew</span>
          {{ isEdit ? t('save') : t('create') }}
        </button>
      </div>
    </template>
  </el-dialog>

  <el-dialog
    v-else
    v-model="editorVisibleModel"
    :title="dialogTitle"
    width="600px"
    align-center
    class="product-dialog"
  >
    <el-form ref="productFormRef" :model="productForm" label-width="100px" :rules="storageRules" class="product-form">
      <div v-if="!isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-primary text-[18px]">hub</span>
          <span>{{ t('selectCluster') }}</span>
        </div>
        <el-form-item :label="t('cluster')" prop="clusterId" class="mt-4">
          <el-select v-model="productForm.clusterId" :placeholder="t('selectCluster')" class="w-full" @change="emit('storage-cluster-change', $event)">
            <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
          </el-select>
        </el-form-item>
      </div>

      <template v-if="isEdit || productForm.clusterId">
        <div class="section-header">
          <span class="material-icons text-purple-500 text-[18px]">folder_open</span>
          <span>{{ t('productInfo') }}</span>
        </div>
        <div class="form-section">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('productName')" prop="name">
                <el-input v-model="productForm.name" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('area')" prop="area">
                <el-input v-model="productForm.area" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="t('storageClass')" prop="storageClass">
                <el-input v-model="productForm.storageClass" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('storagePriceGb')" prop="storagePriceGb">
                <el-input-number v-model="productForm.storagePriceGb" :precision="4" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-slate-400 text-[18px]">tune</span>
          <span>{{ t('otherSettings') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('status')">
            <el-switch v-model="productForm.status" :active-value="1" :inactive-value="0" :active-text="t('onShelf')" :inactive-text="t('offShelf')" />
          </el-form-item>
          <el-form-item :label="t('paramDesc')" class="align-start">
            <el-input v-model="productForm.description" type="textarea" :rows="3" :placeholder="t('inputProductDesc')" />
          </el-form-item>
        </div>
      </template>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <button class="btn-cancel" @click="editorVisibleModel = false">{{ t('cancel') }}</button>
        <button class="btn-primary" :disabled="!canSubmit || submitting" @click="submitProduct">
          <span v-if="submitting" class="material-icons mr-1 animate-spin text-sm">autorenew</span>
          {{ isEdit ? t('save') : t('create') }}
        </button>
      </div>
    </template>
  </el-dialog>

  <el-dialog
    v-model="priceVisibleModel"
    :title="t('adjustPrice')"
    width="520px"
    align-center
    class="product-dialog"
  >
    <div class="product-form">
      <div class="section-header">
        <span class="material-icons text-amber-500 text-[18px]">payments</span>
        <span>{{ t('adjustPrice') }}</span>
      </div>
      <el-form :model="priceForm" label-width="100px" class="form-section">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('priceHourly')">
              <el-input-number v-model="priceForm.priceHourly" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('priceDaily')">
              <el-input-number v-model="priceForm.priceDaily" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('priceWeekly')">
              <el-input-number v-model="priceForm.priceWeekly" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('priceMonthly')">
              <el-input-number v-model="priceForm.priceMonthly" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <button class="btn-cancel" @click="priceVisibleModel = false">{{ t('cancel') }}</button>
        <button class="btn-primary" :disabled="submitting" @click="emit('submit-price')">
          <span v-if="submitting" class="material-icons mr-1 animate-spin text-sm">autorenew</span>
          {{ t('confirm') }}
        </button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

const props = defineProps({
  canSubmit: {
    type: Boolean,
    default: false
  },
  clusterNodes: {
    type: Array,
    default: () => []
  },
  clusters: {
    type: Array,
    default: () => []
  },
  computeRules: {
    type: Object,
    required: true
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  loadingNodes: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  nodeMaxCpu: {
    type: Number,
    default: 256
  },
  nodeMaxGpuCount: {
    type: Number,
    default: 16
  },
  nodeMaxGpuMemory: {
    type: Number,
    default: 256
  },
  nodeMaxMemory: {
    type: Number,
    default: 2048
  },
  nodeMaxVGpuCores: {
    type: Number,
    default: 100
  },
  nodeMaxVGpuCount: {
    type: Number,
    default: 100
  },
  nodeMaxVGpuMemory: {
    type: Number,
    default: 256
  },
  priceForm: {
    type: Object,
    required: true
  },
  priceVisible: {
    type: Boolean,
    default: false
  },
  productForm: {
    type: Object,
    required: true
  },
  resourceType: {
    type: String,
    default: 'cpu'
  },
  storageRules: {
    type: Object,
    required: true
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'cluster-change',
  'node-select',
  'resource-type-change',
  'storage-cluster-change',
  'submit-price',
  'submit-product',
  'update:modelValue',
  'update:price-visible',
  'update:resource-type'
])

const t = inject('t', (key) => key)
const productFormRef = ref(null)

const editorVisibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const priceVisibleModel = computed({
  get: () => props.priceVisible,
  set: (value) => emit('update:price-visible', value)
})

const resourceTypeModel = computed({
  get: () => props.resourceType,
  set: (value) => emit('update:resource-type', value)
})

const submitProduct = async () => {
  if (!productFormRef.value) {
    return
  }

  await productFormRef.value.validate()
  emit('submit-product')
}
</script>

<style scoped>
:deep(.product-dialog .el-dialog) {
  border-radius: 16px;
  overflow: hidden;
  padding: 0;
}

:deep(.product-dialog .el-dialog__header) {
  margin: 0;
  padding: 20px 24px;
  border-bottom: 1px solid rgb(241 245 249);
  background: white;
}

:deep(.dark .product-dialog .el-dialog__header) {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

:deep(.product-dialog .el-dialog__title) {
  font-size: 16px;
  font-weight: 700;
  color: rgb(15 23 42);
}

:deep(.dark .product-dialog .el-dialog__title) {
  color: rgb(226 232 240);
}

:deep(.product-dialog .el-dialog__body) {
  padding: 24px;
  background: rgb(248 250 252);
  max-height: 70vh;
  overflow-y: auto;
}

:deep(.dark .product-dialog .el-dialog__body) {
  background: rgb(9 9 11);
}

:deep(.product-dialog .el-dialog__footer) {
  padding: 16px 24px;
  margin: 0;
  border-top: 1px solid rgb(241 245 249);
  background: white;
}

:deep(.dark .product-dialog .el-dialog__footer) {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgb(241 245 249);
  font-size: 13px;
  font-weight: 700;
  color: rgb(51 65 85);
}

:deep(.dark) .section-header {
  border-color: rgb(39 39 42);
  color: rgb(203 213 225);
}

.form-section {
  margin-bottom: 20px;
  border: 1px solid rgb(241 245 249);
  border-radius: 12px;
  background: white;
  padding: 20px 16px 4px;
}

:deep(.dark) .form-section {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.product-form :deep(.el-form-item) {
  margin-bottom: 18px;
  align-items: center;
}

.product-form :deep(.el-form-item__label) {
  font-size: 13px;
  font-weight: 600;
  color: rgb(71 85 105);
}

:deep(.dark) .product-form :deep(.el-form-item__label) {
  color: rgb(148 163 184);
}

.product-form :deep(.el-form-item.align-start) {
  align-items: flex-start;
}

.product-form :deep(.el-form-item.align-start .el-form-item__label) {
  padding-top: 6px;
}

.node-card {
  padding: 14px;
  border: 1.5px solid rgb(226 232 240);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: white;
}

:deep(.dark) .node-card {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.node-card:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 12px rgba(59, 130, 246, 0.08);
}

.node-card.is-selected {
  border-color: var(--el-color-primary);
  background: rgba(59, 130, 246, 0.04);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

:deep(.dark) .node-card.is-selected {
  background: rgba(59, 130, 246, 0.08);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-cancel {
  padding: 8px 20px;
  border-radius: 10px;
  border: 1px solid rgb(226 232 240);
  background: white;
  color: rgb(100 116 139);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
}

:deep(.dark) .btn-cancel {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
  color: rgb(148 163 184);
}

.btn-cancel:hover {
  background: rgb(248 250 252);
}

:deep(.dark) .btn-cancel:hover {
  background: rgb(39 39 42);
}

.btn-primary {
  display: flex;
  align-items: center;
  padding: 8px 24px;
  border: none;
  border-radius: 10px;
  background: var(--el-color-primary);
  color: white;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
  transition: all 0.15s;
}

.btn-primary:hover:not(:disabled) {
  filter: brightness(1.1);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.resource-type-group :deep(.el-radio-button__inner) {
  font-size: 13px;
  font-weight: 600;
}

.product-form :deep(.el-input-number) {
  width: 100% !important;
}
</style>

<template>
  <el-dialog
    v-if="productForm.productType === 1"
    v-model="editorVisibleModel"
    :title="dialogTitle"
    width="840px"
    class="product-dialog product-dialog--compute"
  >
    <el-form ref="productFormRef" :model="productForm" label-width="96px" :rules="computeRules" class="product-form">
      <div v-if="!isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-primary text-[18px]">hub</span>
          <span>{{ t('selectCluster') }}</span>
        </div>
        <div class="form-section form-section--compact">
          <el-form-item prop="clusterId" class="compact-form-item !mb-0">
            <div class="compact-field">
              <div class="compact-field-label">
                <span class="compact-field-required">*</span>
                <span>{{ t('cluster') }}</span>
              </div>
              <el-select id="compute-cluster-id" v-model="productForm.clusterId" :placeholder="t('selectCluster')" class="w-full" @change="emit('cluster-change', $event)">
                <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
              </el-select>
            </div>
          </el-form-item>
        </div>
      </div>

      <div v-if="productForm.clusterId && !isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-violet-500 text-[18px]">dns</span>
          <span>{{ t('nodeList') }}</span>
        </div>
        <div class="form-section form-section--compact">
          <el-form-item prop="nodeName" class="!m-0 !h-0 overflow-hidden">
            <el-input v-model="productForm.nodeName" class="hidden" />
          </el-form-item>
          <div class="grid max-h-[280px] grid-cols-1 2xl:grid-cols-2 gap-3 overflow-y-auto pr-1" v-loading="loadingNodes">
            <div v-if="clusterNodes.length === 0" class="2xl:col-span-2 py-12 text-center text-sm text-slate-400">
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
      </div>

      <template v-if="isEdit || productForm.nodeName">
        <div class="section-header">
          <span class="material-icons text-blue-500 text-[18px]">inventory_2</span>
          <span>{{ t('productInfo') }}</span>
        </div>
        <div class="form-section">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('productName')" prop="name" for="compute-name">
                <el-input id="compute-name" v-model="productForm.name" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('area')" prop="area" for="compute-area">
                <el-input id="compute-area" v-model="productForm.area" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('cpuCores')" prop="cpu" for="compute-cpu">
                <el-input-number id="compute-cpu" v-model="productForm.cpu" :min="1" :max="nodeMaxCpu" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('memoryGb')" prop="memory" for="compute-memory">
                <el-input-number id="compute-memory" v-model="productForm.memory" :min="1" :max="nodeMaxMemory" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('cpuModel')" for="compute-cpu-model">
                <el-input id="compute-cpu-model" v-model="productForm.cpuModel" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-amber-500 text-[18px]">developer_board</span>
          <span>{{ t('resourceType') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('resourceType')" for="compute-resource-type">
            <el-radio-group id="compute-resource-type" v-model="resourceTypeModel" class="resource-type-group" @change="emit('resource-type-change', $event)">
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
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('gpuModel')" for="compute-gpu-model">
                  <el-input id="compute-gpu-model" v-model="productForm.gpuModel" class="w-full" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('gpuCount')" for="compute-gpu-count">
                  <el-input-number id="compute-gpu-count" v-model="productForm.gpuCount" :min="0" :max="nodeMaxGpuCount" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('gpuMemory')" for="compute-gpu-memory">
                  <el-input-number id="compute-gpu-memory" v-model="productForm.gpuMemory" :min="0" :max="nodeMaxGpuMemory" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-show="resourceType === 'vgpu'">
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('vGpuCount')" for="compute-vgpu-count">
                  <el-input-number id="compute-vgpu-count" v-model="productForm.vGpuCount" :min="0" :max="nodeMaxVGpuCount" class="w-full" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('vGpuMemory')" for="compute-vgpu-memory">
                  <el-input-number id="compute-vgpu-memory" v-model="productForm.vGpuMemory" :min="0" :max="nodeMaxVGpuMemory" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item :label="t('vGpuCores')" for="compute-vgpu-cores">
                  <el-input-number id="compute-vgpu-cores" v-model="productForm.vGpuCores" :min="0" :max="nodeMaxVGpuCores" class="w-full" />
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
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('priceHourly')" prop="priceHourly" for="compute-price-hourly">
                <el-input-number id="compute-price-hourly" v-model="productForm.priceHourly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('priceDaily')" for="compute-price-daily">
                <el-input-number id="compute-price-daily" v-model="productForm.priceDaily" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('priceWeekly')" for="compute-price-weekly">
                <el-input-number id="compute-price-weekly" v-model="productForm.priceWeekly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('priceMonthly')" for="compute-price-monthly">
                <el-input-number id="compute-price-monthly" v-model="productForm.priceMonthly" :precision="2" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-slate-400 text-[18px]">tune</span>
          <span>{{ t('otherSettings') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('status')" for="compute-status">
            <el-switch id="compute-status" v-model="productForm.status" :active-value="1" :inactive-value="0" :active-text="t('onShelf')" :inactive-text="t('offShelf')" />
          </el-form-item>
          <el-form-item :label="t('paramDesc')" for="compute-description" class="align-start">
            <el-input id="compute-description" v-model="productForm.description" type="textarea" :rows="3" :placeholder="t('inputProductDesc')" />
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
    class="product-dialog product-dialog--storage"
  >
    <el-form ref="productFormRef" :model="productForm" label-width="96px" :rules="storageRules" class="product-form">
      <div v-if="!isEdit" class="mb-6">
        <div class="section-header">
          <span class="material-icons text-primary text-[18px]">hub</span>
          <span>{{ t('selectCluster') }}</span>
        </div>
        <div class="form-section form-section--compact">
          <el-form-item prop="clusterId" class="compact-form-item !mb-0">
            <div class="compact-field">
              <div class="compact-field-label">
                <span class="compact-field-required">*</span>
                <span>{{ t('cluster') }}</span>
              </div>
              <el-select id="storage-cluster-id" v-model="productForm.clusterId" :placeholder="t('selectCluster')" class="w-full" @change="emit('storage-cluster-change', $event)">
                <el-option v-for="cluster in clusters" :key="cluster.id" :label="cluster.name" :value="cluster.id" />
              </el-select>
            </div>
          </el-form-item>
        </div>
      </div>

      <template v-if="isEdit || productForm.clusterId">
        <div class="section-header">
          <span class="material-icons text-purple-500 text-[18px]">folder_open</span>
          <span>{{ t('productInfo') }}</span>
        </div>
        <div class="form-section">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('productName')" prop="name" for="storage-name">
                <el-input id="storage-name" v-model="productForm.name" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('area')" prop="area" for="storage-area">
                <el-input id="storage-area" v-model="productForm.area" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('storageClass')" prop="storageClass" for="storage-class">
                <el-input id="storage-class" v-model="productForm.storageClass" class="w-full" />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('storagePriceGb')" prop="storagePriceGb" for="storage-price-gb">
                <el-input-number id="storage-price-gb" v-model="productForm.storagePriceGb" :precision="4" :min="0" class="w-full" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="section-header">
          <span class="material-icons text-slate-400 text-[18px]">tune</span>
          <span>{{ t('otherSettings') }}</span>
        </div>
        <div class="form-section">
          <el-form-item :label="t('status')" for="storage-status">
            <el-switch id="storage-status" v-model="productForm.status" :active-value="1" :inactive-value="0" :active-text="t('onShelf')" :inactive-text="t('offShelf')" />
          </el-form-item>
          <el-form-item :label="t('paramDesc')" for="storage-description" class="align-start">
            <el-input id="storage-description" v-model="productForm.description" type="textarea" :rows="3" :placeholder="t('inputProductDesc')" />
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
    class="product-dialog product-dialog--price"
  >
    <div class="product-form">
      <div class="section-header">
        <span class="material-icons text-amber-500 text-[18px]">payments</span>
        <span>{{ t('adjustPrice') }}</span>
      </div>
      <el-form :model="priceForm" label-width="96px" class="form-section">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceHourly')" for="adjust-price-hourly">
              <el-input-number id="adjust-price-hourly" v-model="priceForm.priceHourly" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceDaily')" for="adjust-price-daily">
              <el-input-number id="adjust-price-daily" v-model="priceForm.priceDaily" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceWeekly')" for="adjust-price-weekly">
              <el-input-number id="adjust-price-weekly" v-model="priceForm.priceWeekly" :precision="2" :min="0" class="w-full" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceMonthly')" for="adjust-price-monthly">
              <el-input-number id="adjust-price-monthly" v-model="priceForm.priceMonthly" :precision="2" :min="0" class="w-full" />
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

<style>
.product-dialog .el-dialog {
  border-radius: 16px;
  overflow: hidden;
  padding: 0;
  margin: 24px auto !important;
  max-width: calc(100vw - 64px);
}

.product-dialog--compute .el-dialog {
  width: min(840px, calc(100vw - 64px)) !important;
}

.product-dialog--storage .el-dialog {
  width: min(600px, calc(100vw - 64px)) !important;
}

.product-dialog--price .el-dialog {
  width: min(520px, calc(100vw - 64px)) !important;
}

.product-dialog .el-dialog__header {
  margin: 0;
  padding: 20px 24px;
  border-bottom: 1px solid rgb(241 245 249);
  background: white;
}

.dark .product-dialog .el-dialog__header {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.product-dialog .el-dialog__title {
  font-size: 16px;
  font-weight: 700;
  color: rgb(15 23 42);
}

.dark .product-dialog .el-dialog__title {
  color: rgb(226 232 240);
}

.product-dialog .el-dialog__body {
  padding: 24px 24px 14px;
  background: rgb(248 250 252);
  max-height: 70vh;
  overflow-y: auto;
  overflow-x: hidden;
  box-sizing: border-box;
}

.dark .product-dialog .el-dialog__body {
  background: rgb(9 9 11);
}

.product-dialog .el-dialog__footer {
  padding: 8px 24px 18px;
  margin: 0;
  border-top: 1px solid rgb(241 245 249);
  background: white;
  overflow-x: hidden;
  box-sizing: border-box;
}

.dark .product-dialog .el-dialog__footer {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.product-dialog .section-header {
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

.dark .product-dialog .section-header {
  border-color: rgb(39 39 42);
  color: rgb(203 213 225);
}

.product-dialog .form-section {
  margin: 0 auto 14px;
  border: 1px solid rgb(241 245 249);
  border-radius: 14px;
  background: white;
  padding: 18px 18px 6px;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
  max-width: 100%;
}

.product-dialog .form-section--compact {
  padding: 18px 20px 16px;
}

.dark .product-dialog .form-section {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.product-dialog .product-form > .mb-6 {
  margin-bottom: 14px;
}

.product-dialog .compact-form-item .el-form-item__content {
  display: block;
  width: 100% !important;
  margin-left: 0 !important;
}

.product-dialog .compact-field {
  width: 100%;
}

.product-dialog .compact-field-label {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 10px;
  color: rgb(51 65 85);
  font-size: 13px;
  font-weight: 700;
  line-height: 1.2;
}

.product-dialog .compact-field-required {
  color: #f87171;
  font-size: 14px;
  line-height: 1;
}

.product-dialog .product-form .el-form-item {
  margin-bottom: 14px;
  align-items: center;
}

.product-dialog .product-form .el-row {
  margin-left: 0 !important;
  margin-right: 0 !important;
}

.product-dialog .product-form .el-row > .el-col {
  padding-left: 6px !important;
  padding-right: 6px !important;
}

.product-dialog .product-form .el-form-item__content {
  flex: 1 1 auto !important;
  min-width: 0 !important;
  max-width: 100% !important;
}

.product-dialog .product-form .el-form-item__label {
  font-size: 13px;
  font-weight: 600;
  color: rgb(71 85 105);
  white-space: nowrap;
  line-height: 40px;
  padding-right: 10px;
}

.dark .product-dialog .product-form .el-form-item__label {
  color: rgb(148 163 184);
}

.dark .product-dialog .compact-field-label {
  color: rgb(203 213 225);
}

.product-dialog .product-form .el-form-item.align-start {
  align-items: flex-start;
}

.product-dialog .product-form .el-form-item.align-start .el-form-item__label {
  padding-top: 6px;
}

.product-dialog .node-card {
  padding: 14px;
  border: 1.5px solid rgb(226 232 240);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: white;
}

.dark .product-dialog .node-card {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
}

.product-dialog .node-card:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 12px rgba(59, 130, 246, 0.08);
}

.product-dialog .node-card.is-selected {
  border-color: var(--el-color-primary);
  background: rgba(59, 130, 246, 0.04);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .product-dialog .node-card.is-selected {
  background: rgba(59, 130, 246, 0.08);
}

.product-dialog .dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.product-dialog .product-form {
  margin-bottom: 0;
}

.product-dialog .product-form > :last-child {
  margin-bottom: 0 !important;
}

.product-dialog .btn-cancel {
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

.dark .product-dialog .btn-cancel {
  border-color: rgb(39 39 42);
  background: rgb(24 24 27);
  color: rgb(148 163 184);
}

.product-dialog .btn-cancel:hover {
  background: rgb(248 250 252);
}

.dark .product-dialog .btn-cancel:hover {
  background: rgb(39 39 42);
}

.product-dialog .btn-primary {
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

.product-dialog .btn-primary:hover:not(:disabled) {
  filter: brightness(1.1);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.3);
}

.product-dialog .btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.product-dialog .resource-type-group .el-radio-button__inner {
  font-size: 13px;
  font-weight: 600;
}

.product-dialog .product-form .el-input,
.product-dialog .product-form .el-input-number,
.product-dialog .product-form .el-select,
.product-dialog .product-form .el-textarea {
  width: 100% !important;
  max-width: 100% !important;
}

.product-dialog .el-dialog__header,
.product-dialog .el-dialog__body,
.product-dialog .el-dialog__footer,
.product-dialog .product-form,
.product-dialog .product-form .el-form,
.product-dialog .product-form .el-form-item,
.product-dialog .product-form .el-form-item__label,
.product-dialog .product-form .el-form-item__content,
.product-dialog .product-form .el-row,
.product-dialog .product-form .el-col {
  box-sizing: border-box;
  max-width: 100%;
  min-width: 0;
}

.product-dialog .product-form .el-input-number {
  --el-input-number-width: 100%;
}

@media (max-width: 1200px) {
  .product-dialog .el-row > .el-col[class*='el-col-sm-12'] {
    flex: 0 0 100% !important;
    max-width: 100% !important;
  }
}

@media (max-width: 768px) {
  .product-dialog .el-dialog__body {
    padding: 18px 16px 10px;
  }

  .product-dialog .form-section {
    max-width: 100%;
    padding: 14px 12px 4px;
  }

  .product-dialog .form-section--compact {
    padding: 12px 12px 12px;
  }

  .product-dialog .product-form .el-row > .el-col {
    padding-left: 4px !important;
    padding-right: 4px !important;
  }
}
</style>

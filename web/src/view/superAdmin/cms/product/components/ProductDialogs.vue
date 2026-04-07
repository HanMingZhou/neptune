<template>
  <BaseFormDialog
    v-if="productForm.productType === 1"
    v-model="editorVisibleModel"
    class="product-dialog product-dialog--compute"
    :model="productForm"
    :rules="computeRules"
    form-class="product-form"
    label-width="124px"
    :title="dialogTitle"
    :shell="false"
    width="840px"
    @submit="emit('submit-product')"
  >
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
            <el-select
              id="compute-cluster-id"
              v-model="productForm.clusterId"
              :placeholder="t('selectCluster')"
              class="w-full"
              @change="emit('cluster-change', $event)"
            >
              <el-option
                v-for="cluster in clusters"
                :key="cluster.id"
                :label="cluster.name"
                :value="cluster.id"
              />
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
        <div
          class="grid max-h-[280px] grid-cols-1 2xl:grid-cols-2 gap-3 overflow-y-auto pr-1"
          v-loading="loadingNodes"
        >
          <div
            v-if="clusterNodes.length === 0"
            class="2xl:col-span-2 py-12 text-center text-sm text-slate-400"
          >
            <span class="material-icons mb-2 block text-3xl opacity-30"
              >dns</span
            >
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
              <span
                class="mr-2 truncate font-mono text-sm font-bold text-slate-800 dark:text-slate-200"
                >{{ node.nodeName }}</span
              >
              <span
                class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-bold"
                :class="
                  node.schedulable
                    ? 'bg-emerald-500/10 text-emerald-500'
                    : 'bg-red-500/10 text-red-500'
                "
              >
                {{ node.schedulable ? t('schedulable') : t('unschedulable') }}
              </span>
            </div>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-slate-500">
              <span class="flex items-center gap-1"
                ><span class="material-icons text-[12px] text-blue-400"
                  >memory</span
                >CPU: {{ node.cpu }}核</span
              >
              <span class="flex items-center gap-1"
                ><span class="material-icons text-[12px] text-cyan-400"
                  >sd_card</span
                >{{ node.memory }}GB</span
              >
              <span
                v-if="node.gpuCount > 0"
                class="flex items-center gap-1 font-bold text-amber-600"
                ><span class="material-icons text-[12px]">developer_board</span
                >{{ node.gpuModel || 'GPU' }} × {{ node.gpuCount }}</span
              >
              <span
                v-if="node.vGpuNumber > 0"
                class="flex items-center gap-1 font-bold text-violet-600"
                ><span class="material-icons text-[12px]">grid_view</span>vGPU:
                {{ node.vGpuNumber }}</span
              >
              <span
                v-if="!node.gpuCount && !node.vGpuNumber"
                class="text-slate-400"
                >CPU Only</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>

    <template v-if="isEdit || productForm.nodeName">
      <div class="section-header">
        <span class="material-icons text-blue-500 text-[18px]"
          >inventory_2</span
        >
        <span>{{ t('productInfo') }}</span>
      </div>
      <div class="form-section">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('productName')"
              prop="name"
              for="compute-name"
            >
              <el-input
                id="compute-name"
                v-model="productForm.name"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('area')" prop="area" for="compute-area">
              <el-input
                id="compute-area"
                v-model="productForm.area"
                class="w-full"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('cpuCores')" prop="cpu" for="compute-cpu">
              <el-input-number
                id="compute-cpu"
                v-model="productForm.cpu"
                :min="1"
                :max="nodeMaxCpu"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('memoryGb')"
              prop="memory"
              for="compute-memory"
            >
              <el-input-number
                id="compute-memory"
                v-model="productForm.memory"
                :min="1"
                :max="nodeMaxMemory"
                class="w-full"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('cpuModel')" for="compute-cpu-model">
              <el-input
                id="compute-cpu-model"
                v-model="productForm.cpuModel"
                class="w-full"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="section-header">
        <span class="material-icons text-amber-500 text-[18px]"
          >developer_board</span
        >
        <span>{{ t('resourceType') }}</span>
      </div>
      <div class="form-section">
        <el-form-item :label="t('resourceType')" for="compute-resource-type">
          <el-radio-group
            id="compute-resource-type"
            v-model="resourceTypeModel"
            class="resource-type-group"
            @change="emit('resource-type-change', $event)"
          >
            <el-radio-button value="cpu">
              <span class="flex items-center gap-1.5"
                ><span class="material-icons text-[14px]">memory</span>CPU
                Only</span
              >
            </el-radio-button>
            <el-radio-button value="gpu">
              <span class="flex items-center gap-1.5"
                ><span class="material-icons text-[14px]">developer_board</span
                >GPU</span
              >
            </el-radio-button>
            <el-radio-button value="vgpu">
              <span class="flex items-center gap-1.5"
                ><span class="material-icons text-[14px]">grid_view</span
                >vGPU</span
              >
            </el-radio-button>
          </el-radio-group>
        </el-form-item>

        <template v-if="resourceType === 'gpu'">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('gpuModel')" for="compute-gpu-model">
                <el-input
                  id="compute-gpu-model"
                  v-model="productForm.gpuModel"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('gpuCount')" for="compute-gpu-count">
                <el-input-number
                  id="compute-gpu-count"
                  v-model="productForm.gpuCount"
                  :min="0"
                  :max="nodeMaxGpuCount"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('gpuMemory')" for="compute-gpu-memory">
                <el-input-number
                  id="compute-gpu-memory"
                  v-model="productForm.gpuMemory"
                  :min="0"
                  :max="nodeMaxGpuMemory"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </template>

        <template v-else-if="resourceType === 'vgpu'">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('vGpuCount')" for="compute-vgpu-count">
                <el-input-number
                  id="compute-vgpu-count"
                  v-model="productForm.vGpuCount"
                  :min="0"
                  :max="nodeMaxVGpuCount"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('vGpuMemory')" for="compute-vgpu-memory">
                <el-input-number
                  id="compute-vgpu-memory"
                  v-model="productForm.vGpuMemory"
                  :min="0"
                  :max="nodeMaxVGpuMemory"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item :label="t('vGpuCores')" for="compute-vgpu-cores">
                <el-input-number
                  id="compute-vgpu-cores"
                  v-model="productForm.vGpuCores"
                  :min="0"
                  :max="nodeMaxVGpuCores"
                  class="w-full"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </template>
      </div>

      <div class="section-header">
        <span class="material-icons text-emerald-500 text-[18px]"
          >payments</span
        >
        <span>{{ t('priceSettingsSingleCard') }}</span>
      </div>
      <div class="form-section">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('priceHourly')"
              prop="priceHourly"
              for="compute-price-hourly"
            >
              <el-input-number
                id="compute-price-hourly"
                v-model="productForm.priceHourly"
                :precision="2"
                :min="0"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceDaily')" for="compute-price-daily">
              <el-input-number
                id="compute-price-daily"
                v-model="productForm.priceDaily"
                :precision="2"
                :min="0"
                class="w-full"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('priceWeekly')" for="compute-price-weekly">
              <el-input-number
                id="compute-price-weekly"
                v-model="productForm.priceWeekly"
                :precision="2"
                :min="0"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('priceMonthly')"
              for="compute-price-monthly"
            >
              <el-input-number
                id="compute-price-monthly"
                v-model="productForm.priceMonthly"
                :precision="2"
                :min="0"
                class="w-full"
              />
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
          <el-switch
            id="compute-status"
            v-model="productForm.status"
            :active-value="1"
            :inactive-value="0"
            :active-text="t('onShelf')"
            :inactive-text="t('offShelf')"
          />
        </el-form-item>
        <el-form-item
          :label="t('paramDesc')"
          for="compute-description"
          class="align-start"
        >
          <el-input
            id="compute-description"
            v-model="productForm.description"
            type="textarea"
            :rows="3"
            :placeholder="t('inputProductDesc')"
          />
        </el-form-item>
      </div>
    </template>

    <template #footer="{ requestClose, submitForm }">
      <div class="dialog-footer">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!canSubmit || submitting"
          @click="submitForm"
        >
          {{ isEdit ? t('save') : t('create') }}
        </el-button>
      </div>
    </template>
  </BaseFormDialog>

  <BaseFormDialog
    v-else
    v-model="editorVisibleModel"
    class="product-dialog product-dialog--storage"
    :model="productForm"
    :rules="storageRules"
    form-class="product-form"
    label-width="96px"
    :title="dialogTitle"
    :shell="false"
    width="600px"
    @submit="emit('submit-product')"
  >
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
            <el-select
              id="storage-cluster-id"
              v-model="productForm.clusterId"
              :placeholder="t('selectCluster')"
              class="w-full"
              @change="emit('storage-cluster-change', $event)"
            >
              <el-option
                v-for="cluster in clusters"
                :key="cluster.id"
                :label="cluster.name"
                :value="cluster.id"
              />
            </el-select>
          </div>
        </el-form-item>
      </div>
    </div>

    <template v-if="isEdit || productForm.clusterId">
      <div class="section-header">
        <span class="material-icons text-purple-500 text-[18px]"
          >folder_open</span
        >
        <span>{{ t('productInfo') }}</span>
      </div>
      <div class="form-section">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('productName')"
              prop="name"
              for="storage-name"
            >
              <el-input
                id="storage-name"
                v-model="productForm.name"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('area')" prop="area" for="storage-area">
              <el-input
                id="storage-area"
                v-model="productForm.area"
                class="w-full"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('storageClass')"
              prop="storageClass"
              for="storage-class"
            >
              <el-input
                id="storage-class"
                v-model="productForm.storageClass"
                class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item
              :label="t('storagePriceGb')"
              prop="storagePriceGb"
              for="storage-price-gb"
            >
              <el-input-number
                id="storage-price-gb"
                v-model="productForm.storagePriceGb"
                :precision="4"
                :min="0"
                class="w-full"
              />
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
          <el-switch
            id="storage-status"
            v-model="productForm.status"
            :active-value="1"
            :inactive-value="0"
            :active-text="t('onShelf')"
            :inactive-text="t('offShelf')"
          />
        </el-form-item>
        <el-form-item
          :label="t('paramDesc')"
          for="storage-description"
          class="align-start"
        >
          <el-input
            id="storage-description"
            v-model="productForm.description"
            type="textarea"
            :rows="3"
            :placeholder="t('inputProductDesc')"
          />
        </el-form-item>
      </div>
    </template>

    <template #footer="{ requestClose, submitForm }">
      <div class="dialog-footer">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!canSubmit || submitting"
          @click="submitForm"
        >
          {{ isEdit ? t('save') : t('create') }}
        </el-button>
      </div>
    </template>
  </BaseFormDialog>

  <BaseFormDialog
    v-model="priceVisibleModel"
    class="product-dialog product-dialog--price"
    :model="priceForm"
    form-class="product-form"
    label-width="96px"
    :title="t('adjustPrice')"
    :shell="false"
    width="520px"
    @submit="emit('submit-price')"
  >
    <div class="section-header">
      <span class="material-icons text-amber-500 text-[18px]">payments</span>
      <span>{{ t('adjustPrice') }}</span>
    </div>
    <div class="form-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('priceHourly')" for="adjust-price-hourly">
            <el-input-number
              id="adjust-price-hourly"
              v-model="priceForm.priceHourly"
              :precision="2"
              :min="0"
              class="w-full"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('priceDaily')" for="adjust-price-daily">
            <el-input-number
              id="adjust-price-daily"
              v-model="priceForm.priceDaily"
              :precision="2"
              :min="0"
              class="w-full"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('priceWeekly')" for="adjust-price-weekly">
            <el-input-number
              id="adjust-price-weekly"
              v-model="priceForm.priceWeekly"
              :precision="2"
              :min="0"
              class="w-full"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('priceMonthly')" for="adjust-price-monthly">
            <el-input-number
              id="adjust-price-monthly"
              v-model="priceForm.priceMonthly"
              :precision="2"
              :min="0"
              class="w-full"
            />
          </el-form-item>
        </el-col>
      </el-row>
    </div>

    <template #footer="{ requestClose, submitForm }">
      <div class="dialog-footer">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="submitting"
          @click="submitForm"
        >
          {{ t('confirm') }}
        </el-button>
      </div>
    </template>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type {
  CmsClusterOption,
  CmsNodeRow,
  CmsProductForm,
  CmsProductPriceForm,
  CmsProductResourceType
} from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    canSubmit?: boolean
    clusterNodes?: CmsNodeRow[]
    clusters?: CmsClusterOption[]
    computeRules: FormRules<CmsProductForm>
    dialogTitle?: string
    isEdit?: boolean
    loadingNodes?: boolean
    modelValue?: boolean
    nodeMaxCpu?: number
    nodeMaxGpuCount?: number
    nodeMaxGpuMemory?: number
    nodeMaxMemory?: number
    nodeMaxVGpuCores?: number
    nodeMaxVGpuCount?: number
    nodeMaxVGpuMemory?: number
    priceForm: CmsProductPriceForm
    priceVisible?: boolean
    productForm: CmsProductForm
    resourceType?: CmsProductResourceType
    storageRules: FormRules<CmsProductForm>
    submitting?: boolean
  }>(),
  {
    canSubmit: false,
    clusterNodes: () => [],
    clusters: () => [],
    dialogTitle: '',
    isEdit: false,
    loadingNodes: false,
    modelValue: false,
    nodeMaxCpu: 256,
    nodeMaxGpuCount: 16,
    nodeMaxGpuMemory: 256,
    nodeMaxMemory: 2048,
    nodeMaxVGpuCores: 100,
    nodeMaxVGpuCount: 100,
    nodeMaxVGpuMemory: 256,
    priceVisible: false,
    resourceType: 'cpu',
    submitting: false
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: CmsProductForm['clusterId']]
  'node-select': [node: CmsNodeRow]
  'resource-type-change': [type: CmsProductResourceType]
  'storage-cluster-change': [clusterId: CmsProductForm['clusterId']]
  'submit-price': []
  'submit-product': []
  'update:modelValue': [value: boolean]
  'update:price-visible': [value: boolean]
  'update:resource-type': [value: CmsProductResourceType]
}>()

const t = inject<Translator>('t', (key: string) => key)

const editorVisibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const priceVisibleModel = computed({
  get: () => props.priceVisible,
  set: (value: boolean) => emit('update:price-visible', value)
})

const resourceTypeModel = computed({
  get: () => props.resourceType,
  set: (value: CmsProductResourceType) => emit('update:resource-type', value)
})
</script>

<style>
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

.dark .product-dialog .compact-field-label {
  color: rgb(203 213 225);
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

@media (max-width: 768px) {
  .product-dialog .form-section {
    max-width: 100%;
    padding: 14px 12px 4px;
  }

  .product-dialog .form-section--compact {
    padding: 12px 12px 12px;
  }
}
</style>

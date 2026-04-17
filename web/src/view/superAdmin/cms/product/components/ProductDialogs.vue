<template>
  <BaseFormDialog
    v-if="productForm.productType === 1"
    v-model="editorVisibleModel"
    class="product-dialog product-dialog--compute"
    :model="productForm"
    :rules="computeRules"
    form-class="product-form"
    :label-position="computeFormLabelPosition"
    :label-width="computeFormLabelWidth"
    :title="dialogTitle"
    :shell="false"
    :width="computeDialogWidth"
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

    <template v-if="isEdit || productForm.clusterId">
      <div
        :class="[
          'compute-dialog-layout',
          { 'compute-dialog-layout--create': !isEdit }
        ]"
      >
        <div class="compute-dialog-layout__main">
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
                <el-form-item
                  :label="t('cpuCores')"
                  prop="cpu"
                  for="compute-cpu"
                >
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
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item
                  :label="t('driverVersion')"
                  for="compute-driver-version"
                >
                  <el-input
                    id="compute-driver-version"
                    v-model="productForm.driverVersion"
                    class="w-full"
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item
                  :label="t('cudaVersion')"
                  for="compute-cuda-version"
                >
                  <el-input
                    id="compute-cuda-version"
                    v-model="productForm.cudaVersion"
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
            <el-form-item
              :label="t('resourceType')"
              for="compute-resource-type"
              class="resource-type-form-item"
            >
              <el-radio-group
                id="compute-resource-type"
                v-model="resourceTypeModel"
                class="resource-type-group"
                @change="emit('resource-type-change', $event)"
              >
                <el-radio-button value="cpu">
                  <span class="resource-type-option"
                    ><span class="material-icons text-[16px]">memory</span
                    ><span>CPU ONLY</span></span
                  >
                </el-radio-button>
                <el-radio-button value="gpu">
                  <span class="resource-type-option"
                    ><span class="material-icons text-[16px]"
                      >developer_board</span
                    ><span>GPU</span></span
                  >
                </el-radio-button>
                <el-radio-button value="vgpu">
                  <span class="resource-type-option"
                    ><span class="material-icons text-[16px]">grid_view</span
                    ><span>vGPU</span></span
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
                  <el-form-item
                    :label="t('gpuCountPerInstance')"
                    for="compute-gpu-count"
                  >
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
                  <el-form-item
                    :label="t('gpuMemoryPerCard')"
                    for="compute-gpu-memory"
                  >
                    <el-input-number
                      id="compute-gpu-memory"
                      v-model="productForm.gpuMemory"
                      :min="0"
                      :max="nodeMaxGpuMemory"
                      class="w-full"
                    />
                    <div class="form-field-hint">
                      {{ t('gpuMemoryPerCardHint') }}
                    </div>
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-else-if="resourceType === 'vgpu'">
              <el-row :gutter="20">
                <el-col :xs="24" :sm="12">
                  <el-form-item
                    :label="t('gpuModel')"
                    for="compute-vgpu-gpu-model"
                  >
                    <el-input
                      id="compute-vgpu-gpu-model"
                      v-model="productForm.gpuModel"
                      class="w-full"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12">
                  <el-form-item
                    :label="t('vGpuCount')"
                    for="compute-vgpu-count"
                  >
                    <el-input-number
                      id="compute-vgpu-count"
                      v-model="productForm.vGpuCount"
                      :min="0"
                      :max="nodeMaxVGpuCount"
                      class="w-full"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :xs="24" :sm="12">
                  <el-form-item
                    :label="t('vGpuMemory')"
                    for="compute-vgpu-memory"
                  >
                    <el-input-number
                      id="compute-vgpu-memory"
                      v-model="productForm.vGpuMemory"
                      :min="0"
                      :max="nodeMaxVGpuMemory"
                      class="w-full"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :sm="12">
                  <el-form-item
                    :label="t('vGpuCores')"
                    for="compute-vgpu-cores"
                  >
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
                    controls-position="right"
                    class="w-full price-input-number"
                  />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item
                  :label="t('priceDaily')"
                  for="compute-price-daily"
                >
                  <el-input-number
                    id="compute-price-daily"
                    v-model="productForm.priceDaily"
                    :precision="2"
                    :min="0"
                    controls-position="right"
                    class="w-full price-input-number"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="12">
                <el-form-item
                  :label="t('priceWeekly')"
                  for="compute-price-weekly"
                >
                  <el-input-number
                    id="compute-price-weekly"
                    v-model="productForm.priceWeekly"
                    :precision="2"
                    :min="0"
                    controls-position="right"
                    class="w-full price-input-number"
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
                    controls-position="right"
                    class="w-full price-input-number"
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
        </div>

        <div v-if="!isEdit" class="compute-dialog-layout__sidebar">
          <div class="section-header compute-dialog-layout__sidebar-heading">
            <span class="material-icons text-violet-500 text-[18px]">dns</span>
            <span>{{ t('nodeComparison') }}</span>
          </div>
          <div class="form-section form-section--comparison compute-side-stage">
            <el-form-item prop="nodeName" class="!m-0 !h-0 overflow-hidden">
              <el-input v-model="productForm.nodeName" class="hidden" />
            </el-form-item>

            <div class="comparison-summary-bar">
              <div class="comparison-summary-main">
                <span class="comparison-summary-pill">
                  {{ t('selectedNodeCount', { count: selectedNodeCount }) }}
                </span>
                <span class="comparison-summary-text">
                  {{ t('nodeSelectionHint') }}
                </span>
              </div>
              <button
                v-if="selectedNodeCount > 0"
                type="button"
                class="comparison-clear-button"
                @click="emit('selection-clear')"
              >
                {{ t('clearSelection') }}
              </button>
            </div>

            <div v-if="selectedNodeCount > 0" class="selected-node-strip">
              <button
                v-for="nodeName in selectedNodeNames"
                :key="nodeName"
                type="button"
                class="selected-node-chip"
                @click="handleSelectedNodeRemove(nodeName)"
              >
                <span class="material-icons">check_circle</span>
                <span class="selected-node-chip__label">{{ nodeName }}</span>
                <span class="material-icons">close</span>
              </button>
            </div>

            <div
              class="node-comparison-shell node-comparison-shell--stacked"
              v-loading="loadingNodes"
            >
              <div class="node-rail">
                <div class="compare-panel-header">
                  <div>
                    <div class="compare-panel-title">{{ t('nodeList') }}</div>
                    <div class="compare-panel-caption">
                      {{ t('clickPreviewNode') }}
                    </div>
                  </div>
                  <span class="compare-panel-badge">
                    {{ selectedNodeCount }}
                  </span>
                </div>

                <div v-if="clusterNodes.length === 0" class="node-rail-empty">
                  <span class="material-icons">dns</span>
                  <span>{{ t('nodeCandidateEmptyHint') }}</span>
                </div>

                <div v-else class="node-card-grid">
                  <article
                    v-for="node in clusterNodes"
                    :key="node.nodeName"
                    class="node-card"
                    :class="{
                      'is-selected': isNodeSelected(node.nodeName),
                      'is-preview': activePreviewNodeName === node.nodeName,
                      'is-disabled': !node.canCreateComputeProduct
                    }"
                    @click="emit('node-preview', node)"
                  >
                    <div class="node-card__header">
                      <div class="node-card__title-group">
                        <span class="node-card__title">{{
                          node.nodeName
                        }}</span>
                        <span class="node-card__meta">{{
                          node.internalIp || '--'
                        }}</span>
                      </div>
                      <button
                        type="button"
                        class="node-card__selector"
                        :class="{
                          'is-selected': isNodeSelected(node.nodeName)
                        }"
                        :disabled="!node.canCreateComputeProduct"
                        @click.stop="emit('node-toggle', node)"
                      >
                        <span class="material-icons">
                          {{
                            isNodeSelected(node.nodeName)
                              ? 'check_circle'
                              : 'add_circle'
                          }}
                        </span>
                        <span>
                          {{
                            isNodeSelected(node.nodeName)
                              ? t('cancelSelection')
                              : t('selectCurrentNode')
                          }}
                        </span>
                      </button>
                    </div>

                    <div class="node-card__badges">
                      <span
                        v-if="activePreviewNodeName === node.nodeName"
                        class="node-badge node-badge--info"
                      >
                        {{ t('previewingNode') }}
                      </span>
                      <span
                        class="node-badge"
                        :class="
                          node.schedulable
                            ? 'node-badge--success'
                            : 'node-badge--warning'
                        "
                      >
                        {{
                          node.schedulable
                            ? t('schedulable')
                            : t('unschedulable')
                        }}
                      </span>
                      <span
                        v-if="node.existingComputeProducts?.length"
                        class="node-badge node-badge--danger"
                      >
                        {{ t('occupiedByExistingProduct') }}
                      </span>
                    </div>

                    <div class="node-card__specs">
                      <span class="node-spec-pill">
                        <span class="material-icons">memory</span>
                        {{ node.cpu }}C
                      </span>
                      <span class="node-spec-pill">
                        <span class="material-icons">sd_card</span>
                        {{ node.memory }}GB
                      </span>
                      <span
                        v-if="node.gpuCount"
                        class="node-spec-pill node-spec-pill--gpu"
                      >
                        {{ node.gpuModel || t('gpu') }} x {{ node.gpuCount }}
                      </span>
                      <span
                        v-else-if="hasVGpuSpec(node)"
                        class="node-spec-pill node-spec-pill--vgpu"
                      >
                        vGPU {{ formatVGpuSpec(node) }}
                      </span>
                      <span v-else class="node-spec-pill node-spec-pill--cpu">
                        CPU ONLY
                      </span>
                    </div>

                    <div class="node-card__footer">
                      <span
                        class="node-card__status"
                        :class="{ 'is-danger': !node.canCreateComputeProduct }"
                      >
                        {{
                          node.canCreateComputeProduct
                            ? t('availableForCreate')
                            : node.disableReason ||
                              t('incompatibleWithCurrentConfig')
                        }}
                      </span>
                      <span class="node-card__preview-hint">
                        {{
                          activePreviewNodeName === node.nodeName
                            ? t('previewNode')
                            : t('clickPreviewNode')
                        }}
                      </span>
                    </div>
                  </article>
                </div>
              </div>

              <div class="compare-panel">
                <div class="compare-panel-header">
                  <div>
                    <div class="compare-panel-title">
                      {{ t('previewNode') }}
                    </div>
                    <div class="compare-panel-caption">
                      {{ t('existingComputeProducts') }}
                    </div>
                  </div>
                </div>

                <div v-if="previewNodeCandidate" class="compare-panel-body">
                  <div class="draft-compare-card">
                    <div class="draft-compare-card__header">
                      <div>
                        <div class="draft-compare-card__title">
                          {{ t('pendingCreateSpec') }}
                        </div>
                        <div class="draft-compare-card__subtitle">
                          {{ selectedNodeLabel }}
                        </div>
                      </div>
                      <span class="draft-compare-card__count">
                        {{
                          t('selectedNodeCount', { count: selectedNodeCount })
                        }}
                      </span>
                    </div>

                    <div class="draft-compare-card__chips">
                      <span
                        v-for="entry in draftSpecEntries"
                        :key="entry.key"
                        class="compare-mini-chip"
                        :class="
                          entry.tone ? `compare-mini-chip--${entry.tone}` : ''
                        "
                      >
                        {{ entry.label }} {{ entry.value }}
                      </span>
                    </div>

                    <div class="draft-compare-card__prices">
                      <span
                        v-for="entry in draftPriceEntries"
                        :key="entry.key"
                        class="draft-price-pill"
                      >
                        {{ entry.label }} ¥{{ entry.value }}
                      </span>
                    </div>
                  </div>

                  <div class="preview-node-hero">
                    <div>
                      <div class="preview-node-hero__title">
                        {{ previewNodeCandidate.nodeName }}
                      </div>
                      <div class="preview-node-hero__subtitle">
                        {{ previewNodeCandidate.clusterName || '--' }} ·
                        {{ previewNodeCandidate.area || '--' }}
                      </div>
                    </div>
                    <div class="preview-node-hero__chips">
                      <span class="preview-chip">
                        CPU {{ previewNodeCandidate.cpu }}
                      </span>
                      <span class="preview-chip">
                        MEM {{ previewNodeCandidate.memory }}GB
                      </span>
                      <span
                        v-if="previewNodeCandidate.gpuCount"
                        class="preview-chip preview-chip--gpu"
                      >
                        GPU {{ previewNodeCandidate.gpuCount }}
                      </span>
                      <span
                        v-else-if="hasVGpuSpec(previewNodeCandidate)"
                        class="preview-chip preview-chip--vgpu"
                      >
                        {{ formatVGpuSpec(previewNodeCandidate) }}
                      </span>
                    </div>
                  </div>

                  <div
                    v-if="previewNodeCandidate.disableReason"
                    class="compare-alert"
                  >
                    <span class="material-icons">info</span>
                    <span>{{ previewNodeCandidate.disableReason }}</span>
                  </div>

                  <div
                    v-if="previewExistingProducts.length"
                    class="compare-product-list"
                  >
                    <article
                      v-for="product in previewExistingProducts"
                      :key="product.id"
                      class="compare-product-card"
                    >
                      <div class="compare-product-card__header">
                        <div>
                          <div class="compare-product-card__title">
                            {{ product.name }}
                          </div>
                          <div class="compare-product-card__subtitle">
                            {{ product.resourceType.toUpperCase() }}
                          </div>
                        </div>
                        <span
                          class="node-badge"
                          :class="
                            product.status === 1
                              ? 'node-badge--success'
                              : 'node-badge--neutral'
                          "
                        >
                          {{
                            product.status === 1 ? t('onShelf') : t('offShelf')
                          }}
                        </span>
                      </div>

                      <div class="compare-product-card__specs">
                        <span class="compare-mini-chip"
                          >CPU {{ product.cpu }}</span
                        >
                        <span class="compare-mini-chip">
                          MEM {{ product.memory }}GB
                        </span>
                        <span
                          v-if="product.gpuCount"
                          class="compare-mini-chip compare-mini-chip--gpu"
                        >
                          {{ product.gpuModel || t('gpu') }} x
                          {{ product.gpuCount }}
                        </span>
                        <span
                          v-if="product.gpuMemory"
                          class="compare-mini-chip compare-mini-chip--gpu"
                        >
                          {{ product.gpuMemory }}GB/card
                        </span>
                        <span
                          v-if="
                            product.vGpuNumber ||
                            product.vGpuMemory ||
                            product.vGpuCores
                          "
                          class="compare-mini-chip compare-mini-chip--vgpu"
                        >
                          {{ formatExistingVGpu(product) }}
                        </span>
                      </div>

                      <div class="compare-product-card__metrics">
                        <span
                          >{{ t('priceHourly') }} ¥{{
                            product.priceHourly?.toFixed(2) || '0.00'
                          }}</span
                        >
                        <span
                          >{{ t('remainingInventory') }}
                          {{ product.available }}</span
                        >
                        <span
                          >{{ t('maxInstances') }}
                          {{ product.maxInstances }}</span
                        >
                      </div>
                    </article>
                  </div>

                  <div v-else class="compare-empty-state">
                    <span class="material-icons">inventory_2</span>
                    <span>{{ t('existingComputeProductEmpty') }}</span>
                  </div>
                </div>

                <div
                  v-else
                  class="compare-empty-state compare-empty-state--panel"
                >
                  <span class="material-icons">travel_explore</span>
                  <span>{{ t('previewNodePlaceholder') }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
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
          {{ submitButtonText }}
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
                controls-position="right"
                class="w-full price-input-number"
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
    width="640px"
    @submit="emit('submit-price')"
  >
    <div class="section-header">
      <span class="material-icons text-amber-500 text-[18px]">payments</span>
      <span>{{ t('adjustPrice') }}</span>
    </div>
    <div class="form-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24">
          <el-form-item :label="t('priceHourly')" for="adjust-price-hourly">
            <el-input-number
              id="adjust-price-hourly"
              v-model="priceForm.priceHourly"
              :precision="2"
              :min="0"
              controls-position="right"
              class="w-full price-input-number"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24">
          <el-form-item :label="t('priceDaily')" for="adjust-price-daily">
            <el-input-number
              id="adjust-price-daily"
              v-model="priceForm.priceDaily"
              :precision="2"
              :min="0"
              controls-position="right"
              class="w-full price-input-number"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24">
          <el-form-item :label="t('priceWeekly')" for="adjust-price-weekly">
            <el-input-number
              id="adjust-price-weekly"
              v-model="priceForm.priceWeekly"
              :precision="2"
              :min="0"
              controls-position="right"
              class="w-full price-input-number"
            />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="24">
          <el-form-item :label="t('priceMonthly')" for="adjust-price-monthly">
            <el-input-number
              id="adjust-price-monthly"
              v-model="priceForm.priceMonthly"
              :precision="2"
              :min="0"
              controls-position="right"
              class="w-full price-input-number"
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
  CmsExistingComputeProductSummary,
  CmsNodeRow,
  CmsProductNodeCandidate,
  CmsProductForm,
  CmsProductPriceForm,
  CmsProductResourceType
} from '@/types/superAdmin'
import { formatVGpuSpec, hasVGpuSpec } from '@/utils/vgpu'

const props = withDefaults(
  defineProps<{
    activePreviewNodeName?: string
    canSubmit?: boolean
    clusterNodes?: CmsProductNodeCandidate[]
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
    previewNode?: CmsNodeRow | CmsProductNodeCandidate | null
    priceForm: CmsProductPriceForm
    priceVisible?: boolean
    productForm: CmsProductForm
    resourceType?: CmsProductResourceType
    selectedNodeCount?: number
    selectedNodeLabel?: string
    selectedNodeNames?: string[]
    storageRules: FormRules<CmsProductForm>
    submitButtonText?: string
    submitting?: boolean
  }>(),
  {
    activePreviewNodeName: '',
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
    previewNode: null,
    priceVisible: false,
    resourceType: 'cpu',
    selectedNodeCount: 0,
    selectedNodeLabel: '',
    selectedNodeNames: () => [],
    submitButtonText: '',
    submitting: false
  }
)

const emit = defineEmits<{
  'cluster-change': [clusterId: CmsProductForm['clusterId']]
  'node-preview': [node: CmsProductNodeCandidate]
  'node-toggle': [node: CmsProductNodeCandidate]
  'selection-clear': []
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

const computeDialogWidth = computed(() =>
  props.isEdit ? '1040px' : 'min(1360px, calc(100vw - 40px))'
)

const computeFormLabelPosition = computed(() => (props.isEdit ? 'left' : 'top'))

const computeFormLabelWidth = computed(() => (props.isEdit ? '168px' : 'auto'))

const previewNodeCandidate = computed<CmsProductNodeCandidate | null>(
  () => (props.previewNode || null) as CmsProductNodeCandidate | null
)

const previewExistingProducts = computed<CmsExistingComputeProductSummary[]>(
  () => previewNodeCandidate.value?.existingComputeProducts || []
)

const isNodeSelected = (nodeName: string): boolean =>
  props.selectedNodeNames.includes(nodeName)

const getNodeCandidateByName = (
  nodeName: string
): CmsProductNodeCandidate | undefined =>
  props.clusterNodes.find((node) => node.nodeName === nodeName)

const handleSelectedNodeRemove = (nodeName: string): void => {
  const node = getNodeCandidateByName(nodeName)
  if (node) {
    emit('node-toggle', node)
  }
}

const formatExistingVGpu = (
  product: Partial<CmsExistingComputeProductSummary>
): string => {
  const parts = [
    product.vGpuNumber ? `${product.vGpuNumber}x` : '',
    product.vGpuMemory ? `${product.vGpuMemory}GB` : '',
    product.vGpuCores ? `${product.vGpuCores}%` : ''
  ].filter(Boolean)

  return parts.join(' / ') || 'vGPU'
}

const draftSpecEntries = computed(() => {
  const entries: Array<{
    key: string
    label: string
    value: string | number
    tone?: 'gpu' | 'vgpu' | 'cpu'
  }> = [
    {
      key: 'cpu',
      label: 'CPU',
      value: `${props.productForm.cpu || 0}C`
    },
    {
      key: 'memory',
      label: 'MEM',
      value: `${props.productForm.memory || 0}GB`
    }
  ]

  if (props.resourceType === 'gpu') {
    entries.push({
      key: 'gpuModel',
      label: t('gpu'),
      value: props.productForm.gpuModel || 'GPU',
      tone: 'gpu'
    })
    entries.push({
      key: 'gpuCount',
      label: t('gpuCountShort'),
      value: props.productForm.gpuCount || 0,
      tone: 'gpu'
    })
    entries.push({
      key: 'gpuMemory',
      label: t('gpuMemoryShort'),
      value: `${props.productForm.gpuMemory || 0}GB`,
      tone: 'gpu'
    })
  } else if (props.resourceType === 'vgpu') {
    entries.push({
      key: 'vgpuCount',
      label: 'vGPU',
      value: props.productForm.vGpuCount || props.productForm.vGpuNumber || 0,
      tone: 'vgpu'
    })
    entries.push({
      key: 'vgpuMemory',
      label: t('vGpuMemoryShort'),
      value: `${props.productForm.vGpuMemory || 0}GB`,
      tone: 'vgpu'
    })
    entries.push({
      key: 'vgpuCores',
      label: t('vGpuCoresShort'),
      value: `${props.productForm.vGpuCores || 0}%`,
      tone: 'vgpu'
    })
  } else {
    entries.push({
      key: 'resource',
      label: t('resourceType'),
      value: 'CPU ONLY',
      tone: 'cpu'
    })
  }

  return entries
})

const draftPriceEntries = computed(() => [
  {
    key: 'hourly',
    label: 'H',
    value: (props.productForm.priceHourly || 0).toFixed(2)
  },
  {
    key: 'daily',
    label: 'D',
    value: (props.productForm.priceDaily || 0).toFixed(2)
  },
  {
    key: 'weekly',
    label: 'W',
    value: (props.productForm.priceWeekly || 0).toFixed(2)
  },
  {
    key: 'monthly',
    label: 'M',
    value: (props.productForm.priceMonthly || 0).toFixed(2)
  }
])
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

.product-dialog .compute-dialog-layout {
  min-width: 0;
}

.product-dialog .compute-dialog-layout--create {
  display: grid;
  grid-template-columns: minmax(0, 520px) minmax(430px, 1fr);
  gap: 28px;
  align-items: start;
}

.product-dialog .compute-dialog-layout__main {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.product-dialog .compute-dialog-layout__sidebar {
  min-width: 0;
  position: sticky;
  top: 0;
  align-self: start;
}

.product-dialog .compute-dialog-layout__sidebar-heading {
  margin-top: 2px;
  margin-bottom: 10px;
}

.product-dialog .compute-side-stage {
  display: flex;
  flex-direction: column;
  padding: 14px;
  max-height: calc(100vh - 190px);
  overflow: hidden;
  border-color: rgba(59, 130, 246, 0.16);
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.11),
      transparent 34%
    ),
    linear-gradient(
      180deg,
      rgba(248, 250, 252, 0.98),
      rgba(255, 255, 255, 0.98)
    );
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.06);
}

.product-dialog .compute-side-stage .node-comparison-shell--stacked {
  min-height: 0;
  flex: 1;
}

.product-dialog .compute-dialog-layout--create .form-section {
  padding: 18px 18px 10px;
}

.product-dialog .compute-dialog-layout--create .el-form-item {
  margin-bottom: 16px;
}

.product-dialog .compute-dialog-layout--create .el-form-item__label {
  margin-bottom: 8px;
  padding: 0 0 2px !important;
  white-space: normal;
  line-height: 1.35;
}

.product-dialog .compute-dialog-layout--create .el-form-item__content {
  margin-left: 0 !important;
}

.product-dialog
  .compute-dialog-layout--create
  .compute-dialog-layout__main
  .el-form-item__content
  > :is(.el-input, .el-select, .el-input-number) {
  width: 100%;
}

.product-dialog
  .compute-dialog-layout--create
  .compute-dialog-layout__main
  .el-form-item__content
  > .el-textarea {
  width: 100%;
  max-width: 460px;
}

.dark .product-dialog .compute-side-stage {
  border-color: rgba(59, 130, 246, 0.26);
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.18),
      transparent 34%
    ),
    linear-gradient(180deg, rgba(24, 24, 27, 0.98), rgba(30, 41, 59, 0.96));
  box-shadow: none;
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

.product-dialog .form-section--comparison {
  padding: 16px;
}

.product-dialog .comparison-summary-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.product-dialog .comparison-summary-main {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.product-dialog .comparison-summary-pill {
  display: inline-flex;
  align-items: center;
  padding: 6px 10px;
  border-radius: 999px;
  background: linear-gradient(
    135deg,
    rgba(15, 23, 42, 0.92),
    rgba(37, 99, 235, 0.9)
  );
  color: white;
  font-size: 12px;
  font-weight: 700;
}

.product-dialog .comparison-summary-text {
  color: rgb(100 116 139);
  font-size: 12px;
}

.dark .product-dialog .comparison-summary-text {
  color: rgb(148 163 184);
}

.product-dialog .comparison-clear-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 30px;
  padding: 0 12px;
  border: 1px solid rgba(148, 163, 184, 0.35);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.92);
  color: rgb(71 85 105);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition:
    border-color 0.18s ease,
    color 0.18s ease,
    background-color 0.18s ease,
    transform 0.18s ease;
}

.product-dialog .comparison-clear-button:hover {
  transform: translateY(-1px);
  border-color: rgba(59, 130, 246, 0.42);
  color: var(--el-color-primary);
}

.dark .product-dialog .comparison-clear-button {
  border-color: rgba(82, 82, 91, 0.9);
  background: rgba(39, 39, 42, 0.94);
  color: rgb(203 213 225);
}

.dark .product-dialog .comparison-clear-button:hover {
  border-color: rgba(96, 165, 250, 0.62);
  color: rgb(191 219 254);
}

.product-dialog .selected-node-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}

.product-dialog .selected-node-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  max-width: 100%;
  height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(59, 130, 246, 0.16);
  border-radius: 999px;
  background: linear-gradient(
    135deg,
    rgba(239, 246, 255, 0.96),
    rgba(255, 255, 255, 0.96)
  );
  color: rgb(30 64 175);
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.product-dialog .selected-node-chip:hover {
  transform: translateY(-1px);
  border-color: rgba(59, 130, 246, 0.32);
  box-shadow: 0 8px 18px rgba(59, 130, 246, 0.12);
}

.product-dialog .selected-node-chip .material-icons {
  font-size: 14px;
}

.product-dialog .selected-node-chip__label {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dark .product-dialog .selected-node-chip {
  border-color: rgba(59, 130, 246, 0.24);
  background: linear-gradient(
    135deg,
    rgba(30, 41, 59, 0.96),
    rgba(15, 23, 42, 0.94)
  );
  color: rgb(191 219 254);
}

.product-dialog .node-comparison-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(0, 0.85fr);
  gap: 16px;
  align-items: stretch;
}

.product-dialog .node-comparison-shell--stacked {
  grid-template-columns: minmax(0, 1fr);
  gap: 14px;
}

.product-dialog .node-rail,
.product-dialog .compare-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  border: 1px solid rgb(226 232 240);
  border-radius: 16px;
  background: linear-gradient(
    180deg,
    rgba(248, 250, 252, 0.96),
    rgba(255, 255, 255, 0.96)
  );
  overflow: hidden;
}

.dark .product-dialog .node-rail,
.dark .product-dialog .compare-panel {
  border-color: rgb(39 39 42);
  background: linear-gradient(
    180deg,
    rgba(24, 24, 27, 0.96),
    rgba(39, 39, 42, 0.96)
  );
}

.product-dialog .compare-panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px 12px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
}

.dark .product-dialog .compare-panel-header {
  border-color: rgba(63, 63, 70, 0.9);
}

.product-dialog .compare-panel-title {
  font-size: 13px;
  font-weight: 800;
  color: rgb(30 41 59);
}

.dark .product-dialog .compare-panel-title {
  color: rgb(241 245 249);
}

.product-dialog .compare-panel-caption {
  margin-top: 4px;
  color: rgb(100 116 139);
  font-size: 12px;
  line-height: 1.4;
}

.dark .product-dialog .compare-panel-caption {
  color: rgb(148 163 184);
}

.product-dialog .compare-panel-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.12);
  color: var(--el-color-primary);
  font-size: 12px;
  font-weight: 700;
}

.product-dialog .node-rail-empty,
.product-dialog .compare-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  min-height: 180px;
  padding: 24px 20px;
  color: rgb(148 163 184);
  text-align: center;
  font-size: 12px;
}

.product-dialog .compare-empty-state--panel {
  flex: 1;
}

.product-dialog .node-rail-empty .material-icons,
.product-dialog .compare-empty-state .material-icons {
  font-size: 32px;
  opacity: 0.45;
}

.product-dialog .node-card-grid {
  display: grid;
  gap: 10px;
  padding: 12px;
  max-height: 360px;
  overflow-y: auto;
}

.product-dialog .node-comparison-shell--stacked .node-card-grid {
  max-height: min(38vh, 420px);
}

.product-dialog .node-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  border: 1px solid rgb(226 232 240);
  border-radius: 14px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.88);
  transition:
    border-color 0.18s ease,
    transform 0.18s ease,
    box-shadow 0.18s ease,
    background-color 0.18s ease;
}

.product-dialog .node-card:hover {
  transform: translateY(-1px);
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 10px 22px rgba(15, 23, 42, 0.08);
}

.product-dialog .node-card.is-preview {
  border-color: rgba(59, 130, 246, 0.55);
}

.product-dialog .node-card.is-selected {
  border-color: var(--el-color-primary);
  background: rgba(59, 130, 246, 0.05);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.08);
}

.product-dialog .node-card.is-disabled {
  background: rgba(248, 250, 252, 0.7);
}

.dark .product-dialog .node-card {
  border-color: rgb(63 63 70);
  background: rgba(24, 24, 27, 0.92);
}

.dark .product-dialog .node-card.is-disabled {
  background: rgba(39, 39, 42, 0.92);
}

.product-dialog .node-card__header,
.product-dialog .compare-product-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.product-dialog .node-card__title-group {
  min-width: 0;
}

.product-dialog .node-card__selector {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  min-height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(148, 163, 184, 0.3);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.92);
  color: rgb(51 65 85);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.01em;
  cursor: pointer;
  transition:
    border-color 0.18s ease,
    background-color 0.18s ease,
    color 0.18s ease,
    transform 0.18s ease,
    box-shadow 0.18s ease;
}

.product-dialog .node-card__selector .material-icons {
  font-size: 15px;
}

.product-dialog .node-card__selector:hover:not(:disabled) {
  transform: translateY(-1px);
  border-color: rgba(59, 130, 246, 0.44);
  color: var(--el-color-primary);
}

.product-dialog .node-card__selector.is-selected {
  border-color: rgba(59, 130, 246, 0.28);
  background: rgba(59, 130, 246, 0.08);
  color: var(--el-color-primary);
  box-shadow: inset 0 0 0 1px rgba(59, 130, 246, 0.08);
}

.product-dialog .node-card__selector:disabled {
  cursor: not-allowed;
  opacity: 0.48;
  transform: none;
}

.dark .product-dialog .node-card__selector {
  border-color: rgba(82, 82, 91, 0.9);
  background: rgba(39, 39, 42, 0.9);
  color: rgb(226 232 240);
}

.dark .product-dialog .node-card__selector:hover:not(:disabled) {
  border-color: rgba(96, 165, 250, 0.62);
  color: rgb(191 219 254);
}

.dark .product-dialog .node-card__selector.is-selected {
  border-color: rgba(59, 130, 246, 0.32);
  background: rgba(59, 130, 246, 0.18);
  color: rgb(191 219 254);
}

.product-dialog .node-card__title {
  display: block;
  font-family:
    ui-monospace, SFMono-Regular, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 13px;
  font-weight: 800;
  color: rgb(15 23 42);
  line-height: 1.35;
}

.dark .product-dialog .node-card__title {
  color: rgb(248 250 252);
}

.product-dialog .node-card__meta {
  display: block;
  margin-top: 3px;
  color: rgb(100 116 139);
  font-size: 11px;
}

.product-dialog .node-card__badges,
.product-dialog .preview-node-hero__chips,
.product-dialog .node-card__specs,
.product-dialog .compare-product-card__specs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.product-dialog .node-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 700;
  line-height: 1;
}

.product-dialog .node-badge--success {
  background: rgba(16, 185, 129, 0.12);
  color: rgb(5 150 105);
}

.product-dialog .node-badge--warning {
  background: rgba(245, 158, 11, 0.14);
  color: rgb(217 119 6);
}

.product-dialog .node-badge--danger {
  background: rgba(239, 68, 68, 0.12);
  color: rgb(220 38 38);
}

.product-dialog .node-badge--info {
  background: rgba(59, 130, 246, 0.12);
  color: rgb(37 99 235);
}

.product-dialog .node-badge--neutral {
  background: rgba(100, 116, 139, 0.14);
  color: rgb(71 85 105);
}

.dark .product-dialog .node-badge--info {
  background: rgba(59, 130, 246, 0.2);
  color: rgb(191 219 254);
}

.product-dialog .node-spec-pill,
.product-dialog .preview-chip,
.product-dialog .compare-mini-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 8px;
  border-radius: 10px;
  background: rgba(15, 23, 42, 0.06);
  color: rgb(51 65 85);
  font-size: 11px;
  font-weight: 600;
  line-height: 1;
}

.product-dialog .node-spec-pill .material-icons {
  font-size: 13px;
}

.product-dialog .node-spec-pill--gpu,
.product-dialog .preview-chip--gpu,
.product-dialog .compare-mini-chip--gpu {
  background: rgba(245, 158, 11, 0.12);
  color: rgb(180 83 9);
}

.product-dialog .node-spec-pill--vgpu,
.product-dialog .preview-chip--vgpu,
.product-dialog .compare-mini-chip--vgpu {
  background: rgba(124, 58, 237, 0.12);
  color: rgb(109 40 217);
}

.product-dialog .node-spec-pill--cpu {
  background: rgba(59, 130, 246, 0.08);
  color: rgb(37 99 235);
}

.dark .product-dialog .node-spec-pill,
.dark .product-dialog .preview-chip,
.dark .product-dialog .compare-mini-chip {
  background: rgba(63, 63, 70, 0.95);
  color: rgb(226 232 240);
}

.product-dialog .node-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.product-dialog .node-card__status {
  color: rgb(71 85 105);
  font-size: 12px;
  line-height: 1.4;
}

.product-dialog .node-card__status.is-danger {
  color: rgb(220 38 38);
}

.product-dialog .node-card__preview-hint {
  flex-shrink: 0;
  color: rgb(148 163 184);
  font-size: 11px;
  font-weight: 600;
  text-align: right;
}

.dark .product-dialog .node-card__preview-hint {
  color: rgb(113 113 122);
}

.product-dialog .node-card__action {
  border: none;
  border-radius: 999px;
  background: rgb(15 23 42);
  color: white;
  padding: 7px 12px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition:
    opacity 0.18s ease,
    transform 0.18s ease;
}

.product-dialog .node-card__action:hover:not(:disabled) {
  transform: translateY(-1px);
}

.product-dialog .node-card__action.is-selected {
  background: var(--el-color-primary);
}

.product-dialog .node-card__action:disabled {
  cursor: not-allowed;
  opacity: 0.4;
}

.product-dialog .compare-panel-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 12px;
  min-height: 0;
  padding: 14px 16px 16px;
}

.product-dialog .draft-compare-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 13px 14px;
  border: 1px solid rgba(59, 130, 246, 0.14);
  border-radius: 16px;
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.14),
      transparent 42%
    ),
    linear-gradient(
      135deg,
      rgba(248, 250, 252, 0.98),
      rgba(255, 255, 255, 0.98)
    );
}

.dark .product-dialog .draft-compare-card {
  border-color: rgba(59, 130, 246, 0.24);
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.22),
      transparent 42%
    ),
    linear-gradient(135deg, rgba(30, 41, 59, 0.96), rgba(15, 23, 42, 0.96));
}

.product-dialog .draft-compare-card__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.product-dialog .draft-compare-card__title {
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  color: rgb(30 64 175);
}

.product-dialog .draft-compare-card__subtitle {
  margin-top: 4px;
  color: rgb(71 85 105);
  font-size: 12px;
  line-height: 1.45;
}

.product-dialog .draft-compare-card__count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.08);
  color: rgb(15 23 42);
  font-size: 11px;
  font-weight: 700;
  white-space: nowrap;
}

.product-dialog .draft-compare-card__chips,
.product-dialog .draft-compare-card__prices {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.product-dialog .draft-price-pill {
  display: inline-flex;
  align-items: center;
  min-height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.06);
  color: rgb(51 65 85);
  font-size: 11px;
  font-weight: 700;
}

.dark .product-dialog .draft-compare-card__title {
  color: rgb(147 197 253);
}

.dark .product-dialog .draft-compare-card__subtitle {
  color: rgb(203 213 225);
}

.dark .product-dialog .draft-compare-card__count,
.dark .product-dialog .draft-price-pill {
  background: rgba(255, 255, 255, 0.08);
  color: rgb(226 232 240);
}

.product-dialog .preview-node-hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 14px;
  border-radius: 14px;
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.16),
      transparent 48%
    ),
    rgba(248, 250, 252, 0.92);
}

.dark .product-dialog .preview-node-hero {
  background:
    radial-gradient(
      circle at top right,
      rgba(59, 130, 246, 0.18),
      transparent 48%
    ),
    rgba(24, 24, 27, 0.92);
}

.product-dialog .preview-node-hero__title {
  font-size: 15px;
  font-weight: 800;
  color: rgb(15 23 42);
}

.dark .product-dialog .preview-node-hero__title {
  color: rgb(248 250 252);
}

.product-dialog .preview-node-hero__subtitle {
  margin-top: 4px;
  color: rgb(100 116 139);
  font-size: 12px;
}

.product-dialog .compare-alert {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(254, 242, 242, 0.96);
  color: rgb(185 28 28);
  font-size: 12px;
  line-height: 1.5;
}

.dark .product-dialog .compare-alert {
  background: rgba(69, 10, 10, 0.58);
  color: rgb(252 165 165);
}

.product-dialog .compare-alert .material-icons {
  font-size: 16px;
  margin-top: 1px;
}

.product-dialog .compare-product-list {
  display: grid;
  gap: 10px;
  max-height: 236px;
  overflow-y: auto;
}

.product-dialog .node-comparison-shell--stacked .compare-product-list {
  max-height: min(34vh, 360px);
}

.product-dialog .compare-product-card {
  padding: 12px;
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.9);
}

.dark .product-dialog .compare-product-card {
  border-color: rgba(63, 63, 70, 0.95);
  background: rgba(24, 24, 27, 0.88);
}

.product-dialog .compare-product-card__title {
  font-size: 13px;
  font-weight: 800;
  color: rgb(30 41 59);
}

.dark .product-dialog .compare-product-card__title {
  color: rgb(248 250 252);
}

.product-dialog .compare-product-card__subtitle {
  margin-top: 4px;
  color: rgb(100 116 139);
  font-size: 11px;
  letter-spacing: 0.04em;
}

.product-dialog .compare-product-card__metrics {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
  color: rgb(71 85 105);
  font-size: 11px;
  font-weight: 600;
}

.product-dialog .resource-type-form-item {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.product-dialog .resource-type-form-item .el-form-item__label {
  width: 100% !important;
  max-width: 100%;
  justify-content: flex-start;
  margin-bottom: 10px;
  padding: 0;
  font-size: 13px;
  font-weight: 700;
  line-height: 1.2;
  color: rgb(51 65 85);
}

.product-dialog .resource-type-form-item .el-form-item__content {
  display: flex;
  align-items: stretch;
  width: 100%;
  margin-left: 0 !important;
}

.product-dialog .resource-type-group {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
}

.product-dialog .resource-type-group .el-radio-button {
  width: 100%;
  margin-left: 0 !important;
}

.product-dialog .resource-type-group .el-radio-button__inner {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 44px;
  padding: 9px 14px;
  border: 1px solid rgb(226 232 240);
  border-radius: 12px !important;
  background: rgb(248 250 252);
  box-shadow: none;
  color: rgb(51 65 85);
  font-weight: 700;
  line-height: 1.2;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease;
}

.product-dialog .resource-type-option {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  text-align: center;
}

.product-dialog .resource-type-option .material-icons {
  flex-shrink: 0;
}

.product-dialog .resource-type-group .el-radio-button__inner:hover {
  border-color: rgba(59, 130, 246, 0.45);
  color: rgb(30 41 59);
}

.product-dialog
  .resource-type-group
  .el-radio-button.is-active
  .el-radio-button__inner {
  border-color: var(--el-color-primary);
  background: rgba(59, 130, 246, 0.08);
  color: var(--el-color-primary);
  box-shadow: 0 0 0 1px rgba(59, 130, 246, 0.14);
}

.dark .product-dialog .resource-type-group .el-radio-button__inner {
  border-color: rgb(63 63 70);
  background: rgb(39 39 42);
  color: rgb(226 232 240);
}

.dark .product-dialog .resource-type-form-item .el-form-item__label {
  color: rgb(203 213 225);
}

.dark .product-dialog .resource-type-group .el-radio-button__inner:hover {
  border-color: rgba(96, 165, 250, 0.6);
  color: white;
}

.dark
  .product-dialog
  .resource-type-group
  .el-radio-button.is-active
  .el-radio-button__inner {
  background: rgba(59, 130, 246, 0.16);
  color: rgb(191 219 254);
}

.product-dialog .dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.product-dialog .form-field-hint {
  margin-top: 6px;
  color: rgb(148 163 184);
  font-size: 12px;
  line-height: 1.2;
}

.dark .product-dialog .form-field-hint {
  color: rgb(113 113 122);
}

.product-dialog .product-form {
  margin-bottom: 0;
}

.product-dialog .price-input-number {
  width: 100%;
  max-width: 100%;
}

.product-dialog .price-input-number .el-input {
  min-width: 0;
}

.product-dialog .price-input-number .el-input__wrapper {
  padding-left: 12px;
  padding-right: 48px !important;
}

.product-dialog .price-input-number .el-input__inner {
  min-width: 0;
  overflow: hidden;
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.product-dialog .price-input-number .el-input-number__increase,
.product-dialog .price-input-number .el-input-number__decrease {
  width: 40px;
  right: 1px;
}

.product-dialog .el-form-item__label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-dialog .product-form > :last-child {
  margin-bottom: 0 !important;
}

@media (max-width: 1360px) {
  .product-dialog .compute-dialog-layout--create {
    grid-template-columns: minmax(0, 1fr);
    gap: 16px;
  }

  .product-dialog .compute-dialog-layout__sidebar {
    position: static;
    margin-top: 0;
  }

  .product-dialog .compute-side-stage {
    max-height: none;
  }
}

@media (max-width: 768px) {
  .product-dialog .form-section {
    max-width: 100%;
    padding: 14px 12px 4px;
  }

  .product-dialog .form-section--compact {
    padding: 12px 12px 12px;
  }

  .product-dialog .resource-type-group {
    grid-template-columns: 1fr;
  }

  .product-dialog .node-comparison-shell {
    grid-template-columns: 1fr;
  }

  .product-dialog .preview-node-hero,
  .product-dialog .draft-compare-card__header,
  .product-dialog .node-card__footer,
  .product-dialog .compare-product-card__header {
    flex-direction: column;
    align-items: stretch;
  }

  .product-dialog .node-card__selector,
  .product-dialog .comparison-clear-button {
    width: 100%;
  }

  .product-dialog .node-card__preview-hint {
    text-align: left;
  }
}
</style>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark pb-32">
    <PageHeader title-key="inference.createTitle" @back="handleCancel" />

    <div class="console-page-container px-6 py-6 space-y-6">
      <InferenceResourceSelectionSection
        :areas="areas"
        :change-filter="changeFilter"
        :filters="filters"
        :form="form"
        :format-price="formatPrice"
        :gpu-models-list="gpuModelsList"
        :pay-types="payTypes"
        :price-unit-text="priceUnitText"
        :products="products"
        @update:pay-type="form.payType = $event"
        @update:product-id="form.productId = $event"
      />

      <InferenceConfigurationSection
        :active-tab="activeTab"
        :add-env="addEnv"
        :add-mount="addMount"
        :auth-types="authTypes"
        :change-image-tab="changeImageTab"
        :deploy-types="deployTypes"
        :form="form"
        :framework-required="frameworkRequired"
        :frameworks="frameworks"
        :image-options="imageOptions"
        :image-tabs="imageTabs"
        :on-deploy-type-change="onDeployTypeChange"
        :pvcs="pvcs"
        :remove-env="removeEnv"
        :remove-mount="removeMount"
        @update:env="form.envs[$event.index][$event.key] = $event.value"
        @update:field="form[$event.key] = $event.value"
        @update:mount="form.mounts[$event.index][$event.key] = $event.value"
      />
    </div>

    <StickyActionBar
      :action-gap-class="'gap-8'"
      :can-submit="canCreate"
      :loading="loading"
      :price-label-class="'text-sm text-slate-500 font-medium'"
      :price-unit-text="priceUnitText"
      :price-value-class="'text-2xl font-black text-red-500'"
      :show-spinner="true"
      :submit-base-class="'px-8 py-2.5 rounded-lg text-sm font-bold transition-all flex items-center gap-2'"
      :submit-label-key="'inference.createService'"
      :total-price="totalPrice"
      :unit-label-class="'text-xs text-slate-400 font-bold uppercase tracking-wider'"
      @back="handleCancel"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/createPage/PageHeader.vue'
import StickyActionBar from '@/components/createPage/StickyActionBar.vue'
import InferenceConfigurationSection from './components/InferenceConfigurationSection.vue'
import InferenceResourceSelectionSection from './components/InferenceResourceSelectionSection.vue'
import { useInferenceCreate } from './composables/useInferenceCreate'

const t = inject('t', (key) => key)
const router = useRouter()

const {
  activeTab,
  addEnv,
  addMount,
  areas,
  authTypes,
  canCreate,
  changeFilter,
  changeImageTab,
  deployTypes,
  filters,
  form,
  formatPrice,
  frameworkRequired,
  frameworks,
  gpuModelsList,
  handleCancel,
  handleSubmit,
  imageOptions,
  imageTabs,
  loading,
  onDeployTypeChange,
  payTypes,
  priceUnitText,
  products,
  pvcs,
  removeEnv,
  removeMount,
  totalPrice
} = useInferenceCreate({ t, router })
</script>

<style scoped>
:deep(.el-form-item__label) {
  @apply text-slate-500 font-medium mb-1 !important;
}

:deep(.el-select-dropdown__item) {
  @apply py-2 h-auto leading-normal !important;
}

:deep(.el-select-dropdown__item.selected) {
  @apply font-bold !important;
}

:deep(.el-input-number) {
  @apply w-full !important;
}
</style>

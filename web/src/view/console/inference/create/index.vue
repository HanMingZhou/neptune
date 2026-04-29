<template>
  <div class="console-create-page">
    <PageHeader :title-key="pageTitleKey" @back="handleCancel" />

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
        :field-errors="fieldErrors"
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
        @update:field="updateField($event)"
        @update:mount="form.mounts[$event.index][$event.key] = $event.value"
        @validate:display-name="validateDisplayNameField"
      />
    </div>

    <StickyActionBar
      :action-gap-class="'gap-8'"
      :can-submit="canCreate"
      :loading="loading || editLoading"
      :price-label-class="'text-sm text-slate-500 font-medium'"
      :price-label-key="'totalPrice'"
      :price-unit-text="priceUnitText"
      :price-value-class="'text-2xl font-black text-red-500'"
      :show-spinner="true"
      :submit-base-class="'px-8 py-2.5 rounded-lg text-sm font-bold transition-all flex items-center gap-2'"
      :submit-label-key="submitLabelKey"
      :total-price="totalPrice"
      :unit-label-class="'text-xs text-slate-400 font-bold uppercase tracking-wider'"
      @back="handleCancel"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/createPage/PageHeader.vue'
import StickyActionBar from '@/components/createPage/StickyActionBar.vue'
import InferenceConfigurationSection from './components/InferenceConfigurationSection.vue'
import InferenceResourceSelectionSection from './components/InferenceResourceSelectionSection.vue'
import { useInferenceCreate } from './composables/useInferenceCreate'
import type { Translator } from '@/types/consoleResource'

const t = inject<Translator>('t', (key: string) => key)
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
  editLoading,
  fieldErrors,
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
  pageTitleKey,
  payTypes,
  priceUnitText,
  products,
  pvcs,
  removeEnv,
  removeMount,
  submitLabelKey,
  totalPrice,
  updateField,
  validateDisplayNameField
} = useInferenceCreate({ t, router })
</script>

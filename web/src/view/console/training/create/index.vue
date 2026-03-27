<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark pb-24">
    <PageHeader title-key="createTrainingJob" @back="goBack" />

    <div class="max-w-[1400px] mx-auto px-6 py-6 space-y-6">
      <TrainingResourceSelectionSection
        :areas="areas"
        :available-capacity="availableCapacity"
        :change-filter="changeFilter"
        :cpu-models="cpuModels"
        :decrease-worker="decreaseWorker"
        :filters="filters"
        :form="form"
        :format-price="formatPrice"
        :gpu-models="gpuModels"
        :increase-worker="increaseWorker"
        :max-worker-count="maxWorkerCount"
        :pay-types="payTypes"
        :price-unit-text="priceUnitText"
        :products="products"
        :select-product="selectProduct"
        :selected-product="selectedProduct"
        :show-worker-count="showWorkerCount"
        :total-price="totalPrice"
        :total-resources="totalResources"
      />

      <TrainingConfigurationSection
        :active-tab="activeTab"
        :add-env="addEnv"
        :add-mount="addMount"
        :change-image-tab="changeImageTab"
        :filtered-images="filteredImages"
        :form="form"
        :framework-types="frameworkTypes"
        :image-tabs="imageTabs"
        :insert-mpi-example="insertMpiExample"
        :insert-pytorch-example="insertPytorchExample"
        :pvcs="pvcs"
        :remove-env="removeEnv"
        :remove-mount="removeMount"
      />
    </div>

    <StickyActionBar
      :can-submit="canCreate"
      :loading="creating"
      :loading-label-key="'submitting'"
      :price-label-key="'configPrice'"
      :price-unit-text="priceUnitText"
      :submit-label-key="'submitOrder'"
      :total-price="totalPrice"
      @back="goBack"
      @submit="handleCreate"
    />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/createPage/PageHeader.vue'
import StickyActionBar from '@/components/createPage/StickyActionBar.vue'
import TrainingConfigurationSection from './components/TrainingConfigurationSection.vue'
import TrainingResourceSelectionSection from './components/TrainingResourceSelectionSection.vue'
import { useTrainingCreate } from './composables/useTrainingCreate'

const t = inject('t', (key) => key)
const router = useRouter()

const {
  activeTab,
  addEnv,
  addMount,
  areas,
  availableCapacity,
  canCreate,
  changeFilter,
  changeImageTab,
  cpuModels,
  creating,
  decreaseWorker,
  filteredImages,
  filters,
  form,
  formatPrice,
  frameworkTypes,
  goBack,
  gpuModels,
  handleCreate,
  imageTabs,
  increaseWorker,
  insertMpiExample,
  insertPytorchExample,
  maxWorkerCount,
  payTypes,
  priceUnitText,
  products,
  pvcs,
  removeEnv,
  removeMount,
  selectProduct,
  selectedProduct,
  showWorkerCount,
  totalPrice,
  totalResources
} = useTrainingCreate({ t, router })
</script>

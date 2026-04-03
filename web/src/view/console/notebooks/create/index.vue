<template>
  <div class="console-create-page pb-24">
    <PageHeader title-key="createInstance" @back="goBack" />

    <div class="console-page-container px-6 py-6 space-y-6">
      <NotebookResourceSelectionSection
        :areas="areas"
        :available-volumes="availableVolumes"
        :change-filter="changeFilter"
        :cpu-models="cpuModels"
        :filters="filters"
        :format-price="formatPrice"
        :gpu-count="gpuCount"
        :gpu-models="gpuModels"
        :on-volume-change="onVolumeChange"
        :pay-type="payType"
        :pay-types="payTypes"
        :price-unit-text="priceUnitText"
        :products="products"
        :selected-product="selectedProduct"
        :selected-volume-id="selectedVolumeId"
        :selected-volume-name="selectedVolumeName"
        :select-product="selectProduct"
        :total-price="totalPrice"
        :volume-mount-path="volumeMountPath"
        @update:gpu-count="gpuCount = $event"
        @update:pay-type="payType = $event"
        @update:selected-volume-id="selectedVolumeId = $event"
        @update:volume-mount-path="volumeMountPath = $event"
      />

      <NotebookConfigurationSection
        :active-tab="activeTab"
        :change-image-tab="changeImageTab"
        :enable-ssh-password="enableSshPassword"
        :enable-tensorboard="enableTensorboard"
        :field-errors="fieldErrors"
        :filtered-images="filteredImages"
        :image-tabs="imageTabs"
        :instance-name="instanceName"
        :selected-ssh-key="selectedSshKey"
        :selected-image="selectedImage"
        :ssh-keys="sshKeys"
        :tensorboard-log-path="tensorboardLogPath"
        @update:selected-image="selectedImage = $event"
        @update:enable-ssh-password="enableSshPassword = $event"
        @update:enable-tensorboard="enableTensorboard = $event"
        @update:instance-name="updateInstanceName($event)"
        @update:selected-ssh-key="selectedSshKey = $event"
        @update:tensorboard-log-path="updateTensorboardLogPath($event)"
        @validate:instance-name="validateInstanceNameField"
        @validate:tensorboard-log-path="validateTensorboardLogPathField"
      />
    </div>

    <StickyActionBar
      :can-submit="canCreate"
      :price-unit-text="priceUnitText"
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
import NotebookConfigurationSection from './components/NotebookConfigurationSection.vue'
import NotebookResourceSelectionSection from './components/NotebookResourceSelectionSection.vue'
import { useNotebookCreate } from './composables/useNotebookCreate'

const t = inject('t', (key) => key)
const router = useRouter()

const {
  activeTab,
  areas,
  availableVolumes,
  canCreate,
  changeFilter,
  changeImageTab,
  cpuModels,
  enableSshPassword,
  enableTensorboard,
  fieldErrors,
  filteredImages,
  filters,
  formatPrice,
  goBack,
  gpuCount,
  gpuModels,
  handleCreate,
  imageTabs,
  instanceName,
  onVolumeChange,
  payType,
  payTypes,
  priceUnitText,
  products,
  selectProduct,
  selectedImage,
  selectedProduct,
  selectedSshKey,
  selectedVolumeId,
  selectedVolumeName,
  sshKeys,
  tensorboardLogPath,
  totalPrice,
  updateInstanceName,
  updateTensorboardLogPath,
  validateInstanceNameField,
  validateTensorboardLogPathField,
  volumeMountPath
} = useNotebookCreate({ t, router })
</script>

<template>
  <el-date-picker
    v-model="rangeModel"
    type="daterange"
    size="small"
    value-format="YYYY-MM-DD"
    :range-separator="rangeSeparator"
    :start-placeholder="startPlaceholder"
    :end-placeholder="endPlaceholder"
    class="!w-[240px] custom-date-picker-standard"
    @change="$emit('change')"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  endPlaceholder: {
    type: String,
    default: ''
  },
  modelValue: {
    type: Array,
    default: () => []
  },
  rangeSeparator: {
    type: String,
    default: ''
  },
  startPlaceholder: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['change', 'update:modelValue'])

const rangeModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

<style scoped>
:deep(.custom-date-picker-standard) {
  border: none !important;
  box-shadow: none !important;
  background-color: transparent !important;
  padding: 0 !important;
  height: 32px !important;
}

:deep(.custom-date-picker-standard .el-input__wrapper) {
  background-color: rgb(248 250 252) !important;
  box-shadow: none !important;
  border: 1px solid rgb(226 232 240) !important;
  border-radius: 0.375rem !important;
  padding-left: 8px !important;
}

:deep(.dark .custom-date-picker-standard .el-input__wrapper) {
  background-color: rgb(24 24 27) !important;
  border-color: rgb(39 39 42) !important;
}

:deep(.custom-date-picker-standard .el-range-input) {
  font-size: 11px !important;
  color: rgb(71 85 105) !important;
  background: transparent !important;
}

:deep(.dark .custom-date-picker-standard .el-range-input) {
  color: rgb(148 163 184) !important;
}

:deep(.custom-date-picker-standard .el-range-separator) {
  font-size: 11px !important;
  color: rgb(148 163 184) !important;
  padding: 0 4px !important;
}
</style>

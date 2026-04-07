<template>
  <el-date-picker
    v-model="rangeModel"
    type="daterange"
    size="small"
    value-format="YYYY-MM-DD"
    :range-separator="rangeSeparator"
    :start-placeholder="startPlaceholder"
    :end-placeholder="endPlaceholder"
    class="!w-[240px] order-date-range-picker"
    @change="$emit('change')"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

type DateRangeValue = Array<string | number | Date>

const props = withDefaults(
  defineProps<{
    endPlaceholder?: string
    modelValue?: DateRangeValue
    rangeSeparator?: string
    startPlaceholder?: string
  }>(),
  {
    endPlaceholder: '',
    modelValue: () => [],
    rangeSeparator: '',
    startPlaceholder: ''
  }
)

const emit = defineEmits<{
  change: []
  'update:modelValue': [value: DateRangeValue]
}>()

const rangeModel = computed({
  get: () => props.modelValue,
  set: (value: DateRangeValue) => emit('update:modelValue', value)
})
</script>

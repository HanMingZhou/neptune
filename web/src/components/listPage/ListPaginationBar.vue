<template>
  <div
    v-if="!hideWhenEmpty || total > 0"
    class="console-pagination-bar text-xs text-slate-500"
    :class="justify === 'end' ? 'justify-end' : ''"
  >
    <span v-if="showTotal && !showPagination">{{ resolvedTotalText }}</span>
    <el-pagination
      v-if="showPagination"
      v-model:current-page="currentPageModel"
      v-model:page-size="pageSizeModel"
      :background="background"
      :layout="resolvedLayout"
      :page-sizes="pageSizes"
      :size="size"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  background: {
    type: Boolean,
    default: true
  },
  currentPage: {
    type: Number,
    default: 1
  },
  hideWhenEmpty: {
    type: Boolean,
    default: false
  },
  justify: {
    type: String,
    default: 'between'
  },
  layout: {
    type: String,
    default: 'sizes, prev, pager, next, jumper'
  },
  pageSize: {
    type: Number,
    default: 10
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  },
  showTotal: {
    type: Boolean,
    default: true
  },
  showPagination: {
    type: Boolean,
    default: true
  },
  size: {
    type: String,
    default: ''
  },
  total: {
    type: Number,
    default: 0
  },
  totalText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'current-change',
  'size-change',
  'update:current-page',
  'update:page-size'
])

const currentPageModel = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:current-page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:page-size', value)
})

const resolvedLayout = computed(() => {
  if (!props.showPagination) {
    return props.layout
  }

  const layoutItems = (props.layout || '')
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)

  if (!layoutItems.includes('total')) {
    layoutItems.unshift('total')
  }

  return layoutItems.join(', ')
})

const resolvedTotalText = computed(() => (
  props.totalText || `共 ${props.total} 条记录`
))

const handleCurrentChange = (value) => {
  emit('current-change', value)
}

const handleSizeChange = (value) => {
  emit('size-change', value)
}
</script>

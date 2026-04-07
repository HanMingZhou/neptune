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

<script setup lang="ts">
import { computed } from 'vue'

type PaginationAlign = 'between' | 'end'
type PaginationSize = '' | 'default' | 'small' | 'large'

const props = withDefaults(
  defineProps<{
    background?: boolean
    currentPage?: number
    hideWhenEmpty?: boolean
    justify?: PaginationAlign
    layout?: string
    pageSize?: number
    pageSizes?: number[]
    showTotal?: boolean
    showPagination?: boolean
    size?: PaginationSize
    total?: number
    totalText?: string
  }>(),
  {
    background: true,
    currentPage: 1,
    hideWhenEmpty: false,
    justify: 'between',
    layout: 'sizes, prev, pager, next, jumper',
    pageSize: 10,
    pageSizes: () => [10, 20, 50, 100],
    showTotal: true,
    showPagination: true,
    size: '',
    total: 0,
    totalText: ''
  }
)

const emit = defineEmits<{
  'current-change': [value: number]
  'size-change': [value: number]
  'update:current-page': [value: number]
  'update:page-size': [value: number]
}>()

const currentPageModel = computed({
  get: () => props.currentPage,
  set: (value: number) => emit('update:current-page', value)
})

const pageSizeModel = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:page-size', value)
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

const resolvedTotalText = computed(
  () => props.totalText || `共 ${props.total} 条记录`
)

const handleCurrentChange = (value: number) => {
  emit('current-change', value)
}

const handleSizeChange = (value: number) => {
  emit('size-change', value)
}
</script>

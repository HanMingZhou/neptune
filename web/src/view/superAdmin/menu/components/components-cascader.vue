<template>
  <div class="flex justify-between items-center gap-2 w-full">
    <el-cascader
      v-if="pathIsSelect"
      placeholder="请选择文件路径"
      :options="pathOptions"
      v-model="activeComponent"
      filterable
      class="!w-full"
      clearable
      @change="emitChange"
    />
    <el-input
      v-else
      v-model="tempPath"
      placeholder="页面:view/xxx/xx.vue 插件:plugin/xx/xx.vue"
      @change="emitChange"
    />
    <el-button @click="togglePathIsSelect"
      >{{ pathIsSelect ? '手动输入' : '快捷选择' }}
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import pathInfo from '@/pathInfo.json'

interface CascaderOption {
  value: string
  label: string
  children?: CascaderOption[]
}

const props = withDefaults(
  defineProps<{
    component?: string
  }>(),
  {
    component: ''
  }
)

const emits = defineEmits<{
  change: [value: string]
}>()

const pathOptions = ref<CascaderOption[]>([])
const tempPath = ref('')
const activeComponent = ref<string[]>([])
const pathIsSelect = ref(true)

const togglePathIsSelect = (): void => {
  if (pathIsSelect.value) {
    tempPath.value = activeComponent.value.join('/') || ''
  } else {
    activeComponent.value = tempPath.value.split('/').filter(Boolean)
  }

  pathIsSelect.value = !pathIsSelect.value
  emitChange()
}

function convertToCascaderOptions(
  data: Record<string, string>
): CascaderOption[] {
  const result: CascaderOption[] = []

  for (const filePath in data) {
    const label = data[filePath]
    const parts = filePath.split('/').filter(Boolean)
    const startIndex = parts[0] === 'src' ? 1 : 0

    let currentLevel = result

    for (let i = startIndex; i < parts.length; i += 1) {
      const part = parts[i]
      let node = currentLevel.find((item) => item.value === part)

      if (!node) {
        node = {
          value: part,
          label: part,
          children: []
        }
        currentLevel.push(node)
      }

      if (i === parts.length - 1) {
        node.label = label
        delete node.children
      }

      currentLevel = node.children || []
    }
  }

  return result
}

watch(
  () => props.component,
  (value) => {
    initCascader(value)
  }
)

onMounted(() => {
  pathOptions.value = convertToCascaderOptions(
    pathInfo as Record<string, string>
  )
  initCascader(props.component)
})

const initCascader = (value: string): void => {
  if (value === '') {
    pathIsSelect.value = true
    return
  }

  if ((pathInfo as Record<string, string>)[`/src/${value}`]) {
    activeComponent.value = value.split('/').filter(Boolean)
    tempPath.value = ''
    pathIsSelect.value = true
    return
  }

  tempPath.value = value
  activeComponent.value = []
  pathIsSelect.value = false
}

const emitChange = (): void => {
  emits(
    'change',
    pathIsSelect.value ? activeComponent.value.join('/') : tempPath.value
  )
}
</script>

<style scoped></style>

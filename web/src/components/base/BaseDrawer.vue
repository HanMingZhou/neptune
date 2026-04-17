<template>
  <el-drawer
    v-bind="forwardedAttrs"
    :model-value="modelValue"
    :title="title"
    :size="size"
    :direction="direction"
    :show-close="showClose"
    :class="drawerClasses"
    :before-close="handleBeforeClose"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    :destroy-on-close="destroyOnClose"
    @update:model-value="handleModelValueChange"
    @closed="emit('closed')"
    @opened="emit('opened')"
  >
    <template v-if="$slots.header" #header>
      <slot name="header" :request-close="requestClose" />
    </template>

    <slot :request-close="requestClose" />

    <template v-if="$slots.footer" #footer>
      <slot name="footer" :request-close="requestClose" />
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed, nextTick, useAttrs } from 'vue'

defineOptions({
  inheritAttrs: false
})

type DrawerDirection = 'rtl' | 'ltr' | 'ttb' | 'btt'

const props = withDefaults(
  defineProps<{
    closeOnClickModal?: boolean
    closeOnPressEscape?: boolean
    destroyOnClose?: boolean
    direction?: DrawerDirection
    modelValue?: boolean
    shell?: boolean
    showClose?: boolean
    size?: string | number
    title?: string
  }>(),
  {
    closeOnClickModal: true,
    closeOnPressEscape: true,
    destroyOnClose: false,
    direction: 'rtl',
    modelValue: false,
    shell: true,
    showClose: true,
    size: '30%',
    title: ''
  }
)

const emit = defineEmits<{
  close: []
  closed: []
  opened: []
  'update:modelValue': [value: boolean]
}>()

const attrs = useAttrs()

const forwardedAttrs = computed(() => {
  const { class: _className, ...rest } = attrs
  return rest
})

const drawerClasses = computed(() => [
  props.shell ? 'drawer-shell' : null,
  attrs.class
])

const requestClose = (): void => {
  emit('update:modelValue', false)
  emit('close')
}

const handleModelValueChange = (value: boolean): void => {
  emit('update:modelValue', value)
}

const handleBeforeClose = (done: () => void): void => {
  emit('update:modelValue', false)
  emit('close')
  void nextTick(() => {
    done()
  })
}
</script>

<template>
  <el-dialog
    v-bind="forwardedAttrs"
    :model-value="modelValue"
    :title="title"
    :width="width"
    :top="top"
    :align-center="alignCenter"
    :before-close="handleBeforeClose"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    :destroy-on-close="destroyOnClose"
    :show-close="showClose"
    :class="dialogClasses"
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
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, useAttrs } from 'vue'

defineOptions({
  inheritAttrs: false
})

type DialogShellSize = 'md' | 'lg' | 'xl'

const props = withDefaults(
  defineProps<{
    alignCenter?: boolean
    closeOnClickModal?: boolean
    closeOnPressEscape?: boolean
    destroyOnClose?: boolean
    modelValue?: boolean
    shell?: boolean
    shellSize?: DialogShellSize
    showClose?: boolean
    title?: string
    top?: string
    width?: string
  }>(),
  {
    alignCenter: false,
    closeOnClickModal: true,
    closeOnPressEscape: true,
    destroyOnClose: false,
    modelValue: false,
    shell: true,
    shellSize: undefined,
    showClose: true,
    title: '',
    top: undefined,
    width: undefined
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

const dialogClasses = computed(() => [
  props.shell ? 'dialog-shell' : null,
  props.shell && props.shellSize ? `dialog-shell--${props.shellSize}` : null,
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
  emit('close')
  done()
}
</script>

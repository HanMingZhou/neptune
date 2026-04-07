<template>
  <div class="console-resource-summary">
    <component
      :is="clickable ? 'button' : 'span'"
      :type="clickable ? 'button' : undefined"
      class="console-resource-summary__primary"
      :class="
        clickable
          ? 'font-bold text-primary hover:underline cursor-pointer text-sm'
          : ''
      "
      @click="handlePrimaryClick"
    >
      {{ primary }}
    </component>
    <span v-if="resolvedSecondary" class="console-resource-summary__secondary">
      {{ resolvedSecondary }}
    </span>
    <button
      v-if="copyValue"
      :title="copyTitle"
      class="console-resource-copy"
      type="button"
      @click.stop="emit('copy', copyValue)"
    >
      <span class="material-icons text-[12px]">content_copy</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    clickable?: boolean
    copyTitle?: string
    copyValue?: string
    primary?: string
    secondary?: string
    showSecondaryWhenSame?: boolean
  }>(),
  {
    clickable: true,
    copyTitle: '',
    copyValue: '',
    primary: '',
    secondary: '',
    showSecondaryWhenSame: false
  }
)

const emit = defineEmits<{
  copy: [value: string]
  'primary-click': []
}>()

const resolvedSecondary = computed(() => {
  if (!props.secondary) {
    return ''
  }

  if (!props.showSecondaryWhenSame && props.secondary === props.primary) {
    return ''
  }

  return props.secondary
})

const handlePrimaryClick = () => {
  if (!props.clickable) {
    return
  }

  emit('primary-click')
}
</script>

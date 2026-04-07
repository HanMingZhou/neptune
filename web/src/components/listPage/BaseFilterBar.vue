<template>
  <div v-if="plain" class="list-filter-bar" :class="barClass">
    <div v-if="wrapMain" class="list-toolbar-main" :class="mainClass">
      <slot />
    </div>
    <slot v-else />

    <div
      v-if="$slots.actions"
      class="list-toolbar-actions"
      :class="actionsClass"
    >
      <slot name="actions" />
    </div>
  </div>

  <div v-else class="console-filter-card" :class="cardClass">
    <div class="list-filter-bar" :class="barClass">
      <div v-if="wrapMain" class="list-toolbar-main" :class="mainClass">
        <slot />
      </div>
      <slot v-else />

      <div
        v-if="$slots.actions"
        class="list-toolbar-actions"
        :class="actionsClass"
      >
        <slot name="actions" />
      </div>
    </div>

    <slot name="footer" />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    actionsClass?: string
    barClass?: string
    cardClass?: string
    mainClass?: string
    plain?: boolean
    wrapMain?: boolean
  }>(),
  {
    actionsClass: '',
    barClass: '',
    cardClass: 'px-5 py-4',
    mainClass: '',
    plain: false,
    wrapMain: false
  }
)
</script>

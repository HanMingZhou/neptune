<template>
  <section class="console-management-shell">
    <div v-if="$slots.filters" class="console-management-shell__filters">
      <slot name="filters" />
    </div>

    <div class="console-management-shell__body">
      <slot />
    </div>
  </section>
</template>

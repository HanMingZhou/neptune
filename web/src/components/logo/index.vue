<template>
  <img :src="logoUrl" alt="Neptune Logo" class="neptune-logo" />
</template>

<script setup lang="ts">
import logoUrl from '@/assets/logo.png'

defineOptions({
  name: 'NeptuneLogo'
})
</script>

<style scoped>
.neptune-logo {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>

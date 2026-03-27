<template>
  <footer
    class="border-t py-12 transition-colors duration-300 bg-transparent"
    :class="isDark ? 'border-white/5' : 'border-slate-200/50'"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col md:flex-row justify-between items-center">
      <div class="mb-4 md:mb-0">
        <span class="font-[900] text-xl tracking-tight transition-colors duration-300 bg-gradient-to-br from-primary via-blue-500 to-indigo-600 bg-clip-text text-transparent">
          机器学习平台
        </span>
        <p class="text-gray-500 text-sm mt-2">© 2025 机器学习平台. All rights reserved.</p>
      </div>
      <div class="flex space-x-6">
        <a
          v-for="link in footerLinks"
          :key="link"
          href="#"
          class="transition-colors"
          :class="isDark ? 'text-gray-400 hover:text-white' : 'text-gray-600 hover:text-gray-900'"
        >
          {{ link }}
        </a>
      </div>
    </div>
  </footer>
</template>

<script setup>
defineProps({
  footerLinks: {
    type: Array,
    default: () => []
  },
  isDark: {
    type: Boolean,
    default: true
  }
})
</script>

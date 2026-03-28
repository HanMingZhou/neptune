<template>
  <nav
    class="fixed top-0 w-full z-50 transition-all duration-300 border-b"
    :class="isDark ? 'bg-[#0B0C15]/70 border-white/5 backdrop-blur-xl supports-[backdrop-filter]:bg-[#0B0C15]/50' : 'bg-white/70 border-slate-200/60 backdrop-blur-xl supports-[backdrop-filter]:bg-white/50'"
  >
    <div class="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
      <router-link to="/" class="flex items-center gap-4 group !no-underline">
        <div class="relative">
          <div class="absolute -inset-1.5 bg-gradient-to-tr from-indigo-500 to-cyan-500 rounded-xl blur opacity-30 group-hover:opacity-60 transition duration-700"></div>
          <div class="relative w-12 h-12 rounded-xl overflow-hidden border border-white/20 shadow-2xl bg-slate-950">
            <img :src="logoUrl" class="w-full h-full object-cover scale-[1.7] group-hover:scale-[1.9] transition duration-1000" alt="Logo" />
          </div>
        </div>
        <div class="flex flex-col">
          <span
            class="text-2xl font-[1000] tracking-tighter uppercase transition-colors"
            :class="isDark ? 'text-white' : 'text-slate-900'"
          >
            机器<span class="text-indigo-500 group-hover:text-cyan-400 transition-colors">学习</span>平台
          </span>
          <div class="flex items-center gap-1.5 -mt-1.5">
            <span class="text-[10px] font-black tracking-[5px] text-slate-500 uppercase">NEPTUNE</span>
            <div class="h-[1.5px] w-8 bg-indigo-500/30 group-hover:w-16 transition-all duration-500"></div>
          </div>
        </div>
      </router-link>

      <div class="flex items-center gap-4">
        <button
          class="flex items-center justify-center w-10 h-10 rounded-full transition-all duration-300 group border-none outline-none"
          :class="isDark ? 'bg-slate-700 hover:bg-slate-600' : 'bg-slate-200 hover:bg-slate-300'"
          @click="emit('toggle-theme')"
        >
          <el-icon
            :size="20"
            class="transition-transform duration-500 group-hover:rotate-180"
            :class="isDark ? 'text-yellow-400' : 'text-indigo-600'"
          >
            <Sunny v-if="isDark" />
            <Moon v-else />
          </el-icon>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import logoUrl from '@/assets/logo.png'

defineProps({
  isDark: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['toggle-theme'])
</script>

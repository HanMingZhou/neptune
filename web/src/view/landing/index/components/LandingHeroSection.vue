<template>
  <div class="relative pt-40 pb-20 sm:pt-48 sm:pb-32 overflow-hidden">
    <div class="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[1000px] pointer-events-none">
      <div
        class="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-[#0B0C15] dark:to-[#0B0C15] z-10"
        :class="isDark ? '' : 'to-[#F8FAFC]'"
      ></div>
      <div
        class="absolute top-[-20%] left-[20%] w-[600px] h-[600px] rounded-full mix-blend-multiply filter blur-[120px] opacity-20 animate-blob"
        :class="isDark ? 'bg-indigo-600' : 'bg-indigo-300'"
      ></div>
      <div
        class="absolute top-[-10%] right-[20%] w-[500px] h-[500px] rounded-full mix-blend-multiply filter blur-[120px] opacity-20 animate-blob animation-delay-2000"
        :class="isDark ? 'bg-cyan-600' : 'bg-cyan-300'"
      ></div>
    </div>

    <div class="max-w-5xl mx-auto px-6 text-center relative z-10">
      <div
        class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border mb-8 animate-fade-in-up"
        :class="isDark ? 'bg-indigo-500/10 border-indigo-500/20 text-indigo-300' : 'bg-indigo-50 border-indigo-100 text-indigo-600'"
      >
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 bg-indigo-500"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span>
        </span>
        <span class="text-xs font-semibold tracking-wide uppercase">AI Infrastructure v2.0</span>
      </div>

      <h1 class="text-5xl sm:text-7xl font-extrabold tracking-tight mb-8 leading-tight animate-fade-in-up animation-delay-100">
        <span class="block mb-2">释放 AI 无限潜能</span>
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-indigo-500 via-cyan-500 to-indigo-500 animate-gradient-x">
          构建未来算力
        </span>
      </h1>

      <p
        class="max-w-2xl mx-auto text-lg sm:text-xl mb-12 leading-relaxed animate-fade-in-up animation-delay-200"
        :class="isDark ? 'text-slate-400' : 'text-slate-600'"
      >
        一站式深度学习开发平台。从 Notebook 开发到大规模模型训练，<br class="hidden sm:block" />
        为您提供极致丝滑、安全可靠的算力获取体验。
      </p>

      <div class="flex flex-col sm:flex-row justify-center gap-5 animate-fade-in-up animation-delay-300">
        <router-link
          :to="isLoggedIn ? '/layout/dashboard' : '/login'"
          class="group relative inline-flex items-center justify-center px-8 py-4 text-base font-semibold text-white transition-all duration-200 bg-indigo-600 rounded-full hover:bg-indigo-700 hover:ring-4 hover:ring-indigo-500/30 shadow-lg shadow-indigo-500/30 !no-underline"
        >
          <span class="mr-2">{{ isLoggedIn ? '进入控制台' : '立即开始' }}</span>
          <svg class="w-5 h-5 transition-transform duration-200 group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </router-link>

        <button
          class="group inline-flex items-center justify-center px-8 py-4 text-base font-semibold transition-all duration-200 rounded-full border-none outline-none"
          :class="isDark ? 'bg-slate-700 text-white hover:bg-slate-600' : 'bg-slate-200 text-slate-800 hover:bg-slate-300'"
        >
          <span class="mr-2">浏览算力市场</span>
          <svg class="w-5 h-5 transition-transform duration-200 group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </button>
      </div>

      <div
        class="mt-20 grid grid-cols-2 gap-8 md:grid-cols-4 border-t pt-10 transition-colors duration-300"
        :class="isDark ? 'border-gray-800' : 'border-gray-100'"
      >
        <div v-for="stat in heroStats" :key="stat.label">
          <div class="text-3xl font-bold transition-colors duration-300" :class="isDark ? 'text-white' : 'text-gray-900'">
            {{ stat.value }}
          </div>
          <div class="text-sm text-gray-500 mt-1">{{ stat.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  heroStats: {
    type: Array,
    default: () => []
  },
  isDark: {
    type: Boolean,
    default: true
  },
  isLoggedIn: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
@keyframes gradient-x {
  0%,
  100% {
    background-size: 200% 200%;
    background-position: left center;
  }

  50% {
    background-size: 200% 200%;
    background-position: right center;
  }
}

.animate-gradient-x {
  animation: gradient-x 3s ease infinite;
}
</style>

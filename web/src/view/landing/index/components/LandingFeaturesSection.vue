<template>
  <div id="features" class="py-24 relative transition-colors duration-300 bg-transparent">
    <div class="absolute inset-x-0 -top-40 h-80 bg-gradient-to-b from-transparent via-indigo-500/[0.02] to-transparent pointer-events-none"></div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
      <div class="text-center mb-16">
        <h2 class="text-3xl font-bold sm:text-4xl transition-colors duration-300" :class="isDark ? 'text-white' : 'text-gray-900'">
          全流程 AI 开发支持
        </h2>
        <p class="mt-4 transition-colors duration-300" :class="isDark ? 'text-gray-400' : 'text-gray-600'">
          专为算法工程师打造的极致体验
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="group p-8 rounded-xl border transition-all hover:shadow-xl"
          :class="getCardClass(feature.accent)"
        >
          <div
            class="w-12 h-12 rounded-lg flex items-center justify-center mb-6 group-hover:scale-110 transition-transform"
            :class="getIconWrapperClass(feature.accent)"
          >
            <svg class="w-6 h-6" :class="getIconClass(feature.accent)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="feature.iconPath"></path>
            </svg>
          </div>
          <h3 class="text-xl font-bold mb-3 transition-colors duration-300" :class="isDark ? 'text-white' : 'text-gray-900'">
            {{ feature.title }}
          </h3>
          <p class="leading-relaxed transition-colors duration-300" :class="isDark ? 'text-gray-400' : 'text-gray-600'">
            {{ feature.description }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  features: {
    type: Array,
    default: () => []
  },
  isDark: {
    type: Boolean,
    default: true
  }
})

const getCardClass = (accent) => {
  const accentMap = {
    blue: 'hover:border-blue-500/30',
    pink: 'hover:border-pink-500/30',
    purple: 'hover:border-purple-500/30'
  }

  const baseClass = props.isDark ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'
  return `${baseClass} ${accentMap[accent] || accentMap.blue}`
}

const getIconWrapperClass = (accent) => {
  const accentMap = {
    blue: 'bg-blue-500/10',
    pink: 'bg-pink-500/10',
    purple: 'bg-purple-500/10'
  }

  return accentMap[accent] || accentMap.blue
}

const getIconClass = (accent) => {
  const accentMap = {
    blue: 'text-blue-500',
    pink: 'text-pink-500',
    purple: 'text-purple-500'
  }

  return accentMap[accent] || accentMap.blue
}
</script>

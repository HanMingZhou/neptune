<template>
  <div
    class="min-h-screen font-sans selection:bg-indigo-500/30 selection:text-indigo-900 dark:selection:bg-indigo-400/30 dark:selection:text-white transition-colors duration-500 relative"
    :class="isDark ? 'bg-[#0B0C15] text-white' : 'bg-[#F8FAFC] text-slate-900'"
  >
    <div v-if="isDark" class="absolute inset-0 z-0 pointer-events-none overflow-hidden">
      <div class="absolute top-0 right-0 w-[80%] h-[50%] bg-indigo-500/5 blur-[120px] rounded-full -translate-y-1/2 translate-x-1/4"></div>
      <div class="absolute bottom-0 left-0 w-[60%] h-[40%] bg-cyan-500/5 blur-[100px] rounded-full translate-y-1/4 -translate-x-1/4"></div>
    </div>

    <LandingNavbar
      :is-dark="isDark"
      @toggle-theme="toggleTheme"
    />

    <LandingHeroSection
      :hero-stats="heroStats"
      :is-dark="isDark"
      :is-logged-in="isLoggedIn"
    />

    <LandingFeaturesSection
      :features="featureCards"
      :is-dark="isDark"
    />

    <LandingFooter
      :footer-links="footerLinks"
      :is-dark="isDark"
    />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import LandingFeaturesSection from './components/LandingFeaturesSection.vue'
import LandingFooter from './components/LandingFooter.vue'
import LandingHeroSection from './components/LandingHeroSection.vue'
import LandingNavbar from './components/LandingNavbar.vue'
import { useLandingPage } from './composables/useLandingPage'

const {
  featureCards,
  footerLinks,
  heroStats,
  initialize,
  isDark,
  isLoggedIn,
  toggleTheme
} = useLandingPage()

onMounted(() => {
  initialize()
})
</script>

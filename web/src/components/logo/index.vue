<template>
  <div class="logo-wrapper">
    <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-svg">
      <!-- Outer circle with gradient -->
      <defs>
        <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#60a5fa;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#3b82f6;stop-opacity:1" />
        </linearGradient>
      </defs>
      
      <!-- Background circle -->
      <circle cx="20" cy="20" r="18" fill="url(#logoGradient)" opacity="0.1"/>
      
      <!-- Trident design -->
      <g stroke="url(#logoGradient)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <!-- Center prong -->
        <path d="M20 10 L20 30" />
        
        <!-- Left prong -->
        <path d="M14 10 L14 22 Q14 24 16 24" />
        
        <!-- Right prong -->
        <path d="M26 10 L26 22 Q26 24 24 24" />
        
        <!-- Top connections -->
        <path d="M14 10 L20 10 L26 10" />
        
        <!-- Circuit nodes -->
        <circle cx="14" cy="10" r="1.5" fill="url(#logoGradient)"/>
        <circle cx="20" cy="10" r="1.5" fill="url(#logoGradient)"/>
        <circle cx="26" cy="10" r="1.5" fill="url(#logoGradient)"/>
        
        <!-- Bottom handle -->
        <circle cx="20" cy="30" r="2" fill="url(#logoGradient)"/>
      </g>
    </svg>
  </div>
</template>

<script setup>
defineOptions({
  name: 'NeptuneLogo'
})
</script>

<style scoped>
.logo-wrapper {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-svg {
  width: 100%;
  height: 100%;
}
</style>

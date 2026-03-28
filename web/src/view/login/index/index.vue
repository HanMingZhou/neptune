<template>
  <div
    class="min-h-screen flex items-start lg:items-center justify-center px-4 py-8 sm:p-6 relative overflow-hidden transition-colors duration-500 bg-gradient-to-br from-indigo-900 via-purple-900 to-slate-900"
  >
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 rounded-full bg-indigo-500/20 blur-[120px] animate-blob"></div>
      <div class="absolute bottom-1/4 right-1/4 w-80 h-80 rounded-full bg-cyan-500/20 blur-[100px] animate-blob animation-delay-2000"></div>
      <div class="absolute top-1/2 left-1/2 w-72 h-72 rounded-full bg-purple-500/20 blur-[100px] animate-blob animation-delay-4000"></div>
    </div>

    <div class="relative z-10 w-full max-w-7xl mt-4 lg:mt-0">
      <div class="grid lg:grid-cols-2 gap-16 lg:gap-32 items-center">
        <LoginHeroPanel />

        <div class="flex items-center justify-center">
          <LoginAuthCard
            v-model:mfa-code="mfaCode"
            v-model:show-password="showPassword"
            :captcha-img="captchaImg"
            :error-msg="errorMsg"
            :field-errors="fieldErrors"
            :form="form"
            :is-register="isRegister"
            :loading="loading"
            :show-mfa-step="showMfaStep"
            @refresh-captcha="refreshCaptcha"
            @submit="handleSubmit"
            @toggle-mode="toggleMode"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import LoginAuthCard from './components/LoginAuthCard.vue'
import LoginHeroPanel from './components/LoginHeroPanel.vue'
import { useLoginPage } from './composables/useLoginPage'

const {
  captchaImg,
  errorMsg,
  fieldErrors,
  form,
  handleSubmit,
  initialize,
  isRegister,
  loading,
  mfaCode,
  refreshCaptcha,
  showMfaStep,
  showPassword,
  toggleMode
} = useLoginPage()

onMounted(() => {
  initialize()
})
</script>

<style scoped>
input:-webkit-autofill {
  -webkit-box-shadow: 0 0 0 30px transparent inset !important;
  -webkit-text-fill-color: inherit !important;
}

@keyframes blob {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }

  33% {
    transform: translate(30px, -50px) scale(1.1);
  }

  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}

@media (prefers-reduced-motion: reduce) {
  .animate-blob {
    animation: none;
  }

  .animation-delay-2000,
  .animation-delay-4000 {
    animation-delay: 0s;
  }
}
</style>

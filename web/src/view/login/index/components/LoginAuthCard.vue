<template>
  <div class="w-full max-w-md">
    <router-link to="/" class="lg:hidden inline-flex items-center gap-4 mb-8 !no-underline group">
      <div class="relative w-12 h-12 rounded-xl overflow-hidden border border-slate-200 shadow-xl">
        <img :src="logoUrl" class="w-full h-full object-cover scale-[1.7]" alt="Logo" />
      </div>
      <div class="flex flex-col">
        <span class="text-2xl font-black text-slate-900 tracking-tight">机器<span class="text-indigo-600">学习</span>平台</span>
      </div>
    </router-link>

    <div class="relative backdrop-blur-xl rounded-3xl p-8 border transition-all duration-300 shadow-2xl bg-white/85 border-white/50 shadow-black/20">
      <div class="mb-8">
        <h1 class="text-3xl font-bold mb-2 text-slate-900">
          {{ titleText }}
        </h1>
        <p class="text-sm text-slate-600">
          {{ descriptionText }}
        </p>
      </div>

      <form class="space-y-5" @submit.prevent="emit('submit')">
        <div v-if="!showMfaStep">
          <label class="block text-sm font-medium mb-2 text-slate-700">用户名</label>
          <input
            v-model="form.username"
            type="text"
            required
            placeholder="请输入用户名"
            class="w-full px-4 py-3 rounded-xl text-base transition-all duration-200 border-none outline-none focus:ring-2 focus:ring-indigo-500 bg-slate-100 text-slate-900 placeholder-slate-400"
          />
        </div>

        <div v-if="isRegister">
          <label class="block text-sm font-medium mb-2 text-slate-700">邮箱</label>
          <input
            v-model="form.email"
            type="email"
            required
            placeholder="请输入邮箱地址"
            class="w-full px-4 py-3 rounded-xl text-base transition-all duration-200 border-none outline-none focus:ring-2 focus:ring-indigo-500 bg-slate-100 text-slate-900 placeholder-slate-400"
          />
        </div>

        <div v-if="!showMfaStep">
          <label class="block text-sm font-medium mb-2 text-slate-700">密码</label>
          <div class="relative">
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="请输入密码"
              class="w-full px-4 py-3 rounded-xl text-base transition-all duration-200 border-none outline-none focus:ring-2 focus:ring-indigo-500 pr-12 bg-slate-100 text-slate-900 placeholder-slate-400"
            />
            <button
              type="button"
              class="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded transition-colors text-slate-500 hover:text-slate-700"
              @click="emit('update:showPassword', !showPassword)"
            >
              <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <div v-if="!isRegister && !showMfaStep">
          <label class="block text-sm font-medium mb-2 text-slate-700">验证码</label>
          <div class="flex gap-3">
            <input
              v-model="form.captcha"
              type="text"
              required
              placeholder="请输入验证码"
              class="flex-1 px-4 py-3 rounded-xl text-base transition-all duration-200 border-none outline-none focus:ring-2 focus:ring-indigo-500 bg-slate-100 text-slate-900 placeholder-slate-400"
            />
            <div
              class="w-28 h-12 rounded-xl cursor-pointer overflow-hidden flex items-center justify-center bg-slate-100"
              @click="emit('refresh-captcha')"
            >
              <img v-if="captchaImg" :src="captchaImg" alt="captcha" class="h-full object-contain" />
              <span v-else class="text-xs text-slate-500">加载中...</span>
            </div>
          </div>
        </div>

        <div v-if="showMfaStep">
          <div class="p-3 bg-indigo-500/10 border border-indigo-500/20 rounded-xl text-indigo-700 text-sm mb-4 flex items-center gap-2">
            <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
            </svg>
            请打开 Google Authenticator 输入6位验证码
          </div>
          <label class="block text-sm font-medium mb-2 text-slate-700">MFA 验证码</label>
          <input
            :value="mfaCode"
            type="text"
            required
            maxlength="6"
            placeholder="请输入6位验证码"
            class="w-full px-4 py-3 rounded-xl text-base transition-all duration-200 border-none outline-none focus:ring-2 focus:ring-indigo-500 bg-slate-100 text-slate-900 placeholder-slate-400 tracking-[0.3em] text-center font-mono text-lg"
            @input="emit('update:mfaCode', $event.target.value)"
          />
        </div>

        <div v-if="errorMsg" class="p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-600 text-sm">
          {{ errorMsg }}
        </div>

        <button
          type="submit"
          :disabled="loading || (showMfaStep && mfaCode.length !== 6)"
          class="w-full py-4 rounded-xl text-base font-semibold text-white bg-gradient-to-r from-indigo-600 to-cyan-600 hover:from-indigo-500 hover:to-cyan-500 transition-all duration-200 shadow-lg shadow-indigo-500/30 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          <svg v-if="loading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ submitText }}
        </button>
      </form>

      <div class="mt-6 text-center">
        <span class="text-sm text-slate-600">
          {{ isRegister ? '已有账号？' : '还没有账号？' }}
        </span>
        <button
          class="text-sm font-medium ml-1 transition-colors text-indigo-600 hover:text-indigo-700"
          @click="emit('toggle-mode')"
        >
          {{ isRegister ? '立即登录' : '立即注册' }}
        </button>
      </div>

      <div class="mt-4 text-center">
        <router-link
          to="/"
          class="text-sm transition-colors !no-underline text-slate-500 hover:text-slate-700"
        >
          ← 返回首页
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import logoUrl from '@/assets/logo.png'

const props = defineProps({
  captchaImg: {
    type: String,
    default: ''
  },
  errorMsg: {
    type: String,
    default: ''
  },
  form: {
    type: Object,
    required: true
  },
  isRegister: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  mfaCode: {
    type: String,
    default: ''
  },
  showMfaStep: {
    type: Boolean,
    default: false
  },
  showPassword: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'refresh-captcha',
  'submit',
  'toggle-mode',
  'update:mfaCode',
  'update:showPassword'
])

const titleText = computed(() => {
  if (props.showMfaStep) {
    return '安全验证'
  }

  return props.isRegister ? '创建账号' : '欢迎回来'
})

const descriptionText = computed(() => {
  if (props.showMfaStep) {
    return '您的账号已开启 MFA 双重验证'
  }

  return props.isRegister ? '开启您的 AI 开发之旅' : '登录以继续使用机器学习平台'
})

const submitText = computed(() => {
  if (props.loading) {
    return '处理中...'
  }

  if (props.showMfaStep) {
    return '验证并登录'
  }

  return props.isRegister ? '注册' : '登录'
})
</script>

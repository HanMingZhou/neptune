<template>
  <div class="hidden lg:flex flex-col justify-center">
    <router-link to="/" class="inline-flex items-center gap-5 mb-12 !no-underline group">
      <div class="relative">
        <div class="absolute -inset-1.5 bg-gradient-to-r from-indigo-500 to-cyan-500 rounded-2xl blur opacity-30 group-hover:opacity-60 transition duration-700"></div>
        <div class="relative w-16 h-16 rounded-2xl overflow-hidden border border-white/20 shadow-2xl bg-slate-950">
          <img :src="logoUrl" class="w-full h-full object-cover scale-[1.7] group-hover:scale-[1.9] transition duration-1000" alt="Logo" />
        </div>
      </div>
      <div class="flex flex-col">
        <span class="text-4xl font-[1000] tracking-tighter text-white uppercase">机器<span class="text-cyan-400">学习</span>平台</span>
        <span class="text-[11px] font-black tracking-[6px] text-slate-400 uppercase mt-1 opacity-80">AI INFRASTRUCTURE</span>
      </div>
    </router-link>

    <h1 class="text-5xl font-extrabold mb-6 leading-tight text-white">
      释放 AI 无限潜能<br />
      <span class="text-cyan-300">构建未来算力</span>
    </h1>

    <p class="text-xl mb-12 max-w-lg leading-relaxed text-slate-200">
      一站式深度学习开发平台。从 容器实例 开发到大规模模型训练，为您提供极致丝滑、安全可靠的算力获取体验。
    </p>

    <div class="space-y-6 max-w-lg">
      <div v-for="feature in features" :key="feature.title" class="flex items-start gap-4 group">
        <div
          class="w-12 h-12 rounded-xl backdrop-blur-sm flex items-center justify-center flex-shrink-0 transition-all duration-300 group-hover:scale-110 border"
          :class="feature.iconWrapperClass"
        >
          <svg class="w-6 h-6" :class="feature.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              :d="feature.iconPath"
            ></path>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-semibold mb-1 text-white">{{ feature.title }}</h3>
          <p class="text-sm text-slate-300">{{ feature.description }}</p>
        </div>
      </div>
    </div>

    <div class="mt-16 grid grid-cols-3 gap-8 max-w-lg">
      <div v-for="stat in stats" :key="stat.label">
        <div class="text-3xl font-bold" :class="stat.valueClass">{{ stat.value }}</div>
        <div class="text-sm mt-1 text-slate-400">{{ stat.label }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import logoUrl from '@/assets/logo.png'

const features = [
  {
    title: '云端 容器实例',
    description: '预置 PyTorch/TensorFlow 环境，开箱即用。支持 JupyterLab 与 SSH 直连。',
    iconClass: 'text-indigo-300',
    iconPath: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4',
    iconWrapperClass: 'bg-indigo-500/20 group-hover:bg-indigo-500/30 border-indigo-400/30'
  },
  {
    title: '弹性算力市场',
    description: '汇聚 H800, A100, 4090 等多种 GPU 资源。按秒计费，随用随停。',
    iconClass: 'text-purple-300',
    iconPath: 'M13 10V3L4 14h7v7l9-11h-7z',
    iconWrapperClass: 'bg-purple-500/20 group-hover:bg-purple-500/30 border-purple-400/30'
  },
  {
    title: '安全可靠',
    description: '99.9% 服务可用性，企业级数据安全保障，让您专注于 AI 创新。',
    iconClass: 'text-cyan-300',
    iconPath: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z',
    iconWrapperClass: 'bg-cyan-500/20 group-hover:bg-cyan-500/30 border-cyan-400/30'
  }
]

const stats = [
  { label: '高性能 GPU', value: 'H800', valueClass: 'text-indigo-300' },
  { label: '环境启动', value: '秒级', valueClass: 'text-purple-300' },
  { label: '服务可用性', value: '99.9%', valueClass: 'text-cyan-300' }
]
</script>

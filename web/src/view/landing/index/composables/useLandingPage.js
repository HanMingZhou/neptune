import { computed, ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'

const heroStats = [
  { label: '高性能 GPU', value: 'H800' },
  { label: '环境启动', value: '秒级' },
  { label: '极速网络', value: '10Gbps' },
  { label: '服务可用性', value: '99.9%' }
]

const featureCards = [
  {
    accent: 'blue',
    description: '预置 PyTorch/TensorFlow 环境，开箱即用。支持 JupyterLab 与 SSH 直连，无缝对接本地 IDE。',
    iconPath: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4',
    title: '云端容器实例'
  },
  {
    accent: 'purple',
    description: '汇聚 H800, A100, 4090 等多种 GPU 资源。按秒计费，随用随停，成本降低 50% 以上。',
    iconPath: 'M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z',
    title: '弹性算力市场'
  },
  {
    accent: 'pink',
    description: '支持模型一键发布为高可用 API 服务。内置自动扩缩容，轻松应对流量洪峰。',
    iconPath: 'M13 10V3L4 14h7v7l9-11h-7z',
    title: '一键模型部署'
  }
]

const footerLinks = ['关于我们', '服务条款', '隐私政策']

export function useLandingPage() {
  const userStore = useUserStore()
  const isDark = ref(true)

  const isLoggedIn = computed(() => !!userStore.token)

  const initialize = () => {
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      isDark.value = savedTheme === 'dark'
    }
  }

  const toggleTheme = () => {
    isDark.value = !isDark.value
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  }

  return {
    featureCards,
    footerLinks,
    heroStats,
    initialize,
    isDark,
    isLoggedIn,
    toggleTheme
  }
}

import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

const normalizeModuleId = (id: string) => id.replace(/\\/g, '/')

const manualChunks = (id: string) => {
  const moduleId = normalizeModuleId(id)

  if (!moduleId.includes('/node_modules/')) {
    return
  }

  if (moduleId.includes('/@element-plus/icons-vue/')) {
    return 'vendor-element-icons'
  }

  // Keep Element Plus and its tightly coupled runtime dependencies in one
  // chunk to avoid initialization order issues in production builds.
  if (
    moduleId.includes('/element-plus/') ||
    moduleId.includes('/@floating-ui/') ||
    moduleId.includes('/@popperjs/') ||
    moduleId.includes('/async-validator/') ||
    moduleId.includes('/dayjs/')
  ) {
    return 'vendor-element-plus'
  }

  if (moduleId.includes('/echarts/')) {
    return 'vendor-echarts'
  }

  if (moduleId.includes('/@xterm/')) {
    return 'vendor-xterm'
  }

  if (
    moduleId.includes('/ace-builds/') ||
    moduleId.includes('/vue3-ace-editor/')
  ) {
    return 'vendor-editor'
  }

  if (
    moduleId.includes('/marked/') ||
    moduleId.includes('/marked-highlight/') ||
    moduleId.includes('/highlight.js/')
  ) {
    return 'vendor-markdown'
  }

  if (
    moduleId.includes('/vue/') ||
    moduleId.includes('/@vue/') ||
    moduleId.includes('/vue-router/') ||
    moduleId.includes('/pinia/') ||
    moduleId.includes('/@vueuse/')
  ) {
    return 'vendor-vue'
  }

  return 'vendor'
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  return {
    server: {
      port: parseInt(env.VITE_CLI_PORT || '5173'),
      host: '0.0.0.0',
      proxy: {
        [env.VITE_BASE_API]: {
          target: `${env.VITE_BASE_PATH}:${env.VITE_SERVER_PORT}/`,
          changeOrigin: true,
          ws: true
        }
      }
    },
    plugins: [
      vue(),
      UnoCSS(),
      Components({
        dts: false,
        resolvers: [
          ElementPlusResolver({
            importStyle: false
          })
        ]
      })
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    build: {
      commonjsOptions: {
        include: [/node_modules/],
        transformMixedEsModules: true
      },
      rollupOptions: {
        output: {
          manualChunks
        }
      }
    },
    optimizeDeps: {
      include: ['vue'],
      esbuildOptions: {
        define: {
          global: 'globalThis'
        }
      }
    }
  }
})

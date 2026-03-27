import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import UnoCSS from 'unocss/vite';

// 空模块用于替换不兼容的 form-create 包
const emptyModule = 'export default {}; export const install = () => {};';
const normalizeModuleId = (id: string) => id.replace(/\\/g, '/');
const elementPlusChunkGroups: Record<string, string[]> = {
  'vendor-element-table': ['pagination', 'scrollbar', 'table', 'table-column', 'virtual-list'],
  'vendor-element-form': [
    'cascader',
    'checkbox',
    'checkbox-group',
    'form',
    'input',
    'input-number',
    'option',
    'option-group',
    'radio',
    'radio-group',
    'select',
    'select-v2',
    'switch',
    'upload'
  ],
  'vendor-element-overlay': [
    'dialog',
    'drawer',
    'dropdown',
    'loading',
    'message',
    'message-box',
    'notification',
    'popover',
    'popconfirm',
    'tooltip'
  ],
  'vendor-element-date': ['calendar', 'date-picker', 'time-picker', 'time-select']
};

const resolveElementPlusChunk = (moduleId: string) => {
  const componentMatch = moduleId.match(/\/element-plus\/es\/components\/([^/]+)\//);
  const componentName = componentMatch?.[1];

  if (componentName) {
    for (const [chunkName, componentNames] of Object.entries(elementPlusChunkGroups)) {
      if (componentNames.includes(componentName)) {
        return chunkName;
      }
    }
  }

  if (moduleId.includes('/@floating-ui/') || moduleId.includes('/@popperjs/')) {
    return 'vendor-element-overlay';
  }

  if (moduleId.includes('/async-validator/')) {
    return 'vendor-element-form';
  }

  if (moduleId.includes('/dayjs/')) {
    return 'vendor-element-date';
  }

  return 'vendor-element-plus';
};

const manualChunks = (id: string) => {
  const moduleId = normalizeModuleId(id);

  if (!moduleId.includes('/node_modules/')) {
    return;
  }

  if (moduleId.includes('/@element-plus/icons-vue/')) {
    return 'vendor-element-icons';
  }

  if (moduleId.includes('/element-plus/')) {
    return resolveElementPlusChunk(moduleId);
  }

  if (moduleId.includes('/echarts/')) {
    return 'vendor-echarts';
  }

  if (moduleId.includes('/@xterm/')) {
    return 'vendor-xterm';
  }

  if (moduleId.includes('/ace-builds/') || moduleId.includes('/vue3-ace-editor/')) {
    return 'vendor-editor';
  }

  if (
    moduleId.includes('/marked/') ||
    moduleId.includes('/marked-highlight/') ||
    moduleId.includes('/highlight.js/')
  ) {
    return 'vendor-markdown';
  }

  if (
    moduleId.includes('/vue/') ||
    moduleId.includes('/@vue/') ||
    moduleId.includes('/vue-router/') ||
    moduleId.includes('/pinia/') ||
    moduleId.includes('/@vueuse/')
  ) {
    return 'vendor-vue';
  }

  return 'vendor';
};

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());

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
      // 虚拟模块插件：替换不兼容的包
      {
        name: 'vite-plugin-form-create-stub',
        resolveId(id) {
          if (id === '@form-create/designer' || id === '@form-create/element-ui') {
            return '\0' + id;
          }
        },
        load(id) {
          if (id === '\0@form-create/designer' || id === '\0@form-create/element-ui') {
            return emptyModule;
          }
        }
      }
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      }
    },
    build: {
      commonjsOptions: {
        include: [/node_modules/],
        transformMixedEsModules: true
      },
      rollupOptions: {
        onwarn(warning, warn) {
          // 忽略 form-create 相关的警告
          if (warning.code === 'MISSING_EXPORT' && warning.exporter?.includes('form-create')) {
            return;
          }
          warn(warning);
        },
        output: {
          manualChunks
        },
        // 将 form-create 包标记为外部依赖
        external: [/@form-create\/.*/]
      }
    },
    optimizeDeps: {
      include: ['vue'],
      exclude: ['@form-create/designer', '@form-create/element-ui'],
      esbuildOptions: {
        define: {
          global: 'globalThis'
        }
      }
    }
  };
});

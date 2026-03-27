import { defineConfig } from '@unocss/vite';
import presetWind3 from '@unocss/preset-wind3';
import transformerDirectives from '@unocss/transformer-directives'

export default defineConfig({
    content: {
        filesystem: [
            'src/**/*.{vue,js,ts,jsx,tsx}',
            'index.html'
        ]
    },
    theme: {
        colors: {
            // 主色 - 来自模版
            primary: {
                DEFAULT: '#165DFF',
                hover: '#0E42D2',
                light: '#E8F3FF'
            },
            // 背景色
            background: {
                light: '#F2F3F5',
                dark: '#0F0F0F'
            },
            // 表面色
            surface: {
                light: '#FFFFFF',
                dark: '#1D1D1D'
            },
            // 边框色
            border: {
                light: '#E5E6EB',
                dark: '#333333'
            }
        },
        fontFamily: {
            sans: 'Inter, PingFang SC, Microsoft YaHei, sans-serif',
        }
    },
    presets: [
        presetWind3({ dark: 'class' })
    ],
    transformers: [
        transformerDirectives(),
    ],
})

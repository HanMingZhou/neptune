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
                DEFAULT: '#3158d4',
                hover: '#2445ad',
                light: '#ecf1ff'
            },
            // 背景色
            background: {
                light: '#f5f7fb',
                dark: '#0f1117'
            },
            // 表面色
            surface: {
                light: '#FFFFFF',
                dark: '#141923'
            },
            // 边框色
            border: {
                light: '#dbe2ee',
                dark: '#2e3647'
            }
        },
        fontFamily: {
            sans: 'Avenir Next, SF Pro Display, SF Pro Text, Segoe UI, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Noto Sans CJK SC, sans-serif',
        }
    },
    presets: [
        presetWind3({ dark: 'class' })
    ],
    transformers: [
        transformerDirectives(),
    ],
})

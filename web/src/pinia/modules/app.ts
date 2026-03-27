import { defineStore } from 'pinia'
import { ref, watchEffect, reactive } from 'vue'
import { setBodyPrimaryColor } from '@/utils/format'
import { useDark, usePreferredDark } from '@vueuse/core'

export interface AppConfig {
    weakness: boolean
    grey: boolean
    primaryColor: string
    showTabs: boolean
    darkMode: 'light' | 'dark' | 'auto'
    layout_side_width: number
    layout_side_collapsed_width: number
    layout_side_item_height: number
    show_watermark: boolean
    side_mode: string
    transition_type: string
    global_size: string
}

export const useAppStore = defineStore('app', () => {
    const device = ref('')
    const drawerSize = ref('')
    const operateMinWith = ref('240')
    const config = reactive<AppConfig>({
        weakness: false,
        grey: false,
        primaryColor: '#3b82f6',
        showTabs: false,
        darkMode: 'light',
        layout_side_width: 256,
        layout_side_collapsed_width: 80,
        layout_side_item_height: 48,
        show_watermark: false,
        side_mode: 'normal',
        transition_type: 'slide',
        global_size: 'default'
    })

    const isDark = useDark({
        selector: 'html',
        attribute: 'class',
        valueDark: 'dark',
        valueLight: 'light'
    })

    const preferredDark = usePreferredDark()

    const toggleTheme = (darkMode: boolean) => {
        isDark.value = darkMode
    }

    const toggleWeakness = (e: boolean) => {
        config.weakness = e
    }

    const toggleGrey = (e: boolean) => {
        config.grey = e
    }

    const togglePrimaryColor = (e: string) => {
        config.primaryColor = e
    }

    const toggleTabs = (e: boolean) => {
        config.showTabs = e
    }

    const toggleDevice = (e: string) => {
        if (e === 'mobile') {
            drawerSize.value = '100%'
            operateMinWith.value = '80'
        } else {
            drawerSize.value = '800'
            operateMinWith.value = '240'
        }
        device.value = e
    }

    const toggleDarkMode = (e: 'light' | 'dark' | 'auto') => {
        config.darkMode = e
    }

    watchEffect(() => {
        if (config.darkMode === 'auto') {
            isDark.value = preferredDark.value
            return
        }
        isDark.value = config.darkMode === 'dark'
    })

    const toggleConfigSideWidth = (e: number) => {
        config.layout_side_width = e
    }

    const toggleConfigSideCollapsedWidth = (e: number) => {
        config.layout_side_collapsed_width = e
    }

    const toggleConfigSideItemHeight = (e: number) => {
        config.layout_side_item_height = e
    }

    const toggleConfigWatermark = (e: boolean) => {
        config.show_watermark = e
    }

    const toggleSideMode = (e: string) => {
        config.side_mode = e
    }

    const toggleTransition = (e: string) => {
        config.transition_type = e
    }

    const toggleGlobalSize = (e: string) => {
        config.global_size = e
    }

    const baseConfig: AppConfig = {
        weakness: false,
        grey: false,
        primaryColor: '#3b82f6',
        showTabs: false,
        darkMode: 'light',
        layout_side_width: 256,
        layout_side_collapsed_width: 80,
        layout_side_item_height: 48,
        show_watermark: false,
        side_mode: 'normal',
        transition_type: 'slide',
        global_size: 'default'
    }

    const resetConfig = () => {
        for (let key in baseConfig) {
            (config as any)[key] = (baseConfig as any)[key]
        }
    }

    watchEffect(() => {
        document.documentElement.classList.toggle('html-weakenss', config.weakness)
        document.documentElement.classList.toggle('html-grey', config.grey)
    })

    watchEffect(() => {
        setBodyPrimaryColor(config.primaryColor, isDark.value ? 'dark' : 'light')
    })

    return {
        isDark,
        device,
        drawerSize,
        operateMinWith,
        config,
        toggleTheme,
        toggleDevice,
        toggleWeakness,
        toggleGrey,
        togglePrimaryColor,
        toggleTabs,
        toggleDarkMode,
        toggleConfigSideWidth,
        toggleConfigSideCollapsedWidth,
        toggleConfigSideItemHeight,
        toggleConfigWatermark,
        toggleSideMode,
        toggleTransition,
        resetConfig,
        toggleGlobalSize
    }
})

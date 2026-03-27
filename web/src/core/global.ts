import config from './config'
import { h, type App } from 'vue'

// 统一导入el-icon图标
import * as ElIconModules from '@element-plus/icons-vue'
import svgIcon from '@/components/svgIcon/svgIcon.vue'

const createIconComponent = (name: string) => ({
    name: 'SvgIcon',
    render() {
        return h(svgIcon, {
            localIcon: name
        })
    }
})

const registerIcons = async (app: App) => {
    const iconModules = import.meta.glob('@/assets/icons/**/*.svg')
    const pluginIconModules = import.meta.glob(
        '@/plugin/**/assets/icons/**/*.svg'
    )
    const mergedIconModules = Object.assign({}, iconModules, pluginIconModules)
    const allKeys: string[] = []
    for (const path in mergedIconModules) {
        let pluginName = ''
        if (path.startsWith('/src/plugin/')) {
            pluginName = `${path.split('/')[3]}-`
        }
        const iconName = path
            .split('/')
            .pop()!
            .replace(/\.svg$/, '')
        if (iconName.indexOf(' ') !== -1) {
            console.error(`icon ${iconName}.svg includes whitespace in ${path}`)
            continue
        }
        const key = `${pluginName}${iconName}`
        const iconComponent = createIconComponent(key)
        config.logs.push({
            key: key,
            label: key
        })
        app.component(key, iconComponent)

        allKeys.push(key)
    }
}

export const register = (app: App) => {
    for (const iconName in ElIconModules) {
        app.component(iconName, (ElIconModules as any)[iconName])
    }
    app.component('SvgIcon', svgIcon)
    registerIcons(app)
    app.config.globalProperties.$GIN_VUE_ADMIN = config
}

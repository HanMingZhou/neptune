import config from './config'
import { h, type App } from 'vue'

import svgIcon from '@/components/svgIcon/svgIcon.vue'

type IconModuleMap = Record<string, string>

const customIconModules = import.meta.glob('@/assets/icons/**/*.svg', {
  eager: true,
  import: 'default',
  query: '?raw'
}) as IconModuleMap

const pluginIconModules = import.meta.glob(
  '@/plugin/**/assets/icons/**/*.svg',
  {
    eager: true,
    import: 'default',
    query: '?raw'
  }
) as IconModuleMap

const createIconComponent = (name: string, svgMarkup: string) => ({
  name: 'SvgIcon',
  render() {
    return h('span', {
      class: 'app-inline-svg-icon',
      innerHTML: svgMarkup,
      role: 'img',
      'aria-label': name,
      style: {
        display: 'inline-flex',
        width: '1em',
        height: '1em',
        lineHeight: '1'
      }
    })
  }
})

const registerIcons = async (app: App) => {
  const mergedIconModules = Object.assign(
    {},
    customIconModules,
    pluginIconModules
  )
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
    const svgMarkup = mergedIconModules[path]
    if (!svgMarkup || svgMarkup.trim() === '') {
      console.error(`icon ${iconName}.svg is empty in ${path}`)
      continue
    }

    const iconComponent = createIconComponent(key, svgMarkup)
    config.logs.push({
      key: key,
      label: key
    })
    app.component(key, iconComponent)
  }
}

export const register = (app: App) => {
  app.component('SvgIcon', svgIcon)
  registerIcons(app)
  app.config.globalProperties.$GIN_VUE_ADMIN = config
}

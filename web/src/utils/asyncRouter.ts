const viewModules = import.meta.glob('../view/**/*.vue')
const pluginModules = import.meta.glob('../plugin/**/*.vue')

export const asyncRouterHandle = (asyncRouter: any[]) => {
    asyncRouter.forEach((item) => {
        if (item.component && typeof item.component === 'string') {
            item.meta.path = '/src/' + item.component
            if (item.component.split('/')[0] === 'view') {
                item.component = dynamicImport(viewModules, item.component)
            } else if (item.component.split('/')[0] === 'plugin') {
                item.component = dynamicImport(pluginModules, item.component)
            }
        }
        if (item.children) {
            asyncRouterHandle(item.children)
        }
    })
}

function dynamicImport(dynamicViewsModules: Record<string, any>, component: string) {
    const keys = Object.keys(dynamicViewsModules)
    const matchKeys = keys.filter((key) => {
        const k = key.replace('../', '')
        return k === component
    })
    let matchKey = matchKeys[0]
    if (!matchKey) {
        if (component === 'view/console/account/records.vue') {
            // @ts-ignore
            matchKey = keys.find(k => k.includes('view/console/account/accesslog.vue'))
        }
        if (component === 'view/console/account/settings.vue' || component === 'view/person/person.vue') {
            // @ts-ignore
            matchKey = keys.find(k => k.includes('view/console/account/person.vue'))
        }
    }

    return dynamicViewsModules[matchKey]
}

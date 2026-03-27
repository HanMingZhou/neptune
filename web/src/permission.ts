import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'
import getPageTitle from '@/utils/page'
import router from '@/router'
import Nprogress from 'nprogress'
import 'nprogress/nprogress.css'
import type { RouteLocationNormalized, RouteRecordRaw } from 'vue-router'

// 配置 NProgress
Nprogress.configure({
    showSpinner: false,
    ease: 'ease',
    speed: 500
})

// 白名单路由（无需登录即可访问）
const WHITE_LIST = ['Login', 'Landing']

function isExternalUrl(val: any): boolean {
    return typeof val === 'string' && /^(https?:)?\/\//.test(val)
}

function normalizeAbsolutePath(p: string): string {
    const s = '/' + String(p || '')
    return s.replace(/\/+/g, '/')
}

function normalizeRelativePath(p: string): string {
    return String(p || '').replace(/^\/+/, '')
}

function addTopLevelIfAbsent(r: RouteRecordRaw): void {
    if (!router.hasRoute(r.name as string)) {
        router.addRoute(r)
    }
}

interface RouteWithMeta {
    path?: string
    name?: string | symbol
    component?: any
    children?: RouteWithMeta[]
    meta?: {
        defaultMenu?: boolean
        [key: string]: any
    }
    parent?: RouteWithMeta
    [key: string]: any
}

function addRouteByChildren(route: RouteWithMeta | undefined, segments: string[] = [], parentName: string | null = null): void {
    if (isExternalUrl(route?.path) || isExternalUrl(route?.name) || isExternalUrl(route?.component)) {
        return
    }

    if (route?.name === 'layout') {
        route.children?.forEach((child) => addRouteByChildren(child, [], null))
        return
    }

    if (route?.meta?.defaultMenu === true && parentName === null) {
        const fullPath = [...segments, route.path].filter(Boolean).join('/')
        const children = route.children ? [...route.children] : []
        const newRoute: any = { ...route, path: fullPath }
        delete newRoute.children
        delete newRoute.parent
        newRoute.path = normalizeAbsolutePath(newRoute.path)

        if (router.hasRoute(newRoute.name)) return
        addTopLevelIfAbsent(newRoute)

        if (children.length) {
            children.forEach((child) => addRouteByChildren(child, [], newRoute.name))
        }
        return
    }

    if (route?.children && route.children.length) {
        const nextSegments = isExternalUrl(route.path) ? segments : [...segments, route.path || '']
        route.children.forEach((child) => addRouteByChildren(child, nextSegments, parentName))
        return
    }

    const fullPath = [...segments, route?.path || ''].filter(Boolean).join('/')
    const newRoute: any = { ...route, path: fullPath }
    delete newRoute.children
    delete newRoute.parent
    newRoute.path = normalizeRelativePath(newRoute.path)

    if (parentName) {
        if (newRoute.component || newRoute.children?.length) {
            router.addRoute(parentName, newRoute)
        }
    } else {
        if (newRoute.component || newRoute.children?.length) {
            router.addRoute('layout', newRoute)
        }
    }
}

const setupRouter = async (userStore: ReturnType<typeof useUserStore>): Promise<boolean> => {
    try {
        const routerStore = useRouterStore()
        await Promise.all([routerStore.SetAsyncRouter(), userStore.GetUserInfo()])

        const baseRouters = routerStore.asyncRouters || []
        const layoutRoute = baseRouters[0]
        if (layoutRoute?.name === 'layout' && !router.hasRoute('layout')) {
            const bareLayout = { ...layoutRoute, children: [] }
            router.addRoute(bareLayout as any)
        }

        const toRegister: RouteWithMeta[] = []
        if (layoutRoute?.children?.length) {
            toRegister.push(...(layoutRoute.children as RouteWithMeta[]))
        }
        if (baseRouters.length > 1) {
            baseRouters.slice(1).forEach((r) => {
                if (r?.name !== 'layout') toRegister.push(r as RouteWithMeta)
            })
        }
        toRegister.forEach((r) => addRouteByChildren(r, [], null))

        return true
    } catch (error) {
        console.error('Setup router failed:', error)
        try {
            await userStore.ClearStorage()
        } catch (e) {
            console.error('Clear storage failed:', e)
        }
        return false
    }
}

const removeLoading = (): void => {
    const element = document.getElementById('gva-loading-box')
    element?.remove()
}

const handleKeepAlive = async (to: RouteLocationNormalized): Promise<void> => {
    if (!to.matched.some((item) => item.meta.keepAlive)) return

    if (to.matched?.length > 2) {
        for (let i = 1; i < to.matched.length; i++) {
            const element = to.matched[i - 1]

            if (element.name === 'layout') {
                to.matched.splice(i, 1)
                await handleKeepAlive(to)
                continue
            }

            if (typeof element.components?.default === 'function') {
                await (element.components.default as () => Promise<any>)()
                await handleKeepAlive(to)
            }
        }
    }
}

const handleRedirect = (to: RouteLocationNormalized, userStore: ReturnType<typeof useUserStore>) => {
    const defaultRouter = userStore.userInfo?.authority?.defaultRouter
    console.log('handleRedirect check:', {
        defaultRouter,
        hasRoute: defaultRouter ? router.hasRoute(defaultRouter) : false,
        allRoutes: router.getRoutes().map(r => r.name)
    })
    if (defaultRouter && router.hasRoute(defaultRouter)) {
        return { name: defaultRouter, replace: true }
    }
    return { path: '/layout/404' }
}

router.beforeEach(async (to, from) => {
    const userStore = useUserStore()
    const routerStore = useRouterStore()
    const token = userStore.token

    Nprogress.start()

        ; (to.meta as any).matched = [...to.matched]
    await handleKeepAlive(to)

    document.title = getPageTitle((to.meta.title as string) || '', to)

    if (to.meta.client) {
        return true
    }

    if (WHITE_LIST.includes(to.name as string)) {
        if (to.name === 'Login' && token) {
            if (routerStore.asyncRouterFlag && userStore.userInfo?.authority?.defaultRouter) {
                return { name: userStore.userInfo.authority.defaultRouter, replace: true }
            }
            return { name: 'Landing', replace: true }
        }
        return true
    }

    if (token) {
        if (sessionStorage.getItem('needToHome') === 'true') {
            sessionStorage.removeItem('needToHome')
            return { path: '/' }
        }

        if (!routerStore.asyncRouterFlag) {
            const setupSuccess = await setupRouter(userStore)

            if (!setupSuccess || !userStore.token) {
                await userStore.ClearStorage()
                return {
                    name: 'Login',
                    query: { redirect: to.fullPath }
                }
            }

            return handleRedirect(to, userStore)
        }

        if (to.matched.length && (userStore.userInfo as any)?.authorityId) {
            return true
        }

        if (!(userStore.userInfo as any)?.authorityId) {
            await userStore.ClearStorage()
            return {
                name: 'Login',
                query: { redirect: to.fullPath }
            }
        }

        return { path: '/layout/404' }
    }

    return {
        name: 'Login',
        query: {
            redirect: to.fullPath
        }
    }
})

router.afterEach(() => {
    document.querySelector('.main-cont.main-right')?.scrollTo(0, 0)
    Nprogress.done()
})

router.onError((error) => {
    console.error('Router error:', error)
    Nprogress.remove()
})

removeLoading()

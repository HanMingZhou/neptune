import { asyncRouterHandle } from '@/utils/asyncRouter'
import { emitter } from '@/utils/bus'
import { asyncMenu } from '@/api/menu'
import { defineStore } from 'pinia'
import { ref, watchEffect } from 'vue'
import pathInfo from '@/pathInfo.json'
import { useRoute, type RouteLocationNormalizedLoaded } from "vue-router"
import { config } from "@/core/config"
import router from "@/router/index"
import type { RouteRecordRaw } from 'vue-router'

export interface RouteItem {
    path: string
    name: string
    component?: string | any
    meta: {
        title: string
        path?: string
        btns?: any
        hidden?: boolean
        defaultMenu?: boolean
        keepAlive?: boolean
        closeTab?: boolean
        [key: string]: any
    }
    children?: RouteItem[]
    parent?: RouteItem
    hidden?: boolean
    btns?: any
    sort?: number
    [key: string]: any
}

const notLayoutRouterArr: RouteItem[] = []
const keepAliveRoutersArr: string[] = []
const nameMap: Record<string, string> = {}

const formatRouter = (routes: RouteItem[] | undefined, routeMap: Record<string, RouteItem>, parent?: RouteItem) => {
    routes &&
        routes.forEach((item) => {
            item.parent = parent
            item.meta.btns = item.btns
            item.meta.hidden = item.hidden
            if (item.meta.defaultMenu === true) {
                if (!parent) {
                    item = { ...item, path: `/${item.path}` }
                    notLayoutRouterArr.push(item)
                }
            }
            routeMap[item.name] = item
            if (item.children && item.children.length > 0) {
                formatRouter(item.children, routeMap, item)
            }
        })
}

const KeepAliveFilter = (routes: RouteItem[] | undefined) => {
    routes &&
        routes.forEach((item) => {
            if (
                (item.children && item.children.some((ch) => ch.meta.keepAlive)) ||
                item.meta.keepAlive
            ) {
                const path = item.meta.path
                if (path && (pathInfo as any)[path]) {
                    keepAliveRoutersArr.push((pathInfo as any)[path])
                    nameMap[item.name] = (pathInfo as any)[path]
                }
            }
            if (item.children && item.children.length > 0) {
                KeepAliveFilter(item.children)
            }
        })
}

export const useRouterStore = defineStore('router', () => {
    const keepAliveRouters = ref<string[]>([])
    const asyncRouterFlag = ref(0)
    const setKeepAliveRouters = (history: any[]) => {
        const keepArrTemp: string[] = []

        keepArrTemp.push(...keepAliveRoutersArr)
        if (config.keepAliveTabs) {
            history.forEach((item) => {
                const routeInfo = routeMap[item.name]
                if (routeInfo && routeInfo.meta && routeInfo.meta.path) {
                    const componentName = (pathInfo as any)[routeInfo.meta.path]
                    if (componentName) {
                        keepArrTemp.push(componentName)
                    }
                }

                if (nameMap[item.name]) {
                    keepArrTemp.push(nameMap[item.name])
                }
            })
        }
        keepAliveRouters.value = Array.from(new Set(keepArrTemp))
    }


    const route = useRoute()

    emitter.on('setKeepAlive', setKeepAliveRouters)

    const asyncRouters = ref<RouteItem[]>([])

    const topMenu = ref<RouteItem[]>([])

    const leftMenu = ref<RouteItem[]>([])

    const menuMap: Record<string, RouteItem> = {}

    const topActive = ref('')

    const setLeftMenu = (name: string) => {
        sessionStorage.setItem('topActive', name)
        topActive.value = name
        leftMenu.value = []
        if (menuMap[name]?.children) {
            leftMenu.value = menuMap[name].children!
        }
        return menuMap[name]?.children
    }

    const findTopActive = (menuMapLocal: Record<string, RouteItem>, routeName: string): string | null => {
        for (let topName in menuMapLocal) {
            const topItem = menuMapLocal[topName]
            if (topItem.children?.some(item => item.name === routeName)) {
                return topName
            }
            const foundName = findTopActive(topItem.children as any || {}, routeName)
            if (foundName) {
                return topName
            }
        }
        return null
    }

    watchEffect(() => {
        let topActiveSession = sessionStorage.getItem('topActive')
        topMenu.value = []
        const children = asyncRouters.value[0]?.children || []
        children.forEach((item) => {
            if (item.hidden) return
            menuMap[item.name] = item
            topMenu.value.push({ ...item, children: [] })
        })
        if (!topActiveSession || topActiveSession === 'undefined' || topActiveSession === 'null') {
            topActiveSession = findTopActive(menuMap, route.name as string)
        }
        setLeftMenu(topActiveSession || '')
    })

    const routeMap: Record<string, RouteItem> = {}

    const SetAsyncRouter = async () => {
        asyncRouterFlag.value++
        const baseRouter: RouteItem[] = [
            {
                path: '/layout',
                name: 'layout',
                component: 'view/layout/index.vue',
                meta: {
                    title: '底层layout'
                },
                children: []
            }
        ]
        const asyncRouterRes = await asyncMenu()
        const asyncRouter = asyncRouterRes.data.menus

        asyncRouter.forEach((level1: RouteItem) => {
            if (level1.children) {
                const rolesParentIndex = level1.children.findIndex(c => c.name === 'rolesParent')
                if (rolesParentIndex !== -1) {
                    const rolesParent = level1.children[rolesParentIndex]
                    if (rolesParent.children && rolesParent.children.length > 0) {
                        rolesParent.children.forEach(child => {
                            level1.children!.push(child)
                        })
                        level1.children.splice(rolesParentIndex, 1)
                        level1.children.sort((a, b) => (a.sort || 0) - (b.sort || 0))
                    }
                }
            }
        })
        asyncRouter &&
            asyncRouter.push({
                path: 'reload',
                name: 'Reload',
                hidden: true,
                meta: {
                    title: '',
                    closeTab: true
                },
                component: 'view/error/reload.vue'
            })
        formatRouter(asyncRouter, routeMap)
        baseRouter[0].children = asyncRouter

        // 追加手动路由（详情页、创建页等隐藏路由），与后端菜单路由一起注册
        const manualRoutes: RouteItem[] = [
            {
                path: 'notebooks/detail',
                name: 'notebookDetail',
                hidden: true,
                component: 'view/console/notebooks/detail.vue',
                meta: { title: '容器实例详情', hidden: true, keepAlive: false }
            },
            {
                path: 'notebooks/create',
                name: 'notebookCreate',
                hidden: true,
                component: 'view/console/notebooks/create.vue',
                meta: { title: '创建容器实例', hidden: true, keepAlive: false }
            },
            {
                path: 'training/detail',
                name: 'trainingDetail',
                hidden: true,
                component: 'view/console/training/detail.vue',
                meta: { title: '训练任务详情', hidden: true, keepAlive: false }
            },
            {
                path: 'training/create',
                name: 'trainingCreate',
                hidden: true,
                component: 'view/console/training/create.vue',
                meta: { title: '创建训练任务', hidden: true, keepAlive: false }
            },
            {
                path: 'inference/detail',
                name: 'inferenceDetail',
                hidden: true,
                component: 'view/console/inference/detail.vue',
                meta: { title: '推理服务详情', hidden: true, keepAlive: false }
            },
            {
                path: 'inference/create',
                name: 'inferenceCreate',
                hidden: true,
                component: 'view/console/inference/create.vue',
                meta: { title: '创建推理服务', hidden: true, keepAlive: false }
            }
        ]
        baseRouter[0].children = [...(baseRouter[0].children || []), ...manualRoutes]

        if (notLayoutRouterArr.length !== 0) {
            baseRouter.push(...notLayoutRouterArr)
        }
        asyncRouterHandle(baseRouter)
        KeepAliveFilter(asyncRouter)
        asyncRouters.value = baseRouter

        // Make sure the router knows about these new dynamic routes immediately
        baseRouter.forEach(r => router.addRoute(r as unknown as RouteRecordRaw))

        return true
    }

    return {
        topActive,
        setLeftMenu,
        topMenu,
        leftMenu,
        asyncRouters,
        keepAliveRouters,
        asyncRouterFlag,
        SetAsyncRouter,
        routeMap
    }
})

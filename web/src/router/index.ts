import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'

/**
 * 静态路由表
 * - Landing: 网站介绍首页（公开页面，无需登录）
 * - Login: 登录/注册页面
 * - NotFound: 404错误页面
 * 
 * 注意：控制台相关页面通过动态路由加载（/layout 下的子路由）
 */
const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Landing',
        component: () => import('@/view/landing/index.vue'),
        meta: {
            title: '机器学习平台 - 深度学习平台'
        }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/view/login/index.vue'),
        meta: {
            title: '登录 - 机器学习平台'
        }
    },
    {
        path: '/:catchAll(.*)',
        name: 'NotFound',
        meta: {
            closeTab: true
        },
        component: () => import('@/view/error/index.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router

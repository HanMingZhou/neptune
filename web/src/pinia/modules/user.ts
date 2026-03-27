import { login, getUserInfo } from '@/api/user'
import { mfaLogin } from '@/api/account'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouterStore } from './router'
import { useCookies } from '@vueuse/integrations/useCookies'
import { useStorage } from '@vueuse/core'

import { useAppStore } from '@/pinia'

export interface UserInfo {
    uuid: string
    nickName: string
    headerImg: string
    authority: {
        defaultRouter?: string
        [key: string]: any
    }
    originSetting?: Record<string, any>
    [key: string]: any
}

export interface LoginInfo {
    username: string
    password: string
    captcha?: string
    captchaId?: string
}

export const useUserStore = defineStore('user', () => {
    const appStore = useAppStore()
    const loadingInstance = ref<any>(null)

    const userInfo = ref<UserInfo>({
        uuid: '',
        nickName: '',
        headerImg: '',
        authority: {}
    })
    const token = useStorage('token', '')
    const xToken = useCookies(['x-token'])
    const currentToken = computed(() => token.value || xToken.get('x-token') || '')

    const setUserInfo = (val: UserInfo) => {
        userInfo.value = val
        if (val.originSetting) {
            Object.keys(appStore.config).forEach((key) => {
                if (val.originSetting![key] !== undefined) {
                    (appStore.config as any)[key] = val.originSetting![key]
                }
            })
        }
    }

    const setToken = (val: string) => {
        token.value = val
        xToken.set('x-token', val)
    }

    const NeedInit = async () => {
        await ClearStorage()
        await router.push({ name: 'Init', replace: true })
    }

    const ResetUserInfo = (value: Partial<UserInfo> = {}) => {
        userInfo.value = {
            ...userInfo.value,
            ...value
        }
    }

    const GetUserInfo = async () => {
        const res = await getUserInfo()
        if (res.code === 0) {
            setUserInfo(res.data.userInfo)
        }
        return res
    }

    const LoginIn = async (loginInfo: LoginInfo) => {
        try {
            loadingInstance.value = ElLoading.service({
                fullscreen: true,
                text: '登录中，请稍候...'
            })

            const res = await login(loginInfo)

            if (res.code !== 0) {
                return false
            }

            // 检查是否需要MFA二次验证
            if (res.data.needMfa) {
                return { needMfa: true, mfaToken: res.data.mfaToken }
            }

            setUserInfo(res.data.user)
            setToken(res.data.token)

            const routerStore = useRouterStore()
            await routerStore.SetAsyncRouter()
            const asyncRouters = routerStore.asyncRouters

            asyncRouters.forEach((asyncRouter: any) => {
                router.addRoute(asyncRouter)
            })

            if (router.currentRoute.value.query.redirect) {
                await router.replace(router.currentRoute.value.query.redirect as string)
                return true
            }

            if (router.hasRoute(userInfo.value.authority.defaultRouter!)) {
                await router.replace({ name: userInfo.value.authority.defaultRouter })
            } else {
                const firstMenu = asyncRouters[0]?.children?.[0]
                if (firstMenu?.name) {
                    await router.replace({ name: firstMenu.name })
                } else {
                    console.error('未找到可用的菜单路由，请检查菜单配置')
                    await router.replace({ name: 'Landing' })
                }
            }

            const isWindows = /windows/i.test(navigator.userAgent)
            window.localStorage.setItem('osType', isWindows ? 'WIN' : 'MAC')

            return true
        } catch (error) {
            console.error('LoginIn error:', error)
            return false
        } finally {
            loadingInstance.value?.close()
        }
    }

    // MFA二次验证登录
    const MfaLoginIn = async (mfaToken: string, code: string) => {
        try {
            loadingInstance.value = ElLoading.service({
                fullscreen: true,
                text: 'MFA 验证中...'
            })

            const res = await mfaLogin({ mfaToken, code })

            if (res.code !== 0) {
                return false
            }

            setUserInfo(res.data.user)
            setToken(res.data.token)

            const routerStore = useRouterStore()
            await routerStore.SetAsyncRouter()
            const asyncRouters = routerStore.asyncRouters

            asyncRouters.forEach((asyncRouter: any) => {
                router.addRoute(asyncRouter)
            })

            if (router.currentRoute.value.query.redirect) {
                await router.replace(router.currentRoute.value.query.redirect as string)
                return true
            }

            if (router.hasRoute(userInfo.value.authority.defaultRouter!)) {
                await router.replace({ name: userInfo.value.authority.defaultRouter })
            } else {
                const firstMenu = asyncRouters[0]?.children?.[0]
                if (firstMenu?.name) {
                    await router.replace({ name: firstMenu.name })
                } else {
                    await router.replace({ name: 'Landing' })
                }
            }

            return true
        } catch (error) {
            console.error('MfaLoginIn error:', error)
            return false
        } finally {
            loadingInstance.value?.close()
        }
    }

    const LoginOut = async () => {
        const res = await jsonInBlacklist()

        if (res.code !== 0) {
            return
        }

        await ClearStorage()

        router.push({ name: 'Login', replace: true })
        window.location.reload()
    }

    const ClearStorage = async () => {
        token.value = ''
        xToken.remove('x-token')
        sessionStorage.clear()
        localStorage.removeItem('originSetting')
        localStorage.removeItem('token')
    }

    return {
        userInfo,
        token: currentToken,
        NeedInit,
        ResetUserInfo,
        GetUserInfo,
        LoginIn,
        MfaLoginIn,
        LoginOut,
        setToken,
        loadingInstance,
        ClearStorage
    }
})

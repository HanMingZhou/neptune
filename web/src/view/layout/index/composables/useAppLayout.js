import { computed, ref } from 'vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import en from 'element-plus/dist/locale/en.mjs'
import { getOrderOverview } from '@/api/order'
import { setUserAuthority } from '@/api/user'
import { translations } from '@/locales'
import { useRouterStore } from '@/pinia/modules/router'
import { useUserStore } from '@/pinia/modules/user'

const iconMap = {
  dashboard: 'dashboard',
  notebooks: 'terminal',
  notebookrouter: 'terminal',
  notebooklist: 'terminal',
  training: 'model_training',
  trainingparent: 'model_training',
  traininglist: 'model_training',
  inference: 'rocket_launch',
  storage: 'folder_open',
  sshkeys: 'vpn_key',
  order: 'payments',
  transactions: 'receipt',
  usage: 'data_usage',
  invoice: 'description',
  images: 'photo_library',
  account: 'person',
  security: 'shield',
  accessrecords: 'history',
  settings: 'settings',
  admin: 'admin_panel_settings',
  superadmin: 'admin_panel_settings',
  products: 'inventory_2',
  cmsproduct: 'inventory_2',
  rolesparent: 'admin_panel_settings',
  authority: 'badge',
  roles: 'badge',
  menu: 'menu_open',
  menus: 'menu_open',
  api: 'api',
  apis: 'api',
  user: 'group',
  users: 'group',
  operation: 'receipt_long',
  operations: 'receipt_long',
  person: 'person',
  clustermanage: 'cloud',
  nodemanage: 'hub',
  default: 'circle'
}

const iconAliasMap = {
  'circle-check': 'check_circle',
  'user-filled': 'person',
  user: 'person',
  setting: 'settings',
  delete: 'delete',
  edit: 'edit',
  search: 'search',
  upload: 'upload',
  download: 'download',
  plus: 'add',
  minus: 'remove',
  close: 'close',
  check: 'check',
  'arrow-right': 'arrow_forward',
  'arrow-left': 'arrow_back',
  warning: 'warning',
  'info-filled': 'info',
  'question-filled': 'help',
  'star-filled': 'star',
  document: 'description',
  folder: 'folder',
  picture: 'image',
  'video-camera': 'videocam',
  phone: 'phone',
  message: 'email',
  lock: 'lock',
  unlock: 'lock_open',
  key: 'vpn_key',
  monitor: 'monitor',
  cpu: 'memory',
  cloudy: 'cloud'
}

const computeRoutes = ['dashboard', 'notebooks', 'notebookrouter', 'training', 'trainingparent', 'inference']
const resourceRoutes = ['sshkeys', 'order', 'storage', 'images']
const managementRoutes = ['account', 'person']

const getIcon = (name) => {
  const key = name?.toLowerCase() || ''
  return iconMap[key] || iconMap.default
}

const resolveIcon = (routeName, metaIcon) => {
  if (metaIcon && iconAliasMap[metaIcon]) {
    return iconAliasMap[metaIcon]
  }

  if (metaIcon) {
    return metaIcon
  }

  return getIcon(routeName)
}

const createTranslator = (langRef) => (key, params = {}) => {
  if (!key) {
    return ''
  }

  const langDict = translations[langRef.value] || {}
  let value = langDict[key] || langDict[key.toLowerCase()]

  if (!value && key.includes('.')) {
    const keys = key.split('.')
    let current = langDict

    for (const currentKey of keys) {
      if (current === undefined || current === null) {
        break
      }

      current = current[currentKey]
    }

    if (typeof current === 'string') {
      value = current
    }
  }

  if (!value) {
    value = key
  }

  Object.keys(params).forEach((paramKey) => {
    value = value.replace(`{${paramKey}}`, params[paramKey])
  })

  return value
}

const convertRouteToNavItem = (route) => {
  const visibleChildren = route.children?.filter((child) => !child.hidden && child.meta) || []
  const hasVisibleChildren = visibleChildren.length > 0
  const item = {
    key: route.name,
    titleKey: route.name,
    title: route.meta?.title || route.name,
    icon: resolveIcon(route.name, route.meta?.icon),
    routeName: hasVisibleChildren ? undefined : route.name
  }

  if (hasVisibleChildren) {
    item.children = visibleChildren.map((child) => convertRouteToNavItem(child))
  }

  return item
}

export function useAppLayout() {
  const routerStore = useRouterStore()
  const userStore = useUserStore()

  const lang = ref('zh')
  const isDark = ref(false)
  const userBalance = ref(0)
  const elLocale = computed(() => (lang.value === 'zh' ? zhCn : en))
  const t = createTranslator(lang)
  const userInfo = computed(() => userStore.userInfo || {})

  const dynamicNavigation = computed(() => {
    const asyncRouters = routerStore.asyncRouters[0]?.children || []
    const groups = {
      admin: { title: 'admin', items: [], sort: 4 },
      compute: { title: 'compute', items: [], sort: 1 },
      management: { title: 'management', items: [], sort: 3 },
      resources: { title: 'resources', items: [], sort: 2 }
    }

    asyncRouters.forEach((route) => {
      if (route.hidden || !route.meta) {
        return
      }

      const routeName = route.name?.toLowerCase() || ''
      const item = convertRouteToNavItem(route)

      if (computeRoutes.some((name) => routeName.includes(name))) {
        groups.compute.items.push(item)
      } else if (resourceRoutes.some((name) => routeName.includes(name))) {
        groups.resources.items.push(item)
      } else if (managementRoutes.some((name) => routeName.includes(name))) {
        groups.management.items.push(item)
      } else {
        groups.admin.items.push(item)
      }
    })

    return Object.values(groups)
      .filter((group) => group.items.length > 0)
      .sort((left, right) => left.sort - right.sort)
  })

  const toggleLang = () => {
    lang.value = lang.value === 'en' ? 'zh' : 'en'
  }

  const toggleTheme = () => {
    const nextDark = document.documentElement.classList.toggle('dark')
    isDark.value = nextDark
    sessionStorage.setItem('theme', nextDark ? 'dark' : 'light')
  }

  const changeUserAuth = async (authorityId) => {
    const res = await setUserAuthority({ authorityId })
    if (res.code === 0) {
      window.sessionStorage.setItem('needCloseAll', 'true')
      window.sessionStorage.setItem('needToHome', 'true')
      window.location.reload()
    }
  }

  const logout = () => {
    userStore.LoginOut()
  }

  const fetchBalance = async () => {
    try {
      const res = await getOrderOverview()
      if (res.code === 0) {
        userBalance.value = res.data?.balance ?? 0
      }
    } catch (error) {
      console.error('Failed to fetch balance:', error)
    }
  }

  const initialize = async () => {
    const savedTheme = sessionStorage.getItem('theme')
    if (savedTheme === 'dark') {
      document.documentElement.classList.add('dark')
      isDark.value = true
    } else if (savedTheme === 'light') {
      document.documentElement.classList.remove('dark')
      isDark.value = false
    } else {
      isDark.value = document.documentElement.classList.contains('dark')
    }

    if (userStore.loadingInstance) {
      userStore.loadingInstance.close()
    }

    await fetchBalance()
  }

  return {
    changeUserAuth,
    dynamicNavigation,
    elLocale,
    initialize,
    isDark,
    lang,
    logout,
    routerStore,
    t,
    toggleLang,
    toggleTheme,
    userBalance,
    userInfo
  }
}

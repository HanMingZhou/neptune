import type { RouteItem } from './router'

const MANUAL_ROUTE_DEFINITIONS = [
  {
    keyword: 'notebook',
    routes: [
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
      }
    ]
  },
  {
    keyword: 'training',
    routes: [
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
      }
    ]
  },
  {
    keyword: 'inference',
    routes: [
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
  }
] satisfies Array<{
  keyword: string
  routes: RouteItem[]
}>

const collectRouteCandidates = (
  route: Pick<RouteItem, 'component' | 'meta' | 'name' | 'path'>
): string[] =>
  [
    route.name,
    route.path,
    route.meta?.path,
    typeof route.component === 'string' ? route.component : ''
  ]
    .map((value) => String(value || '').toLowerCase())
    .filter(Boolean)

export const routeMatchesKeyword = (
  route: Pick<RouteItem, 'component' | 'meta' | 'name' | 'path'>,
  keyword: string
): boolean => {
  const normalizedKeyword = keyword.toLowerCase()
  return collectRouteCandidates(route).some((candidate) =>
    candidate.includes(normalizedKeyword)
  )
}

export const hasRouteKeyword = (
  routes: RouteItem[] | undefined,
  keyword: string
): boolean => {
  if (!routes || routes.length === 0) {
    return false
  }

  return routes.some(
    (route) =>
      routeMatchesKeyword(route, keyword) ||
      hasRouteKeyword(route.children, keyword)
  )
}

export const createConsoleManualRoutes = (
  routes: RouteItem[] | undefined
): RouteItem[] =>
  MANUAL_ROUTE_DEFINITIONS.flatMap((definition) =>
    hasRouteKeyword(routes, definition.keyword) ? definition.routes : []
  )

import type { RouteItem } from '@/pinia/modules/router'
import type { LayoutGroupKey } from '@/types/layout'

const ROUTE_GROUP_KEYWORDS: Record<
  Exclude<LayoutGroupKey, 'admin'>,
  string[]
> = {
  compute: ['dashboard', 'notebook', 'training', 'inference'],
  resources: [
    'sshkey',
    'order',
    'storage',
    'image',
    'invoice',
    'transaction',
    'usage'
  ],
  management: ['account', 'person', 'security', 'accesslog', 'accessrecord']
}

const GROUP_ALIAS_MAP: Record<string, LayoutGroupKey> = {
  admin: 'admin',
  compute: 'compute',
  management: 'management',
  resource: 'resources',
  resources: 'resources'
}

const normalizeValue = (value: unknown): string =>
  String(value || '')
    .trim()
    .toLowerCase()

const collectRouteTokens = (
  route: Pick<RouteItem, 'component' | 'meta' | 'name' | 'path'>
): string[] =>
  [
    route.name,
    route.path,
    route.meta?.path,
    typeof route.component === 'string' ? route.component : ''
  ]
    .map(normalizeValue)
    .filter(Boolean)

const resolveExplicitGroup = (value: unknown): LayoutGroupKey | null => {
  const normalizedValue = normalizeValue(value)
  return GROUP_ALIAS_MAP[normalizedValue] || null
}

export const resolveLayoutGroupKey = (
  route: Pick<RouteItem, 'component' | 'meta' | 'name' | 'path'>
): LayoutGroupKey => {
  const explicitGroup =
    resolveExplicitGroup(route.meta?.navGroup) ||
    resolveExplicitGroup(route.meta?.group)
  if (explicitGroup) {
    return explicitGroup
  }

  const routeTokens = collectRouteTokens(route)
  const matchedGroup = (
    Object.entries(ROUTE_GROUP_KEYWORDS) as Array<
      [Exclude<LayoutGroupKey, 'admin'>, string[]]
    >
  ).find(([, keywords]) =>
    routeTokens.some((token) =>
      keywords.some((keyword) => token.includes(keyword))
    )
  )

  return matchedGroup?.[0] || 'admin'
}

import { describe, expect, it } from 'vitest'
import {
  createConsoleManualRoutes,
  hasRouteKeyword,
  routeMatchesKeyword
} from './routerAccess'
import type { RouteItem } from './router'

const createRoute = (overrides: Partial<RouteItem>): RouteItem => ({
  path: 'sample',
  name: 'sample',
  meta: {
    title: 'sample'
  },
  ...overrides
})

describe('router access helpers', () => {
  it('matches keywords from route component paths', () => {
    expect(
      routeMatchesKeyword(
        createRoute({
          component: 'view/console/notebooks/index.vue'
        }),
        'notebook'
      )
    ).toBe(true)
  })

  it('recursively detects module permissions from nested routes', () => {
    const routes: RouteItem[] = [
      createRoute({
        name: 'computeRoot',
        children: [
          createRoute({
            name: 'jobs',
            path: 'training/list'
          })
        ]
      })
    ]

    expect(hasRouteKeyword(routes, 'training')).toBe(true)
    expect(hasRouteKeyword(routes, 'inference')).toBe(false)
  })

  it('only creates hidden manual routes for modules present in backend menus', () => {
    const routes: RouteItem[] = [
      createRoute({
        name: 'notebookrouter',
        path: 'notebooks'
      }),
      createRoute({
        name: 'trainingRoot',
        path: 'training'
      })
    ]

    const manualRoutes = createConsoleManualRoutes(routes)
    expect(manualRoutes.map((route) => route.name)).toEqual([
      'notebookDetail',
      'notebookCreate',
      'trainingDetail',
      'trainingCreate'
    ])
  })
})

import { describe, expect, it } from 'vitest'
import { resolveLayoutGroupKey } from './navigationGrouping'

describe('resolveLayoutGroupKey', () => {
  it('uses explicit backend nav group when provided', () => {
    expect(
      resolveLayoutGroupKey({
        name: 'anything',
        path: '/misc',
        meta: {
          group: 'resources',
          title: 'Misc'
        }
      } as never)
    ).toBe('resources')
  })

  it('classifies compute routes using path tokens, not only route name', () => {
    expect(
      resolveLayoutGroupKey({
        name: 'modelCenter',
        path: '/console/training/jobs',
        meta: {
          title: 'Training'
        }
      } as never)
    ).toBe('compute')
  })

  it('classifies management routes using meta path tokens', () => {
    expect(
      resolveLayoutGroupKey({
        name: 'profileCenter',
        path: '/console/profile',
        meta: {
          path: '/layout/account/security',
          title: 'Security'
        }
      } as never)
    ).toBe('management')
  })

  it('falls back to admin when no known tokens exist', () => {
    expect(
      resolveLayoutGroupKey({
        name: 'customAdminPage',
        path: '/console/custom',
        meta: {
          title: 'Custom'
        }
      } as never)
    ).toBe('admin')
  })
})

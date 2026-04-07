import { afterEach, describe, expect, it } from 'vitest'
import config from '@/core/config'
import { useMenuIconOptions } from './useMenuIconOptions'

const originalLogs = [...config.logs]

afterEach(() => {
  config.logs = [...originalLogs]
})

describe('useMenuIconOptions', () => {
  it('includes element plus icons by default', () => {
    const { iconOptions } = useMenuIconOptions()

    expect(
      iconOptions.value.some((option) => option.key === 'alarm-clock')
    ).toBe(true)
  })

  it('appends custom registered icons without duplicates', () => {
    config.logs = [
      { key: 'server', label: 'Server' },
      { key: 'server', label: 'Server Duplicate' },
      { key: 'alarm-clock', label: 'Alarm Clock Custom' }
    ]

    const { iconOptions } = useMenuIconOptions()
    const serverOptions = iconOptions.value.filter(
      (option) => option.key === 'server'
    )
    const alarmClockOptions = iconOptions.value.filter(
      (option) => option.key === 'alarm-clock'
    )

    expect(serverOptions).toHaveLength(1)
    expect(serverOptions[0]?.label).toBe('Server')
    expect(alarmClockOptions).toHaveLength(1)
    expect(alarmClockOptions[0]?.label).toBe('alarm-clock')
  })
})

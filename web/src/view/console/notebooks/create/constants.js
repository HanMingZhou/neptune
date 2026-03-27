export const DEFAULT_VOLUME_MOUNT_PATH = '/home/notebook/neptune'

export const NOTEBOOK_PAY_TYPES = [
  { value: 1, labelKey: 'payHourly' },
  { value: 2, labelKey: 'payDaily' },
  { value: 3, labelKey: 'payWeekly' },
  { value: 4, labelKey: 'payMonthly' }
]

export const NOTEBOOK_IMAGE_TABS = [
  { value: 'base', labelKey: 'baseImage' },
  { value: 'community', labelKey: 'communityImage' },
  { value: 'my', labelKey: 'myImage' }
]

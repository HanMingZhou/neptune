export const TRAINING_PAY_TYPES = [
  { value: 1, labelKey: 'payHourly' },
  { value: 2, labelKey: 'payDaily' },
  { value: 3, labelKey: 'payWeekly' },
  { value: 4, labelKey: 'payMonthly' }
]

export const TRAINING_IMAGE_TABS = [
  { value: 'base', labelKey: 'baseImage' },
  { value: 'community', labelKey: 'communityImage' },
  { value: 'my', labelKey: 'myImage' }
]

export const TRAINING_FRAMEWORK_TYPES = [
  { value: 'STANDALONE', label: 'Standalone', hintKey: 'standaloneHint' },
  { value: 'PYTORCH_DDP', label: 'PyTorch DDP', hintKey: 'pytorchDdpHint' },
  { value: 'MPI', label: 'MPI', hintKey: 'mpiModeHint' }
]

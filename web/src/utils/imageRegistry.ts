import type { ConsoleImage } from '@/types/consoleResource'

const trimImagePath = (value?: string): string =>
  String(value || '')
    .trim()
    .replace(/^\/+/, '')

export const composeImageAddr = (
  harborAddr?: string,
  imagePath?: string
): string => {
  const registry = String(harborAddr || '')
    .trim()
    .replace(/\/+$/, '')
  const path = trimImagePath(imagePath)

  if (!registry || !path) {
    return ''
  }

  return `${registry}/${path}`
}

export const buildImageDescription = (image: ConsoleImage): string => {
  const parts = [image.clusterName, image.imagePath || image.image]
    .map((item) => String(item || '').trim())
    .filter(Boolean)

  return parts.join(' · ')
}

export const decorateConsoleImages = (
  images: ConsoleImage[] = []
): ConsoleImage[] =>
  images.map((image) => ({
    ...image,
    description: buildImageDescription(image)
  }))

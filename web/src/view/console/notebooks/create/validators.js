export const validateTensorBoardPath = (path, t = (key) => key) => {
  if (!path) return null

  const normalizedPath = path.trim()
  if (normalizedPath.includes('..')) return t('tensorboardPathError')

  const invalidChars = [';', '&', '|', '`', '$', '(', ')', '{', '}', '[', ']', '<', '>', '\\', '"', "'"]
  for (const char of invalidChars) {
    if (normalizedPath.includes(char)) return t('tensorboardPathError')
  }

  if (normalizedPath.startsWith('~')) return t('tensorboardPathError')
  if (normalizedPath.length > 200) return t('tensorboardPathError')

  const validPattern = /^[a-zA-Z0-9\-_/.]+$/
  if (!validPattern.test(normalizedPath)) return t('tensorboardPathError')

  return null
}

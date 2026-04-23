const RFC1123_REGEX =
  /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/

type TranslationParams = Record<string, string | number>
type Translator = (key: string, params?: TranslationParams) => string

interface ResourceNameValidationOptions {
  example?: string
  fieldKey?: string
  fieldLabel?: string
  t?: Translator
}

interface ErrorWithResponse {
  message?: string
  response?: {
    data?: {
      message?: string
      msg?: string
    }
  }
}

export const RESOURCE_NAME_MAX_LENGTH = 63

export const validateK8sResourceName = (
  value: unknown,
  {
    t = (key) => key,
    fieldKey = 'name',
    fieldLabel,
    example = 'my-resource'
  }: ResourceNameValidationOptions = {}
): string | null => {
  const label = fieldLabel || t(fieldKey)
  const normalizedValue = String(value || '').trim()

  if (!normalizedValue) {
    return t('resourceNameRequired', { field: label })
  }

  if (normalizedValue.length > RESOURCE_NAME_MAX_LENGTH) {
    return t('resourceNameTooLong', {
      field: label,
      max: RESOURCE_NAME_MAX_LENGTH
    })
  }

  if (!RFC1123_REGEX.test(normalizedValue)) {
    return t('resourceNameRule', { field: label, example })
  }

  return null
}

export const isResourceNameErrorMessage = (message = ''): boolean =>
  /((名称|name).*(不能为空|不能超过|小写字母|数字|lowercase|start|end|exceed))|(my-notebook|my-training|my-service)/i.test(
    message
  )

export const getErrorMessage = (
  error: unknown,
  fallbackMessage: string
): string => {
  const normalizedError = error as ErrorWithResponse | undefined

  return (
    normalizedError?.response?.data?.msg ||
    normalizedError?.response?.data?.message ||
    normalizedError?.message ||
    fallbackMessage
  )
}

export const getSubmitErrorMessage = (
  error: unknown,
  fallbackMessage: string
): string => getErrorMessage(error, fallbackMessage)

export const validateTensorBoardPath = (
  path: string,
  t: Translator = (key) => key
): string | null => {
  if (!path) {
    return null
  }

  const normalizedPath = path.trim()
  if (normalizedPath.includes('..')) {
    return t('tensorboardPathError')
  }

  const invalidChars = [
    ';',
    '&',
    '|',
    '`',
    '$',
    '(',
    ')',
    '{',
    '}',
    '[',
    ']',
    '<',
    '>',
    '\\',
    '"',
    "'"
  ]
  for (const char of invalidChars) {
    if (normalizedPath.includes(char)) {
      return t('tensorboardPathError')
    }
  }

  if (normalizedPath.startsWith('~')) {
    return t('tensorboardPathError')
  }
  if (normalizedPath.length > 200) {
    return t('tensorboardPathError')
  }

  const validPattern = /^[a-zA-Z0-9\-_/.]+$/
  if (!validPattern.test(normalizedPath)) {
    return t('tensorboardPathError')
  }

  return null
}

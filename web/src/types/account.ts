import type { ResourceId } from './consoleResource'

export interface AccessLogSession {
  id: ResourceId
  current?: boolean
  device?: string
  location?: string
  ip?: string
  time?: string
  [key: string]: unknown
}

export interface AccessLogRecord {
  id: ResourceId
  time?: string
  ip?: string
  location?: string
  device?: string
  status?: string
  [key: string]: unknown
}

export interface SecurityStatusData {
  phone?: string
  phoneStatus?: string
  email?: string
  emailStatus?: string
  mfaEnabled?: boolean
  mfaStatus?: string
  githubBound?: boolean
  githubUsername?: string
  githubStatus?: string
  accessKeyId?: string
  accessKeyStatus?: string
  securityScore?: number
  [key: string]: unknown
}

export interface AccountPasswordForm {
  password: string
  newPassword: string
  confirmPassword: string
}

export interface AccountPhoneForm {
  phone: string
  code: string
}

export interface AccountEmailForm {
  email: string
  code: string
}

export interface MfaSetupData {
  qr?: string
  secret?: string
}

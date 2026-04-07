import type { ConsoleSshKey } from './consoleResource'

export interface SshKeyListItem extends ConsoleSshKey {
  fingerprint?: string
  createdAt?: string | number
}

export interface SshKeyCreateForm {
  name: string
  publicKey: string
}

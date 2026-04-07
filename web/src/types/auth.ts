export interface LoginFormState {
  username: string
  password: string
  email: string
  captcha: string
}

export interface LoginFieldErrors {
  username: string
  password: string
  email: string
  captcha: string
  mfaCode: string
}

export interface CaptchaData {
  picPath?: string
  captchaId?: string
}

export interface LoginMfaResult {
  needMfa: true
  mfaToken: string
}

export interface RegisterPayload {
  userName: string
  passWord: string
  email: string
}

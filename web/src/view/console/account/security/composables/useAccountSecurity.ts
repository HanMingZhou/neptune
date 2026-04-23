import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormRules } from 'element-plus'
import {
  bindAccount,
  generateAccessKey,
  getSecurityStatus,
  mfaActivate,
  mfaSetup,
  toggleMfa,
  updatePassword
} from '@/api/account'
import type { ApiResponse } from '@/utils/request'
import { getErrorMessage } from '@/utils/resourceValidators'
import type { Translator } from '@/types/consoleResource'
import type {
  AccountEmailForm,
  AccountPasswordForm,
  AccountPhoneForm,
  MfaSetupData,
  SecurityStatusData
} from '@/types/account'

type ValidateCallback = (error?: Error) => void

interface UseAccountSecurityOptions {
  t?: Translator
}

export function useAccountSecurity({ t }: UseAccountSecurityOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const createDefaultAccountInfo = (): Required<
    Omit<SecurityStatusData, 'githubUsername'>
  > & {
    githubUsername: string
  } => ({
    phone: '',
    phoneStatus: translate('security.loading'),
    email: '',
    emailStatus: translate('security.loading'),
    mfaEnabled: false,
    mfaStatus: translate('security.loading'),
    githubBound: false,
    githubUsername: '',
    githubStatus: translate('security.loading'),
    accessKeyId: '',
    accessKeyStatus: translate('security.loading'),
    securityScore: 60
  })

  const accountInfo = reactive(createDefaultAccountInfo())

  const showPassword = ref(false)
  const changePhoneFlag = ref(false)
  const changeEmailFlag = ref(false)
  const showMfaSetup = ref(false)
  const showMfaDisable = ref(false)

  const mfaLoading = ref(false)
  const pwdLoading = ref(false)
  const phoneLoading = ref(false)
  const emailLoading = ref(false)
  const akLoading = ref(false)

  const pwdModify = reactive<AccountPasswordForm>({
    password: '',
    newPassword: '',
    confirmPassword: ''
  })

  const phoneForm = reactive<AccountPhoneForm>({
    phone: '',
    code: ''
  })

  const emailForm = reactive<AccountEmailForm>({
    email: '',
    code: ''
  })

  const time = ref(0)
  const emailTime = ref(0)

  const mfaQr = ref('')
  const mfaSecret = ref('')
  const mfaCode = ref('')
  const mfaDisableCode = ref('')

  let phoneTimer: ReturnType<typeof setInterval> | null = null
  let emailTimer: ReturnType<typeof setInterval> | null = null

  const rules = reactive<FormRules<AccountPasswordForm>>({
    password: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      {
        min: 6,
        message: translate('minCharacters', { min: 6 }),
        trigger: 'blur'
      }
    ],
    newPassword: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      {
        min: 6,
        message: translate('minCharacters', { min: 6 }),
        trigger: 'blur'
      }
    ],
    confirmPassword: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      {
        min: 6,
        message: translate('minCharacters', { min: 6 }),
        trigger: 'blur'
      },
      {
        validator: (
          _rule: unknown,
          value: string,
          callback: ValidateCallback
        ) => {
          if (value !== pwdModify.newPassword) {
            callback(new Error(translate('passwordMismatch')))
            return
          }
          callback()
        },
        trigger: 'blur'
      }
    ]
  })

  const phoneRules = reactive<FormRules<AccountPhoneForm>>({
    phone: [
      { required: true, message: translate('fillAllFields'), trigger: 'blur' },
      {
        pattern: /^1[3-9]\d{9}$/,
        message: translate('illegalPhone'),
        trigger: 'blur'
      }
    ],
    code: [
      { required: true, message: translate('fillAllFields'), trigger: 'blur' }
    ]
  })

  const emailRules = reactive<FormRules<AccountEmailForm>>({
    email: [
      { required: true, message: translate('fillAllFields'), trigger: 'blur' },
      { type: 'email', message: translate('illegalEmail'), trigger: 'blur' }
    ],
    code: [
      { required: true, message: translate('fillAllFields'), trigger: 'blur' }
    ]
  })

  const isPhoneValid = computed(() => /^1[3-9]\d{9}$/.test(phoneForm.phone))
  const isEmailValid = computed(() =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailForm.email)
  )

  const applyAccountInfo = (data: SecurityStatusData = {}): void => {
    Object.assign(accountInfo, createDefaultAccountInfo(), data)
  }

  const loadStatus = async (): Promise<void> => {
    const res = (await getSecurityStatus()) as ApiResponse<SecurityStatusData>
    if (res.code === 0) {
      applyAccountInfo(res.data || {})
    }
  }

  const maskString = (value?: string): string | undefined => {
    if (!value || value.length < 7) {
      return value
    }

    return `${value.substring(0, 3)}****${value.substring(value.length - 4)}`
  }

  const getStatusColor = (status?: string): string => {
    if (status && (status.includes('已') || status === 'Success')) {
      return 'bg-emerald-500/10 text-emerald-500'
    }

    return 'bg-amber-500/10 text-amber-500'
  }

  const clearPassword = (): void => {
    Object.assign(pwdModify, {
      password: '',
      newPassword: '',
      confirmPassword: ''
    })
  }

  const savePassword = async (): Promise<void> => {
    pwdLoading.value = true

    try {
      const res = await updatePassword({
        oldPassword: pwdModify.password,
        newPassword: pwdModify.newPassword
      })

      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        showPassword.value = false
        clearPassword()
        return
      }

      ElMessage.error(getErrorMessage(res, translate('changeFailed')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('changeFailed')))
    } finally {
      pwdLoading.value = false
    }
  }

  const resetTimer = (key: 'phone' | 'email'): void => {
    if (key === 'phone' && phoneTimer) {
      clearInterval(phoneTimer)
      phoneTimer = null
    }

    if (key === 'email' && emailTimer) {
      clearInterval(emailTimer)
      emailTimer = null
    }
  }

  const startCountdown = (
    key: 'phone' | 'email',
    counter: typeof time
  ): void => {
    resetTimer(key)
    counter.value = 60

    const timer = setInterval(() => {
      counter.value -= 1
      if (counter.value <= 0) {
        resetTimer(key)
      }
    }, 1000)

    if (key === 'phone') {
      phoneTimer = timer
    } else {
      emailTimer = timer
    }
  }

  const getCode = (): void => {
    if (!isPhoneValid.value) {
      return
    }

    startCountdown('phone', time)
  }

  const closeChangePhone = (): void => {
    Object.assign(phoneForm, {
      phone: '',
      code: ''
    })
  }

  const changePhone = async (): Promise<void> => {
    phoneLoading.value = true

    try {
      const res = await bindAccount({
        type: 1,
        value: phoneForm.phone,
        code: phoneForm.code
      })
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        changePhoneFlag.value = false
        closeChangePhone()
        await loadStatus()
        return
      }

      ElMessage.error(getErrorMessage(res, translate('changeFailed')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('changeFailed')))
    } finally {
      phoneLoading.value = false
    }
  }

  const getEmailCode = (): void => {
    if (!isEmailValid.value) {
      return
    }

    startCountdown('email', emailTime)
  }

  const closeChangeEmail = (): void => {
    Object.assign(emailForm, {
      email: '',
      code: ''
    })
  }

  const changeEmail = async (): Promise<void> => {
    emailLoading.value = true

    try {
      const res = await bindAccount({
        type: 2,
        value: emailForm.email,
        code: emailForm.code
      })
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        changeEmailFlag.value = false
        closeChangeEmail()
        await loadStatus()
        return
      }

      ElMessage.error(getErrorMessage(res, translate('changeFailed')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('changeFailed')))
    } finally {
      emailLoading.value = false
    }
  }

  const closeMfaSetup = (): void => {
    mfaQr.value = ''
    mfaSecret.value = ''
    mfaCode.value = ''
  }

  const closeMfaDisable = (): void => {
    mfaDisableCode.value = ''
  }

  const enableMfa = async (): Promise<void> => {
    if (mfaLoading.value) {
      return
    }

    mfaLoading.value = true
    try {
      const res = (await mfaSetup()) as ApiResponse<MfaSetupData>
      if (res.code === 0) {
        mfaQr.value = res.data?.qr || ''
        mfaSecret.value = res.data?.secret || ''
        showMfaSetup.value = true
        return
      }

      ElMessage.error(getErrorMessage(res, translate('error')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('error')))
    } finally {
      mfaLoading.value = false
    }
  }

  const handleMfaAction = (): void => {
    if (accountInfo.mfaEnabled) {
      showMfaDisable.value = true
      return
    }

    void enableMfa()
  }

  const confirmMfaActivate = async (): Promise<void> => {
    if (mfaCode.value.length !== 6) {
      return
    }

    mfaLoading.value = true
    try {
      const res = await mfaActivate({ code: mfaCode.value })
      if (res.code === 0) {
        ElMessage.success('MFA 已成功开启')
        showMfaSetup.value = false
        closeMfaSetup()
        await loadStatus()
        return
      }

      ElMessage.error(getErrorMessage(res, '验证码错误'))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, '激活失败'))
    } finally {
      mfaLoading.value = false
    }
  }

  const confirmMfaDisable = async (): Promise<void> => {
    if (mfaDisableCode.value.length !== 6) {
      return
    }

    mfaLoading.value = true
    try {
      const res = await toggleMfa({
        enabled: false,
        code: mfaDisableCode.value
      })
      if (res.code === 0) {
        ElMessage.success('MFA 已关闭')
        showMfaDisable.value = false
        closeMfaDisable()
        await loadStatus()
        return
      }

      ElMessage.error(getErrorMessage(res, '验证码错误'))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, '关闭失败'))
    } finally {
      mfaLoading.value = false
    }
  }

  const writeClipboard = async (
    value: string | undefined,
    successMessage: string
  ): Promise<void> => {
    if (!value) {
      return
    }

    try {
      await navigator.clipboard.writeText(value)
      ElMessage.success(successMessage)
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('error')))
    }
  }

  const copySecret = (): Promise<void> =>
    writeClipboard(mfaSecret.value, '密钥已复制')
  const copyAccessKeyId = (): Promise<void> =>
    writeClipboard(
      typeof accountInfo.accessKeyId === 'string'
        ? accountInfo.accessKeyId
        : undefined,
      'Access Key ID 已复制'
    )

  const handleGenerateAK = async (): Promise<void> => {
    if (akLoading.value) {
      return
    }

    try {
      await ElMessageBox.confirm(
        translate('security.akDesc'),
        translate('security.manageAk'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )
    } catch (error) {
      return
    }

    akLoading.value = true
    try {
      const res = (await generateAccessKey()) as ApiResponse<unknown>
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        await loadStatus()
        return
      }

      ElMessage.error(getErrorMessage(res, translate('error')))
    } catch (error) {
      ElMessage.error(getErrorMessage(error, translate('error')))
    } finally {
      akLoading.value = false
    }
  }

  const setupNotification = (): void => {
    ElMessage.info(translate('order.comingSoon'))
  }

  const bindGithub = (): void => {
    ElMessage.info(translate('order.comingSoon'))
  }

  const passwordDialog = reactive({
    clear: clearPassword,
    form: pwdModify,
    loading: pwdLoading,
    rules,
    submit: savePassword,
    title: computed(() => translate('security.changePwd')),
    visible: showPassword
  })

  const phoneDialog = reactive({
    close: closeChangePhone,
    disableRequestCode: computed(() => time.value > 0 || !isPhoneValid.value),
    form: phoneForm,
    loading: phoneLoading,
    requestCode: getCode,
    requestCodeText: computed(() =>
      time.value > 0 ? `(${time.value}s)后重新获取` : '获取验证码'
    ),
    rules: phoneRules,
    submit: changePhone,
    title: computed(() => translate('security.changePhone')),
    visible: changePhoneFlag
  })

  const emailDialog = reactive({
    close: closeChangeEmail,
    disableRequestCode: computed(
      () => emailTime.value > 0 || !isEmailValid.value
    ),
    form: emailForm,
    loading: emailLoading,
    requestCode: getEmailCode,
    requestCodeText: computed(() =>
      emailTime.value > 0 ? `(${emailTime.value}s)后重新获取` : '获取验证码'
    ),
    rules: emailRules,
    submit: changeEmail,
    title: computed(() => translate('security.emailBind')),
    visible: changeEmailFlag
  })

  const mfaSetupDialog = reactive({
    close: closeMfaSetup,
    code: mfaCode,
    confirm: confirmMfaActivate,
    copySecret,
    loading: mfaLoading,
    qr: mfaQr,
    secret: mfaSecret,
    visible: showMfaSetup
  })

  const mfaDisableDialog = reactive({
    close: closeMfaDisable,
    code: mfaDisableCode,
    confirm: confirmMfaDisable,
    loading: mfaLoading,
    visible: showMfaDisable
  })

  onMounted(() => {
    void loadStatus()
  })

  onBeforeUnmount(() => {
    resetTimer('phone')
    resetTimer('email')
  })

  return {
    accountInfo,
    akLoading,
    bindGithub,
    changeEmail,
    changeEmailFlag,
    changePhone,
    changePhoneFlag,
    clearPassword,
    closeChangeEmail,
    closeChangePhone,
    closeMfaDisable,
    closeMfaSetup,
    confirmMfaActivate,
    confirmMfaDisable,
    copyAccessKeyId,
    copySecret,
    emailDialog,
    emailForm,
    emailLoading,
    emailRules,
    emailTime,
    getCode,
    getEmailCode,
    getStatusColor,
    handleGenerateAK,
    handleMfaAction,
    isEmailValid,
    isPhoneValid,
    maskString,
    mfaCode,
    mfaDisableDialog,
    mfaDisableCode,
    mfaLoading,
    mfaQr,
    mfaSecret,
    mfaSetupDialog,
    passwordDialog,
    phoneDialog,
    phoneForm,
    phoneLoading,
    phoneRules,
    pwdLoading,
    pwdModify,
    rules,
    savePassword,
    setupNotification,
    showMfaDisable,
    showMfaSetup,
    showPassword,
    time
  }
}

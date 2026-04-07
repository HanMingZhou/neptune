import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { ElMessage, type FormRules } from 'element-plus'
import { changePassword, setSelfInfo } from '@/api/user'
import { useUserStore } from '@/pinia/modules/user'
import type { Translator } from '@/types/consoleResource'
import type {
  AccountEmailForm,
  AccountPasswordForm,
  AccountPhoneForm
} from '@/types/account'

type ValidateCallback = (error?: Error) => void

interface UseAccountProfileOptions {
  t?: Translator
}

export function useAccountProfile({ t }: UseAccountProfileOptions = {}) {
  const translate: Translator = t || ((key: string) => key)
  const userStore = useUserStore()

  const editFlag = ref(false)
  const nickName = ref('')
  const showPassword = ref(false)
  const changePhoneFlag = ref(false)
  const changeEmailFlag = ref(false)

  const passwordLoading = ref(false)
  const phoneLoading = ref(false)
  const emailLoading = ref(false)

  const time = ref(0)
  const emailTime = ref(0)

  let phoneTimer: ReturnType<typeof setInterval> | null = null
  let emailTimer: ReturnType<typeof setInterval> | null = null

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

  const phoneRules = reactive<FormRules<AccountPhoneForm>>({})
  const emailRules = reactive<FormRules<AccountEmailForm>>({})

  const clearPassword = (): void => {
    Object.assign(pwdModify, {
      password: '',
      newPassword: '',
      confirmPassword: ''
    })
  }

  const savePassword = async (): Promise<void> => {
    passwordLoading.value = true

    try {
      const res = await changePassword({
        password: pwdModify.password,
        newPassword: pwdModify.newPassword
      })

      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        showPassword.value = false
        clearPassword()
        return
      }

      ElMessage.error(res.msg || translate('changeFailed'))
    } catch (_error) {
      ElMessage.error(translate('changeFailed'))
    } finally {
      passwordLoading.value = false
    }
  }

  const openEdit = (): void => {
    nickName.value = userStore.userInfo.nickName
    editFlag.value = true
  }

  const closeEdit = (): void => {
    nickName.value = ''
    editFlag.value = false
  }

  const enterEdit = async (): Promise<void> => {
    const res = await setSelfInfo({
      nickName: nickName.value
    })

    if (res.code === 0) {
      userStore.ResetUserInfo({ nickName: nickName.value })
      ElMessage.success(translate('changeSuccess'))
    }

    closeEdit()
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
      const res = await setSelfInfo({ phone: phoneForm.phone })
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        userStore.ResetUserInfo({ phone: phoneForm.phone })
        changePhoneFlag.value = false
        closeChangePhone()
        return
      }

      ElMessage.error(res.msg || translate('changeFailed'))
    } catch (_error) {
      ElMessage.error(translate('changeFailed'))
    } finally {
      phoneLoading.value = false
    }
  }

  const getEmailCode = (): void => {
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
      const res = await setSelfInfo({ email: emailForm.email })
      if (res.code === 0) {
        ElMessage.success(translate('changeSuccess'))
        userStore.ResetUserInfo({ email: emailForm.email })
        changeEmailFlag.value = false
        closeChangeEmail()
        return
      }

      ElMessage.error(res.msg || translate('changeFailed'))
    } catch (_error) {
      ElMessage.error(translate('changeFailed'))
    } finally {
      emailLoading.value = false
    }
  }

  const passwordDialog = reactive({
    clear: clearPassword,
    form: pwdModify,
    iconClass: 'text-purple-500',
    loading: passwordLoading,
    rules,
    submit: savePassword,
    title: computed(() => translate('changePasswordTitle')),
    visible: showPassword
  })

  const phoneDialog = reactive({
    close: closeChangePhone,
    disableRequestCode: computed(() => time.value > 0),
    form: phoneForm,
    loading: phoneLoading,
    requestCode: getCode,
    requestCodeText: computed(() =>
      time.value > 0 ? `${time.value}s` : '获取验证码'
    ),
    rules: phoneRules,
    submit: changePhone,
    title: '修改手机号',
    visible: changePhoneFlag
  })

  const emailDialog = reactive({
    close: closeChangeEmail,
    disableRequestCode: computed(() => emailTime.value > 0),
    form: emailForm,
    loading: emailLoading,
    requestCode: getEmailCode,
    requestCodeText: computed(() =>
      emailTime.value > 0 ? `${emailTime.value}s` : '获取验证码'
    ),
    rules: emailRules,
    submit: changeEmail,
    title: '修改邮箱',
    visible: changeEmailFlag
  })

  watch(
    () => userStore.userInfo.headerImg,
    async (value) => {
      const res = await setSelfInfo({ headerImg: value })
      if (res.code === 0) {
        userStore.ResetUserInfo({ headerImg: value })
        ElMessage({
          type: 'success',
          message: translate('changeSuccess')
        })
      }
    }
  )

  onBeforeUnmount(() => {
    resetTimer('phone')
    resetTimer('email')
  })

  return {
    closeEdit,
    editFlag,
    emailDialog,
    enterEdit,
    nickName,
    openEdit,
    passwordDialog,
    phoneDialog,
    userStore
  }
}

import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { changePassword, setSelfInfo } from '@/api/user'
import { useUserStore } from '@/pinia/modules/user'

export function useAccountProfile({ t }) {
  const translate = t || ((key) => key)
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

  let phoneTimer = null
  let emailTimer = null

  const pwdModify = reactive({
    password: '',
    newPassword: '',
    confirmPassword: ''
  })

  const phoneForm = reactive({
    phone: '',
    code: ''
  })

  const emailForm = reactive({
    email: '',
    code: ''
  })

  const rules = reactive({
    password: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      { min: 6, message: translate('minCharacters', { min: 6 }), trigger: 'blur' }
    ],
    newPassword: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      { min: 6, message: translate('minCharacters', { min: 6 }), trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, message: translate('inputPassword'), trigger: 'blur' },
      { min: 6, message: translate('minCharacters', { min: 6 }), trigger: 'blur' },
      {
        validator: (rule, value, callback) => {
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

  const phoneRules = reactive({})
  const emailRules = reactive({})

  const clearPassword = () => {
    Object.assign(pwdModify, {
      password: '',
      newPassword: '',
      confirmPassword: ''
    })
  }

  const savePassword = async () => {
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
    } catch (error) {
      ElMessage.error(translate('changeFailed'))
    } finally {
      passwordLoading.value = false
    }
  }

  const openEdit = () => {
    nickName.value = userStore.userInfo.nickName
    editFlag.value = true
  }

  const closeEdit = () => {
    nickName.value = ''
    editFlag.value = false
  }

  const enterEdit = async () => {
    const res = await setSelfInfo({
      nickName: nickName.value
    })

    if (res.code === 0) {
      userStore.ResetUserInfo({ nickName: nickName.value })
      ElMessage.success(translate('changeSuccess'))
    }

    closeEdit()
  }

  const resetTimer = (key) => {
    if (key === 'phone' && phoneTimer) {
      clearInterval(phoneTimer)
      phoneTimer = null
    }

    if (key === 'email' && emailTimer) {
      clearInterval(emailTimer)
      emailTimer = null
    }
  }

  const startCountdown = (key, counter) => {
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

  const getCode = () => {
    startCountdown('phone', time)
  }

  const closeChangePhone = () => {
    Object.assign(phoneForm, {
      phone: '',
      code: ''
    })
  }

  const changePhone = async () => {
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
    } catch (error) {
      ElMessage.error(translate('changeFailed'))
    } finally {
      phoneLoading.value = false
    }
  }

  const getEmailCode = () => {
    startCountdown('email', emailTime)
  }

  const closeChangeEmail = () => {
    Object.assign(emailForm, {
      email: '',
      code: ''
    })
  }

  const changeEmail = async () => {
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
    } catch (error) {
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
    requestCodeText: computed(() => (time.value > 0 ? `${time.value}s` : '获取验证码')),
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
    requestCodeText: computed(() => (emailTime.value > 0 ? `${emailTime.value}s` : '获取验证码')),
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

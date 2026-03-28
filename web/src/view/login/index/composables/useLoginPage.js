import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { captcha, register } from '@/api/user'
import { useUserStore } from '@/pinia/modules/user'

const createDefaultForm = () => ({
  username: '',
  password: '',
  email: '',
  captcha: ''
})

const createDefaultFieldErrors = () => ({
  username: '',
  password: '',
  email: '',
  captcha: '',
  mfaCode: ''
})

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const MFA_CODE_REGEX = /^\d{6}$/

const getErrorMessage = (error, fallbackMessage) =>
  error?.response?.data?.msg || error?.message || fallbackMessage

export function useLoginPage() {
  const userStore = useUserStore()

  const isRegister = ref(false)
  const showPassword = ref(false)
  const loading = ref(false)
  const errorMsg = ref('')
  const captchaImg = ref('')
  const captchaId = ref('')
  const showMfaStep = ref(false)
  const mfaToken = ref('')
  const mfaCode = ref('')
  const form = reactive(createDefaultForm())
  const fieldErrors = reactive(createDefaultFieldErrors())

  const resetForm = () => {
    Object.assign(form, createDefaultForm())
  }

  const resetFieldErrors = () => {
    Object.assign(fieldErrors, createDefaultFieldErrors())
  }

  const resetMfaState = () => {
    showMfaStep.value = false
    mfaToken.value = ''
    mfaCode.value = ''
    fieldErrors.mfaCode = ''
  }

  const validateForm = () => {
    resetFieldErrors()

    if (showMfaStep.value) {
      const normalizedMfaCode = mfaCode.value.trim()
      mfaCode.value = normalizedMfaCode

      if (!MFA_CODE_REGEX.test(normalizedMfaCode)) {
        fieldErrors.mfaCode = '请输入 6 位数字验证码'
        return false
      }

      return true
    }

    form.username = form.username.trim()
    form.email = form.email.trim()
    form.captcha = form.captcha.trim()

    let valid = true

    if (!form.username) {
      fieldErrors.username = '请输入用户名'
      valid = false
    }

    if (!form.password) {
      fieldErrors.password = '请输入密码'
      valid = false
    }

    if (isRegister.value) {
      if (!form.email) {
        fieldErrors.email = '请输入邮箱地址'
        valid = false
      } else if (!EMAIL_REGEX.test(form.email)) {
        fieldErrors.email = '邮箱格式不正确'
        valid = false
      }
    } else if (!form.captcha) {
      fieldErrors.captcha = '请输入验证码'
      valid = false
    }

    return valid
  }

  const refreshCaptcha = async () => {
    try {
      const res = await captcha()
      if (res.code === 0) {
        captchaImg.value = res.data.picPath
        captchaId.value = res.data.captchaId
      }
    } catch (error) {
      console.error('Failed to load captcha:', error)
    }
  }

  const initialize = async () => {
    await refreshCaptcha()
  }

  const toggleMode = async () => {
    isRegister.value = !isRegister.value
    showPassword.value = false
    errorMsg.value = ''
    resetForm()
    resetFieldErrors()
    resetMfaState()

    if (!isRegister.value) {
      await refreshCaptcha()
    }
  }

  const handleMfaSubmit = async () => {
    const success = await userStore.MfaLoginIn(mfaToken.value, mfaCode.value.trim())
    if (!success) {
      errorMsg.value = 'MFA 验证码错误，请重试'
      fieldErrors.mfaCode = '验证码错误，请重新输入'
    }
  }

  const handleRegisterSubmit = async () => {
    const res = await register({
      userName: form.username.trim(),
      passWord: form.password,
      email: form.email.trim()
    })

    if (res.code === 0) {
      ElMessage.success('注册成功，请登录')
      await toggleMode()
      return
    }

    errorMsg.value = res.msg || '注册失败'
  }

  const handleLoginSubmit = async () => {
    const result = await userStore.LoginIn({
      username: form.username,
      password: form.password,
      captcha: form.captcha.trim(),
      captchaId: captchaId.value
    })

    if (result && typeof result === 'object' && result.needMfa) {
      showMfaStep.value = true
      mfaToken.value = result.mfaToken
      return
    }

    if (!result) {
      errorMsg.value = '登录失败，请检查用户名和密码'
      await refreshCaptcha()
    }
  }

  const handleSubmit = async () => {
    errorMsg.value = ''
    resetFieldErrors()

    if (!validateForm()) {
      return
    }

    loading.value = true

    try {
      if (showMfaStep.value) {
        await handleMfaSubmit()
        return
      }

      if (isRegister.value) {
        await handleRegisterSubmit()
        return
      }

      await handleLoginSubmit()
    } catch (error) {
      errorMsg.value = getErrorMessage(error, '请求失败，请重试')

      if (!isRegister.value && !showMfaStep.value) {
        await refreshCaptcha()
      }
    } finally {
      loading.value = false
    }
  }

  return {
    captchaImg,
    errorMsg,
    fieldErrors,
    form,
    handleSubmit,
    initialize,
    isRegister,
    loading,
    mfaCode,
    refreshCaptcha,
    showMfaStep,
    showPassword,
    toggleMode
  }
}

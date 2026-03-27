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

  const resetForm = () => {
    Object.assign(form, createDefaultForm())
  }

  const resetMfaState = () => {
    showMfaStep.value = false
    mfaToken.value = ''
    mfaCode.value = ''
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
    resetMfaState()

    if (!isRegister.value) {
      await refreshCaptcha()
    }
  }

  const handleMfaSubmit = async () => {
    const success = await userStore.MfaLoginIn(mfaToken.value, mfaCode.value)
    if (!success) {
      errorMsg.value = 'MFA 验证码错误，请重试'
    }
  }

  const handleRegisterSubmit = async () => {
    const res = await register({
      userName: form.username,
      passWord: form.password,
      email: form.email
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
      captcha: form.captcha,
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
    loading.value = true
    errorMsg.value = ''

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

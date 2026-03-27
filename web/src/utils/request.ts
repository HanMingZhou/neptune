import axios, { type AxiosInstance, type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage, type MessageHandler } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 99999
})

let activeAxios: number = 0
let timer: ReturnType<typeof setTimeout> | undefined

const showLoading = () => {
  activeAxios++
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    if (activeAxios > 0) {
      // emitter.emit('showLoading')
    }
  }, 300)
}

const closeLoading = () => {
  activeAxios--
  if (activeAxios <= 0) {
    if (timer) {
      clearTimeout(timer)
    }
    // emitter.emit('closeLoading')
  }
}

// http request 拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // @ts-ignore - donNotShowLoading is a custom property
    if (!config.donNotShowLoading) {
      showLoading()
    }
    const userStore = useUserStore()
    config.headers = config.headers || {}
    config.headers['Content-Type'] = 'application/json'
    config.headers['x-token'] = userStore.token
    config.headers['x-user-id'] = userStore.userInfo?.ID || ''

    return config
  },
  (error: AxiosError) => {
    closeLoading()
    ElMessage({
      showClose: true,
      message: error.message,
      type: 'error'
    })
    return Promise.reject(error)
  }
)

// http response 拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    closeLoading()
    if (response.headers['new-token']) {
      const userStore = useUserStore()
      userStore.setToken(response.headers['new-token'])
    }
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response.data
    }
    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else {
      ElMessage({
        showClose: true,
        message: response.data.msg || (response.headers.msg ? decodeURI(response.headers.msg) : '未知错误'),
        type: 'error'
      })
      if (response.data.data && response.data.data.reload) {
        const userStore = useUserStore()
        userStore.ClearStorage()
      }
      // Return the error data so the caller can handle it if needed
      return response.data.code ? response.data : response
    }
  },
  (error: AxiosError | any) => {
    closeLoading()

    if (!error.response) {
      ElMessage({
        showClose: true,
        message: error.message === 'Network Error' ? '网络连接失败' : error.message,
        type: 'error'
      })
      return Promise.reject(error)
    }

    switch (error.response.status) {
      case 500:
        ElMessage({
          showClose: true,
          message: error.response.data.msg || '服务器错误',
          type: 'error'
        })
        break
      case 404:
        ElMessage({
          showClose: true,
          message: '请求路径找不到',
          type: 'error'
        })
        break
    }

    return Promise.reject(error)
  }
)
// 自定义API响应类型（拦截器已将response.data直接返回）
export interface ApiResponse<T = any> {
  code: number
  data: T
  msg: string
}

export default service as unknown as {
  (config: any): Promise<ApiResponse>
}

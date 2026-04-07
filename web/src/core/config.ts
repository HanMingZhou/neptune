export interface AppConfig {
  appName: string
  keepAliveTabs: boolean
  logs: Array<{ key: string; label: string }>
}

export const config: AppConfig = {
  appName: '机器学习平台',
  keepAliveTabs: false,
  logs: []
}

export default config

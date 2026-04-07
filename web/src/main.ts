import 'element-plus/theme-chalk/dark/css-vars.css'
import 'uno.css'
import { createApp } from 'vue'
import { vLoading } from 'element-plus'

import 'element-plus/dist/index.css'
import './styles/tokens.css'
// 引入 Element Plus 统一覆盖层
import './styles/element-plus-overrides.css'
// 引入全局 CSS (Tailwind + 重置 + 页面基座样式)
import './index.css'
// 引入gin-vue-admin前端初始化相关内容
import './core/gin-vue-admin'
// 引入封装的router
import router from '@/router/index'
import '@/permission'
import run from '@/core/gin-vue-admin'
import auth from '@/directive/auth'
import clickOutSide from '@/directive/clickOutSide'
import { store } from '@/pinia'
import App from './App.vue'
import '@/core/error-handel'

const app = createApp(App)

app
  .use(run)
  .use(store)
  .use(auth)
  .use(clickOutSide)
  .directive('loading', vLoading)
  .use(router)
  .mount('#app')

export default app

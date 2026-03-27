/*
 * gin-vue-admin web框架组
 *
 * */
import { register } from './global'
import type { App } from 'vue'

export default {
    install: (app: App) => {
        register(app)
    }
}

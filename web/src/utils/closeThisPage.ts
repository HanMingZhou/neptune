import { emitter } from '@/utils/bus'

export const closeThisPage = () => {
    emitter.emit('closeThisPage')
}

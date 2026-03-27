import type { App, DirectiveBinding } from 'vue'

interface ClickOutsideBinding {
    handler?: (e: Event) => void
    exclude?: HTMLElement[]
}

interface HTMLElementWithHandler extends HTMLElement {
    __clickOutsideHandler__?: (e: Event) => void
}

export default {
    install: (app: App) => {
        app.directive('click-outside', {
            mounted(el: HTMLElementWithHandler, binding: DirectiveBinding<ClickOutsideBinding | ((e: Event) => void)>) {
                const handler = (e: Event) => {
                    if (!el || el.contains(e.target as Node) || e.target === el) return
                    const value = binding.value
                    if (value && typeof value === 'object') {
                        if (
                            value.exclude &&
                            value.exclude.some(
                                (ex) => ex && ex.contains && ex.contains(e.target as Node)
                            )
                        )
                            return
                        if (typeof value.handler === 'function') value.handler(e)
                    } else if (typeof value === 'function') {
                        value(e)
                    }
                }

                el.__clickOutsideHandler__ = handler

                setTimeout(() => {
                    document.addEventListener('mousedown', handler)
                    document.addEventListener('touchstart', handler)
                }, 0)
            },
            unmounted(el: HTMLElementWithHandler) {
                const h = el.__clickOutsideHandler__
                if (h) {
                    document.removeEventListener('mousedown', h)
                    document.removeEventListener('touchstart', h)
                    delete el.__clickOutsideHandler__
                }
            }
        })
    }
}

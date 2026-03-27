export function addEventListen(
    target: any,
    event: string,
    handler: EventListenerOrEventListenerObject,
    capture: boolean | AddEventListenerOptions = false
) {
    if (
        target.addEventListener &&
        typeof target.addEventListener === 'function'
    ) {
        target.addEventListener(event, handler, capture)
    }
}

export function removeEventListen(
    target: any,
    event: string,
    handler: EventListenerOrEventListenerObject,
    capture: boolean | EventListenerOptions = false
) {
    if (
        target.removeEventListener &&
        typeof target.removeEventListener === 'function'
    ) {
        target.removeEventListener(event, handler, capture)
    }
}

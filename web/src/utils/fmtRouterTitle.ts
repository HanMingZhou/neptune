export const fmtTitle = (title: string, now: any) => {
    const reg = /\$\{(.+?)\}/
    const reg_g = /\$\{(.+?)\}/g
    const result = title.match(reg_g)
    if (result) {
        result.forEach((item) => {
            const match = item.match(reg)
            if (match) {
                const key = match[1]
                const value = now.params[key] || now.query[key]
                title = title.replace(item, value)
            }
        })
    }
    return title
}

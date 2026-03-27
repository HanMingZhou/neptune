import { fmtTitle } from '@/utils/fmtRouterTitle'
import config from '@/core/config'

export default function getPageTitle(pageTitle: string, route: any) {
    if (pageTitle) {
        const title = fmtTitle(pageTitle, route)
        return `${title} - ${config.appName}`
    }
    return `${config.appName}`
}

import service from '@/utils/request'

// @Tags jwt
// @Summary jwt加入黑名单
export const jsonInBlacklist = () => {
    return service({
        url: '/api/v1/jwt/blacklist',
        method: 'post'
    })
}

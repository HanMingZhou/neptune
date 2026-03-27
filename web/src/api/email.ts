import service from '@/utils/request'

// @Tags email
// @Summary 发送测试邮件
export const emailTest = (data: { to: string; subject: string; body: string }) => {
    return service({
        url: '/api/v1/email/test',
        method: 'post',
        data
    })
}

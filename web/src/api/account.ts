import service from '@/utils/request';

// 获取访问日志列表
export const getAccessLogList = (data?: any) => {
    return service({
        url: '/api/v1/account/access/log/list',
        method: 'post',
        data
    });
};

// 获取活跃会话列表
export const getActiveSessionList = (data?: any) => {
    return service({
        url: '/api/v1/account/active/session/list',
        method: 'post',
        data
    });
};

// 获取安全状态
export const getSecurityStatus = (params?: any) => {
    return service({
        url: '/api/v1/account/security/status',
        method: 'get',
        params
    });
};

// 修改密码
export const updatePassword = (data: any) => {
    return service({
        url: '/api/v1/account/password/update',
        method: 'post',
        data
    });
};

// 统一绑定(手机/邮箱)
export const bindAccount = (data: { type: number, value: string, code: string }) => {
    return service({
        url: '/api/v1/account/bind',
        method: 'post',
        data
    });
};

// 开启/关闭MFA
export const toggleMfa = (data: { enabled: boolean, code?: string }) => {
    return service({
        url: '/api/v1/account/mfa/toggle',
        method: 'post',
        data
    });
};

// MFA初始化（生成密钥和二维码）
export const mfaSetup = () => {
    return service({
        url: '/api/v1/account/mfa/setup',
        method: 'post'
    });
};

// MFA激活（验证码确认绑定）
export const mfaActivate = (data: { code: string }) => {
    return service({
        url: '/api/v1/account/mfa/activate',
        method: 'post',
        data
    });
};

// MFA二次验证登录（公共接口，无需JWT）
export const mfaLogin = (data: { mfaToken: string, code: string }) => {
    return service({
        url: '/api/v1/base/mfa/login',
        method: 'post',
        data
    });
};

// 生成AccessKey
export const generateAccessKey = () => {
    return service({
        url: '/api/v1/account/ak/generate',
        method: 'post'
    });
};

// 强制下线
export const killSession = (data: { logId: number }) => {
    return service({
        url: '/api/v1/account/active/session/kill',
        method: 'post',
        data
    });
};

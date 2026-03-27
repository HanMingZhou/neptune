import service from '@/utils/request'

export interface LoginData {
    username?: string
    password?: string
    captcha?: string
    img?: string
    [key: string]: any
}

export interface RegisterData {
    username?: string
    password?: string
    nickName?: string
    headerImg?: string
    authorityId?: string
    enable?: number
    phone?: string
    email?: string
}

export interface ChangePasswordData {
    password?: string
    newPassword?: string
    username?: string
}

export interface UserListParams {
    page: number
    pageSize: number
    ID?: number
    username?: string
    nickName?: string
    phone?: string
    email?: string
    enable?: number
}

export interface SetUserAuthParams {
    authorityId: number | string
}

export interface UserInfo {
    ID: number
    uuid: string
    nickName: string
    headerImg: string
    authority: any
    sideMode: string
    activeColor: string
    baseColor: string
    email: string
    phone: string
    enable: number
    [key: string]: any
}

// @Summary 用户登录
export const login = (data: LoginData) => {
    return service({
        url: '/api/v1/base/login',
        method: 'post',
        data: data
    })
}

// @Summary 获取验证码
export const captcha = () => {
    return service({
        url: '/api/v1/base/captcha',
        method: 'post'
    })
}

// @Summary 用户注册
export const register = (data: RegisterData) => {
    return service({
        url: '/api/v1/user/register',
        method: 'post',
        data: data
    })
}

// @Summary 修改密码
export const changePassword = (data: ChangePasswordData) => {
    return service({
        url: '/api/v1/user/password/update',
        method: 'post',
        data: data
    })
}

// @Tags User
// @Summary 分页获取用户列表
export const getUserList = (data: UserListParams) => {
    return service({
        url: '/api/v1/user/list',
        method: 'post',
        data: data
    })
}

// @Tags User
// @Summary 设置用户权限
export const setUserAuthority = (data: SetUserAuthParams) => {
    return service({
        url: '/api/v1/user/authority/update',
        method: 'post',
        data: data
    })
}

// @Tags SysUser
// @Summary 删除用户
export const deleteUser = (data: { id: number }) => {
    return service({
        url: '/api/v1/user/delete',
        method: 'post',
        data: data
    })
}

// @Tags SysUser
// @Summary 设置用户信息
export const setUserInfo = (data: UserInfo) => {
    return service({
        url: '/api/v1/user/info/update',
        method: 'post',
        data: data
    })
}

// @Tags SysUser
// @Summary 设置用户信息（自身）
export const setSelfInfo = (data: Partial<UserInfo>) => {
    return service({
        url: '/api/v1/user/self/info/update',
        method: 'post',
        data: data
    })
}

// @Tags SysUser
// @Summary 设置自身界面配置
export const setSelfSetting = (data: any) => {
    return service({
        url: '/api/v1/user/self/setting/update',
        method: 'post',
        data: data
    })
}

// @Tags User
// @Summary 设置用户权限(批量)
export const setUserAuthorities = (data: { id: number; authorityIds: number[] }) => {
    return service({
        url: '/api/v1/user/authorities/update',
        method: 'post',
        data: data
    })
}

// @Tags User
// @Summary 获取用户信息
export const getUserInfo = () => {
    return service({
        url: '/api/v1/user/info',
        method: 'get'
    })
}

export const resetPassword = (data: { ID: number }) => {
    return service({
        url: '/api/v1/user/password/reset',
        method: 'post',
        data: data
    })
}

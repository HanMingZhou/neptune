import service from '@/utils/request'

export interface MenuMeta {
    title: string
    icon?: string
    keepAlive?: boolean
    closeTab?: boolean
    defaultMenu?: boolean
}

export interface Menu {
    ID: number
    path: string
    name: string
    hidden?: boolean
    component: string
    sort: number
    meta: MenuMeta
    authoritys?: any[] // Authority[]
    children?: Menu[]
    parameters?: any[]
    menuBtn?: any[]
    parentId?: string
}

export interface AddMenuParams {
    path: string
    name: string
    hidden: boolean
    component: string
    sort: number
    meta: MenuMeta
    parentId: string
    menuBtn: any[]
    parameters: any[]
}

export interface MenuListParams {
    page: number
    pageSize: number
}

// @Summary 用户登录 获取动态路由
export const asyncMenu = () => {
    return service({
        url: '/api/v1/menu/get',
        method: 'post'
    })
}

// @Summary 获取menu列表
export const getMenuList = (data: MenuListParams) => {
    return service({
        url: '/api/v1/menu/list',
        method: 'post',
        data
    })
}

// @Summary 新增基础menu
export const addBaseMenu = (data: AddMenuParams) => {
    return service({
        url: '/api/v1/menu/add',
        method: 'post',
        data
    })
}

// @Summary 获取基础路由列表
export const getBaseMenuTree = () => {
    return service({
        url: '/api/v1/menu/tree',
        method: 'post'
    })
}

// @Summary 添加用户menu关联关系
export const addMenuAuthority = (data: { menus: Menu[]; authorityId: string }) => {
    return service({
        url: '/api/v1/menu/authority/update',
        method: 'post',
        data: data
    })
}

// @Summary 获取用户menu关联关系
export const getMenuAuthority = (data: { authorityId: string }) => {
    return service({
        url: '/api/v1/menu/authority/get',
        method: 'post',
        data: data
    })
}

// @Summary 删除menu
export const deleteBaseMenu = (data: { ID: number }) => {
    return service({
        url: '/api/v1/menu/delete',
        method: 'post',
        data: data
    })
}

// @Summary 修改menu列表
export const updateBaseMenu = (data: AddMenuParams & { ID: number }) => {
    return service({
        url: '/api/v1/menu/update',
        method: 'post',
        data: data
    })
}

// @Summary 根据id获取菜单
export const getBaseMenuById = (data: { id: number | string }) => {
    return service({
        url: '/api/v1/menu/get/id',
        method: 'post',
        data
    })
}

import service from '@/utils/request'

export interface AuthorityBtn {
    authorityId: number
    menuId: number
    selected: number[]
}

export const getAuthorityBtnApi = (data: { menuId: number; authorityId: number }) => {
    return service({
        url: '/api/v1/authority/btn/get',
        method: 'post',
        data
    })
}

export const setAuthorityBtnApi = (data: AuthorityBtn) => {
    return service({
        url: '/api/v1/authority/btn/update',
        method: 'post',
        data
    })
}

export const canRemoveAuthorityBtnApi = (params: { id: number }) => {
    return service({
        url: '/api/v1/authority/btn/delete',
        method: 'post',
        params
    })
}

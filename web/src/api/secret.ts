import service from '@/utils/request';

// 获取Secret列表
export const getSecretList = (params?: any) => {
    return service({
        url: '/api/v1/secret/list',
        method: 'get',
        params
    });
};

// 创建Secret
export const createSecret = (data: any) => {
    return service({
        url: '/api/v1/secret/add',
        method: 'post',
        data
    });
};

// 删除Secret
export const deleteSecret = (data: any) => {
    return service({
        url: '/api/v1/secret/delete',
        method: 'post',
        data
    });
};

// 更新Secret
export const updateSecret = (data: any) => {
    return service({
        url: '/api/v1/secret/update',
        method: 'post',
        data
    });
};

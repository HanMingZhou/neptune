import service from '@/utils/request';

// 获取SSH密钥列表
export const getSSHKeyList = (data?: any) => {
    return service({
        url: '/api/v1/sshkey/list',
        method: 'post',
        data
    });
};

// 创建SSH密钥
export const createSSHKey = (data: any) => {
    return service({
        url: '/api/v1/sshkey/add',
        method: 'post',
        data
    });
};

// 删除SSH密钥
export const deleteSSHKey = (data: any) => {
    return service({
        url: '/api/v1/sshkey/delete',
        method: 'post',
        data
    });
};

// 设置默认密钥
export const setDefaultSSHKey = (data: any) => {
    return service({
        url: '/api/v1/sshkey/default/update',
        method: 'post',
        data
    });
};

import service from '@/utils/request';

// 获取PVC列表
export const getPVCList = (params?: any) => {
    return service({
        url: '/api/v1/pvc/list',
        method: 'get',
        params
    });
};

// 创建PVC
export const createPVC = (data: any) => {
    return service({
        url: '/api/v1/pvc/add',
        method: 'post',
        data
    });
};

// 删除PVC
export const deletePVC = (data: any) => {
    return service({
        url: '/api/v1/pvc/delete',
        method: 'post',
        data
    });
};

// 更新PVC
export const updatePVC = (data: any) => {
    return service({
        url: '/api/v1/pvc/update',
        method: 'post',
        data
    });
};

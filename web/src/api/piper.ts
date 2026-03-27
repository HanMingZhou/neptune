import service from '@/utils/request';

// 获取Pipe列表
export const getPipeList = (params?: any) => {
    return service({
        url: '/api/v1/piper/list',
        method: 'get',
        params
    });
};

// 创建Pipe
export const createPipe = (data: any) => {
    return service({
        url: '/api/v1/piper/add',
        method: 'post',
        data
    });
};

// 删除Pipe
export const deletePipe = (data: any) => {
    return service({
        url: '/api/v1/piper/delete',
        method: 'post',
        data
    });
};

// 更新Pipe
export const updatePipe = (data: any) => {
    return service({
        url: '/api/v1/piper/update',
        method: 'post',
        data
    });
};

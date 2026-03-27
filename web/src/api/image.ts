import service from '@/utils/request';

// 获取镜像列表
export const getImageList = (params?: any) => {
    return service({
        url: '/api/v1/image/list',
        method: 'get',
        params
    });
};

// 获取镜像详情
export const getImageDetail = (id: number | string) => {
    return service({
        url: '/api/v1/image/detail',
        method: 'get',
        params: { id }
    });
};

// 创建镜像
export const createImage = (data: any) => {
    return service({
        url: '/api/v1/image/add',
        method: 'post',
        data
    });
};

// 更新镜像
export const updateImage = (data: any) => {
    return service({
        url: '/api/v1/image/update',
        method: 'post',
        data
    });
};

// 删除镜像
export const deleteImage = (data: any) => {
    return service({
        url: '/api/v1/image/delete',
        method: 'post',
        data
    });
};

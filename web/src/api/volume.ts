import service from '@/utils/request';

// 获取文件存储列表
export const getVolumeList = (params?: any) => {
    return service({
        url: '/api/v1/volume/list',
        method: 'get',
        params
    });
};

// 创建文件存储
export const createVolume = (data: any) => {
    return service({
        url: '/api/v1/volume/add',
        method: 'post',
        data
    });
};

// 扩容文件存储
export const expandVolume = (data: any) => {
    return service({
        url: '/api/v1/volume/expand',
        method: 'post',
        data
    });
};

// 删除文件存储
export const deleteVolume = (data: any) => {
    return service({
        url: '/api/v1/volume/delete',
        method: 'post',
        data
    });
};

// 获取可用区域列表
export const getAreaList = () => {
    return service({
        url: '/api/v1/volume/area/list',
        method: 'get'
    });
};

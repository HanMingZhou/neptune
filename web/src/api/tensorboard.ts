import service from '@/utils/request';

// 获取TensorBoard列表
export const getTensorBoardList = (params?: any) => {
    return service({
        url: '/api/v1/tensorboard/list',
        method: 'get',
        params
    });
};

// 创建TensorBoard
export const createTensorBoard = (data: any) => {
    return service({
        url: '/api/v1/tensorboard/add',
        method: 'post',
        data
    });
};

// 删除TensorBoard
export const deleteTensorBoard = (data: any) => {
    return service({
        url: '/api/v1/tensorboard/delete',
        method: 'post',
        data
    });
};

// 更新TensorBoard
export const updateTensorBoard = (data: any) => {
    return service({
        url: '/api/v1/tensorboard/update',
        method: 'post',
        data
    });
};

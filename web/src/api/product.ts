import service from '@/utils/request';

// 获取产品列表
export const getProductList = (params?: any) => {
    return service({
        url: '/api/v1/product/list',
        method: 'get',
        params
    });
};

// 获取产品详情
export const getProductDetail = (id: number | string) => {
    return service({
        url: '/api/v1/product/get',
        method: 'get',
        params: { id }
    });
};

// 获取产品筛选条件（地区、GPU型号等）
export const getProductFilters = (params?: any) => {
    return service({
        url: '/api/v1/product/filter/list',
        method: 'get',
        params
    });
};

import service from '@/utils/request';

// 获取仪表盘数据
export const getDashboardData = (data: { days?: number }) => {
    return service({
        url: '/api/v1/dashboard/get',
        method: 'post',
        data
    });
};

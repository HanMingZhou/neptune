import service from '@/utils/request';

// 获取训练任务列表
export const getTrainingJobList = (params: any) => {
    return service({
        url: '/api/v1/training/list',
        method: 'get',
        params
    });
};

// 创建训练任务
export const createTrainingJob = (data: any) => {
    return service({
        url: '/api/v1/training/add',
        method: 'post',
        data
    });
};

// 删除训练任务
export const deleteTrainingJob = (data: any) => {
    return service({
        url: '/api/v1/training/delete',
        method: 'post',
        data
    });
};

// 停止训练任务
export const stopTrainingJob = (data: any) => {
    return service({
        url: '/api/v1/training/stop',
        method: 'post',
        data
    });
};

// 获取训练任务详情
export const getTrainingJobDetail = (params: any) => {
    return service({
        url: '/api/v1/training/get',
        method: 'get',
        params
    });
};

// 获取训练任务 Pod 列表
export const getTrainingJobPods = (params: any) => {
    return service({
        url: '/api/v1/training/pod/list',
        method: 'get',
        params
    });
};

// 获取训练任务日志
export const getTrainingJobLogs = (params: any) => {
    return service({
        url: '/api/v1/training/log/list',
        method: 'get',
        params
    });
};

// 下载完整训练日志
export const downloadTrainingJobLogs = (params: any) => {
    return service({
        url: '/api/v1/training/log/download',
        method: 'get',
        params,
        responseType: 'blob'
    });
};

// 训练任务终端 (WebSocket)
export const getTrainingTerminal = (params: any) => {
    return service({
        url: '/api/v1/training/terminal',
        method: 'get',
        params
    });
};

// 训练任务日志流 (WebSocket)
export const getTrainingStreamLogs = (params: any) => {
    return service({
        url: '/api/v1/training/log/stream',
        method: 'get',
        params
    });
};

// 获取训练任务日志流 WebSocket URL
export const getTrainingLogStreamWsUrl = (id: number, token: string, taskName?: string, podIndex?: number) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const baseApi = import.meta.env.VITE_BASE_API || '';
    let url = `${protocol}//${window.location.host}${baseApi}/api/v1/training/log/stream?id=${id}&token=${token}`;
    if (taskName) url += `&taskName=${taskName}`;
    if (typeof podIndex === 'number') url += `&podIndex=${podIndex}`;
    return url;
};

// 获取训练任务终端 WebSocket URL
export const getTrainingTerminalWsUrl = (id: number, token: string, taskName?: string) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const baseApi = import.meta.env.VITE_BASE_API || '';
    let url = `${protocol}//${window.location.host}${baseApi}/api/v1/training/terminal?id=${id}&token=${token}`;
    if (taskName) url += `&taskName=${taskName}`;
    return url;
};

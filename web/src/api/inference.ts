import service from '@/utils/request';

// 获取推理服务列表
export const getInferenceServiceList = (data: any) => {
    return service({
        url: '/api/v1/inference/list',
        method: 'post',
        data
    });
};

// 创建推理服务
export const createInferenceService = (data: any) => {
    return service({
        url: '/api/v1/inference/add',
        method: 'post',
        data
    });
};

// 删除推理服务
export const deleteInferenceService = (data: any) => {
    return service({
        url: '/api/v1/inference/delete',
        method: 'post',
        data
    });
};

// 停止推理服务
export const stopInferenceService = (data: any) => {
    return service({
        url: '/api/v1/inference/stop',
        method: 'post',
        data
    });
};

// 启动推理服务
export const startInferenceService = (data: any) => {
    return service({
        url: '/api/v1/inference/start',
        method: 'post',
        data
    });
};

// 获取推理服务详情
export const getInferenceServiceDetail = (params: any) => {
    return service({
        url: '/api/v1/inference/get',
        method: 'get',
        params
    });
};

// 获取推理服务容器列表
export const getInferenceServicePods = (params: any) => {
    return service({
        url: '/api/v1/inference/pod/list',
        method: 'get',
        params
    });
};

// 获取推理服务日志
export const getInferenceServiceLogs = (params: any) => {
    return service({
        url: '/api/v1/inference/log/list',
        method: 'get',
        params
    });
};

// 创建 API Key
export const createInferenceApiKey = (data: any) => {
    return service({
        url: '/api/v1/inference/api/key/add',
        method: 'post',
        data
    });
};

// 获取 API Key 列表
export const getInferenceApiKeyList = (data: any) => {
    return service({
        url: '/api/v1/inference/api/key/list',
        method: 'post',
        data
    });
};

// 删除 API Key
export const deleteInferenceApiKey = (data: any) => {
    return service({
        url: '/api/v1/inference/api/key/delete',
        method: 'post',
        data
    });
};

// JWT Token 认证 (APISIX forward-auth)
export const inferenceAuth = () => {
    return service({
        url: '/api/v1/inference/auth',
        method: 'post'
    });
};

// API Key 认证 (APISIX forward-auth)
export const inferenceKeyAuth = () => {
    return service({
        url: '/api/v1/inference/key-auth',
        method: 'post'
    });
};


// 获取推理服务日志流 WebSocket URL
export const getInferenceLogStreamWsUrl = (id: number, token: string, podName?: string, container?: string) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const baseApi = import.meta.env.VITE_BASE_API || '';
    let url = `${protocol}//${window.location.host}${baseApi}/api/v1/inference/log/stream?id=${id}&token=${token}`;
    if (podName) url += `&podName=${podName}`;
    if (container) url += `&container=${container}`;
    return url;
};

// 获取推理服务终端 WebSocket URL
export const getInferenceTerminalWsUrl = (id: number, token: string, podName?: string, container?: string) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const baseApi = import.meta.env.VITE_BASE_API || '';
    let url = `${protocol}//${window.location.host}${baseApi}/api/v1/inference/terminal?id=${id}&token=${token}`;
    if (podName) url += `&podName=${podName}`;
    if (container) url += `&container=${container}`;
    return url;
};

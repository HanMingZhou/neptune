import service from '@/utils/request';

// 获取 Notebook 列表
export const getNotebookList = (params: any) => {
    return service({
        url: '/api/v1/notebook/list',
        method: 'get',
        params
    });
};

// 创建 Notebook
export const createNotebook = (data: any) => {
    return service({
        url: '/api/v1/notebook/add',
        method: 'post',
        data
    });
};

// 删除 Notebook
export const deleteNotebook = (data: any) => {
    return service({
        url: '/api/v1/notebook/delete',
        method: 'post',
        data
    });
};

// 停止 Notebook
export const stopNotebook = (data: any) => {
    return service({
        url: '/api/v1/notebook/stop',
        method: 'post',
        data
    });
};

// 启动 Notebook
export const startNotebook = (data: any) => {
    return service({
        url: '/api/v1/notebook/start',
        method: 'post',
        data
    });
};

// 获取 Notebook 详情
export const getNotebookDetail = (params: any) => {
    return service({
        url: '/api/v1/notebook/get',
        method: 'get',
        params
    });
};

// 获取 Notebook Pod 列表
export const getNotebookPods = (params: any) => {
    return service({
        url: '/api/v1/notebook/pod/list',
        method: 'get',
        params
    });
};

// 获取 Notebook 日志
export const getNotebookLogs = (params: any) => {
    return service({
        url: '/api/v1/notebook/log/list',
        method: 'get',
        params
    });
};

// 更新 Notebook
export const updateNotebook = (data: any) => {
    return service({
        url: '/api/v1/notebook/update',
        method: 'post',
        data
    });
};

// Notebook 终端 (WebSocket)
export const getNotebookTerminal = (params: any) => {
    return service({
        url: '/api/v1/notebook/terminal',
        method: 'get',
        params
    });
};

// 终端及身份认证
export const getNotebookAuth = (params: any) => {
    return service({
        url: '/api/v1/notebook/auth',
        method: 'get',
        params
    });
};

// Notebook 日志流 (WebSocket)
export const getNotebookLogStream = (params: any) => {
    return service({
        url: '/api/v1/notebook/log/stream',
        method: 'get',
        params
    });
};

import service from '@/utils/request'

// 获取财务概览
export const getOrderOverview = (params?: any) => {
    return service({
        url: '/api/v1/order/overview',
        method: 'get',
        params
    })
}

// 获取使用详情列表
export const getUsageList = (data?: any) => {
    return service({
        url: '/api/v1/order/usage/list',
        method: 'post',
        data
    })
}

// 获取收支明细列表
export const getTransactionList = (data?: any) => {
    return service({
        url: '/api/v1/order/transaction/list',
        method: 'post',
        data
    })
}

// 获取订单列表
export const getOrderList = (data?: any) => {
    return service({
        url: '/api/v1/order/order/list',
        method: 'post',
        data
    })
}

// 获取发票记录
export const getInvoiceList = (data?: any) => {
    return service({
        url: '/api/v1/order/invoice/list',
        method: 'post',
        data
    })
}

// 申请发票
export const applyInvoice = (data?: any) => {
    return service({
        url: '/api/v1/order/invoice/apply',
        method: 'post',
        data
    })
}

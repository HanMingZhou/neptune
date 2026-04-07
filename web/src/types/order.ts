import type { ResourceId } from './consoleResource'

export interface OverviewRankingItem {
  name?: string
  percent?: number
  amount?: number
  [key: string]: unknown
}

export interface ConsumptionTrendPoint {
  date?: string
  value?: number
  amount?: number
  [key: string]: unknown
}

export interface OrderOverviewData {
  balance?: number
  mtdEstimate?: number
  lastMonthSettlement?: number
  monthlyRanking?: OverviewRankingItem[]
  consumptionTrend?: ConsumptionTrendPoint[]
  [key: string]: unknown
}

export interface OrderOverviewMetric {
  amount: number
  change: number
  lastMonth: number
}

export interface OrderSettlementMetric {
  amount: number
  period: string
}

export interface InvoiceListItem {
  requestId: ResourceId
  amount?: number
  status?: string
  type?: string
  CreatedAt?: string
  title?: string
  [key: string]: unknown
}

export interface InvoiceRecord {
  id: ResourceId
  amount: string
  status: string
  type: string
  date: string
  title: string
}

export interface ApplyInvoiceForm {
  amount: number
  titleId: number
  addressId: number
}

export interface TransactionListItem {
  id?: ResourceId
  transactionNo?: string
  createdAt?: string | number
  type?: number
  amount?: number
  method?: number
  orderStatus?: number
  remark?: string
  resourceName?: string
  resourceType?: number
  orderNo?: string
  balanceAfter?: number
  [key: string]: unknown
}

export interface TransactionRecord {
  id: string
  time: string
  type: string
  typeLabel: string
  typeStyle: string
  amount: string
  method: string
  rawAmount: number
  orderStatus?: number
  remark?: string
  resourceName?: string
  resourceTypeText: string
  orderNo?: string
  balanceAfter: number
  statusStyle: string
  statusText: string
}

export interface RechargeForm {
  amount: number
  method: number
  remark: string
}

export interface UsageListItem {
  id?: ResourceId
  source?: string
  duration?: number
  quantity?: number
  chargeType?: number
  chargeTypeName?: string
  status?: number
  statusText?: string
  resourceType?: number
  resourceName?: string
  orderNo?: string
  unitPrice?: number
  createdAt?: string
  updatedAt?: string
  amount?: number
  [key: string]: unknown
}

export interface UsageMainItem {
  name: string
  percent: number
}

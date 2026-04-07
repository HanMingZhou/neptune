import { reactive, ref } from 'vue'
import { ElMessage, type FormRules } from 'element-plus'
import { applyInvoice, getInvoiceList } from '@/api/order'
import type { ApiResponse } from '@/utils/request'
import type { PageListData, Translator } from '@/types/consoleResource'
import type {
  ApplyInvoiceForm,
  InvoiceListItem,
  InvoiceRecord
} from '@/types/order'

interface UseInvoicePageOptions {
  t?: Translator
}

type InvoiceTab = 'records' | 'titles' | 'addresses'

export function useInvoicePage({ t }: UseInvoicePageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const loading = ref(false)
  const invoices = ref<InvoiceRecord[]>([])
  const page = ref(1)
  const pageSize = ref(10)
  const activeTab = ref<InvoiceTab>('records')
  const showApplyDialog = ref(false)
  const submitting = ref(false)

  const applyForm = reactive<ApplyInvoiceForm>({
    amount: 0,
    titleId: 1,
    addressId: 1
  })

  const applyRules = reactive<FormRules<ApplyInvoiceForm>>({
    amount: [
      {
        required: true,
        message: `${translate('general.pleaseInput')}${translate('order.amount')}`,
        trigger: 'blur'
      }
    ],
    titleId: [
      {
        required: true,
        message: translate('general.pleaseInput'),
        trigger: 'blur'
      }
    ],
    addressId: [
      {
        required: true,
        message: translate('general.pleaseInput'),
        trigger: 'blur'
      }
    ]
  })

  const fetchData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getInvoiceList({
        page: page.value,
        pageSize: pageSize.value
      })) as ApiResponse<PageListData<InvoiceListItem>>

      if (res.code === 0) {
        invoices.value = (res.data?.list || []).map((item) => ({
          id: item.requestId,
          amount: `¥${Number(item.amount || 0).toFixed(2)}`,
          status: item.status === 'SENT' ? 'Sent' : 'Processing',
          type:
            item.type === 'ENTERPRISE'
              ? translate('enterprise')
              : translate('personal'),
          date: item.CreatedAt ? item.CreatedAt.substring(0, 10) : '-',
          title: item.title || '-'
        }))
      }
    } catch (error) {
      console.error('获取发票列表失败:', error)
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const resetApplyForm = (): void => {
    applyForm.amount = 0
    applyForm.titleId = 1
    applyForm.addressId = 1
  }

  const submitApply = async (): Promise<void> => {
    submitting.value = true

    try {
      const res = await applyInvoice(applyForm)
      if (res.code === 0) {
        ElMessage.success(
          translate('order.applyInvoiceSuccess') || '申请发票成功'
        )
        showApplyDialog.value = false
        resetApplyForm()
        await fetchData()
      }
    } catch (error) {
      console.error('申请发票失败:', error)
    } finally {
      submitting.value = false
    }
  }

  return {
    activeTab,
    applyForm,
    applyRules,
    fetchData,
    invoices,
    loading,
    page,
    pageSize,
    resetApplyForm,
    showApplyDialog,
    submitApply,
    submitting
  }
}

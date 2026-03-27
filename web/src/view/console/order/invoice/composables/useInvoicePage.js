import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { applyInvoice, getInvoiceList } from '@/api/order'

export function useInvoicePage({ t }) {
  const translate = t || ((key) => key)

  const loading = ref(false)
  const invoices = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const activeTab = ref('records')
  const showApplyDialog = ref(false)
  const submitting = ref(false)

  const applyForm = reactive({
    amount: 0,
    titleId: 1,
    addressId: 1
  })

  const applyRules = {
    amount: [{ required: true, message: `${translate('general.pleaseInput')}${translate('order.amount')}`, trigger: 'blur' }],
    titleId: [{ required: true, message: translate('general.pleaseInput'), trigger: 'blur' }],
    addressId: [{ required: true, message: translate('general.pleaseInput'), trigger: 'blur' }]
  }

  const fetchData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getInvoiceList({
        page: page.value,
        pageSize: pageSize.value
      })

      if (res.code === 0) {
        invoices.value = (res.data?.list || []).map((item) => ({
          id: item.requestId,
          amount: `¥${item.amount.toFixed(2)}`,
          status: item.status === 'SENT' ? 'Sent' : 'Processing',
          type: item.type === 'ENTERPRISE' ? translate('enterprise') : translate('personal'),
          date: item.CreatedAt ? item.CreatedAt.substring(0, 10) : '-',
          title: item.title
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

  const resetApplyForm = () => {
    applyForm.amount = 0
    applyForm.titleId = 1
    applyForm.addressId = 1
  }

  const submitApply = async () => {
    submitting.value = true

    try {
      const res = await applyInvoice(applyForm)
      if (res.code === 0) {
        ElMessage.success(translate('order.applyInvoiceSuccess') || '申请发票成功')
        showApplyDialog.value = false
        resetApplyForm()
        fetchData()
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

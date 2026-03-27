import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  addBaseMenu,
  deleteBaseMenu,
  getBaseMenuById,
  getMenuList,
  updateBaseMenu
} from '@/api/menu'
import { canRemoveAuthorityBtnApi } from '@/api/authorityBtn'
import pathInfo from '@/pathInfo.json'
import { toLowerCase } from '@/utils/stringFun'

const createDefaultMenuForm = (parentId = 0) => ({
  ID: 0,
  path: '',
  name: '',
  hidden: false,
  parentId,
  component: '',
  sort: undefined,
  meta: {
    activeName: '',
    title: '',
    icon: '',
    defaultMenu: false,
    closeTab: false,
    keepAlive: false,
    transitionType: ''
  },
  parameters: [],
  menuBtn: []
})

const normalizeMenuForm = (menu = {}) => {
  const defaults = createDefaultMenuForm(menu.parentId ?? 0)

  return {
    ...defaults,
    ...menu,
    meta: {
      ...defaults.meta,
      ...(menu.meta || {})
    },
    parameters: Array.isArray(menu.parameters)
      ? menu.parameters.map((item) => ({ ...item }))
      : [],
    menuBtn: Array.isArray(menu.menuBtn)
      ? menu.menuBtn.map((item) => ({ ...item }))
      : []
  }
}

const buildMenuOptions = (menuData = [], currentId = 0, disabled = false) => {
  return menuData.map((item) => {
    const currentDisabled = disabled || item.ID === currentId
    const option = {
      title: item.meta?.title,
      ID: item.ID,
      disabled: currentDisabled
    }

    if (item.children?.length) {
      option.children = buildMenuOptions(item.children, currentId, currentDisabled)
    }

    return option
  })
}

export function useMenuManagementPage({ t }) {
  const translate = t || ((key) => key)

  const checkFlag = ref(false)
  const dialogFormVisible = ref(false)
  const form = reactive(createDefaultMenuForm())
  const isEdit = ref(false)
  const loading = ref(false)
  const searchKeyword = ref('')
  const tableData = ref([])

  const rules = reactive({
    path: [{ required: true, message: translate('inputRouteName'), trigger: 'blur' }],
    component: [{ required: true, message: translate('inputComponentPath'), trigger: 'blur' }],
    'meta.title': [{ required: true, message: translate('inputDisplayName'), trigger: 'blur' }]
  })

  const dialogTitle = computed(() => isEdit.value ? translate('edit') : translate('add'))

  const filteredTableData = computed(() => {
    const keyword = searchKeyword.value?.trim().toLowerCase()
    if (!keyword) {
      return tableData.value
    }

    const filterTree = (nodes = []) => {
      return nodes.reduce((acc, node) => {
        const title = node.meta?.title || ''
        const name = node.name || ''
        const isMatch = title.toLowerCase().includes(keyword) || name.toLowerCase().includes(keyword)
        const filteredChildren = node.children?.length ? filterTree(node.children) : []

        if (isMatch || filteredChildren.length > 0) {
          acc.push({
            ...node,
            children: filteredChildren
          })
        }

        return acc
      }, [])
    }

    return filterTree(tableData.value)
  })

  const menuOptions = computed(() => ([
    {
      ID: 0,
      title: translate('rootDirectory')
    },
    ...buildMenuOptions(tableData.value, form.ID, false)
  ]))

  const resetForm = (parentId = 0) => {
    checkFlag.value = false
    Object.assign(form, createDefaultMenuForm(parentId))
  }

  const getTableData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getMenuList()
      if (res.code === 0) {
        tableData.value = res.data || []
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error) {
      console.error('Failed to fetch menu list:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await getTableData()
  }

  const handleSearch = () => {}

  const handleResetSearch = () => {
    searchKeyword.value = ''
  }

  const addParameter = () => {
    if (!Array.isArray(form.parameters)) {
      form.parameters = []
    }

    form.parameters.push({
      type: 'query',
      key: '',
      value: ''
    })
  }

  const deleteParameter = (index) => {
    form.parameters.splice(index, 1)
  }

  const addButton = () => {
    if (!Array.isArray(form.menuBtn)) {
      form.menuBtn = []
    }

    form.menuBtn.push({
      name: '',
      desc: ''
    })
  }

  const deleteButton = async (index) => {
    const button = form.menuBtn[index]
    if (!button?.ID) {
      form.menuBtn.splice(index, 1)
      return
    }

    const res = await canRemoveAuthorityBtnApi({ id: button.ID })
    if (res.code === 0) {
      form.menuBtn.splice(index, 1)
    }
  }

  const fmtComponent = (component) => {
    const normalizedComponent = component.replace(/\\/g, '/')
    form.component = normalizedComponent

    const routeName = pathInfo[`/src/${normalizedComponent}`]
    if (routeName) {
      form.name = toLowerCase(routeName)
      form.path = form.name
    }
  }

  const changeName = () => {
    form.path = form.name
  }

  const closeDialog = () => {
    resetForm()
    dialogFormVisible.value = false
  }

  const openCreateDialog = (parentId = 0) => {
    resetForm(parentId)
    isEdit.value = false
    dialogFormVisible.value = true
  }

  const openEditDialog = async (id) => {
    const res = await getBaseMenuById({ id })
    if (res.code === 0 && res.data?.menu) {
      resetForm(res.data.menu.parentId ?? 0)
      Object.assign(form, normalizeMenuForm(res.data.menu))
      isEdit.value = true
      dialogFormVisible.value = true
    }
  }

  const handleSubmitMenu = async () => {
    const api = isEdit.value ? updateBaseMenu : addBaseMenu
    const res = await api(normalizeMenuForm(form))

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate('success')
      })
      await getTableData()
      closeDialog()
    }
  }

  const handleDeleteMenu = async (id) => {
    try {
      await ElMessageBox.confirm(
        translate('confirmDeleteMenu'),
        translate('tip'),
        {
          confirmButtonText: translate('confirm'),
          cancelButtonText: translate('cancel'),
          type: 'warning'
        }
      )

      const res = await deleteBaseMenu({ ID: id })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: translate('success')
        })
        await getTableData()
      }
    } catch {
      ElMessage({
        type: 'info',
        message: translate('deleteCancel')
      })
    }
  }

  return {
    addButton,
    addParameter,
    changeName,
    checkFlag,
    closeDialog,
    deleteButton,
    deleteParameter,
    dialogFormVisible,
    dialogTitle,
    filteredTableData,
    fmtComponent,
    form,
    getTableData,
    handleDeleteMenu,
    handleResetSearch,
    handleSearch,
    handleSubmitMenu,
    initialize,
    isEdit,
    loading,
    menuOptions,
    openCreateDialog,
    openEditDialog,
    rules,
    searchKeyword
  }
}

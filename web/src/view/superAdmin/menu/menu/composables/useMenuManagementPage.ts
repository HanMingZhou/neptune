import { computed, reactive, ref } from 'vue'
import type { FormRules } from 'element-plus'
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
import type { Translator } from '@/types/consoleResource'
import type { ApiResponse } from '@/utils/request'
import type {
  MenuForm,
  MenuListData,
  MenuOption,
  MenuTreeNode
} from '@/types/superAdmin'
import { getErrorMessage } from '@/utils/resourceValidators'
import { toLowerCase } from '@/utils/stringFun'

interface UseMenuManagementPageOptions {
  t?: Translator
}

type MenuFormSource = Partial<
  Omit<MenuForm, 'meta' | 'parameters' | 'menuBtn'>
> & {
  meta?: Partial<MenuForm['meta']>
  parameters?: MenuForm['parameters']
  menuBtn?: MenuForm['menuBtn']
}

const createDefaultMenuForm = (
  parentId: MenuForm['parentId'] = 0
): MenuForm => ({
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

const normalizeMenuForm = (menu: MenuFormSource = {}): MenuForm => {
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

const buildMenuOptions = (
  menuData: MenuTreeNode[] = [],
  currentId = 0,
  disabled = false
): MenuOption[] =>
  menuData.map((item) => {
    const currentDisabled = disabled || item.ID === currentId
    const option: MenuOption = {
      title: item.meta?.title,
      ID: item.ID,
      disabled: currentDisabled
    }

    if (item.children?.length) {
      option.children = buildMenuOptions(
        item.children,
        currentId,
        currentDisabled
      )
    }

    return option
  })

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useMenuManagementPage({
  t
}: UseMenuManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const checkFlag = ref(false)
  const dialogFormVisible = ref(false)
  const form = reactive<MenuForm>(createDefaultMenuForm())
  const isEdit = ref(false)
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(15)
  const searchKeyword = ref('')
  const tableData = ref<MenuTreeNode[]>([])
  const pagedTableData = ref<MenuTreeNode[]>([])
  const total = ref(0)

  const rules = computed<FormRules<MenuForm>>(() => ({
    path: [
      { required: true, message: translate('inputRouteName'), trigger: 'blur' }
    ],
    component: [
      {
        required: true,
        message: translate('inputComponentPath'),
        trigger: 'blur'
      }
    ],
    'meta.title': [
      {
        required: true,
        message: translate('inputDisplayName'),
        trigger: 'blur'
      }
    ]
  }))

  const dialogTitle = computed(() =>
    isEdit.value ? translate('edit') : translate('add')
  )

  const menuOptions = computed<MenuOption[]>(() => [
    {
      ID: 0,
      title: translate('rootDirectory')
    },
    ...buildMenuOptions(tableData.value, form.ID, false)
  ])

  const resetForm = (parentId: MenuForm['parentId'] = 0): void => {
    checkFlag.value = false
    Object.assign(form, createDefaultMenuForm(parentId))
  }

  const loadMenuTree = async (): Promise<void> => {
    const res = (await getMenuList()) as ApiResponse<MenuTreeNode[]>
    if (res.code === 0) {
      tableData.value = res.data ?? []
      return
    }

    ElMessage.error(res.msg || translate('failed'))
  }

  const getTableData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getMenuList({
        keyword: searchKeyword.value || undefined,
        page: page.value,
        pageSize: pageSize.value
      })) as ApiResponse<MenuListData>

      if (res.code === 0) {
        pagedTableData.value = res.data?.list ?? []
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? page.value
        pageSize.value = res.data?.pageSize ?? pageSize.value
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch menu list:', error)
      ElMessage.error(getErrorMessage(error, translate('failed')))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await Promise.all([getTableData(), loadMenuTree()])
  }

  const handleSearch = (): void => {
    searchKeyword.value = searchKeyword.value.trim()
    page.value = 1
    void getTableData()
  }

  const handleResetSearch = (): void => {
    searchKeyword.value = ''
    page.value = 1
    void getTableData()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void getTableData()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void getTableData()
  }

  const addParameter = (): void => {
    form.parameters.push({
      type: 'query',
      key: '',
      value: ''
    })
  }

  const deleteParameter = (index: number): void => {
    form.parameters.splice(index, 1)
  }

  const addButton = (): void => {
    form.menuBtn.push({
      name: '',
      desc: ''
    })
  }

  const deleteButton = async (index: number): Promise<void> => {
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

  const fmtComponent = (component: string): void => {
    const normalizedComponent = component.replace(/\\/g, '/')
    form.component = normalizedComponent

    const routeName = (pathInfo as Record<string, string>)[
      `/src/${normalizedComponent}`
    ]
    if (routeName) {
      form.name = toLowerCase(routeName)
      form.path = form.name
    }
  }

  const changeName = (): void => {
    form.path = form.name
  }

  const closeDialog = (): void => {
    resetForm()
    dialogFormVisible.value = false
  }

  const openCreateDialog = (parentId: MenuForm['parentId'] = 0): void => {
    resetForm(parentId)
    isEdit.value = false
    dialogFormVisible.value = true
  }

  const openEditDialog = async (id: number | string): Promise<void> => {
    const res = (await getBaseMenuById({ id })) as ApiResponse<{
      menu?: MenuTreeNode
    }>
    if (res.code === 0 && res.data?.menu) {
      resetForm(res.data.menu.parentId ?? 0)
      Object.assign(form, normalizeMenuForm(res.data.menu))
      isEdit.value = true
      dialogFormVisible.value = true
    }
  }

  const refreshAfterMutation = async (): Promise<void> => {
    await Promise.all([getTableData(), loadMenuTree()])
  }

  const handleSubmitMenu = async (): Promise<void> => {
    const api = isEdit.value ? updateBaseMenu : addBaseMenu
    const res = await api(normalizeMenuForm(form))

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate('success')
      })
      if (!isEdit.value) {
        page.value = 1
      }
      await refreshAfterMutation()
      closeDialog()
    }
  }

  const handleDeleteMenu = async (id: number): Promise<void> => {
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
        if (pagedTableData.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await refreshAfterMutation()
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage({
          type: 'info',
          message: translate('deleteCancel')
        })
      }
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
    fmtComponent,
    form,
    getTableData,
    handleCurrentChange,
    handleDeleteMenu,
    handleResetSearch,
    handleSearch,
    handleSizeChange,
    handleSubmitMenu,
    initialize,
    isEdit,
    loading,
    menuOptions,
    openCreateDialog,
    openEditDialog,
    page,
    pageSize,
    pagedTableData,
    rules,
    searchKeyword,
    total
  }
}


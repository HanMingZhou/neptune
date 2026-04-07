import { computed, reactive, ref } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  copyAuthority,
  createAuthority,
  deleteAuthority,
  getAuthorityList,
  updateAuthority
} from '@/api/authority'
import type { Translator } from '@/types/consoleResource'
import type {
  AuthorityCopyPayload,
  AuthorityDialogType,
  AuthorityForm,
  AuthorityOption,
  AuthorityTreeNode
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

interface UseAuthorityManagementPageOptions {
  t?: Translator
}

type EditableAuthorityRow = AuthorityTreeNode & Record<string, unknown>

const createDefaultForm = (): AuthorityForm => ({
  authorityId: '',
  authorityName: '',
  parentId: 0
})

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

const buildFilteredAuthorities = (
  nodes: AuthorityTreeNode[] = [],
  keyword = ''
): AuthorityTreeNode[] => {
  const normalizedKeyword = keyword.trim().toLowerCase()
  if (!normalizedKeyword) {
    return nodes
  }

  return nodes.reduce<AuthorityTreeNode[]>((accumulator, node) => {
    const matchesKeyword =
      String(node.authorityName || '')
        .toLowerCase()
        .includes(normalizedKeyword) ||
      String(node.authorityId ?? '').includes(keyword.trim())

    const filteredChildren = node.children?.length
      ? buildFilteredAuthorities(node.children, keyword)
      : []

    if (matchesKeyword || filteredChildren.length > 0) {
      accumulator.push({
        ...node,
        children: filteredChildren
      })
    }

    return accumulator
  }, [])
}

const buildAuthorityOptions = (
  authorityData: AuthorityTreeNode[] = [],
  disabledAuthorityId: AuthorityTreeNode['authorityId'] | null = null,
  inheritedDisabled = false
): AuthorityOption[] => {
  const currentId =
    disabledAuthorityId === null ||
    disabledAuthorityId === undefined ||
    disabledAuthorityId === ''
      ? null
      : Number(disabledAuthorityId)

  return authorityData.map((item) => {
    const isCurrent =
      currentId !== null && Number(item.authorityId) === currentId
    const option: AuthorityOption = {
      authorityId: item.authorityId,
      authorityName: item.authorityName,
      disabled: inheritedDisabled || isCurrent
    }

    if (item.children?.length) {
      option.children = buildAuthorityOptions(
        item.children,
        disabledAuthorityId,
        inheritedDisabled || isCurrent
      )
    }

    return option
  })
}

const createCopyPayload = (
  form: AuthorityForm,
  sourceAuthority: AuthorityTreeNode | null
): AuthorityCopyPayload => ({
  authority: {
    authorityId: Number(form.authorityId),
    authorityName: form.authorityName,
    dataAuthorityId: sourceAuthority?.dataAuthorityId || [],
    parentId: Number(form.parentId)
  },
  oldAuthorityId: Number(sourceAuthority?.authorityId || 0)
})

export function useAuthorityManagementPage({
  t
}: UseAuthorityManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const activeRow = ref<EditableAuthorityRow | Record<string, never>>({})
  const authorityFormVisible = ref(false)
  const copySource = ref<AuthorityTreeNode | null>(null)
  const dialogType = ref<AuthorityDialogType>('add')
  const disabledAuthorityId = ref<AuthorityTreeNode['authorityId'] | null>(null)
  const drawer = ref(false)
  const form = reactive<AuthorityForm>(createDefaultForm())
  const loading = ref(false)
  const searchKeyword = ref('')
  const submitting = ref(false)
  const tableData = ref<AuthorityTreeNode[]>([])

  const dialogTitle = computed(() => {
    if (dialogType.value === 'edit') {
      return translate('editRole')
    }

    if (dialogType.value === 'copy') {
      return translate('copyRole')
    }

    return translate('addRole')
  })

  const rules = computed<FormRules<AuthorityForm>>(() => ({
    authorityId: [
      { required: true, message: translate('roleId'), trigger: 'blur' },
      {
        pattern: /^[0-9]*[1-9][0-9]*$/,
        message: translate('mustBePositiveInteger'),
        trigger: 'blur'
      }
    ],
    authorityName: [
      { required: true, message: translate('roleName'), trigger: 'blur' }
    ],
    parentId: [
      {
        required: true,
        message: translate('selectParentRole'),
        trigger: 'change'
      }
    ]
  }))

  const filteredTableData = computed(() =>
    buildFilteredAuthorities(tableData.value, searchKeyword.value)
  )

  const authorityOptions = computed<AuthorityOption[]>(() => [
    {
      authorityId: 0,
      authorityName: translate('rootRoleTip')
    },
    ...buildAuthorityOptions(tableData.value, disabledAuthorityId.value)
  ])

  const resetForm = (): void => {
    Object.assign(form, createDefaultForm())
    copySource.value = null
    disabledAuthorityId.value = null
  }

  const getTableData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getAuthorityList()) as ApiResponse<AuthorityTreeNode[]>

      if (res.code === 0) {
        tableData.value = res.data ?? []
      } else {
        ElMessage.error(res.msg || translate('getRoleListFailed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch authority list:', error)
      ElMessage.error(translate('getRoleListFailedDetail'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async (): Promise<void> => {
    await getTableData()
  }

  const handleSearch = (): void => {
    searchKeyword.value = searchKeyword.value.trim()
  }

  const handleResetSearch = (): void => {
    searchKeyword.value = ''
  }

  const openDrawer = (row: AuthorityTreeNode): void => {
    activeRow.value = row as EditableAuthorityRow
    drawer.value = true
  }

  const closeDrawer = (): void => {
    drawer.value = false
    activeRow.value = {}
  }

  const changeRow = (key: string, value: unknown): void => {
    if (!activeRow.value) {
      return
    }

    ;(activeRow.value as Record<string, unknown>)[key] = value
  }

  const addAuthority = (parentId: number | string = 0): void => {
    resetForm()
    dialogType.value = 'add'
    form.parentId = Number(parentId)
    authorityFormVisible.value = true
  }

  const editAuthority = (row: AuthorityTreeNode): void => {
    resetForm()
    dialogType.value = 'edit'
    disabledAuthorityId.value = row.authorityId
    Object.assign(form, {
      authorityId: row.authorityId,
      authorityName: row.authorityName,
      parentId: row.parentId
    })
    authorityFormVisible.value = true
  }

  const copyAuthorityFunc = (row: AuthorityTreeNode): void => {
    resetForm()
    dialogType.value = 'copy'
    copySource.value = row
    disabledAuthorityId.value = row.authorityId
    Object.assign(form, {
      authorityId: row.authorityId,
      authorityName: row.authorityName,
      parentId: row.parentId
    })
    authorityFormVisible.value = true
  }

  const closeAuthorityForm = (): void => {
    authorityFormVisible.value = false
    resetForm()
  }

  const submitAuthorityForm = async (): Promise<void> => {
    submitting.value = true

    try {
      const authorityPayload = {
        authorityId: Number(form.authorityId),
        authorityName: form.authorityName,
        parentId: Number(form.parentId)
      }

      let res

      if (dialogType.value === 'add') {
        res = await createAuthority(authorityPayload)
      } else if (dialogType.value === 'edit') {
        res = await updateAuthority(authorityPayload)
      } else {
        res = await copyAuthority(createCopyPayload(form, copySource.value))
      }

      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        closeAuthorityForm()
        await getTableData()
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error: unknown) {
      console.error('Failed to submit authority form:', error)
      ElMessage.error(translate('failed'))
    } finally {
      submitting.value = false
    }
  }

  const deleteAuth = async (row: AuthorityTreeNode): Promise<void> => {
    try {
      await ElMessageBox.confirm(translate('deleteConfirm'), translate('tip'), {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      })

      const res = await deleteAuthority({
        authorityId: Number(row.authorityId)
      })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        await getTableData()
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error: unknown) {
      if (isDialogCancel(error)) {
        ElMessage.info(translate('deleteCancel'))
        return
      }

      ElMessage.error(translate('failed'))
    }
  }

  return {
    activeRow,
    addAuthority,
    authorityFormVisible,
    authorityOptions,
    changeRow,
    closeAuthorityForm,
    closeDrawer,
    copyAuthorityFunc,
    deleteAuth,
    dialogTitle,
    dialogType,
    drawer,
    editAuthority,
    filteredTableData,
    form,
    getTableData,
    handleResetSearch,
    handleSearch,
    initialize,
    loading,
    openDrawer,
    rules,
    searchKeyword,
    submitAuthorityForm,
    submitting,
    tableData
  }
}

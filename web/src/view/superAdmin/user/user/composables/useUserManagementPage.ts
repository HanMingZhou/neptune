import { computed, nextTick, reactive, ref } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAuthorityList } from '@/api/authority'
import {
  deleteUser,
  getUserList,
  register,
  resetPassword,
  setUserAuthorities,
  setUserInfo
} from '@/api/user'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type {
  AuthorityTreeNode,
  UserAuthority,
  UserForm,
  UserListData,
  UserResetPasswordForm,
  UserRow,
  UserSearchInfo
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

interface UseUserManagementPageOptions {
  t?: Translator
}

interface AuthorityChangeEvent {
  row: UserRow
  flag: boolean
}

interface SwitchEnableEvent {
  row: UserRow
  value: number
}

const PHONE_PATTERN = /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/
const EMAIL_PATTERN =
  /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/

const createDefaultSearchInfo = (): UserSearchInfo => ({
  username: '',
  nickName: '',
  phone: '',
  email: ''
})

const createDefaultUserForm = (): UserForm => ({
  ID: 0,
  userName: '',
  password: '',
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
  enable: 1,
  phone: '',
  email: ''
})

const createDefaultResetPasswordForm = (): UserResetPasswordForm => ({
  ID: '',
  userName: '',
  nickName: '',
  password: ''
})

const buildAuthorityOptions = (
  authorityData: AuthorityTreeNode[] = []
): UserAuthority[] =>
  authorityData.map((item) => {
    const option: UserAuthority = {
      authorityId: item.authorityId,
      authorityName: item.authorityName
    }

    if (item.children?.length) {
      option.children = buildAuthorityOptions(item.children)
    }

    return option
  })

const normalizeUserRow = (
  user: Partial<UserForm> & { authorities?: UserAuthority[] } = {}
): UserRow => ({
  ...createDefaultUserForm(),
  ...user,
  authorityIds: Array.isArray(user.authorities)
    ? user.authorities.map((item) => item.authorityId)
    : Array.isArray(user.authorityIds)
      ? [...user.authorityIds]
      : [],
  _authorityDirty: false
})

const cloneValue = <T>(value: T): T => JSON.parse(JSON.stringify(value)) as T

const getRowId = (row: Partial<UserRow>): ResourceId => row.id || row.ID || 0

const createRandomPassword = (length = 12): string => {
  const chars =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
  let password = ''

  for (let index = 0; index < length; index += 1) {
    password += chars.charAt(Math.floor(Math.random() * chars.length))
  }

  return password
}

const copyText = async (text: string): Promise<boolean> => {
  if (typeof navigator === 'undefined' || !navigator.clipboard?.writeText) {
    return false
  }

  await navigator.clipboard.writeText(text)
  return true
}

const isDialogCancel = (error: unknown): error is 'cancel' | 'close' =>
  error === 'cancel' || error === 'close'

export function useUserManagementPage({
  t
}: UseUserManagementPageOptions = {}) {
  const translate: Translator = t || ((key: string) => key)

  const addUserDialog = ref(false)
  const authOptions = ref<UserAuthority[]>([])
  const dialogFlag = ref<'add' | 'edit'>('add')
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(15)
  const resetPwdDialog = ref(false)
  const searchInfo = reactive<UserSearchInfo>(createDefaultSearchInfo())
  const tableData = ref<UserRow[]>([])
  const total = ref(0)
  const userInfo = reactive<UserForm>(createDefaultUserForm())
  const resetPwdInfo = reactive<UserResetPasswordForm>(
    createDefaultResetPasswordForm()
  )
  const authoritySnapshots = new Map<ResourceId, ResourceId[]>()

  const rules = computed<FormRules<UserForm>>(() => ({
    userName: [
      { required: true, message: translate('username'), trigger: 'blur' },
      { min: 5, message: translate('minChars', { min: 5 }), trigger: 'blur' }
    ],
    password: [
      { required: true, message: translate('password'), trigger: 'blur' },
      { min: 6, message: translate('minChars', { min: 6 }), trigger: 'blur' }
    ],
    nickName: [
      { required: true, message: translate('nickname'), trigger: 'blur' }
    ],
    phone: [
      {
        pattern: PHONE_PATTERN,
        message: translate('illegalPhone'),
        trigger: 'blur'
      }
    ],
    email: [
      {
        pattern: EMAIL_PATTERN,
        message: translate('illegalEmail'),
        trigger: 'blur'
      }
    ],
    authorityIds: [
      {
        type: 'array',
        required: true,
        message: translate('userRole'),
        trigger: 'change'
      }
    ]
  }))

  const resetSearchInfo = (): void => {
    Object.assign(searchInfo, createDefaultSearchInfo())
  }

  const resetUserForm = (): void => {
    Object.assign(userInfo, createDefaultUserForm())
  }

  const resetResetPasswordForm = (): void => {
    Object.assign(resetPwdInfo, createDefaultResetPasswordForm())
  }

  const getTableData = async (silent = false): Promise<void> => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = (await getUserList({
        page: page.value,
        pageSize: pageSize.value,
        ...searchInfo
      })) as ApiResponse<UserListData>

      if (res.code === 0) {
        tableData.value = (res.data?.list ?? []).map((item) =>
          normalizeUserRow(item)
        )
        total.value = res.data?.total ?? 0
        page.value = res.data?.page ?? 1
        pageSize.value = res.data?.pageSize ?? 20
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch user list:', error)
      ElMessage.error(translate('failed'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const getAuthorityOptions = async (): Promise<void> => {
    try {
      const res = (await getAuthorityList()) as ApiResponse<AuthorityTreeNode[]>
      if (res.code === 0) {
        authOptions.value = buildAuthorityOptions(res.data ?? [])
      } else {
        ElMessage.error(res.msg || translate('failed'))
      }
    } catch (error: unknown) {
      console.error('Failed to fetch authority list:', error)
      ElMessage.error(translate('failed'))
    }
  }

  const initialize = async (): Promise<void> => {
    await Promise.all([getTableData(), getAuthorityOptions()])
  }

  const onSubmit = (): void => {
    page.value = 1
    void getTableData()
  }

  const onReset = (): void => {
    page.value = 1
    resetSearchInfo()
    void getTableData()
  }

  const handleSizeChange = (value: number): void => {
    page.value = 1
    pageSize.value = value
    void getTableData()
  }

  const handleCurrentChange = (value: number): void => {
    page.value = value
    void getTableData()
  }

  const openCreateDialog = (): void => {
    resetUserForm()
    dialogFlag.value = 'add'
    addUserDialog.value = true
  }

  const openEditDialog = (row: UserRow): void => {
    const nextUserInfo = normalizeUserRow(cloneValue(row))

    resetUserForm()
    Object.assign(userInfo, createDefaultUserForm(), nextUserInfo, {
      password: '',
      authorityIds: [...(nextUserInfo.authorityIds || [])]
    })

    dialogFlag.value = 'edit'
    addUserDialog.value = true
  }

  const closeAddUserDialog = (): void => {
    addUserDialog.value = false
    resetUserForm()
  }

  const submitUserDialog = async (): Promise<void> => {
    const payload = {
      ...cloneValue(userInfo),
      authorityId: userInfo.authorityIds[0] || ''
    }

    const api = dialogFlag.value === 'add' ? register : setUserInfo
    const res = await api(payload)

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate('success')
      })
      closeAddUserDialog()
      await getTableData()
      return
    }

    ElMessage.error(res.msg || translate('failed'))
  }

  const openResetPasswordDialog = (row: UserRow): void => {
    resetResetPasswordForm()
    Object.assign(resetPwdInfo, {
      ID: getRowId(row),
      userName: row.userName,
      nickName: row.nickName
    })
    resetPwdDialog.value = true
  }

  const closeResetPwdDialog = (): void => {
    resetPwdDialog.value = false
    resetResetPasswordForm()
  }

  const generateRandomPassword = async (): Promise<void> => {
    const password = createRandomPassword()
    resetPwdInfo.password = password

    try {
      const copied = await copyText(password)
      ElMessage({
        type: copied ? 'success' : 'warning',
        message: copied ? translate('passwordCopied') : translate('copyFailed')
      })
    } catch {
      ElMessage({
        type: 'error',
        message: translate('copyFailed')
      })
    }
  }

  const confirmResetPassword = async (): Promise<void> => {
    if (!resetPwdInfo.password) {
      ElMessage({
        type: 'warning',
        message: translate('inputOrGeneratePassword')
      })
      return
    }

    const res = await resetPassword({
      ID: resetPwdInfo.ID,
      password: resetPwdInfo.password
    })

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg || translate('resetPasswordSuccess')
      })
      closeResetPwdDialog()
      return
    }

    ElMessage({
      type: 'error',
      message: res.msg || translate('resetPasswordFailed')
    })
  }

  const deleteUserFunc = async (row: UserRow): Promise<void> => {
    try {
      await ElMessageBox.confirm(translate('deleteConfirm'), translate('tip'), {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      })

      const res = await deleteUser({ id: getRowId(row) })
      if (res.code === 0) {
        ElMessage.success(translate('success'))
        if (tableData.value.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await getTableData()
      }
    } catch (error: unknown) {
      if (!isDialogCancel(error)) {
        ElMessage.error(translate('failed'))
      }
    }
  }

  const markAuthorityDirty = (row: UserRow): void => {
    row._authorityDirty = true
  }

  const changeAuthority = async ({
    row,
    flag
  }: AuthorityChangeEvent): Promise<void> => {
    const rowId = getRowId(row)

    if (flag) {
      authoritySnapshots.set(rowId, [...(row.authorityIds || [])])
      row._authorityDirty = false
      return
    }

    if (!row._authorityDirty) {
      authoritySnapshots.delete(rowId)
      return
    }

    await nextTick()
    const res = await setUserAuthorities({
      id: rowId,
      authorityIds: [...(row.authorityIds || [])]
    })

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate('roleSetSuccess')
      })
    } else {
      row.authorityIds = [...(authoritySnapshots.get(rowId) || [])]
      ElMessage.error(res.msg || translate('failed'))
    }

    row._authorityDirty = false
    authoritySnapshots.delete(rowId)
  }

  const switchEnable = async ({
    row,
    value
  }: SwitchEnableEvent): Promise<void> => {
    const previousEnable = value === 1 ? 2 : 1
    const res = await setUserInfo({
      ...cloneValue(row),
      enable: value
    })

    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: translate(value === 2 ? 'disableSuccess' : 'enableSuccess')
      })
      await getTableData(true)
      return
    }

    row.enable = previousEnable
    ElMessage.error(res.msg || translate('failed'))
  }

  return {
    addUserDialog,
    authOptions,
    changeAuthority,
    closeAddUserDialog,
    closeResetPwdDialog,
    confirmResetPassword,
    deleteUserFunc,
    dialogFlag,
    generateRandomPassword,
    getTableData,
    handleCurrentChange,
    handleSizeChange,
    initialize,
    loading,
    markAuthorityDirty,
    onReset,
    onSubmit,
    openCreateDialog,
    openEditDialog,
    openResetPasswordDialog,
    page,
    pageSize,
    resetPwdDialog,
    resetPwdInfo,
    rules,
    searchInfo,
    submitUserDialog,
    switchEnable,
    tableData,
    total,
    userInfo
  }
}


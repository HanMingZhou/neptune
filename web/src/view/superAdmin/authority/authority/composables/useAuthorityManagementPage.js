import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  copyAuthority,
  createAuthority,
  deleteAuthority,
  getAuthorityList,
  updateAuthority
} from '@/api/authority'

const createDefaultForm = () => ({
  authorityId: '',
  authorityName: '',
  parentId: 0
})

const isDialogCancel = (error) => error === 'cancel' || error === 'close'

const buildFilteredAuthorities = (nodes = [], keyword = '') => {
  const normalizedKeyword = keyword.trim().toLowerCase()
  if (!normalizedKeyword) {
    return nodes
  }

  return nodes.reduce((accumulator, node) => {
    const matchesKeyword =
      String(node.authorityName || '').toLowerCase().includes(normalizedKeyword) ||
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

const buildAuthorityOptions = (authorityData = [], disabledAuthorityId = null, inheritedDisabled = false) => {
  const currentId =
    disabledAuthorityId === null || disabledAuthorityId === undefined || disabledAuthorityId === ''
      ? null
      : Number(disabledAuthorityId)

  return authorityData.map((item) => {
    const isCurrent = currentId !== null && Number(item.authorityId) === currentId
    const option = {
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

const createCopyPayload = (form, sourceAuthority) => ({
  authority: {
    authorityId: Number(form.authorityId),
    authorityName: form.authorityName,
    dataAuthorityId: sourceAuthority?.dataAuthorityId || [],
    parentId: Number(form.parentId)
  },
  oldAuthorityId: Number(sourceAuthority?.authorityId || 0)
})

export function useAuthorityManagementPage({ t }) {
  const translate = t || ((key) => key)

  const activeRow = ref({})
  const authorityFormVisible = ref(false)
  const copySource = ref(null)
  const dialogType = ref('add')
  const disabledAuthorityId = ref(null)
  const drawer = ref(false)
  const form = reactive(createDefaultForm())
  const loading = ref(false)
  const searchKeyword = ref('')
  const submitting = ref(false)
  const tableData = ref([])

  const dialogTitle = computed(() => {
    if (dialogType.value === 'edit') {
      return translate('editRole')
    }

    if (dialogType.value === 'copy') {
      return translate('copyRole')
    }

    return translate('addRole')
  })

  const rules = computed(() => ({
    authorityId: [
      { required: true, message: translate('roleId'), trigger: 'blur' },
      {
        validator: (_, value, callback) => {
          if (!/^[0-9]*[1-9][0-9]*$/.test(String(value || ''))) {
            callback(new Error(translate('mustBePositiveInteger')))
            return
          }

          callback()
        },
        trigger: 'blur'
      }
    ],
    authorityName: [{ required: true, message: translate('roleName'), trigger: 'blur' }],
    parentId: [{ required: true, message: translate('selectParentRole'), trigger: 'change' }]
  }))

  const filteredTableData = computed(() => buildFilteredAuthorities(tableData.value, searchKeyword.value))

  const authorityOptions = computed(() => [
    {
      authorityId: 0,
      authorityName: translate('rootRoleTip')
    },
    ...buildAuthorityOptions(tableData.value, disabledAuthorityId.value)
  ])

  const resetForm = () => {
    Object.assign(form, createDefaultForm())
    copySource.value = null
    disabledAuthorityId.value = null
  }

  const getTableData = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }

    try {
      const res = await getAuthorityList()

      if (res.code === 0) {
        tableData.value = res.data || []
      } else {
        ElMessage.error(res.msg || translate('getRoleListFailed'))
      }
    } catch (error) {
      console.error('Failed to fetch authority list:', error)
      ElMessage.error(translate('getRoleListFailedDetail'))
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const initialize = async () => {
    await getTableData()
  }

  const handleSearch = () => {
    searchKeyword.value = searchKeyword.value.trim()
  }

  const handleResetSearch = () => {
    searchKeyword.value = ''
  }

  const openDrawer = (row) => {
    activeRow.value = row
    drawer.value = true
  }

  const closeDrawer = () => {
    drawer.value = false
    activeRow.value = {}
  }

  const changeRow = (key, value) => {
    if (!activeRow.value) {
      return
    }

    activeRow.value[key] = value
  }

  const addAuthority = (parentId = 0) => {
    resetForm()
    dialogType.value = 'add'
    form.parentId = Number(parentId)
    authorityFormVisible.value = true
  }

  const editAuthority = (row) => {
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

  const copyAuthorityFunc = (row) => {
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

  const closeAuthorityForm = () => {
    authorityFormVisible.value = false
    resetForm()
  }

  const submitAuthorityForm = async () => {
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
    } catch (error) {
      console.error('Failed to submit authority form:', error)
      ElMessage.error(translate('failed'))
    } finally {
      submitting.value = false
    }
  }

  const deleteAuth = async (row) => {
    try {
      await ElMessageBox.confirm(translate('deleteConfirm'), translate('tip'), {
        confirmButtonText: translate('confirm'),
        cancelButtonText: translate('cancel'),
        type: 'warning'
      })

      const res = await deleteAuthority({ authorityId: Number(row.authorityId) })
      if (res.code === 0) {
        ElMessage.success(res.msg || translate('success'))
        await getTableData()
        return
      }

      ElMessage.error(res.msg || translate('failed'))
    } catch (error) {
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

<template>
  <div>
    <div class="sticky top-0.5 z-10">
      <div class="flex flex-wrap items-center gap-3">
        <el-input v-model="filterText" class="min-w-[220px] flex-1" :placeholder="t('filter')" />
        <el-button class="shrink-0" type="primary" @click="relation">{{ t('confirm') }}</el-button>
      </div>
    </div>
    <div class="tree-content pt-3">
      <el-scrollbar>
        <el-tree
          ref="menuTree"
          :data="menuTreeData"
          :default-checked-keys="menuTreeIds"
          :props="menuDefaultProps"
          default-expand-all
          highlight-current
          node-key="ID"
          show-checkbox
          :filter-node-method="filterNode"
          @check="nodeChange"
        >
          <template #default="{ node, data }">
            <div class="custom-tree-node">
              <span class="custom-tree-node__label">{{ node.label }}</span>
              <div v-if="shouldShowActions(node, data)" class="custom-tree-node__actions">
                <el-tag
                  v-if="isDefaultRoute(data)"
                  effect="plain"
                  round
                  size="small"
                  type="warning"
                >
                  {{ t('homePage') }}
                </el-tag>
                <el-link
                  v-else-if="canSetDefaultRoute(node, data)"
                  :underline="false"
                  class="custom-tree-node__link"
                  type="primary"
                  @click.stop="setDefault(data)"
                >
                  {{ t('setHomePage') }}
                </el-link>
                <el-link
                  v-if="hasMenuButtons(data)"
                  :underline="false"
                  class="custom-tree-node__link"
                  type="primary"
                  @click.stop="OpenBtn(data)"
                >
                  {{ t('assignBtn') }}
                </el-link>
              </div>
            </div>
          </template>
        </el-tree>
      </el-scrollbar>
    </div>
    <el-dialog v-model="btnVisible" :title="t('assignBtn')" destroy-on-close>
      <el-table
        ref="btnTableRef"
        :data="btnData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column :label="t('buttonName')" prop="name" />
        <el-table-column :label="t('buttonDesc')" prop="desc" />
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">{{ t('cancel') }}</el-button>
          <el-button type="primary" @click="enterDialog">{{ t('confirm') }}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    getBaseMenuTree,
    getMenuAuthority,
    addMenuAuthority
  } from '@/api/menu'
  import { updateAuthority } from '@/api/authority'
  import { getAuthorityBtnApi, setAuthorityBtnApi } from '@/api/authorityBtn'
  import { nextTick, ref, watch, inject } from 'vue'
  import { ElMessage } from 'element-plus'

  const t = inject('t')

  defineOptions({
    name: 'Menus'
  })

  const props = defineProps({
    row: {
      default: function () {
        return {}
      },
      type: Object
    }
  })

  const emit = defineEmits(['changeRow'])
  const filterText = ref('')
  const menuTreeData = ref([])
  const menuTreeIds = ref([])
  const needConfirm = ref(false)
  const menuDefaultProps = ref({
    children: 'children',
    label: function (data) {
      return data.meta.title
    },
    disabled: function (data) {
      return props.row.defaultRouter === data.name
    }
  })

  const init = async () => {
    // 获取所有菜单树
    const res = await getBaseMenuTree()
    menuTreeData.value = res.data.menus
    const res1 = await getMenuAuthority({ authorityId: props.row.authorityId })
    const menus = res1.data.menus || []
    const arr = []
    menus.forEach((item) => {
      // 防止直接选中父级造成全选
      if (!menus.some((same) => same.parentId === item.menuId)) {
        arr.push(Number(item.menuId))
      }
    })
    menuTreeIds.value = arr
  }

  init()

  const setDefault = async (data) => {
    const res = await updateAuthority({
      authorityId: props.row.authorityId,
      AuthorityName: props.row.authorityName,
      parentId: props.row.parentId,
      defaultRouter: data.name
    })
    if (res.code === 0) {
      relation()
      emit('changeRow', 'defaultRouter', res.data.authority.defaultRouter)
    }
  }

  const isExternalRoute = (data) => {
    const name = String(data?.name || '')
    return name.startsWith('http://') || name.startsWith('https://')
  }

  const isLeafMenu = (data) => !Array.isArray(data?.children) || data.children.length === 0

  const hasMenuButtons = (data) => Array.isArray(data?.menuBtn) && data.menuBtn.length > 0

  const isDefaultRoute = (data) => props.row.defaultRouter === data?.name

  const canSetDefaultRoute = (node, data) => {
    if (!node?.checked) {
      return false
    }
    if (!data?.name || data.hidden) {
      return false
    }
    return isLeafMenu(data) && !isExternalRoute(data) && !isDefaultRoute(data)
  }

  const shouldShowActions = (node, data) => {
    if (hasMenuButtons(data)) {
      return true
    }
    return canSetDefaultRoute(node, data) || isDefaultRoute(data)
  }

  const nodeChange = () => {
    needConfirm.value = true
  }
  // 暴露给外层使用的切换拦截统一方法
  const enterAndNext = () => {
    relation()
  }
  // 关联树 确认方法
  const menuTree = ref(null)
  const relation = async () => {
    const checkArr = menuTree.value.getCheckedNodes(false, true)
    const res = await addMenuAuthority({
      menus: checkArr,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: t('menuSetSuccess')
      })
    }
  }

  defineExpose({ enterAndNext, needConfirm })

  const btnVisible = ref(false)

  const btnData = ref([])
  const multipleSelection = ref([])
  const btnTableRef = ref()
  let menuID = ''
  const OpenBtn = async (data) => {
    menuID = data.ID
    const res = await getAuthorityBtnApi({
      menuID: menuID,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      openDialog(data)
      await nextTick()
      if (res.data.selected) {
        res.data.selected.forEach((id) => {
          btnData.value.some((item) => {
            if (item.ID === id) {
              btnTableRef.value.toggleRowSelection(item, true)
            }
          })
        })
      }
    }
  }

  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  const openDialog = (data) => {
    btnVisible.value = true
    btnData.value = data.menuBtn
  }

  const closeDialog = () => {
    btnVisible.value = false
  }
  const enterDialog = async () => {
    const selected = multipleSelection.value.map((item) => item.ID)
    const res = await setAuthorityBtnApi({
      menuID,
      selected,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: t('setSuccess') })
      btnVisible.value = false
    }
  }

  const filterNode = (value, data) => {
    if (!value) return true
    // console.log(data.mate.title)
    return data.meta.title.indexOf(value) !== -1
  }

  watch(filterText, (val) => {
    menuTree.value.filter(val)
  })
</script>

<style scoped>
  .custom-tree-node {
    @apply flex w-full items-center gap-3 pr-3;
  }

  .custom-tree-node__label {
    @apply min-w-0 flex-1 truncate;
  }

  .custom-tree-node__actions {
    @apply flex shrink-0 items-center gap-3;
  }

  .custom-tree-node__link {
    @apply text-xs font-medium;
  }
</style>

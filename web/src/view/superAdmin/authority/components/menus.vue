<template>
  <div>
    <div class="sticky top-0.5 z-10">
      <div class="flex flex-wrap items-center gap-3">
        <el-input
          v-model="filterText"
          class="min-w-[220px] flex-1"
          :placeholder="t('filter')"
        />
        <el-button class="shrink-0" type="primary" @click="relation">{{
          t('confirm')
        }}</el-button>
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
              <div
                v-if="shouldShowActions(node, data)"
                class="custom-tree-node__actions"
              >
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
    <BaseDialog
      v-model="btnVisible"
      :title="t('assignBtn')"
      destroy-on-close
      width="720px"
    >
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
          <el-button type="primary" @click="enterDialog">{{
            t('confirm')
          }}</el-button>
        </div>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { addMenuAuthority, getBaseMenuTree, getMenuAuthority } from '@/api/menu'
import { updateAuthority } from '@/api/authority'
import { getAuthorityBtnApi, setAuthorityBtnApi } from '@/api/authorityBtn'
import BaseDialog from '@/components/base/BaseDialog.vue'
import { ElMessage } from 'element-plus'
import type { TableInstance, TreeInstance } from 'element-plus'
import { inject, nextTick, ref, watch } from 'vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type {
  AuthorityTreeNode,
  MenuButton,
  MenuTreeNode
} from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

defineOptions({
  name: 'Menus'
})

interface AuthorityPermissionRow extends Partial<AuthorityTreeNode> {
  authorityId?: ResourceId
  authorityName?: string
  parentId?: ResourceId
  defaultRouter?: string
}

interface MenuAuthorityRelation {
  menuId: ResourceId
  parentId?: ResourceId
}

interface TreeNodeState {
  checked?: boolean
}

const t = inject<Translator>('t', (key: string) => key)

const props = withDefaults(
  defineProps<{
    row?: AuthorityPermissionRow
  }>(),
  {
    row: () => ({})
  }
)

const emit = defineEmits<{
  changeRow: [key: string, value: unknown]
}>()

const filterText = ref('')
const menuTreeData = ref<MenuTreeNode[]>([])
const menuTreeIds = ref<number[]>([])
const needConfirm = ref(false)
const menuTree = ref<TreeInstance | null>(null)
const btnVisible = ref(false)
const btnData = ref<MenuButton[]>([])
const multipleSelection = ref<MenuButton[]>([])
const btnTableRef = ref<TableInstance | null>(null)

let menuID: ResourceId | '' = ''

const menuDefaultProps = {
  children: 'children',
  label: (data: MenuTreeNode) => data.meta?.title || '',
  disabled: (data: MenuTreeNode) => props.row.defaultRouter === data.name
}

const applyCheckedKeys = async (keys: number[]): Promise<void> => {
  await nextTick()
  menuTree.value?.setCheckedKeys([])
  if (keys.length > 0) {
    menuTree.value?.setCheckedKeys(keys)
  }
}

const init = async (): Promise<void> => {
  menuTreeIds.value = []
  await applyCheckedKeys([])

  const menuTreeRes = (await getBaseMenuTree()) as ApiResponse<{
    menus?: MenuTreeNode[]
  }>
  menuTreeData.value = menuTreeRes.data?.menus ?? []

  if (
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === ''
  ) {
    await applyCheckedKeys([])
    return
  }

  const menuAuthorityRes = (await getMenuAuthority({
    authorityId: props.row.authorityId
  })) as ApiResponse<{ menus?: MenuAuthorityRelation[] }>

  const menus = menuAuthorityRes.data?.menus ?? []
  const nextIds: number[] = []

  menus.forEach((item) => {
    if (!menus.some((same) => same.parentId === item.menuId)) {
      nextIds.push(Number(item.menuId))
    }
  })

  menuTreeIds.value = nextIds
  await applyCheckedKeys(nextIds)
}

watch(
  () => props.row.authorityId,
  () => {
    void init()
  },
  { immediate: true }
)

const setDefault = async (data: MenuTreeNode): Promise<void> => {
  if (
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === ''
  ) {
    return
  }

  const res = (await updateAuthority({
    authorityId: props.row.authorityId,
    authorityName: props.row.authorityName,
    parentId: props.row.parentId,
    defaultRouter: data.name
  })) as ApiResponse<{ authority?: { defaultRouter?: string } }>

  if (res.code === 0) {
    await relation()
    emit('changeRow', 'defaultRouter', res.data?.authority?.defaultRouter)
  }
}

const isExternalRoute = (data: MenuTreeNode): boolean => {
  const name = String(data.name || '')
  return name.startsWith('http://') || name.startsWith('https://')
}

const isLeafMenu = (data: MenuTreeNode): boolean =>
  !Array.isArray(data.children) || data.children.length === 0

const hasMenuButtons = (data: MenuTreeNode): boolean =>
  Array.isArray(data.menuBtn) && data.menuBtn.length > 0

const isDefaultRoute = (data: MenuTreeNode): boolean =>
  props.row.defaultRouter === data.name

const canSetDefaultRoute = (
  node: TreeNodeState | undefined,
  data: MenuTreeNode
): boolean => {
  if (!node?.checked) {
    return false
  }

  if (!data.name || data.hidden) {
    return false
  }

  return isLeafMenu(data) && !isExternalRoute(data) && !isDefaultRoute(data)
}

const shouldShowActions = (
  node: TreeNodeState | undefined,
  data: MenuTreeNode
): boolean => {
  if (hasMenuButtons(data)) {
    return true
  }

  return canSetDefaultRoute(node, data) || isDefaultRoute(data)
}

const nodeChange = (): void => {
  needConfirm.value = true
}

const enterAndNext = (): void => {
  void relation()
}

const relation = async (): Promise<void> => {
  if (
    !menuTree.value ||
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === ''
  ) {
    return
  }

  const checkedNodes = (menuTree.value.getCheckedNodes(false, true) ??
    []) as MenuTreeNode[]
  const res = await addMenuAuthority({
    menus: checkedNodes,
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

const OpenBtn = async (data: MenuTreeNode): Promise<void> => {
  if (
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === ''
  ) {
    return
  }

  menuID = data.ID
  const res = (await getAuthorityBtnApi({
    menuID,
    authorityId: props.row.authorityId
  })) as ApiResponse<{ selected?: ResourceId[] }>

  if (res.code === 0) {
    openDialog(data)
    await nextTick()

    res.data?.selected?.forEach((id) => {
      const target = btnData.value.find((item) => item.ID === id)
      if (target) {
        btnTableRef.value?.toggleRowSelection(target, true)
      }
    })
  }
}

const handleSelectionChange = (val: MenuButton[]): void => {
  multipleSelection.value = val
}

const openDialog = (data: MenuTreeNode): void => {
  btnVisible.value = true
  btnData.value = data.menuBtn ?? []
}

const closeDialog = (): void => {
  btnVisible.value = false
}

const enterDialog = async (): Promise<void> => {
  if (
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === '' ||
    menuID === ''
  ) {
    return
  }

  const selected = multipleSelection.value
    .map((item) => item.ID)
    .filter((id): id is number => typeof id === 'number')

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

const filterNode = (value: string, data: MenuTreeNode): boolean => {
  if (!value) {
    return true
  }

  return String(data.meta?.title || '').includes(value)
}

watch(filterText, (val: string) => {
  menuTree.value?.filter(val)
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

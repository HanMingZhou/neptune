<template>
  <div>
    <div
      class="sticky top-0.5 z-10 flex flex-col gap-3 bg-surface-light pb-2 dark:bg-surface-dark"
    >
      <div class="flex flex-wrap items-center gap-3">
        <el-input
          v-model="filterTextName"
          class="min-w-[180px] flex-1"
          :placeholder="t('filterName')"
        />
        <el-input
          v-model="filterTextPath"
          class="min-w-[180px] flex-1"
          :placeholder="t('filterPath')"
        />
        <el-button class="shrink-0" type="primary" @click="authApiEnter">{{
          t('confirm')
        }}</el-button>
      </div>
      <div class="flex justify-between items-center">
        <div class="flex flex-wrap gap-2">
          <el-button size="small" @click="expandAll">{{
            isAllExpanded ? t('collapseAll') : t('expandAll')
          }}</el-button>
          <el-button size="small" @click="selectAll">{{
            isAllSelected ? t('deselectAll') : t('selectAll')
          }}</el-button>
        </div>
      </div>
    </div>
    <div class="tree-content">
      <el-scrollbar>
        <el-tree
          ref="apiTree"
          :data="apiTreeData"
          :default-checked-keys="apiTreeIds"
          :props="apiDefaultProps"
          :default-expand-all="isAllExpanded"
          highlight-current
          node-key="onlyId"
          show-checkbox
          :filter-node-method="filterNode"
          @check="nodeChange"
        >
          <template #default="{ data }">
            <div class="flex items-center justify-between w-full pr-1">
              <span>{{ data.description }} </span>
              <el-tooltip :content="data.path">
                <span
                  class="max-w-[240px] break-all overflow-ellipsis overflow-hidden"
                  >{{ data.path }}</span
                >
              </el-tooltip>
            </div>
          </template>
        </el-tree>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getAllApis } from '@/api/api'
import { UpdateCasbin, getPolicyPathByAuthorityId } from '@/api/casbin'
import { ElMessage } from 'element-plus'
import type { TreeInstance } from 'element-plus'
import { inject, nextTick, ref, watch } from 'vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type { ApiListItem, AuthorityTreeNode } from '@/types/superAdmin'
import type { ApiResponse } from '@/utils/request'

defineOptions({
  name: 'Apis'
})

interface AuthorityPermissionRow extends Partial<AuthorityTreeNode> {
  authorityId?: ResourceId
}

interface ApiTreeLeaf extends ApiListItem {
  onlyId: string
}

interface ApiTreeGroupNode {
  ID: string
  onlyId: string
  description: string
  children: ApiTreeLeaf[]
}

type ApiTreeNode = ApiTreeGroupNode | ApiTreeLeaf

interface CasbinInfo {
  path: string
  method: string
}

interface TreeStoreNode {
  expanded: boolean
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

const apiDefaultProps = {
  children: 'children',
  label: 'description'
}

const filterTextName = ref('')
const filterTextPath = ref('')
const apiTreeData = ref<ApiTreeGroupNode[]>([])
const apiTreeIds = ref<string[]>([])
const activeUserId = ref<ResourceId | ''>('')
const needConfirm = ref(false)
const apiTree = ref<TreeInstance | null>(null)
const isAllExpanded = ref(true)
const isAllSelected = ref(false)

const buildApiTree = (apis: ApiListItem[] = []): ApiTreeGroupNode[] => {
  const apiGroups: Record<string, ApiTreeLeaf[]> = {}

  apis.forEach((item) => {
    const normalizedItem: ApiTreeLeaf = {
      ...item,
      onlyId: `p:${item.path}m:${item.method}`
    }

    if (Object.prototype.hasOwnProperty.call(apiGroups, item.apiGroup)) {
      apiGroups[item.apiGroup].push(normalizedItem)
      return
    }

    apiGroups[item.apiGroup] = [normalizedItem]
  })

  return Object.keys(apiGroups).map((key) => ({
    ID: key,
    onlyId: key,
    description: `${key} ${t('group')}`,
    children: apiGroups[key]
  }))
}

const init = async (): Promise<void> => {
  const allApisRes = (await getAllApis()) as ApiResponse<{
    apis?: ApiListItem[]
  }>
  apiTreeData.value = buildApiTree(allApisRes.data?.apis ?? [])

  if (
    props.row.authorityId === undefined ||
    props.row.authorityId === null ||
    props.row.authorityId === ''
  ) {
    activeUserId.value = ''
    apiTreeIds.value = []
    return
  }

  const policyRes = (await getPolicyPathByAuthorityId({
    authorityId: props.row.authorityId
  })) as ApiResponse<{ paths?: CasbinInfo[] }>

  activeUserId.value = props.row.authorityId
  apiTreeIds.value = (policyRes.data?.paths ?? []).map(
    (item) => `p:${item.path}m:${item.method}`
  )
}

void init()

const nodeChange = (): void => {
  needConfirm.value = true
}

const enterAndNext = (): void => {
  void authApiEnter()
}

const authApiEnter = async (): Promise<void> => {
  if (!apiTree.value || activeUserId.value === '') {
    return
  }

  const checkedNodes = (apiTree.value.getCheckedNodes(true) ??
    []) as ApiTreeLeaf[]
  const casbinInfos: CasbinInfo[] = checkedNodes
    .filter((item) => Boolean(item.path) && Boolean(item.method))
    .map((item) => ({
      path: item.path,
      method: item.method
    }))

  const res = await UpdateCasbin({
    authorityId: activeUserId.value,
    casbinInfos
  })

  if (res.code === 0) {
    ElMessage({ type: 'success', message: t('apiSetSuccess') })
  }
}

defineExpose({
  needConfirm,
  enterAndNext
})

const filterNode = (_value: string, data: ApiTreeNode): boolean => {
  if (!filterTextName.value && !filterTextPath.value) {
    return true
  }

  const matchesName =
    !filterTextName.value ||
    String(data.description || '').includes(filterTextName.value)
  const matchesPath =
    !filterTextPath.value ||
    String((data as ApiTreeLeaf).path || '').includes(filterTextPath.value)

  return matchesName && matchesPath
}

watch([filterTextName, filterTextPath], () => {
  apiTree.value?.filter('')
})

const expandAll = (): void => {
  isAllExpanded.value = !isAllExpanded.value

  nextTick(() => {
    const store = (
      apiTree.value as
        | (TreeInstance & {
            store?: {
              _getAllNodes?: () => TreeStoreNode[]
            }
          })
        | null
    )?.store

    const nodes = store?._getAllNodes?.() ?? []
    nodes.forEach((node) => {
      node.expanded = isAllExpanded.value
    })
  })
}

const extractApiTreeKeys = (nodes: ApiTreeNode[], keys: string[]): void => {
  nodes.forEach((node) => {
    keys.push(node.onlyId)
    if ('children' in node && Array.isArray(node.children)) {
      extractApiTreeKeys(node.children, keys)
    }
  })
}

const selectAll = (): void => {
  isAllSelected.value = !isAllSelected.value

  nextTick(() => {
    if (!apiTree.value) {
      return
    }

    if (isAllSelected.value) {
      const allKeys: string[] = []
      extractApiTreeKeys(apiTreeData.value, allKeys)
      apiTree.value.setCheckedKeys(allKeys)
    } else {
      apiTree.value.setCheckedKeys([])
    }

    nodeChange()
  })
}
</script>

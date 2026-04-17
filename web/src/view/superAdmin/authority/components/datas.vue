<template>
  <div>
    <warning-bar
      :title="t('dataPermissionTip')"
      href="https://plugin.gin-vue-admin.com/#/layout/newPluginInfo?id=36"
    />
    <div class="sticky top-0.5 z-10 my-4">
      <div class="flex flex-wrap items-center gap-3">
        <div class="flex flex-wrap gap-3">
          <el-button type="primary" @click="all">{{
            t('selectAll')
          }}</el-button>
          <el-button type="primary" @click="self">{{
            t('selfRole')
          }}</el-button>
          <el-button type="primary" @click="selfAndChildren">{{
            t('selfAndChildrenRole')
          }}</el-button>
        </div>
        <el-button class="sm:ml-auto" type="primary" @click="authDataEnter">{{
          t('confirm')
        }}</el-button>
      </div>
    </div>
    <div class="pt-4">
      <el-checkbox-group v-model="dataAuthorityId" @change="selectAuthority">
        <el-checkbox
          v-for="(item, key) in authoritys"
          :key="key"
          :label="item"
          >{{ item.authorityName }}</el-checkbox
        >
      </el-checkbox-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { setDataAuthority } from '@/api/authority'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ElMessage } from 'element-plus'
import { inject, ref, watch } from 'vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type { AuthorityTreeNode } from '@/types/superAdmin'

defineOptions({
  name: 'Datas'
})

interface AuthoritySelectionItem {
  authorityId: ResourceId
  authorityName: string
}

interface AuthorityPermissionRow extends Partial<AuthorityTreeNode> {
  authorityId?: ResourceId
  dataAuthorityId?: AuthoritySelectionItem[]
  children?: AuthorityTreeNode[]
}

const t = inject<Translator>('t', (key: string) => key)

const props = withDefaults(
  defineProps<{
    authority?: AuthorityTreeNode[]
    row?: AuthorityPermissionRow
  }>(),
  {
    authority: () => [],
    row: () => ({})
  }
)

const authoritys = ref<AuthoritySelectionItem[]>([])
const needConfirm = ref(false)
const dataAuthorityId = ref<AuthoritySelectionItem[]>([])

const roundAuthority = (authoritysData: AuthorityTreeNode[] = []): void => {
  authoritysData.forEach((item) => {
    authoritys.value.push({
      authorityId: item.authorityId,
      authorityName: item.authorityName
    })

    if (item.children?.length) {
      roundAuthority(item.children)
    }
  })
}

const init = (): void => {
  authoritys.value = []
  dataAuthorityId.value = []
  roundAuthority(props.authority)

  props.row.dataAuthorityId?.forEach((item) => {
    const selectedAuthority = authoritys.value.find(
      (authority) => authority.authorityId === item.authorityId
    )
    if (selectedAuthority) {
      dataAuthorityId.value.push(selectedAuthority)
    }
  })
}

watch(
  () => [props.row.authorityId, props.row.dataAuthorityId, props.authority],
  () => {
    init()
  },
  {
    immediate: true,
    deep: true
  }
)

const emit = defineEmits<{
  changeRow: [key: string, value: AuthoritySelectionItem[]]
}>()

const enterAndNext = (): void => {
  void authDataEnter()
}

const all = (): void => {
  dataAuthorityId.value = [...authoritys.value]
  emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
  needConfirm.value = true
}

const self = (): void => {
  dataAuthorityId.value = authoritys.value.filter(
    (item) => item.authorityId === props.row.authorityId
  )
  emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
  needConfirm.value = true
}

const getChildrenId = (
  row: AuthorityPermissionRow | AuthorityTreeNode,
  arrBox: ResourceId[]
): void => {
  if (
    row.authorityId !== undefined &&
    row.authorityId !== null &&
    row.authorityId !== ''
  ) {
    arrBox.push(row.authorityId)
  }

  row.children?.forEach((item) => {
    getChildrenId(item, arrBox)
  })
}

const selfAndChildren = (): void => {
  const authorityIds: ResourceId[] = []
  getChildrenId(props.row, authorityIds)
  dataAuthorityId.value = authoritys.value.filter((item) =>
    authorityIds.includes(item.authorityId)
  )
  emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
  needConfirm.value = true
}

const authDataEnter = async (): Promise<void> => {
  const res = await setDataAuthority(props.row)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: t('resourceSetSuccess') })
  }
}

const selectAuthority = (): void => {
  dataAuthorityId.value = dataAuthorityId.value.filter(Boolean)
  emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
  needConfirm.value = true
}

defineExpose({
  enterAndNext,
  needConfirm
})
</script>

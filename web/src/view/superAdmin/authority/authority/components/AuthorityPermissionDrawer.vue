<template>
  <BaseDrawer
    v-model="visibleModel"
    destroy-on-close
    :size="600"
    :title="t('roleConfig')"
    @close="emit('close')"
  >
    <el-tabs :before-leave="handleBeforeLeave" type="border-card">
      <el-tab-pane :label="t('roleMenus')" name="menus">
        <Menus ref="menusRef" :row="activeRow" @change-row="handleChangeRow" />
      </el-tab-pane>
      <el-tab-pane :label="t('roleApis')" name="apis">
        <Apis ref="apisRef" :row="activeRow" @change-row="handleChangeRow" />
      </el-tab-pane>
      <el-tab-pane :label="t('dataPermission')" name="datas">
        <Datas
          ref="datasRef"
          :authority="authority"
          :row="activeRow"
          @change-row="handleChangeRow"
        />
      </el-tab-pane>
    </el-tabs>
  </BaseDrawer>
</template>

<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import BaseDrawer from '@/components/base/BaseDrawer.vue'
import Apis from '../../components/apis.vue'
import Datas from '../../components/datas.vue'
import Menus from '../../components/menus.vue'
import type { ResourceId, Translator } from '@/types/consoleResource'
import type { AuthorityTreeNode } from '@/types/superAdmin'

interface AuthorityPermissionRef {
  authorityId: ResourceId
}

interface AuthorityPermissionRow extends AuthorityTreeNode {
  dataAuthorityId?: AuthorityPermissionRef[]
  defaultRouter?: string
}

interface PermissionPaneHandle extends ComponentPublicInstance {
  needConfirm?: boolean
  enterAndNext?: () => void
}

const props = withDefaults(
  defineProps<{
    activeRow?: AuthorityPermissionRow
    authority?: AuthorityTreeNode[]
    modelValue?: boolean
  }>(),
  {
    activeRow: () => ({}) as AuthorityPermissionRow,
    authority: () => [],
    modelValue: false
  }
)

const emit = defineEmits<{
  'change-row': [key: string, value: unknown]
  close: []
  'update:modelValue': [value: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const menusRef = ref<PermissionPaneHandle | null>(null)
const apisRef = ref<PermissionPaneHandle | null>(null)
const datasRef = ref<PermissionPaneHandle | null>(null)

const handleChangeRow = (key: string, value: unknown): void => {
  emit('change-row', key, value)
}

const paneMap: Record<string, typeof menusRef> = {
  apis: apisRef,
  datas: datasRef,
  menus: menusRef
}

const handleBeforeLeave = (
  _nextPaneName: string | number,
  previousPaneName: string | number
): void => {
  const currentPane = paneMap[String(previousPaneName)]
  if (currentPane?.value?.needConfirm) {
    currentPane.value.enterAndNext?.()
    currentPane.value.needConfirm = false
  }
}
</script>

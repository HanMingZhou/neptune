<template>
  <el-drawer
    v-if="visibleModel"
    v-model="visibleModel"
    :before-close="handleBeforeClose"
    :size="600"
    :title="t('roleConfig')"
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
  </el-drawer>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import Apis from '../../components/apis.vue'
import Datas from '../../components/datas.vue'
import Menus from '../../components/menus.vue'

const props = defineProps({
  activeRow: {
    type: Object,
    default: () => ({})
  },
  authority: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['change-row', 'close', 'update:modelValue'])
const t = inject('t', (key) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const menusRef = ref(null)
const apisRef = ref(null)
const datasRef = ref(null)

const handleChangeRow = (key, value) => {
  emit('change-row', key, value)
}

const handleBeforeClose = (done) => {
  emit('close')
  done()
}

const handleBeforeLeave = (_, previousPaneName) => {
  const paneMap = {
    apis: apisRef,
    datas: datasRef,
    menus: menusRef
  }

  const currentPane = paneMap[previousPaneName]
  if (currentPane?.value?.needConfirm) {
    currentPane.value.enterAndNext()
    currentPane.value.needConfirm = false
  }
}
</script>

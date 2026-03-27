<template>
  <el-drawer
    v-model="visibleModel"
    :size="600"
    :before-close="handleBeforeClose"
    :show-close="false"
  >
    <template #header>
      <div class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ dialogTitle }}</span>
        <div class="flex gap-3">
          <el-button @click="requestClose">{{ t('cancel') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ t('confirm') }}</el-button>
        </div>
      </div>
    </template>

    <div class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3 mb-6">
      <span class="material-icons text-amber-500">info</span>
      <span class="text-sm text-amber-700 dark:text-amber-400">{{ t('menuTip') }}</span>
    </div>

    <el-form
      v-if="visibleModel"
      ref="menuFormRef"
      :inline="true"
      :model="form"
      :rules="rules"
      label-position="top"
    >
      <div class="border-b border-border-light dark:border-border-dark pb-6 mb-6">
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">{{ t('baseInfo') }}</h3>
        <el-row class="w-full">
          <el-col :span="24">
            <el-form-item :label="t('componentPath')" prop="component">
              <ComponentsCascader :component="form.component" @change="$emit('component-change', $event)" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="12">
            <el-form-item :label="t('displayName')" prop="meta.title">
              <el-input v-model="form.meta.title" autocomplete="off" :placeholder="t('displayName')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('routeName')" prop="path">
              <el-input
                v-model="form.name"
                autocomplete="off"
                :placeholder="t('englishOnly')"
                @change="$emit('name-change')"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="border-b border-border-light dark:border-border-dark pb-6 mb-6">
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">{{ t('routeConfig') }}</h3>
        <el-row class="w-full">
          <el-col :span="12">
            <el-form-item :label="t('parentMenu')">
              <el-cascader
                v-model="form.parentId"
                style="width: 100%"
                :disabled="!isEdit"
                :options="menuOptions"
                :props="{
                  checkStrictly: true,
                  label: 'title',
                  value: 'ID',
                  disabled: 'disabled',
                  emitPath: false
                }"
                :show-all-levels="false"
                filterable
                :placeholder="t('parentMenu')"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="path">
              <template #label>
                <div class="inline-flex items-center h-4">
                  <span>{{ t('path') }}</span>
                  <el-checkbox class="ml-2" v-model="checkFlagModel">{{ t('addParam') }}</el-checkbox>
                </div>
              </template>
              <el-input
                v-model="form.path"
                :disabled="!checkFlagModel"
                autocomplete="off"
                :placeholder="t('paramHint')"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="border-b border-border-light dark:border-border-dark pb-6 mb-6">
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">{{ t('displaySetting') }}</h3>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item :label="t('menuIcon')" prop="meta.icon">
              <IconPicker v-model="form.meta.icon" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('sort')" prop="sort">
              <el-input v-model.number="form.sort" autocomplete="off" :placeholder="t('sort')" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('status')">
              <el-select v-model="form.hidden" style="width: 100%" :placeholder="t('status')">
                <el-option :value="false" :label="t('visible')" />
                <el-option :value="true" :label="t('hidden')" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="border-b border-border-light dark:border-border-dark pb-6 mb-6">
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">{{ t('advancedConfig') }}</h3>
        <el-row class="w-full">
          <el-col :span="12">
            <el-form-item prop="meta.activeName">
              <template #label>
                <div class="flex items-center gap-2">
                  <span>{{ t('activeMenu') }}</span>
                  <el-tooltip :content="t('activeMenuTip')" placement="top" effect="light">
                    <span class="material-icons text-slate-400 text-[16px] cursor-help">help</span>
                  </el-tooltip>
                </div>
              </template>
              <el-input v-model="form.meta.activeName" :placeholder="form.name || t('activeMenu')" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('keepAlive')" prop="meta.keepAlive">
              <el-select v-model="form.meta.keepAlive" style="width: 100%" :placeholder="t('keepAlive')">
                <el-option :value="false" :label="t('disable')" />
                <el-option :value="true" :label="t('enable')" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item :label="t('closeTab')" prop="meta.closeTab">
              <el-select v-model="form.meta.closeTab" style="width: 100%" :placeholder="t('closeTab')">
                <el-option :value="false" :label="t('disable')" />
                <el-option :value="true" :label="t('enable')" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item>
              <template #label>
                <div class="flex items-center gap-2">
                  <span>{{ t('basePage') }}</span>
                  <el-tooltip :content="t('basePageTip')" placement="top" effect="light">
                    <span class="material-icons text-slate-400 text-[16px] cursor-help">help</span>
                  </el-tooltip>
                </div>
              </template>
              <el-select v-model="form.meta.defaultMenu" style="width: 100%" :placeholder="t('basePage')">
                <el-option :value="false" :label="t('disable')" />
                <el-option :value="true" :label="t('enable')" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item>
              <template #label>
                <div class="flex items-center gap-2">
                  <span>{{ t('transitionAnim') }}</span>
                </div>
              </template>
              <el-select v-model="form.meta.transitionType" style="width: 100%" :placeholder="t('all')" clearable>
                <el-option value="fade" :label="t('fade')" />
                <el-option value="slide" :label="t('slide')" />
                <el-option value="zoom" :label="t('zoom')" />
                <el-option value="none" label="None" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <div class="border-b border-border-light dark:border-border-dark pb-6 mb-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-slate-700 dark:text-slate-300">{{ t('parameterConfig') }}</h3>
          <el-button type="primary" size="small" @click="$emit('add-parameter')">
            <span class="material-icons text-[16px] mr-1">add</span>
            {{ t('addParameter') }}
          </el-button>
        </div>

        <el-table :data="form.parameters" style="width: 100%" class="rounded-lg overflow-hidden">
          <el-table-column align="center" prop="type" :label="t('parameterType')" width="150">
            <template #default="scope">
              <el-select v-model="scope.row.type" :placeholder="t('pleaseSelect')" size="small">
                <el-option key="query" value="query" label="query" />
                <el-option key="params" value="params" label="params" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="key" :label="t('parameterKey')" width="150">
            <template #default="scope">
              <el-input v-model="scope.row.key" size="small" :placeholder="t('parameterKey')" />
            </template>
          </el-table-column>
          <el-table-column align="center" prop="value" :label="t('parameterValue')">
            <template #default="scope">
              <el-input v-model="scope.row.value" size="small" :placeholder="t('parameterValue')" />
            </template>
          </el-table-column>
          <el-table-column align="center" :label="t('actions')" width="100">
            <template #default="scope">
              <el-button type="danger" size="small" @click="$emit('delete-parameter', scope.$index)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="mb-2">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-slate-700 dark:text-slate-300">{{ t('buttonConfig') }}</h3>
          <el-button type="primary" size="small" @click="$emit('add-button')">
            <span class="material-icons text-[16px] mr-1">add</span>
            {{ t('addButton') }}
          </el-button>
        </div>

        <el-table :data="form.menuBtn" style="width: 100%" class="rounded-lg overflow-hidden">
          <el-table-column align="center" prop="name" :label="t('buttonName')" width="150">
            <template #default="scope">
              <el-input v-model="scope.row.name" size="small" :placeholder="t('buttonName')" />
            </template>
          </el-table-column>
          <el-table-column align="center" prop="desc" :label="t('remark')">
            <template #default="scope">
              <el-input v-model="scope.row.desc" size="small" :placeholder="t('remark')" />
            </template>
          </el-table-column>
          <el-table-column align="center" :label="t('actions')" width="100">
            <template #default="scope">
              <el-button type="danger" size="small" @click="$emit('delete-button', scope.$index)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-form>
  </el-drawer>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import IconPicker from '@/view/superAdmin/menu/icon.vue'
import ComponentsCascader from '@/view/superAdmin/menu/components/components-cascader.vue'

const props = defineProps({
  checkFlag: {
    type: Boolean,
    default: false
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  form: {
    type: Object,
    required: true
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  menuOptions: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  rules: {
    type: Object,
    required: true
  }
})

const emit = defineEmits([
  'add-button',
  'add-parameter',
  'close',
  'component-change',
  'delete-button',
  'delete-parameter',
  'name-change',
  'submit',
  'update:check-flag',
  'update:modelValue'
])

const t = inject('t', (key) => key)
const menuFormRef = ref(null)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const checkFlagModel = computed({
  get: () => props.checkFlag,
  set: (value) => emit('update:check-flag', value)
})

const requestClose = () => {
  emit('close')
}

const handleBeforeClose = (done) => {
  emit('close')
  done()
}

const submitForm = async () => {
  if (!menuFormRef.value) {
    return
  }

  await menuFormRef.value.validate()
  emit('submit')
}
</script>

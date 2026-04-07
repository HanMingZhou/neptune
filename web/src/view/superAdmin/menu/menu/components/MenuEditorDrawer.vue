<template>
  <BaseFormDrawer
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :model="form"
    :rules="rules"
    :size="600"
    :submit-text="t('confirm')"
    :title="dialogTitle"
    label-position="top"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <template #prepend>
      <div
        class="bg-amber-500/10 border border-amber-500/20 rounded-lg px-4 py-3 flex items-center gap-3 mb-6"
      >
        <span class="material-icons text-amber-500">info</span>
        <span class="text-sm text-amber-700 dark:text-amber-400">{{
          t('menuTip')
        }}</span>
      </div>
    </template>

    <div class="menu-editor-drawer">
      <div
        class="border-b border-border-light dark:border-border-dark pb-6 mb-6"
      >
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">
          {{ t('baseInfo') }}
        </h3>
        <div class="menu-editor-grid menu-editor-grid--single">
          <el-form-item
            :label="t('componentPath')"
            prop="component"
            class="menu-editor-field--full"
          >
            <ComponentsCascader
              :component="form.component"
              @change="emit('component-change', $event)"
            />
          </el-form-item>
        </div>
        <div class="menu-editor-grid menu-editor-grid--regular">
          <el-form-item :label="t('displayName')" prop="meta.title">
            <el-input
              v-model="form.meta.title"
              autocomplete="off"
              :placeholder="t('displayName')"
            />
          </el-form-item>
          <el-form-item :label="t('routeName')" prop="path">
            <el-input
              v-model="form.name"
              autocomplete="off"
              :placeholder="t('englishOnly')"
              @change="emit('name-change')"
            />
          </el-form-item>
        </div>
      </div>

      <div
        class="border-b border-border-light dark:border-border-dark pb-6 mb-6"
      >
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">
          {{ t('routeConfig') }}
        </h3>
        <div class="menu-editor-grid menu-editor-grid--regular">
          <el-form-item :label="t('parentMenu')">
            <el-cascader
              v-model="form.parentId"
              style="width: 100%"
              :disabled="!isEdit"
              :options="menuOptions"
              :props="menuCascaderProps"
              :show-all-levels="false"
              filterable
              :placeholder="t('parentMenu')"
            />
          </el-form-item>
          <el-form-item prop="path">
            <template #label>
              <div class="menu-editor-label menu-editor-label--between">
                <span class="menu-editor-label__text">{{ t('path') }}</span>
                <el-checkbox v-model="checkFlagModel">{{
                  t('addParam')
                }}</el-checkbox>
              </div>
            </template>
            <el-input
              v-model="form.path"
              :disabled="!checkFlagModel"
              autocomplete="off"
              :placeholder="t('paramHint')"
            />
          </el-form-item>
        </div>
      </div>

      <div
        class="border-b border-border-light dark:border-border-dark pb-6 mb-6"
      >
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">
          {{ t('displaySetting') }}
        </h3>
        <div class="menu-editor-grid menu-editor-grid--compact">
          <el-form-item :label="t('menuIcon')" prop="meta.icon">
            <IconPicker v-model="form.meta.icon" />
          </el-form-item>
          <el-form-item :label="t('sort')" prop="sort">
            <el-input
              v-model.number="form.sort"
              autocomplete="off"
              :placeholder="t('sort')"
            />
          </el-form-item>
          <el-form-item :label="t('status')">
            <el-select
              v-model="form.hidden"
              style="width: 100%"
              :placeholder="t('status')"
            >
              <el-option :value="false" :label="t('visible')" />
              <el-option :value="true" :label="t('hidden')" />
            </el-select>
          </el-form-item>
        </div>
      </div>

      <div
        class="border-b border-border-light dark:border-border-dark pb-6 mb-6"
      >
        <h3 class="font-bold text-slate-700 dark:text-slate-300 mb-4">
          {{ t('advancedConfig') }}
        </h3>
        <div class="menu-editor-grid menu-editor-grid--regular mb-4">
          <el-form-item prop="meta.activeName">
            <template #label>
              <div class="menu-editor-label">
                <span class="menu-editor-label__text">{{
                  t('activeMenu')
                }}</span>
                <el-tooltip
                  :content="t('activeMenuTip')"
                  placement="top"
                  effect="light"
                >
                  <span
                    class="material-icons text-slate-400 text-[16px] cursor-help"
                    >help</span
                  >
                </el-tooltip>
              </div>
            </template>
            <el-input
              v-model="form.meta.activeName"
              :placeholder="form.name || t('activeMenu')"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item :label="t('keepAlive')" prop="meta.keepAlive">
            <el-select
              v-model="form.meta.keepAlive"
              style="width: 100%"
              :placeholder="t('keepAlive')"
            >
              <el-option :value="false" :label="t('disable')" />
              <el-option :value="true" :label="t('enable')" />
            </el-select>
          </el-form-item>
        </div>
        <div class="menu-editor-grid menu-editor-grid--compact">
          <el-form-item :label="t('closeTab')" prop="meta.closeTab">
            <el-select
              v-model="form.meta.closeTab"
              style="width: 100%"
              :placeholder="t('closeTab')"
            >
              <el-option :value="false" :label="t('disable')" />
              <el-option :value="true" :label="t('enable')" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="menu-editor-label">
                <span class="menu-editor-label__text">{{ t('basePage') }}</span>
                <el-tooltip
                  :content="t('basePageTip')"
                  placement="top"
                  effect="light"
                >
                  <span
                    class="material-icons text-slate-400 text-[16px] cursor-help"
                    >help</span
                  >
                </el-tooltip>
              </div>
            </template>
            <el-select
              v-model="form.meta.defaultMenu"
              style="width: 100%"
              :placeholder="t('basePage')"
            >
              <el-option :value="false" :label="t('disable')" />
              <el-option :value="true" :label="t('enable')" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('transitionAnim')">
            <el-select
              v-model="form.meta.transitionType"
              style="width: 100%"
              :placeholder="t('all')"
              clearable
            >
              <el-option value="fade" :label="t('fade')" />
              <el-option value="slide" :label="t('slide')" />
              <el-option value="zoom" :label="t('zoom')" />
              <el-option value="none" label="None" />
            </el-select>
          </el-form-item>
        </div>
      </div>

      <div
        class="border-b border-border-light dark:border-border-dark pb-6 mb-6"
      >
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-slate-700 dark:text-slate-300">
            {{ t('parameterConfig') }}
          </h3>
          <el-button type="primary" size="small" @click="emit('add-parameter')">
            <span class="material-icons text-[16px] mr-1">add</span>
            {{ t('addParameter') }}
          </el-button>
        </div>

        <el-table
          :data="form.parameters"
          style="width: 100%"
          class="rounded-lg overflow-hidden"
        >
          <el-table-column
            align="center"
            prop="type"
            :label="t('parameterType')"
            width="150"
          >
            <template #default="scope">
              <el-select
                v-model="scope.row.type"
                :placeholder="t('pleaseSelect')"
                size="small"
              >
                <el-option key="query" value="query" label="query" />
                <el-option key="params" value="params" label="params" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column
            align="center"
            prop="key"
            :label="t('parameterKey')"
            width="150"
          >
            <template #default="scope">
              <el-input
                v-model="scope.row.key"
                size="small"
                :placeholder="t('parameterKey')"
              />
            </template>
          </el-table-column>
          <el-table-column
            align="center"
            prop="value"
            :label="t('parameterValue')"
          >
            <template #default="scope">
              <el-input
                v-model="scope.row.value"
                size="small"
                :placeholder="t('parameterValue')"
              />
            </template>
          </el-table-column>
          <el-table-column align="center" :label="t('actions')" width="100">
            <template #default="scope">
              <el-button
                type="danger"
                size="small"
                @click="emit('delete-parameter', scope.$index)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="mb-2">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-slate-700 dark:text-slate-300">
            {{ t('buttonConfig') }}
          </h3>
          <el-button type="primary" size="small" @click="emit('add-button')">
            <span class="material-icons text-[16px] mr-1">add</span>
            {{ t('addButton') }}
          </el-button>
        </div>

        <el-table
          :data="form.menuBtn"
          style="width: 100%"
          class="rounded-lg overflow-hidden"
        >
          <el-table-column
            align="center"
            prop="name"
            :label="t('buttonName')"
            width="150"
          >
            <template #default="scope">
              <el-input
                v-model="scope.row.name"
                size="small"
                :placeholder="t('buttonName')"
              />
            </template>
          </el-table-column>
          <el-table-column align="center" prop="desc" :label="t('remark')">
            <template #default="scope">
              <el-input
                v-model="scope.row.desc"
                size="small"
                :placeholder="t('remark')"
              />
            </template>
          </el-table-column>
          <el-table-column align="center" :label="t('actions')" width="100">
            <template #default="scope">
              <el-button
                type="danger"
                size="small"
                @click="emit('delete-button', scope.$index)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </BaseFormDrawer>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import BaseFormDrawer from '@/components/base/BaseFormDrawer.vue'
import type { Translator } from '@/types/consoleResource'
import type { MenuForm, MenuOption } from '@/types/superAdmin'
import IconPicker from '@/view/superAdmin/menu/icon.vue'
import ComponentsCascader from '@/view/superAdmin/menu/components/components-cascader.vue'

const props = withDefaults(
  defineProps<{
    checkFlag?: boolean
    dialogTitle?: string
    form: MenuForm
    isEdit?: boolean
    menuOptions?: MenuOption[]
    modelValue?: boolean
    rules: FormRules<MenuForm>
  }>(),
  {
    checkFlag: false,
    dialogTitle: '',
    isEdit: false,
    menuOptions: () => [],
    modelValue: false
  }
)

const emit = defineEmits<{
  'add-button': []
  'add-parameter': []
  close: []
  'component-change': [component: string]
  'delete-button': [index: number]
  'delete-parameter': [index: number]
  'name-change': []
  submit: []
  'update:check-flag': [value: boolean]
  'update:modelValue': [value: boolean]
}>()

const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const checkFlagModel = computed({
  get: () => props.checkFlag,
  set: (value: boolean) => emit('update:check-flag', value)
})

const menuCascaderProps = {
  checkStrictly: true,
  label: 'title',
  value: 'ID',
  disabled: 'disabled',
  emitPath: false
}
</script>

<style scoped>
.menu-editor-grid {
  display: grid;
  gap: 1rem 1.25rem;
  align-items: start;
}

.menu-editor-grid--single {
  grid-template-columns: minmax(0, 1fr);
  margin-bottom: 1rem;
}

.menu-editor-grid--regular {
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
}

.menu-editor-grid--compact {
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
}

.menu-editor-field--full {
  grid-column: 1 / -1;
}

.menu-editor-label {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
  width: 100%;
  min-height: 1.5rem;
  line-height: 1.375;
}

.menu-editor-label--between {
  justify-content: space-between;
}

.menu-editor-label__text {
  min-width: 0;
}

.menu-editor-label--between .menu-editor-label__text {
  flex: 1 1 auto;
}

.menu-editor-drawer :deep(.el-form-item) {
  margin-bottom: 0;
}

.menu-editor-drawer :deep(.el-form-item__label) {
  display: flex;
  align-items: center;
  min-height: 1.5rem;
  margin-bottom: 0.5rem;
  line-height: 1.375;
  white-space: normal;
}

.menu-editor-drawer :deep(.el-checkbox) {
  margin-right: 0;
}
</style>

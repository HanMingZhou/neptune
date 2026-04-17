<template>
  <BaseFormDrawer
    v-model="visibleModel"
    :cancel-text="t('cancel')"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    :model="form"
    :rules="rules"
    :size="480"
    :submit-text="t('confirm')"
    :title="dialogFlag === 'add' ? t('addUser') : t('editUser')"
    @close="emit('close')"
    @submit="emit('submit')"
  >
    <el-form-item
      v-if="dialogFlag === 'add'"
      :label="t('username')"
      prop="userName"
    >
      <el-input v-model="form.userName" :placeholder="t('username')" />
    </el-form-item>
    <el-form-item
      v-if="dialogFlag === 'add'"
      :label="t('password')"
      prop="password"
    >
      <el-input
        v-model="form.password"
        :placeholder="t('password')"
        show-password
      />
    </el-form-item>
    <el-form-item :label="t('nickname')" prop="nickName">
      <el-input v-model="form.nickName" :placeholder="t('nickname')" />
    </el-form-item>
    <el-form-item :label="t('phone')" prop="phone">
      <el-input v-model="form.phone" :placeholder="t('phone')" />
    </el-form-item>
    <el-form-item :label="t('email')" prop="email">
      <el-input v-model="form.email" :placeholder="t('email')" />
    </el-form-item>
    <el-form-item :label="t('userRole')" prop="authorityIds">
      <el-cascader
        v-model="form.authorityIds"
        style="width: 100%"
        :options="authOptions"
        :show-all-levels="false"
        collapse-tags
        :props="cascaderProps"
        :clearable="false"
        :placeholder="t('userRole')"
      />
    </el-form-item>
    <el-form-item :label="t('enable')" prop="enable">
      <el-switch
        v-model="form.enable"
        inline-prompt
        :active-value="1"
        :inactive-value="2"
      />
    </el-form-item>
    <el-form-item :label="t('avatar')">
      <SelectImage v-model="form.headerImg" />
    </el-form-item>
  </BaseFormDrawer>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDrawer from '@/components/base/BaseFormDrawer.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import type { Translator } from '@/types/consoleResource'
import type { UserAuthority, UserForm } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    authOptions?: UserAuthority[]
    dialogFlag?: 'add' | 'edit'
    form: UserForm
    modelValue?: boolean
    rules: FormRules<UserForm>
  }>(),
  {
    authOptions: () => [],
    dialogFlag: 'add',
    modelValue: false
  }
)

const emit = defineEmits<{
  close: []
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const cascaderProps = {
  multiple: true,
  checkStrictly: true,
  label: 'authorityName',
  value: 'authorityId',
  disabled: 'disabled',
  emitPath: false
}
</script>

<template>
  <el-drawer
    v-model="visibleModel"
    :size="480"
    :before-close="handleBeforeClose"
    :show-close="false"
    :close-on-press-escape="false"
    :close-on-click-modal="false"
  >
    <template #header>
      <div class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ dialogFlag === 'add' ? t('addUser') : t('editUser') }}</span>
        <div class="flex gap-3">
          <el-button @click="requestClose">{{ t('cancel') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ t('confirm') }}</el-button>
        </div>
      </div>
    </template>

    <el-form
      v-if="visibleModel"
      ref="userFormRef"
      :model="form"
      :rules="rules"
      label-width="auto"
    >
      <el-form-item v-if="dialogFlag === 'add'" :label="t('username')" prop="userName">
        <el-input v-model="form.userName" :placeholder="t('username')" />
      </el-form-item>
      <el-form-item v-if="dialogFlag === 'add'" :label="t('password')" prop="password">
        <el-input v-model="form.password" :placeholder="t('password')" show-password />
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
    </el-form>
  </el-drawer>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import SelectImage from '@/components/selectImage/selectImage.vue'

const props = defineProps({
  authOptions: {
    type: Array,
    default: () => []
  },
  dialogFlag: {
    type: String,
    default: 'add'
  },
  form: {
    type: Object,
    required: true
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

const emit = defineEmits(['close', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)
const userFormRef = ref(null)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const cascaderProps = {
  multiple: true,
  checkStrictly: true,
  label: 'authorityName',
  value: 'authorityId',
  disabled: 'disabled',
  emitPath: false
}

const requestClose = () => {
  emit('close')
}

const handleBeforeClose = (done) => {
  emit('close')
  done()
}

const submitForm = async () => {
  if (!userFormRef.value) {
    return
  }

  try {
    await userFormRef.value.validate()
    emit('submit')
  } catch {
    // Validation feedback is already shown by Element Plus.
  }
}
</script>

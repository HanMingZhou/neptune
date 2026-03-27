<template>
  <el-drawer
    v-model="visibleModel"
    :before-close="handleBeforeClose"
    :show-close="false"
    :size="480"
  >
    <template #header>
      <div class="flex justify-between items-center w-full">
        <span class="text-lg font-bold">{{ dialogTitle }}</span>
        <div class="flex gap-3">
          <el-button @click="requestClose">{{ t('cancel') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="submitForm">{{ t('confirm') }}</el-button>
        </div>
      </div>
    </template>

    <el-form
      v-if="visibleModel"
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="auto"
    >
      <el-form-item :label="t('parentMenu')" prop="parentId">
        <el-cascader
          v-model="form.parentId"
          style="width: 100%"
          :disabled="dialogType === 'add'"
          :options="authorityOptions"
          :props="cascaderProps"
          :show-all-levels="false"
          filterable
          :placeholder="t('selectParentRole')"
        />
      </el-form-item>
      <el-form-item :label="t('roleId')" prop="authorityId">
        <el-input
          v-model="form.authorityId"
          :disabled="dialogType === 'edit'"
          autocomplete="off"
          maxlength="15"
          :placeholder="t('inputRoleId')"
        />
      </el-form-item>
      <el-form-item :label="t('roleName')" prop="authorityName">
        <el-input v-model="form.authorityName" autocomplete="off" :placeholder="t('inputRoleName')" />
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script setup>
import { computed, inject, nextTick, ref, watch } from 'vue'

const props = defineProps({
  authorityOptions: {
    type: Array,
    required: true
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  dialogType: {
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
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'submit', 'update:modelValue'])
const formRef = ref(null)
const t = inject('t', (key) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const cascaderProps = {
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
  if (!formRef.value) {
    return
  }

  try {
    await formRef.value.validate()
    emit('submit')
  } catch {
    // Element Plus will show validation feedback.
  }
}

watch(
  () => props.modelValue,
  async (value) => {
    if (!value) {
      return
    }

    await nextTick()
    formRef.value?.clearValidate()
  }
)
</script>

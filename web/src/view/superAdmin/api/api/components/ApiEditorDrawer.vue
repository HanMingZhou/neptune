<template>
  <el-drawer
    v-model="visibleModel"
    :size="500"
    :before-close="requestClose"
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

    <el-form ref="apiFormRef" :model="form" :rules="rules" label-width="auto">
      <el-form-item :label="t('path')" prop="path">
        <el-input v-model="form.path" :placeholder="t('apiPath')" />
      </el-form-item>
      <el-form-item :label="t('method')" prop="method">
        <el-select v-model="form.method" :placeholder="t('select')" style="width: 100%">
          <el-option
            v-for="item in methodOptions"
            :key="item.value"
            :label="`${item.label}(${item.value})`"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('apiGroup')" prop="apiGroup">
        <el-select v-model="form.apiGroup" :placeholder="t('apiGroup')" allow-create filterable>
          <el-option v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('apiDesc')" prop="description">
        <el-input v-model="form.description" :placeholder="t('apiDesc')" />
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

const props = defineProps({
  apiGroupOptions: {
    type: Array,
    default: () => []
  },
  dialogTitle: {
    type: String,
    default: ''
  },
  form: {
    type: Object,
    required: true
  },
  methodOptions: {
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

const emit = defineEmits(['close', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)
const apiFormRef = ref(null)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const requestClose = () => {
  emit('close')
}

const submitForm = async () => {
  if (!apiFormRef.value) {
    return
  }

  await apiFormRef.value.validate()
  emit('submit')
}
</script>

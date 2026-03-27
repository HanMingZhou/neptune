<template>
  <el-dialog
    v-model="visibleModel"
    :before-close="handleBeforeClose"
    :title="t('newSshKey')"
    width="600px"
  >
    <el-form
      v-if="visibleModel"
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item :label="t('sshKeyName')" prop="name">
        <el-input v-model="form.name" :placeholder="t('sshKeyNamePlaceholder')" maxlength="50" />
        <div class="text-xs text-slate-400 mt-1">{{ t('sshKeyNamePlaceholder') }}</div>
      </el-form-item>

      <el-form-item :label="t('publicKeyContent')" prop="publicKey">
        <el-input
          v-model="form.publicKey"
          type="textarea"
          :rows="6"
          :placeholder="t('publicKeyPlaceholder')"
          class="public-key-textarea"
        />
        <div class="text-xs text-slate-400 mt-1">{{ t('publicKeyPlaceholder') }}</div>
      </el-form-item>
    </el-form>

    <div class="bg-slate-50 dark:bg-zinc-800 rounded-lg p-4 mt-4 text-xs">
      <h4 class="font-bold text-sm mb-2 flex items-center gap-2">
        <span class="material-icons text-[16px]">tips_and_updates</span>
        {{ t('sshKeyHintTitle') }}
      </h4>
      <ol class="text-slate-600 dark:text-slate-400 list-decimal list-inside space-y-1">
        <li>{{ t('sshKeyHintStep1') }}</li>
        <li>{{ t('sshKeyHintStep2') }}</li>
        <li>{{ t('sshKeyHintStep3') }}</li>
      </ol>
    </div>

    <template #footer>
      <el-button @click="requestClose">{{ t('cancel') }}</el-button>
      <el-button type="primary" :loading="loading" @click="submitForm">{{ t('add') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, nextTick, ref, watch } from 'vue'

const props = defineProps({
  form: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
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
const formRef = ref(null)
const t = inject('t', (key) => key)

const visibleModel = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

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

<style scoped>
:deep(.public-key-textarea .el-textarea__inner) {
  border: 1px solid var(--el-border-color);
  background-color: transparent;
}

:deep(.public-key-textarea .el-textarea__inner:focus) {
  border-color: var(--el-color-primary);
}
</style>

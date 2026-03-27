<template>
  <el-dialog
    v-model="visibleModel"
    :title="t('resetPassword')"
    width="500px"
    :before-close="handleBeforeClose"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form :model="form" label-width="auto">
      <el-form-item :label="t('username')">
        <el-input v-model="form.userName" disabled />
      </el-form-item>
      <el-form-item :label="t('nickname')">
        <el-input v-model="form.nickName" disabled />
      </el-form-item>
      <el-form-item :label="t('password')">
        <div class="flex w-full gap-3">
          <el-input v-model="form.password" class="flex-1" :placeholder="t('password')" show-password />
          <el-button type="primary" @click="$emit('generate-password')">
            {{ t('refresh') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button type="primary" @click="$emit('submit')">{{ t('confirm') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  form: {
    type: Object,
    required: true
  },
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'generate-password', 'submit', 'update:modelValue'])
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
</script>

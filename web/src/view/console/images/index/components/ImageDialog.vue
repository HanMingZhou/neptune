<template>
  <el-dialog v-model="dialogVisible" :title="isEdit ? t('imageEdit') : t('imageAdd')" width="600px" align-center @closed="handleClosed">
    <el-form ref="formRef" :model="form" label-position="top" :rules="rules" class="space-y-1">
      <div class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">{{ t('basicInfo') }}</div>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item :label="t('name')" prop="name">
            <el-input v-model="form.name" :placeholder="t('inputName')" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('imageArea')">
            <el-input v-model="form.area" :placeholder="t('inputArea')" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item :label="t('imageType')">
            <el-select v-model="form.type" class="!w-full">
              <el-option :label="t('systemImage')" :value="1" />
              <el-option :label="t('customImage')" :value="2" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('imageUsageType')" prop="usageType">
            <el-select v-model="form.usageType" class="!w-full">
              <el-option :label="t('usageNotebook')" :value="1" />
              <el-option :label="t('usageTrain')" :value="2" />
              <el-option :label="t('usageInfer')" :value="3" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <div class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark">
        {{ t('resourceConfig') }}
      </div>
      <el-form-item :label="t('imageAddr')">
        <el-input v-model="form.imageAddr" placeholder="registry.example.com/image:tag" />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item :label="t('imagePath')">
            <el-input v-model="form.imagePath" placeholder="/path/to/image" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('imageSize')">
            <el-input v-model="form.size" placeholder="e.g. 2.5GB" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? t('save') : t('create') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

const props = defineProps({
  form: {
    type: Object,
    required: true
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  rules: {
    type: Object,
    default: () => ({})
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['closed', 'submit', 'update:modelValue'])
const t = inject('t', (key) => key)
const formRef = ref(null)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('submit')
  } catch (error) {
    // validation error
  }
}

const handleClosed = () => {
  formRef.value?.clearValidate()
  emit('closed')
}
</script>

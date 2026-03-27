<template>
  <el-dialog
    v-model="visibleModel"
    :title="dialogTitle"
    width="680px"
    align-center
    :before-close="handleBeforeClose"
  >
    <el-form
      v-if="visibleModel"
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="space-y-1"
    >
      <div class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">{{ t('basicInfo') }}</div>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item :label="t('clusterName')" prop="name">
            <el-input v-model="form.name" :placeholder="t('inputName')" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('area')" prop="area">
            <el-input v-model="form.area" :placeholder="t('inputArea')" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item :label="t('desc')">
        <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('inputDesc')" />
      </el-form-item>

      <div class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark">
        {{ t('resourceConfig') }}
      </div>
      <el-form-item :label="t('clusterApiServer')">
        <el-input v-model="form.apiServer" placeholder="https://x.x.x.x:6443" />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item :label="t('clusterHarbor')">
            <el-input v-model="form.harborAddr" placeholder="harbor.example.com" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('storageClass')">
            <el-input v-model="form.storageClass" :placeholder="t('inputStorageClass')" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item :label="t('clusterKubeconfig')">
        <el-input
          v-model="form.kubeconfig"
          type="textarea"
          :rows="6"
          placeholder="kubeconfig YAML"
          class="font-mono"
        />
      </el-form-item>

      <div class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2 pt-2 border-t border-border-light dark:border-border-dark">
        {{ t('status') }}
      </div>
      <el-form-item>
        <el-switch
          v-model="form.status"
          :active-value="1"
          :inactive-value="0"
          :active-text="t('enable')"
          :inactive-text="t('disable')"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <el-button @click="requestClose">{{ t('cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">
          {{ isEdit ? t('save') : t('create') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'

const props = defineProps({
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
    // Element Plus will display validation feedback.
  }
}
</script>

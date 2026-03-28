<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-4 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('otherConfig') }}
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-slate-500 mb-2">{{ t('instanceName') }}</label>
          <input
            :value="instanceName"
            :placeholder="t('enterInstanceName')"
            class="create-form-input"
            maxlength="50"
            type="text"
            @input="$emit('update:instanceName', $event.target.value)"
          />
        </div>
        <div>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              :checked="enableTensorboard"
              class="rounded border-slate-300"
              type="checkbox"
              @change="$emit('update:enableTensorboard', $event.target.checked)"
            />
            <span class="text-sm">{{ t('enableTensorBoard') }}</span>
          </label>
          <div v-if="enableTensorboard" class="mt-2">
            <input
              :value="tensorboardLogPath"
              :placeholder="t('enterLogPath')"
              class="create-form-input"
              type="text"
              @input="$emit('update:tensorboardLogPath', $event.target.value)"
            />
          </div>
        </div>
      </div>

      <div class="space-y-4">
        <div>
          <label class="block text-sm text-slate-500 mb-2">{{ t('sshKeyLogin') }}</label>
          <el-select
            :model-value="selectedSshKey"
            :placeholder="t('selectSshKey')"
            clearable
            class="w-full"
            @update:model-value="$emit('update:selectedSshKey', $event)"
          >
            <el-option
              v-for="key in sshKeys"
              :key="key.id"
              :label="key.name"
              :value="key.id"
            />
          </el-select>
        </div>
        <div>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              :checked="enableSshPassword"
              class="rounded border-slate-300"
              type="checkbox"
              @change="$emit('update:enableSshPassword', $event.target.checked)"
            />
            <span class="text-sm">{{ t('enableSshPassword') }}</span>
          </label>
          <p v-if="enableSshPassword" class="mt-2 text-xs text-slate-400">✓ {{ t('sshPasswordHint') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  enableSshPassword: {
    type: Boolean,
    default: false
  },
  enableTensorboard: {
    type: Boolean,
    default: false
  },
  instanceName: {
    type: String,
    default: ''
  },
  selectedSshKey: {
    type: [Number, String],
    default: null
  },
  sshKeys: {
    type: Array,
    default: () => []
  },
  tensorboardLogPath: {
    type: String,
    default: ''
  }
})

defineEmits([
  'update:enableSshPassword',
  'update:enableTensorboard',
  'update:instanceName',
  'update:selectedSshKey',
  'update:tensorboardLogPath'
])

const t = inject('t', (key) => key)
</script>

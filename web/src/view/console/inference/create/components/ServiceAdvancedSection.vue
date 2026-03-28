<template>
  <div class="space-y-6">
    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-6 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.serviceConfig') }}
      </h3>
      <el-form :model="form" label-position="top" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <el-form-item :label="t('inference.servicePort')">
          <el-input-number
            :model-value="form.servicePort"
            :max="65535"
            :min="1024"
            class="w-full"
            @update:model-value="$emit('update:field', { key: 'servicePort', value: $event })"
          />
        </el-form-item>
        <el-form-item>
          <template #label>
            <span>{{ t('inference.authType') }}</span>
          </template>
          <div class="flex gap-3">
            <button
              v-for="item in authTypes"
              :key="item.value"
              :class="[
                'px-5 py-2 rounded-lg text-sm font-medium border transition-all',
                form.authType === item.value
                  ? 'bg-primary text-white border-primary shadow-lg shadow-primary/20'
                  : 'bg-white dark:bg-zinc-800 border-border-light dark:border-border-dark hover:border-primary hover:text-primary'
              ]"
              type="button"
              @click="$emit('update:field', { key: 'authType', value: item.value })"
            >
              {{ item.label }}
            </button>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-6 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.advancedConfig') }}
      </h3>
      <div class="space-y-6">
        <el-form :model="form" label-position="top">
          <el-form-item :label="t('inference.customCommand')" required>
            <el-input
              :model-value="form.command"
              :rows="3"
              :placeholder="'python3 -m vllm.entrypoints.openai.api_server --model /model/Qwen2.5-0.5B-Instruct --port 8000 --dtype float32'"
              class="w-full font-mono"
              type="textarea"
              @update:model-value="$emit('update:field', { key: 'command', value: $event })"
            />
            <div class="text-xs text-slate-400 mt-1">完整的启动命令，包含所有参数。挂载路径与命令中的 --model 路径需对应。</div>
          </el-form-item>
          <el-form-item :label="t('inference.customArgs')">
            <el-input
              :model-value="form.args"
              :placeholder="t('inference.customArgsHint')"
              :rows="4"
              class="w-full font-mono"
              type="textarea"
              @update:model-value="$emit('update:field', { key: 'args', value: $event })"
            />
            <div class="text-xs text-slate-400 mt-1">{{ t('inference.customArgsDesc') }}</div>
          </el-form-item>
        </el-form>

        <div>
          <label class="block text-sm text-slate-500 mb-2">{{ t('inference.dataMount') }}</label>
          <div class="bg-slate-50 dark:bg-zinc-800/50 rounded-lg p-4 border border-border-light dark:border-border-dark">
            <div v-if="form.mounts.length > 0" class="mb-3">
              <div class="grid grid-cols-12 gap-4 text-xs text-slate-400 font-medium mb-2 px-1">
                <div class="col-span-3">{{ t('pvc') }}</div>
                <div class="col-span-3">{{ t('inference.mountPath') }}</div>
                <div class="col-span-3">{{ t('inference.subPath') }}</div>
                <div class="col-span-2 text-center">{{ t('readOnly') }}</div>
                <div class="col-span-1 text-center">{{ t('actions') }}</div>
              </div>
              <div v-for="(mount, index) in form.mounts" :key="index" class="grid grid-cols-12 gap-4 items-center mb-2">
                <div class="col-span-3">
                  <el-select
                    :model-value="mount.pvcId"
                    :placeholder="t('selectPvc')"
                    class="w-full"
                    size="default"
                    @update:model-value="$emit('update:mount', { index, key: 'pvcId', value: $event })"
                  >
                    <el-option v-for="pvc in pvcs" :key="pvc.id" :label="pvc.name" :value="pvc.id" />
                  </el-select>
                </div>
                <div class="col-span-3">
                  <input
                    :value="mount.mountPath"
                    class="create-form-input"
                    placeholder="/data"
                    type="text"
                    @input="$emit('update:mount', { index, key: 'mountPath', value: $event.target.value })"
                  />
                </div>
                <div class="col-span-3">
                  <input
                    :value="mount.subPath"
                    class="create-form-input"
                    placeholder="sub-dir"
                    type="text"
                    @input="$emit('update:mount', { index, key: 'subPath', value: $event.target.value })"
                  />
                </div>
                <div class="col-span-2 text-center">
                  <el-checkbox
                    :model-value="mount.readOnly"
                    @update:model-value="$emit('update:mount', { index, key: 'readOnly', value: $event })"
                  />
                </div>
                <div class="col-span-1 text-center">
                  <button class="text-red-500 hover:text-red-600" @click="$emit('remove-mount', index)">
                    <span class="material-icons text-lg">delete</span>
                  </button>
                </div>
              </div>
            </div>
            <button
              class="w-full py-2 border border-dashed border-border-light dark:border-border-dark rounded-lg text-sm text-slate-500 hover:border-primary hover:text-primary transition-colors flex items-center justify-center gap-1"
              @click="$emit('add-mount')"
            >
              <span class="material-icons text-base">add</span>
              {{ t('addMount') }}
            </button>
          </div>
        </div>

        <div>
          <label class="block text-sm text-slate-500 mb-2">{{ t('envVars') }}</label>
          <div class="bg-slate-50 dark:bg-zinc-800/50 rounded-lg p-4 border border-border-light dark:border-border-dark">
            <div v-if="form.envs.length > 0" class="mb-3">
              <div class="grid grid-cols-12 gap-4 text-xs text-slate-400 font-medium mb-2 px-1">
                <div class="col-span-4">{{ t('variableName') }}</div>
                <div class="col-span-7">{{ t('variableValue') }}</div>
                <div class="col-span-1 text-center">{{ t('actions') }}</div>
              </div>
              <div v-for="(env, index) in form.envs" :key="index" class="grid grid-cols-12 gap-4 items-center mb-2">
                <div class="col-span-4">
                  <input
                    :value="env.name"
                    :placeholder="t('enterVariableName')"
                    class="create-form-input font-mono"
                    type="text"
                    @input="$emit('update:env', { index, key: 'name', value: $event.target.value })"
                  />
                </div>
                <div class="col-span-7">
                  <input
                    :value="env.value"
                    :placeholder="t('enterVariableValue')"
                    class="create-form-input"
                    type="text"
                    @input="$emit('update:env', { index, key: 'value', value: $event.target.value })"
                  />
                </div>
                <div class="col-span-1 text-center">
                  <button class="text-red-500 hover:text-red-600" @click="$emit('remove-env', index)">
                    <span class="material-icons text-lg">delete</span>
                  </button>
                </div>
              </div>
            </div>
            <button
              class="w-full py-2 border border-dashed border-border-light dark:border-border-dark rounded-lg text-sm text-slate-500 hover:border-primary hover:text-primary transition-colors flex items-center justify-center gap-1"
              @click="$emit('add-env')"
            >
              <span class="material-icons text-base">add</span>
              {{ t('addEnv') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  authTypes: {
    type: Array,
    default: () => []
  },
  form: {
    type: Object,
    required: true
  },
  pvcs: {
    type: Array,
    default: () => []
  }
})

defineEmits([
  'add-env',
  'add-mount',
  'remove-env',
  'remove-mount',
  'update:env',
  'update:field',
  'update:mount'
])

const t = inject('t', (key) => key)
</script>

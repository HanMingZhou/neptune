<template>
  <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
    <h3 class="text-base font-bold mb-6 flex items-center gap-2">
      <span class="w-1 h-4 bg-primary rounded"></span>
      {{ t('jobConfig') }}
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-slate-500 mb-2">
            {{ t('jobName') }} <span class="text-red-500">*</span>
          </label>
          <input
            :value="form.name"
            :placeholder="t('enterJobName')"
            class="create-form-input"
            type="text"
            @input="$emit('update:field', { key: 'name', value: $event.target.value })"
          />
        </div>
        <div>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              :checked="form.enableTensorboard"
              class="rounded border-slate-300"
              type="checkbox"
              @change="$emit('update:field', { key: 'enableTensorboard', value: $event.target.checked })"
            />
            <span class="text-sm">{{ t('enableTensorBoard') }}</span>
          </label>
          <div v-if="form.enableTensorboard" class="mt-2">
            <input
              :value="form.tensorboardLogPath"
              :placeholder="t('enterLogPath')"
              class="create-form-input"
              type="text"
              @input="$emit('update:field', { key: 'tensorboardLogPath', value: $event.target.value })"
            />
          </div>
        </div>
      </div>

      <div class="space-y-4">
        <div>
          <label class="block text-sm text-slate-500 mb-2">
            {{ t('trainingFramework') }} <span class="text-red-500">*</span>
          </label>
          <el-select
            :model-value="form.frameworkType"
            class="w-full"
            @update:model-value="$emit('update:field', { key: 'frameworkType', value: $event })"
          >
            <el-option
              v-for="framework in frameworkTypes"
              :key="framework.value"
              :label="framework.label"
              :value="framework.value"
            />
          </el-select>
        </div>

        <WorkerCountSection
          :available-capacity="availableCapacity"
          :embedded="true"
          :framework-type="form.frameworkType"
          :max-worker-count="maxWorkerCount"
          :schedule-strategy="form.scheduleStrategy"
          :selected-product="selectedProduct"
          :show-worker-count="showWorkerCount"
          :worker-count="form.workerCount"
          @decrease-worker="$emit('decrease-worker')"
          @increase-worker="$emit('increase-worker')"
          @update:schedule-strategy="$emit('update:schedule-strategy', $event)"
          @update:worker-count="$emit('update:workerCount', $event)"
        />
      </div>
    </div>

    <div class="mb-6">
      <label class="block text-sm text-slate-500 mb-2">
        {{ t('startupCommand') }} <span class="text-red-500">*</span>
      </label>
      <textarea
        :value="form.startupCommand"
        :placeholder="t('enterStartupCommand')"
        class="create-form-textarea font-mono"
        rows="3"
        @input="$emit('update:field', { key: 'startupCommand', value: $event.target.value })"
      ></textarea>
      <div v-if="form.frameworkType === 'MPI'" class="mt-2 flex items-center justify-between">
        <div class="flex items-center gap-2 text-sm text-amber-600 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg">
          <span class="material-icons text-base">info</span>
          {{ t('mpiCommandHint') }}
        </div>
        <button class="text-sm text-primary hover:underline" @click="$emit('insert-mpi-example')">{{ t('insertExample') }}</button>
      </div>
      <div v-if="form.frameworkType === 'PYTORCH_DDP'" class="mt-2 flex items-center justify-between">
        <div class="flex items-center gap-2 text-sm text-amber-600 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg">
          <span class="material-icons text-base">info</span>
          {{ t('pytorchDdpHint') }}
        </div>
        <button class="text-sm text-primary hover:underline" @click="$emit('insert-pytorch-example')">{{ t('insertExample') }}</button>
      </div>
    </div>

    <div class="mb-6">
      <label class="block text-sm text-slate-500 mb-2">{{ t('dataMount') }}</label>
      <div class="bg-slate-50 dark:bg-zinc-800/50 rounded-lg p-4 border border-border-light dark:border-border-dark">
        <div v-if="form.mounts.length > 0" class="mb-3">
          <div class="grid grid-cols-12 gap-4 text-xs text-slate-400 font-medium mb-2 px-1">
            <div class="col-span-4">{{ t('pvc') }}</div>
            <div class="col-span-5">{{ t('mountPath') }}</div>
            <div class="col-span-2 text-center">{{ t('readOnly') }}</div>
            <div class="col-span-1 text-center">{{ t('actions') }}</div>
          </div>
          <div v-for="(mount, index) in form.mounts" :key="index" class="grid grid-cols-12 gap-4 items-center mb-2">
            <div class="col-span-4">
              <el-select
                :model-value="mount.pvcId"
                :placeholder="t('selectPvc')"
                class="w-full"
                @update:model-value="$emit('update:mount', { index, key: 'pvcId', value: $event })"
              >
                <el-option v-for="pvc in pvcs" :key="pvc.id" :label="pvc.name" :value="pvc.id" />
              </el-select>
            </div>
            <div class="col-span-5">
              <input
                :value="mount.mountPath"
                :placeholder="t('enterMountPath')"
                class="create-form-input"
                type="text"
                @input="$emit('update:mount', { index, key: 'mountPath', value: $event.target.value })"
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
</template>

<script setup>
import { inject } from 'vue'
import WorkerCountSection from './WorkerCountSection.vue'

defineProps({
  availableCapacity: {
    type: Number,
    default: 0
  },
  form: {
    type: Object,
    required: true
  },
  frameworkTypes: {
    type: Array,
    default: () => []
  },
  maxWorkerCount: {
    type: Number,
    default: 2
  },
  pvcs: {
    type: Array,
    default: () => []
  },
  selectedProduct: {
    type: Object,
    default: null
  },
  showWorkerCount: {
    type: Boolean,
    default: false
  }
})

defineEmits([
  'add-env',
  'add-mount',
  'decrease-worker',
  'insert-mpi-example',
  'insert-pytorch-example',
  'increase-worker',
  'remove-env',
  'remove-mount',
  'update:env',
  'update:field',
  'update:mount',
  'update:schedule-strategy',
  'update:workerCount'
])

const t = inject('t', (key) => key)
</script>

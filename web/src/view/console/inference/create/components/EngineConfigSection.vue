<template>
  <div class="space-y-6">
    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-6 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('engineConfig') }}
      </h3>
      <el-form :model="form" label-position="top" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <el-form-item :label="t('name')" :error="fieldErrors.displayName || ''" required>
          <el-input
            :model-value="form.displayName"
            maxlength="63"
            :placeholder="t('inference.enterName')"
            class="w-full"
            @update:model-value="$emit('update:field', { key: 'displayName', value: $event })"
            @blur="$emit('validate:display-name')"
          />
          <div class="text-xs text-slate-400 mt-1">{{ t('resourceNameHint', { max: 63 }) }}</div>
        </el-form-item>
        <el-form-item :label="t('inference.framework')" :required="frameworkRequired">
          <el-select
            :model-value="form.framework"
            :placeholder="t('inference.selectFramework')"
            clearable
            class="w-full"
            @update:model-value="$emit('update:field', { key: 'framework', value: $event })"
          >
            <el-option
              v-for="item in frameworks"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
              <div class="flex items-center gap-3">
                <span class="material-icons text-base opacity-70">{{ item.icon }}</span>
                <span class="font-medium">{{ item.label }}</span>
              </div>
            </el-option>
          </el-select>
          <div v-if="!frameworkRequired" class="text-xs text-slate-400 mt-1">{{ t('inference.frameworkOptionalHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('deployMode')" required>
          <el-select
            :model-value="form.deployType"
            :placeholder="t('inference.selectDeployMode')"
            class="w-full"
            @update:model-value="$emit('deploy-type-change', $event)"
          >
            <el-option v-for="item in deployTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="form.deployType === 'DISTRIBUTED'" class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-6 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.distributedConfig') }}
      </h3>
      <div class="flex items-center gap-2 text-sm text-amber-600 bg-amber-50 dark:bg-amber-900/20 px-3 py-2 rounded-lg mb-6">
        <span class="material-icons text-base">info</span>
        {{ t('inference.distributedHint') }}
      </div>
      <el-form :model="form" label-position="top">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <el-form-item :label="t('inference.workerCount')" required>
            <div class="create-form-stepper">
              <button
                class="create-form-stepper__button"
                type="button"
                @click="$emit('update:field', { key: 'workerCount', value: Math.max(2, form.workerCount - 1) })"
              >
                -
              </button>
              <input
                :value="form.workerCount"
                class="create-form-stepper__input"
                min="2"
                type="number"
                @input="$emit('update:field', { key: 'workerCount', value: Number($event.target.value) || 2 })"
              />
              <button
                class="w-9 h-9 border border-border-light dark:border-border-dark rounded-r-lg hover:bg-slate-50 dark:hover:bg-zinc-800"
                type="button"
                @click="$emit('update:field', { key: 'workerCount', value: form.workerCount + 1 })"
              >
                +
              </button>
            </div>
          </el-form-item>
          <el-form-item label="调度策略" required>
            <el-select
              :model-value="form.scheduleStrategy"
              class="w-full"
              @update:model-value="$emit('update:field', { key: 'scheduleStrategy', value: $event || 'BALANCED' })"
            >
              <el-option label="智能均衡（默认）" value="BALANCED" />
              <el-option label="严格分布式（节点不足则失败）" value="STRICT" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('inference.autoRestart')">
            <div class="flex items-center gap-4">
              <el-switch
                :model-value="form.autoRestart"
                @update:model-value="$emit('update:field', { key: 'autoRestart', value: $event })"
              />
              <template v-if="form.autoRestart">
                <span class="text-xs text-slate-500">{{ t('inference.maxRestarts') }}</span>
                <div class="create-form-stepper">
                  <button
                    class="create-form-stepper__button"
                    type="button"
                    @click="$emit('update:field', { key: 'maxRestarts', value: Math.max(1, form.maxRestarts - 1) })"
                  >
                    -
                  </button>
                  <input
                    :value="form.maxRestarts"
                    class="create-form-stepper__input"
                    max="10"
                    min="1"
                    type="number"
                    @input="$emit('update:field', { key: 'maxRestarts', value: Number($event.target.value) || 1 })"
                  />
                  <button
                    class="w-8 h-8 border border-border-light dark:border-border-dark rounded-r-lg hover:bg-slate-50 dark:hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                    type="button"
                    @click="$emit('update:field', { key: 'maxRestarts', value: Math.min(10, form.maxRestarts + 1) })"
                  >
                    +
                  </button>
                </div>
              </template>
            </div>
          </el-form-item>
        </div>
      </el-form>
    </div>

    <div class="bg-surface-light dark:bg-surface-dark border border-border-light dark:border-border-dark rounded-xl p-6">
      <h3 class="text-base font-bold mb-6 flex items-center gap-2">
        <span class="w-1 h-4 bg-primary rounded"></span>
        {{ t('inference.modelConfig') }}
      </h3>
      <el-form :model="form" label-position="top" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <el-form-item :label="t('inference.modelVolume')" required>
          <el-select
            :model-value="form.modelPvcId"
            :placeholder="t('inference.selectModelPvc')"
            class="w-full"
            @update:model-value="$emit('update:field', { key: 'modelPvcId', value: $event })"
          >
            <el-option v-for="pvc in pvcs" :key="pvc.id" :label="pvc.name" :value="pvc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('inference.modelMountPath')" required>
          <el-input
            :model-value="form.modelMountPath"
            class="w-full"
            placeholder="/model"
            @update:model-value="$emit('update:field', { key: 'modelMountPath', value: $event })"
          />
          <div class="text-xs text-slate-400 mt-1">PVC 挂载到容器内的路径</div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

defineProps({
  deployTypes: {
    type: Array,
    default: () => []
  },
  form: {
    type: Object,
    required: true
  },
  fieldErrors: {
    type: Object,
    default: () => ({})
  },
  frameworkRequired: {
    type: Boolean,
    default: false
  },
  frameworks: {
    type: Array,
    default: () => []
  },
  pvcs: {
    type: Array,
    default: () => []
  }
})

defineEmits(['deploy-type-change', 'update:field', 'validate:display-name'])

const t = inject('t', (key) => key)
</script>

<template>
  <BaseFormDialog
    v-model="dialogVisible"
    :cancel-text="t('cancel')"
    class="image-dialog"
    form-class="image-dialog__form"
    :model="form"
    :rules="rules"
    :shell-size="'lg'"
    :submit-text="isEdit ? t('save') : t('create')"
    :submitting="submitting"
    :title="isEdit ? t('imageEdit') : t('imageAdd')"
    label-position="top"
    top="6vh"
    @closed="emit('closed')"
    @submit="emit('submit')"
  >
    <div
      class="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1.08fr)_minmax(0,0.92fr)]"
    >
      <section
        class="rounded-[22px] border border-slate-200/70 bg-[linear-gradient(180deg,rgba(255,255,255,0.98)_0%,rgba(248,250,252,0.96)_100%)] p-4 shadow-[0_18px_48px_-36px_rgba(15,23,42,0.45)] dark:border-zinc-700/70 dark:bg-[linear-gradient(180deg,rgba(39,39,42,0.96)_0%,rgba(24,24,27,0.96)_100%)]"
      >
        <div class="mb-4 flex items-start gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-sky-500/12 text-sky-600 ring-1 ring-sky-500/15 dark:bg-sky-500/15 dark:text-sky-300 dark:ring-sky-400/20"
          >
            <span class="material-icons text-[18px]">inventory_2</span>
          </div>
          <div class="min-w-0">
            <div
              class="text-sm font-semibold text-slate-900 dark:text-slate-100"
            >
              {{ t('basicInfo') }}
            </div>
            <div
              class="mt-2 flex flex-wrap gap-2 text-[11px] font-medium text-slate-500 dark:text-slate-400"
            >
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-slate-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('name') }}
              </span>
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-slate-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('imageType') }}
              </span>
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-slate-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('imageUsageType') }}
              </span>
            </div>
          </div>
        </div>

        <div class="image-dialog__fields grid grid-cols-1 gap-4 md:grid-cols-2">
          <el-form-item :label="t('name')" prop="name">
            <el-input v-model="form.name" :placeholder="t('inputName')" />
          </el-form-item>
          <el-form-item :label="t('cluster')" prop="clusterId">
            <el-select
              v-model="form.clusterId"
              :placeholder="t('selectCluster')"
            >
              <el-option
                v-for="item in clusterOptions"
                :key="item.id"
                :label="item.area ? `${item.area} · ${item.name}` : item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('imageType')">
            <el-select v-model="form.type">
              <el-option :label="t('systemImage')" :value="1" />
              <el-option :label="t('customImage')" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('imageUsageType')" prop="usageType">
            <el-select v-model="form.usageType">
              <el-option :label="t('usageNotebook')" :value="1" />
              <el-option :label="t('usageTrain')" :value="2" />
              <el-option :label="t('usageInfer')" :value="3" />
            </el-select>
          </el-form-item>
        </div>

        <div
          class="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2"
          v-if="selectedCluster"
        >
          <div
            class="rounded-2xl border border-slate-200/80 bg-white/80 px-4 py-3 text-sm dark:border-white/10 dark:bg-white/5"
          >
            <div
              class="text-[11px] font-semibold uppercase tracking-wider text-slate-400"
            >
              {{ t('imageArea') }}
            </div>
            <div class="mt-1 font-medium text-slate-700 dark:text-slate-200">
              {{ selectedCluster.area || '-' }}
            </div>
          </div>
          <div
            class="rounded-2xl border border-slate-200/80 bg-white/80 px-4 py-3 text-sm dark:border-white/10 dark:bg-white/5"
          >
            <div
              class="text-[11px] font-semibold uppercase tracking-wider text-slate-400"
            >
              {{ t('clusterHarbor') }}
            </div>
            <div
              class="mt-1 break-all font-medium text-slate-700 dark:text-slate-200"
            >
              {{ selectedCluster.harborAddr || '-' }}
            </div>
          </div>
        </div>
      </section>

      <section
        class="rounded-[22px] border border-emerald-200/70 bg-[linear-gradient(180deg,rgba(240,253,250,0.96)_0%,rgba(236,253,245,0.92)_100%)] p-4 shadow-[0_18px_48px_-36px_rgba(6,95,70,0.32)] dark:border-emerald-900/60 dark:bg-[linear-gradient(180deg,rgba(6,78,59,0.28)_0%,rgba(24,24,27,0.96)_100%)]"
      >
        <div class="mb-4 flex items-start gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-emerald-500/12 text-emerald-600 ring-1 ring-emerald-500/15 dark:bg-emerald-500/15 dark:text-emerald-300 dark:ring-emerald-400/20"
          >
            <span class="material-icons text-[18px]">dns</span>
          </div>
          <div class="min-w-0">
            <div
              class="text-sm font-semibold text-slate-900 dark:text-slate-100"
            >
              {{ t('resourceConfig') }}
            </div>
            <div
              class="mt-2 flex flex-wrap gap-2 text-[11px] font-medium text-slate-500 dark:text-slate-400"
            >
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-emerald-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('imageAddr') }}
              </span>
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-emerald-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('imagePath') }}
              </span>
              <span
                class="rounded-full bg-white/80 px-2.5 py-1 ring-1 ring-emerald-200/80 dark:bg-white/5 dark:ring-white/10"
              >
                {{ t('imageSize') }}
              </span>
            </div>
          </div>
        </div>

        <div class="image-dialog__fields grid grid-cols-1 gap-4 md:grid-cols-2">
          <el-form-item
            :label="t('imagePath')"
            prop="imagePath"
            class="md:col-span-2"
          >
            <el-input
              v-model="form.imagePath"
              placeholder="project/image:tag"
            />
          </el-form-item>
          <el-form-item :label="t('imageAddr')" class="md:col-span-2">
            <el-input
              v-model="form.imageAddr"
              :placeholder="
                generatedImageAddr || 'registry.example.com/project/image:tag'
              "
              @input="emit('image-addr-input', String($event || ''))"
            />
          </el-form-item>
          <el-form-item :label="t('imageSize')">
            <el-input v-model="form.size" placeholder="e.g. 2.5GB" />
          </el-form-item>
        </div>
      </section>
    </div>
  </BaseFormDialog>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import type { FormRules } from 'element-plus'
import BaseFormDialog from '@/components/base/BaseFormDialog.vue'
import type { Translator } from '@/types/consoleResource'
import type { ImageForm } from '@/types/image'
import type { CmsClusterOption } from '@/types/superAdmin'

const props = withDefaults(
  defineProps<{
    clusterOptions?: CmsClusterOption[]
    form: ImageForm
    generatedImageAddr?: string
    isEdit?: boolean
    modelValue?: boolean
    rules?: FormRules<ImageForm>
    submitting?: boolean
  }>(),
  {
    clusterOptions: () => [],
    generatedImageAddr: '',
    isEdit: false,
    modelValue: false,
    rules: () => ({}),
    submitting: false
  }
)

const emit = defineEmits<{
  closed: []
  'image-addr-input': [value: string]
  submit: []
  'update:modelValue': [value: boolean]
}>()
const t = inject<Translator>('t', (key: string) => key)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const selectedCluster = computed<CmsClusterOption | null>(
  () =>
    props.clusterOptions.find((item) => item.id === props.form.clusterId) ||
    null
)
</script>

<style scoped>
.image-dialog__form {
  padding-top: 0.25rem;
  padding-bottom: 0.125rem;
}

.image-dialog__fields :deep(.el-form-item) {
  margin-bottom: 0;
}

.image-dialog__fields :deep(.el-form-item__label) {
  padding-bottom: 0.45rem;
  font-size: 0.8125rem;
  font-weight: 600;
  line-height: 1.3;
}

.image-dialog__fields :deep(.el-input__wrapper),
.image-dialog__fields :deep(.el-select__wrapper) {
  min-height: 44px;
}
</style>

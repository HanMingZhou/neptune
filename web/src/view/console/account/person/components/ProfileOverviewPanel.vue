<template>
  <div class="space-y-4">
    <ProfileHero
      :edit-flag="editFlag"
      :nick-name="nickName"
      :user-info="userInfo"
      @change-email="emit('change-email')"
      @change-password="emit('change-password')"
      @change-phone="emit('change-phone')"
      @close-edit="emit('close-edit')"
      @confirm-edit="emit('confirm-edit')"
      @open-edit="emit('open-edit')"
      @update:nick-name="emit('update:nick-name', $event)"
    />

    <section class="console-detail-card rounded-xl p-5 md:p-6">
      <header class="mb-5 space-y-2">
        <div class="space-y-2">
          <p class="text-xs font-semibold uppercase tracking-[0.08em] text-slate-500">{{ t('accountOverview') }}</p>
          <h2 class="text-2xl font-semibold tracking-tight text-slate-950 dark:text-slate-100">
            {{ t('person') }}
          </h2>
          <p class="max-w-3xl text-sm leading-6 text-slate-500">
            {{ t('profileCenterDesc') }}
          </p>
        </div>
      </header>

      <div class="detail-info-grid">
        <article
          v-for="item in overviewItems"
          :key="item.key"
          class="detail-info-item gap-4"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex min-w-0 items-start gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg border border-border-light bg-surface-light text-slate-500 dark:border-border-dark dark:bg-surface-dark">
              <el-icon><component :is="item.icon" /></el-icon>
              </div>
              <div class="min-w-0">
                <p class="detail-info-label">{{ item.label }}</p>
                <p class="detail-info-value">{{ item.value }}</p>
                <p class="mt-1 text-sm leading-6 text-slate-500">{{ item.description }}</p>
              </div>
            </div>

            <button class="list-row-button list-row-button--neutral shrink-0" @click="item.action()">
              {{ t('edit') }}
            </button>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { Key, Message, Phone, User } from '@element-plus/icons-vue'

import ProfileHero from './ProfileHero.vue'

const props = defineProps({
  editFlag: {
    type: Boolean,
    default: false
  },
  nickName: {
    type: String,
    default: ''
  },
  userInfo: {
    type: Object,
    required: true
  }
})

const emit = defineEmits([
  'change-email',
  'change-password',
  'change-phone',
  'close-edit',
  'confirm-edit',
  'open-edit',
  'update:nick-name'
])

const t = inject('t', (key) => key)

const overviewItems = computed(() => [
  {
    key: 'nickname',
    icon: User,
    label: t('nickname'),
    value: props.userInfo.nickName || t('nicknameEmpty'),
    description: t('nicknameUpdateHint'),
    action: () => emit('open-edit')
  },
  {
    key: 'phone',
    icon: Phone,
    label: t('phoneLabel'),
    value: props.userInfo.phone || t('notSet'),
    description: t(props.userInfo.phone ? 'phoneReadyHint' : 'phonePendingHint'),
    action: () => emit('change-phone')
  },
  {
    key: 'email',
    icon: Message,
    label: t('emailLabel'),
    value: props.userInfo.email || t('notSet'),
    description: t(props.userInfo.email ? 'emailReadyHint' : 'emailPendingHint'),
    action: () => emit('change-email')
  },
  {
    key: 'password',
    icon: Key,
    label: t('changePasswordTitle'),
    value: t('isSet'),
    description: t('passwordReadyHint'),
    action: () => emit('change-password')
  }
])
</script>

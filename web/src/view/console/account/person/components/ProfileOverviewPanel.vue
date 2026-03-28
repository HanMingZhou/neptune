<template>
  <div class="profile-container">
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

    <section class="surface-card detail-surface">
      <header class="detail-header">
        <div class="space-y-3">
          <p class="section-kicker">{{ t('accountOverview') }}</p>
          <h2 class="text-3xl font-semibold tracking-tight text-slate-950">
            {{ t('person') }}
          </h2>
          <p class="max-w-3xl text-sm leading-7 text-slate-500">
            {{ t('profileCenterDesc') }}
          </p>
        </div>
      </header>

      <div class="detail-list">
        <article
          v-for="item in overviewItems"
          :key="item.key"
          class="detail-row"
        >
          <div class="detail-main">
            <div class="info-icon">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
            <div class="detail-copy">
              <p class="detail-label">{{ item.label }}</p>
              <p class="detail-value">{{ item.value }}</p>
              <p class="detail-note">{{ item.description }}</p>
            </div>
          </div>

          <el-button class="action-link" @click="item.action()">
            {{ t('edit') }}
          </el-button>
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

<style scoped>
.profile-container {
  padding: 16px;
  background: #f8fafc;
}

.surface-card {
  border: 1px solid #e2e8f0;
  border-radius: 24px;
  background: #ffffff;
  padding: 24px;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.05);
}

.detail-surface {
  display: flex;
  flex-direction: column;
  gap: 26px;
}

.detail-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.detail-list {
  border-top: 1px solid #e2e8f0;
}

.detail-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 22px 0;
  border-bottom: 1px solid #e2e8f0;
}

.detail-main {
  display: flex;
  min-width: 0;
  align-items: flex-start;
  gap: 16px;
}

.detail-copy {
  min-width: 0;
}

.detail-label {
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #64748b;
}

.detail-value {
  font-size: 19px;
  font-weight: 600;
  color: #0f172a;
  word-break: break-word;
}

.detail-note {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.7;
  color: #64748b;
}

.section-kicker {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #64748b;
}

.info-icon {
  display: inline-flex;
  width: 42px;
  height: 42px;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  background: #f8fafc;
  color: #475569;
}

.action-link {
  flex-shrink: 0;
  border-radius: 9999px;
  border-color: #dbe2ea;
  color: #0f172a;
}

@media (min-width: 1024px) {
  .profile-container {
    padding: 24px;
  }
}

@media (max-width: 767px) {
  .surface-card {
    padding: 18px;
    border-radius: 20px;
  }

  .detail-row {
    flex-direction: column;
    align-items: stretch;
  }

  .action-link {
    width: fit-content;
  }
}
</style>

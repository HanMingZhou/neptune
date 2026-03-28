<template>
  <section class="hero-card">
    <div class="space-y-6">
      <p class="hero-kicker">{{ t('person') }}</p>

      <div class="flex flex-col gap-6 md:flex-row md:items-start">
        <div class="avatar-shell">
          <el-avatar v-if="userInfo.headerImg" :src="userInfo.headerImg" :size="88" />
          <div v-else class="avatar-fallback">
            {{ userInitials }}
          </div>
        </div>

        <div class="min-w-0 flex-1 space-y-5">
          <div class="space-y-3">
            <div
              v-if="!editFlag"
              class="flex flex-col gap-3 md:flex-row md:items-center md:gap-4"
            >
              <h1 class="truncate text-3xl font-semibold tracking-tight text-slate-950">
                {{ displayName }}
              </h1>
              <el-button class="hero-action-btn" @click="$emit('open-edit')">
                <el-icon><Edit /></el-icon>
                {{ t('editNickname') }}
              </el-button>
            </div>

            <div v-else class="flex flex-col gap-3 sm:flex-row sm:items-center">
              <el-input v-model="nickNameModel" class="max-w-sm" />
              <div class="flex gap-3">
                <el-button type="primary" @click="$emit('confirm-edit')">
                  {{ t('confirm') }}
                </el-button>
                <el-button @click="$emit('close-edit')">
                  {{ t('cancel') }}
                </el-button>
              </div>
            </div>

            <p class="max-w-2xl text-sm leading-7 text-slate-600">
              {{ t('profileCenterDesc') }}
            </p>
          </div>

          <div class="flex flex-wrap gap-3">
            <div class="hero-meta-chip">
              <el-icon><Message /></el-icon>
              <span>{{ userInfo.email || t('notSet') }}</span>
            </div>
            <div class="hero-meta-chip">
              <el-icon><Phone /></el-icon>
              <span>{{ userInfo.phone || t('notSet') }}</span>
            </div>
            <div class="hero-meta-chip">
              <el-icon><Postcard /></el-icon>
              <span>{{ shortAccountId }}</span>
            </div>
          </div>

          <div class="flex flex-wrap gap-3">
            <el-button type="primary" @click="$emit('change-password')">
              <el-icon><Lock /></el-icon>
              {{ t('changePasswordTitle') }}
            </el-button>
            <el-button @click="$emit('change-email')">
              <el-icon><Message /></el-icon>
              {{ t('emailLabel') }}
            </el-button>
            <el-button @click="$emit('change-phone')">
              <el-icon><Phone /></el-icon>
              {{ t('phoneLabel') }}
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, inject } from 'vue'
import { Edit, Lock, Message, Phone, Postcard } from '@element-plus/icons-vue'

const props = defineProps({
  userInfo: {
    type: Object,
    required: true
  },
  nickName: {
    type: String,
    default: ''
  },
  editFlag: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'change-email',
  'change-password',
  'change-phone',
  'open-edit',
  'close-edit',
  'confirm-edit',
  'update:nick-name'
])

const t = inject('t', (key) => key)

const displayName = computed(() => props.userInfo.nickName || t('nicknameEmpty'))

const shortAccountId = computed(() => {
  const accountId = props.userInfo.uuid || props.userInfo.id || props.userInfo.userName
  if (!accountId) {
    return t('notSet')
  }

  const value = String(accountId)
  return value.length > 12 ? `${value.slice(0, 6)}...${value.slice(-4)}` : value
})

const userInitials = computed(() => {
  const rawName = props.userInfo.nickName || props.userInfo.userName || props.userInfo.username || 'U'
  const text = String(rawName).trim()

  if (!text) {
    return 'U'
  }

  return text.slice(0, 2).toUpperCase()
})

const nickNameModel = computed({
  get: () => props.nickName,
  set: (value) => {
    emit('update:nick-name', value)
  }
})
</script>

<style scoped>
.hero-card {
  border: 1px solid #e2e8f0;
  border-radius: 24px;
  padding: 28px;
  background: #ffffff;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.05);
}

.hero-kicker {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.08em;
  color: #64748b;
  text-transform: uppercase;
}

.avatar-shell {
  display: flex;
  width: 108px;
  height: 108px;
  align-items: center;
  justify-content: center;
  border-radius: 24px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
}

.avatar-fallback {
  display: flex;
  width: 88px;
  height: 88px;
  align-items: center;
  justify-content: center;
  border-radius: 20px;
  background: #e2e8f0;
  font-size: 24px;
  font-weight: 700;
  color: #334155;
}

.hero-meta-chip {
  display: inline-flex;
  max-width: 100%;
  align-items: center;
  gap: 10px;
  border-radius: 9999px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  padding: 10px 14px;
  font-size: 13px;
  color: #334155;
}

.hero-meta-chip :deep(.el-icon) {
  color: #64748b;
}

.hero-action-btn {
  border-radius: 14px;
  border-color: #dbe2ea;
  background: #ffffff;
  color: #0f172a;
}

@media (max-width: 768px) {
  .hero-card {
    padding: 22px;
    border-radius: 20px;
  }
}
</style>

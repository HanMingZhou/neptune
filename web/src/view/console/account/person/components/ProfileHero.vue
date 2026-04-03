<template>
  <section class="console-detail-card rounded-xl p-5 md:p-6">
    <div class="flex flex-col gap-5 md:flex-row md:items-start">
      <div class="avatar-shell">
        <el-avatar v-if="userInfo.headerImg" :src="userInfo.headerImg" :size="88" />
        <div v-else class="avatar-fallback">
          {{ userInitials }}
        </div>
      </div>

      <div class="min-w-0 flex-1 space-y-4">
        <div class="space-y-2">
          <p class="text-xs font-semibold uppercase tracking-[0.08em] text-slate-500">{{ t('person') }}</p>

          <div v-if="!editFlag" class="flex flex-col gap-3 md:flex-row md:items-center">
            <h2 class="truncate text-[1.75rem] font-semibold tracking-tight text-slate-950 dark:text-slate-100">
              {{ displayName }}
            </h2>
            <button class="list-row-button list-row-button--neutral" @click="$emit('open-edit')">
              <el-icon><Edit /></el-icon>
              {{ t('editNickname') }}
            </button>
          </div>

          <div v-else class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <el-input v-model="nickNameModel" class="max-w-sm" />
            <div class="flex gap-2">
              <button class="list-toolbar-button list-toolbar-button--primary" @click="$emit('confirm-edit')">
                {{ t('confirm') }}
              </button>
              <button class="list-toolbar-button list-toolbar-button--secondary" @click="$emit('close-edit')">
                {{ t('cancel') }}
              </button>
            </div>
          </div>

          <p class="max-w-3xl text-sm leading-6 text-slate-500">
            {{ t('profileCenterDesc') }}
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
          <div class="detail-inline-chip">
            <el-icon><Message /></el-icon>
            <span>{{ userInfo.email || t('notSet') }}</span>
          </div>
          <div class="detail-inline-chip">
            <el-icon><Phone /></el-icon>
            <span>{{ userInfo.phone || t('notSet') }}</span>
          </div>
          <div class="detail-inline-chip">
            <el-icon><Postcard /></el-icon>
            <span>{{ shortAccountId }}</span>
          </div>
        </div>

        <div class="flex flex-wrap gap-2.5">
          <button class="list-toolbar-button list-toolbar-button--primary" @click="$emit('change-password')">
            <el-icon><Lock /></el-icon>
            {{ t('changePasswordTitle') }}
          </button>
          <button class="list-toolbar-button list-toolbar-button--secondary" @click="$emit('change-email')">
            <el-icon><Message /></el-icon>
            {{ t('emailLabel') }}
          </button>
          <button class="list-toolbar-button list-toolbar-button--secondary" @click="$emit('change-phone')">
            <el-icon><Phone /></el-icon>
            {{ t('phoneLabel') }}
          </button>
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
.avatar-shell {
  display: flex;
  width: 108px;
  height: 108px;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  border: 1px solid #dbe2ea;
  background: #f8fafc;
}

.avatar-fallback {
  display: flex;
  width: 88px;
  height: 88px;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: #e2e8f0;
  font-size: 24px;
  font-weight: 700;
  color: #334155;
}
</style>

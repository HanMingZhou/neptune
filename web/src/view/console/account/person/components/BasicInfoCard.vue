<template>
  <div class="bg-white dark:bg-slate-800 rounded-xl p-6 mb-6 profile-card">
    <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
      <el-icon class="text-blue-500"><InfoFilled /></el-icon>
      {{ t('basicInfo') }}
    </h2>
    <div class="space-y-4">
      <div class="flex items-center gap-1 lg:gap-3 text-gray-600 dark:text-gray-300">
        <el-icon class="text-blue-500"><Phone /></el-icon>
        <span class="font-medium">{{ t('phoneLabel') }}：</span>
        <span>{{ userInfo.phone || t('notSet') }}</span>
        <el-button link type="primary" class="ml-auto" @click="$emit('change-phone')">
          {{ t('edit') }}
        </el-button>
      </div>
      <div class="flex items-center gap-1 lg:gap-3 text-gray-600 dark:text-gray-300">
        <el-icon class="text-green-500"><Message /></el-icon>
        <span class="font-medium flex-shrink-0">{{ t('emailLabel') }}：</span>
        <span>{{ userInfo.email || t('notSet') }}</span>
        <el-button link type="primary" class="ml-auto" @click="$emit('change-email')">
          {{ t('edit') }}
        </el-button>
      </div>
      <div class="flex items-center gap-1 lg:gap-3 text-gray-600 dark:text-gray-300">
        <el-icon class="text-purple-500"><Lock /></el-icon>
        <span class="font-medium">{{ t('confirmPassword') }}：</span>
        <span>{{ t('isSet') }}</span>
        <el-button link type="primary" class="ml-auto" @click="$emit('change-password')">
          {{ t('edit') }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { InfoFilled, Lock, Message, Phone } from '@element-plus/icons-vue'

defineProps({
  userInfo: {
    type: Object,
    required: true
  }
})

defineEmits(['change-phone', 'change-email', 'change-password'])

const t = inject('t', (key) => key)
</script>

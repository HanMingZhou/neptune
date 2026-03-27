<template>
  <div class="bg-white dark:bg-slate-800 rounded-xl p-6 profile-card">
    <el-tabs class="custom-tabs">
      <el-tab-pane>
        <template #label>
          <div class="flex items-center gap-2">
            <el-icon><DataLine /></el-icon>
            {{ t('dataStats') }}
          </div>
        </template>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 lg:gap-6 py-6">
          <div v-for="stat in stats" :key="stat.labelKey" class="stat-card">
            <div class="text-2xl lg:text-4xl font-bold mb-2" :class="stat.valueClass">
              {{ stat.value }}
            </div>
            <div class="text-gray-500 text-sm">{{ t(stat.labelKey) }}</div>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <template #label>
          <div class="flex items-center gap-2">
            <el-icon><Calendar /></el-icon>
            {{ t('recentActivities') }}
          </div>
        </template>
        <div class="py-6">
          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in activities"
              :key="index"
              :type="activity.type"
              :timestamp="activity.timestamp"
              :hollow="true"
              class="pb-6"
            >
              <h3 class="text-base font-medium mb-1">{{ activity.title }}</h3>
              <p class="text-gray-500 text-sm">{{ activity.content }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { Calendar, DataLine } from '@element-plus/icons-vue'

defineProps({
  stats: {
    type: Array,
    default: () => []
  },
  activities: {
    type: Array,
    default: () => []
  }
})

const t = inject('t', (key) => key)
</script>

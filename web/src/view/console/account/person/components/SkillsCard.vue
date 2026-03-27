<template>
  <div class="bg-white dark:bg-slate-800 rounded-xl p-6 profile-card">
    <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
      <el-icon class="text-blue-500"><Medal /></el-icon>
      {{ t('skills') }}
    </h2>
    <div class="flex flex-wrap gap-2">
      <el-tag v-for="skill in skills" :key="skill.label" effect="plain" :type="skill.type">
        {{ skill.label }}
      </el-tag>
      <el-button link class="text-sm">
        <el-icon><Plus /></el-icon>
        {{ t('skillsAdd') }}
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { Medal, Plus } from '@element-plus/icons-vue'

defineProps({
  skills: {
    type: Array,
    default: () => []
  }
})

const t = inject('t', (key) => key)
</script>

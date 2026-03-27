<template>
  <div class="profile-container">
    <ProfileHero
      :edit-flag="editFlag"
      :nick-name="nickName"
      :user-info="userInfo"
      @close-edit="emit('close-edit')"
      @confirm-edit="emit('confirm-edit')"
      @open-edit="emit('open-edit')"
      @update:nick-name="emit('update:nick-name', $event)"
    />

    <div class="grid lg:grid-cols-12 md:grid-cols-1 gap-8">
      <div class="lg:col-span-4">
        <BasicInfoCard
          :user-info="userInfo"
          @change-email="emit('change-email')"
          @change-password="emit('change-password')"
          @change-phone="emit('change-phone')"
        />

        <SkillsCard :skills="skills" />
      </div>

      <div class="lg:col-span-8">
        <ProfileStatsTabs :activities="activities" :stats="stats" />
      </div>
    </div>
  </div>
</template>

<script setup>
import BasicInfoCard from './BasicInfoCard.vue'
import ProfileHero from './ProfileHero.vue'
import ProfileStatsTabs from './ProfileStatsTabs.vue'
import SkillsCard from './SkillsCard.vue'

defineProps({
  activities: {
    type: Array,
    default: () => []
  },
  editFlag: {
    type: Boolean,
    default: false
  },
  nickName: {
    type: String,
    default: ''
  },
  skills: {
    type: Array,
    default: () => []
  },
  stats: {
    type: Array,
    default: () => []
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
</script>

<style>
.profile-container {
  @apply p-4 lg:p-6 min-h-screen bg-gray-50 dark:bg-slate-900;

  .bg-pattern {
    background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23000000' fill-opacity='0.1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }

  .profile-card {
    @apply shadow-sm hover:shadow-md transition-shadow duration-300;
  }

  .stat-card {
    @apply p-4 lg:p-6 rounded-lg bg-gray-50 dark:bg-slate-700/50 text-center hover:shadow-md transition-all duration-300;
  }

  .custom-tabs {
    :deep(.el-tabs__nav-wrap::after) {
      @apply h-0.5 bg-gray-100 dark:bg-gray-700;
    }

    :deep(.el-tabs__active-bar) {
      @apply h-0.5 bg-blue-500;
    }

    :deep(.el-tabs__item) {
      @apply text-base font-medium px-6;

      .el-icon {
        @apply mr-1 text-lg;
      }

      &.is-active {
        @apply text-blue-500;
      }
    }

    :deep(.el-timeline-item__node--normal) {
      @apply left-[-2px];
    }

    :deep(.el-timeline-item__wrapper) {
      @apply pl-8;
    }

    :deep(.el-timeline-item__timestamp) {
      @apply text-gray-400 text-sm;
    }
  }
}
</style>

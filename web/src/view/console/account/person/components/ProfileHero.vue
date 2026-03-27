<template>
  <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm mb-8">
    <div class="h-48 bg-blue-50 dark:bg-slate-600 relative">
      <div class="absolute inset-0 bg-pattern opacity-7"></div>
    </div>

    <div class="px-8 -mt-20 pb-8">
      <div class="flex flex-col lg:flex-row items-start gap-8">
        <div class="profile-avatar-wrapper flex-shrink-0 mx-auto lg:mx-0">
          <SelectImage v-model="userInfo.headerImg" file-type="image" rounded />
        </div>

        <div class="flex-1 pt-12 lg:pt-20 w-full">
          <div class="flex flex-col lg:flex-row items-start lg:items-start justify-between gap-4">
            <div class="lg:mt-4">
              <div class="flex items-center gap-4 mb-4">
                <div
                  v-if="!editFlag"
                  class="text-2xl font-bold flex items-center gap-3 text-gray-800 dark:text-gray-100"
                >
                  {{ userInfo.nickName }}
                  <el-icon
                    class="cursor-pointer text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors duration-200"
                    @click="$emit('open-edit')"
                  >
                    <Edit />
                  </el-icon>
                </div>
                <div v-else class="flex items-center">
                  <el-input v-model="nickNameModel" class="w-48 mr-4" />
                  <el-button type="primary" plain @click="$emit('confirm-edit')">
                    {{ t('confirm') }}
                  </el-button>
                  <el-button type="danger" plain @click="$emit('close-edit')">
                    {{ t('cancel') }}
                  </el-button>
                </div>
              </div>

              <div class="flex flex-col lg:flex-row items-start lg:items-center gap-4 lg:gap-8 text-gray-500 dark:text-gray-400">
                <div class="flex items-center gap-2 mb-2">
                  <el-icon class="text-slate-400"><Location /></el-icon>
                  <span>{{ t('mock.location') }}</span>
                </div>
                <div class="flex items-center gap-2 mb-2">
                  <el-icon class="text-slate-400"><OfficeBuilding /></el-icon>
                  <span>{{ t('mock.company') }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <el-icon class="text-slate-400"><Management /></el-icon>
                  <span>{{ t('mock.department') }}</span>
                </div>
              </div>
            </div>

            <div class="flex gap-4 mt-4">
              <el-button type="primary" plain icon="message">
                {{ t('sendMsg') }}
              </el-button>
              <el-button icon="share">
                {{ t('sharePage') }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { Edit, Location, Management, OfficeBuilding } from '@element-plus/icons-vue'
import SelectImage from '@/components/selectImage/selectImage.vue'

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

const emit = defineEmits(['open-edit', 'close-edit', 'confirm-edit', 'update:nick-name'])

const t = inject('t', (key) => key)

const nickNameModel = computed({
  get: () => props.nickName,
  set: (value) => {
    emit('update:nick-name', value)
  }
})
</script>

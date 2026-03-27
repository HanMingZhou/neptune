<template>
  <div class="w-full">
    <el-select
      v-model="value"
      clearable
      filterable
      :placeholder="t('pleaseSelect')"
      class="w-full"
    >
      <template #prefix>
        <el-icon>
          <component v-if="value" :is="value" />
        </el-icon>
      </template>
      <el-option
        v-for="item in iconOptions"
        :key="item.key"
        class="select__option_item"
        :label="item.key"
        :value="item.key"
      >
        <span class="gva-icon" style="padding: 3px 0 0" :class="item.label">
          <el-icon>
            <component v-if="item.label" :is="item.label" />
          </el-icon>
        </span>
        <span style="text-align: left">{{ item.key }}</span>
      </el-option>
    </el-select>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useMenuIconOptions } from './composables/useMenuIconOptions'

defineOptions({
  name: 'Icon'
})

const t = inject('t', (key) => key)
const value = defineModel()
const { iconOptions } = useMenuIconOptions()
</script>

<style>
.gva-icon {
  color: rgb(132, 146, 166);
  font-size: 14px;
  margin-right: 10px;
}

.select__option_item {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}
</style>

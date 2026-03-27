<template>
  <div class="space-y-1">
    <!-- 菜单项：有路由名则渲染为链接，无路由名则渲染为折叠组 -->
    <template v-if="item.routeName">
      <div
        @click="navigateTo(item.routeName)"
        class="flex items-center gap-3 px-3 py-2 rounded transition-all text-[13px] font-semibold cursor-pointer"
        :class="isActive ? 'sidebar-active shadow-sm' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-zinc-800'"
        :style="{ paddingLeft: `${level * 16 + 12}px` }"
      >
        <span v-if="item.icon" class="material-icons text-[18px] opacity-80">{{ item.icon }}</span>
        <span class="flex-1 truncate">{{ displayTitle }}</span>
      </div>
    </template>
    
    <template v-else>
      <div 
        @click="handleClick"
        class="flex items-center gap-3 px-3 py-2 rounded transition-all text-[13px] font-bold cursor-pointer"
        :class="hasActiveChild ? 'text-primary' : 'text-slate-500 dark:text-slate-500 hover:bg-slate-100 dark:hover:bg-zinc-800'"
        :style="{ paddingLeft: `${level * 16 + 12}px` }"
      >
        <span v-if="item.icon" class="material-icons text-[18px] opacity-80">{{ item.icon }}</span>
        <span class="flex-1 truncate">{{ displayTitle }}</span>
        <span v-if="hasChildren" class="material-icons text-[16px] transition-transform duration-200" :class="{ '-rotate-90': !isOpen }">expand_more</span>
      </div>
    </template>
    
    <!-- 递归渲染子菜单 -->
    <div v-if="hasChildren && isOpen" class="space-y-1">
      <NavMenuItem 
        v-for="child in item.children" 
        :key="child.key" 
        :item="child" 
        :level="level + 1" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const props = defineProps({
  item: Object,
  level: { type: Number, default: 0 }
});

const t = inject('t');
const route = useRoute();
const router = useRouter();
const isOpen = ref(true);

// 菜单标题：优先使用 i18n 翻译，fallback 到后端返回的 title
const displayTitle = computed(() => {
  if (props.item.titleKey && t) {
    const translated = t(props.item.titleKey);
    // 如果翻译结果不等于 key 本身，说明找到了翻译
    if (translated !== props.item.titleKey) {
      return translated;
    }
  }
  return props.item.title;
});

const isActive = computed(() => props.item.routeName === route.name);
const hasChildren = computed(() => props.item.children && props.item.children.length > 0);

// 导航到指定路由
const navigateTo = (routeName) => {
  console.log('Navigating to:', routeName);
  if (routeName) {
    router.push({ name: routeName }).catch(err => {
      console.error('Navigation failed:', err);
    });
  }
};

// 点击处理（折叠/展开）
const handleClick = () => {
  if (hasChildren.value) {
    isOpen.value = !isOpen.value;
  }
};

// 判断是否有子项目处于激活状态
const hasActiveChild = computed(() => {
  if (!props.item.children) return false;
  const check = (items) => items.some(i => i.routeName === route.name || (i.children && check(i.children)));
  return check(props.item.children);
});
</script>

<style scoped>
a {
  text-decoration: none !important;
}
</style>

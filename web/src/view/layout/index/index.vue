<template>
  <el-config-provider :locale="elLocale">
    <div class="flex h-screen w-screen overflow-hidden text-slate-900 dark:text-slate-100 font-sans">
      <LayoutSidebar
        :dynamic-navigation="dynamicNavigation"
        @toggle-theme="toggleTheme"
      />

      <div class="flex-1 flex flex-col min-w-0">
        <LayoutHeaderBar
          :lang="lang"
          :user-balance="userBalance"
          :user-info="userInfo"
          @change-user-auth="changeUserAuth"
          @logout="logout"
          @toggle-lang="toggleLang"
        />

        <main class="flex-1 overflow-y-auto custom-scrollbar p-8 bg-background-light dark:bg-background-dark/50">
          <router-view v-slot="{ Component, route: currentRoute }">
            <keep-alive :include="routerStore.keepAliveRouters">
              <component :is="Component" :key="currentRoute.fullPath" />
            </keep-alive>
          </router-view>
        </main>
      </div>
    </div>
  </el-config-provider>
</template>

<script setup>
import { onMounted, provide } from 'vue'
import LayoutHeaderBar from './components/LayoutHeaderBar.vue'
import LayoutSidebar from './components/LayoutSidebar.vue'
import { useAppLayout } from './composables/useAppLayout'

const {
  changeUserAuth,
  dynamicNavigation,
  elLocale,
  initialize,
  lang,
  logout,
  routerStore,
  t,
  toggleLang,
  toggleTheme,
  userBalance,
  userInfo
} = useAppLayout()

provide('lang', lang)
provide('t', t)

onMounted(() => {
  initialize()
})
</script>

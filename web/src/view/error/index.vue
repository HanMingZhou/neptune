<template>
  <div class="error-page">
    <div class="error-page__card">
      <div class="error-page__badge">
        <span class="material-icons">error_outline</span>
      </div>
      <div class="error-page__code">404</div>
      <p class="error-page__title">页面被神秘力量吸走了</p>
      <p class="error-page__description">
        常见问题为当前此角色无当前路由，如果确定要使用本路由，请到角色管理进行分配。
      </p>
      <p class="error-page__link">
        项目地址：
        <a
          href="https://github.com/flipped-aurora/gin-vue-admin"
          target="_blank"
          rel="noreferrer"
        >
          https://github.com/flipped-aurora/gin-vue-admin
        </a>
      </p>
      <el-button type="primary" @click="toDashboard">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { emitter } from '@/utils/bus'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'Error'
})

const userStore = useUserStore()
const router = useRouter()

const toDashboard = (): void => {
  try {
    void router.push({ name: userStore.userInfo.authority.defaultRouter })
  } catch {
    emitter.emit('show-error', {
      code: '401',
      message: '检测到其他用户修改了路由权限，请重新登录',
      fn: () => {
        void userStore.ClearStorage()
        void router.push({ name: 'Login', replace: true })
      }
    })
  }
}
</script>

<style scoped>
.error-page {
  display: flex;
  min-height: 100vh;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background:
    radial-gradient(circle at top, rgba(59, 130, 246, 0.08), transparent 48%),
    linear-gradient(180deg, #f8fafc 0%, #eef2ff 100%);
}

.error-page__card {
  display: flex;
  width: min(100%, 40rem);
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  border: 1px solid rgba(148, 163, 184, 0.24);
  border-radius: 1.5rem;
  padding: 3rem 2rem;
  text-align: center;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.08);
}

.error-page__badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 5rem;
  height: 5rem;
  border-radius: 9999px;
  color: #1d4ed8;
  background: rgba(59, 130, 246, 0.12);
}

.error-page__badge .material-icons {
  font-size: 2.5rem;
}

.error-page__code {
  font-size: 4rem;
  font-weight: 900;
  line-height: 1;
  letter-spacing: -0.08em;
  color: #0f172a;
}

.error-page__title {
  margin: 0;
  font-size: 1.375rem;
  font-weight: 700;
  color: #0f172a;
}

.error-page__description,
.error-page__link {
  margin: 0;
  font-size: 0.95rem;
  line-height: 1.75;
  color: #475569;
}

.error-page__link a {
  color: #2563eb;
}

.dark .error-page {
  background:
    radial-gradient(circle at top, rgba(59, 130, 246, 0.18), transparent 42%),
    linear-gradient(180deg, #09090b 0%, #111827 100%);
}

.dark .error-page__card {
  border-color: rgba(71, 85, 105, 0.5);
  background: rgba(15, 23, 42, 0.9);
  box-shadow: 0 24px 60px rgba(2, 6, 23, 0.45);
}

.dark .error-page__badge {
  color: #93c5fd;
  background: rgba(59, 130, 246, 0.18);
}

.dark .error-page__code,
.dark .error-page__title {
  color: #f8fafc;
}

.dark .error-page__description,
.dark .error-page__link {
  color: #cbd5e1;
}
</style>

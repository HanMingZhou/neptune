<template>
  <div class="space-y-8">
    <header>
      <h1 class="text-3xl font-black tracking-tight">{{ t('security.center') }}</h1>
      <p class="text-slate-500 mt-2">{{ t('security.centerDesc') }}</p>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="md:col-span-2 space-y-6">
        <SecuritySection :title="t('security.basicTitle')">
          <SecurityItem
            icon="lock"
            :title="t('security.loginPwd')"
            :desc="t('security.pwdDesc')"
            :status-text="t('security.high')"
            status-color="bg-emerald-500/10 text-emerald-500"
            :action-text="t('security.changePwd')"
            @click="emit('change-password')"
          />
          <SecurityItem
            icon="smartphone"
            :title="t('security.phoneBind')"
            :desc="t('security.phoneDesc')"
            :status-text="phoneStatusText"
            :status-color="getStatusColor(accountInfo.phoneStatus)"
            :action-text="accountInfo.phone ? t('security.changePhone') : t('security.bindNow')"
            @click="emit('change-phone')"
          />
          <SecurityItem
            icon="mail"
            :title="t('security.emailBind')"
            :desc="t('security.emailDesc')"
            :status-text="emailStatusText"
            :status-color="getStatusColor(accountInfo.emailStatus)"
            :action-text="accountInfo.email ? t('security.changeEmail') : t('security.bindNow')"
            @click="emit('change-email')"
          />
        </SecuritySection>

        <SecuritySection :title="t('security.mfaTitle')">
          <SecurityItem
            icon="security"
            :title="t('security.mfaDevice')"
            :desc="t('security.mfaDesc')"
            :status-text="accountInfo.mfaStatus"
            :status-color="getStatusColor(accountInfo.mfaStatus)"
            :action-text="accountInfo.mfaEnabled ? t('security.disable') : t('security.enableNow')"
            :loading="mfaLoading"
            @click="emit('mfa-action')"
          />
          <SecurityItem
            icon="notifications_active"
            :title="t('security.loginNotify')"
            :desc="t('security.notifyDesc')"
            :status-text="t('security.enabled')"
            status-color="bg-emerald-500/10 text-emerald-500"
            :action-text="t('security.setup')"
            @click="emit('setup-notification')"
          />
        </SecuritySection>

        <SecuritySection :title="t('security.thirdTitle')">
          <SecurityItem
            icon="hub"
            :title="t('security.githubBind')"
            :desc="t('security.githubDesc')"
            :status-text="githubStatusText"
            :status-color="getStatusColor(accountInfo.githubStatus)"
            :action-text="t('security.unbind')"
            @click="emit('bind-github')"
          />
        </SecuritySection>
      </div>

      <div class="space-y-6">
        <SecurityScoreCard :score="accountInfo.securityScore" />

        <AccessKeyCard
          :access-key-id="accountInfo.accessKeyId"
          :loading="akLoading"
          @copy="emit('copy-access-key')"
          @generate="emit('generate-access-key')"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import SecurityItem from '../../components/SecurityItem.vue'
import AccessKeyCard from './AccessKeyCard.vue'
import SecurityScoreCard from './SecurityScoreCard.vue'
import SecuritySection from './SecuritySection.vue'

const props = defineProps({
  accountInfo: {
    type: Object,
    required: true
  },
  akLoading: {
    type: Boolean,
    default: false
  },
  getStatusColor: {
    type: Function,
    required: true
  },
  maskString: {
    type: Function,
    required: true
  },
  mfaLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'bind-github',
  'change-email',
  'change-password',
  'change-phone',
  'copy-access-key',
  'generate-access-key',
  'mfa-action',
  'setup-notification'
])

const t = inject('t', (key) => key)

const phoneStatusText = computed(() =>
  props.accountInfo.phone
    ? `${props.accountInfo.phoneStatus}: ${props.maskString(props.accountInfo.phone)}`
    : props.accountInfo.phoneStatus
)

const emailStatusText = computed(() =>
  props.accountInfo.email
    ? `${props.accountInfo.emailStatus}: ${props.accountInfo.email}`
    : props.accountInfo.emailStatus
)

const githubStatusText = computed(() =>
  props.accountInfo.githubUsername
    ? `${props.accountInfo.githubStatus}: ${props.accountInfo.githubUsername}`
    : props.accountInfo.githubStatus
)
</script>

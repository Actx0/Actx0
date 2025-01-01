<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full px-4 sm:px-6 lg:px-8 py-8">
      <header class="mb-8">
        <h1 class="text-2xl sm:text-3xl font-semibold text-theme-text tracking-tight">
          {{ $t('dashboard.title') }}
        </h1>
        <p class="mt-1 text-sm text-theme-textLight">
          {{ $t('dashboard.subtitle') }}
        </p>
        <p class="mt-2 text-sm text-theme-text">
          {{ $t('dashboard.welcome_back') }}
          <span class="font-medium">{{ user.name || user.email }}</span>
        </p>
      </header>

      <section class="mb-8" :aria-label="$t('dashboard.key_metrics')">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.memories_stored') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.memoriesStored.toLocaleString() }}</p>
          </div>
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.api_calls_month') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.apiCallsMonth.toLocaleString() }}</p>
          </div>
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.active_agents') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.activeAgents.toLocaleString() }}</p>
          </div>
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.prompts') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.prompts.toLocaleString() }}</p>
          </div>
        </div>
      </section>

      <section class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div class="section overflow-hidden !p-0">
          <div class="border-b border-theme-border px-6 py-4">
            <div>
              <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.attention_title') }}</h2>
              <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.attention_subtitle') }}</p>
            </div>
          </div>

          <p v-if="notificationsLoading" class="px-6 py-8 text-sm text-theme-textLight">{{ $t('common.loading') }}</p>
          <p v-else-if="attentionItems.length === 0" class="px-6 py-8 text-sm text-theme-textLight">
            {{ $t('dashboard.attention_empty') }}
          </p>
          <ul v-else class="divide-y divide-theme-border max-h-[28rem] overflow-y-auto">
            <li
              v-for="item in attentionItems"
              :key="item.id"
              class="flex gap-4 px-6 py-4"
            >
              <span
                class="mt-1.5 h-2 w-2 shrink-0 rounded-full"
                :class="attentionDotClass(item.level)"
                aria-hidden="true"
              />
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-theme-text">{{ item.title }}</p>
                <p class="mt-0.5 text-sm text-theme-textLight leading-relaxed">{{ item.description }}</p>
                <router-link
                  v-if="item.to"
                  :to="item.to"
                  class="mt-2 inline-flex text-sm font-medium text-primary-800 hover:underline"
                >
                  {{ item.action }}
                </router-link>
              </div>
            </li>
          </ul>
        </div>

        <div class="section overflow-hidden !p-0">
          <div class="border-b border-theme-border px-6 py-4">
            <div>
              <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.activity_title') }}</h2>
              <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.activity_subtitle') }}</p>
            </div>
          </div>

          <p v-if="loading" class="px-6 py-8 text-sm text-theme-textLight">{{ $t('common.loading') }}</p>
          <p v-else-if="recentActivity.length === 0" class="px-6 py-8 text-sm text-theme-textLight">
            {{ $t('dashboard.activity_empty') }}
          </p>
          <ul v-else class="divide-y divide-theme-border">
            <li
              v-for="event in recentActivity"
              :key="event.id"
              class="flex gap-3 px-6 py-3.5"
            >
              <span
                class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-theme-hover text-theme-textLight"
                aria-hidden="true"
              >
                <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" :d="activityIcons[event.icon]" />
                </svg>
              </span>
              <div class="min-w-0 flex-1">
                <p class="text-sm text-theme-text">
                  <span class="font-medium">{{ event.actor }}</span>
                  {{ event.action }}
                  <span v-if="event.target" class="font-medium">{{ event.target }}</span>
                </p>
                <p class="mt-0.5 text-xs text-theme-textLight">{{ event.time }}</p>
              </div>
            </li>
          </ul>

          <div v-if="canManage" class="border-t border-theme-border px-6 py-3">
            <router-link to="/audits" class="text-sm font-medium text-primary-800 hover:underline">
              {{ $t('dashboard.activity_view_all') }}
            </router-link>
          </div>
        </div>
      </section>

      <section class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div class="space-y-6 lg:col-span-2">
          <div
            v-if="showWhatsNew"
            class="section overflow-hidden !p-0"
          >
            <div class="border-b border-theme-border px-6 py-4">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="min-w-0">
                  <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.whats_new_title') }}</h2>
                  <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.whats_new_subtitle') }}</p>
                </div>
                <button
                  type="button"
                  class="rounded-md px-2 py-1 text-sm text-theme-textLight transition hover:bg-theme-hover hover:text-theme-text"
                  @click="dismissWhatsNew"
                >
                  {{ $t('dashboard.whats_new_dismiss') }}
                </button>
              </div>
            </div>
            <div class="px-6 py-5">
              <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">
                {{ $t('dashboard.whats_new_release_label') }}
              </p>
              <p class="mt-2 text-base font-medium text-theme-text">{{ $t('dashboard.whats_new_release_title') }}</p>
              <ul class="mt-3 space-y-2 text-sm leading-relaxed text-theme-textLight">
                <li
                  v-for="featureKey in whatsNewFeatures"
                  :key="featureKey"
                  class="flex gap-2"
                >
                  <span class="mt-2 h-1 w-1 shrink-0 rounded-full bg-theme-textLight" aria-hidden="true" />
                  <span>{{ $t(featureKey) }}</span>
                </li>
              </ul>
              <router-link
                to="/docs"
                class="btn-secondary mt-4 inline-flex items-center px-3 py-1.5 text-sm"
              >
                {{ $t('dashboard.whats_new_release_action') }}
              </router-link>
            </div>
          </div>

          <div class="section overflow-hidden !p-0">
            <div class="border-b border-theme-border px-6 py-4">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.top_agents_title') }}</h2>
                  <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.top_agents_subtitle') }}</p>
                </div>
                <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
                  {{ $t('dashboard.mock_badge') }}
                </span>
              </div>
            </div>

            <ul class="divide-y divide-theme-border">
              <li
                v-for="(agent, index) in topAgents"
                :key="agent.id"
                class="flex items-center gap-4 px-6 py-4"
              >
                <span class="w-6 shrink-0 text-sm font-medium text-theme-textLight tabular-nums">
                  {{ index + 1 }}
                </span>
                <div class="min-w-0 flex-1">
                  <router-link
                    :to="agentPath(agent.id)"
                    class="text-sm font-medium text-theme-text hover:text-primary-800 hover:underline"
                  >
                    {{ agent.name }}
                  </router-link>
                </div>
                <div class="hidden sm:block text-right shrink-0">
                  <p class="text-sm font-medium text-theme-text tabular-nums">
                    {{ agent.sessions.toLocaleString() }}
                  </p>
                  <p class="text-xs text-theme-textLight">{{ $t('dashboard.top_agents_sessions') }}</p>
                </div>
                <div class="text-right shrink-0 min-w-[4.5rem]">
                  <p class="text-sm font-medium text-theme-text tabular-nums">
                    {{ agent.memories.toLocaleString() }}
                  </p>
                  <p class="text-xs text-theme-textLight">{{ $t('dashboard.top_agents_memories') }}</p>
                </div>
              </li>
            </ul>

            <div class="border-t border-theme-border px-6 py-3">
              <router-link to="/agents" class="text-sm font-medium text-primary-800 hover:underline">
                {{ $t('dashboard.top_agents_view_all') }}
              </router-link>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="section h-fit overflow-hidden">
            <div class="mx-auto mb-4 flex h-20 w-full items-center justify-center" aria-hidden="true">
              <svg class="h-20 w-auto" viewBox="0 0 120 96" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="8" y="16" width="72" height="64" rx="10" fill="#F5F0EB" stroke="#D4C4B0" stroke-width="1.5"/>
                <rect x="8" y="16" width="72" height="18" rx="10" fill="#37352F"/>
                <circle cx="24" cy="25" r="3" fill="#E8DDD0"/>
                <circle cx="36" cy="25" r="3" fill="#E8DDD0"/>
                <rect x="20" y="44" width="20" height="16" rx="4" fill="#FFFFFF" stroke="#D4C4B0" stroke-width="1.2"/>
                <rect x="48" y="44" width="20" height="16" rx="4" fill="#FFFFFF" stroke="#D4C4B0" stroke-width="1.2"/>
                <rect x="20" y="66" width="20" height="8" rx="3" fill="#E8DDD0"/>
                <rect x="48" y="66" width="20" height="8" rx="3" fill="#37352F" opacity="0.85"/>
                <path d="M88 34c8 0 14 6 14 14v10c0 8-6 14-14 14h-2l-6 8v-8h-6c-8 0-14-6-14-14V48c0-8 6-14 14-14h14z" fill="#FFFFFF" stroke="#37352F" stroke-width="1.5"/>
                <circle cx="82" cy="52" r="2" fill="#37352F"/>
                <circle cx="92" cy="52" r="2" fill="#37352F"/>
                <circle cx="102" cy="52" r="2" fill="#37352F"/>
              </svg>
            </div>
            <h2 class="section-title">{{ $t('dashboard.book_call_title') }}</h2>
            <p class="text-sm text-theme-textLight leading-relaxed">
              {{ $t('dashboard.book_call_desc') }}
            </p>
            <a
              :href="calendlyUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="btn-primary mt-4 inline-flex w-full items-center justify-center"
            >
              {{ $t('dashboard.book_call_cta') }}
            </a>
          </div>

          <div class="section h-fit">
            <h2 class="section-title">{{ $t('dashboard.workspace_info') }}</h2>
            <dl class="space-y-0">
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_name') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right">{{ workspace.name }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_role') }}</dt>
                <dd class="text-sm font-medium text-theme-text capitalize text-right">{{ roleLabel }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_plan') }}</dt>
                <dd class="text-sm font-medium text-theme-text capitalize text-right">{{ planLabel }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_members') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right">{{ memberCount }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  stats_api,
  billing_api,
  workspace_api,
  audit_api,
  notification_api,
} from '@/api'
import AppNav from '@/components/AppNav.vue'
import { user } from '@/lib/auth'
import { agentPath } from '@/lib/agent'
import { saveWorkspaceToStorage } from '@/utils/storage'
import { useWorkspaceContext } from '@/lib/permission'
import { formatRelative } from '@/utils/dt'

const { t, te } = useI18n()
const { currentWorkspace, canManage } = useWorkspaceContext()

const calendlyUrl = 'https://calendly.com/'
const WHATS_NEW_KEY = 'dashboard_whats_new_v1'
const ACTIVITY_LIMIT = 6

const whatsNewFeatures = [
  'dashboard.whats_new_feature_memory',
  'dashboard.whats_new_feature_agents',
  'dashboard.whats_new_feature_knowledge',
  'dashboard.whats_new_feature_prompts',
  'dashboard.whats_new_feature_keys',
  'dashboard.whats_new_feature_observability',
  'dashboard.whats_new_feature_sdks',
]

const workspace = ref({ ...currentWorkspace })
const loading = ref(true)
const notificationsLoading = ref(true)
const plan = ref('hobby')
const showWhatsNew = ref(localStorage.getItem(WHATS_NEW_KEY) !== '1')
const audits = ref([])
const notifications = ref([])

const topAgents = [
  { id: 'agt_support', name: 'Support bot', handle: 'support-bot', sessions: 1842, memories: 9210 },
  { id: 'agt_sales', name: 'Sales assistant', handle: 'sales-assistant', sessions: 963, memories: 4104 },
  { id: 'agt_onboard', name: 'Onboarding guide', handle: 'onboarding-guide', sessions: 511, memories: 1880 },
  { id: 'agt_docs', name: 'Docs helper', handle: 'docs-helper', sessions: 287, memories: 940 },
  { id: 'agt_billing', name: 'Billing FAQ', handle: 'billing-faq', sessions: 124, memories: 356 },
]

const stats = ref({
  memoriesStored: 0,
  apiCallsMonth: 0,
  activeAgents: 0,
  prompts: 0,
})

const planLabel = computed(() => t(`billing_page.plans.${plan.value}.name`))

const roleLabel = computed(() => t(`members_page.role_${workspace.value.role}`))

const memberCount = computed(() => {
  return loading.value ? '…' : workspace.value.membersCount.toLocaleString()
})

const notificationMeta = {
  email_unverified: {
    level: 'warn',
    to: '/profile',
    action: 'dashboard.attention_email_unverified_action',
  },
}

const attentionItems = computed(() =>
  notifications.value.map((item) => {
    const meta = notificationMeta[item.type] || {}
    return {
      id: item.id,
      level: meta.level || 'info',
      title: item.title,
      description: item.body || '',
      action: meta.action ? t(meta.action) : '',
      to: meta.to || null,
    }
  }),
)

const activityIcons = {
  prompt: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
  agent: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z',
  doc: 'M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z',
  member: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
  key: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z',
  memory: 'M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4',
  workspace: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4',
  billing: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z',
  session: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
}

const activityIconByResource = {
  prompt: 'prompt',
  prompt_version: 'prompt',
  agent: 'agent',
  document: 'doc',
  member: 'member',
  invite: 'member',
  access_key: 'key',
  memory: 'memory',
  message: 'memory',
  session: 'session',
  workspace: 'workspace',
  billing: 'billing',
  subscription: 'billing',
}

const recentActivity = computed(() =>
  audits.value.map((event) => ({
    id: event.id,
    icon: activityIconByResource[event.resourceType] || activityIconByResource[event.action?.split('.')[0]] || 'workspace',
    actor: event.actorName || t('common.em_dash'),
    action: activityActionLabel(event.action),
    target: event.meta?.name || event.meta?.handle || event.meta?.email || '',
    time: formatRelative(event.createdAt),
  })),
)

function activityActionLabel(action) {
  const key = `dashboard.action_${action.replaceAll('.', '_')}`
  return te(key) ? t(key) : action
}

function dismissWhatsNew() {
  showWhatsNew.value = false
  localStorage.setItem(WHATS_NEW_KEY, '1')
}

function attentionDotClass(level) {
  return {
    error: 'bg-red-500',
    warn: 'bg-amber-500',
    info: 'bg-sky-500',
  }[level] || 'bg-theme-textLight'
}

async function loadDashboard() {
  loading.value = true
  notificationsLoading.value = true

  const tasks = [
    stats_api.get(currentWorkspace.id)
      .then((res) => {
        stats.value = {
          memoriesStored: res.data.memoriesStored,
          apiCallsMonth: res.data.apiCallsMonth,
          activeAgents: res.data.activeAgents,
          prompts: res.data.prompts,
        }
      })
      .catch(() => {}),

    workspace_api.get(currentWorkspace.id)
      .then((res) => {
        workspace.value = { ...currentWorkspace, ...res.data, role: res.data.role }
        saveWorkspaceToStorage(workspace.value)
      })
      .catch(() => {}),

    notification_api.list(currentWorkspace.id, { limit: 50, offset: 0 })
      .then((res) => { notifications.value = res.data.notifications || [] })
      .catch(() => { notifications.value = [] })
      .finally(() => { notificationsLoading.value = false }),
  ]

  if (canManage.value) {
    tasks.push(
      billing_api.status(currentWorkspace.id)
        .then((res) => { plan.value = res.data.plan })
        .catch(() => {}),
      audit_api.list(currentWorkspace.id, { limit: ACTIVITY_LIMIT, offset: 0 })
        .then((res) => { audits.value = res.data.audits })
        .catch(() => { audits.value = [] }),
    )
  }

  await Promise.allSettled(tasks)
  loading.value = false
}

onMounted(loadDashboard)
</script>

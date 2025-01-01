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
        <p v-if="user?.email" class="mt-2 text-sm text-theme-text">
          {{ $t('dashboard.welcome_back') }}
          <span class="font-medium">{{ user.name || user.email }}</span>
        </p>
      </header>

      <section class="mb-8" aria-label="Key metrics">
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

      <section class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div class="section lg:col-span-2 overflow-hidden !p-0">
          <div class="border-b border-theme-border bg-primary-50/40 px-6 py-5">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.getting_started') }}</h2>
                <p class="mt-1 max-w-2xl text-sm text-theme-textLight leading-relaxed">
                  {{ $t('dashboard.getting_started_desc') }}
                </p>
              </div>
              <span class="rounded-full border border-theme-border bg-white px-3 py-1 text-xs font-medium text-theme-textLight">
                {{ steps.length }} {{ $t('dashboard.steps') }}
              </span>
            </div>
          </div>

          <ol class="relative divide-y divide-theme-border">
            <li
              v-for="(step, index) in steps"
              :key="step.title"
              class="group relative flex gap-4 px-6 py-5 transition-colors hover:bg-theme-hover/40"
            >
              <div class="relative flex flex-col items-center">
                <span
                  class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-theme-border bg-white text-theme-text shadow-sm transition-colors group-hover:border-primary-300 group-hover:bg-primary-50"
                  :class="stepIconClass(step.icon)"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" :d="stepIcons[step.icon]" />
                  </svg>
                </span>
                <span
                  v-if="index < steps.length - 1"
                  class="mt-2 w-px flex-1 min-h-[1.5rem] bg-theme-border"
                  aria-hidden="true"
                />
              </div>
              <div class="min-w-0 flex-1 pb-1">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="text-xs font-medium uppercase tracking-wide text-theme-textLight">
                    {{ $t('dashboard.step_label', { number: step.number }) }}
                  </span>
                </div>
                <p class="mt-1 text-base font-medium text-theme-text">{{ step.title }}</p>
                <p class="mt-1 text-sm leading-relaxed text-theme-textLight">{{ step.description }}</p>
                <router-link
                  v-if="step.link"
                  :to="step.link.to"
                  class="btn-secondary mt-3 inline-flex items-center px-3 py-1.5 text-sm"
                >
                  {{ step.link.label }}
                </router-link>
              </div>
            </li>
          </ol>

          <div class="border-t border-theme-border bg-theme-bg px-6 py-5">
            <p class="text-sm font-medium text-theme-text">{{ $t('dashboard.next_steps') }}</p>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('dashboard.next_steps_desc') }}</p>
            <div class="mt-4 flex flex-wrap gap-2">
              <router-link
                to="/prompts"
                class="rounded-lg border border-theme-border bg-white px-3 py-2 text-sm font-medium text-theme-text transition-colors hover:border-primary-300 hover:bg-primary-50/50"
              >
                {{ $t('dashboard.link_prompts') }}
              </router-link>
              <router-link
                v-if="canManage"
                to="/members"
                class="rounded-lg border border-theme-border bg-white px-3 py-2 text-sm font-medium text-theme-text transition-colors hover:border-primary-300 hover:bg-primary-50/50"
              >
                {{ $t('dashboard.link_members') }}
              </router-link>
              <router-link
                to="/docs"
                class="rounded-lg border border-theme-border bg-white px-3 py-2 text-sm font-medium text-theme-text transition-colors hover:border-primary-300 hover:bg-primary-50/50"
              >
                {{ $t('dashboard.link_docs') }}
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
                <dd class="text-sm font-medium text-theme-text text-right">{{ workspace?.name || '—' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_handle') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right font-mono">{{ workspace?.handle || '—' }}</dd>
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
  statsAPI,
  billingAPI,
  workspaceAPI,
} from '@/api'
import AppNav from '@/components/AppNav.vue'
import { user } from '@/lib/auth'
import { saveWorkspaceToStorage } from '@/utils/storage'
import { useWorkspaceContext } from '@/lib/permission'

const { t } = useI18n()

const calendlyUrl = 'https://calendly.com/'

const { currentWorkspace, canManage } = useWorkspaceContext()
const workspace = ref({ ...currentWorkspace })
const loading = ref(true)
const plan = ref('hobby')

const stats = ref({
  memoriesStored: 0,
  apiCallsMonth: 0,
  activeAgents: 0,
  prompts: 0,
})

const PLAN_LABELS = {
  hobby: 'Hobby',
  starter: 'Starter',
  growth: 'Growth',
  pro: 'Pro',
}

const planLabel = computed(() => PLAN_LABELS[plan.value] || plan.value || 'Hobby')

const roleLabel = computed(() => workspace.value?.role || '—')

const memberCount = computed(() => {
  const count = workspace.value?.membersCount
  if (count == null) return loading.value ? '…' : '—'
  return count.toLocaleString()
})

const stepIcons = {
  key: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z',
  package: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4',
  plug: 'M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1',
  memory: 'M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4',
  search: 'M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z',
}

function stepIconClass(icon) {
  const map = {
    key: 'text-amber-700',
    package: 'text-violet-700',
    plug: 'text-sky-700',
    memory: 'text-emerald-700',
    search: 'text-primary-800',
  }
  return map[icon] ?? 'text-theme-text'
}

const steps = computed(() => {
  const list = [
    {
      number: 1,
      icon: 'key',
      title: t('dashboard.step_access_key'),
      description: t('dashboard.step_access_key_desc'),
      link: canManage.value
        ? { to: '/settings', label: t('dashboard.step_access_key_link') }
        : null,
    },
    {
      number: 2,
      icon: 'package',
      title: t('dashboard.step_install_sdk'),
      description: t('dashboard.step_install_sdk_desc'),
    },
    {
      number: 3,
      icon: 'plug',
      title: t('dashboard.step_init_client'),
      description: t('dashboard.step_init_client_desc'),
    },
    {
      number: 4,
      icon: 'memory',
      title: t('dashboard.step_store_memory'),
      description: t('dashboard.step_store_memory_desc'),
    },
    {
      number: 5,
      icon: 'search',
      title: t('dashboard.step_search_memory'),
      description: t('dashboard.step_search_memory_desc'),
    },
  ]
  return list
})

async function loadDashboard() {
  loading.value = true

  const tasks = [
    statsAPI.get(currentWorkspace.id)
      .then((res) => {
        stats.value = {
          memoriesStored: res.data.memoriesStored ?? 0,
          apiCallsMonth: res.data.apiCallsMonth ?? 0,
          activeAgents: res.data.activeAgents ?? 0,
          prompts: res.data.prompts ?? 0,
        }
      })
      .catch(() => {}),

    workspaceAPI.get(currentWorkspace.id)
      .then((res) => {
        workspace.value = { ...currentWorkspace, ...res.data, role: res.data.role ?? currentWorkspace.role }
        saveWorkspaceToStorage(workspace.value)
      })
      .catch(() => {}),
  ]

  if (canManage.value) {
    tasks.push(
      billingAPI.status(currentWorkspace.id)
        .then((res) => { plan.value = res.data.plan || 'hobby' })
        .catch(() => {}),
    )
  }

  await Promise.allSettled(tasks)
  loading.value = false
}

onMounted(loadDashboard)
</script>

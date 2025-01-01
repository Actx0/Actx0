<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full px-4 py-8 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">Billing</h1>
        <p class="mt-2 text-theme-textLight">Manage your workspace plan and subscription.</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
        {{ errorMessage }}
      </div>

      <div v-if="billingLoading" class="space-y-6">
        <div class="h-72 animate-pulse rounded-xl border border-theme-border bg-white" />
        <div class="h-96 animate-pulse rounded-xl border border-theme-border bg-white" />
      </div>

      <template v-else>
        <div class="grid grid-cols-1 gap-6 xl:grid-cols-3">
          <section class="overflow-hidden rounded-xl border border-theme-border bg-white shadow-sm xl:col-span-2">
            <div class="border-b border-theme-border bg-primary-50/40 px-6 py-4">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <p class="text-xs font-medium uppercase tracking-wider text-theme-textLight">Your plan</p>
                  <h2 class="mt-1 text-2xl font-semibold text-theme-text">{{ currentPlan.name }}</h2>
                </div>
                <span
                  class="rounded-full px-3 py-1 text-xs font-medium"
                  :class="isPaidPlan ? 'bg-primary-100 text-theme-text' : 'bg-theme-hover text-theme-textLight'"
                >
                  {{ isPaidPlan ? 'Paid plan' : 'Free plan' }}
                </span>
              </div>
            </div>

            <div class="p-6">
              <p class="text-3xl font-semibold text-theme-text">
                {{ currentPlan.price }}
                <span class="text-base font-normal text-theme-textLight">{{ currentPlan.unit }}</span>
              </p>
              <p class="mt-3 text-sm leading-6 text-theme-textLight">{{ currentPlan.description }}</p>

              <h3 class="mt-6 text-sm font-semibold text-theme-text">Included in this plan</h3>
              <ul class="mt-3 space-y-2.5">
                <li
                  v-for="feature in currentPlan.features"
                  :key="feature"
                  class="flex items-start gap-2.5 text-sm text-theme-textLight"
                >
                  <svg class="mt-0.5 h-4 w-4 shrink-0 text-theme-text" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>{{ feature }}</span>
                </li>
              </ul>

              <div class="mt-8 flex flex-wrap gap-3">
                <button type="button" class="btn-primary" @click="planModalOpen = true">
                  {{ isPaidPlan ? 'Change plan' : 'Upgrade plan' }}
                </button>
                <button
                  v-if="canManageStripe"
                  type="button"
                  class="btn-secondary"
                  :disabled="portalLoading"
                  @click="openPortal"
                >
                  {{ portalLoading ? 'Opening...' : 'Manage billing' }}
                </button>
              </div>
            </div>
          </section>

          <aside class="space-y-6">
            <section class="rounded-xl border border-theme-border bg-white p-6 shadow-sm">
              <h2 class="text-lg font-semibold text-theme-text mb-4">Subscription</h2>
              <dl class="space-y-0">
                <div class="flex justify-between gap-4 border-b border-theme-border py-3">
                  <dt class="text-sm text-theme-textLight">Plan</dt>
                  <dd class="text-sm font-medium text-theme-text">{{ currentPlan.name }}</dd>
                </div>
                <div class="flex justify-between gap-4 border-b border-theme-border py-3">
                  <dt class="text-sm text-theme-textLight">Status</dt>
                  <dd class="text-sm font-medium capitalize text-theme-text">{{ billingStatus.status || 'active' }}</dd>
                </div>
                <div v-if="billingStatus.updatedAt" class="flex justify-between gap-4 py-3">
                  <dt class="text-sm text-theme-textLight">Last updated</dt>
                  <dd class="text-sm font-medium text-theme-text">{{ formatDate(billingStatus.updatedAt) }}</dd>
                </div>
              </dl>
            </section>

            <section class="rounded-xl border border-theme-border bg-white p-6 shadow-sm">
              <h2 class="text-lg font-semibold text-theme-text mb-4">Payment & cancellation</h2>
              <p class="text-sm leading-6 text-theme-textLight">
                <template v-if="canManageStripe">
                  Update your payment method, review charges, download receipts, or cancel your subscription.
                </template>
                <template v-else-if="isPaidPlan">
                  Billing details are still syncing. Refresh in a moment or contact support if this persists.
                </template>
                <template v-else>
                  Upgrade to a paid plan to manage payment and cancellation.
                </template>
              </p>
              <button
                type="button"
                class="btn-secondary mt-4 w-full"
                :class="{ 'cursor-not-allowed opacity-60': !canManageStripe || portalLoading }"
                :disabled="!canManageStripe || portalLoading"
                @click="openPortal"
              >
                {{ canManageStripe ? (portalLoading ? 'Opening...' : 'Open billing portal') : 'Available after upgrade' }}
              </button>
            </section>

            <section v-if="!isPaidPlan" class="rounded-lg border border-theme-border bg-primary-50/50 p-5">
              <h2 class="text-sm font-semibold text-theme-text">Ready for more?</h2>
              <p class="mt-2 text-sm leading-6 text-theme-textLight">
                Starter begins at $19/mo with higher limits and email support.
              </p>
              <button type="button" class="btn-primary mt-4 w-full" @click="planModalOpen = true">
                View plans
              </button>
            </section>
          </aside>
        </div>

        <section class="mt-6 overflow-hidden rounded-xl border border-theme-border bg-white shadow-sm">
          <div class="border-b border-theme-border bg-primary-50/40 px-6 py-4">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold text-theme-text">Usage</h2>
                <p class="mt-1 text-sm text-theme-textLight">
                  Track usage against your {{ currentPlan.name }} plan limits.
                </p>
              </div>
              <span class="rounded-full bg-theme-hover px-3 py-1 text-xs font-medium text-theme-textLight">
                {{ currentPlan.name }} limits
              </span>
            </div>
          </div>

          <div class="space-y-6 px-6 py-5">
            <div
              v-for="metric in usageMetrics"
              :key="metric.id"
            >
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="text-sm font-medium text-theme-text">{{ metric.label }}</p>
                    <span
                      class="rounded px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide"
                      :class="metric.metricBadgeClass"
                    >
                      {{ metric.badgeShort }}
                    </span>
                  </div>
                  <p v-if="metric.hint" class="mt-0.5 text-xs text-theme-textLight">{{ metric.hint }}</p>
                </div>
                <p class="text-sm font-medium tabular-nums text-theme-text">
                  <span>{{ metric.displayUsed }}</span>
                  <span class="font-normal text-theme-textLight"> / {{ metric.displayLimit }}</span>
                </p>
              </div>

              <div class="mt-3 h-2 overflow-hidden rounded-full bg-theme-hover">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :style="metric.barStyle"
                />
              </div>

              <p class="mt-2 text-xs text-theme-textLight">
                {{ metric.percentLabel }} {{ metric.percentSuffix }}
              </p>
            </div>
          </div>

          <div class="border-t border-theme-border bg-theme-hover/30 px-6 py-4">
            <p class="text-xs leading-5 text-theme-textLight">
              <span class="font-medium text-theme-text">Monthly</span> metrics reset on {{ usagePeriodReset }}.
              <span class="font-medium text-theme-text">Ongoing</span> metrics reflect total workspace storage and do not reset — delete data or upgrade to free capacity.
            </p>
          </div>
        </section>
      </template>
    </main>

    <Teleport to="body">
      <Transition name="modal">
        <div v-if="planModalOpen" class="fixed inset-0 z-50 overflow-y-auto">
          <div class="fixed inset-0 bg-black/40 backdrop-blur-sm" @click="planModalOpen = false"></div>
          <div class="relative flex min-h-full items-center justify-center p-4">
            <section class="relative w-full max-w-6xl rounded-xl border border-theme-border bg-white shadow-xl">
              <button
                type="button"
                class="absolute right-5 top-5 text-theme-textLight hover:text-theme-text"
                aria-label="Close plan dialog"
                @click="planModalOpen = false"
              >
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18 18 6M6 6l12 12" />
                </svg>
              </button>

              <div class="px-6 pb-6 pt-8 sm:px-10">
                <div class="text-center">
                  <h2 class="text-3xl font-semibold text-theme-text">Adjust your plan</h2>
                  <p class="mt-3 text-sm text-theme-textLight">Plans are billed monthly.</p>
                </div>

                <div class="mt-8 grid gap-4 lg:grid-cols-4">
                  <article
                    v-for="plan in plans"
                    :key="plan.name"
                    class="flex min-h-[28rem] flex-col rounded-xl border border-theme-border bg-white p-5 shadow-sm"
                    :class="{ 'bg-theme-hover ring-1 ring-theme-border': isCurrentPlan(plan) }"
                  >
                    <div class="flex items-center justify-between gap-3">
                      <h3 class="text-lg font-semibold text-theme-text">{{ plan.name }}</h3>
                      <span v-if="isCurrentPlan(plan)" class="rounded-md bg-primary-100 px-2 py-1 text-xs font-medium text-theme-textLight">
                        Current plan
                      </span>
                    </div>
                    <p class="mt-4">
                      <span class="text-3xl font-semibold text-theme-text">{{ plan.price }}</span>
                      <span class="text-sm text-theme-textLight">{{ plan.unit }}</span>
                    </p>
                    <p class="mt-4 text-sm leading-5 text-theme-textLight">{{ plan.description }}</p>
                    <ul class="mt-5 space-y-3">
                      <li v-for="feature in plan.features" :key="feature" class="flex gap-2 text-sm leading-5 text-theme-textLight">
                        <span class="text-theme-textLight">&check;</span>
                        <span>{{ feature }}</span>
                      </li>
                    </ul>
                    <button
                      v-if="plan.id !== 'hobby'"
                      type="button"
                      class="mt-auto w-full"
                      :class="isCurrentPlan(plan) ? 'rounded-md border border-theme-border px-4 py-2 text-sm font-medium text-theme-textLight cursor-default' : 'btn-primary'"
                      :disabled="isCurrentPlan(plan) || checkoutPlan === plan.id"
                      @click="choosePlan(plan)"
                    >
                      {{ isCurrentPlan(plan) ? 'Your current plan' : (checkoutPlan === plan.id ? 'Opening...' : plan.cta) }}
                    </button>
                  </article>
                </div>
              </div>

              <div class="border-t border-theme-border px-6 py-5 text-center text-sm text-theme-textLight">
                Need more capabilities for your business?
                <button type="button" class="text-theme-text underline underline-offset-2">Learn more about our Enterprise plans.</button>
              </div>
            </section>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watchEffect } from 'vue'
import AppNav from '@/components/AppNav.vue'
import { billingAPI } from '@/api'
import { useWorkspaceContext } from '@/lib/permission'

const { currentWorkspace } = useWorkspaceContext()

const planModalOpen = ref(false)
const billingLoading = ref(false)
const portalLoading = ref(false)
const checkoutPlan = ref('')
const errorMessage = ref('')
const billingStatus = ref({
  plan: 'hobby',
  status: 'active',
  providerCustomerID: '',
  updatedAt: ''
})

const usageData = ref({
  used: {
    apiCalls: 0,
    workspaceMembers: 0,
    documentsCount: 0,
    promptsCount: 0,
    storageGB: 0,
    historyRecords: 0,
    memoryRecords: 0,
    aiTokens: 0
  },
  limits: {
    apiCalls: 10_000,
    workspaceMembers: 3,
    documentsCount: 100,
    promptsCount: 50,
    storageGB: 5,
    historyRecords: 1_000,
    memoryRecords: 500,
    aiTokens: 500_000
  },
  periodReset: ''
})

const plans = [
  {
    id: 'hobby',
    name: 'Hobby',
    price: '$0',
    unit: '/mo.',
    description: 'Free workspace plan to get started with core capabilities.',
    cta: 'Choose plan',
    features: [
      'Core workspace access',
      'Basic usage limits',
      'Community support',
      'Upgrade anytime'
    ]
  },
  {
    id: 'starter',
    name: 'Starter',
    price: '$19',
    unit: '/mo.',
    description: 'A low-cost plan for individuals and small projects that need more room.',
    cta: 'Choose plan',
    features: [
      'Everything in Hobby',
      'Higher usage limits',
      'Monthly billing',
      'Email support'
    ]
  },
  {
    id: 'growth',
    name: 'Growth',
    price: '$69',
    unit: '/mo.',
    description: 'More capacity for growing teams and heavier workloads.',
    cta: 'Choose plan',
    features: [
      'Everything in Starter',
      'Expanded usage limits',
      'Priority capacity',
      'Faster support'
    ]
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '€219',
    unit: '/mo.',
    description: 'Advanced plan for teams that need the highest monthly capacity.',
    cta: 'Choose plan',
    features: [
      'Everything in Growth',
      'Highest usage limits',
      'Premium capacity',
      'Priority support'
    ]
  }
]

const currentPlan = computed(() => {
  return plans.find((plan) => plan.id === billingStatus.value.plan) || plans[0]
})

const isPaidPlan = computed(() => billingStatus.value.plan !== 'hobby')

const canManageStripe = computed(() => {
  return billingStatus.value.providerCustomerID !== '' && billingStatus.value.plan !== 'hobby'
})

const isCurrentPlan = (plan) => plan.id === billingStatus.value.plan

const usagePeriodReset = computed(() => {
  if (usageData.value.periodReset) {
    return formatDate(usageData.value.periodReset)
  }

  const now = new Date()
  const next = new Date(now.getFullYear(), now.getMonth() + 1, 1)
  return next.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
})

const formatNumber = (value) => value.toLocaleString(undefined, { maximumFractionDigits: 0 })

const formatCompact = (value) => {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1).replace(/\.0$/, '')}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(1).replace(/\.0$/, '')}K`
  return formatNumber(value)
}

const usagePercent = (used, limit) => Math.min(100, Math.round((used / limit) * 100))

const usageBarColor = (percent) => {
  const p = Math.min(100, Math.max(0, percent)) / 100
  return `color-mix(in srgb, var(--primary-400) ${Math.round((1 - p) * 100)}%, var(--theme-text))`
}

const usageBarStyle = (percent) => ({
  width: `${Math.min(100, Math.max(0, percent))}%`,
  backgroundColor: usageBarColor(percent)
})

const usageMetricDisplay = {
  monthly: {
    badgeShort: 'Monthly',
    metricBadgeClass: 'bg-primary-50 text-theme-textLight',
    percentSuffix: 'of monthly limit used'
  },
  capacity: {
    badgeShort: 'Ongoing',
    metricBadgeClass: 'bg-theme-hover text-theme-textLight',
    percentSuffix: 'of plan capacity used'
  }
}

const buildUsageMetrics = (limits, used) => {
  const definitions = [
    {
      id: 'apiCalls',
      resetPolicy: 'monthly',
      label: 'API Calls',
      hint: 'Total API requests this billing period',
      used: used.apiCalls,
      limit: limits.apiCalls,
      displayUsed: formatNumber(used.apiCalls),
      displayLimit: formatNumber(limits.apiCalls)
    },
    {
      id: 'aiTokens',
      resetPolicy: 'monthly',
      label: 'AI Tokens',
      hint: 'Tokens consumed by AI features this period',
      used: used.aiTokens,
      limit: limits.aiTokens,
      displayUsed: formatCompact(used.aiTokens),
      displayLimit: formatCompact(limits.aiTokens)
    },
    {
      id: 'workspaceMembers',
      resetPolicy: 'capacity',
      label: 'Workspace Members',
      hint: 'Active members in this workspace',
      used: used.workspaceMembers,
      limit: limits.workspaceMembers,
      displayUsed: formatNumber(used.workspaceMembers),
      displayLimit: formatNumber(limits.workspaceMembers)
    },
    {
      id: 'documentsCount',
      resetPolicy: 'capacity',
      label: 'Documents Count',
      hint: 'Documents stored in this workspace',
      used: used.documentsCount,
      limit: limits.documentsCount,
      displayUsed: formatNumber(used.documentsCount),
      displayLimit: formatNumber(limits.documentsCount)
    },
    {
      id: 'promptsCount',
      resetPolicy: 'capacity',
      label: 'Prompts Count',
      hint: 'Prompts stored in this workspace',
      used: used.promptsCount,
      limit: limits.promptsCount,
      displayUsed: formatNumber(used.promptsCount),
      displayLimit: formatNumber(limits.promptsCount)
    },
    {
      id: 'storageGB',
      resetPolicy: 'capacity',
      label: 'Storage Used',
      hint: 'Total file storage consumed',
      used: used.storageGB,
      limit: limits.storageGB,
      displayUsed: `${used.storageGB.toFixed(1)} GB`,
      displayLimit: `${limits.storageGB} GB`
    },
    {
      id: 'historyRecords',
      resetPolicy: 'capacity',
      label: 'History Records',
      hint: 'Agent context rows stored across sessions',
      used: used.historyRecords,
      limit: limits.historyRecords,
      displayUsed: formatNumber(used.historyRecords),
      displayLimit: formatNumber(limits.historyRecords)
    },
    {
      id: 'memoryRecords',
      resetPolicy: 'capacity',
      label: 'Memory Records',
      hint: 'Agent memory rows stored across sessions',
      used: used.memoryRecords,
      limit: limits.memoryRecords,
      displayUsed: formatNumber(used.memoryRecords),
      displayLimit: formatNumber(limits.memoryRecords)
    }
  ]

  return definitions.map((metric) => {
    const percent = usagePercent(metric.used, metric.limit)
    const display = usageMetricDisplay[metric.resetPolicy]
    return {
      ...metric,
      ...display,
      percent,
      percentLabel: `${percent}%`,
      barStyle: usageBarStyle(percent)
    }
  })
}

const usageMetrics = computed(() => {
  return buildUsageMetrics(usageData.value.limits, usageData.value.used)
})

const formatDate = (value) => {
  if (!value) return ''
  return new Date(value).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const loadBilling = async () => {
  billingLoading.value = true
  errorMessage.value = ''
  try {
    const [statusRes, usageRes] = await Promise.all([
      billingAPI.status(currentWorkspace.id),
      billingAPI.usage(currentWorkspace.id)
    ])
    billingStatus.value = statusRes.data
    usageData.value = usageRes.data
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load billing.'
  } finally {
    billingLoading.value = false
  }
}

const choosePlan = async (plan) => {
  if (isCurrentPlan(plan) || plan.id === 'hobby') return

  checkoutPlan.value = plan.id
  errorMessage.value = ''
  try {
    const res = await billingAPI.checkout(currentWorkspace.id, { plan: plan.id })
    window.location.href = res.data.url
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to start checkout.'
  } finally {
    checkoutPlan.value = ''
  }
}

const openPortal = async () => {
  if (!canManageStripe.value) return

  portalLoading.value = true
  errorMessage.value = ''
  try {
    const res = await billingAPI.portal(currentWorkspace.id)
    window.location.href = res.data.url
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to open billing portal.'
  } finally {
    portalLoading.value = false
  }
}

onMounted(loadBilling)

watchEffect(() => {
  if (!planModalOpen.value) {
    document.body.style.overflow = ''
    return
  }

  document.body.style.overflow = 'hidden'
  const onEsc = (event) => {
    if (event.key === 'Escape') planModalOpen.value = false
  }
  document.addEventListener('keydown', onEsc)

  return () => {
    document.body.style.overflow = ''
    document.removeEventListener('keydown', onEsc)
  }
})
</script>

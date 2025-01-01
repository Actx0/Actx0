<template>
  <div ref="rootRef" class="landing-page relative min-h-screen overflow-x-hidden">
    <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
    <div class="pointer-events-none absolute inset-x-0 top-0 h-[520px] hero-grid opacity-60" />

    <LandingNav home-page />

    <main>
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-16 md:pt-24">
        <div class="mx-auto max-w-3xl text-center">
          <div class="opacity-0 animate-fade-up">
            <h1 class="hero-gradient-text text-4xl font-semibold leading-tight tracking-tight md:text-6xl md:leading-[1.1]">
              {{ $t('landing.hero.title') }}
            </h1>
          </div>
          <p class="mx-auto mt-6 max-w-2xl text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
            {{ $t('landing.hero.subtitle') }}
          </p>
          <div class="mt-8 flex flex-wrap items-center justify-center gap-3 opacity-0 animate-fade-up animate-delay-300">
            <router-link :to="startPath" class="landing-btn">{{ $t('landing.hero.cta') }}</router-link>
          </div>
        </div>

        <div class="reveal mx-auto mt-16 max-w-3xl">
          <div class="code-window">
            <div class="flex items-center justify-between border-b border-white/10 bg-primary-800 px-4 py-3">
              <div class="flex items-center gap-2">
                <span class="h-3 w-3 rounded-full bg-red-500/80" />
                <span class="h-3 w-3 rounded-full bg-amber-500/80" />
                <span class="h-3 w-3 rounded-full bg-emerald-500/80" />
              </div>
              <div class="flex gap-1">
                <button
                  v-for="tab in codeTabs"
                  :key="tab"
                  type="button"
                  class="code-tab"
                  :class="{ active: codeTab === tab }"
                  @click="codeTab = tab"
                >
                  {{ $t(`landing.code.tabs.${tab}`) }}
                </button>
              </div>
              <button type="button" class="landing-btn-outline px-3 py-1.5 text-xs" @click="copyCode">
                {{ copyLabel }}
              </button>
            </div>
            <div class="overflow-x-auto p-5">
              <pre class="code-line" v-html="highlightedCode" />
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-12">
        <div class="mx-auto grid max-w-6xl grid-cols-2 gap-8 px-6 md:grid-cols-4">
          <div v-for="stat in stats" :key="stat.labelKey" class="reveal text-center">
            <p class="text-3xl font-semibold md:text-4xl">{{ $t(stat.valueKey) }}</p>
            <p class="mt-1 text-sm text-muted">{{ $t(stat.labelKey) }}</p>
          </div>
        </div>
      </section>

      <section id="developers" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">{{ $t('landing.developers.label') }}</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('landing.developers.title') }}</h2>
          <p class="mt-4 text-muted">{{ $t('landing.developers.subtitle') }}</p>
        </div>

        <div class="reveal mt-12">
          <div class="mx-auto flex w-fit gap-2 rounded-xl border border-theme-border bg-theme-hover p-1">
            <button
              v-for="tab in featureTabs"
              :key="tab"
              type="button"
              class="feature-tab"
              :class="{ active: featureTab === tab }"
              @click="featureTab = tab"
            >
              {{ $t(`landing.developers.tabs.${tab}`) }}
            </button>
          </div>

          <div class="mt-8 grid gap-8 md:grid-cols-2">
            <template v-if="featureTab === 'efficiency'">
              <div class="glass-card bg-card-glow p-8">
                <h3 class="text-xl font-semibold">{{ $t('landing.developers.efficiency.title') }}</h3>
                <p class="mt-3 text-muted">{{ $t('landing.developers.efficiency.description') }}</p>
                <ul class="mt-6 space-y-3 text-sm text-theme-text">
                  <li v-for="key in efficiencyItems" :key="key" class="flex items-start gap-2">
                    <span class="mt-1 text-primary-700">→</span>
                    {{ $t(key) }}
                  </li>
                </ul>
              </div>
              <div class="glass-card flex items-center justify-center p-8">
                <div class="w-full space-y-3 font-mono text-xs">
                  <div class="rounded-lg bg-theme-hover p-4 text-primary-500">{{ $t('landing.developers.efficiency.before') }}</div>
                  <div class="rounded-lg border border-primary-300 bg-primary-50 p-4 text-primary-800">{{ $t('landing.developers.efficiency.after') }}</div>
                  <div class="h-2 overflow-hidden rounded-full bg-primary-200">
                    <div class="h-full w-[10%] rounded-full bg-primary-800" />
                  </div>
                  <p class="text-center text-primary-500">{{ $t('landing.developers.efficiency.compression') }}</p>
                </div>
              </div>
            </template>

            <div v-else-if="featureTab === 'visibility'" class="glass-card bg-card-glow p-8 md:col-span-2">
              <div class="grid gap-8 md:grid-cols-2">
                <div>
                  <h3 class="text-xl font-semibold">{{ $t('landing.developers.visibility.title') }}</h3>
                  <p class="mt-3 text-muted">{{ $t('landing.developers.visibility.description') }}</p>
                </div>
                <div class="space-y-2">
                  <div class="log-line">{{ $t('landing.developers.visibility.log_1') }}</div>
                  <div class="log-line">{{ $t('landing.developers.visibility.log_2') }}</div>
                  <div class="log-line-muted">{{ $t('landing.developers.visibility.log_3') }}</div>
                </div>
              </div>
            </div>

            <div v-else class="glass-card bg-card-glow p-8 md:col-span-2">
              <div class="grid gap-8 md:grid-cols-2">
                <div>
                  <h3 class="text-xl font-semibold">{{ $t('landing.developers.control.title') }}</h3>
                  <p class="mt-3 text-muted">{{ $t('landing.developers.control.description') }}</p>
                </div>
                <div class="grid grid-cols-2 gap-3 text-sm">
                  <div v-for="item in controlItems" :key="item.titleKey" class="inner-card">
                    <p class="font-medium">{{ $t(item.titleKey) }}</p>
                    <p class="mt-1 text-muted">{{ $t(item.textKey) }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">{{ $t('landing.how_it_works.label') }}</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('landing.how_it_works.title') }}</h2>
          </div>
          <div class="mt-16 grid gap-6 md:grid-cols-3">
            <div v-for="(step, index) in howSteps" :key="step.titleKey" class="reveal step-card">
              <div class="step-badge">{{ index + 1 }}</div>
              <h3 class="text-lg font-semibold">{{ $t(step.titleKey) }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t(step.textKey) }}</p>
            </div>
          </div>
        </div>
      </section>

      <section id="use-cases" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">{{ $t('landing.use_cases.label') }}</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('landing.use_cases.title') }}</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <router-link
            v-for="card in useCaseCards"
            :key="card.to"
            :to="card.to"
            class="reveal use-case-card block p-6 transition hover:border-primary-300"
          >
            <p class="use-case-label">{{ $t(card.labelKey) }}</p>
            <h3 class="mt-3 text-lg font-semibold">{{ $t(card.titleKey) }}</h3>
            <p class="mt-2 text-sm text-muted">{{ $t(card.textKey) }}</p>
          </router-link>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">{{ $t('landing.enterprise.label') }}</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('landing.enterprise.title') }}</h2>
            <p class="mt-4 text-muted">{{ $t('landing.enterprise.subtitle') }}</p>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div v-for="item in enterpriseItems" :key="item.titleKey" class="reveal glass-card p-6">
              <h3 class="font-semibold">{{ $t(item.titleKey) }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t(item.textKey) }}</p>
            </div>
          </div>
        </div>
      </section>

      <section id="pricing" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">{{ $t('landing.pricing.label') }}</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('landing.pricing.title') }}</h2>
          <p class="mt-4 text-muted">{{ $t('landing.pricing.subtitle') }}</p>
        </div>

        <div class="mt-12 grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="plan in plans"
            :key="plan.id"
            class="reveal glass-card relative flex flex-col p-8"
            :class="{ 'pricing-popular': plan.popular }"
          >
            <span v-if="plan.popular" class="pricing-badge">{{ $t('landing.pricing.popular') }}</span>
            <p class="text-sm font-medium text-muted">{{ $t(plan.nameKey) }}</p>
            <p class="mt-2 text-3xl font-semibold">
              {{ plan.price }}<span class="text-lg font-normal text-muted">{{ $t('landing.pricing.unit') }}</span>
            </p>
            <p class="mt-2 text-sm text-muted">{{ $t(plan.descriptionKey) }}</p>
            <ul class="mt-6 flex-1 space-y-2 text-sm text-theme-text">
              <li v-for="featureKey in plan.featureKeys" :key="featureKey">{{ $t(featureKey) }}</li>
            </ul>
            <router-link
              :to="startPath"
              class="mt-8 w-full"
              :class="plan.popular ? 'landing-btn' : 'landing-btn-outline'"
            >
              {{ $t(plan.ctaKey) }}
            </router-link>
          </div>
        </div>

        <div class="reveal glass-card mx-auto mt-8 flex max-w-4xl flex-col items-start justify-between gap-6 p-8 md:flex-row md:items-center">
          <div>
            <p class="text-sm font-medium text-muted">{{ $t('landing.pricing.enterprise_label') }}</p>
            <h3 class="mt-1 text-2xl font-semibold">{{ $t('landing.pricing.enterprise_title') }}</h3>
            <p class="mt-2 text-sm text-muted">{{ $t('landing.pricing.enterprise_description') }}</p>
          </div>
          <a href="mailto:hello@clivern.com" class="landing-btn-outline shrink-0">{{ $t('landing.pricing.enterprise_cta') }}</a>
        </div>
      </section>

      <section id="get-started" class="relative overflow-hidden py-24">
        <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
        <div class="relative mx-auto max-w-3xl px-6 text-center">
          <h2 class="reveal text-4xl font-semibold tracking-tight md:text-5xl">
            {{ $t('landing.cta.title_before') }} <span class="hero-gradient-text">{{ $t('landing.cta.title_highlight') }}</span>
          </h2>
          <p class="reveal mt-4 text-lg text-muted">{{ $t('landing.cta.subtitle') }}</p>
          <div class="reveal mt-8 flex flex-wrap items-center justify-center gap-3">
            <router-link :to="startPath" class="landing-btn">{{ $t('landing.cta.get_started') }}</router-link>
            <a href="#pricing" class="landing-btn-outline">{{ $t('landing.cta.see_pricing') }}</a>
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-theme-border bg-theme-hover">
      <div class="mx-auto max-w-6xl px-6 py-12">
        <div class="flex flex-col gap-8 md:flex-row md:justify-between">
          <div>
            <router-link to="/" class="flex items-center gap-2.5">
              <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-7 w-auto" />
              <span class="font-semibold">{{ $t('common.brand') }}</span>
            </router-link>
            <p class="mt-3 max-w-xs text-sm text-muted">{{ $t('landing.footer.tagline') }}</p>
          </div>
          <div class="grid grid-cols-2 gap-8 sm:grid-cols-4">
            <div>
              <p class="text-sm font-medium">{{ $t('landing.footer.product') }}</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><a href="#developers" class="hover:text-theme-text">{{ $t('landing.footer.developers') }}</a></li>
                <li><a href="#pricing" class="hover:text-theme-text">{{ $t('landing.footer.pricing') }}</a></li>
                <li><a href="#use-cases" class="hover:text-theme-text">{{ $t('landing.footer.use_cases') }}</a></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">{{ $t('landing.footer.use_cases') }}</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><router-link to="/cases/customer-support" class="hover:text-theme-text">{{ $t('landing.use_cases.customer_support.label') }}</router-link></li>
                <li><router-link to="/cases/sales" class="hover:text-theme-text">{{ $t('landing.use_cases.sales.label') }}</router-link></li>
                <li><router-link to="/cases/healthcare" class="hover:text-theme-text">{{ $t('landing.use_cases.healthcare.label') }}</router-link></li>
                <li><router-link to="/cases/education" class="hover:text-theme-text">{{ $t('landing.use_cases.education.label') }}</router-link></li>
                <li><router-link to="/cases/devtools" class="hover:text-theme-text">{{ $t('landing.use_cases.devtools.label') }}</router-link></li>
                <li><router-link to="/cases/e-commerce" class="hover:text-theme-text">{{ $t('landing.use_cases.ecommerce.label') }}</router-link></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">{{ $t('landing.footer.resources') }}</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><router-link to="/docs" class="hover:text-theme-text">{{ $t('landing.footer.docs') }}</router-link></li>
                <li><router-link to="/status" class="hover:text-theme-text">{{ $t('landing.footer.status') }}</router-link></li>
                <li><a href="https://github.com/actx0/Actx0/issues" class="hover:text-theme-text" target="_blank" rel="noopener">{{ $t('landing.footer.support') }}</a></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">{{ $t('landing.footer.company') }}</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><a href="mailto:hello@clivern.com" class="hover:text-theme-text">{{ $t('landing.footer.contact') }}</a></li>
                <li>
                  <button type="button" class="hover:text-theme-text" @click="openCookiePreferences">
                    {{ $t('cookies.settings_link') }}
                  </button>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="mt-10 border-t border-theme-border pt-6 text-center text-sm text-primary-500">
          {{ $t('common.copyright') }}
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { highlightCode } from '@/lib/code'
import { openCookiePreferences } from '@/lib/cookies'
import { loadUserFromStorage } from '@/utils/storage'
import LandingNav from '@/components/LandingNav.vue'

const { t } = useI18n()

const codeTab = ref('python')
const featureTab = ref('efficiency')
const copyLabel = ref(t('landing.copy'))
const rootRef = ref(null)

const codeTabs = ['python', 'curl', 'go']
const featureTabs = ['efficiency', 'visibility', 'control']

const efficiencyItems = [
  'landing.developers.efficiency.item_1',
  'landing.developers.efficiency.item_2',
  'landing.developers.efficiency.item_3',
]

const controlItems = [
  { titleKey: 'landing.developers.control.access_keys_title', textKey: 'landing.developers.control.access_keys_text' },
  { titleKey: 'landing.developers.control.prompt_versions_title', textKey: 'landing.developers.control.prompt_versions_text' },
  { titleKey: 'landing.developers.control.knowledge_title', textKey: 'landing.developers.control.knowledge_text' },
  { titleKey: 'landing.developers.control.sessions_title', textKey: 'landing.developers.control.sessions_text' },
]

const howSteps = [
  { titleKey: 'landing.how_it_works.add_title', textKey: 'landing.how_it_works.add_text' },
  { titleKey: 'landing.how_it_works.learn_title', textKey: 'landing.how_it_works.learn_text' },
  { titleKey: 'landing.how_it_works.retrieve_title', textKey: 'landing.how_it_works.retrieve_text' },
]

const stats = [
  { valueKey: 'landing.stats.latency_value', labelKey: 'landing.stats.latency_label' },
  { valueKey: 'landing.stats.token_value', labelKey: 'landing.stats.token_label' },
  { valueKey: 'landing.stats.infra_value', labelKey: 'landing.stats.infra_label' },
  { valueKey: 'landing.stats.billing_value', labelKey: 'landing.stats.billing_label' },
]

const useCaseCards = [
  {
    to: '/cases/customer-support',
    labelKey: 'landing.use_cases.customer_support.label',
    titleKey: 'landing.use_cases.customer_support.title',
    textKey: 'landing.use_cases.customer_support.text',
  },
  {
    to: '/cases/sales',
    labelKey: 'landing.use_cases.sales.label',
    titleKey: 'landing.use_cases.sales.title',
    textKey: 'landing.use_cases.sales.text',
  },
  {
    to: '/cases/healthcare',
    labelKey: 'landing.use_cases.healthcare.label',
    titleKey: 'landing.use_cases.healthcare.title',
    textKey: 'landing.use_cases.healthcare.text',
  },
  {
    to: '/cases/education',
    labelKey: 'landing.use_cases.education.label',
    titleKey: 'landing.use_cases.education.title',
    textKey: 'landing.use_cases.education.text',
  },
  {
    to: '/cases/devtools',
    labelKey: 'landing.use_cases.devtools.label',
    titleKey: 'landing.use_cases.devtools.title',
    textKey: 'landing.use_cases.devtools.text',
  },
  {
    to: '/cases/e-commerce',
    labelKey: 'landing.use_cases.ecommerce.label',
    titleKey: 'landing.use_cases.ecommerce.title',
    textKey: 'landing.use_cases.ecommerce.text',
  },
]

const enterpriseItems = [
  { titleKey: 'landing.enterprise.managed_title', textKey: 'landing.enterprise.managed_text' },
  { titleKey: 'landing.enterprise.auditable_title', textKey: 'landing.enterprise.auditable_text' },
  { titleKey: 'landing.enterprise.multi_tenant_title', textKey: 'landing.enterprise.multi_tenant_text' },
]

const plans = [
  {
    id: 'hobby',
    nameKey: 'landing.pricing.plans.hobby.name',
    price: '$0',
    descriptionKey: 'landing.pricing.plans.hobby.description',
    featureKeys: [
      'landing.pricing.plans.hobby.feature_1',
      'landing.pricing.plans.hobby.feature_2',
      'landing.pricing.plans.hobby.feature_3',
      'landing.pricing.plans.hobby.feature_4',
    ],
    ctaKey: 'landing.pricing.cta_free',
  },
  {
    id: 'starter',
    nameKey: 'landing.pricing.plans.starter.name',
    price: '$19',
    descriptionKey: 'landing.pricing.plans.starter.description',
    popular: true,
    featureKeys: [
      'landing.pricing.plans.starter.feature_1',
      'landing.pricing.plans.starter.feature_2',
      'landing.pricing.plans.starter.feature_3',
      'landing.pricing.plans.starter.feature_4',
    ],
    ctaKey: 'landing.pricing.cta_trial',
  },
  {
    id: 'growth',
    nameKey: 'landing.pricing.plans.growth.name',
    price: '$69',
    descriptionKey: 'landing.pricing.plans.growth.description',
    featureKeys: [
      'landing.pricing.plans.growth.feature_1',
      'landing.pricing.plans.growth.feature_2',
      'landing.pricing.plans.growth.feature_3',
      'landing.pricing.plans.growth.feature_4',
    ],
    ctaKey: 'landing.pricing.cta_trial',
  },
  {
    id: 'pro',
    nameKey: 'landing.pricing.plans.pro.name',
    price: '€219',
    descriptionKey: 'landing.pricing.plans.pro.description',
    featureKeys: [
      'landing.pricing.plans.pro.feature_1',
      'landing.pricing.plans.pro.feature_2',
      'landing.pricing.plans.pro.feature_3',
      'landing.pricing.plans.pro.feature_4',
    ],
    ctaKey: 'landing.pricing.cta_trial',
  },
]

const highlightedCode = computed(() => highlightCode(codeSamples[codeTab.value], codeTab.value))

const startPath = computed(() => {
  if (loadUserFromStorage()) return '/select-workspace'
  return '/login'
})

const codeSamples = {
  python: `# pip install actx0
from actx0 import MemoryClient

client = MemoryClient(
    api_key="your-access-key",
    workspace="your-workspace",
)

# Store a memory from conversation
client.add(
    messages=[
        {"role": "user", "content": "I'm allergic to nuts."},
        {"role": "assistant", "content": "Noted. I'll remember that."},
    ],
    agent_id="support-bot",
)

# Retrieve relevant context at inference time
results = client.search(
    "What are my dietary restrictions?",
    agent_id="support-bot",
    limit=5,
)
print(results)`,
  curl: `# Search agent memory
curl -X POST "https://your-instance/api/v1/memory/search" \\
  -H "Authorization: Bearer YOUR_ACCESS_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "query": "What are my dietary restrictions?",
    "agent_id": "support-bot",
    "limit": 5
  }'`,
  go: `// go get github.com/actx0/Actx0/sdk/go
client, _ := actx0.New(actx0.Config{
    APIKey:    "your-access-key",
    Workspace: "your-workspace",
})

_, _ = client.Add(ctx, actx0.AddRequest{
    AgentID: "support-bot",
    Messages: []actx0.Message{
        {Role: "user", Content: "I'm allergic to nuts."},
    },
})

results, _ := client.Search(ctx, actx0.SearchRequest{
    Query:   "What are my dietary restrictions?",
    AgentID: "support-bot",
})`,
}

let revealObserver

async function copyCode() {
  await navigator.clipboard.writeText(codeSamples[codeTab.value])
  copyLabel.value = t('landing.copied')
  setTimeout(() => {
    copyLabel.value = t('landing.copy')
  }, 1600)
}

onMounted(() => {
  const items = rootRef.value?.querySelectorAll('.reveal') ?? []
  revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
          revealObserver.unobserve(entry.target)
        }
      })
    },
    { threshold: 0.12, rootMargin: '0px 0px -40px 0px' }
  )
  items.forEach((item) => revealObserver.observe(item))
})

onUnmounted(() => {
  revealObserver?.disconnect()
})
</script>

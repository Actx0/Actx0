<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('audits_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('audits_page.subtitle') }}</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4 flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('audits_page.event_log') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('audits_page.event_log_desc') }}</p>
          </div>
          <div class="flex flex-wrap items-center gap-2 shrink-0">
            <select v-model="categoryFilter" class="input-field text-sm py-1.5 w-auto min-w-[10rem]">
              <option value="">{{ $t('audits_page.all_categories') }}</option>
              <option v-for="cat in categoryOptions" :key="cat" :value="cat">{{ categoryLabel(cat) }}</option>
            </select>
            <select v-model="actionFilter" class="input-field text-sm py-1.5 w-auto min-w-[10rem]">
              <option value="">{{ $t('audits_page.all_actions') }}</option>
              <option v-for="action in actionOptions" :key="action" :value="action">{{ action }}</option>
            </select>
          </div>
        </div>

        <div v-if="loading" class="p-12 text-center">
          <svg class="animate-spin h-8 w-8 mx-auto text-theme-text" fill="none" viewBox="0 0 24 24" aria-hidden="true">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-theme-textLight mt-3">{{ $t('common.loading') }}</p>
        </div>

        <div v-else-if="audits.length === 0" class="p-12 text-center text-sm text-theme-textLight">
          {{ $t('audits_page.no_events') }}
        </div>

        <div v-else-if="filteredAudits.length === 0" class="p-12 text-center text-sm text-theme-textLight">
          {{ $t('audits_page.no_matching_events') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.time') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.action') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.actor') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.resource') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.source') }}</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('audits_page.details') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <template v-for="event in filteredAudits" :key="event.id">
                <tr class="hover:bg-theme-hover">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-theme-text tabular-nums">{{ formatAuditTime(event.createdAt, locale) }}</div>
                    <div class="text-xs text-theme-textLight">{{ formatAuditDate(event.createdAt, locale) }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span :class="['inline-flex rounded-md px-2 py-0.5 text-xs font-mono', actionPillClass(event.action)]">
                      {{ event.action }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <template v-if="event.actorName">
                      <div class="text-sm font-medium text-theme-text">{{ event.actorName }}</div>
                      <div class="text-xs font-mono text-theme-textLight">{{ actorSecondary(event) }}</div>
                    </template>
                    <span v-else class="text-sm text-theme-textLight">{{ $t('common.em_dash') }}</span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <template v-if="event.resourceType || resourceTitle(event)">
                      <div class="text-sm text-theme-text">{{ resourceTitle(event) }}</div>
                      <div v-if="event.resourceId" class="text-xs font-mono text-theme-textLight">{{ shortId(event.resourceId) }}</div>
                    </template>
                    <span v-else class="text-sm text-theme-textLight">{{ $t('common.em_dash') }}</span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-mono text-theme-text">{{ event.ipAddress || $t('common.em_dash') }}</div>
                    <div class="text-xs text-theme-textLight max-w-[14rem] truncate" :title="event.userAgent || undefined">
                      {{ truncateAgent(event.userAgent) }}
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      type="button"
                      class="text-primary-600 hover:text-primary-700 hover:underline"
                      :aria-expanded="isDetailOpen(event.id)"
                      @click="toggleDetail(event.id)"
                    >
                      {{ isDetailOpen(event.id) ? $t('audits_page.hide') : $t('audits_page.view') }}
                    </button>
                  </td>
                </tr>
                <tr v-if="isDetailOpen(event.id)" class="bg-theme-hover/50">
                  <td colspan="6" class="px-6 py-4">
                    <dl class="grid gap-x-8 gap-y-2 text-sm sm:grid-cols-2 lg:grid-cols-4">
                      <div>
                        <dt class="text-xs font-medium uppercase tracking-wider text-theme-textLight">{{ $t('audits_page.event_id') }}</dt>
                        <dd class="mt-1 font-mono text-xs text-theme-text break-all">{{ event.id }}</dd>
                      </div>
                      <div>
                        <dt class="text-xs font-medium uppercase tracking-wider text-theme-textLight">{{ $t('audits_page.category') }}</dt>
                        <dd class="mt-1 text-theme-text">{{ categoryLabel(categoryOf(event.action)) }}</dd>
                      </div>
                      <div>
                        <dt class="text-xs font-medium uppercase tracking-wider text-theme-textLight">{{ $t('audits_page.resource_type') }}</dt>
                        <dd class="mt-1 text-theme-text capitalize">{{ event.resourceType || $t('common.em_dash') }}</dd>
                      </div>
                      <div>
                        <dt class="text-xs font-medium uppercase tracking-wider text-theme-textLight">{{ $t('audits_page.source') }}</dt>
                        <dd class="mt-1 text-theme-text break-all">{{ event.userAgent || $t('common.em_dash') }}</dd>
                      </div>
                    </dl>
                    <div v-if="hasMetadata(event.meta)" class="mt-4">
                      <p class="text-xs font-medium uppercase tracking-wider text-theme-textLight">{{ $t('audits_page.metadata') }}</p>
                      <pre class="mt-1 overflow-x-auto rounded-lg border border-theme-border bg-white px-3 py-2 text-xs font-mono text-theme-text">{{ formatMetadata(event.meta) }}</pre>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>

        <div v-if="total > 0" class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between">
          <div class="text-sm text-theme-textLight">
            {{ $t('audits_page.showing', { from: offset + 1, to: Math.min(offset + limit, total), total }) }}
          </div>
          <div class="flex items-center gap-2">
            <button type="button" @click="goToPage(offset - limit)" :disabled="offset === 0" class="btn-secondary text-sm disabled:opacity-50">{{ $t('audits_page.previous') }}</button>
            <button type="button" @click="goToPage(offset + limit)" :disabled="offset + limit >= total" class="btn-secondary text-sm disabled:opacity-50">{{ $t('audits_page.next') }}</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { audit_api } from '@/api'
import { useWorkspaceContext } from '@/lib/permission'
import { useUrlPagination } from '@/lib/pagination'
import { hasMetadata, formatMetadata } from '@/lib/agent'
import { formatAuditDate, formatAuditTime } from '@/utils/dt'

const { t, te, locale } = useI18n()
const route = useRoute()
const { currentWorkspace } = useWorkspaceContext()
const { offset, limit, goToOffset } = useUrlPagination({ limit: 50 })

const audits = ref([])
const total = ref(0)
const loading = ref(false)
const errorMessage = ref(null)
const categoryFilter = ref('')
const actionFilter = ref('')
const openDetailIds = ref({})

const PILL_CLASSES = {
  memory: 'bg-emerald-50 text-emerald-800',
  message: 'bg-emerald-50 text-emerald-800',
  session: 'bg-emerald-50 text-emerald-800',
  prompt: 'bg-violet-50 text-violet-800',
  document: 'bg-violet-50 text-violet-800',
  access_key: 'bg-amber-50 text-amber-900',
  member: 'bg-blue-50 text-blue-800',
  invite: 'bg-blue-50 text-blue-800',
  agent: 'bg-teal-50 text-teal-800',
  workspace: 'bg-slate-100 text-slate-700',
  billing: 'bg-slate-100 text-slate-700',
}

function categoryOf(action) {
  return action.split('.')[0]
}

function categoryLabel(category) {
  const key = `audits_page.cat_${category}`
  return te(key) ? t(key) : category
}

function shortId(id) {
  return String(id).slice(0, 8)
}

function actionPillClass(action) {
  return PILL_CLASSES[categoryOf(action)] || 'bg-slate-100 text-slate-700'
}

function actorSecondary(event) {
  const id = event.actorId ? shortId(event.actorId) : ''
  if (event.actorType === 'access_key') {
    return id ? `access-key:${id}` : t('audits_page.actor_access_key')
  }
  return id || t('audits_page.actor_user')
}

function parseMeta(meta) {
  if (!meta) return null
  if (typeof meta === 'object') return meta
  return JSON.parse(meta)
}

function resourceTitle(event) {
  const meta = parseMeta(event.meta)
  if (meta) {
    const label = meta.name || meta.title || meta.handle || meta.email
    if (label) return label
  }
  if (event.resourceType) {
    return event.resourceType.replaceAll('_', ' ')
  }
  return ''
}

function truncateAgent(ua) {
  if (!ua) return t('common.em_dash')
  if (ua.length <= 28) return ua
  return `${ua.slice(0, 28)}…`
}

const categoryOptions = computed(() => {
  const set = new Set(audits.value.map((e) => categoryOf(e.action)))
  return [...set].sort()
})

const actionOptions = computed(() => {
  const list = categoryFilter.value
    ? audits.value.filter((e) => categoryOf(e.action) === categoryFilter.value)
    : audits.value
  const set = new Set(list.map((e) => e.action))
  return [...set].sort()
})

const filteredAudits = computed(() => {
  return audits.value.filter((event) => {
    if (categoryFilter.value && categoryOf(event.action) !== categoryFilter.value) return false
    if (actionFilter.value && event.action !== actionFilter.value) return false
    return true
  })
})

watch(categoryFilter, () => {
  if (actionFilter.value && !actionOptions.value.includes(actionFilter.value)) {
    actionFilter.value = ''
  }
})

function isDetailOpen(eventId) {
  return Boolean(openDetailIds.value[eventId])
}

function toggleDetail(eventId) {
  openDetailIds.value = {
    ...openDetailIds.value,
    [eventId]: !openDetailIds.value[eventId],
  }
}

async function loadPage() {
  loading.value = true
  errorMessage.value = null

  try {
    const res = await audit_api.list(currentWorkspace.id, { limit, offset: offset.value })
    audits.value = res.data.audits
    total.value = res.data._meta.total
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('audits_page.failed_load')
    audits.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function goToPage(newOffset) {
  goToOffset(newOffset)
}

watch(() => route.query.page, loadPage)
onMounted(loadPage)
</script>

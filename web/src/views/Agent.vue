<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-6">
        <router-link to="/agents" class="text-sm text-primary-600 hover:text-primary-700 hover:underline">
          ← {{ $t('agents_page.back_to_agents') }}
        </router-link>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div v-if="loading" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('common.loading') }}
      </div>

      <div v-else-if="!agent" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('agents_page.not_found') }}
      </div>

      <template v-else>
        <div class="mb-8 rounded-xl border border-theme-border bg-white shadow-sm p-6">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div class="min-w-0 flex-1">
              <h1 class="text-2xl font-semibold text-theme-text">{{ agent.name }}</h1>
              <p class="mt-2 text-sm text-theme-textLight">{{ agent.description }}</p>
              <p class="mt-2 text-xs font-mono text-theme-textLight break-all">{{ agent.id }}</p>
              <dl class="mt-4 grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
                <div>
                  <dt class="text-theme-textLight">{{ $t('agents_page.created_at') }}</dt>
                  <dd class="mt-0.5 text-theme-text">{{ formatDate(agent.createdAt) }}</dd>
                </div>
                <div>
                  <dt class="text-theme-textLight">{{ $t('agents_page.memory_pipeline') }}</dt>
                  <dd class="mt-0.5 text-theme-text">
                    {{ agent.configs.memoryPipeline ? $t('agents_page.memory_pipeline_on') : $t('agents_page.memory_pipeline_off') }}
                  </dd>
                </div>
              </dl>
            </div>
            <button v-if="canEdit" type="button" class="btn-secondary shrink-0" @click="openEditModal">
              {{ $t('common.edit') }}
            </button>
          </div>
        </div>

        <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6 py-4">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.sessions_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.sessions_desc') }}</p>
          </div>

          <div v-if="sessionsLoading" class="p-12 text-center text-sm text-theme-textLight">
            {{ $t('common.loading') }}
          </div>

          <div v-else-if="sessionsError" class="p-6 text-center text-sm text-red-800">
            {{ sessionsError }}
          </div>

          <div v-else-if="sessions.length === 0" class="p-12 text-center text-sm text-theme-textLight">
            {{ $t('agents_page.no_sessions') }}
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-theme-border">
              <thead class="bg-theme-hover">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.session') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.status') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.labels') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.created_at') }}</th>
                  <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.actions') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-theme-border bg-white">
                <tr v-for="session in sessions" :key="session.id" class="hover:bg-theme-hover/60">
                  <td class="px-6 py-4">
                    <router-link :to="sessionPath(agentId, session.id)" class="text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">
                      {{ session.title || $t('agents_page.untitled_session') }}
                    </router-link>
                    <p class="mt-0.5 text-xs font-mono text-theme-textLight break-all">{{ session.externalId }}</p>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium" :class="sessionStatusClass(session.status)">
                      {{ sessionStatusLabel(session.status) }}
                    </span>
                  </td>
                  <td class="px-6 py-4">
                    <div v-if="formatLabels(session.labels).length" class="flex flex-wrap gap-1">
                      <span
                        v-for="label in formatLabels(session.labels)"
                        :key="label"
                        class="inline-flex items-center rounded-full bg-theme-hover px-2 py-0.5 text-xs font-mono text-theme-text"
                      >
                        {{ label }}
                      </span>
                    </div>
                    <span v-else class="text-sm text-theme-textLight">{{ $t('common.em_dash') }}</span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ formatDateTime(session.createdAt) }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <span class="inline-flex items-center gap-3">
                      <router-link :to="sessionPath(agentId, session.id)" class="text-primary-600 hover:text-primary-700 hover:underline">
                        {{ $t('agents_page.view') }}
                      </router-link>
                      <button
                        v-if="canEdit"
                        type="button"
                        class="text-red-600 hover:text-red-700 hover:underline disabled:opacity-50"
                        :disabled="deletingSessionId === session.id"
                        @click="openDeleteSessionModal(session)"
                      >
                        {{ deletingSessionId === session.id ? $t('agents_page.deleting') : $t('agents_page.delete') }}
                      </button>
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="sessionsTotal > 0" class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between">
            <div class="text-sm text-theme-textLight">
              {{ $t('agents_page.showing', { from: sessionsOffset + 1, to: Math.min(sessionsOffset + sessionsLimit, sessionsTotal), total: sessionsTotal }) }}
            </div>
            <div class="flex items-center gap-2">
              <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="sessionsOffset === 0" @click="goToSessionsPage(sessionsOffset - sessionsLimit)">
                {{ $t('agents_page.previous') }}
              </button>
              <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="sessionsOffset + sessionsLimit >= sessionsTotal" @click="goToSessionsPage(sessionsOffset + sessionsLimit)">
                {{ $t('agents_page.next') }}
              </button>
            </div>
          </div>
        </section>
      </template>
    </main>

    <!-- Delete session modal -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-if="deleteSessionTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="!deletingSession && closeDeleteSessionModal()">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.delete_session_title') }}</h2>
            <div v-if="deleteSessionError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ deleteSessionError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">
              {{ $t('agents_page.delete_session_confirm', { name: deleteSessionTarget.title || $t('agents_page.untitled_session') }) }}
            </p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deletingSession" @click="closeDeleteSessionModal">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deletingSession" @click="confirmDeleteSession">
                {{ deletingSession ? $t('agents_page.deleting') : $t('agents_page.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Edit modal -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="!saveLoading && (showEditModal = false)">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.edit_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.edit_desc') }}</p>
            <div v-if="editError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ editError }}</div>
            <form class="mt-4 space-y-4" @submit.prevent="handleSave">
              <div>
                <label for="edit-agent-name" class="form-label">{{ $t('agents_page.name') }}</label>
                <input id="edit-agent-name" v-model="editForm.name" type="text" required class="input-field" :placeholder="$t('agents_page.name_placeholder')" :disabled="saveLoading" />
              </div>
              <div>
                <label for="edit-agent-description" class="form-label">{{ $t('agents_page.description') }}</label>
                <textarea id="edit-agent-description" v-model="editForm.description" rows="2" required class="input-field" :placeholder="$t('agents_page.description_placeholder')" :disabled="saveLoading" />
              </div>
              <div class="rounded-lg border border-theme-border bg-theme-hover/40 p-4">
                <label class="flex items-start gap-3 cursor-pointer">
                  <input v-model="editForm.memoryPipeline" type="checkbox" class="mt-1 h-4 w-4 rounded border-theme-border text-primary-800 focus:ring-primary-500" :disabled="saveLoading || isHobbyPlan" />
                  <span>
                    <span class="block text-sm font-medium text-theme-text">{{ $t('agents_page.memory_pipeline') }}</span>
                    <span class="mt-0.5 block text-xs text-theme-textLight">{{ $t('agents_page.memory_pipeline_hint') }}</span>
                  </span>
                </label>
              </div>

              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="saveLoading" @click="showEditModal = false">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="saveLoading">
                  {{ saveLoading ? $t('agents_page.saving') : $t('common.save') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { useWorkspaceContext } from '@/lib/permission'
import { agent_api, session_api, billing_api } from '@/api'
import { sessionPath, formatLabels } from '@/lib/agent'
import { formatDate, formatDateTime } from '@/utils/dt'

const { t } = useI18n()
const route = useRoute()
const { currentWorkspace, canManage: canEdit } = useWorkspaceContext()
const { offset: sessionsOffset, limit: sessionsLimit, goToOffset: goToSessionsOffset, goToPage: setSessionsPage } = useUrlPagination({ limit: 50 })

const agent = ref(null)
const loading = ref(false)
const errorMessage = ref(null)

const sessions = ref([])
const sessionsLoading = ref(false)
const sessionsError = ref(null)
const sessionsTotal = ref(0)

const showEditModal = ref(false)
const editError = ref(null)
const saveLoading = ref(false)
const editForm = reactive({ name: '', description: '', memoryPipeline: false })
const isHobbyPlan = ref(true)

const deleteSessionTarget = ref(null)
const deleteSessionError = ref(null)
const deletingSession = ref(false)
const deletingSessionId = ref(null)

const agentId = computed(() => route.params.agentId)

async function loadAgent() {
  loading.value = true
  errorMessage.value = null
  agent.value = null
  sessions.value = []
  sessionsTotal.value = 0
  try {
    const res = await agent_api.get(currentWorkspace.id, agentId.value)
    agent.value = res.data
    await loadSessions()
  } catch (err) {
    if (err.response?.status === 404) {
      agent.value = null
    } else {
      errorMessage.value = err.response?.data?.errorMessage || t('agents_page.failed_load')
    }
  } finally {
    loading.value = false
  }
}

async function loadSessions() {
  sessionsLoading.value = true
  sessionsError.value = null
  try {
    const res = await session_api.list(currentWorkspace.id, agentId.value, {
      limit: sessionsLimit,
      offset: sessionsOffset.value,
    })
    sessions.value = res.data.sessions
    sessionsTotal.value = res.data._meta.total
  } catch (err) {
    sessions.value = []
    sessionsTotal.value = 0
    sessionsError.value = err.response?.data?.errorMessage || t('agents_page.failed_load_sessions')
  } finally {
    sessionsLoading.value = false
  }
}

function goToSessionsPage(offset) {
  goToSessionsOffset(offset)
}

function sessionStatusLabel(status) {
  if (status === 'closed') return t('agents_page.session_closed')
  return t('agents_page.session_active')
}

function sessionStatusClass(status) {
  if (status === 'closed') return 'bg-theme-hover text-theme-textLight'
  return 'bg-sky-100 text-sky-800'
}

async function openEditModal() {
  editForm.name = agent.value.name
  editForm.description = agent.value.description
  editForm.memoryPipeline = agent.value.configs.memoryPipeline
  editError.value = null
  showEditModal.value = true

  const res = await billing_api.status(currentWorkspace.id)
  isHobbyPlan.value = res.data.plan === 'hobby'
  if (isHobbyPlan.value) {
    editForm.memoryPipeline = false
  }
}

function openDeleteSessionModal(session) {
  deleteSessionTarget.value = session
  deleteSessionError.value = null
}

function closeDeleteSessionModal() {
  deleteSessionTarget.value = null
  deleteSessionError.value = null
}

async function handleSave() {
  const name = editForm.name.trim()
  const description = editForm.description.trim()
  if (!name) {
    editError.value = t('agents_page.name_required')
    return
  }
  if (!description) {
    editError.value = t('agents_page.description_required')
    return
  }

  saveLoading.value = true
  editError.value = null
  try {
    const res = await agent_api.update(currentWorkspace.id, agentId.value, {
      name,
      description,
      configs: {
        memoryPipeline: editForm.memoryPipeline,
      },
    })
    agent.value = res.data
    showEditModal.value = false
    showFlash(t('common.saved'))
  } catch (err) {
    editError.value = err.response?.data?.errorMessage || t('agents_page.failed_save')
  } finally {
    saveLoading.value = false
  }
}

async function confirmDeleteSession() {
  const session = deleteSessionTarget.value

  deletingSession.value = true
  deletingSessionId.value = session.id
  deleteSessionError.value = null

  try {
    await session_api.delete(currentWorkspace.id, agentId.value, { id: session.externalId })
    deleteSessionTarget.value = null

    if (sessionsOffset.value >= sessionsTotal.value - 1 && sessionsOffset.value > 0) {
      goToSessionsOffset(Math.max(0, sessionsOffset.value - sessionsLimit))
    }

    await loadSessions()
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteSessionError.value = err.response?.data?.errorMessage || t('agents_page.failed_delete_session')
  } finally {
    deletingSession.value = false
    deletingSessionId.value = null
  }
}

watch(() => route.query.page, () => {
  if (agent.value) loadSessions()
})

watch(agentId, async (newId, oldId) => {
  if (oldId && newId !== oldId) {
    await setSessionsPage(1)
  }
  await loadAgent()
})

onMounted(loadAgent)
</script>

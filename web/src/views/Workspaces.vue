<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('workspaces.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('workspaces_page.subtitle') }}</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="flex justify-end border-b border-theme-border bg-white px-6 py-4">
          <button
            type="button"
            class="btn-primary inline-flex items-center justify-center gap-2"
            @click="openCreateModal"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ $t('workspaces_page.create') }}
          </button>
        </div>

        <div v-if="loading" class="p-12 text-center">
          <svg class="animate-spin h-8 w-8 mx-auto text-theme-text" fill="none" viewBox="0 0 24 24" aria-hidden="true">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-theme-textLight mt-3">{{ $t('workspaces_page.loading') }}</p>
        </div>

        <div v-else-if="workspaces.length === 0" class="p-12 text-center">
          <div class="flex justify-center mb-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-primary-100">
              <svg class="h-6 w-6 text-theme-text" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
          </div>
          <p class="text-theme-textLight">{{ $t('workspaces_page.no_workspaces') }}</p>
          <p class="text-sm text-theme-textLight mt-1">{{ $t('workspaces_page.no_workspaces_hint') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">
                  {{ $t('workspaces_page.name') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">
                  {{ $t('workspaces_page.created') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">
                  {{ $t('workspaces_page.members') }}
                </th>
                <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">
                  {{ $t('workspaces_page.actions') }}
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <tr
                v-for="workspace in workspaces"
                :key="workspace.id"
                class="hover:bg-theme-hover"
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center gap-3">
                    <span class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-primary-100 text-theme-text">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                      </svg>
                    </span>
                    <span class="text-sm font-medium text-theme-text">{{ workspace.name }}</span>
                    <span
                      v-if="currentWorkspace.id === workspace.id"
                      class="text-xs text-theme-textLight"
                    >
                      {{ $t('workspaces_page.current') }}
                    </span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">
                  {{ formatDate(workspace.createdAt) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text">
                  {{ workspace.membersCount }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <span class="inline-flex items-center gap-3">
                    <button
                      type="button"
                      class="text-emerald-600 hover:text-emerald-700 hover:underline disabled:opacity-50 text-sm"
                      :disabled="switchLoadingId === workspace.id"
                      :title="$t('workspaces_page.switch_into_title')"
                      @click="switchToWorkspace(workspace)"
                    >
                      {{ switchLoadingId === workspace.id ? $t('workspaces_page.switching') : $t('workspaces_page.switch_into') }}
                    </button>
                    <button
                      v-if="canManageWorkspace(currentUser, workspace)"
                      type="button"
                      class="text-primary-600 hover:text-primary-700 hover:underline disabled:opacity-50"
                      :disabled="updateLoadingId === workspace.id"
                      :title="$t('workspaces_page.edit_title')"
                      @click="openEditModal(workspace)"
                    >
                      {{ updateLoadingId === workspace.id ? $t('workspaces_page.saving') : $t('workspaces_page.edit') }}
                    </button>
                    <button
                      v-if="canManageWorkspace(currentUser, workspace)"
                      type="button"
                      class="text-red-600 hover:text-red-700 hover:underline disabled:opacity-50"
                      :disabled="deletingId === workspace.id"
                      :title="$t('workspaces_page.delete_title')"
                      @click="openDeleteModal(workspace)"
                    >
                      {{ deletingId === workspace.id ? $t('workspaces_page.deleting') : $t('workspaces_page.delete') }}
                    </button>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div
          v-if="total > 0"
          class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between"
        >
          <div class="text-sm text-theme-textLight">
            {{ $t('workspaces_page.showing', { from: offset + 1, to: Math.min(offset + limit, total), total }) }}
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="offset === 0"
              @click="goToPage(offset - limit)"
            >
              {{ $t('workspaces_page.previous') }}
            </button>
            <button
              type="button"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="offset + limit >= total"
              @click="goToPage(offset + limit)"
            >
              {{ $t('workspaces_page.next') }}
            </button>
          </div>
        </div>
      </div>
    </main>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showCreateModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="create-workspace-title"
          @click.self="closeCreateModal"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="create-workspace-title" class="text-lg font-semibold text-theme-text">
              {{ $t('workspaces_page.create_title') }}
            </h2>
            <p class="mt-1 text-sm text-theme-textLight">
              {{ $t('workspaces_page.create_desc') }}
            </p>
            <div v-if="createError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800" role="alert">
              {{ createError }}
            </div>
            <form class="mt-4 space-y-4" @submit.prevent="handleCreate">
              <div>
                <label for="create-workspace-name" class="form-label">{{ $t('workspaces_page.name') }}</label>
                <input
                  id="create-workspace-name"
                  v-model="createForm.name"
                  type="text"
                  required
                  class="input-field"
                  :placeholder="$t('workspaces_page.name_placeholder')"
                  :disabled="createLoading"
                  autofocus
                />
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="createLoading" @click="closeCreateModal">
                  {{ $t('common.cancel') }}
                </button>
                <button
                  type="submit"
                  class="btn-primary"
                  :disabled="createLoading || !createForm.name.trim()"
                >
                  {{ createLoading ? $t('workspaces_page.creating') : $t('common.create') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="editTarget"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="edit-workspace-title"
          @click.self="closeEditModal"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="edit-workspace-title" class="text-lg font-semibold text-theme-text">
              {{ $t('workspaces_page.edit_title') }}
            </h2>
            <p class="mt-1 text-sm text-theme-textLight">
              {{ $t('workspaces_page.edit_desc') }}
            </p>
            <div v-if="editError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800" role="alert">
              {{ editError }}
            </div>
            <form class="mt-4 space-y-4" @submit.prevent="handleUpdate">
              <div>
                <label for="edit-workspace-name" class="form-label">{{ $t('workspaces_page.name') }}</label>
                <input
                  id="edit-workspace-name"
                  v-model="editForm.name"
                  type="text"
                  required
                  class="input-field"
                  :placeholder="$t('workspaces_page.name_placeholder')"
                  :disabled="updateLoadingId"
                  autofocus
                />
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="updateLoadingId" @click="closeEditModal">
                  {{ $t('common.cancel') }}
                </button>
                <button
                  type="submit"
                  class="btn-primary"
                  :disabled="updateLoadingId || !editForm.name.trim() || editForm.name.trim() === editTarget.name"
                >
                  {{ updateLoadingId ? $t('workspaces_page.saving') : $t('common.save') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="deleteTarget"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="delete-workspace-title"
          @click.self="closeDeleteModal"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="delete-workspace-title" class="text-lg font-semibold text-theme-text">
              {{ $t('workspaces_page.delete_title') }}
            </h2>
            <div v-if="deleteError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800" role="alert">
              {{ deleteError }}
            </div>
            <p class="mt-2 text-sm text-theme-textLight">
              {{ $t('workspaces_page.delete_confirm', { name: deleteTarget.name }) }}
            </p>
            <div class="mt-4">
              <label for="delete-confirm-name" class="form-label">
                {{ $t('workspaces_page.type_to_confirm', { name: deleteTarget.name }) }}
              </label>
              <input
                id="delete-confirm-name"
                v-model="deleteForm.confirmName"
                type="text"
                class="input-field"
                :placeholder="deleteTarget.name"
                autocomplete="off"
                :disabled="deleting"
              />
            </div>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deleting" @click="closeDeleteModal">
                {{ $t('common.cancel') }}
              </button>
              <button
                type="button"
                class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500"
                :disabled="deleting || deleteForm.confirmName.trim() !== deleteTarget.name"
                @click="confirmDelete"
              >
                {{ deleting ? $t('workspaces_page.deleting') : $t('workspaces_page.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { workspace_api } from '@/api'
import { saveWorkspaceToStorage, removeWorkspaceFromStorage } from '@/utils/storage'
import { canManageWorkspace, useWorkspaceContext } from '@/lib/permission'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { currentUser, currentWorkspace } = useWorkspaceContext()
const { offset, limit, goToOffset, goToPage: setPage } = useUrlPagination({ limit: 50 })

const workspaces = ref([])
const loading = ref(false)
const total = ref(0)
const errorMessage = ref(null)

const showCreateModal = ref(false)
const createLoading = ref(false)
const createError = ref(null)
const createForm = reactive({ name: '' })

const editTarget = ref(null)
const updateLoadingId = ref(null)
const editError = ref(null)
const editForm = reactive({ name: '' })

const deleteTarget = ref(null)
const deleting = ref(false)
const deletingId = ref(null)
const deleteError = ref(null)
const deleteForm = reactive({ confirmName: '' })

const switchLoadingId = ref(null)

function formatDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString(
      'en-US',
      { month: 'short', day: 'numeric', year: 'numeric' }
    )
  } catch {
    return iso
  }
}

async function loadPage() {
  loading.value = true
  errorMessage.value = null

  try {
    const res = await workspace_api.list({
      limit,
      offset: offset.value
    })
    workspaces.value = res.data.workspaces
    total.value = res.data._meta.total
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('workspaces_page.failed_load')
    workspaces.value = []
  } finally {
    loading.value = false
  }
}

function goToPage(newOffset) {
  if (newOffset < 0) return
  if (newOffset >= total.value) return
  goToOffset(newOffset)
}

async function switchToWorkspace(workspace) {
  switchLoadingId.value = workspace.id
  errorMessage.value = null

  try {
    const res = await workspace_api.get(workspace.id)
    saveWorkspaceToStorage({ ...workspace, ...res.data, role: res.data.role })
    router.push('/dashboard')
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('workspaces_page.failed_switch')
  } finally {
    switchLoadingId.value = null
  }
}

function openCreateModal() {
  createForm.name = ''
  createError.value = null
  showCreateModal.value = true
}

function closeCreateModal() {
  if (createLoading.value) return
  showCreateModal.value = false
  createError.value = null
  createForm.name = ''
}

async function handleCreate() {
  const name = createForm.name.trim()
  if (!name) return

  createLoading.value = true
  errorMessage.value = null
  createError.value = null

  try {
    await workspace_api.create({ name })
    showCreateModal.value = false
    createForm.name = ''
    setPage(1)
    await loadPage()
    showFlash(t('common.created'))
  } catch (err) {
    createError.value = err.response?.data?.errorMessage || t('workspaces_page.failed_create')
  } finally {
    createLoading.value = false
  }
}

function openEditModal(workspace) {
  editTarget.value = workspace
  editForm.name = workspace.name
  editError.value = null
}

function closeEditModal() {
  if (updateLoadingId.value) return
  editTarget.value = null
  editForm.name = ''
  editError.value = null
}

async function handleUpdate() {
  const workspace = editTarget.value
  const name = editForm.name.trim()

  if (!workspace || !name) return
  updateLoadingId.value = workspace.id
  errorMessage.value = null
  editError.value = null

  try {
    const res = await workspace_api.update(workspace.id, { name })
    if (currentWorkspace.id === workspace.id) {
      saveWorkspaceToStorage({ ...currentWorkspace, ...res.data })
    }
    editTarget.value = null
    editForm.name = ''
    editError.value = null
    await loadPage()
    showFlash(t('common.saved'))
  } catch (err) {
    editError.value = err.response?.data?.errorMessage || t('workspaces_page.failed_update')
  } finally {
    updateLoadingId.value = null
  }
}

function openDeleteModal(workspace) {
  deleteTarget.value = workspace
  deleteForm.confirmName = ''
  deleteError.value = null
}

function closeDeleteModal() {
  if (!deleting.value) {
    deleteTarget.value = null
    deleteForm.confirmName = ''
    deleteError.value = null
  }
}

async function confirmDelete() {
  const workspace = deleteTarget.value
  if (!workspace) return
  deleting.value = true
  deletingId.value = workspace.id
  errorMessage.value = null
  deleteError.value = null

  try {
    await workspace_api.delete(workspace.id)
    if (currentWorkspace.id === workspace.id) {
      removeWorkspaceFromStorage()
    }
    deleteTarget.value = null
    deleteForm.confirmName = ''
    deleteError.value = null
    await loadPage()
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteError.value = err.response?.data?.errorMessage || t('workspaces_page.failed_delete')
  } finally {
    deleting.value = false
    deletingId.value = null
  }
}

watch(() => route.query.page, loadPage)

onMounted(loadPage)
</script>

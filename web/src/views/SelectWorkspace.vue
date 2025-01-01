<template>
  <div class="min-h-screen bg-theme-bg flex flex-col items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
      <div class="flex justify-center mb-6">
        <div class="flex h-14 w-14 items-center justify-center rounded-full bg-primary-200">
          <svg class="h-8 w-8 text-theme-text" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
        </div>
      </div>

      <h1 class="text-xl font-semibold text-theme-text text-center">
        {{ workspaces.length ? $t('select_workspace.select_title') : $t('select_workspace.create_title') }}
      </h1>
      <p class="mt-2 text-sm text-theme-textLight text-center">
        {{ workspaces.length ? $t('select_workspace.select_description') : $t('select_workspace.no_workspaces') }}
      </p>

      <div v-if="workspaces.length" class="mt-6 space-y-3">
        <div
          v-for="workspace in workspaces"
          :key="workspace.id"
          class="flex items-center gap-2 rounded-lg border border-theme-border bg-white shadow-sm overflow-hidden"
        >
          <button
            type="button"
            class="flex-1 flex items-center justify-between px-4 py-3 text-left hover:bg-theme-hover transition-colors focus:outline-none focus:ring-2 focus:ring-primary-800 focus:ring-inset min-w-0"
            @click="selectWorkspace(workspace)"
          >
            <div class="flex min-w-0 flex-1 items-center gap-3">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-theme-text truncate">{{ workspace.name }}</p>
                <p class="text-xs text-theme-textLight mt-0.5">{{ formatDate(workspace.createdAt) }}</p>
              </div>
              <span
                v-if="workspace.status === 'active' || storedWorkspace?.id === workspace.id"
                class="inline-flex flex-shrink-0 rounded-full bg-primary-200 px-2.5 py-0.5 text-xs font-medium text-theme-text"
              >
                {{ $t('select_workspace.active') }}
              </span>
            </div>
            <svg class="ml-3 h-5 w-5 flex-shrink-0 text-theme-textLight" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>

      <div class="mt-6">
        <button
          type="button"
          :class="[
            'w-full flex items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-primary-800 focus:ring-offset-1 transition-colors',
            workspaces.length ? 'border border-theme-border bg-white text-theme-text hover:bg-theme-hover' : 'border-2 border-dashed border-primary-400 bg-primary-50/50 text-primary-800 hover:bg-primary-100'
          ]"
          @click="openCreateModal"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ workspaces.length ? $t('select_workspace.create_new') : $t('select_workspace.create_first') }}
        </button>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600 text-center">{{ error }}</p>
    </div>

    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="showCreateModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="create-workspace-title"
          @click.self="!createLoading && closeCreateModal()"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="create-workspace-title" class="text-lg font-semibold text-theme-text">{{ $t('select_workspace.create_modal_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('select_workspace.create_modal_description') }}</p>
            <div v-if="createModalError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3" role="alert">
              <p class="text-sm text-red-800">{{ createModalError }}</p>
            </div>
            <form class="mt-4 space-y-4" @submit.prevent="handleCreate">
              <div>
                <label for="workspace-name" class="form-label">{{ $t('select_workspace.workspace_name') }}</label>
                <input
                  id="workspace-name"
                  v-model="createForm.name"
                  type="text"
                  required
                  class="input-field"
                  :placeholder="$t('select_workspace.workspace_name_placeholder')"
                  :disabled="createLoading"
                  autofocus
                >
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="createLoading" @click="closeCreateModal">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="createLoading || !createForm.name.trim()">
                  {{ createLoading ? $t('common.loading') : $t('common.create') }}
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
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { workspace_api } from '@/api'
import { loadWorkspaceFromStorage, saveWorkspaceToStorage } from '@/utils/storage'
import { formatDate } from '@/utils/dt'

const { t } = useI18n()
const router = useRouter()

const storedWorkspace = loadWorkspaceFromStorage()
const workspaces = ref([])
const error = ref(null)
const showCreateModal = ref(false)
const createModalError = ref(null)
const createLoading = ref(false)
const createForm = reactive({ name: '' })

function openCreateModal() {
  createModalError.value = null
  createForm.name = ''
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  createModalError.value = null
  createForm.name = ''
}

function selectWorkspace(workspace) {
  saveWorkspaceToStorage(workspace)
  router.push('/dashboard')
}

async function loadWorkspaces() {
  error.value = null
  try {
    const res = await workspace_api.list({})
    workspaces.value = res.data.workspaces
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('select_workspace.failed_load')
  }
}

async function handleCreate() {
  createLoading.value = true
  createModalError.value = null

  try {
    const res = await workspace_api.create({ name: createForm.name })
    saveWorkspaceToStorage({ ...res.data, role: 'owner', status: 'active' })
    showCreateModal.value = false
    createForm.name = ''
    router.push('/dashboard')
  } catch (err) {
    createModalError.value = err.response?.data?.errorMessage || t('select_workspace.failed_create')
  } finally {
    createLoading.value = false
  }
}

onMounted(async () => {
  await loadWorkspaces()
  if (workspaces.value.length === 0 && !error.value) {
    openCreateModal()
  }
})
</script>

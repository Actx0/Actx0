<template>
  <div class="min-h-screen flex items-center justify-center bg-theme-bg px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="flex justify-center mb-8">
          <router-link to="/">
            <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-24 w-auto">
          </router-link>
        </div>
      </div>

      <div class="bg-white rounded-lg border border-theme-border p-8 shadow-sm text-center">
        <h2 class="text-2xl font-semibold text-theme-text mb-6">{{ $t('verify_email.title') }}</h2>

        <div v-if="loading" class="flex items-center justify-center py-4">
          <svg class="animate-spin mr-3 h-5 w-5 text-theme-text" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <span class="text-sm text-theme-textLight">{{ $t('verify_email.verifying') }}</span>
        </div>

        <div v-else-if="error" class="rounded-md border border-red-200 bg-red-50 p-3">
          <p class="text-sm text-red-800">{{ error }}</p>
        </div>

        <div v-else class="rounded-md border border-green-200 bg-green-50 p-3">
          <p class="text-sm text-green-800">{{ $t('verify_email.success') }}</p>
        </div>

        <p class="text-center mt-6">
          <router-link to="/login" class="text-sm text-primary-800 hover:text-primary-900 font-medium">
            {{ $t('verify_email.back_to_login') }}
          </router-link>
        </p>
      </div>

      <p class="text-center text-xs text-theme-textLight mt-8">
        {{ $t('common.copyright') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { auth_api } from '@/api'
import { showFlash } from '@/lib/flash'
import { user, saveUser } from '@/lib/auth'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  try {
    await auth_api.verifyEmail(route.params.token)
    if (user.value) {
      saveUser({ ...user.value, isEmailVerified: true })
    }
    showFlash(t('common.saved'))
    setTimeout(() => router.push('/login'), 2000)
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('verify_email.error_default')
  } finally {
    loading.value = false
  }
})
</script>

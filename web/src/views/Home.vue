<template>
  <div class="min-h-screen bg-theme-bg" />
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { setup_api } from '@/api'

const router = useRouter()

onMounted(async () => {
  const response = await setup_api.checkInstalled()
  router.replace(response.data.installed ? '/login' : '/setup')
})
</script>

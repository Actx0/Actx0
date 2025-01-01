<template>
  <div class="landing-page min-h-screen bg-theme-bg">
    <header class="landing-header sticky top-0 z-50">
      <nav class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <router-link to="/" class="flex items-center gap-2.5">
          <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-8 w-auto" />
          <span class="text-lg font-semibold tracking-tight">{{ $t('common.brand') }}</span>
        </router-link>
        <router-link to="/" class="landing-link">{{ $t('docs_page.back_to_home') }}</router-link>
      </nav>
    </header>

    <main class="mx-auto flex max-w-2xl flex-col items-center px-6 py-24 text-center">
      <p class="section-label">{{ $t('docs_page.label') }}</p>
      <h1 class="mt-3 text-3xl font-semibold tracking-tight text-theme-text">{{ $t('docs_page.title') }}</h1>
      <p class="mt-4 text-muted">{{ $t('docs_page.description') }}</p>
    </main>
  </div>
</template>

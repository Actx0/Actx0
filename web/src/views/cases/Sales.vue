<template>
  <CaseLayout
    :title="$t('cases.sales.meta_title')"
    :description="$t('cases.sales.meta_description')"
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">{{ $t('cases.shared.breadcrumb') }}</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">{{ $t('cases.sales.label') }}</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">{{ $t('cases.sales.label') }}</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              {{ $t('cases.sales.hero_title') }}
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              {{ $t('cases.sales.hero_subtitle') }}
            </p>
            <div class="mt-8 opacity-0 animate-fade-up animate-delay-300">
              <router-link :to="startPath" class="landing-btn">{{ $t('cases.shared.get_started') }}</router-link>
            </div>
          </div>

          <div class="reveal">
            <div class="glass-card p-6">
              <div class="flex items-start gap-3">
                <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-primary-100 text-sm font-medium text-primary-800">U</div>
                <div class="rounded-2xl rounded-tl-sm bg-theme-hover px-4 py-3 text-sm text-theme-text">
                  {{ $t('cases.sales.demo_user') }}
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  {{ $t('cases.sales.demo_agent') }}
                </div>
                <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-primary-800 text-sm font-medium text-white">A</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">{{ $t('cases.shared.problems_label') }}</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('cases.sales.problems_title') }}</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div v-for="i in 3" :key="i" class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">{{ $t(`cases.sales.problem_${i}_title`) }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t(`cases.sales.problem_${i}_text`) }}</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">{{ $t('cases.shared.features_label') }}</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('cases.sales.features_title') }}</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div v-for="i in 4" :key="i" class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">{{ $t(`cases.sales.feature_${i}_title`) }}</h3>
            <p class="mt-2 text-sm text-muted">{{ $t(`cases.sales.feature_${i}_text`) }}</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div v-for="i in 2" :key="i" class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">{{ $t(`cases.sales.benefit_${i}_title`) }}</h3>
              <p class="mt-3 text-muted">{{ $t(`cases.sales.benefit_${i}_text`) }}</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li v-for="j in 3" :key="j" class="flex items-start gap-2">
                  <span class="mt-1 text-primary-700">→</span>
                  {{ $t(`cases.sales.benefit_${i}_item_${j}`) }}
                </li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">{{ $t('cases.shared.get_started') }}</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">{{ $t('cases.shared.why_label') }}</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('cases.sales.why_title') }}</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="i in 6" :key="i" class="reveal step-card">
            <h3 class="font-semibold">{{ $t(`cases.sales.why_${i}_title`) }}</h3>
            <p class="mt-2 text-sm text-muted">{{ $t(`cases.sales.why_${i}_text`) }}</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">{{ $t('cases.shared.more_label') }}</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">{{ $t('cases.shared.more_title') }}</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <router-link to="/cases/customer-support" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">{{ $t('landing.use_cases.customer_support.label') }}</p>
              <h3 class="mt-3 text-lg font-semibold">{{ $t('landing.use_cases.customer_support.title') }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t('landing.use_cases.customer_support.text') }}</p>
            </router-link>
            <router-link to="/cases/healthcare" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">{{ $t('landing.use_cases.healthcare.label') }}</p>
              <h3 class="mt-3 text-lg font-semibold">{{ $t('landing.use_cases.healthcare.title') }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t('landing.use_cases.healthcare.text') }}</p>
            </router-link>
            <router-link to="/cases/devtools" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">{{ $t('landing.use_cases.devtools.label') }}</p>
              <h3 class="mt-3 text-lg font-semibold">{{ $t('landing.use_cases.devtools.title') }}</h3>
              <p class="mt-2 text-sm text-muted">{{ $t('landing.use_cases.devtools.text') }}</p>
            </router-link>
          </div>
        </div>
      </section>
    </template>
  </CaseLayout>
</template>

<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

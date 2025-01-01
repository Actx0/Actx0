<template>
  <Teleport to="body">
    <Transition name="global-loading">
      <div
        v-if="globalLoading"
        class="global-loading-overlay"
        role="status"
        aria-live="polite"
        aria-busy="true"
        aria-label="Loading"
      >
        <svg class="global-loading-spinner" fill="none" viewBox="0 0 24 24" aria-hidden="true">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { globalLoading } from '@/lib/loading'
</script>

<style scoped>
.global-loading-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--theme-bg);
  opacity: 0.96;
  pointer-events: all;
}

.global-loading-spinner {
  width: 2rem;
  height: 2rem;
  color: var(--theme-text);
  animation: spin 1s linear infinite;
}

.global-loading-enter-active,
.global-loading-leave-active {
  transition: opacity 0.15s ease;
}

.global-loading-enter-from,
.global-loading-leave-to {
  opacity: 0;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>

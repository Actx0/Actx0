import { ref } from 'vue'

const pending = ref(0)
export const globalLoading = ref(false)

let hideTimer = null

export function startLoading() {
  pending.value++
  clearTimeout(hideTimer)
  globalLoading.value = true
}

export function stopLoading() {
  pending.value = Math.max(0, pending.value - 1)
  if (pending.value > 0) return

  clearTimeout(hideTimer)
  hideTimer = setTimeout(() => {
    if (pending.value === 0) globalLoading.value = false
  }, 150)
}

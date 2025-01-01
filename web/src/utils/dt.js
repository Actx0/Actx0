/** Short date (e.g. "Mar 27, 2026"). Pass `locale` from `useI18n()` when needed. */
export function formatDate(iso, locale) {
  return new Date(iso).toLocaleDateString(locale, { month: 'short', day: 'numeric', year: 'numeric' })
}

/** Date + time (e.g. "Mar 27, 2026, 3:45 PM"). */
export function formatDateTime(iso, locale) {
  return new Date(iso).toLocaleString(locale, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** Numeric 24h datetime with seconds (e.g. "03/27/2026, 15:45:01"). */
export function formatPromptDateTime(iso, locale) {
  return new Date(iso).toLocaleString(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

export function formatRelative(dateString) {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins} ${diffMins === 1 ? 'minute' : 'minutes'} ago`
  if (diffHours < 24) return `${diffHours} ${diffHours === 1 ? 'hour' : 'hours'} ago`
  if (diffDays < 7) return `${diffDays} ${diffDays === 1 ? 'day' : 'days'} ago`

  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const month = months[date.getMonth()]
  const day = date.getDate()
  const year = date.getFullYear()
  return year === now.getFullYear() ? `${month} ${day}` : `${month} ${day}, ${year}`
}

export function formatExpires(dateString) {
  const date = new Date(dateString)
  const now = new Date()
  if (date < now) return formatRelative(dateString)
  const diffMs = date - now
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  if (diffMins < 60) return `in ${diffMins} ${diffMins === 1 ? 'minute' : 'minutes'}`
  if (diffHours < 24) return `in ${diffHours} ${diffHours === 1 ? 'hour' : 'hours'}`
  if (diffDays < 7) return `in ${diffDays} ${diffDays === 1 ? 'day' : 'days'}`
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  return `${months[date.getMonth()]} ${date.getDate()}`
}

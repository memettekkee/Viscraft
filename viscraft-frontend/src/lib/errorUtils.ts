import { ERROR_MESSAGES } from '../constants'

/**
 * Describes the action the UI should take in response to an error code.
 * - toast: show a toast notification
 * - inline: display inline in the current component (e.g., form banner)
 * - redirect: navigate to another route (with a toast)
 */
export interface ErrorAction {
  action: 'toast' | 'inline' | 'redirect'
  message: string
  redirectTo?: string
}

/**
 * Resolves an error code to a UI action, message, and optional redirect target.
 *
 * Centralizes the error routing logic so property tests can verify
 * that every known error code maps to the correct user-facing behavior.
 *
 * Rules:
 * - ERR_01 and ERR_08 → redirect to /workspace with toast
 * - ERR_09 → redirect to / (login) with toast (handled by axios interceptor)
 * - ERR_02 → inline (rate limit shown in modal)
 * - Network errors → toast with "Unable to connect to server"
 * - All others → toast with mapped message or generic fallback
 *
 * Validates: Requirements 12.3, 12.4, 12.5, 12.6
 */
export function resolveErrorAction(errorCode: string): ErrorAction {
  const message = ERROR_MESSAGES[errorCode] ?? 'An unexpected error occurred'

  switch (errorCode) {
    case 'ERR_01':
      return { action: 'redirect', message, redirectTo: '/workspace' }

    case 'ERR_08':
      return { action: 'redirect', message, redirectTo: '/workspace' }

    case 'ERR_09':
      return { action: 'redirect', message, redirectTo: '/' }

    case 'ERR_02':
      return { action: 'inline', message }

    case 'NETWORK_ERROR':
      return { action: 'toast', message }

    default:
      return { action: 'toast', message }
  }
}

import { useCallback } from 'react'
import { showToast, type ToastType } from '../components/CustomToast'
import { ERROR_MESSAGES } from '../constants'

/**
 * Options for showing a toast notification.
 */
export interface ShowToastOptions {
  /** The message to display as the toast title */
  title: string
  /** Optional description below the title */
  description?: string
  /** Toast variant: error (oxblood), success (moss), info (amber) */
  type?: ToastType
  /** Override auto-dismiss duration in ms (default: 5000 error, 3000 others) */
  duration?: number
}

/**
 * Return type for the useToast hook.
 */
export interface UseToastReturn {
  /** Show a toast with custom options */
  toast: (options: ShowToastOptions) => void
  /** Show an error toast by error code (looks up ERROR_MESSAGES) */
  toastError: (errorCode: string) => void
  /** Show a custom error toast with a direct message */
  toastErrorMessage: (message: string, description?: string) => void
  /** Show a success toast */
  toastSuccess: (title: string, description?: string) => void
  /** Show an info toast */
  toastInfo: (title: string, description?: string) => void
}

/**
 * Custom hook wrapping the Viscraft toast system.
 *
 * Provides a simple API for showing toast notifications styled with the
 * Cartographer's Atlas design system tokens:
 * - Error toasts: oxblood (#8B2E2E) — 5s auto-dismiss
 * - Success toasts: moss (#3E5C4E) — 3s auto-dismiss
 * - Info toasts: amber (#C9762C) — 3s auto-dismiss
 *
 * Supports showing error messages by error code (from ERROR_MESSAGES constant)
 * or with custom messages directly.
 *
 * Validates: Requirements 12.3, 12.4, 12.5
 */
export function useToast(): UseToastReturn {
  const toast = useCallback((options: ShowToastOptions) => {
    showToast({
      type: options.type ?? 'info',
      title: options.title,
      description: options.description,
      duration: options.duration,
    })
  }, [])

  const toastError = useCallback((errorCode: string) => {
    const message = ERROR_MESSAGES[errorCode] ?? 'An unexpected error occurred'
    showToast({
      type: 'error',
      title: message,
    })
  }, [])

  const toastErrorMessage = useCallback((message: string, description?: string) => {
    showToast({
      type: 'error',
      title: message,
      description,
    })
  }, [])

  const toastSuccess = useCallback((title: string, description?: string) => {
    showToast({
      type: 'success',
      title,
      description,
    })
  }, [])

  const toastInfo = useCallback((title: string, description?: string) => {
    showToast({
      type: 'info',
      title,
      description,
    })
  }, [])

  return {
    toast,
    toastError,
    toastErrorMessage,
    toastSuccess,
    toastInfo,
  }
}

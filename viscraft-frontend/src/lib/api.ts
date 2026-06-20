import axios from 'axios'
import type { InternalAxiosRequestConfig, AxiosError } from 'axios'
import { useAuthStore } from '../store/authStore'
import { showToast } from '../components/CustomToast'
import { ERROR_MESSAGES } from '../constants'
import type { ApiResponse } from '../types'

// Augment the global Window interface for runtime config
declare global {
  interface Window {
    __VISCRAFT_CONFIG__?: {
      API_BASE_URL?: string
    }
  }
}

const api = axios.create({
  baseURL: window.__VISCRAFT_CONFIG__?.API_BASE_URL || 'http://localhost:8089',
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor — inject Bearer token from authStore
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor — handle global error codes and network failures
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse>) => {
    const errorCode = error.response?.data?.errorCode

    if (errorCode === 'ERR_09') {
      // Session expired — only redirect if user was authenticated (has token)
      const token = useAuthStore.getState().token
      if (token) {
        useAuthStore.getState().clearAuth()
        console.warn('[Viscraft] Session expired. Redirecting to login.')
        window.location.href = '/'
      }
      // kalau tidak ada token = login attempt gagal, biarkan error bubble ke caller
    } else if (errorCode === 'ERR_01') {
      // Resource not found — toast and redirect to workspace
      showToast({ type: 'error', title: ERROR_MESSAGES.ERR_01 })
      window.location.href = '/workspace'
    } else if (errorCode === 'ERR_08') {
      // Project not found — toast and redirect to workspace
      showToast({ type: 'error', title: ERROR_MESSAGES.ERR_08 })
      window.location.href = '/workspace'
    } else if (!error.response) {
      // Network error (no response from server)
      showToast({ type: 'error', title: ERROR_MESSAGES.NETWORK_ERROR })
    }

    return Promise.reject(error)
  }
)

export { api }
export default api

import axios from 'axios'
import type { InternalAxiosRequestConfig, AxiosError } from 'axios'
import { useAuthStore } from '../store/authStore'
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
  baseURL: window.__VISCRAFT_CONFIG__?.API_BASE_URL || 'http://localhost:8080',
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

// Response interceptor — handle ERR_09 (session expired / 401)
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.data?.errorCode === 'ERR_09') {
      useAuthStore.getState().clearAuth()
      console.warn('[Viscraft] Session expired. Redirecting to login.')
      window.location.href = '/'
    }
    return Promise.reject(error)
  }
)

export { api }
export default api

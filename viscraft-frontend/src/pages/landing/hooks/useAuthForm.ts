import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login, createUser } from '../../../service/auth'
import { useAuthStore } from '../../../store/authStore'
import { ERROR_MESSAGES } from '../../../constants'
import type { AxiosError } from 'axios'
import type { ApiResponse } from '../../../types'

type TabType = 'login' | 'register'

function validateEmail(email: string): string | null {
  if (!email) return 'Email is required'
  if (!email.includes('@') || !email.includes('.')) return 'Invalid email format'
  return null
}

function validatePassword(password: string): string | null {
  if (!password) return 'Password is required'
  if (password.length < 8) return 'Password must be at least 8 characters'
  if (password.length > 72) return 'Password must be at most 72 characters'
  return null
}

export function useAuthForm(onSuccess: () => void) {
  const navigate = useNavigate()
  const setAuth = useAuthStore((s) => s.setAuth)

  const [activeTab, setActiveTab] = useState<TabType>('login')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [emailError, setEmailError] = useState<string | null>(null)
  const [passwordError, setPasswordError] = useState<string | null>(null)
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  function resetForm() {
    setEmail('')
    setPassword('')
    setName('')
    setEmailError(null)
    setPasswordError(null)
    setApiError(null)
  }

  function handleTabSwitch(tab: TabType) {
    setActiveTab(tab)
    resetForm()
  }

  function validate(): boolean {
    const eErr = validateEmail(email)
    const pErr = validatePassword(password)
    setEmailError(eErr)
    setPasswordError(pErr)
    return !eErr && !pErr
  }

  const hasValidationErrors = !!validateEmail(email) || !!validatePassword(password)

  async function handleSubmit() {
    if (!validate()) return

    setIsSubmitting(true)
    setApiError(null)

    try {
      let response: ApiResponse<unknown>

      if (activeTab === 'login') {
        response = await login({ email, password })
      } else {
        response = await createUser({ email, password, name: name || undefined })
      }

      const token = response.token
      const user = response.data

      if (token && user) {
        setAuth(token, user as Parameters<typeof setAuth>[1])
        onSuccess()
        navigate('/workspace')
      }
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>

      if (axiosError.response?.data?.errorCode) {
        const code = axiosError.response.data.errorCode
        const backendMessage = axiosError.response.data.message
        setApiError(backendMessage ?? (ERROR_MESSAGES[code] ?? 'An error occurred'))
      } else {
        setApiError(ERROR_MESSAGES.NETWORK_ERROR)
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    activeTab,
    email,
    setEmail,
    password,
    setPassword,
    name,
    setName,
    emailError,
    setEmailError,
    passwordError,
    setPasswordError,
    apiError,
    isSubmitting,
    hasValidationErrors,
    handleTabSwitch,
    handleSubmit,
  }
}

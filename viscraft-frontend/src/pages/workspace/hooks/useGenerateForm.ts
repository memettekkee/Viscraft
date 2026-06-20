import { useState } from 'react'
import { useSWRConfig } from 'swr'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { generateScene } from '../../../service/scene'
import { ERROR_MESSAGES } from '../../../constants'
import type { AxiosError } from 'axios'
import type { ApiResponse, Scene } from '../../../types'

/**
 * Custom hook for Generate Modal form state, validation, and submission logic.
 * Simplified — prompt-only, no reference images.
 */
export function useGenerateForm() {
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const closeModal = useWorkspaceStore((s) => s.closeModal)
  const { mutate } = useSWRConfig()

  const [prompt, setPrompt] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  function resetForm() {
    setPrompt('')
    setErrors({})
    setApiError(null)
  }

  function validate(): boolean {
    const trimmed = prompt.trim()
    const newErrors: Record<string, string> = {}

    if (trimmed.length < 3) {
      newErrors.prompt = 'Prompt must be at least 3 characters'
    } else if (trimmed.length > 300) {
      newErrors.prompt = 'Prompt must be 300 characters or less'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  async function handleSubmit() {
    if (!validate()) return
    if (!activeProjectId) return

    setIsSubmitting(true)
    setApiError(null)

    try {
      await generateScene({
        projectId: activeProjectId,
        prompt: prompt.trim(),
      })

      // Mutate the SWR cache to revalidate scene list
      await mutate(['/scenes/list', { projectId: activeProjectId }])

      // Close modal and reset form
      closeModal()
      resetForm()
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse<Scene>>

      if (axiosError.response?.data?.errorCode === 'ERR_02') {
        setApiError(ERROR_MESSAGES.ERR_02)
      } else if (axiosError.response?.data?.errorCode) {
        const code = axiosError.response.data.errorCode
        setApiError(ERROR_MESSAGES[code] ?? 'An error occurred')
      } else {
        setApiError(ERROR_MESSAGES.NETWORK_ERROR)
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    prompt,
    setPrompt,
    errors,
    apiError,
    isSubmitting,
    handleSubmit,
    resetForm,
  }
}

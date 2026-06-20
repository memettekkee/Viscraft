import { useEffect, useState, useCallback, useMemo } from 'react'
import useSWR from 'swr'
import { useSWRConfig } from 'swr'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { generateScene } from '../../../service/scene'
import { fetchAllPromptOptions, groupByCategory } from '../../../service/promptOptions'
import { buildPrompt } from '../utils/buildPrompt'
import { ERROR_MESSAGES } from '../../../constants'
import type { PromptOption } from '../../../service/promptOptions'
import type { AxiosError } from 'axios'
import type { ApiResponse } from '../../../types'

/**
 * Hook for Generate Ad Shot modal — manages form state, prompt options,
 * generated prompt computation, and submission.
 */
export function useGenerateAdShot(isOpen: boolean) {
  const prefillPrompt = useWorkspaceStore((s) => s.prefillPrompt)
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const regenerateSceneId = useWorkspaceStore((s) => s.regenerateSceneId)
  const regenerateFileUrl = useWorkspaceStore((s) => s.regenerateFileUrl)
  const closeModal = useWorkspaceStore((s) => s.closeModal)
  const { mutate } = useSWRConfig()

  const [userPrompt, setUserPrompt] = useState<string>('')
  const [selectedOptions, setSelectedOptions] = useState<PromptOption[]>([])
  const [uploadedReferenceImage, setUploadedReferenceImage] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)

  // Fetch ALL prompt options in a single request
  const { data: allOptions } = useSWR(isOpen ? 'prompt-options-all' : null, fetchAllPromptOptions)

  const grouped = useMemo(() => groupByCategory(allOptions ?? []), [allOptions])

  const categories = useMemo(() => [
    { label: 'Background', category: 'background', options: grouped['background'] ?? [], multi: false },
    { label: 'Lighting', category: 'lighting', options: grouped['lighting'] ?? [], multi: false },
    { label: 'Mood', category: 'mood', options: grouped['mood'] ?? [], multi: true },
    { label: 'Angle', category: 'angle', options: grouped['angle'] ?? [], multi: false },
    { label: 'Props', category: 'props', options: grouped['props'] ?? [], multi: true },
  ], [grouped])

  // Pre-fill prompt when modal opens
  useEffect(() => {
    if (isOpen && prefillPrompt && typeof prefillPrompt === 'string') {
      setUserPrompt(prefillPrompt)
    }
  }, [isOpen, prefillPrompt])

  const resetForm = useCallback(() => {
    setUserPrompt('')
    setSelectedOptions([])
    setUploadedReferenceImage(null)
    setApiError(null)
  }, [])

  const handleUserPromptChange = useCallback((value: string) => {
    setUserPrompt(String(value ?? ''))
  }, [])

  const toggleOption = useCallback((option: PromptOption, multi: boolean) => {
    setSelectedOptions((prev) => {
      const exists = prev.find((o) => o.id === option.id)
      if (exists) {
        return prev.filter((o) => o.id !== option.id)
      }
      if (multi) {
        return [...prev, option]
      }
      return [...prev.filter((o) => o.category !== option.category), option]
    })
  }, [])

  const isSelected = useCallback((optionId: string): boolean => {
    return selectedOptions.some((o) => o.id === optionId)
  }, [selectedOptions])

  const handleUploadReference = useCallback((base64: string | null) => {
    setUploadedReferenceImage(base64)
  }, [])

  const handleFileSelect = useCallback(async (file: File | null) => {
    if (!file) return
    const { validateImageFile, fileToBase64 } = await import('../utils/referenceImage')
    const validationError = validateImageFile(file)
    if (validationError) return
    try {
      const base64 = await fileToBase64(file)
      setUploadedReferenceImage(base64)
    } catch { /* ignore */ }
  }, [])

  // Compute generated prompt
  const promptText = String(userPrompt ?? '')
  const isPromptValid = promptText.trim().length >= 3
  const generatedPrompt = isPromptValid ? buildPrompt(promptText, selectedOptions) : ''

  const handleSubmit = useCallback(async () => {
    if (!activeProjectId || !isPromptValid) return

    setIsSubmitting(true)
    setApiError(null)

    try {
      await generateScene({
        projectId: activeProjectId,
        prompt: promptText.trim(),
        generatedPrompt,
        ...(uploadedReferenceImage ? { uploadedReferenceImage } : {}),
        ...(regenerateSceneId && !uploadedReferenceImage ? { referenceSceneId: regenerateSceneId } : {}),
      })

      await mutate(['/scenes/list', { projectId: activeProjectId }])
      closeModal()
      resetForm()

      // Show success toast
      const { showToast } = await import('../../../components/CustomToast')
      showToast({ type: 'success', title: 'Ad shot generation started!' })
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>
      const code = axiosError.response?.data?.errorCode
      const backendMessage = axiosError.response?.data?.message
      // Prefer backend message, fall back to local error map, then generic
      const msg = backendMessage ?? (code ? (ERROR_MESSAGES[code] ?? 'An error occurred') : 'Network error')
      setApiError(msg)

      const { showToast } = await import('../../../components/CustomToast')
      showToast({ type: 'error', title: msg })
    } finally {
      setIsSubmitting(false)
    }
  }, [activeProjectId, isPromptValid, promptText, generatedPrompt, mutate, closeModal, resetForm])

  return {
    userPrompt: promptText,
    handleUserPromptChange,
    selectedOptions,
    toggleOption,
    isSelected,
    categories,
    generatedPrompt,
    isPromptValid,
    isSubmitting,
    apiError,
    handleSubmit,
    resetForm,
    isRegenerate: !!prefillPrompt,
    uploadedReferenceImage,
    handleUploadReference,
    handleFileSelect,
    regenerateFileUrl,
  }
}

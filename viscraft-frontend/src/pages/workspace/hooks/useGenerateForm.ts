import { useState } from 'react'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { generateImage } from '../../../service/image'
import { validateGenerateForm } from '../../../lib/inputValidation'
import { ERROR_MESSAGES } from '../../../constants'
import { useGallery } from './useGallery'
import type { Genre, AssetType, Mood } from '../../../types'
import type { AxiosError } from 'axios'
import type { ApiResponse } from '../../../types'

export type GenerateMode = 'create' | 'from-reference'

/**
 * Custom hook for Generate Modal form state, validation, and submission logic.
 * Extracts all non-UI concerns from the GenerateModal component.
 *
 * Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 12.1
 */
export function useGenerateForm() {
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const closeModal = useWorkspaceStore((s) => s.closeModal)
  const { mutate } = useGallery(activeProjectId)

  const [mode, setMode] = useState<GenerateMode>('create')
  const [prompt, setPrompt] = useState('')
  const [genre, setGenre] = useState<Genre | ''>('')
  const [assetType, setAssetType] = useState<AssetType | ''>('')
  const [mood, setMood] = useState<Mood | ''>('')
  const [referenceImage, setReferenceImage] = useState<string | undefined>(undefined)

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [fromCache, setFromCache] = useState(false)

  function resetForm() {
    setPrompt('')
    setGenre('')
    setAssetType('')
    setMood('')
    setReferenceImage(undefined)
    setErrors({})
    setApiError(null)
    setFromCache(false)
  }

  function handleModeSwitch(newMode: GenerateMode) {
    setMode(newMode)
    if (newMode === 'create') {
      setReferenceImage(undefined)
    }
  }

  function validate(): boolean {
    const result = validateGenerateForm({
      prompt,
      genre,
      assetType,
      mood,
      referenceImage,
    })
    setErrors(result.errors)
    return result.valid
  }

  async function handleSubmit() {
    if (!validate()) return
    if (!activeProjectId) return

    setIsSubmitting(true)
    setApiError(null)
    setFromCache(false)

    try {
      const response = await generateImage({
        projectId: activeProjectId,
        prompt: prompt.trim(),
        genre: genre as Genre,
        assetType: assetType as AssetType,
        mood: mood as Mood,
        referenceImage,
      })

      // Both 202 (processing) and 200 (cache hit) return success
      // We detect cache hit by checking if the image already has a completed status
      if (response.data?.status === 'completed') {
        setFromCache(true)
      }

      // Mutate the SWR cache to revalidate gallery
      await mutate()

      // Close modal and reset
      closeModal()
      resetForm()
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>

      if (axiosError.response?.data?.errorCode === 'ERR_02') {
        // Rate limit — keep modal open, show inline banner
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
    mode,
    prompt,
    setPrompt,
    genre,
    setGenre,
    assetType,
    setAssetType,
    mood,
    setMood,
    referenceImage,
    setReferenceImage,
    errors,
    setErrors,
    apiError,
    isSubmitting,
    fromCache,
    handleModeSwitch,
    handleSubmit,
    resetForm,
  }
}

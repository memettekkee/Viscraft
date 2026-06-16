import { useState, useEffect, useCallback } from 'react'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { generateImage } from '../../../service/image'
import { validateGenerateForm } from '../../../lib/inputValidation'
import { imageUrlToBase64 } from '../utils/referenceImage'
import { ERROR_MESSAGES } from '../../../constants'
import { useGallery } from './useGallery'
import type { Genre, AssetType, Mood } from '../../../types'
import type { AxiosError } from 'axios'
import type { ApiResponse } from '../../../types'

export type GenerateMode = 'create' | 'from-reference'

/**
 * Custom hook for Regenerate Modal form state, validation, and submission logic.
 * Pre-fills all fields from the source image and auto-selects "From Reference" mode.
 *
 * Validates: Requirements 10.1, 10.2, 10.3
 */
export function useRegenerateForm() {
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const regenerateSource = useWorkspaceStore((s) => s.regenerateSource)
  const closeModal = useWorkspaceStore((s) => s.closeModal)
  const { mutate } = useGallery(activeProjectId)

  const [mode, setMode] = useState<GenerateMode>('from-reference')
  const [prompt, setPrompt] = useState('')
  const [genre, setGenre] = useState<Genre | ''>('')
  const [assetType, setAssetType] = useState<AssetType | ''>('')
  const [mood, setMood] = useState<Mood | ''>('')
  const [referenceImage, setReferenceImage] = useState<string | undefined>(undefined)

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isLoadingSourceReference, setIsLoadingSourceReference] = useState(false)

  /**
   * Pre-fill form fields from source image when regenerateSource changes.
   * Auto-selects "From Reference" mode and converts source image URL to base64.
   */
  const prefillFromSource = useCallback(async () => {
    if (!regenerateSource) return

    setPrompt(regenerateSource.prompt)
    setGenre(regenerateSource.genre)
    setAssetType(regenerateSource.assetType)
    setMood(regenerateSource.mood)
    setMode('from-reference')
    setErrors({})
    setApiError(null)

    // Convert source image fileUrl to base64 for reference
    if (regenerateSource.fileUrl) {
      try {
        setIsLoadingSourceReference(true)
        const base64 = await imageUrlToBase64(regenerateSource.fileUrl)
        setReferenceImage(base64)
      } catch {
        setErrors((prev) => ({ ...prev, referenceImage: 'Failed to load source image as reference' }))
      } finally {
        setIsLoadingSourceReference(false)
      }
    }
  }, [regenerateSource])

  useEffect(() => {
    prefillFromSource()
  }, [prefillFromSource])

  function resetForm() {
    setPrompt('')
    setGenre('')
    setAssetType('')
    setMood('')
    setReferenceImage(undefined)
    setMode('from-reference')
    setErrors({})
    setApiError(null)
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

    try {
      await generateImage({
        projectId: activeProjectId,
        prompt: prompt.trim(),
        genre: genre as Genre,
        assetType: assetType as AssetType,
        mood: mood as Mood,
        referenceImage,
      })

      // Mutate the SWR cache to revalidate gallery — new image appears
      await mutate()

      // Close modal and reset
      closeModal()
      resetForm()
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>

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
    isLoadingSourceReference,
    handleModeSwitch,
    handleSubmit,
    resetForm,
  }
}

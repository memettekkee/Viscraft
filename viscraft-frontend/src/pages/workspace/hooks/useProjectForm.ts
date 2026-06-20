import { useState, useCallback } from 'react'
import { useSWRConfig } from 'swr'
import { createProject } from '../../../service/project'
import { useWorkspaceStore } from '../../../store/workspaceStore'

/**
 * Hook for project (campaign) creation form state and submission.
 * Extracted from ProjectModal component.
 */
export function useProjectForm(onSuccess: () => void) {
  const { mutate } = useSWRConfig()
  const setActiveProject = useWorkspaceStore((s) => s.setActiveProject)

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [productCategory, setProductCategory] = useState('general')
  const [visualStyle, setVisualStyle] = useState('')
  const [nameError, setNameError] = useState<string | null>(null)
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const resetForm = useCallback(() => {
    setName('')
    setDescription('')
    setProductCategory('general')
    setVisualStyle('')
    setNameError(null)
    setApiError(null)
    setIsSubmitting(false)
  }, [])

  const handleNameChange = useCallback((value: string) => {
    setName(value)
    setNameError(null)
  }, [])

  const handleSubmit = useCallback(async () => {
    const trimmed = name.trim()
    if (trimmed.length === 0 || trimmed.length > 255) {
      setNameError('Campaign name is required (max 255 characters)')
      return
    }

    setNameError(null)
    setApiError(null)
    setIsSubmitting(true)

    try {
      const response = await createProject({
        name: trimmed,
        description: description.trim() || undefined,
        productCategory,
        visualStyle: visualStyle.trim() || undefined,
      })

      if (response.success && response.data) {
        await mutate(['/projects/list'])
        setActiveProject(response.data.id)
        onSuccess()
        resetForm()
      } else {
        setApiError(response.message || 'Failed to create campaign')
      }
    } catch {
      setApiError('Unable to connect to server')
    } finally {
      setIsSubmitting(false)
    }
  }, [name, description, productCategory, visualStyle, mutate, setActiveProject, onSuccess, resetForm])

  const isValid = name.trim().length > 0

  return {
    name,
    setName: handleNameChange,
    description,
    setDescription,
    productCategory,
    setProductCategory,
    visualStyle,
    setVisualStyle,
    nameError,
    apiError,
    isSubmitting,
    isValid,
    handleSubmit,
    resetForm,
  }
}

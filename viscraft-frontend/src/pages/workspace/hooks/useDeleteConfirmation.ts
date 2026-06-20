import { useState, useCallback } from 'react'
import { useSceneActions } from './useSceneActions'

/**
 * Hook for delete confirmation modal state and actions.
 * Extracted from StoryboardGrid component.
 */
export function useDeleteConfirmation() {
  const { handleDelete } = useSceneActions()
  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const onDeleteRequest = useCallback((sceneId: string) => {
    setDeleteTargetId(sceneId)
  }, [])

  const onDeleteCancel = useCallback(() => {
    setDeleteTargetId(null)
  }, [])

  const onDeleteConfirm = useCallback(async () => {
    if (!deleteTargetId) return
    setIsDeleting(true)
    await handleDelete(deleteTargetId)
    setIsDeleting(false)
    setDeleteTargetId(null)
  }, [deleteTargetId, handleDelete])

  return {
    deleteTargetId,
    isDeleting,
    isDeleteModalOpen: deleteTargetId !== null,
    onDeleteRequest,
    onDeleteCancel,
    onDeleteConfirm,
  }
}

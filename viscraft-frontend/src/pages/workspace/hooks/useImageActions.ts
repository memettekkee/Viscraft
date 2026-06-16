import { useSWRConfig } from 'swr'
import { generateImage, deleteImage } from '../../../service/image'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import type { Image } from '../../../types'

/**
 * Extracts image card action logic: retry (re-generate same payload),
 * delete (remove from gallery), and regenerate (open modal with prefill).
 *
 * Validates: Requirements 8.5, 8.6, 11.2
 */
export function useImageActions() {
  const { mutate } = useSWRConfig()
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const openRegenerateModal = useWorkspaceStore((s) => s.openRegenerateModal)

  /**
   * Retry: resubmits the same generation payload for a failed image,
   * transitions the card back to processing state in the SWR cache.
   */
  async function retry(image: Image) {
    if (!activeProjectId) return

    try {
      const response = await generateImage({
        projectId: activeProjectId,
        prompt: image.prompt,
        genre: image.genre,
        assetType: image.assetType,
        mood: image.mood,
      })

      // Mutate gallery cache to reflect new processing card
      mutate(['/images/list', { projectId: activeProjectId }], undefined, {
        revalidate: true,
      })

      return response
    } catch {
      // Error handling is left to the caller or global interceptor
    }
  }

  /**
   * Delete: removes an image and mutates the gallery cache to remove the card.
   */
  async function handleDelete(imageId: string) {
    if (!activeProjectId) return

    try {
      await deleteImage({ id: imageId })

      // Optimistically remove the card from the SWR cache
      mutate(
        ['/images/list', { projectId: activeProjectId }],
        undefined,
        { revalidate: true }
      )
    } catch {
      // Error handling is left to the caller or global interceptor
    }
  }

  /**
   * Regenerate: opens the regenerate modal with the source image prefilled.
   */
  function handleRegenerate(image: Image) {
    openRegenerateModal(image)
  }

  return { retry, handleDelete, handleRegenerate }
}

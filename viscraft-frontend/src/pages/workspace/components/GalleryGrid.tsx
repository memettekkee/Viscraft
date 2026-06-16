import { useState } from 'react'
import { SimpleGrid } from '@chakra-ui/react'
import { useGallery } from '../hooks/useGallery'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { useImageActions } from '../hooks/useImageActions'
import { EmptyState } from '../../../components/common/EmptyState'
import { ConfirmationModal } from '../../../components/common/ConfirmationModal'
import { ImageCardSkeleton } from '../../../components/skeleton/ImageCardSkeleton'
import { ImageCard } from './ImageCard'
import type { Image } from '../../../types'

/**
 * Responsive grid of image cards for the active project.
 * Renders EmptyState when no images exist, skeleton cards when loading,
 * and ImageCard components for completed/failed images.
 *
 * Validates: Requirements 5.1, 11.1, 11.2, 11.3, 13.3
 */

interface GalleryGridProps {
  projectId: string
}

export function GalleryGrid({ projectId }: GalleryGridProps) {
  const { images, isLoading } = useGallery(projectId)
  const openGenerateModal = useWorkspaceStore((s) => s.openGenerateModal)
  const { handleRegenerate, handleDelete } = useImageActions()

  // Confirmation modal state for image deletion
  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  /** Open ConfirmationModal when user clicks Delete on a card */
  function onDeleteRequest(imageId: string) {
    setDeleteTargetId(imageId)
  }

  /** Close ConfirmationModal without deleting */
  function onDeleteCancel() {
    setDeleteTargetId(null)
  }

  /** Confirm deletion: call deleteImage, close modal */
  async function onDeleteConfirm() {
    if (!deleteTargetId) return
    setIsDeleting(true)
    await handleDelete(deleteTargetId)
    setIsDeleting(false)
    setDeleteTargetId(null)
  }

  // Loading state — show skeleton placeholders
  if (isLoading) {
    return (
      <SimpleGrid
        columns={{ base: 1, sm: 1, md: 2, lg: 3, xl: 4 }}
        gap="4"
      >
        {Array.from({ length: 6 }).map((_, i) => (
          <ImageCardSkeleton key={i} />
        ))}
      </SimpleGrid>
    )
  }

  // Empty state — no images in this project
  if (images.length === 0) {
    return <EmptyState onAction={openGenerateModal} />
  }

  // Render image grid + deletion confirmation modal
  return (
    <>
      <SimpleGrid
        columns={{ base: 1, sm: 1, md: 2, lg: 3, xl: 4 }}
        gap="4"
      >
        {images.map((image: Image) => {
          if (image.status === 'processing') {
            return <ImageCardSkeleton key={image.id} />
          }

          return (
            <ImageCard
              key={image.id}
              image={image}
              onRegenerate={handleRegenerate}
              onDelete={onDeleteRequest}
            />
          )
        })}
      </SimpleGrid>

      <ConfirmationModal
        isOpen={deleteTargetId !== null}
        onClose={onDeleteCancel}
        onConfirm={onDeleteConfirm}
        title="Delete Image"
        message="Are you sure you want to delete this image? This action cannot be undone."
        confirmLabel="Delete"
        isLoading={isDeleting}
      />
    </>
  )
}

import { SimpleGrid } from '@chakra-ui/react'
import { useGallery } from '../hooks/useGallery'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { useImageActions } from '../hooks/useImageActions'
import { EmptyState } from '../../../components/common/EmptyState'
import { ImageCardSkeleton } from '../../../components/skeleton/ImageCardSkeleton'
import { ImageCard } from './ImageCard'
import type { Image } from '../../../types'

/**
 * Responsive grid of image cards for the active project.
 * Renders EmptyState when no images exist, skeleton cards when loading,
 * and ImageCard components for completed/failed images.
 *
 * Validates: Requirements 5.1, 13.3
 */

interface GalleryGridProps {
  projectId: string
}

export function GalleryGrid({ projectId }: GalleryGridProps) {
  const { images, isLoading } = useGallery(projectId)
  const openGenerateModal = useWorkspaceStore((s) => s.openGenerateModal)
  const { handleRegenerate, handleDelete } = useImageActions()

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

  // Render image grid
  return (
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
            onDelete={handleDelete}
          />
        )
      })}
    </SimpleGrid>
  )
}

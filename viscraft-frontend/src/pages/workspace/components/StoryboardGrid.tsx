import { Box, SimpleGrid, Text } from '@chakra-ui/react'
import { useSceneList } from '../hooks/useSceneList'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { useSceneActions } from '../hooks/useSceneActions'
import { useDeleteConfirmation } from '../hooks/useDeleteConfirmation'
import { ConfirmationModal } from '../../../components/common/ConfirmationModal'
import { ImageCardSkeleton } from '../../../components/skeleton/ImageCardSkeleton'
import { SceneCard } from './SceneCard'
import { SceneDetailModal } from './SceneDetailModal'
import { PollingSceneCard } from './PollingSceneCard'
import { GenerateModal } from './GenerateModal'
import type { Scene } from '../../../types'

interface StoryboardGridProps {
  projectId: string
}

export function StoryboardGrid({ projectId }: StoryboardGridProps) {
  const { scenes, isLoading } = useSceneList(projectId)
  const openGenerateModal = useWorkspaceStore((s) => s.openGenerateModal)
  const generateModalOpen = useWorkspaceStore((s) => s.generateModalOpen)
  const closeModal = useWorkspaceStore((s) => s.closeModal)
  const { handleRegenerate } = useSceneActions()
  const {
    isDeleteModalOpen,
    isDeleting,
    onDeleteRequest,
    onDeleteCancel,
    onDeleteConfirm,
  } = useDeleteConfirmation()

  if (isLoading) {
    return (
      <SimpleGrid columns={{ base: 1, md: 2, lg: 3, xl: 4 }} gap="4">
        {Array.from({ length: 8 }).map((_, i) => (
          <ImageCardSkeleton key={i} />
        ))}
      </SimpleGrid>
    )
  }

  if (scenes.length === 0) {
    return (
      <>
        <SimpleGrid columns={{ base: 1, md: 2, lg: 3, xl: 4 }} gap="4">
          <GeneratePlaceholderCard onClick={() => openGenerateModal()} />
        </SimpleGrid>
        <GenerateModal isOpen={generateModalOpen} onClose={closeModal} />
      </>
    )
  }

  return (
    <>
      <SimpleGrid columns={{ base: 1, md: 2, lg: 3, xl: 4 }} gap="4">
        {scenes.map((scene: Scene) => {
          if (scene.status === 'processing') {
            return (
              <PollingSceneCard
                key={scene.id}
                scene={scene}
                scenes={scenes}
                projectId={projectId}
                onDelete={onDeleteRequest}
              />
            )
          }
          return (
            <SceneCard
              key={scene.id}
              scene={scene}
              scenes={scenes}
              onDelete={onDeleteRequest}
              onRegenerate={handleRegenerate}
            />
          )
        })}
        <GeneratePlaceholderCard onClick={() => openGenerateModal()} />
      </SimpleGrid>

      <SceneDetailModal onDelete={onDeleteRequest} onRegenerate={handleRegenerate} />

      <ConfirmationModal
        isOpen={isDeleteModalOpen}
        onClose={onDeleteCancel}
        onConfirm={onDeleteConfirm}
        title="Delete Ad Shot"
        message="Are you sure you want to delete this ad shot? This action cannot be undone."
        confirmLabel="Delete"
        isLoading={isDeleting}
      />
    </>
  )
}

function GeneratePlaceholderCard({ onClick }: { onClick: () => void }) {
  return (
    <Box
      as="button"
      w="100%"
      aspectRatio="4/3"
      border="2px dashed"
      borderColor="amber"
      borderRadius="md"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      gap="2"
      cursor="pointer"
      bg="transparent"
      opacity={0.6}
      _hover={{ opacity: 1, bg: 'rgba(201, 118, 44, 0.04)' }}
      transition="all 0.2s"
      onClick={onClick}
      aria-label="Generate new ad shot"
      data-tour="generate-card"
    >
      <Text fontSize="2xl" color="amber" aria-hidden="true">＋</Text>
      <Text fontFamily="body" fontSize="sm" color="amber" fontWeight="medium">
        Generate new ad shot
      </Text>
    </Box>
  )
}

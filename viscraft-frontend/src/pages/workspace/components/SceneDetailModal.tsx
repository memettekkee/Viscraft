import { Box, Button, HStack, Text, VStack } from '@chakra-ui/react'
import { ReusableModal } from '../../../components/ReusableModal'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import type { Scene } from '../../../types'

interface SceneDetailModalProps {
  onDelete: (sceneId: string) => void
  onRegenerate: (scene: Scene) => void
}

export function SceneDetailModal({ onDelete, onRegenerate }: SceneDetailModalProps) {
  const selectedScene = useWorkspaceStore((s) => s.selectedScene)
  const closeSceneDetail = useWorkspaceStore((s) => s.closeSceneDetail)

  if (!selectedScene) return null

  const sceneNumber = selectedScene.orderIndex + 1
  const imageUrl = selectedScene.fileUrl ?? undefined

  const statusConfig: Record<string, { bg: string; label: string }> = {
    completed: { bg: 'moss', label: 'Completed' },
    failed: { bg: 'oxblood', label: 'Failed' },
    processing: { bg: 'amber', label: 'Processing' },
  }

  const status = statusConfig[selectedScene.status] ?? { bg: 'warmgray', label: selectedScene.status }

  function handleDelete() {
    onDelete(selectedScene!.id)
    closeSceneDetail()
  }

  function handleRegenerate() {
    onRegenerate(selectedScene!)
    closeSceneDetail()
  }

  return (
    <ReusableModal
      isOpen={selectedScene !== null}
      onClose={closeSceneDetail}
      title={`Scene #${sceneNumber}`}
      size="lg"
    >
      <VStack gap="4" align="stretch" maxW="650px" mx="auto">
        {/* Image */}
        <Box
          w="100%"
          borderRadius="md"
          overflow="hidden"
          borderWidth="1px"
          borderColor="amber"
          bg="ink"
        >
          {imageUrl ? (
            <img
              src={imageUrl}
              alt={`Scene ${sceneNumber}: ${selectedScene.prompt}`}
              style={{
                width: '100%',
                maxHeight: '60vh',
                objectFit: 'contain',
                display: 'block',
              }}
            />
          ) : (
            <Box py="12" textAlign="center">
              <Text color="warmgray" fontSize="sm">No image available</Text>
            </Box>
          )}
        </Box>

        {/* Status row */}
        <HStack gap="3" align="center">
          <Box
            px="2"
            py="0.5"
            borderRadius="sm"
            bg={status.bg}
          >
            <Text fontSize="2xs" fontWeight="bold" color="white" textTransform="uppercase" fontFamily="mono">
              {status.label}
            </Text>
          </Box>
          <Text fontFamily="mono" fontSize="xs" color="warmgray">
            Scene #{sceneNumber}
          </Text>
        </HStack>

        {/* Prompt */}
        <Text fontFamily="mono" fontSize="sm" color="warmgray" lineHeight="tall">
          {selectedScene.generated_prompt}
        </Text>

        {/* Actions */}
        <HStack gap="3" justify="flex-end" pt="1" borderTop="1px solid" borderColor="rgba(201,118,44,0.2)">
          <Button
            variant="outline"
            size="sm"
            minH="36px"
            onClick={handleRegenerate}
          >
            Regenerate
          </Button>
          <Button
            variant="ghost"
            size="sm"
            minH="36px"
            color="oxblood"
            onClick={handleDelete}
          >
            Delete
          </Button>
        </HStack>
      </VStack>
    </ReusableModal>
  )
}

import { Box, Button, HStack, Text, VStack } from '@chakra-ui/react'
import { motion } from 'framer-motion'
import { ERROR_MESSAGES } from '../../../constants'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import type { Scene } from '../../../types'

interface SceneCardProps {
  scene: Scene
  scenes: Scene[]
  onDelete: (sceneId: string) => void
  onRegenerate: (scene: Scene) => void
}

function getErrorMessage(errorCode?: string): string {
  if (!errorCode) return 'Scene generation failed'
  return ERROR_MESSAGES[errorCode] || 'Scene generation failed'
}

const MotionBox = motion.create(Box)

export function SceneCard({ scene, onDelete, onRegenerate }: SceneCardProps) {
  const sceneNumber = scene.orderIndex + 1
  const openSceneDetail = useWorkspaceStore((s) => s.openSceneDetail)
  console.log("tesss",scene)

  // Failed state
  if (scene.status === 'failed') {
    const errorMessage = getErrorMessage(scene.errorCode)

    return (
      <Box
        bg="parchment"
        borderWidth="1px"
        borderColor="amber"
        borderRadius="md"
        overflow="hidden"
        w="100%"
        position="relative"
      >
        <Box
          aspectRatio="4/3"
          display="flex"
          alignItems="center"
          justifyContent="center"
          flexDirection="column"
          gap="3"
          px="4"
        >
          <Text fontSize="xl" aria-hidden="true">🗺️✕</Text>
          <Text fontFamily="mono" fontSize="xs" color="oxblood" textAlign="center">
            {errorMessage}
          </Text>
          <HStack gap="2">
            <Button
              variant="outline"
              size="xs"
              minH="32px"
              onClick={() => onRegenerate(scene)}
            >
              Retry
            </Button>
            <Button
              variant="ghost"
              size="xs"
              minH="32px"
              color="oxblood"
              onClick={() => onDelete(scene.id)}
            >
              Delete
            </Button>
          </HStack>
        </Box>
        {/* Stamp badge */}
        <StampBadge number={sceneNumber} />
      </Box>
    )
  }

  // Completed state
  const imageUrl = scene.fileUrl ?? undefined

  return (
    <MotionBox
      position="relative"
      bg="parchment"
      borderWidth="1px"
      borderColor="amber"
      borderRadius="md"
      overflow="hidden"
      w="100%"
      whileHover={{ scale: 1.02 }}
      transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      cursor="pointer"
      _hover={{ borderColor: 'ink' }}
    >
      {/* Image — clickable */}
      <Box
        position="relative"
        aspectRatio="4/3"
        overflow="hidden"
        onClick={() => openSceneDetail(scene)}
      >
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={`Scene ${sceneNumber}: ${scene.prompt}`}
            style={{
              width: '100%',
              height: '100%',
              objectFit: 'cover',
              display: 'block',
            }}
          />
        ) : (
          <Box
            w="100%"
            h="100%"
            bg="warmgray"
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            <Text color="parchment" fontSize="sm">No image</Text>
          </Box>
        )}

        {/* Cartographer stamp badge */}
        <StampBadge number={sceneNumber} />
      </Box>

      {/* Footer: prompt + actions */}
      <VStack gap="1" p="2.5" align="stretch">
        <Text
          fontFamily="mono"
          fontSize="xs"
          color="ink"
          lineHeight="short"
          lineClamp={2}
        >
          {scene.generated_prompt}
        </Text>

        <HStack gap="1" justify="flex-end" pt="1">
          <Button
            variant="ghost"
            size="xs"
            fontSize="2xs"
            minH="28px"
            px="2"
            color="amber"
            onClick={(e) => {
              e.stopPropagation()
              onRegenerate(scene)
            }}
          >
            Regenerate
          </Button>
          <Button
            variant="ghost"
            size="xs"
            fontSize="2xs"
            minH="28px"
            px="2"
            color="oxblood"
            onClick={(e) => {
              e.stopPropagation()
              onDelete(scene.id)
            }}
          >
            Delete
          </Button>
        </HStack>
      </VStack>
    </MotionBox>
  )
}

/**
 * Cartographer stamp badge — scene number in top-left corner,
 * slightly rotated with mono font and amber border.
 */
function StampBadge({ number }: { number: number }) {
  return (
    <Box
      position="absolute"
      top="2"
      left="2"
      bg="parchment"
      borderWidth="1px"
      borderColor="amber"
      borderRadius="sm"
      px="1.5"
      py="0.5"
      transform="rotate(-2deg)"
    >
      <Text
        fontFamily="mono"
        fontSize="2xs"
        fontWeight="bold"
        color="ink"
        lineHeight="1"
      >
        #{number}
      </Text>
    </Box>
  )
}

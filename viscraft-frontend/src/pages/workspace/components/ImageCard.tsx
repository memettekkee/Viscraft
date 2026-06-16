import { useState } from 'react'
import { Box, Button, HStack, Text, VStack } from '@chakra-ui/react'
import { ImageCardSkeleton } from '../../../components/skeleton/ImageCardSkeleton'
import { ERROR_MESSAGES } from '../../../constants'
import { useImageActions } from '../hooks/useImageActions'
import type { Image } from '../../../types'

/**
 * Renders a single image in one of its lifecycle states:
 * - processing → ImageCardSkeleton with shimmer
 * - completed → image display, truncated prompt, genre·mood stamp badge, actions
 * - failed → broken-map icon, error message, Retry button
 *
 * Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 14.3, 14.5
 */

interface ImageCardProps {
  image: Image
  onRegenerate: (image: Image) => void
  onDelete: (imageId: string) => void
}

const API_BASE_URL =
  window.__VISCRAFT_CONFIG__?.API_BASE_URL || 'http://localhost:8080'

/** Truncate a string to maxLen chars with ellipsis */
function truncatePrompt(prompt: string, maxLen = 60): string {
  if (prompt.length <= maxLen) return prompt
  return prompt.slice(0, maxLen).trimEnd() + '…'
}

/** Get user-facing error message for a given error code */
function getErrorMessage(errorCode?: string): string {
  if (!errorCode) return 'Image generation failed'
  return ERROR_MESSAGES[errorCode] || 'Image generation failed'
}

export function ImageCard({ image, onRegenerate, onDelete }: ImageCardProps) {
  const { retry } = useImageActions()
  const [isRetrying, setIsRetrying] = useState(false)

  // Processing state — show skeleton
  if (image.status === 'processing') {
    return <ImageCardSkeleton />
  }

  // Failed state
  if (image.status === 'failed') {
    const errorMessage = getErrorMessage(image.errorCode)

    async function handleRetry() {
      setIsRetrying(true)
      await retry(image)
      setIsRetrying(false)
    }

    return (
      <Box
        position="relative"
        bg="parchment"
        borderWidth="1px"
        borderColor="amber"
        borderRadius="md"
        overflow="hidden"
        aspectRatio="4/3"
        display="flex"
        alignItems="center"
        justifyContent="center"
        data-testid={`image-card-failed-${image.id}`}
      >
        <VStack gap="3">
          {/* Broken-map icon */}
          <Text fontSize="2xl" aria-hidden="true">
            🗺️✕
          </Text>

          {/* Error message */}
          <Text
            fontFamily="mono"
            fontSize="sm"
            color="oxblood"
            textAlign="center"
            px="4"
          >
            {errorMessage}
          </Text>

          {/* Action buttons */}
          <HStack gap="2">
            <Button
              variant="outline"
              size="sm"
              minH="44px"
              minW="44px"
              onClick={handleRetry}
              loading={isRetrying}
              disabled={isRetrying}
            >
              Retry
            </Button>
            <Button
              variant="ghost"
              size="sm"
              minH="44px"
              minW="44px"
              color="oxblood"
              onClick={() => onDelete(image.id)}
            >
              Delete
            </Button>
          </HStack>
        </VStack>
      </Box>
    )
  }

  // Completed state
  const imageUrl = image.fileUrl
    ? `${API_BASE_URL}${image.fileUrl}`
    : undefined

  const stampText = `${image.genre} · ${image.mood}`.toUpperCase()

  return (
    <Box
      position="relative"
      bg="parchment"
      borderWidth="1px"
      borderColor="amber"
      borderRadius="md"
      overflow="hidden"
      data-testid={`image-card-completed-${image.id}`}
    >
      {/* Image display */}
      <Box position="relative" aspectRatio="4/3" overflow="hidden">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={truncatePrompt(image.prompt, 100)}
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
            <Text color="parchment" fontSize="sm">
              No image
            </Text>
          </Box>
        )}

        {/* Stamp badge — genre · mood */}
        <Box
          position="absolute"
          top="2"
          right="2"
          bg="parchment"
          borderWidth="1px"
          borderColor="amber"
          borderRadius="sm"
          px="2"
          py="1"
          transform="rotate(-3deg)"
          data-testid={`stamp-badge-${image.id}`}
        >
          <Text
            fontFamily="mono"
            fontSize="xs"
            fontWeight="normal"
            color="ink"
            textTransform="uppercase"
            letterSpacing="wider"
            lineHeight="1"
          >
            {stampText}
          </Text>
        </Box>
      </Box>

      {/* Card footer: prompt + actions */}
      <VStack gap="2" p="3" align="stretch">
        {/* Truncated prompt */}
        <Text
          fontFamily="body"
          fontSize="sm"
          color="ink"
          lineHeight="short"
        >
          {truncatePrompt(image.prompt)}
        </Text>

        {/* Action buttons */}
        <HStack gap="2" justify="flex-end">
          <Button
            variant="outline"
            size="sm"
            minH="44px"
            minW="44px"
            onClick={() => onRegenerate(image)}
          >
            Regenerate
          </Button>
          <Button
            variant="ghost"
            size="sm"
            minH="44px"
            minW="44px"
            color="oxblood"
            onClick={() => onDelete(image.id)}
          >
            Delete
          </Button>
        </HStack>
      </VStack>
    </Box>
  )
}

import { Button, Text, VStack } from '@chakra-ui/react'

/**
 * Empty state placeholder displayed when a project has no images
 * or no project is selected. Provides an inviting call-to-action
 * to guide the user toward generating their first image.
 *
 * Validates: Requirements 11.1
 */

export interface EmptyStateProps {
  onAction: () => void
  title?: string
  description?: string
}

export function EmptyState({
  onAction,
  title = 'No maps charted yet',
  description = 'Start generating concept art to fill your collection.',
}: EmptyStateProps) {
  return (
    <VStack
      gap="4"
      py="16"
      px="6"
      align="center"
      justify="center"
      textAlign="center"
    >
      <Text
        fontFamily="display"
        fontSize="2xl"
        color="ink"
        fontWeight="medium"
      >
        {title}
      </Text>
      <Text
        fontFamily="body"
        fontSize="md"
        color="warmgray"
        maxW="sm"
      >
        {description}
      </Text>
      <Button
        variant="solid"
        onClick={onAction}
        mt="2"
        minW="44px"
        minH="44px"
      >
        Generate your first image
      </Button>
    </VStack>
  )
}

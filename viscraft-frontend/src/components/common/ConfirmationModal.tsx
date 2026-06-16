import { Button, HStack, Text, VStack } from '@chakra-ui/react'
import { ReusableModal } from '../ReusableModal'

/**
 * Confirmation dialog for destructive actions (project deletion, image deletion).
 *
 * Wraps ReusableModal with confirm/cancel action buttons and an optional loading state
 * on the confirm button to prevent double-submission.
 *
 * Validates: Requirements 4.5, 11.1
 */

export interface ConfirmationModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  message: string
  confirmLabel?: string
  isLoading?: boolean
}

export function ConfirmationModal({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmLabel = 'Delete',
  isLoading = false,
}: ConfirmationModalProps) {
  return (
    <ReusableModal isOpen={isOpen} onClose={onClose} title={title} size={{ base: 'full', md: 'md' }}>
      <VStack gap="6" align="stretch">
        <Text fontFamily="body" color="ink" fontSize="md">
          {message}
        </Text>
        <HStack gap="3" justify="flex-end">
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isLoading}
            minW="44px"
            minH="44px"
          >
            Cancel
          </Button>
          <Button
            variant="solid"
            bg="oxblood"
            color="white"
            _hover={{ opacity: 0.9 }}
            onClick={onConfirm}
            loading={isLoading}
            disabled={isLoading}
            minW="44px"
            minH="44px"
          >
            {confirmLabel}
          </Button>
        </HStack>
      </VStack>
    </ReusableModal>
  )
}

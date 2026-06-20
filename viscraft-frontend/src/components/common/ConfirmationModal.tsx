import { Button, HStack, Text, VStack } from '@chakra-ui/react'
import { ReusableModal } from '../ReusableModal'

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
    <ReusableModal isOpen={isOpen} onClose={onClose} title={title} size="sm">
      <VStack gap="4" align="stretch">
        <Text fontFamily="body" color="ink" fontSize="sm">
          {message}
        </Text>
        <HStack gap="3" justify="flex-end">
          <Button
            variant="outline"
            size="sm"
            onClick={onClose}
            disabled={isLoading}
            minW="44px"
            minH="36px"
          >
            Cancel
          </Button>
          <Button
            variant="solid"
            size="sm"
            bg="oxblood"
            color="white"
            _hover={{ opacity: 0.9 }}
            onClick={onConfirm}
            loading={isLoading}
            disabled={isLoading}
            minW="44px"
            minH="36px"
          >
            {confirmLabel}
          </Button>
        </HStack>
      </VStack>
    </ReusableModal>
  )
}

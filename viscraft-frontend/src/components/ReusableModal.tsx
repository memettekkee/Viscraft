import {
  DialogRoot,
  DialogContent,
  DialogHeader,
  DialogBody,
  DialogBackdrop,
  DialogPositioner,
  IconButton,
  Box,
} from '@chakra-ui/react'
import type { ReactNode } from 'react'

export interface ReusableModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  children: ReactNode
}

export function ReusableModal({
  isOpen,
  onClose,
  title,
  size = 'md',
  children,
}: ReusableModalProps) {
  const widthMap = {
    sm: '24rem',
    md: '32rem',
    lg: '40rem',
    xl: '48rem',
  }

  return (
    <DialogRoot
      open={isOpen}
      onOpenChange={({ open }) => {
        if (!open) onClose()
      }}
      placement="center"
      size={size}
    >
      <DialogBackdrop
        bg="blackAlpha.700"
        backdropFilter="blur(4px)"
      />
      <DialogPositioner>
        <DialogContent
          width={{ base: '100vw', md: widthMap[size] }}
          height={{ base: '100vh', md: 'auto' }}
          maxHeight={{ base: '100vh', md: '85vh' }}
          borderRadius={{ base: '0', md: 'lg' }}
          overflow="auto"
          position="relative"
        >
          {/* Close button - always visible */}
          <Box position="absolute" top="3" right="3" zIndex="10">
            <IconButton
              aria-label="Close dialog"
              onClick={onClose}
              variant="ghost"
              size="sm"
              minW="36px"
              minH="36px"
              borderRadius="full"
              _hover={{ bg: 'blackAlpha.100' }}
            >
              ✕
            </IconButton>
          </Box>

          <DialogHeader fontFamily="display" fontSize="xl" pr="12">
            {title}
          </DialogHeader>
          <DialogBody pb="6">{children}</DialogBody>
        </DialogContent>
      </DialogPositioner>
    </DialogRoot>
  )
}

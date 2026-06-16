import {
  DialogRoot,
  DialogContent,
  DialogHeader,
  DialogBody,
  DialogCloseTrigger,
  DialogBackdrop,
} from '@chakra-ui/react'
import type { ReactNode } from 'react'

/**
 * Reusable modal component wrapping Chakra UI v3 Dialog.
 *
 * Responsive behavior:
 *  - Mobile (base): renders as a full-screen overlay
 *  - Desktop (md+): renders as a centered modal
 *
 * Theme: uses the Cartographer's Atlas dialog slot recipe —
 * parchment background, amber border, ink text, Fraunces display font for title.
 *
 * Validates: Requirements 13.4, 13.5
 */

export interface ReusableModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  size?: { base: string; md: string }
  children: ReactNode
}

export function ReusableModal({
  isOpen,
  onClose,
  title,
  size,
  children,
}: ReusableModalProps) {
  const _size = size ?? { base: 'full', md: 'lg' }

  return (
    <DialogRoot
      open={isOpen}
      onOpenChange={({ open }) => {
        if (!open) onClose()
      }}
      placement="center"
      size={_size.md as 'sm' | 'md' | 'lg' | 'xl' | 'full'}
    >
      <DialogBackdrop />
      <DialogContent
        width={{ base: '100vw', md: 'auto' }}
        height={{ base: '100vh', md: 'auto' }}
        maxHeight={{ base: '100vh', md: '85vh' }}
        maxWidth={{ md: _size.md === 'lg' ? '32rem' : undefined }}
        borderRadius={{ base: '0', md: 'md' }}
        overflow="auto"
      >
        <DialogHeader fontFamily="display" fontSize="xl">
          {title}
        </DialogHeader>
        <DialogBody>{children}</DialogBody>
        <DialogCloseTrigger aria-label="Close dialog" />
      </DialogContent>
    </DialogRoot>
  )
}

import { Box, Flex, Text, IconButton } from '@chakra-ui/react'
import {
  Toaster,
  createToaster,
  type CreateToasterReturn,
} from '@chakra-ui/react'
import type { ReactNode } from 'react'

/**
 * Viscraft toast type variants mapped to design system colors.
 * - error: oxblood (#8B2E2E) — session expired, resource not found
 * - success: moss (#3E5C4E) — operation confirmations
 * - info: amber (#C9762C) — informational/warning notifications
 *
 * Validates: Requirements 12.3, 12.4, 12.5
 */
export type ToastType = 'error' | 'success' | 'info'

/** Design system color mapping for toast variants */
const toastColors: Record<ToastType, { bg: string; border: string; icon: string }> = {
  error: { bg: '#8B2E2E', border: '#8B2E2E', icon: '✕' },
  success: { bg: '#3E5C4E', border: '#3E5C4E', icon: '✓' },
  info: { bg: '#C9762C', border: '#C9762C', icon: 'ℹ' },
}

/**
 * Global toaster instance configured for Viscraft.
 * Placement: top-end (top-right) — works well for desktop.
 * Mobile positioning is handled via responsive styles in the Toaster render.
 */
export const toaster: CreateToasterReturn = createToaster({
  placement: 'top-end',
  max: 5,
  gap: 12,
  offsets: '1rem',
})

interface CustomToastContentProps {
  type: ToastType
  title: string
  description?: string
  onClose?: () => void
}

/**
 * CustomToastContent renders the toast body with design system styling.
 * Shows a colored indicator, message content, and a close button.
 */
function CustomToastContent({ type, title, description, onClose }: CustomToastContentProps) {
  const colors = toastColors[type]

  return (
    <Box
      bg={colors.bg}
      color="#FAF6EC"
      borderWidth="1px"
      borderColor={colors.border}
      borderRadius="sm"
      px="4"
      py="3"
      minW={{ base: '280px', md: '320px' }}
      maxW={{ base: '90vw', md: '400px' }}
      fontFamily="'Inter', sans-serif"
      role="alert"
      aria-live="assertive"
    >
      <Flex align="flex-start" gap="3">
        {/* Type indicator icon */}
        <Flex
          align="center"
          justify="center"
          w="6"
          h="6"
          borderRadius="full"
          bg="whiteAlpha.200"
          flexShrink={0}
          mt="0.5"
          fontSize="xs"
          fontWeight="bold"
          aria-hidden="true"
        >
          {colors.icon}
        </Flex>

        {/* Content */}
        <Box flex="1" minW="0">
          <Text fontSize="sm" fontWeight="medium" lineHeight="short">
            {title}
          </Text>
          {description && (
            <Text fontSize="xs" opacity={0.85} mt="1" lineHeight="short">
              {description}
            </Text>
          )}
        </Box>

        {/* Close button */}
        <IconButton
          aria-label="Dismiss notification"
          size="xs"
          variant="ghost"
          color="#FAF6EC"
          opacity={0.7}
          _hover={{ opacity: 1, bg: 'whiteAlpha.200' }}
          onClick={onClose}
          minW="6"
          h="6"
          flexShrink={0}
        >
          ✕
        </IconButton>
      </Flex>
    </Box>
  )
}

/**
 * ViscraftToaster is the Toaster component that must be rendered once in the app tree.
 * It uses the global `toaster` instance and renders each toast using CustomToastContent.
 */
export function ViscraftToaster() {
  return (
    <Toaster
      toaster={toaster}
    >
      {(toast) => {
        const type = (toast.meta?.toastType as ToastType) ?? 'info'
        return (
          <CustomToastContent
            type={type}
            title={(toast.title as string) ?? ''}
            description={toast.description as string | undefined}
            onClose={() => toaster.dismiss(toast.id)}
          />
        )
      }}
    </Toaster>
  )
}

/**
 * Helper to create a toast via the global toaster store.
 * Used internally by the useToast hook.
 */
export function showToast(options: {
  type: ToastType
  title: string
  description?: string
  duration?: number
}) {
  const { type, title, description, duration } = options
  const defaultDuration = type === 'error' ? 5000 : 3000

  toaster.create({
    title: title as ReactNode,
    description: description as ReactNode,
    duration: duration ?? defaultDuration,
    type: type === 'info' ? 'info' : type,
    meta: { toastType: type },
  })
}

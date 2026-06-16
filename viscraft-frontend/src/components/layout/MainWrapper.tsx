import { Box } from '@chakra-ui/react'

interface MainWrapperProps {
  children: React.ReactNode
}

/**
 * Scrollable content wrapper for the main area.
 * Takes full remaining height and enables vertical scrolling.
 *
 * Validates: Requirements 13.1, 13.2
 */
export function MainWrapper({ children }: MainWrapperProps) {
  return (
    <Box
      flex="1"
      overflowY="auto"
      height="100%"
      bg="surface.bg"
      borderRadius={{ base: '0', md: 'md' }}
      p="4"
    >
      {children}
    </Box>
  )
}

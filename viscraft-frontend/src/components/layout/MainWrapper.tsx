import { Box } from '@chakra-ui/react'

interface MainWrapperProps {
  children: React.ReactNode
}

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

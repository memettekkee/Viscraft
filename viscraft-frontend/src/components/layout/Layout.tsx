import { Box, Flex } from '@chakra-ui/react'
import { ProjectSidebar } from '../sidebar/ProjectSidebar'
import { MainWrapper } from './MainWrapper'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  return (
    <Flex
      direction={{ base: 'column', md: 'row' }}
      height="100vh"
      bg="shell.bg"
    >
      <Box
        width={{ base: '100%', md: '240px' }}
        minWidth={{ base: 'unset', md: '240px' }}
        height={{ base: 'auto', md: '100vh' }}
        overflowX={{ base: 'auto', md: 'hidden' }}
        overflowY={{ base: 'hidden', md: 'auto' }}
        whiteSpace={{ base: 'nowrap', md: 'normal' }}
        borderRight={{ base: 'none', md: '1px solid' }}
        borderBottom={{ base: '1px solid', md: 'none' }}
        borderColor="border.accent"
        flexShrink={0}
      >
        <ProjectSidebar />
      </Box>

      {/* Main content area */}
      <MainWrapper>{children}</MainWrapper>
    </Flex>
  )
}

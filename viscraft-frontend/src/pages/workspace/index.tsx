import { Box, Button, Flex, Text } from '@chakra-ui/react'
import { useWorkspaceStore } from '../../store/workspaceStore'
import { GalleryGrid } from './components/GalleryGrid'

/**
 * Workspace page shell with gallery area and floating Generate button.
 * Shows EmptyState when no project is selected; otherwise renders the GalleryGrid placeholder.
 *
 * Validates: Requirements 4.1, 4.4
 */
export function WorkspacePage() {
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const openGenerateModal = useWorkspaceStore((s) => s.openGenerateModal)

  if (!activeProjectId) {
    return (
      <Flex
        direction="column"
        align="center"
        justify="center"
        minH="100%"
        flex="1"
        px={{ base: '4', md: '8' }}
        py="16"
      >
        <Text
          fontFamily="display"
          fontSize={{ base: '2xl', md: '3xl' }}
          color="parchment"
          mb="3"
        >
          No region selected
        </Text>
        <Text
          fontFamily="body"
          fontSize={{ base: 'sm', md: 'md' }}
          color="warmgray"
          textAlign="center"
          maxW="400px"
        >
          Select a project from the sidebar or create a new region to start
          generating concept art.
        </Text>
      </Flex>
    )
  }

  return (
    <Box position="relative" minH="100%" flex="1" p={{ base: '4', md: '6' }}>
      {/* GalleryGrid — fetches and displays images for the active project */}
      <GalleryGrid projectId={activeProjectId} />

      {/* Floating Generate button */}
      <Button
        position="fixed"
        bottom={{ base: '6', md: '8' }}
        right={{ base: '6', md: '8' }}
        bg="amber"
        color="white"
        fontFamily="body"
        fontWeight="medium"
        fontSize="md"
        px="6"
        py="5"
        borderRadius="sm"
        zIndex="10"
        _hover={{ opacity: 0.9 }}
        _focusVisible={{
          outline: '2px solid',
          outlineColor: 'amber',
          outlineOffset: '2px',
        }}
        onClick={openGenerateModal}
        aria-label="Generate new concept art"
      >
        Generate
      </Button>
    </Box>
  )
}

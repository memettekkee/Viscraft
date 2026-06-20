import { Box, Flex, Text } from '@chakra-ui/react'
import useSWR from 'swr'
import { useWorkspaceStore } from '../../store/workspaceStore'
import { StoryboardGrid } from './components/StoryboardGrid'
import { GenerateModal } from './components/GenerateModal'
import { OnboardingTour } from '../../components/OnboardingTour'
import { useSeedSampleCampaign } from './hooks/useSeedSampleCampaign'
import { postFetcher } from '../../helper/fetcher'
import type { ApiResponse, Project } from '../../types'

export function WorkspacePage() {
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const generateModalOpen = useWorkspaceStore((s) => s.generateModalOpen)
  const closeModal = useWorkspaceStore((s) => s.closeModal)

  // Auto-seed sample campaign for new users
  useSeedSampleCampaign()

  const { data: projectData } = useSWR<ApiResponse<Project>>(
    activeProjectId ? ['/projects/get', { id: activeProjectId }] : null,
    postFetcher
  )
  const project = projectData?.data

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
        <Text fontFamily="display" fontSize={{ base: '2xl', md: '3xl' }} color="parchment" mb="3">
          No region selected
        </Text>
        <Text fontFamily="body" fontSize={{ base: 'sm', md: 'md' }} color="warmgray" textAlign="center" maxW="400px">
          Select a project from the sidebar or create a new region to start generating concept art.
        </Text>
      </Flex>
    )
  }

  return (
    <Box
      position="relative"
      minH="100%"
      flex="1"
      p={{ base: '4', md: '6' }}
      _before={{
        content: '""',
        position: 'absolute',
        inset: 0,
        opacity: 0.02,
        backgroundImage: 'radial-gradient(circle, #B8860B 1px, transparent 1px)',
        backgroundSize: '20px 20px',
        pointerEvents: 'none',
      }}
    >
      {/* Project identity header */}
      {project && (
        <Flex align="center" gap="2" mb="5" wrap="wrap">
          <Text fontFamily="display" fontSize="xl" color="black" fontWeight="medium">
            {project.name}
          </Text>
          <Flex gap="2" align="center" ml="2">
            {project.productCategory && (
              <Box
                px="2.5"
                py="0.5"
                bg="rgba(201,118,44,0.12)"
                borderWidth="1px"
                borderColor="amber"
                borderRadius="full"
              >
                <Text fontFamily="mono" fontSize="2xs" color="amber" textTransform="uppercase" letterSpacing="wider">
                  {project.productCategory}
                </Text>
              </Box>
            )}
            {project.visualStyle && (
              <Box
                px="2.5"
                py="0.5"
                bg="rgba(107,101,85,0.1)"
                borderWidth="1px"
                borderColor="warmgray"
                borderRadius="full"
              >
                <Text fontFamily="mono" fontSize="2xs" color="warmgray" textTransform="uppercase" letterSpacing="wider">
                  {project.visualStyle}
                </Text>
              </Box>
            )}
          </Flex>
        </Flex>
      )}

      <StoryboardGrid projectId={activeProjectId} />

      <GenerateModal isOpen={generateModalOpen} onClose={closeModal} />
      <OnboardingTour />
    </Box>
  )
}

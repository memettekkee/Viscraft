import { useState } from 'react'
import { Box, Button, Flex, HStack, Text, VStack } from '@chakra-ui/react'
import useSWR from 'swr'
import { useSWRConfig } from 'swr'
import { postFetcher } from '../../helper/fetcher'
import { useWorkspaceStore } from '../../store/workspaceStore'
import { useAuthStore } from '../../store/authStore'
import { toRomanNumeral } from '../../pages/workspace/utils/romanNumeral'
import { deleteProject } from '../../service/project'
import { showToast } from '../CustomToast'
import { ERROR_MESSAGES } from '../../constants'
import { ConfirmationModal } from '../common/ConfirmationModal'
import { ProjectModal } from '../../pages/workspace/components/ProjectModal'
import type { AxiosError } from 'axios'
import type { ApiResponse, Project } from '../../types'

export function ProjectSidebar() {
  const { mutate } = useSWRConfig()
  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const setActiveProject = useWorkspaceStore((s) => s.setActiveProject)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const user = useAuthStore((s) => s.user)

  const { data, isLoading } = useSWR<ApiResponse<Project[]>>(
    ['/projects/list'],
    postFetcher
  )

  const projects = data?.data ?? []

  const [projectModalOpen, setProjectModalOpen] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Project | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    if (!deleteTarget) return
    setIsDeleting(true)
    try {
      await deleteProject({ id: deleteTarget.id })
      await mutate(['/projects/list'])
      if (activeProjectId === deleteTarget.id) {
        const remaining = projects.filter((p) => p.id !== deleteTarget.id)
        if (remaining.length > 0) {
          setActiveProject(remaining[0].id)
        }
      }
    } catch (err: unknown) {
      const axiosError = err as AxiosError<ApiResponse>
      const code = axiosError.response?.data?.errorCode
      const backendMessage = axiosError.response?.data?.message
      const message = backendMessage ?? (code ? (ERROR_MESSAGES[code] ?? 'An error occurred') : ERROR_MESSAGES.NETWORK_ERROR)
      showToast({ type: 'error', title: message })
    } finally {
      setIsDeleting(false)
      setDeleteTarget(null)
    }
  }

  return (
    <>
      <VStack
        display={{ base: 'none', md: 'flex' }}
        align="stretch"
        gap="1"
        p="3"
        height="100%"
        data-tour="sidebar"
      >
        <Text
          fontFamily="display"
          fontSize="sm"
          color="warmgray"
          textTransform="uppercase"
          letterSpacing="wider"
          mb="2"
          px="2"
        >
          Campaigns
        </Text>

        {isLoading && (
          <Text fontFamily="body" fontSize="sm" color="warmgray" px="2">
            Loading…
          </Text>
        )}

        <VStack align="stretch" gap="0" flex="1" overflow="auto">
          {projects.map((project, index) => {
            const isActive = project.id === activeProjectId
            return (
              <Flex
                key={project.id}
                align="center"
                justify="space-between"
                px="2"
                py="2"
                borderRadius="sm"
                cursor="pointer"
                bg={isActive ? 'rgba(201, 118, 44, 0.12)' : 'transparent'}
                borderLeft="3px solid"
                borderColor={isActive ? 'amber' : 'transparent'}
                _hover={{ bg: 'rgba(201, 118, 44, 0.08)' }}
                onClick={() => setActiveProject(project.id)}
                role="button"
                aria-label={`Select project ${project.name}`}
                aria-current={isActive ? 'true' : undefined}
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault()
                    setActiveProject(project.id)
                  }
                }}
              >
                <Text
                  fontFamily="display"
                  fontSize="sm"
                  color={isActive ? 'parchment' : 'warmgray'}
                  fontWeight={isActive ? 'semibold' : 'normal'}
                  truncate
                  flex="1"
                >
                  {toRomanNumeral(index + 1)}. {project.name}
                </Text>
                <Button
                  variant="plain"
                  size="xs"
                  color="warmgray"
                  minW={{ base: '44px', md: '32px' }}
                  minH={{ base: '44px', md: '32px' }}
                  p="0"
                  _hover={{ color: 'oxblood' }}
                  onClick={(e) => {
                    e.stopPropagation()
                    setDeleteTarget(project)
                  }}
                  aria-label={`Delete project ${project.name}`}
                >
                  ✕
                </Button>
              </Flex>
            )
          })}
        </VStack>

        <Button
          variant="outline"
          size="sm"
          mt="3"
          fontFamily="body"
          fontSize="sm"
          onClick={() => setProjectModalOpen(true)}
          minH="44px"
        >
          + New Campaign
        </Button>

        {/* User info + Logout */}
        <Box mt="auto" pt="4" borderTop="1px solid" borderColor="border.accent">
          {user && (
            <Text fontFamily="body" fontSize="xs" color="warmgray" px="2" mb="2" truncate>
              {user.email}
            </Text>
          )}
          <Button
            variant="ghost"
            size="sm"
            width="100%"
            fontFamily="body"
            fontSize="sm"
            color="oxblood"
            minH="44px"
            _hover={{ bg: 'rgba(139, 0, 0, 0.08)' }}
            onClick={() => {
              clearAuth()
              window.location.href = '/'
            }}
          >
            Log out
          </Button>
        </Box>
      </VStack>

      {/* Mobile layout: horizontal scrollable chip bar */}
      <HStack
        display={{ base: 'flex', md: 'none' }}
        gap="2"
        px="3"
        py="2"
        overflow="visible"
        flexWrap="nowrap"
      >
        {isLoading && (
          <Text fontFamily="body" fontSize="xs" color="warmgray" flexShrink={0}>
            Loading…
          </Text>
        )}

        {projects.map((project, index) => {
          const isActive = project.id === activeProjectId
          return (
            <Box
              key={project.id}
              as="button"
              flexShrink={0}
              px="3"
              minH="44px"
              display="flex"
              alignItems="center"
              gap="1"
              borderRadius="full"
              bg={isActive ? 'amber' : 'transparent'}
              borderWidth="1px"
              borderColor={isActive ? 'amber' : 'warmgray'}
              cursor="pointer"
              onClick={() => setActiveProject(project.id)}
              _hover={{ borderColor: 'amber' }}
              aria-label={`Select project ${project.name}`}
              aria-current={isActive ? 'true' : undefined}
            >
              <Text
                fontFamily="display"
                fontSize="xs"
                color={isActive ? 'white' : 'parchment'}
                whiteSpace="nowrap"
              >
                {toRomanNumeral(index + 1)}. {project.name}
              </Text>

              {/* Delete button shown on active chip */}
              {isActive && (
                <Box
                  as="span"
                  ml="1"
                  px="1"
                  fontFamily="body"
                  fontSize="xs"
                  fontWeight="bold"
                  color="white"
                  lineHeight="1"
                  borderRadius="full"
                  _hover={{ bg: 'rgba(255,255,255,0.2)' }}
                  onClick={(e: React.MouseEvent) => {
                    e.stopPropagation()
                    setDeleteTarget(project)
                  }}
                  aria-label={`Delete project ${project.name}`}
                >
                  ✕
                </Box>
              )}
            </Box>
          )
        })}

        <Box
          as="button"
          flexShrink={0}
          px="3"
          minH="44px"
          display="flex"
          alignItems="center"
          borderRadius="full"
          borderWidth="1px"
          borderColor="amber"
          bg="transparent"
          cursor="pointer"
          onClick={() => setProjectModalOpen(true)}
          aria-label="Create new project"
        >
          <Text fontFamily="body" fontSize="xs" color="amber" whiteSpace="nowrap">
            + New
          </Text>
        </Box>
      </HStack>

      {/* Project creation modal */}
      <ProjectModal
        isOpen={projectModalOpen}
        onClose={() => setProjectModalOpen(false)}
      />

      {/* Delete confirmation modal */}
      <ConfirmationModal
        isOpen={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        title="Delete Campaign"
        message={`Are you sure you want to delete "${deleteTarget?.name ?? ''}"? All ad shots in this campaign will be permanently removed.`}
        confirmLabel="Delete"
        isLoading={isDeleting}
      />
    </>
  )
}

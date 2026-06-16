import { useState } from 'react'
import { Box, Button, Input, Text, Textarea } from '@chakra-ui/react'
import { useSWRConfig } from 'swr'
import { ReusableModal } from '../../../components/ReusableModal'
import { createProject } from '../../../service/project'
import { useWorkspaceStore } from '../../../store/workspaceStore'

/**
 * ProjectModal — create a new project (region).
 *
 * Name field: required, 1-255 characters.
 * Description field: optional.
 * On success: calls createProject, mutates SWR project list, sets new project as active, closes modal.
 *
 * Validates: Requirements 4.2, 4.3
 */

export interface ProjectModalProps {
  isOpen: boolean
  onClose: () => void
}

export function ProjectModal({ isOpen, onClose }: ProjectModalProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [nameError, setNameError] = useState<string | null>(null)
  const [apiError, setApiError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const { mutate } = useSWRConfig()
  const setActiveProject = useWorkspaceStore((s) => s.setActiveProject)

  function resetForm() {
    setName('')
    setDescription('')
    setNameError(null)
    setApiError(null)
    setIsSubmitting(false)
  }

  function handleClose() {
    resetForm()
    onClose()
  }

  function validateName(value: string): string | null {
    const trimmed = value.trim()
    if (trimmed.length === 0) {
      return 'Project name is required'
    }
    if (trimmed.length > 255) {
      return 'Project name must not exceed 255 characters'
    }
    return null
  }

  async function handleSubmit() {
    const error = validateName(name)
    if (error) {
      setNameError(error)
      return
    }

    setNameError(null)
    setApiError(null)
    setIsSubmitting(true)

    try {
      const response = await createProject({
        name: name.trim(),
        description: description.trim() || undefined,
      })

      if (response.success && response.data) {
        // Mutate SWR project list cache to refetch
        await mutate(['/projects/list'])
        // Set new project as active
        setActiveProject(response.data.id)
        // Close and reset
        handleClose()
      } else {
        setApiError(response.message || 'Failed to create project')
      }
    } catch {
      setApiError('Unable to connect to server')
    } finally {
      setIsSubmitting(false)
    }
  }

  const isValid = name.trim().length >= 1 && name.trim().length <= 255

  return (
    <ReusableModal isOpen={isOpen} onClose={handleClose} title="New Region">
      {/* API error */}
      {apiError && (
        <Box
          bg="red.50"
          border="1px solid"
          borderColor="oxblood"
          borderRadius="sm"
          p="3"
          mb="4"
        >
          <Text color="oxblood" fontSize="sm" fontFamily="mono">
            {apiError}
          </Text>
        </Box>
      )}

      {/* Name field */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Name
        </Text>
        <Input
          type="text"
          placeholder="Enter project name"
          value={name}
          minHeight="44px"
          onChange={(e) => {
            setName(e.target.value)
            setNameError(null)
          }}
          borderColor={nameError ? 'oxblood' : undefined}
          aria-required="true"
          aria-invalid={!!nameError}
          aria-describedby={nameError ? 'name-error' : undefined}
        />
        {nameError && (
          <Text id="name-error" color="oxblood" fontSize="xs" mt="1">
            {nameError}
          </Text>
        )}
      </Box>

      {/* Description field */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Description (optional)
        </Text>
        <Textarea
          placeholder="Describe your project"
          value={description}
          minHeight="80px"
          onChange={(e) => setDescription(e.target.value)}
        />
      </Box>

      {/* Submit button */}
      <Button
        width="full"
        variant="solid"
        mt="2"
        minHeight="44px"
        disabled={!isValid || isSubmitting}
        onClick={handleSubmit}
      >
        {isSubmitting ? 'Creating...' : 'Create Region'}
      </Button>
    </ReusableModal>
  )
}

import { useRef, useState } from 'react'
import { Box, Button, Flex, Text, Textarea, Image } from '@chakra-ui/react'
import { ReusableModal } from '../../../components/ReusableModal'
import { useGenerateForm } from '../hooks/useGenerateForm'
import { useGallery } from '../hooks/useGallery'
import { useWorkspaceStore } from '../../../store/workspaceStore'
import { GENRE_OPTIONS, ASSET_TYPE_OPTIONS, MOOD_OPTIONS } from '../../../constants'
import { validateImageFile, fileToBase64, imageUrlToBase64 } from '../utils/referenceImage'
import type { Genre, AssetType, Mood } from '../../../types'

const MAX_PROMPT_LENGTH = 300

/**
 * GenerateModal — guided form for creating concept art.
 * Supports "Create" and "From Reference" modes.
 * UI-only component — logic lives in useGenerateForm hook.
 *
 * Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 12.1, 13.4, 13.5
 */

interface GenerateModalProps {
  isOpen: boolean
  onClose: () => void
}

export function GenerateModal({ isOpen, onClose }: GenerateModalProps) {
  const {
    mode,
    prompt,
    setPrompt,
    genre,
    setGenre,
    assetType,
    setAssetType,
    mood,
    setMood,
    referenceImage,
    setReferenceImage,
    errors,
    setErrors,
    apiError,
    isSubmitting,
    handleModeSwitch,
    handleSubmit,
    resetForm,
  } = useGenerateForm()

  const activeProjectId = useWorkspaceStore((s) => s.activeProjectId)
  const { images } = useGallery(activeProjectId)
  const completedImages = images.filter((img) => img.status === 'completed')

  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFileName, setSelectedFileName] = useState<string | null>(null)
  const [isLoadingReference, setIsLoadingReference] = useState(false)

  function handleClose() {
    resetForm()
    setSelectedFileName(null)
    onClose()
  }

  async function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return

    // Validate file type and size before processing
    const validationError = validateImageFile(file)
    if (validationError) {
      setErrors((prev) => ({ ...prev, referenceImage: validationError.message }))
      // Reset the input so the same file can be re-selected
      if (fileInputRef.current) fileInputRef.current.value = ''
      return
    }

    // Clear any previous reference image error
    setErrors((prev) => {
      const next = { ...prev }
      delete next.referenceImage
      return next
    })

    try {
      setIsLoadingReference(true)
      const base64 = await fileToBase64(file)
      setReferenceImage(base64)
      setSelectedFileName(file.name)
    } catch {
      setErrors((prev) => ({ ...prev, referenceImage: 'Failed to read file' }))
    } finally {
      setIsLoadingReference(false)
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  async function handleGalleryImageSelect(fileUrl: string) {
    setErrors((prev) => {
      const next = { ...prev }
      delete next.referenceImage
      return next
    })

    try {
      setIsLoadingReference(true)
      const base64 = await imageUrlToBase64(fileUrl)
      setReferenceImage(base64)
      setSelectedFileName(null) // Gallery images don't have a file name
    } catch {
      setErrors((prev) => ({ ...prev, referenceImage: 'Failed to load gallery image' }))
    } finally {
      setIsLoadingReference(false)
    }
  }

  function handleClearReference() {
    setReferenceImage(undefined)
    setSelectedFileName(null)
    setErrors((prev) => {
      const next = { ...prev }
      delete next.referenceImage
      return next
    })
  }

  return (
    <ReusableModal isOpen={isOpen} onClose={handleClose} title="Generate Concept Art">
      {/* Mode toggle */}
      <Flex gap="2" mb="6">
        <Button
          flex="1"
          variant={mode === 'create' ? 'solid' : 'outline'}
          onClick={() => handleModeSwitch('create')}
          disabled={isSubmitting}
        >
          Create
        </Button>
        <Button
          flex="1"
          variant={mode === 'from-reference' ? 'solid' : 'outline'}
          onClick={() => handleModeSwitch('from-reference')}
          disabled={isSubmitting}
        >
          From Reference
        </Button>
      </Flex>

      {/* Reference image section — only visible in "From Reference" mode */}
      {mode === 'from-reference' && (
        <Box mb="4" p="4" borderWidth="1px" borderColor="amber" borderRadius="sm" bg="parchment">
          <Text fontSize="sm" fontWeight="medium" color="ink" mb="3">
            Reference Image
          </Text>

          {/* Show selected reference preview or upload controls */}
          {referenceImage ? (
            <Box>
              <Box
                borderWidth="1px"
                borderColor="amber"
                borderRadius="sm"
                overflow="hidden"
                mb="2"
                maxH="160px"
              >
                <Image
                  src={`data:image/png;base64,${referenceImage}`}
                  alt="Reference image preview"
                  width="100%"
                  maxH="160px"
                  objectFit="contain"
                />
              </Box>
              {selectedFileName && (
                <Text fontSize="xs" color="warmgray" mb="2" fontFamily="mono">
                  {selectedFileName}
                </Text>
              )}
              <Button
                size="sm"
                variant="outline"
                onClick={handleClearReference}
                disabled={isSubmitting}
              >
                Clear Reference
              </Button>
            </Box>
          ) : (
            <Box>
              {/* File upload */}
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/webp"
                onChange={handleFileSelect}
                style={{ display: 'none' }}
                aria-label="Upload reference image"
              />
              <Button
                size="sm"
                variant="outline"
                onClick={() => fileInputRef.current?.click()}
                disabled={isSubmitting || isLoadingReference}
                mb="3"
              >
                {isLoadingReference ? 'Loading...' : 'Upload Image'}
              </Button>
              <Text fontSize="xs" color="warmgray" mb="3">
                JPEG, PNG, or WEBP — max 5MB
              </Text>

              {/* Gallery image picker */}
              {completedImages.length > 0 && (
                <Box>
                  <Text fontSize="xs" fontWeight="medium" color="ink" mb="2">
                    Or pick from gallery
                  </Text>
                  <Flex gap="2" flexWrap="wrap">
                    {completedImages.map((img) => (
                      <Box
                        key={img.id}
                        as="button"
                        width="60px"
                        height="60px"
                        borderWidth="1px"
                        borderColor="amber"
                        borderRadius="sm"
                        overflow="hidden"
                        cursor="pointer"
                        opacity={isSubmitting || isLoadingReference ? 0.5 : 1}
                        _hover={{ borderWidth: '2px' }}
                        onClick={() => {
                          if (!isSubmitting && !isLoadingReference && img.fileUrl) {
                            handleGalleryImageSelect(img.fileUrl)
                          }
                        }}
                        disabled={isSubmitting || isLoadingReference}
                        aria-label={`Use image "${img.prompt}" as reference`}
                      >
                        <Image
                          src={`${window.__VISCRAFT_CONFIG__?.API_BASE_URL || 'http://localhost:8080'}${img.fileUrl}`}
                          alt={img.prompt}
                          width="100%"
                          height="100%"
                          objectFit="cover"
                        />
                      </Box>
                    ))}
                  </Flex>
                </Box>
              )}
            </Box>
          )}

          {/* Reference image error */}
          {errors.referenceImage && (
            <Text color="oxblood" fontSize="xs" mt="2">
              {errors.referenceImage}
            </Text>
          )}
        </Box>
      )}

      {/* Rate limit banner (ERR_02) */}
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

      {/* Prompt textarea */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Description
        </Text>
        <Textarea
          placeholder="Describe the concept art you want to generate..."
          value={prompt}
          onChange={(e) => {
            setPrompt(e.target.value)
            if (errors.prompt) {
              setErrors((prev) => {
                const next = { ...prev }
                delete next.prompt
                return next
              })
            }
          }}
          fontFamily="mono"
          fontSize="sm"
          rows={4}
          resize="vertical"
          borderColor={errors.prompt ? 'oxblood' : undefined}
          disabled={isSubmitting}
        />
        <Flex justify="space-between" mt="1">
          {errors.prompt ? (
            <Text color="oxblood" fontSize="xs">
              {errors.prompt}
            </Text>
          ) : (
            <Box />
          )}
          <Text
            fontSize="xs"
            color={prompt.length > MAX_PROMPT_LENGTH ? 'oxblood' : 'warmgray'}
            fontFamily="mono"
          >
            {prompt.length}/{MAX_PROMPT_LENGTH}
          </Text>
        </Flex>
      </Box>

      {/* Genre select */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Genre
        </Text>
        <Box
          as="select"
          value={genre}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
            setGenre(e.target.value as Genre | '')
            if (errors.genre) {
              setErrors((prev) => {
                const next = { ...prev }
                delete next.genre
                return next
              })
            }
          }}
          width="100%"
          height="40px"
          px="3"
          fontFamily="body"
          fontSize="sm"
          bg="parchment"
          color="ink"
          borderWidth="1px"
          borderColor={errors.genre ? 'oxblood' : 'amber'}
          borderRadius="sm"
          disabled={isSubmitting}
        >
          <option value="">Select genre...</option>
          {GENRE_OPTIONS.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </Box>
        {errors.genre && (
          <Text color="oxblood" fontSize="xs" mt="1">
            {errors.genre}
          </Text>
        )}
      </Box>

      {/* Asset Type select */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Asset Type
        </Text>
        <Box
          as="select"
          value={assetType}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
            setAssetType(e.target.value as AssetType | '')
            if (errors.assetType) {
              setErrors((prev) => {
                const next = { ...prev }
                delete next.assetType
                return next
              })
            }
          }}
          width="100%"
          height="40px"
          px="3"
          fontFamily="body"
          fontSize="sm"
          bg="parchment"
          color="ink"
          borderWidth="1px"
          borderColor={errors.assetType ? 'oxblood' : 'amber'}
          borderRadius="sm"
          disabled={isSubmitting}
        >
          <option value="">Select asset type...</option>
          {ASSET_TYPE_OPTIONS.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </Box>
        {errors.assetType && (
          <Text color="oxblood" fontSize="xs" mt="1">
            {errors.assetType}
          </Text>
        )}
      </Box>

      {/* Mood select */}
      <Box mb="4">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Mood
        </Text>
        <Box
          as="select"
          value={mood}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
            setMood(e.target.value as Mood | '')
            if (errors.mood) {
              setErrors((prev) => {
                const next = { ...prev }
                delete next.mood
                return next
              })
            }
          }}
          width="100%"
          height="40px"
          px="3"
          fontFamily="body"
          fontSize="sm"
          bg="parchment"
          color="ink"
          borderWidth="1px"
          borderColor={errors.mood ? 'oxblood' : 'amber'}
          borderRadius="sm"
          disabled={isSubmitting}
        >
          <option value="">Select mood...</option>
          {MOOD_OPTIONS.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </Box>
        {errors.mood && (
          <Text color="oxblood" fontSize="xs" mt="1">
            {errors.mood}
          </Text>
        )}
      </Box>

      {/* Submit button */}
      <Button
        width="full"
        variant="solid"
        mt="2"
        disabled={isSubmitting}
        onClick={handleSubmit}
      >
        {isSubmitting ? 'Generating...' : 'Generate'}
      </Button>
    </ReusableModal>
  )
}

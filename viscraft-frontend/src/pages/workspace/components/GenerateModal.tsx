import { useRef } from 'react'
import { Box, Button, Text, Textarea, HStack, Image, Flex } from '@chakra-ui/react'
import { ReusableModal } from '../../../components/ReusableModal'
import { useGenerateAdShot } from '../hooks/useGenerateAdShot'

/**
 * Generate Modal for product ad photography.
 * Purely presentational — all logic lives in useGenerateAdShot hook.
 */

interface GenerateModalProps {
  isOpen: boolean
  onClose: () => void
}

export function GenerateModal({ isOpen, onClose }: GenerateModalProps) {
  const {
    userPrompt,
    handleUserPromptChange,
    toggleOption,
    isSelected,
    categories,
    generatedPrompt,
    isPromptValid,
    isSubmitting,
    apiError,
    handleSubmit,
    resetForm,
    isRegenerate,
    uploadedReferenceImage,
    handleUploadReference,
    handleFileSelect,
    regenerateFileUrl,
  } = useGenerateAdShot(isOpen)

  const fileInputRef = useRef<HTMLInputElement>(null)

  function handleClose() {
    resetForm()
    onClose()
  }

  return (
    <ReusableModal
      isOpen={isOpen}
      onClose={handleClose}
      title={isRegenerate ? 'Regenerate Ad Shot' : 'Generate Ad Shot'}
      size="lg"
    >
      {apiError && (
        <Box bg="red.50" border="1px solid" borderColor="oxblood" borderRadius="sm" p="3" mb="4">
          <Text color="oxblood" fontSize="sm" fontFamily="mono">{apiError}</Text>
        </Box>
      )}

      {/* Product description */}
      <Box mb="4" data-tour="generate-form">
        <Text as="label" fontSize="sm" fontWeight="medium" color="ink" mb="1" display="block">
          Product Description *
        </Text>
        <Textarea
          placeholder="e.g. A bottle of cold-pressed orange juice with natural pulp..."
          value={userPrompt}
          onChange={(e) => handleUserPromptChange(e.target.value)}
          fontFamily="mono"
          fontSize="sm"
          rows={3}
          resize="vertical"
          disabled={isSubmitting}
        />
        <Text fontSize="xs" color="warmgray" mt="1" fontFamily="mono">
          {userPrompt.length}/300
        </Text>
      </Box>

      {/* Prompt option categories */}
      <Box data-tour="prompt-options">
      {categories.map((cat) => (
        <Box key={cat.category} mb="3">
          <Text fontSize="xs" fontWeight="medium" color="ink" mb="1.5">
            {cat.label} {cat.multi && <Text as="span" color="warmgray">(multiple)</Text>}
          </Text>
          <HStack gap="2" flexWrap="wrap">
            {cat.options.map((opt) => (
              <Button
                key={opt.id}
                size="xs"
                minH="28px"
                px="3"
                variant={isSelected(opt.id) ? 'solid' : 'outline'}
                onClick={() => toggleOption(opt, cat.multi)}
                disabled={isSubmitting}
                fontSize="2xs"
              >
                {opt.label}
              </Button>
            ))}
          </HStack>
        </Box>
      ))}
      </Box>

      {/* Side-by-side: Generated Prompt Preview + Reference Photo */}
      <Flex gap="4" mt="4" direction={{ base: 'column', md: 'row' }}>
        {/* Left: Generated prompt preview */}
        <Box flex="1" p="3" bg="rgba(201,118,44,0.05)" borderRadius="sm" border="1px solid" borderColor="rgba(201,118,44,0.2)" data-tour="prompt-preview">
          <Text fontSize="xs" fontWeight="medium" color="warmgray" mb="1">
            Generated Prompt
          </Text>
          <Text fontFamily="mono" fontSize="xs" color="ink" lineHeight="tall">
            {generatedPrompt || 'Start typing to see the prompt preview...'}
          </Text>
        </Box>

        {/* Right: Reference photo */}
        <Box flex="1" p="3" borderWidth="1px" borderColor="amber" borderRadius="sm">
          <Text fontSize="xs" fontWeight="medium" color="ink" mb="2">
            Reference Photo {isRegenerate ? '(from original)' : '(optional)'}
          </Text>

          {/* Regenerate source image */}
          {isRegenerate && regenerateFileUrl && !uploadedReferenceImage && (
            <Box mb="2">
              <Box borderWidth="1px" borderColor="amber" borderRadius="sm" overflow="hidden" maxH="100px">
                <Image src={regenerateFileUrl} alt="Original reference" width="100%" maxH="100px" objectFit="contain" />
              </Box>
            </Box>
          )}

          {/* Uploaded reference */}
          {uploadedReferenceImage && (
            <Box mb="2">
              <Box borderWidth="1px" borderColor="amber" borderRadius="sm" overflow="hidden" maxH="100px">
                <Image src={`data:image/png;base64,${uploadedReferenceImage}`} alt="Uploaded reference" width="100%" maxH="100px" objectFit="contain" />
              </Box>
              <Button size="xs" variant="ghost" color="oxblood" mt="1" onClick={() => handleUploadReference(null)} disabled={isSubmitting}>
                Remove
              </Button>
            </Box>
          )}

          {/* Upload button */}
          {!uploadedReferenceImage && (
            <Box>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/webp"
                onChange={(e) => { handleFileSelect(e.target.files?.[0] ?? null); if (fileInputRef.current) fileInputRef.current.value = '' }}
                style={{ display: 'none' }}
              />
              <Button size="xs" variant="outline" onClick={() => fileInputRef.current?.click()} disabled={isSubmitting}>
                {isRegenerate ? 'Upload Different' : 'Choose File'}
              </Button>
              <Text fontSize="2xs" color="warmgray" mt="1">JPEG, PNG, WebP — max 5MB</Text>
            </Box>
          )}
        </Box>
      </Flex>

      {/* Submit */}
      <Button
        width="full"
        variant="solid"
        mt="4"
        minH="44px"
        disabled={isSubmitting || !isPromptValid}
        onClick={handleSubmit}
      >
        {isSubmitting ? 'Generating...' : 'Generate'}
      </Button>
    </ReusableModal>
  )
}
